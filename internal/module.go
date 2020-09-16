package internal

import (
	"megpoid.xyz/go/go-skel/internal/config"
	"megpoid.xyz/go/go-skel/internal/delivery"
	"megpoid.xyz/go/go-skel/internal/repository"
	"megpoid.xyz/go/go-skel/internal/usecases"
	"megpoid.xyz/go/go-skel/pkg/sql/connection"
)

type Module struct {
	Usecase    usecase.APIUsecase
	Repository *repository.Repository
	Handler    *delivery.HTTPHandler
	Config     *config.Config
}

func New(db connection.SQLConnection, cfg *config.Config) *Module {
	repo := repository.NewRepository(db)
	uc := usecase.NewAPIUsecase(repo, cfg)
	handler := delivery.NewHTTPHandler(uc, cfg)

	return &Module{
		Usecase:    uc,
		Repository: repo,
		Handler:    handler,
		Config:     cfg,
	}
}
