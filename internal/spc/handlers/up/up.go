package up

import (
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/SpecularL2/specular-cli/internal/service/config"
	"github.com/SpecularL2/specular-cli/internal/spc/handlers/workspace"
)

type UpHandler struct {
	cfg *config.Config
	log *logrus.Logger
	workspace *workspace.WorkspaceHandler
}

func (u *UpHandler) Cmd() error {
	switch {
	case u.cfg.Args.Up.SpGeth != nil:
		return u.StartSpGeth()
	case u.cfg.Args.Up.L1Geth != nil:
		return u.StartL1Geth()
	}
	u.log.Warn("no command found, exiting...")
	return nil
}

func (u *UpHandler) StartSpGeth() error {
	u.log.Warn("NOT IMPLEMENT - overidden flags:", u.cfg.Args.Up.SpGeth.Flags)

	err := u.workspace.LoadWorkspaceEnvVars()
	if err != nil {
		return err
	}

	// TODO: 
	//	- all of the flag values should be changable
	//	- inject values directly instead of loading via env? 
	spGeth := ".$SPC_SP_GETH_BIN " +
	"--datadir $SPC_DATA_DIR " +
	"--networkid $SPC_NETWORK_ID " +
	"--http " +
	"--http.addr $SPC_ADDRESS " +
	"--http.port $SPC_HTTP_PORT " +
	"--http.api engine,personal,eth,net,web3,txpool,miner,debug " +
	"--http.corsdomain=* " +
	"--http.vhosts=* " +
	"--ws " +
	"--ws.addr $SPC_ADDRESS " +
	"--ws.port $SPC_WS_PORT " +
	"--ws.api engine,personal,eth,net,web3,txpool,miner,debug " +
	"--ws.origins=* " +
	"--authrpc.vhosts=* " +
	"--authrpc.addr $SPC_ADDRESS " +
	"--authrpc.port $SPC_AUTH_PORT " +
	"--authrpc.jwtsecret $SPC_JWT_SECRET_PATH " +
	"--miner.recommit 0 " +
	"--nodiscover " +
	"--maxpeers 0 " +
	"--syncmode full "

	commandToRun := os.ExpandEnv(spGeth)
	args := strings.Fields(commandToRun)
	u.log.Infof("Running: %s", args)

	if len(args) > 0 {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			u.log.Warn("Up failed with:", err)
			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			}
		}
	}
	return nil
}

func (u *UpHandler) StartL1Geth() error {
	u.log.Warn("NOT IMPLEMENT - overidden flags:", u.cfg.Args.Up.SpGeth.Flags)

	err := u.workspace.LoadWorkspaceEnvVars()
	if err != nil {
		return err
	}

	// TODO: 
	//	- all of the flag values should be changable
	//	- inject values directly instead of loading via env? 
	//	- save L1 GETH config in workspace (currently it's in start L1 script
	spGeth := ".$SPC_L1_GETH_BIN " +
	"--dev " +
	"--dev.period $L1_PERIOD " +
	"--verbosity 0 " +
	"--http " +
	"--http.api eth,web3,net " +
	"--http.addr 0.0.0.0 " +
	"--ws " +
	"--ws.api eth,net,web3 " +
	"--ws.addr 0.0.0.0 " +
	"--ws.port $L1_PORT &>$LOG_FILE &"

	commandToRun := os.ExpandEnv(spGeth)
	args := strings.Fields(commandToRun)
	u.log.Infof("Running: %s", args)

	if len(args) > 0 {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			u.log.Warn("Up failed with:", err)
			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			}
		}
	}
	return nil
}

func NewUpHandler(cfg *config.Config, log *logrus.Logger, workspace *workspace.WorkspaceHandler) *UpHandler {
	return &UpHandler{
		cfg: cfg,
		log: log,
		workspace: workspace,
	}
}
