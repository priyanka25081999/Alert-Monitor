package types

import "time"

type Event struct {
	Client    string
	EventType string
	Timestamp time.Time
}
