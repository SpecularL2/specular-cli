package integration

import (
	"github.com/SpecularL2/specular-cli/internal/service/config"
)

type SpcCluster interface {
	Close() error
	Workspace() *config.Workspace
}

type ServerInstance interface {
	Port() (int, error)
	Close() error
	Address() (string, error)
	Prep() error
}
