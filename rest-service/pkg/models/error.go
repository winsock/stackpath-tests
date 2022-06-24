package models

import "time"

type Error struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
