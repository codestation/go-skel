// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

import (
	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"time"
)

// ID is used as an alias for the model primary key to avoid using some int by mistake.
type ID uint

// Model is the base that will be used by other entity who need a primary key and timestamps.
type Model struct {
	ID         ID          `json:"-"`
	ExternalID pgtype.UUID `json:"external_id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	DeletedAt  *time.Time  `json:"-"`
}

type Modelable[T any] interface {
	*T
	GetID() ID
	SetID(id ID)
}

// SetTimestamps configures the time on created/updated fields. Only call this method on new entity.
func (m *Model) SetTimestamps(now time.Time) {
	m.CreatedAt = now
	m.UpdatedAt = now
}

func (m *Model) GetID() ID {
	return m.ID
}

func (m *Model) SetID(id ID) {
	m.ID = id
}

type Option func(m *Model)

func WithTime(now time.Time) Option {
	return func(m *Model) {
		m.SetTimestamps(now)
	}
}

func WithUUID(id uuid.UUID) Option {
	return func(m *Model) {
		_ = m.ExternalID.Set(id)
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
	if e.ExternalID.Status == pgtype.Undefined {
		externalId := uuid.Must(uuid.NewV6())
		_ = e.ExternalID.Set(externalId)
	}
	return e
}
