// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	appmodel "go.megpoid.dev/go-skel/app/model"
	"go.megpoid.dev/go-skel/app/repository"
	"go.megpoid.dev/go-skel/app/repository/uow"
	"go.megpoid.dev/go-skel/pkg/model"
	"go.megpoid.dev/go-skel/pkg/paginator"
	"go.megpoid.dev/go-skel/pkg/repo"
	"go.megpoid.dev/go-skel/pkg/request"
	"go.megpoid.dev/go-skel/pkg/response"
)

func TestProfileList(t *testing.T) {
	profiles := []*appmodel.Profile{
		{
			Model: model.Model{ID: 1},
		},
	}

	mockResponse := response.NewListResponse(profiles, &paginator.Cursor{})

	r := repository.NewMockProfileRepo(t)
	r.EXPECT().List(mock.Anything, mock.Anything).Return(mockResponse, nil)

	store := uow.NewMockUnitOfWorkStore(t)
	store.EXPECT().Profiles().Return(r)

	u := uow.NewMockUnitOfWork(t)
	u.EXPECT().Store().Return(store)
	uc := NewProfile(u)

	result, err := uc.ListProfiles(context.Background(), &request.QueryParams{})
	assert.NoError(t, err)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, int64(1), result.Items[0].ID)
}

func TestProfileGet(t *testing.T) {
	mockProfile := appmodel.Profile{
		Model: model.Model{ID: 1},
	}

	r := repository.NewMockProfileRepo(t)
	r.EXPECT().Get(mock.Anything, int64(1)).Return(&mockProfile, nil)

	store := uow.NewMockUnitOfWorkStore(t)
	store.EXPECT().Profiles().Return(r)

	u := uow.NewMockUnitOfWork(t)
	u.EXPECT().Store().Return(store)
	uc := NewProfile(u)

	profile, err := uc.GetProfile(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
}

func TestProfileSave(t *testing.T) {
	mockSave := func(ctx context.Context, profile *appmodel.Profile) error {
		profile.ID = 1
		return nil
	}

	r := repository.NewMockProfileRepo(t)
	r.EXPECT().Insert(mock.Anything, mock.Anything).RunAndReturn(mockSave)

	store := uow.NewMockUnitOfWorkStore(t)
	store.EXPECT().Profiles().Return(r)

	u := uow.NewMockUnitOfWork(t)
	u.EXPECT().Store().Return(store)
	uc := NewProfile(u)

	req := appmodel.ProfileRequest{
		Email:     "john.doe@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	newProfile, err := uc.SaveProfile(context.Background(), &req)

	assert.NoError(t, err)
	assert.NotNil(t, newProfile)
	assert.Equal(t, int64(1), newProfile.ID)
}

func TestProfileUpdate(t *testing.T) {
	r := repository.NewMockProfileRepo(t)
	r.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)

	store := uow.NewMockUnitOfWorkStore(t)
	store.EXPECT().Profiles().Return(r)

	u := uow.NewMockUnitOfWork(t)
	u.EXPECT().Store().Return(store)
	uc := NewProfile(u)

	updateRequest := &appmodel.ProfileRequest{
		Email: "test@test.com",
	}

	updated, err := uc.UpdateProfile(context.Background(), 1, updateRequest)
	assert.NoError(t, err)
	assert.Equal(t, "test@test.com", updated.Email)
}

func TestProfileDelete(t *testing.T) {
	r := repository.NewMockProfileRepo(t)
	r.EXPECT().Delete(mock.Anything, int64(1)).Return(nil)

	store := uow.NewMockUnitOfWorkStore(t)
	store.EXPECT().Profiles().Return(r)

	u := uow.NewMockUnitOfWork(t)
	u.EXPECT().Store().Return(store)
	uc := NewProfile(u)

	err := uc.RemoveProfile(context.Background(), 1)
	assert.NoError(t, err)
}

func TestProfileError(t *testing.T) {
	r := repository.NewMockProfileRepo(t)
	r.EXPECT().Get(mock.Anything, int64(1)).Return(nil, repo.ErrNotFound)

	store := uow.NewMockUnitOfWorkStore(t)
	store.EXPECT().Profiles().Return(r)

	u := uow.NewMockUnitOfWork(t)
	u.EXPECT().Store().Return(store)
	uc := NewProfile(u)

	_, err := uc.GetProfile(context.Background(), 1)
	assert.Error(t, err)
}
