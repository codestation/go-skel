// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

import (
	"database/sql"
	"time"
)

// ID is used as an alias for the model primary key to avoid using some int by mistake.
type ID int64

// Model is the base that will be used by other entity who need a primary key and timestamps.
type Model struct {
	ID        ID           `json:"-"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"-"`
}

// SetTimestamps configures the time on created/updated fields. Only call this method on new entity.
func (m *Model) SetTimestamps(now time.Time) {
	m.CreatedAt = now
	m.UpdatedAt = now
}

func NewModel(now time.Time) *Model {
	e := &Model{}
	e.SetTimestamps(now)
	return e
}
