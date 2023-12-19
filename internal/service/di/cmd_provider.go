package di

import (
	"github.com/SpecularL2/specular-cli/internal/spc/workspace"
	"github.com/google/wire"
)

var CmdProvider = wire.NewSet( //nolint:gochecknoglobals
	workspace.NewWorkspaceHandler,
)
