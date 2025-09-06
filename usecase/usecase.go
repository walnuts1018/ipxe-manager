package usecase

import (
	"context"

	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/ipxe-manager/config"
	"github.com/walnuts1018/ipxe-manager/domain/entity"
	"github.com/walnuts1018/ipxe-manager/util/clock"
	"github.com/walnuts1018/ipxe-manager/util/random"
)

type OSRepository interface {
	SetOS(ctx context.Context, os string) error
	GetOS(ctx context.Context) (string, error)
}

type AuthService interface {
	IntrospectToken(ctx context.Context, token string) (entity.IntrospectionResponse, error)
}

type Usecase struct {
	cfg *config.Config

	// Repository
	OSRepository OSRepository

	// Service
	authService AuthService

	random random.Random
	clock  clock.Clock[tz.AsiaTokyo]
}

func NewUsecase(
	cfg *config.Config,
	OSRepository OSRepository,
	authService AuthService,
	random random.Random,
	clock clock.Clock[tz.AsiaTokyo],
) *Usecase {
	return &Usecase{
		cfg:          cfg,
		OSRepository: OSRepository,
		authService:  authService,
		random:       random,
		clock:        clock,
	}
}
