// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.megpoid.dev/go-skel/app/usecase"
	"go.megpoid.dev/go-skel/config"
	"go.megpoid.dev/go-skel/oapi"
)

func TestHealthcheckController(t *testing.T) {
	suite.Run(t, &healthcheckSuite{})
}

type healthcheckSuite struct {
	suite.Suite
	cfg config.ServerSettings
}

func (s *healthcheckSuite) TestLive() {
	uc := usecase.NewMockHealthcheck(s.T())
	ctrl := NewHealthCheck(s.cfg, uc)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := ctrl.LiveCheck(ctx, oapi.LiveCheckParams{})
	s.NoError(err)
	s.Equal("ok", rec.Body.String())
	s.Equal(200, rec.Result().StatusCode)
}

func (s *healthcheckSuite) TestReady() {
	uc := usecase.NewMockHealthcheck(s.T())
	uc.EXPECT().Execute(mock.Anything).Return(nil)

	ctrl := NewHealthCheck(s.cfg, uc)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := ctrl.ReadyCheck(ctx, oapi.ReadyCheckParams{})
	s.NoError(err)
	s.Equal("ok", rec.Body.String())
	s.Equal(200, rec.Result().StatusCode)
}

func (s *healthcheckSuite) TestReadyFailed() {
	uc := usecase.NewMockHealthcheck(s.T())
	uc.EXPECT().Execute(mock.Anything).Return(errors.New("an error occurred"))

	ctrl := NewHealthCheck(s.cfg, uc)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := ctrl.ReadyCheck(ctx, oapi.ReadyCheckParams{})
	s.Error(err)
}
