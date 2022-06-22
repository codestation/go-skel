// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

import (
	"time"
)

// ID is used as an alias for the model primary key to avoid using some int by mistake.
type ID uint

// Model is the base that will be used by other entity who need a primary key and timestamps.
type Model struct {
	ID        ID         `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// SetTimestamps configures the time on created/updated fields. Only call this method on new entity.
func (m *Model) SetTimestamps(now time.Time) {
	m.CreatedAt = now
	m.UpdatedAt = now
}

type Option func(m *Model)

func WithTime(now time.Time) Option {
	return func(m *Model) {
		m.SetTimestamps(now)
	}
}

func NewModel(opts ...Option) Model {
	e := Model{}
	for _, opt := range opts {
		opt(&e)
	}
	if e.CreatedAt.IsZero() {
		e.SetTimestamps(time.Now())
	}
	return e
}
