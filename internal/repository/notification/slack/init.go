package slack

import (
	sl "github.com/slack-go/slack"
)

// Slack contains dependencies needed by Slack integration to send notification
type Slack struct {
	client *sl.Client
}

// NewNotification will return slack object used to do slack integration
func NewNotification(token string) (*Slack, error) {
	return &Slack{
		client: sl.New(token),
	}, nil
}
