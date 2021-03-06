/*
Copyright © 2020 codestation <codestation404@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package entity

import (
	"database/sql"
	"time"
)

// ID is used as a alias for the model primary key to avoid using some int by mistake.
type ID int

// Entity is the base that will be used by other entity who need a primary key and timestamps.
type Entity struct {
	ID        ID           `json:"-"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"-"`
}

// SetTimestamps configures the time on created/updated fields. Only call this method on new entity.
func (m *Entity) SetTimestamps(now time.Time) {
	m.CreatedAt = now
	m.UpdatedAt = now
}
