package sqlstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"megpoid.xyz/go/go-skel/model"
)

func TestSqlHealthCheckStore(t *testing.T) {
	db := &FakeDbConn{}
	ss := New(db, model.SqlSettings{})

	err := ss.HealthCheck().HealthCheck(context.Background())
	assert.NoError(t, err)
}
