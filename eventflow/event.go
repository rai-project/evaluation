package eventflow

import "time"

type Event struct {
	ID        string    `json:"EVENT_ID"`
	Name      time.Time `json:"EVENT_NAME"`
	MetaData  string    `json:"META"`
	TimeStamp time.Time `json:"TS"`
}
