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

// Model is the base that will be used by other model who need a primary key and timestamps.
type Model struct {
	ID        int64      `json:"id" goqu:"skipinsert,skipupdate"`
	CreatedAt time.Time  `json:"created_at" goqu:"skipupdate"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type Modelable interface {
	GetID() int64
	SetID(id int64)
}

type Tabler interface {
	TableName() string
}

// SetTimestamps configures the time on created/updated fields. Only call this method on new model.
func (m *Model) SetTimestamps(now time.Time) {
	m.CreatedAt = now
	m.UpdatedAt = now
}

func (m *Model) GetID() int64 {
	return m.ID
}

func (m *Model) SetID(id int64) {
	m.ID = id
}

func (m *Model) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(m)
	}
}

func GetModelName[T any](m T) string {
	var t reflect.Type
	if t = reflect.TypeOf(m); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}

	return t.Name()
}

func GetTableName[T Modelable](m T) string {
	if m, ok := any(m).(Tabler); ok {
		return m.TableName()
	}

	name := GetModelName[T](m)

	return inflection.Plural(dbscan.SnakeCaseMapper(name))
}

type Option func(m *Model)

func WithID(id int64) Option {
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
