package local_docker_run

import (
	"testing"
	"time"

	"github.com/SpecularL2/specular-cli/internal/spc"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func (s *LocalDockerTestSuite) SkipTestExample() {
	cfg := s.service.TestCluster.Cfg
	var err error
	defer s.service.CancelFunc()
	stopApp := s.startApplication()
	defer stopApp()

	fromAddress, err := s.service.TestCluster.GetAccount(0)
	require.NoError(s.T(), err)

	fromPrvKey, err := s.service.TestCluster.GetPrivateKey(0)
	require.NoError(s.T(), err)

	toAddress, err := s.service.TestCluster.GetAccount(1)
	require.NoError(s.T(), err)

	client, err := ethclient.Dial(cfg.Workspace.L1GethURL)
	require.NoError(s.T(), err)

	s.T().Run("send transaction", func(t *testing.T) {
		// TODO: make so that the TestClusterSetup already has built in account funding (for the accounts existing),
		//       - we should abstract accounts to be dynamic and automatically funded
		//	     - we should make some labeling system instead of getting them by index :/

		// TODO: this will fail because the accounts are not yet funded
		fromBalance, err := spc.CheckBalance(client, fromAddress)
		require.NoError(s.T(), err)

		toBalance, err := spc.CheckBalance(client, toAddress)
		require.NoError(s.T(), err)

		require.Equal(s.T(), 0, fromBalance.Cmp(spc.Eth_1K), "from address incorrect initial balance")
		require.Equal(s.T(), 0, toBalance.Cmp(spc.Eth_1K), "from address incorrect initial balance")

		txSigned, err := spc.NewTransaction(client, fromPrvKey, spc.Mwei, toAddress)
		require.NoError(t, err)

		_, err = spc.SendTransaction(client, txSigned)
		require.NoError(t, err)

		fromBalance2, err := spc.CheckBalance(client, fromAddress)
		require.NoError(s.T(), err)

		// TODO: manually call mine or make sure it's mined
		time.Sleep(4 * time.Second)

		require.Equal(s.T(), 0, fromBalance2.Cmp(spc.Eth_1K), "from address incorrect initial balance")
	})
}
