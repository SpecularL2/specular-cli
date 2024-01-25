package exec

import (
	"github.com/SpecularL2/specular-cli/internal/spc/handlers/workspace"

	"github.com/sirupsen/logrus"

	"github.com/SpecularL2/specular-cli/internal/service/config"
)

type RunHandler struct {
	cfg       *config.Config
	log       *logrus.Logger
	workspace *workspace.WorkspaceHandler
}

func (r *RunHandler) Cmd() error {
	strCmd := r.cfg.Args.Exec.Command
	cmd, err := r.workspace.RunStringCommand(strCmd)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func NewRunHandler(cfg *config.Config, log *logrus.Logger, w *workspace.WorkspaceHandler) *RunHandler {
	return &RunHandler{
		cfg:       cfg,
		log:       log,
		workspace: w,
	}
}
