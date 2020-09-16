package models

import (
	"database/sql"
	"time"
)

// ModelID is used as a alias for the model primary key to avoid using some int by mistake.
type ModelID int

// Model is the base that will be used by other models who need a primary key and timestamps.
type Model struct {
	ID        ModelID      `json:"-"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"-"`
}

// SetTimestamps configures the time on created/updated fields. Only call this method on new models.
func (m *Model) SetTimestamps(now time.Time) {
	m.CreatedAt = now
	m.UpdatedAt = now
}
