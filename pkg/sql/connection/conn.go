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

package connection

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Driver handles database connections from sqlx.tx or sqlx.tx
type Driver interface {
	sqlx.ExecerContext
	sqlx.QueryerContext
	sqlx.PreparerContext
}

// SQLConnection adds transaction support for sql databases
type SQLConnection interface {
	Driver
	TxBegin(ctx context.Context) (SQLConnection, error)
	TxEnd(txFunc func() error) error
	Commit() error
	Rollback() error
	PingContext(ctx context.Context) error
}
