// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"
	"errors"

	"go.megpoid.dev/go-skel/app/model"
	"go.megpoid.dev/go-skel/app/repository"
	"go.megpoid.dev/go-skel/app/repository/uow"
	"go.megpoid.dev/go-skel/pkg/apperror"
	"go.megpoid.dev/go-skel/pkg/clause"
	"go.megpoid.dev/go-skel/pkg/repo"
	"go.megpoid.dev/go-skel/pkg/request"
	"go.megpoid.dev/go-skel/pkg/response"
)

// used to validate that the implementation matches the interface
var _ Profile = &ProfileInteractor{}

type ProfileInteractor struct {
	common
	uow         uow.UnitOfWork
	profileRepo repository.ProfileRepo
}

func (u *ProfileInteractor) GetProfile(ctx context.Context, id int64) (*model.Profile, error) {
	t := u.printer(ctx)

	profile, err := u.profileRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, apperror.NewAppError(t.Sprintf("Profile not found"), err)
		}

		return nil, apperror.NewAppError(t.Sprintf("Failed to get profile"), err)
	}

	return profile, nil
}

func (u *ProfileInteractor) ListProfiles(ctx context.Context, query *request.QueryParams) (*response.ListResponse[*model.Profile], error) {
	t := u.printer(ctx)

	result, err := u.profileRepo.List(ctx, clause.WithFilter(query))
	if err != nil {
		return nil, apperror.NewAppError(t.Sprintf("Failed to list profiles"), err)
	}

	return result, nil
}

func (u *ProfileInteractor) SaveProfile(ctx context.Context, req *model.ProfileRequest) (*model.Profile, error) {
	t := u.printer(ctx)

	profile := req.Profile()
	err := u.profileRepo.Insert(ctx, profile)
	if err != nil {
		if errors.Is(err, repo.ErrDuplicated) {
			return nil, apperror.NewAppError(t.Sprintf("Email is already registered with another profile"), err)
		}

		return nil, apperror.NewAppError(t.Sprintf("Failed to save profile"), err)
	}

	return profile, nil
}

func (u *ProfileInteractor) UpdateProfile(ctx context.Context, id int64, req *model.ProfileRequest) (*model.Profile, error) {
	t := u.printer(ctx)

	profile := req.Profile()
	profile.ID = id
	err := u.profileRepo.Update(ctx, profile)
	if err != nil {
		return nil, apperror.NewAppError(t.Sprintf("Failed to update profile"), err)
	}

	return profile, nil
}

func (u *ProfileInteractor) RemoveProfile(ctx context.Context, id int64) error {
	t := u.printer(ctx)

	if err := u.profileRepo.Delete(ctx, id); err != nil {
		return apperror.NewAppError(t.Sprintf("Failed to remove profile"), err)
	}

	return nil
}

func NewProfile(uow uow.UnitOfWork) *ProfileInteractor {
	return &ProfileInteractor{
		common:      newCommon(),
		uow:         uow,
		profileRepo: uow.Store().Profiles(),
	}
}
