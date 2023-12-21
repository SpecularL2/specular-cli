package di

import (
	"github.com/google/wire"

	"github.com/SpecularL2/specular-cli/internal/spc/workspace"
)

var CmdProvider = wire.NewSet( //nolint:gochecknoglobals
	workspace.NewWorkspaceHandler,
)
