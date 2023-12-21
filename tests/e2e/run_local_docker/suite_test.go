package run_local_docker

import (
	"context"
	"syscall"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/SpecularL2/specular-cli/internal/service/di"
	"github.com/SpecularL2/specular-cli/internal/spc/infra/integration/local_docker"
	"github.com/SpecularL2/specular-cli/tests"
)

type SpcIntegrationService struct {
	CancelFunc  context.CancelFunc
	AppInstance *di.TestApplication
	TestCluster *tests.IntegrationTestSetup
}

type LocalDockerTestSuite struct {
	suite.Suite
	service SpcIntegrationService
}

func (s *LocalDockerTestSuite) SetupTest() {
	s.service = initTestService(s.T())
}

func (s *LocalDockerTestSuite) TearDownTest() {
	s.service.TestCluster.Terminate(s.T())
	s.service.CancelFunc()
}

func TestSpcServiceSuit(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(LocalDockerTestSuite))
}

func initTestService(t *testing.T) SpcIntegrationService {
	ctx := context.TODO()
	log := logrus.New()

	cluster, err := local_docker.NewSpcDockerCluster(ctx, log, local_docker.WithAll())
	require.NoError(t, err)

	testSetup := tests.NewIntegrationTestSetup(t, ctx, cluster)

	app, cleanup, err := di.SetupApplicationForIntegrationTests(testSetup.Cfg)
	require.NoError(t, err)

	return SpcIntegrationService{
		CancelFunc:  cleanup,
		AppInstance: app,
		TestCluster: testSetup,
	}
}

func (s *LocalDockerTestSuite) startApplication() func() {
	s.T().Helper()

	go func() {
		err := s.service.AppInstance.Run()
		require.NoError(s.T(), err)
	}()

	return func() {
		err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		require.NoError(s.T(), err)
	}
}
