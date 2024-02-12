package test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/rand"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"

	// TODO: import specular fork - requires changing go.mod of the fork to correct URL
	"github.com/MariusVanDerWijden/FuzzyVM/filler"
	txfuzz "github.com/MariusVanDerWijden/tx-fuzz"
	"github.com/MariusVanDerWijden/tx-fuzz/mutator"

	"github.com/SpecularL2/specular-cli/internal/service/config"
	"github.com/SpecularL2/specular-cli/internal/spc/handlers/workspace"
)

type TestHandler struct {
	cfg       *config.Config
	log       *logrus.Logger
	workspace *workspace.WorkspaceHandler
}

func (t *TestHandler) Cmd() error {
	switch {
	case t.cfg.Args.Test.Fuzzer != nil:
		return t.RunFuzzer(t.cfg.Args.Test.Fuzzer.NumTx)
	}

	t.log.Warn("no command found, exiting...")
	return nil
}

func (t *TestHandler) RunFuzzer(numTx int) error {
	t.log.Infof("sending %d transactions", numTx)

	err := t.workspace.LoadWorkspaceEnvVars()
	if err != nil {
		return fmt.Errorf("could not load active workspace: %s", err)
	}
	url := os.ExpandEnv("$SPC_L2_ENDPOINT")

	client, err := ethclient.Dial(url)
	if err != nil {
		return fmt.Errorf("could not dial client: %s", err)
	}
	defer client.Close()

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("could not geth chain ID: %s", err)
	}

	// TODO: manage accounts through spc
	privateKeyBuf, err := os.ReadFile(os.ExpandEnv("$SPC_DEPLOYER_PK_PATH"))
	if err != nil {
		return fmt.Errorf("could not read private key: %s", err)
	}
	privateKeyString := string(privateKeyBuf)[2:]

	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return fmt.Errorf("could not read private key: %s", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("could not cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// TODO: support fixed seed for reproducable tests
	source := rand.NewSource(rand.Int63())
	rng := rand.New(source)
	mut := mutator.NewMutator(rng)
	random := make([]byte, 10000)
	mut.FillBytes(&random)

	for i := 0; i < numTx; i++ {
		mut.MutateBytes(&random)
		filler := filler.NewFiller(random)

		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			return fmt.Errorf("could not get nonce: %s", err)
		}

		rndTx, err := txfuzz.RandomValidTx(client.Client(), filler, fromAddress, nonce, nil, nil, false)
		if err != nil {
			return fmt.Errorf("could not get random tx: %s", err)
		}

		signedTx, err := types.SignTx(rndTx, types.NewCancunSigner(chainId), privateKey)
		if err != nil {
			return fmt.Errorf("could not sign tx: %s", err)
		}

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			// valid transactions might still fail to get executed
			// for example if they are underpriced
			// the test will still fail if the sequencer crashes (since we won't get a nonce)
			t.log.Warnf("could not send tx: %s", err)
		}
	}

	return nil
}

func NewTestHandler(cfg *config.Config, log *logrus.Logger, workspace *workspace.WorkspaceHandler) *TestHandler {
	return &TestHandler{
		cfg:       cfg,
		log:       log,
		workspace: workspace,
	}
}
