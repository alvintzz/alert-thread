package gmap

import (
	"context"
	"testing"

	"github.com/alvintzz/alert-thread/internal/entity"
)

var messageFailedObj = "Failed to create notification object: %s"
var messageNotError = "Failed in %s for %s. Not expecting error %s"
var messageNotExpect = "Failed in %s for %s. Result is not matched expected \"%+v\", got \"%+v\""

var flowGetIncident = "get incident flow"
var flowRegisterIncident = "register incident flow"
var flowRemoveIncident = "remove incident flow"

func TestGetIncident(t *testing.T) {
	storage, _ := NewStorage()
	storage.Incidents["incident_1"] = entity.Incident{Title: "Incident 1"}

	ctx := context.Background()
	incident, err := storage.GetIncident(ctx, "incident_1")
	if err != nil {
		t.Errorf(messageNotError, flowGetIncident, "incident_1", err)
	} else if incident.Title != "Incident 1" {
		t.Errorf(messageNotExpect, flowGetIncident, "incident_1", "Incident 1", incident.Title)
	}

	incident, _ = storage.GetIncident(ctx, "incident_2")
	if err != nil {
		t.Errorf(messageNotError, flowGetIncident, "incident_2", err)
	} else if incident.Title != "" {
		t.Errorf(messageNotExpect, flowGetIncident, "incident_2", "", incident.Title)
	}
}

func TestRegisterIncident(t *testing.T) {
	storage, _ := NewStorage()
	storage.Incidents["incident_1"] = entity.Incident{Title: "Incident 1"}

	ctx := context.Background()
	err := storage.RegisterIncident(ctx, "incident_2", entity.Incident{Title: "Incident 2"})
	if err != nil {
		t.Errorf(messageNotError, flowRegisterIncident, "incident_2", err)
	} else if value, ok := storage.Incidents["incident_2"]; !(ok && value.Title == "Incident 2") {
		t.Errorf(messageNotExpect, flowGetIncident, "incident_2", "Incident 2", value.Title)
	}

	err = storage.RegisterIncident(ctx, "incident_1", entity.Incident{Title: "Incident 1 New"})
	if err != nil {
		t.Errorf(messageNotError, flowRegisterIncident, "incident_1", err)
	} else if value, ok := storage.Incidents["incident_1"]; !(ok && value.Title == "Incident 1 New") {
		t.Errorf(messageNotExpect, flowGetIncident, "incident_1", "Incident 1 New", value.Title)
	}
}

func TestRemoveIncident(t *testing.T) {
	storage, _ := NewStorage()
	storage.Incidents["incident_1"] = entity.Incident{Title: "Incident 1"}

	ctx := context.Background()
	err := storage.RemoveIncident(ctx, "incident_1")
	if err != nil {
		t.Errorf(messageNotError, flowRemoveIncident, "incident_1", err)
	} else if value, ok := storage.Incidents["incident_1"]; ok {
		t.Errorf(messageNotExpect, flowRemoveIncident, "incident_1", "", value.Title)
	}

	err = storage.RemoveIncident(ctx, "incident_1")
	if err != nil {
		t.Errorf(messageNotError, flowRemoveIncident, "incident_1", err)
	}
}
