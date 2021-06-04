package entity

import (
	"testing"
)

func TestGetColor(t *testing.T) {
	empty := &Notification{}
	if empty.GetColor() != defaultColor {
		t.Errorf("Empty Notification expecting %s, got %s instead", defaultColor, empty.GetColor())
	}

	color := "FF0000"
	filled := &Notification{Color: color}
	if filled.GetColor() != color {
		t.Errorf("Notification expecting %s, got %s instead", color, filled.GetColor())
	}
}
