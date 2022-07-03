// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type LocalCompany struct {
	Model
	Name int
}

type GlobalCompany struct {
	LocalCompany
}

func (g GlobalCompany) TableName() string {
	return "company"
}

func NewLocalCompany(opts ...Option) *LocalCompany {
	return &LocalCompany{Model: NewModel(opts...)}
}

func TestNewModel(t *testing.T) {
	now := time.Now()
	model := NewLocalCompany(WithTime(now))
	assert.Equal(t, model.CreatedAt, now)
	assert.Equal(t, model.UpdatedAt, now)
}

func TestTableName(t *testing.T) {
	name := GetTableName(&LocalCompany{})
	assert.Equal(t, "local_companies", name)

	name = GetTableName(&GlobalCompany{})
	assert.Equal(t, "company", name)
}
