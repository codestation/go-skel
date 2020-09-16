package usecase

import "context"

func (uc *mainUsecase) HealthCheck(ctx context.Context) error {
	return uc.repo.HealthCheck.HealthCheck(ctx)
}
