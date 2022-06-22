// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package store

import (
	"context"
)

// Store lists all the other stores
type Store interface {
	HealthCheck() HealthCheckStore
	WithTransaction(ctx context.Context, f func(s Store) error) error
}

// HealthCheckStore handles all healthCheck related operations on the store
type HealthCheckStore interface {
	HealthCheck(ctx context.Context) error
}
