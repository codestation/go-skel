// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"megpoid.dev/go/go-skel/app/model"
	"megpoid.dev/go/go-skel/pkg/paginator"
	"megpoid.dev/go/go-skel/pkg/repo"
	"megpoid.dev/go/go-skel/pkg/request"
	"megpoid.dev/go/go-skel/pkg/response"
	"megpoid.dev/go/go-skel/repository"
	"megpoid.dev/go/go-skel/repository/uow"
)

func TestProfileList(t *testing.T) {
	profiles := []*model.Profile{
		{
			Model: model.Model{ID: 1},
		},
	}

	mockResponse := response.NewListResponse(profiles, &paginator.Cursor{})
	r := repository.NewMockProfileRepo(t)
	r.EXPECT().List(mock.Anything, mock.Anything).Return(mockResponse, nil)

	u := uow.NewMockUnitOfWork(t)
	uc := NewProfile(u, r)

	result, err := uc.ListProfiles(context.Background(), &request.QueryParams{})
	assert.NoError(t, err)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, model.ID(1), result.Items[0].ID)
}

func TestProfileGet(t *testing.T) {
	mockProfile := model.Profile{
		Model: model.Model{ID: 1},
	}

	r := repository.NewMockProfileRepo(t)
	r.EXPECT().Get(mock.Anything, model.ID(1)).Return(&mockProfile, nil)

	u := uow.NewMockUnitOfWork(t)
	uc := NewProfile(u, r)

	profile, err := uc.GetProfile(context.Background(), model.ID(1))
	assert.NoError(t, err)
	assert.NotNil(t, profile)
}

func TestProfileSave(t *testing.T) {
	mockSave := func(ctx context.Context, profile *model.Profile) error {
		profile.ID = 1
		return nil
	}

	r := repository.NewMockProfileRepo(t)
	r.EXPECT().Save(mock.Anything, mock.Anything).RunAndReturn(mockSave)

	u := uow.NewMockUnitOfWork(t)
	uc := NewProfile(u, r)

	req := model.ProfileRequest{
		Email:     "john.doe@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	newProfile, err := uc.SaveProfile(context.Background(), &req)

	assert.NoError(t, err)
	assert.NotNil(t, newProfile)
	assert.Equal(t, model.ID(1), newProfile.ID)
}

func TestProfileUpdate(t *testing.T) {
	r := repository.NewMockProfileRepo(t)
	r.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)

	u := uow.NewMockUnitOfWork(t)
	uc := NewProfile(u, r)

	updateRequest := &model.ProfileRequest{
		Email: "test@test.com",
	}

	updated, err := uc.UpdateProfile(context.Background(), model.ID(1), updateRequest)
	assert.NoError(t, err)
	assert.Equal(t, "test@test.com", updated.Email)
}

func TestProfileDelete(t *testing.T) {
	r := repository.NewMockProfileRepo(t)
	r.EXPECT().Delete(mock.Anything, model.ID(1)).Return(nil)

	u := uow.NewMockUnitOfWork(t)
	uc := NewProfile(u, r)

	err := uc.RemoveProfile(context.Background(), model.ID(1))
	assert.NoError(t, err)
}

func TestProfileError(t *testing.T) {
	r := repository.NewMockProfileRepo(t)
	r.EXPECT().Get(mock.Anything, model.ID(1)).Return(nil, repo.ErrNotFound)

	u := uow.NewMockUnitOfWork(t)
	uc := NewProfile(u, r)

	_, err := uc.GetProfile(context.Background(), 1)
	assert.Error(t, err)
}
