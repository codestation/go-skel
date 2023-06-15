// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	appmodel "megpoid.dev/go/go-skel/app/model"
	"megpoid.dev/go/go-skel/app/usecase"
	"megpoid.dev/go/go-skel/oapi"
	"megpoid.dev/go/go-skel/pkg/model"
	"megpoid.dev/go/go-skel/pkg/paginator"
	"megpoid.dev/go/go-skel/pkg/response"
)

func TestProfileController(t *testing.T) {
	suite.Run(t, &profileSuite{})
}

type profileSuite struct {
	suite.Suite
}

func (s *profileSuite) TestGet() {
	mockProfile := appmodel.Profile{
		Model:     model.Model{ID: 1},
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	uc := usecase.NewMockProfile(s.T())
	uc.EXPECT().GetProfile(mock.Anything, int64(1)).Return(&mockProfile, nil)

	ctrl := NewProfile(nil, uc)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := ctrl.GetProfile(ctx, 1)
	s.NoError(err)
}

func (s *profileSuite) TestList() {
	mockProfiles := []*appmodel.Profile{{
		Model:     model.Model{ID: 1},
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}}

	resp := response.NewListResponse(mockProfiles, &paginator.Cursor{})

	uc := usecase.NewMockProfile(s.T())
	uc.EXPECT().ListProfiles(mock.Anything, mock.Anything).Return(resp, nil)

	ctrl := NewProfile(nil, uc)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := ctrl.ListProfiles(ctx, oapi.ListProfilesParams{})
	s.NoError(err)
}
