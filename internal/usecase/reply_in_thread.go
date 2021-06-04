package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/alvintzz/alert-thread/internal/entity"

	log "github.com/sirupsen/logrus"
)

// ReplyInThread will check whether the incident is already notified, create a new thread for new incident and reply in the thread for existing incident
func (u *Usecase) ReplyInThread(ctx context.Context, param entity.ReplyInThread) error {
	incident, err := u.storage.GetIncident(ctx, param.GetKey())
	if err != nil {
		log.Errorf("Failed to get incident from storage because %s", err)
		return err
	}

	threadID := incident.ThreadID
	if threadID == "" {
		// Sending Main Thread
		threadID, err = u.sendMessage(ctx, "", param)
		if err != nil {
			log.Error(err)
			return err
		}

		incident = entity.Incident{
			Title:      param.GetTitle(),
			ThreadID:   threadID,
			Vendor:     param.GetVendor(),
			Status:     param.GetStatus(),
			LastUpdate: time.Now(),
		}
	} else {
		// Update Main Thread
		err = u.updateMessage(ctx, threadID, incident, param)
		if err != nil {
			log.Error(err)
			return err
		}

		incident = entity.Incident{
			Status:     param.GetStatus(),
			LastUpdate: time.Now(),
		}
	}

	// Register Incident to Storage
	err = u.storage.RegisterIncident(ctx, param.GetKey(), incident)
	if err != nil {
		log.Errorf("Failed to register incident into storage because %s", err)
		return fmt.Errorf("Failed to register incident into storage because %s", err)
	}

	// Sending Thread
	_, err = u.sendMessage(ctx, threadID, param)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (u *Usecase) sendMessage(ctx context.Context, threadID string, param entity.ReplyInThread) (string, error) {
	message := param.GetSummary()
	image := ""
	if threadID != "" {
		message = param.GetDetail()
		image = param.GetImage()
	}

	threadID, err := u.notification.SendMessage(ctx, entity.Notification{
		Channel: param.GetChannel(),
		Title:   param.GetTitle(),
		Message: message,
		Color:   param.GetStatus().Color,
		Image:   image,
		Metadata: map[string]string{
			"timestamp": threadID,
		},
	})
	if err != nil {
		return "", fmt.Errorf("Failed to send message because %s", err)
	}

	return threadID, nil
}

func (u *Usecase) updateMessage(ctx context.Context, threadID string, incident entity.Incident, param entity.ReplyInThread) error {
	err := u.notification.UpdateMessage(ctx, entity.Notification{
		Channel: param.GetChannel(),
		Title:   incident.Title,
		Message: param.GetSummary(),
		Color:   param.GetStatus().Color,
		Image:   param.GetImage(),
		Metadata: map[string]string{
			"timestamp": threadID,
		},
	})
	if err != nil {
		return fmt.Errorf("Failed to send message because %s", err)
	}

	return nil
}
