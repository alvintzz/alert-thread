package usecase

import (
	"context"

	"github.com/alvintzz/alert-thread/internal/entity"
)

// Storage is interface of storage use to save incidents
type Storage interface {
	GetIncident(ctx context.Context, key string) (entity.Incident, error)
	RegisterIncident(ctx context.Context, key string, incident entity.Incident) error
	RemoveIncident(ctx context.Context, key string) error
}

// Notification is interface of notification channel used to notify an incident
type Notification interface {
	SendMessage(ctx context.Context, param entity.Notification) (string, error)
	UpdateMessage(ctx context.Context, param entity.Notification) error
}

// Usecase contains all dependencies for slack-alert flow
type Usecase struct {
	storage      Storage
	notification Notification
}

// New will return object contain the usecases served by this service
func New(store Storage, notif Notification) *Usecase {
	return &Usecase{
		storage:      store,
		notification: notif,
	}
}
