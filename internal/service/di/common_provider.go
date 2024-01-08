package di

import (
	"github.com/google/wire"

	"github.com/SpecularL2/specular-cli/internal/service/config"
)

var CommonProvider = wire.NewSet( //nolint:gochecknoglobals
	config.NewLogger,
	config.NewCancelChannel,
	config.NewContext,
)
