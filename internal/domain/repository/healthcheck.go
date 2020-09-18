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

package repository

import (
	"context"

	"megpoid.xyz/go/go-skel/pkg/sql/connection"
)

type healthRepo struct {
	db connection.SQLConnection
}

// NewHealthCheck creates a new repository with methods to check its health
func NewHealthCheck(db connection.SQLConnection) HealthCheck {
	return &healthRepo{db}
}

// HealthCheck returns an error if the database doesn't respond
func (r *healthRepo) HealthCheck(ctx context.Context) error {
	if err := r.db.PingContext(ctx); err != nil {
		return NewRepoError(ErrBackend, err)
	}
	return nil
}
