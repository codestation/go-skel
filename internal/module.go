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

package internal

import (
	"megpoid.xyz/go/go-skel/internal/config"
	"megpoid.xyz/go/go-skel/internal/controller/http"
	"megpoid.xyz/go/go-skel/internal/repository"
	"megpoid.xyz/go/go-skel/internal/usecase"
	"megpoid.xyz/go/go-skel/pkg/sql/connection"
)

type Module struct {
	Usecase    usecase.UseCase
	Repository *repository.Repository
	Handler    *http.HTTPHandler
	Config     *config.Config
}

func New(db connection.SQLConnection, cfg *config.Config) *Module {
	repo := repository.NewRepository(db)
	uc := usecase.NewUseCase(repo, cfg)
	handler := http.NewHTTPHandler(uc, cfg)

	return &Module{
		Usecase:    uc,
		Repository: repo,
		Handler:    handler,
		Config:     cfg,
	}
}
