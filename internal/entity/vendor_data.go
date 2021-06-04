package entity

// ReplyInThread is an interface of vendor data for slack-alert flow purpose
type ReplyInThread interface {
	GetVendor() string

	GetKey() string

	GetTitle() string
	GetSummary() string
	GetDetail() string
	GetStatus() IncidentStatus
	GetImage() string
	GetChannel() string
}
