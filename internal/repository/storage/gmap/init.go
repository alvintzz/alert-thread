package gmap

import (
	"sync"

	"github.com/alvintzz/alert-thread/internal/entity"
)

// Storage is object Storage using Golang's Map
type Storage struct {
	Mutex     sync.Mutex
	Incidents map[string]entity.Incident
}

// NewStorage will return storage implementation using Golang's Map
func NewStorage() (*Storage, error) {
	return &Storage{
		Incidents: map[string]entity.Incident{},
	}, nil
}
