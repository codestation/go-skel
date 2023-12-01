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
	appmodel "go.megpoid.dev/go-skel/app/model"
	"go.megpoid.dev/go-skel/app/usecase"
	"go.megpoid.dev/go-skel/config"
	"go.megpoid.dev/go-skel/oapi"
	"go.megpoid.dev/go-skel/pkg/model"
	"go.megpoid.dev/go-skel/pkg/paginator"
	"go.megpoid.dev/go-skel/pkg/response"
)

func TestProfileController(t *testing.T) {
	suite.Run(t, &profileSuite{})
}

type profileSuite struct {
	suite.Suite
	cfg config.ServerSettings
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

	ctrl := NewProfile(s.cfg, uc)

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

	ctrl := NewProfile(s.cfg, uc)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := ctrl.ListProfiles(ctx, oapi.ListProfilesParams{})
	s.NoError(err)
}
