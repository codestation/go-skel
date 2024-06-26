// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repo

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"go.megpoid.dev/go-skel/pkg/clause"
	"go.megpoid.dev/go-skel/pkg/model"
	"go.megpoid.dev/go-skel/pkg/paginator"
	"go.megpoid.dev/go-skel/pkg/types"
)

type jsonField struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func (j *jsonField) Scan(src any) error {
	var source []byte

	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	case nil:
		source = nil
	default:
		return errors.New("incompatible type")
	}

	return json.Unmarshal(source, j)
}

func (j jsonField) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type jsonFieldSlice []jsonField

func (j *jsonFieldSlice) Scan(src any) error {
	var source []byte

	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	case nil:
		source = nil
	default:
		return errors.New("incompatible type")
	}

	return json.Unmarshal(source, j)
}

func (j jsonFieldSlice) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type testProfile struct {
	model.Model
	ExternalID uuid.UUID `json:"external_id"`
	Avatar     string
}

type testUser struct {
	model.Model
	Name       string
	ExternalID uuid.UUID      `json:"external_id"`
	ProfileID  int64          `goqu:"skipupdate"`
	Code       int            `goqu:"skipinsert,skipupdate"`
	Data1      jsonField      `json:"data1"`
	Data2      jsonFieldSlice `json:"data2"`
	Profile    *testProfile
}

func (t *testUser) AttachProfile(p *testProfile) {
	t.ProfileID = 0
	t.Profile = p
}

func newUser(name string, profileId int64) *testUser {
	u := &testUser{
		Model:     model.NewModel(),
		Name:      name,
		ProfileID: profileId,
	}
	return u
}

type userStore struct {
	*GenericStoreImpl[*testUser]
	profile *profileStore
}

type profileStore struct {
	*GenericStoreImpl[*testProfile]
}

func (s *userStore) Attach(ctx context.Context, results []*testUser, relation string) error {
	var err error
	switch relation {
	case "profile":
		err = AttachRelation(ctx, results,
			func(m *testUser) *int64 { return types.AsPointer(m.ProfileID) },
			func(m *testUser, r *testProfile) { m.AttachProfile(r) },
			s.profile.ListByIDs)
	default:
		return ErrBackend
	}
	return err
}

func TestStore(t *testing.T) {
	suite.Run(t, &storeSuite{})
}

type storeSuite struct {
	suite.Suite
	conn *Connection
}

func (s *storeSuite) SetupTest() {
	s.conn = NewTestConnection(s.T(), true)
}

func (s *storeSuite) TearDownTest() {
	if s.conn != nil {
		s.conn.Close(s.T())
	}
}

func (s *storeSuite) TestNewStore() {
	st := NewStore[*testUser](s.conn.Store)
	s.Equal("test_users", st.Table)
	s.Equal([]any{"*"}, st.selectFields)
}

func (s *storeSuite) TestStoreFind() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		id  int64
		err error
	}{
		{1, nil},
		{0, ErrNotFound},
	}

	for _, test := range tests {
		s.Run("Find", func() {
			var user testUser
			err := st.Find(context.Background(), &user, test.id)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
				s.NotZero(user.ID)
				s.NotZero(user.CreatedAt)
			}
		})
	}
}

func (s *storeSuite) TestStoreCount() {
	st := NewStore[*testUser](s.conn.Store)
	count, err := st.CountBy(context.Background(), Ex{})
	s.NoError(err)
	s.Equal(int64(5), count)
}

func (s *storeSuite) TestStoreCountBy() {
	st := NewStore[*testUser](s.conn.Store)
	count, err := st.CountBy(context.Background(), Ex{"name": "John Doe 1"})
	s.NoError(err)
	s.Equal(int64(1), count)
}

func (s *storeSuite) TestStoreFirst() {
	st := NewStore[*testUser](s.conn.Store)
	user, err := st.First(context.Background(), Ex{}, C("id").Asc())
	s.NoError(err)
	s.NotNil(user)
}

func (s *storeSuite) TestStoreFirstWithFilters() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		name string
		err  error
	}{
		{"John Doe 1", nil},
		{"John Doe 6", ErrNotFound},
	}

	for _, test := range tests {
		s.Run("First", func() {
			user, err := st.First(context.Background(), Ex{"name": test.name}, C("id").Asc())
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
				s.NotZero(user.ID)
				s.Equal(test.name, user.Name)
			}
		})
	}
}

func (s *storeSuite) TestStoreGet() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		id  int64
		err error
	}{
		{1, nil},
		{0, ErrNotFound},
	}

	for _, test := range tests {
		s.Run("Get", func() {
			user, err := st.Get(context.Background(), test.id)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
				s.NotZero(user.ID)
				s.NotZero(user.CreatedAt)
			}
		})
	}
}

func (s *storeSuite) TestStoreGetForUpdate() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		expr  Expression
		err   error
		order OrderedExpression
	}{
		{goqu.C("id").Gt(0), nil, goqu.C("id").Asc()},
		{goqu.C("id").Eq(0), ErrNotFound, goqu.C("id").Asc()},
	}

	for _, test := range tests {
		s.Run("GetForUpdate", func() {
			user, err := st.GetForUpdate(context.Background(), test.expr, test.order)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
				s.NotZero(user.ID)
				s.NotZero(user.CreatedAt)
			}
		})
	}
}

func (s *storeSuite) TestStoreList() {
	st := NewStore[*testUser](s.conn.Store)
	users, err := st.List(context.Background())
	if s.NoError(err) {
		s.GreaterOrEqual(len(users.Items), 0)
	}
}

func (s *storeSuite) TestStoreSave() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		name      string
		profileId int64
		data1     jsonField
		data2     []jsonField
		err       error
	}{
		{"Some user", 1, jsonField{"foo", 1}, []jsonField{{"foo", 1}}, nil},
		{"Some user", 1, jsonField{"foo", 1}, []jsonField{{"foo", 1}}, ErrDuplicated}, // do not run more tests after a constraint error
	}

	for _, test := range tests {
		s.Run("Insert", func() {
			user := newUser(test.name, test.profileId)
			user.Data1 = test.data1
			user.Data2 = test.data2
			err := st.Insert(context.Background(), user)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
				s.NotZero(user.ID)
			}
		})
	}
}

func (s *storeSuite) TestStoreSaveExtraField() {
	st := NewStore[*testUser](s.conn.Store, WithReturnFields[*testUser]("code"))
	tests := []struct {
		name      string
		profileId int64
		err       error
	}{
		{"Some user", 1, nil},
		{"Some user", 1, ErrDuplicated}, // do not run more tests after a constraint error
	}

	for _, test := range tests {
		s.Run("Insert", func() {
			user := newUser(test.name, test.profileId)
			err := st.Insert(context.Background(), user)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
				s.NotZero(user.ID)
				s.NotZero(user.Code)
			}
		})
	}
}

func (s *storeSuite) TestStoreUpsert() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		name      string
		profileId int64
		created   bool
		err       error
	}{
		{"Some user 1", 1, true, nil},
		{"Some user 2", 1, true, nil},
		{"Some user 2", 1, false, nil},
	}

	for _, test := range tests {
		s.Run("Insert", func() {
			user := newUser(test.name, test.profileId)
			created, err := st.Upsert(context.Background(), user, "name")
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
				s.NotZero(user.ID)
				s.Equal(test.created, created)
			}
		})
	}
}

func (s *storeSuite) TestStoreUpsertExtraField() {
	st := NewStore[*testUser](s.conn.Store, WithReturnFields[*testUser]("code"))
	tests := []struct {
		name      string
		profileId int64
		created   bool
		err       error
	}{
		{"Some user 1", 1, true, nil},
		{"Some user 2", 1, true, nil},
		{"Some user 2", 1, false, nil},
	}

	for _, test := range tests {
		s.Run("Insert", func() {
			user := newUser(test.name, test.profileId)
			created, err := st.Upsert(context.Background(), user, "name")
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
				s.NotZero(user.ID)
				s.Equal(test.created, created)
				s.NotZero(user.Code)
			}
		})
	}
}

func (s *storeSuite) TestStoreUpdate() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		id  int64
		err error
	}{
		{1, nil},
		{0, ErrNotFound},
	}

	for _, test := range tests {
		s.Run("Update", func() {
			user := newUser("John Doe", 1)
			user.ID = test.id
			user.UpdatedAt = time.Now()
			err := st.Update(context.Background(), user)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *storeSuite) TestStoreUpdateMap() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		id  int64
		err error
	}{
		{1, nil},
		{0, ErrNotFound},
	}

	for _, test := range tests {
		s.Run("UpdateMap", func() {
			user := Ex{
				"name":       "John Doe",
				"profile_id": 1,
				"updated_at": time.Now(),
			}
			err := st.UpdateMap(context.Background(), test.id, user)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *storeSuite) TestStoreUpdateMapBy() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		expr Ex
		err  error
		n    int64
	}{
		{Ex{"id": 1}, nil, 1},
		{Ex{"id": 0}, nil, 0},
	}

	for _, test := range tests {
		s.Run("UpdateMapBy", func() {
			user := Ex{
				"name":       "John Doe",
				"profile_id": 1,
				"updated_at": time.Now(),
			}
			n, err := st.UpdateMapBy(context.Background(), user, test.expr)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else if s.NoError(err) {
				s.Equal(test.n, n)
			}
		})
	}
}

func (s *storeSuite) TestStoreDelete() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		id  int64
		err error
	}{
		{1, nil},
		{0, ErrNotFound},
	}

	for _, test := range tests {
		s.Run("Delete", func() {
			err := st.Delete(context.Background(), test.id)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *storeSuite) TestStoreGetExternal() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		id  uuid.UUID
		err error
	}{
		{uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")), nil},
		{uuid.Must(uuid.NewV7()), ErrNotFound},
	}

	for _, test := range tests {
		s.Run("GetExternal", func() {
			user, err := st.GetBy(context.Background(), Ex{"external_id": test.id})
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else if s.NoError(err) {
				s.NotNil(user)
				s.NotZero(user.ID)
				s.NotZero(user.CreatedAt)
			}
		})
	}
}

func (s *storeSuite) TestStoreGetByExpr() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		name string
		err  error
	}{
		{"John Doe 1", nil},
		{"John Doe 6", ErrNotFound},
	}

	for _, test := range tests {
		s.Run("GetBy", func() {
			user, err := st.GetBy(context.Background(), Ex{"name": test.name})
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else if s.NoError(err) {
				s.NotNil(user)
				s.NotZero(user.ID)
				s.NotZero(user.CreatedAt)
			}
		})
	}
}

func (s *storeSuite) TestStoreExists() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		name   string
		exists bool
	}{
		{"John Doe 1", true},
		{"John Doe 6", false},
	}

	for _, test := range tests {
		s.Run("Exists", func() {
			exists, err := st.Exists(context.Background(), Ex{"name": test.name})
			if s.NoError(err) {
				s.Equal(test.exists, exists)
			}
		})
	}
}

func (s *storeSuite) TestStoreExpression() {
	st := NewStore[*testUser](s.conn.Store)
	exists, err := st.Exists(context.Background(), And(
		I("name").Eq("John Doe 1"),
		Or(
			C("name").IsNotNull(),
			L("name").IsNull(),
		),
	))
	if s.NoError(err) {
		s.True(exists)
	}
}

func (s *storeSuite) TestStoreDeleteExternal() {
	st := NewStore[*testUser](s.conn.Store)
	tests := []struct {
		id    uuid.UUID
		err   error
		count int64
	}{
		{uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")), nil, 1},
		{uuid.Must(uuid.NewV7()), nil, 0},
	}

	for _, test := range tests {
		s.Run("DeleteExternal", func() {
			n, err := st.DeleteBy(context.Background(), Ex{"external_id": test.id})
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
				s.Equal(n, test.count)
			}
		})
	}
}

func (s *storeSuite) TestBackendError() {
	db := &fakeDatabase{
		Error: errors.New("not implemented"),
	}

	st := NewStore[*testUser](db)
	ctx := context.Background()

	_, err := st.Get(ctx, 1)
	s.ErrorIs(err, ErrBackend)
	_, err = st.List(ctx)
	s.ErrorIs(err, ErrBackend)
	err = st.Insert(ctx, newUser("John Doe", 1))
	s.ErrorIs(err, ErrBackend)
	err = st.Update(ctx, newUser("John Doe", 1))
	s.ErrorIs(err, ErrBackend)
	err = st.Delete(ctx, 1)
	s.ErrorIs(err, ErrBackend)
	_, err = st.GetBy(ctx, Ex{"external_id": uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001"))})
	s.ErrorIs(err, ErrBackend)
	_, err = st.DeleteBy(ctx, Ex{"external_id": uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001"))})
	s.ErrorIs(err, ErrBackend)

	db.Error = errors.New("not implemented")
	err = st.Update(ctx, newUser("John Doe", 1))
	s.ErrorIs(err, ErrBackend)
	err = st.Delete(ctx, 1)
	s.ErrorIs(err, ErrBackend)
	_, err = st.DeleteBy(ctx, Ex{"external_id": uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001"))})
	s.ErrorIs(err, ErrBackend)
}

func (s *storeSuite) TestIncludes() {
	st := userStore{
		GenericStoreImpl: NewStore[*testUser](s.conn.Store,
			WithIncludes[*testUser]("profile"),
		),
	}

	st.profile = &profileStore{
		GenericStoreImpl: NewStore[*testProfile](s.conn.Store),
	}

	st.AttachFunc(st.Attach)

	users, err := st.List(context.Background(), clause.WithIncludes("profile"))
	if s.NoError(err) {
		s.Equal(5, len(users.Items))
		user := users.Items[0]
		s.Zero(user.ProfileID)
		s.NotNil(user.Profile)
		s.Equal(int64(1), user.Profile.ID)
	}
}

func (s *storeSuite) TestEach() {
	st := userStore{GenericStoreImpl: NewStore[*testUser](s.conn.Store,
		WithPaginatorOptions[*testUser](
			paginator.WithLimit(2),
		),
	)}

	count := 0
	err := st.ListEach(context.Background(), func(entry *testUser) error {
		count += 1
		return nil
	})
	s.NoError(err)
	s.Equal(5, count)
}

func (s *storeSuite) TestWithFilters() {
	st := userStore{GenericStoreImpl: NewStore[*testUser](s.conn.Store,
		WithExpressions[*testUser](goqu.Ex{"name": "John Doe 3"}),
	)}

	response, err := st.List(context.Background())
	s.NoError(err)
	s.Equal(1, len(response.Items))
}

func (s *storeSuite) TestEmptyResult() {
	st := userStore{GenericStoreImpl: NewStore[*testUser](s.conn.Store)}

	response, err := st.ListBy(context.Background(), Ex{"name": "Not Found"})
	s.NoError(err)
	s.Equal(0, len(response.Items))
}

func (s *storeSuite) TestPrefix() {
	st := NewStore[*testUser](s.conn.Store, WithTablePrefix[*testUser]("app_"))
	s.Equal("app_test_users", st.Table)
}
