package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/alvintzz/alert-thread/internal/entity"

	log "github.com/sirupsen/logrus"
)

// DatadogReplyThread is object to cater Datadog request parameter
type DatadogReplyThread struct {
	ID        string `json:"id"`
	Body      string `json:"body"`
	Title     string `json:"title"`
	Key       string `json:"cycle_key"`
	StatusStr string `json:"alert_status"`
	Status    string `json:"alert_type"`
	AlertID   string `json:"alert_id"`
	Snapshot  string `json:"snapshot"`
	Vendor    string
	Channel   string `json:"channel"`
	Tags      string `json:"tags"`
}

// GetKey will return the incident unique id as identifier whether the incident should go to same thread or not
func (r DatadogReplyThread) GetKey() string {
	return r.Key
}

// GetChannel will return specified Notification Channel ID from Datadog Object
func (r DatadogReplyThread) GetChannel() string {
	return r.Channel
}

// GetVendor will return from which vendor this parameter is which is always Datadog in this case
func (r DatadogReplyThread) GetVendor() string {
	return r.Vendor
}

// GetTitle will return title string for notification message. In Slack Notification, this title will be visible in the slack notification pop-up
func (r DatadogReplyThread) GetTitle() string {
	if strings.Contains(r.Title, ">=") {
		r.Title = strings.Replace(r.Title, ">=", "more/equal than", -1)
	}
	if strings.Contains(r.Title, ">") {
		r.Title = strings.Replace(r.Title, ">", "more than", -1)
	}
	if strings.Contains(r.Title, "<=") {
		r.Title = strings.Replace(r.Title, "<=", "less/equal than", -1)
	}
	if strings.Contains(r.Title, "<") {
		r.Title = strings.Replace(r.Title, "<", "less than", -1)
	}

	return fmt.Sprintf("*<%s|%s>*", r.GetURL(), r.Title)
}

// GetStatus will return status of the incident
func (r DatadogReplyThread) GetStatus() entity.IncidentStatus {
	if r.Status == "success" {
		return entity.StatusRecovered
	} else if r.Status == "warning" {
		return entity.StatusWarning
	}

	return entity.StatusTriggered
}

// GetSummary will return summary string for notification header. This message will be shown in the main thread and should contain at-glance summary of incident
func (r DatadogReplyThread) GetSummary() string {
	hangoutLink := fmt.Sprintf("http://g.co/meet/tkpd-%s", r.GetKey())
	summary := fmt.Sprintf("*Hangout Link* : %s\n\nFrom: *%s*\nCurrent Status: *%s*", hangoutLink, r.GetVendor(), r.GetStatus().Message)

	return summary
}

// GetURL is helper function to get datadog monitor URL
func (r DatadogReplyThread) GetURL() string {
	url := fmt.Sprintf("https://app.datadoghq.com/monitors#%s", r.AlertID)

	arr := strings.Split(r.Body, "Metric Graph:")
	if len(arr) > 1 {
		arr = strings.Split(arr[1], "\u00b7")
		if len(arr) > 1 {
			url = strings.TrimSpace(arr[0])
		}
	}

	return url
}

// GetDetail will return detail string for notification thread. This message will be shown in the threads and can contain more detail incident data
func (r DatadogReplyThread) GetDetail() string {
	body := r.Body

	arr := strings.Split(body, "Metric Graph:")
	if len(arr) > 1 {
		body = arr[0]
	}

	index := strings.LastIndex(body, "`")
	if index > 0 {
		body = body[:index]
	}
	index = strings.LastIndex(body, "`")
	if index > 0 {
		body = body[:index]
	}

	return body
}

// GetImage will return the image url of metrics snapshot. For Datadog, we put delay as Datadog need time to upload the image
func (r DatadogReplyThread) GetImage() string {
	if r.Snapshot == "" || r.Snapshot == "null" {
		return ""
	}
	time.Sleep(2 * time.Second)
	return r.Snapshot
}

// DdogReplyInThread receive callback request from Datadog and parse the request into ReplyInThread interface and call the usecase function
func (s *Handler) DdogReplyInThread(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to read body because %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request := DatadogReplyThread{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Errorf("Failed to unmarshal body because %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	request.Vendor = "Datadog"

	go s.usecase.ReplyInThread(context.Background(), request)

	w.WriteHeader(http.StatusOK)
}
