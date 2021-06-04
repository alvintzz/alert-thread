package usecase

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/alvintzz/alert-thread/internal/entity"
)

var errorDefault = fmt.Errorf("timeout")

const (
	getErrorKey   = "mock_error"
	getEmptyKey   = "mock_empty"
	getSuccessKey = "mock_success"
)

var successIncident = entity.Incident{
	Title:      "Success Mock",
	ThreadID:   "success_channel",
	Vendor:     "datadog",
	Status:     entity.StatusWarning,
	LastUpdate: time.Now(),
}

type storageMock struct{}

func (s *storageMock) GetIncident(ctx context.Context, key string) (entity.Incident, error) {
	if strings.HasPrefix(key, "mock_success") {
		return successIncident, nil
	} else if strings.HasPrefix(key, "mock_empty") {
		return entity.Incident{}, nil
	}
	return entity.Incident{}, errorDefault
}
func (s *storageMock) RegisterIncident(ctx context.Context, key string, incident entity.Incident) error {
	if strings.HasSuffix(key, "register_success") {
		return nil
	}
	return errorDefault
}
func (s *storageMock) RemoveIncident(ctx context.Context, key string) error {
	return nil
}

type notificationMock struct{}

func (n *notificationMock) SendMessage(ctx context.Context, param entity.Notification) (string, error) {
	if strings.HasPrefix(param.Channel, "success_channel") {
		return "success_channel", nil
	}
	return "", errorDefault
}
func (n *notificationMock) UpdateMessage(ctx context.Context, param entity.Notification) error {
	if strings.HasSuffix(param.Channel, "update_success") {
		return nil
	}

	return errorDefault
}

type Parameter struct {
	get      string
	register bool
	send     bool
	update   bool
	thread   bool
	count    int
}

func (p *Parameter) GetVendor() string {
	return "datadog"
}
func (p *Parameter) GetKey() string {
	str := p.get
	if p.register {
		str += "_register_success"
	}
	return str
}
func (p *Parameter) GetTitle() string {
	return "title"
}
func (p *Parameter) GetSummary() string {
	return "summary"
}
func (p *Parameter) GetDetail() string {
	return "detail"
}
func (p *Parameter) GetStatus() entity.IncidentStatus {
	return entity.StatusWarning
}
func (p *Parameter) GetImage() string {
	return "http://image.com"
}
func (p *Parameter) GetChannel() string {
	str := "failed_channel"
	if p.send {
		str = "success_channel"
	}
	if p.update {
		str += "_update_success"
	}
	return str
}

func TestReplyInThread(t *testing.T) {
	ctx := context.Background()
	uc := New(&storageMock{}, &notificationMock{})

	usecase := []Parameter{
		Parameter{get: getErrorKey, register: false, send: false},
		Parameter{get: getEmptyKey, register: false, send: false},
		Parameter{get: getEmptyKey, register: true, send: false},
		Parameter{get: getEmptyKey, register: false, send: true},
		Parameter{get: getEmptyKey, register: true, send: true},
		Parameter{get: getSuccessKey, register: false, send: false},
		Parameter{get: getSuccessKey, register: true, send: false},
		Parameter{get: getSuccessKey, register: false, send: true},
		Parameter{get: getSuccessKey, register: true, send: true, update: false},
		Parameter{get: getSuccessKey, register: true, send: true, update: true},
	}

	expected := []error{
		errorDefault,
		errorDefault,
		errorDefault,
		errorDefault,
		nil,
		errorDefault,
		errorDefault,
		errorDefault,
		errorDefault,
		nil,
	}

	for k, v := range usecase {
		err := uc.ReplyInThread(ctx, &v)
		if err != nil && expected[k] != nil {
			//Expected
			continue
		}
		if err != expected[k] {
			t.Errorf("Reply for key %s [reg:%t][send:%t][update:%t][thread:%t] expecting %+v but got %+v", v.get, v.register, v.send, v.update, v.thread, expected[k], err)
		}
	}
}
