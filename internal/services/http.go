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

const (
	readTimeout     = 30 * time.Second
	writeTimeout    = 30 * time.Second
	shutdownTimeout = 30 * time.Second
)

type httpServer struct {
	HTTPServer *http.Server
	e          *echo.Echo
	notify     chan error
}

func NewHTTPServer(addr string, e *echo.Echo) *httpServer {
	svr := &http.Server{
		Addr:         addr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	server := &httpServer{HTTPServer: svr, e: e}
	server.start()

	return server
}

func (h *httpServer) Serve() error {
	if err := h.e.StartServer(h.HTTPServer); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (h *httpServer) start() {
	go func() {
		if err := h.e.StartServer(h.HTTPServer); err != nil && err != http.ErrServerClosed {
			h.notify <- err
		}
		close(h.notify)
	}()
}

func (h *httpServer) Notify() <-chan error {
	return h.notify
}

func (h *httpServer) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	return h.e.Shutdown(shutdownCtx)
}
