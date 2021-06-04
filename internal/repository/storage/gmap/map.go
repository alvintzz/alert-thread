package gmap

import (
	"context"

	"github.com/alvintzz/alert-thread/internal/entity"
)

// GetIncident will return object of Incident saved inside the chosen storage
func (m *Storage) GetIncident(ctx context.Context, key string) (entity.Incident, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	value, ok := m.Incidents[key]
	if ok {
		return value, nil
	}

	return entity.Incident{}, nil
}

// RegisterIncident will register or update object of Incident saved inside the chosen storage
func (m *Storage) RegisterIncident(ctx context.Context, key string, incident entity.Incident) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	m.Incidents[key] = incident

	return nil
}

// RemoveIncident will remove object of Incident saved inside the chosen storage
func (m *Storage) RemoveIncident(ctx context.Context, key string) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	delete(m.Incidents, key)

	return nil
}
