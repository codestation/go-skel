//go:generate mockgen -source=$GOFILE -destination=mocks/${GOFILE} -package=mocks
//go:generate mockgen  -destination=mocks/uuid.go -package=mocks github.com/gofrs/uuid Generator
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

package usecase

import (
	"context"

	"github.com/gofrs/uuid"
	"megpoid.xyz/go/go-skel/internal/config"
	"megpoid.xyz/go/go-skel/internal/repository"
)

type HealthCheck interface {
	HealthCheck(ctx context.Context) error
}

type UseCase interface {
	HealthCheck
}

type usecase struct {
	repo    *repository.Repository
	cfg     *config.Config
	uuidGen uuid.Generator
}

func NewUseCase(repo *repository.Repository, cfg *config.Config) UseCase {
	return &usecase{
		repo:    repo,
		cfg:     cfg,
		uuidGen: uuid.NewGen(),
	}
}
