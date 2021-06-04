package entity

import (
	"testing"
)

func TestIsRecovered(t *testing.T) {
	if StatusWarning.IsRecovered() {
		t.Errorf("status %s got %t instead", StatusWarning.Message, StatusWarning.IsRecovered())
	}
	if StatusTriggered.IsRecovered() {
		t.Errorf("status %s got %t instead", StatusWarning.Message, StatusWarning.IsRecovered())
	}
	if !StatusRecovered.IsRecovered() {
		t.Errorf("status %s got %t instead", StatusWarning.Message, StatusWarning.IsRecovered())
	}
}
