package config

import (
	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
)

const (
	defaultLogLevel = logrus.InfoLevel
	serviceName     = "spc"
)

type GethConfig struct {
}

type SpMagiConfig struct {
	Network             string
	L1RpcURL            string
	L2RpcURL            string
	SyncMode            string
	L2EngineURL         string
	JWTSecretPath       string
	RpcPort             string
	SequencerMaxSafeLag string
	SequencerPkFile     string
	CheckpointSyncUrl   string
	CheckpointHash      string
}

type Workspace struct {
	// TODO: abstract structure to hold all values loaded from the active workspace/test or set values via CLI tool run
	L1GethURL    string
	SpMagiConfig SpMagiConfig
}

type WorkspaceCmd struct {
	Command string `arg:"positional"`
	Name    string `arg:"positional"`
}

type Config struct {
	LogLevel     string
	WorkspaceCmd *WorkspaceCmd `arg:"subcommand:workspace"`
	Workspace    *Workspace
}

func (c *Config) Description() string {
	return "Specular CLI - toolkit for L2 integration and testing"
}

func (c *Config) GetLogLevel(defaultLevel logrus.Level) logrus.Level {
	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		level = defaultLevel
	}
	return level
}

func NewConfig() (*Config, error) {
	cfg := Config{}
	arg.MustParse(&cfg)
	return &cfg, nil
}
