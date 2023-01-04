package usecase

import (
	"context"
	"gateway-api/domain"
	"gateway-api/helpers/response"
	"time"
)

type GatewayUsecase interface {
	Dynamic(ctx context.Context, payload domain.DefaultPayload) response.Base
	GetConfig() domain.Config
}

type gatewayUsecase struct {
	contextTimeout time.Duration
	config         domain.Config
}

func NewUsecase(timeout time.Duration, config domain.Config) GatewayUsecase {
	return &gatewayUsecase{
		contextTimeout: timeout,
		config:         config,
	}
}

func (u *gatewayUsecase) GetConfig() domain.Config {
	return u.config
}
