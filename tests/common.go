package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/SpecularL2/specular-cli/internal/service/config"
	"github.com/SpecularL2/specular-cli/internal/spc"
	"github.com/SpecularL2/specular-cli/internal/spc/infra/integration"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

type IntegrationTestSetup struct {
	Cfg         *config.Config
	privateKeys []string
	accounts    []common.Address
	TestCluster integration.SpcCluster
}

func (i *IntegrationTestSetup) Terminate(t *testing.T) {
	err := i.TestCluster.Close()
	require.NoError(t, err)
}

func (i *IntegrationTestSetup) GetAccount(idx int) (common.Address, error) {
	if idx >= 0 && idx < len(i.privateKeys) {
		return i.accounts[idx], nil
	}
	return common.Address{}, fmt.Errorf("invalid index %d", idx)
}

func (i *IntegrationTestSetup) GetPrivateKey(idx int) (string, error) {
	if idx >= 0 && idx < len(i.privateKeys) {
		return i.privateKeys[idx], nil
	}
	return "", fmt.Errorf("invalid index %d", idx)
}

func NewIntegrationTestSetup(t *testing.T, ctx context.Context, cluster integration.SpcCluster) *IntegrationTestSetup {
	t.Helper()

	cfg := &config.Config{
		Workspace: cluster.Workspace(),
	}

	ts := IntegrationTestSetup{
		Cfg: cfg,
		privateKeys: []string{
			"12aaa41f8c755a8576fb67122e4d167bf10083f9381372755c522d778f4993e3",
			"f893af1ff2cc1a46915cf9ff7d038f5a6b3472c9bf92ac43f75088479a192c6f",
			"f893af1ff2cc1a46915cf9ff7d038f5a6b3472c9bf92ac43f75088479a192c6f",
		},
		TestCluster: cluster,
	}

	ts.accounts = make([]common.Address, len(ts.privateKeys))
	var err error

	for i, pk := range ts.privateKeys {
		ts.accounts[i], err = spc.PrivateKeyHexToAccountAddress(pk)
		require.NoError(t, err)
	}

	return &ts
}
