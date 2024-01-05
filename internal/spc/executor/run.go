package executor

import (
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/SpecularL2/specular-cli/internal/service/config"
)

type RunHandler struct {
	cfg *config.Config
	log *logrus.Logger
}

func (r *RunHandler) Cmd() error {
	commandToRun := os.ExpandEnv(r.cfg.Args.Run.Command)
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

func NewRunHandler(cfg *config.Config, log *logrus.Logger) *RunHandler {
	return &RunHandler{
		cfg: cfg,
		log: log,
	}
}
