// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"megpoid.xyz/go/go-skel/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlHealthCheckStore(t *testing.T) {
	db := &FakeDbConn{}
	ss := New(db, config.SqlSettings{})

	err := ss.HealthCheck().HealthCheck(context.Background())
	assert.NoError(t, err)
}
