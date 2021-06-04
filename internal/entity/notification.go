package entity

// Notification contains information of what we want to send to notification channel
type Notification struct {
	Channel  string
	Title    string
	Message  string
	Color    string
	Image    string
	Metadata map[string]string
}

const defaultColor = "FFFFFF"

// GetColor will return default color white if the color is not specified
func (n Notification) GetColor() string {
	color := defaultColor
	if n.Color != "" {
		color = n.Color
	}

	return color
}
