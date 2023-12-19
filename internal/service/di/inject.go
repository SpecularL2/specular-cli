//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	"github.com/SpecularL2/specular-cli/internal/service/config"
)

func SetupApplication() (*Application, func(), error) {
	panic(wire.Build(wire.NewSet(
		CommonProvider,
		ConfigProvider,
		CmdProvider,
		wire.Struct(new(Application), "*"))),
	)
}

func SetupApplicationForIntegrationTests(cfg *config.Config) (*TestApplication, func(), error) {
	panic(wire.Build(wire.NewSet(
		CommonProvider,
		CmdProvider,
		wire.Struct(new(Application), "*"),
		wire.Struct(new(TestApplication), "*"))),
	)
}
