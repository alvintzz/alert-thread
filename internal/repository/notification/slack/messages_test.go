package slack

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/alvintzz/alert-thread/internal/entity"
	"gopkg.in/h2non/gock.v1"
)

var slackURL = "https://slack.com"
var slackToken = "this_is_token"
var sendMessageEndpoint = "/api/chat.postMessage"
var updateMessageEndpoint = "/api/chat.update"

var expectedThreadID = "1503435956.000247"
var responseSuccess = fmt.Sprintf(`{"ok": true,"channel": "C1H9RESGL","ts": "%s","message": {"text": "Here's a message for you","username": "ecto1","bot_id": "B19LU7CSY","attachments": [{"text": "This is an attachment","id": 1,"fallback": "This is an attachment's fallback"}],"type": "message","subtype": "bot_message","ts": "%s"}}`, expectedThreadID, expectedThreadID)
var responseFailed = fmt.Sprintf(`{"ok":false,"error":"%s"}`, errorNotAuth)

var errorTimeout = fmt.Errorf("timeout exceeded")
var errorNotAuth = fmt.Errorf("not_authed")

var messageFailedObj = "Failed to create notification object: %s"
var messageNotError = "Failed in %s for %s. Not expecting error %s"
var messageNotExpect = "Failed in %s for %s. Result is not matched expected %+v, got %+v"

var flowSendMessage = "send message flow"
var flowUpdateMessage = "update message flow"

func createMatcher(params map[string]string) func(req *http.Request, ereq *gock.Request) (bool, error) {
	return func(req *http.Request, ereq *gock.Request) (bool, error) {
		match := true
		for key, value := range params {
			if req.PostFormValue(key) == value {
				continue
			}

			match = false
			break
		}
		return match, nil
	}
}

func gockSendMessage() {
	gock.New(slackURL).
		Post(sendMessageEndpoint).
		AddMatcher(createMatcher(map[string]string{"channel": "channel_1"})).
		Reply(200).
		BodyString(responseSuccess)

	gock.New(slackURL).
		Post(sendMessageEndpoint).
		AddMatcher(createMatcher(map[string]string{"channel": "channel_2"})).
		ReplyError(errorTimeout)

	gock.New(slackURL).
		Post(sendMessageEndpoint).
		AddMatcher(createMatcher(map[string]string{"channel": "channel_3"})).
		Reply(200).
		BodyString(responseFailed)
}

func TestSendMessage(t *testing.T) {
	ctx := context.Background()
	obj, err := NewNotification(slackToken)
	if err != nil {
		t.Errorf(messageFailedObj, err)
	}

	gockSendMessage()
	defer gock.Off()

	threadID, err := obj.SendMessage(ctx, entity.Notification{Channel: "channel_1", Image: "http://xxx.com"})
	if err != nil {
		t.Errorf(messageNotError, flowSendMessage, "channel_1", err)
	} else if threadID != expectedThreadID {
		t.Errorf(messageNotExpect, flowSendMessage, "channel_1", expectedThreadID, threadID)
	}

	_, err = obj.SendMessage(ctx, entity.Notification{Channel: "channel_2", Metadata: map[string]string{"timestamp": "XXX"}})
	if !strings.Contains(err.Error(), errorTimeout.Error()) {
		t.Errorf(messageNotExpect, flowSendMessage, "channel_2", errorTimeout, err)
	}

	_, err = obj.SendMessage(ctx, entity.Notification{Channel: "channel_3"})
	if !strings.Contains(err.Error(), errorNotAuth.Error()) {
		t.Errorf(messageNotExpect, flowSendMessage, "channel_3", errorNotAuth, err)
	}
}

func gockUpdateMessage() {
	gock.New(slackURL).
		Post(updateMessageEndpoint).
		AddMatcher(createMatcher(map[string]string{"channel": "channel_1"})).
		Reply(200).
		BodyString(responseSuccess)

	gock.New(slackURL).
		Post(updateMessageEndpoint).
		AddMatcher(createMatcher(map[string]string{"channel": "channel_2"})).
		ReplyError(errorTimeout)

	gock.New(slackURL).
		Post(updateMessageEndpoint).
		AddMatcher(createMatcher(map[string]string{"channel": "channel_3"})).
		Reply(200).
		BodyString(responseFailed)
}

func TestUpdateMessage(t *testing.T) {
	ctx := context.Background()
	obj, err := NewNotification(slackToken)
	if err != nil {
		t.Errorf(messageFailedObj, err)
	}

	gockUpdateMessage()
	defer gock.Off()

	err = obj.UpdateMessage(ctx, entity.Notification{Channel: "channel_1"})
	if err == nil {
		t.Errorf(messageNotExpect, flowUpdateMessage, "channel_1", "nil", err)
	} else if err != errorEmptyThreadID {
		t.Errorf(messageNotExpect, flowUpdateMessage, "channel_1", errorEmptyThreadID, err)
	}

	err = obj.UpdateMessage(ctx, entity.Notification{Channel: "channel_1", Metadata: map[string]string{"timestamp": "XXX"}, Image: "http://xxx.com"})
	if err != nil {
		t.Errorf(messageNotError, flowUpdateMessage, "channel_1", err)
	}

	err = obj.UpdateMessage(ctx, entity.Notification{Channel: "channel_2", Metadata: map[string]string{"timestamp": "XXX"}})
	if err == nil {
		t.Errorf(messageNotExpect, flowSendMessage, "channel_2", "nil", err)
	} else if !strings.Contains(err.Error(), errorTimeout.Error()) {
		t.Errorf(messageNotExpect, flowUpdateMessage, "channel_2", errorTimeout, err)
	}

	err = obj.UpdateMessage(ctx, entity.Notification{Channel: "channel_3", Metadata: map[string]string{"timestamp": "XXX"}})
	if err == nil {
		t.Errorf(messageNotExpect, flowSendMessage, "channel_3", "nil", err)
	} else if !strings.Contains(err.Error(), errorNotAuth.Error()) {
		t.Errorf(messageNotExpect, flowUpdateMessage, "channel_3", errorNotAuth, err)
	}
}
