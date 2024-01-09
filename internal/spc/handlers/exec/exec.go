package exec 

import (
	"os"
	"os/exec"
	"strings"

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
	return r.RunStringCommand()
}

func (r *RunHandler) RunStringCommand() error {
	// TODO: handle case where there is no active workspace
	err := r.workspace.LoadWorkspaceEnvVars()
	if err != nil {
		return err
	}

	commandToRun := os.ExpandEnv(r.cfg.Args.Exec.Command)
	args := strings.Fields(commandToRun)

	if len(args) > 0 {
		r.log.Debugf("Running: %s %v", commandToRun, args)
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			}
		}
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
