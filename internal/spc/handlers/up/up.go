package up

import (
	"github.com/sirupsen/logrus"

	"github.com/SpecularL2/specular-cli/internal/service/config"
)

type UpHandler struct {
	cfg *config.Config
	log *logrus.Logger
}

func (u *UpHandler) Cmd() error {
	switch {
	case u.cfg.Args.Up.SpGeth != nil:
		return u.StartSpGeth()
	}
	u.log.Warn("no command found, exiting...")
	return nil
}

func (u *UpHandler) StartSpGeth() error {
	u.log.Warn("overidden flags:", u.cfg.Args.Up.SpGeth.Flags)
	return nil
}

func NewUpHandler(cfg *config.Config, log *logrus.Logger) *UpHandler {
	return &UpHandler{
		cfg: cfg,
		log: log,
	}
}
