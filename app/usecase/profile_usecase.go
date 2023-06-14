// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"
	"errors"

	"megpoid.dev/go/go-skel/app/model"
	"megpoid.dev/go/go-skel/app/repository/uow"
	"megpoid.dev/go/go-skel/pkg/clause"
	"megpoid.dev/go/go-skel/pkg/repo"
	"megpoid.dev/go/go-skel/pkg/request"
	"megpoid.dev/go/go-skel/pkg/response"
)

// used to validate that the implementation matches the interface
var _ Profile = &ProfileInteractor{}

type ProfileInteractor struct {
	common
	uow uow.UnitOfWork
}

func (u *ProfileInteractor) GetProfile(ctx context.Context, id int64) (*model.Profile, error) {
	t := u.printer(ctx)

	profile, err := u.uow.Store().Profiles().Get(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, NewAppError(t.Sprintf("Profile not found"), err)
		} else {
			return nil, NewAppError(t.Sprintf("Failed to get profile"), err)
		}
	}

	return profile, nil
}

func (u *ProfileInteractor) ListProfiles(ctx context.Context, query *request.QueryParams) (*response.ListResponse[*model.Profile], error) {
	t := u.printer(ctx)

	result, err := u.uow.Store().Profiles().List(ctx, clause.WithFilter(query))
	if err != nil {
		return nil, NewAppError(t.Sprintf("Failed to list profiles"), err)
	}

	return result, nil
}

func (u *ProfileInteractor) SaveProfile(ctx context.Context, req *model.ProfileRequest) (*model.Profile, error) {
	t := u.printer(ctx)

	profile := req.Profile()
	err := u.uow.Store().Profiles().Save(ctx, profile)
	if err != nil {
		if errors.Is(err, repo.ErrDuplicated) {
			return nil, NewAppError(t.Sprintf("Email is already registered with another profile"), err)
		}

		return nil, NewAppError(t.Sprintf("Failed to save profile"), err)
	}

	return profile, nil
}

func (u *ProfileInteractor) UpdateProfile(ctx context.Context, id int64, req *model.ProfileRequest) (*model.Profile, error) {
	t := u.printer(ctx)

	profile := req.Profile()
	profile.ID = id
	err := u.uow.Store().Profiles().Update(ctx, profile)
	if err != nil {
		return nil, NewAppError(t.Sprintf("Failed to update profile"), err)
	}

	return profile, nil
}

func (u *ProfileInteractor) RemoveProfile(ctx context.Context, id int64) error {
	t := u.printer(ctx)

	if err := u.uow.Store().Profiles().Delete(ctx, id); err != nil {
		return NewAppError(t.Sprintf("Failed to remove profile"), err)
	}

	return nil
}

func NewProfile(uow uow.UnitOfWork) *ProfileInteractor {
	return &ProfileInteractor{
		common: newCommon(),
		uow:    uow,
	}
}
