package di

import (
	"github.com/google/wire"

	"github.com/SpecularL2/specular-cli/internal/spc/handlers/exec"
	"github.com/SpecularL2/specular-cli/internal/spc/handlers/up"
	"github.com/SpecularL2/specular-cli/internal/spc/handlers/workspace"
)

var CmdProvider = wire.NewSet( //nolint:gochecknoglobals
	workspace.NewWorkspaceHandler,
	exec.NewRunHandler,
	up.NewUpHandler,
)
