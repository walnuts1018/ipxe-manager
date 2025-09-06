//go:build wireinject
// +build wireinject

package wire

import (
	"context"

	"github.com/Code-Hex/synchro/tz"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/walnuts1018/ipxe-manager/config"
	"github.com/walnuts1018/ipxe-manager/infrastructure/auth"
	"github.com/walnuts1018/ipxe-manager/infrastructure/database"
	"github.com/walnuts1018/ipxe-manager/router"
	"github.com/walnuts1018/ipxe-manager/router/handler"
	"github.com/walnuts1018/ipxe-manager/usecase"
	"github.com/walnuts1018/ipxe-manager/util/clock"
	"github.com/walnuts1018/ipxe-manager/util/random"
)

func CreateRouter(
	ctx context.Context,
	cfg *config.Config,
	db *database.DB,
	clock clock.Clock[tz.AsiaTokyo],
) (*echo.Echo, error) {
	wire.Build(
		handler.NewHandler,
		router.NewRouter,
		usecase.NewUsecase,
		configSet,
		random.New,
		authSet,
		wire.Bind(new(usecase.OSRepository), new(*database.DB)),
	)
	return &echo.Echo{}, nil
}

var authSet = wire.NewSet(
	auth.NewAuthService,
	wire.Bind(new(usecase.AuthService), new(*auth.AuthService)),
)

var configSet = wire.FieldsOf(new(*config.Config),
	"App",
	"OAuth2",
)
