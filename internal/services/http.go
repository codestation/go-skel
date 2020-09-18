/*
Copyright © 2020 codestation <codestation404@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package services

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type httpServer struct {
	HTTPServer *http.Server
	e          *echo.Echo
}

func NewHTTPServer(addr string, e *echo.Echo) *httpServer {
	svr := &http.Server{
		Addr:         addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return &httpServer{HTTPServer: svr, e: e}
}

func (h *httpServer) Serve() error {
	if err := h.e.StartServer(h.HTTPServer); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (h *httpServer) Shutdown(ctx context.Context) error {
	return h.e.Shutdown(ctx)
}
