package entity

import (
	"time"
)

// IncidentStatus contain the information of each incident status
type IncidentStatus struct {
	Code    string
	Message string
	Color   string
}

// IsRecovered check whether this incident is marked as recovered or not
func (s IncidentStatus) IsRecovered() bool {
	return s == StatusRecovered
}

var (
	//StatusWarning is a status when the vendor mark an alert as a WARNING
	StatusWarning = IncidentStatus{Code: "warning", Message: "Warning", Color: "FF9C00"}

	//StatusTriggered is a status when the vendor mark an alert as an ERROR
	StatusTriggered = IncidentStatus{Code: "error", Message: "Triggered", Color: "FF0000"}

	//StatusRecovered is a status when the vendor mark an alert as RECOVER
	StatusRecovered = IncidentStatus{Code: "recover", Message: "Recovered", Color: "00BF85"}
)

// Incident contain information of incident got from vendor data
type Incident struct {
	Title      string
	ThreadID   string
	Vendor     string
	Status     IncidentStatus
	LastUpdate time.Time
}
