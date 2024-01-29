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
	GethPortHTTP string `arg:"env:GETH_PORT_HTTP"`
	GethPortTCP  string
}

func (gc *GethConfig) Args() []string {
	return []string{
		"--dev",
		// TODO: instead of dev.period we should probably have a function to trigger the mine manually via geth API,
		//       but it may depend on the type of test. For L2 and others we could have similar capability.
		"--dev.period", "2",
		"--http",
		"--http.addr", "0.0.0.0",
		"--http.port", gc.GethPortHTTP,
		"--http.api", "engine,personal,eth,net,web3,txpool,miner,debug",
		"--http.corsdomain", "*",
		"--http.vhosts", "*",
		"--ws",
		"--ws.addr", "0.0.0.0",
		"--ws.port", gc.GethPortTCP,
		"--ws.api", "engine,personal,eth,net,web3,txpool,miner,debug",
		"--ws.origins", "\"*\"",
		"--authrpc.vhosts", "*",
		"--authrpc.addr", "0.0.0.0",
		// "--authrpc.port", "$AUTH_PORT",
		// "--authrpc.jwtsecret", "$JWT_SECRET_PATH",
		"--miner.recommit", "0",
	}
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

func (smc *SpMagiConfig) Args() []string {
	return []string{
		"--devnet",
		"--network", smc.Network,
		"--l1-rpc-url", smc.L1RpcURL,
		"--l2-rpc-url", smc.L2RpcURL,
		"--sync-mode", smc.SyncMode,
		"--l2-engine-url", smc.L2EngineURL,
		"--jwt-file", smc.JWTSecretPath,
		"--rpc-port", smc.RpcPort,
		"--sequencer",
		"--sequencer-max-safe-lag", smc.SequencerMaxSafeLag,
		"--sequencer-pk-file", smc.SequencerPkFile,
		"--checkpoint-sync-url", smc.CheckpointSyncUrl,
		"--checkpoint-hash", smc.CheckpointHash,
	}
}

type Workspace struct {
	// TODO: abstract structure to hold all values loaded from the active workspace/test or set values via CLI tool run
	L1GethURL    string
	GethConfig   GethConfig
	SpMagiConfig SpMagiConfig
}

type DownloadCmd struct {
	ConfigPath string `arg:"--config-path" default:"config/local_devnet" help:"path of the workspace config"`
	ConfigRepo string `arg:"--config-repository" default:"specularL2/specular" help:"github repository to pull config from"`
}

type SetCmd struct {
	Name string `arg:"positional"`
}

type ActivateCmd struct {
}

type ListCmd struct {
}

type WorkspaceCmd struct {
	Download *DownloadCmd `arg:"subcommand:download"`
	Set      *SetCmd      `arg:"subcommand:set"`
	Activate *ActivateCmd `arg:"subcommand:activate"`
	List     *ListCmd     `arg:"subcommand:list"`
	Name     string       `arg:"-n,--name" default:"default" help:"name of the workspace"`
}

type ExecCmd struct {
	Command string `arg:"positional" placeholder:"CMD" default:"echo $SPC_WORKSPACE" help:"command to run in SPC environment"`
}

type L1GethCmd struct {
	Flags map[string]string `arg:"positional"`
}

type SpGethCmd struct {
	Flags map[string]string `arg:"positional"`
}

type SpMagiCmd struct {
	Flags map[string]string `arg:"positional"`
}

type SidecarCmd struct {
	Flags map[string]string `arg:"positional"`
}

type UpCmd struct {
	L1Geth  *L1GethCmd  `arg:"subcommand:l1geth"`
	SpGeth  *SpGethCmd  `arg:"subcommand:spgeth"`
	SpMagi  *SpMagiCmd  `arg:"subcommand:spmagi"`
	Sidecar *SidecarCmd `arg:"subcommand:sidecar"`
}

type TestCmd struct {
	Fuzzer *FuzzerCmd `arg:"subcommand:fuzzer"`
}

type FuzzerCmd struct {
	NumTx int `arg:"--num-tx" default:"10"`
}

type Args struct {
	Workspace *WorkspaceCmd `arg:"subcommand:workspace"`
	Exec      *ExecCmd      `arg:"subcommand:exec"`
	Up        *UpCmd        `arg:"subcommand:up"`
	Test      *TestCmd      `arg:"subcommand:test"`
	LogLevel  string        `arg:"-v,--verbosity" help:"set the log level"`
}

type Config struct {
	Workspace *Workspace
	Args      Args
}

func (c *Config) Description() string {
	return "Specular CLI - toolkit for L2 integration and testing"
}

func (c *Config) GetLogLevel(defaultLevel logrus.Level) logrus.Level {
	level, err := logrus.ParseLevel(c.Args.LogLevel)
	if err != nil {
		level = defaultLevel
	}
	return level
}

func NewConfig() (*Config, error) {
	cfg := Config{}
	arg.MustParse(&cfg.Args)
	return &cfg, nil
}
