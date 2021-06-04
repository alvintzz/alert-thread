package handler

import (
	"context"

	"github.com/alvintzz/alert-thread/internal/entity"
)

// Usecase is interface of slack-alert flow to send notification in thread
type Usecase interface {
	ReplyInThread(ctx context.Context, param entity.ReplyInThread) error
}

// Handler contains all dependencies for handler endpoint
type Handler struct {
	usecase Usecase
}

// New will return object contain the handlers served by this service
func New(usecase Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}
