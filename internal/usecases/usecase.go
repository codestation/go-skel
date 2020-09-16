//go:generate mockgen -source=$GOFILE -destination=mocks/${GOFILE} -package=mocks
//go:generate mockgen  -destination=mocks/uuid.go -package=mocks github.com/gofrs/uuid Generator

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

type APIUsecase interface {
	HealthCheck
}

type mainUsecase struct {
	repo    *repository.Repository
	cfg     *config.Config
	uuidGen uuid.Generator
}

func NewAPIUsecase(repo *repository.Repository, cfg *config.Config) APIUsecase {
	return &mainUsecase{
		repo:    repo,
		cfg:     cfg,
		uuidGen: uuid.NewGen(),
	}
}
