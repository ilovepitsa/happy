package entity

import "time"

type Notification struct {
	Target           int32
	NotificationText string
	Date             time.Time
}
