package models

import "time"

type Trigger struct {
	Id            int       `json:"id"`
	SystemId      string    `json:"system_id"`
	Url           string    `json:"url"`
	LastFailedAt  *int64    `json:"last_failed_at"`
	BufferSeconds int       `json:"buffer_seconds"`
	Status        string    `json:"status"` // ok, error, triggered
	LastCheckedAt *int64    `json:"last_checked_at"`
	RetryCount    int       `json:"retry_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
