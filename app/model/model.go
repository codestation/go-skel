// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

import (
	"reflect"
	"time"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/jinzhu/inflection"
)

// ID is used as an alias for the model primary key to avoid using some int by mistake.
type ID uint

// Model is the base that will be used by other model who need a primary key and timestamps.
type Model struct {
	ID        ID         `json:"id" goqu:"skipinsert,skipupdate"`
	CreatedAt time.Time  `json:"created_at" goqu:"skipupdate"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type Modelable interface {
	GetID() ID
	SetID(id ID)
}

type Tabler interface {
	TableName() string
}

// SetTimestamps configures the time on created/updated fields. Only call this method on new model.
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

func (m *Model) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(m)
	}
}

func GetModelName[T any](m T) string {
	if t := reflect.TypeOf(m); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func GetTableName[T Modelable](m T) string {
	if m, ok := any(m).(Tabler); ok {
		return m.TableName()
	}

	name := GetModelName[T](m)

	return inflection.Plural(dbscan.SnakeCaseMapper(name))
}

type Option func(m *Model)

func WithID(id ID) Option {
	return func(m *Model) {
		m.ID = id
	}
}

func WithTime(now time.Time) Option {
	return func(m *Model) {
		m.SetTimestamps(now)
	}
}

func NewModel(opts ...Option) Model {
	e := Model{}
	e.Apply(opts...)
	if e.CreatedAt.IsZero() {
		e.SetTimestamps(time.Now())
	}

	return e
}
