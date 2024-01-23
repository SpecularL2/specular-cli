package up

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/SpecularL2/specular-cli/internal/service/config"
	"github.com/SpecularL2/specular-cli/internal/spc/handlers/workspace"
)

type UpHandler struct {
	cfg       *config.Config
	log       *logrus.Logger
	workspace *workspace.WorkspaceHandler
}

func (u *UpHandler) Cmd() error {
	switch {
	case u.cfg.Args.Up.SpGeth != nil:
		return u.StartSpGeth()
	case u.cfg.Args.Up.L1Geth != nil:
		return u.StartL1Geth()
	case u.cfg.Args.Up.SpMagi != nil:
		return u.StartSpMagi()
	case u.cfg.Args.Up.Sidecar != nil:
		return u.StartSidecar()
	}
	u.log.Warn("no command found, exiting...")
	return nil
}

func (u *UpHandler) StartSpGeth() error {
	// TODO: implement overriding flags
	u.log.Warn("NOT IMPLEMENT - overidden flags:", u.cfg.Args.Up.SpGeth.Flags)

	// TODO:
	//	- all of the flag values should be changable
	//	- inject values directly instead of loading via env?
	spGethCommand := ".$SPC_SP_GETH_BIN " +
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
		"--syncmode full"

	cmd, err := u.workspace.RunStringCommand(spGethCommand)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func (u *UpHandler) StartL1Geth() error {
	// TODO: implement overriding flags
	u.log.Warn("NOT IMPLEMENT - overidden flags:", u.cfg.Args.Up.L1Geth.Flags)

	period := 3
	value, ok := os.LookupEnv("L1_PERIOD")
	if ok {
		parsedVal, err := strconv.Atoi(value)
		period = parsedVal
		if err != nil {
			return fmt.Errorf("invalid L1_PERIOD: %s", err)
		}
	} else {
		os.Setenv("L1_PERIOD", fmt.Sprint(period))
	}
	u.log.Debugf("set block time to %ds", period)

	// TODO: try to read this from $SPC_L1_ENDPOINT which is already set in ENV
	if _, ok := os.LookupEnv("L1_PORT"); !ok {
		os.Setenv("L1_PORT", "8545")
	}

	l1GethCommand := ".$SPC_L1_GETH_BIN " +
		"--dev " +
		"--dev.period $L1_PERIOD " +
		"--verbosity 3 " +
		"--http " +
		"--http.api eth,web3,net " +
		"--http.addr 0.0.0.0 " +
		"--ws " +
		"--ws.api eth,net,web3 " +
		"--ws.addr 0.0.0.0 " +
		"--ws.port $L1_PORT"

	u.log.Info("starting L1 geth")
	cmd, err := u.workspace.RunStringCommand(l1GethCommand)
	if err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

// TODO
// func (u *UpHandler) fundL1Accounts() error {
// 	u.log.Info("funding accounts")
//
// 	var addressesToFund = []string{
// 		"$SPC_SEQUENCER_ADDRESS",
// 		"$SPC_VALIDATOR_ADDRESS",
// 		"$SPC_DEPLOYER_ADDRESS",
// 	}
//
// 	for _, addressName := range addressesToFund {
// 		address, ok := os.LookupEnv(addressName)
// 		if !ok {
// 			return fmt.Errorf("could not get address: %s", addressName)
// 		}
// 		u.log.Infof("funding %s", address)
// 	}
// 	return nil
// }

func (u *UpHandler) StartSpMagi() error {
	// TODO: implement overriding flags
	u.log.Warn("NOT IMPLEMENT - overidden flags:", u.cfg.Args.Up.SpMagi.Flags)

	// TODO: handle sync, devnet, sequencer settings here, not in sbin
	spMagiCommand := ".$SPC_SP_MAGI_BIN" +
		"--network $SPC_NETWORK " +
		"--l1-rpc-url $SPC_L1_RPC_URL " +
		"--l2-rpc-url $SPC_L2_RPC_URL " +
		"--sync-mode $SPC_SYNC_MODE " +
		"--l2-engine-url $SPC_L2_ENGINE_URL " +
		"--jwt-file $SPC_JWT_SECRET_PATH " +
		"--rpc-port $SPC_RPC_PORT " +
		"$SYNC_FLAGS $DEVNET_FLAGS $SEQUENCER_FLAGS $@"

	cmd, err := u.workspace.RunStringCommand(spMagiCommand)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func (u *UpHandler) StartSidecar() error {
	// TODO: implement overriding flags
	u.log.Warn("NOT IMPLEMENT - overidden flags:", u.cfg.Args.Up.SpMagi.Flags)

	// TODO: easily toggle disseminator & toggle
	sidecarCommand := ".$SPC_SIDECAR_BIN" +
		"--l1.endpoint $SPC_L1_ENDPOINT" +
		"--l2.endpoint $SPC_L2_ENDPOINT" +
		"--protocol.rollup-cfg-path $SPC_ROLLUP_CFG_PATH" +
		"--disseminator" +
		"--disseminator.private-key $SPC_DISSEMINATOR_PRIV_KEY" +
		"--disseminator.sub-safety-margin $SPC_DISSEMINATOR_SUB_SAFETY_MARGIN" +
		"--disseminator.target-batch-size $SPC_DISSEMINATOR_TARGET_BATCH_SIZE" +
		"--validator" +
		"--validator.private-key $SPC_VALIDATOR_PRIV_KEY"

	cmd, err := u.workspace.RunStringCommand(sidecarCommand)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func NewUpHandler(cfg *config.Config, log *logrus.Logger, workspace *workspace.WorkspaceHandler) *UpHandler {
	return &UpHandler{
		cfg:       cfg,
		log:       log,
		workspace: workspace,
	}
}
