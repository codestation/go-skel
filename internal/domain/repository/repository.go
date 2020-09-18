//go:generate mockgen -source=$GOFILE -destination=mocks/${GOFILE} -package=mocks
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

// HealthCheck handles all health-related operations on the repository
type HealthCheck interface {
	HealthCheck(ctx context.Context) error
}

// Transaction allow to call multiple repository methods and make their execution atomic
type Transaction interface {
	WithTransaction(ctx context.Context, txFunc func(repo *Repository) error) error
}

// Repository holds the rest of the repository interfaces
type Repository struct {
	HealthCheck
	Transaction
}

// NewRepository created a new repository for the connection
func NewRepository(db connection.SQLConnection) *Repository {
	return &Repository{
		HealthCheck: NewHealthCheck(db),
		Transaction: NewTransaction(db),
	}
}
