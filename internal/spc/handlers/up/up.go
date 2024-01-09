package up 

import (
	"github.com/sirupsen/logrus"

	"github.com/SpecularL2/specular-cli/internal/service/config"
)

const githubUrl = "https://api.github.com/repos/%s/contents/%s"

type UpHandler struct {
	cfg *config.Config
	log *logrus.Logger
}

func (u *UpHandler) Cmd() error {
	u.log.Warn("no command found, exiting...")
	return nil
}

func NewUpHandler(cfg *config.Config, log *logrus.Logger) *UpHandler {
	return &UpHandler{
		cfg: cfg,
		log: log,
	}
}

