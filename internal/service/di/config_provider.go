package di

import (
	"github.com/google/wire"

	"github.com/SpecularL2/specular-cli/internal/service/config"
)

var ConfigProvider = wire.NewSet( //nolint:gochecknoglobals
	config.NewConfig,
)
