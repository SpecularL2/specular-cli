package local_docker

import (
	"context"
	"fmt"
	"github.com/SpecularL2/specular-cli/internal/service/config"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	SpMagiHost          = "sp-magi"
	SpMagiImage         = "specular/sp-magi:latest"
	SpMagiContainerName = "sp-magi"
)

type SpMagiServer struct {
	log      *logrus.Logger
	Instance testcontainers.Container
}

func (s SpMagiServer) Port() (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s SpMagiServer) Close() error {
	//TODO implement me
	panic("implement me")
}

func (s SpMagiServer) Address() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SpMagiServer) Prep() error {
	//TODO implement me
	panic("implement me")
}

func NewSpMagiServer(
	ctx context.Context,
	log *logrus.Logger,
	dockerNetwork *testcontainers.DockerNetwork,
	spMagiConfig *config.SpMagiConfig,
) (*SpMagiServer, error) {
	ctx, cancel := context.WithTimeout(ctx, ContainerContextTimeout)
	defer cancel()

	request := testcontainers.ContainerRequest{
		Name:     lo.Ternary(ReuseContainers, SpMagiContainerName, ""),
		Hostname: SpMagiHost,
		Image:    SpMagiImage,
		ExposedPorts: []string{
			fmt.Sprintf("%s/tcp", GethPortHTTP),
		},
		AutoRemove:     false,
		SkipReaper:     ReuseContainers,
		Env:            map[string]string{},
		Networks:       []string{dockerNetwork.Name},
		NetworkAliases: map[string][]string{dockerNetwork.Name: {GethHost}},
		Entrypoint: []string{
			"geth",
			"--devnet",
			"--network", spMagiConfig.Network,
			"--l1-rpc-url", spMagiConfig.L1RpcURL,
			"--l2-rpc-url", spMagiConfig.L2RpcURL,
			"--sync-mode", spMagiConfig.SyncMode,
			"--l2-engine-url", spMagiConfig.L2EngineURL,
			"--jwt-file", spMagiConfig.JWTSecretPath,
			"--rpc-port", spMagiConfig.RpcPort,
			"--sequencer",
			"--sequencer-max-safe-lag", spMagiConfig.SequencerMaxSafeLag,
			"--sequencer-pk-file", spMagiConfig.SequencerPkFile,
			"--checkpoint-sync-url", spMagiConfig.CheckpointSyncUrl,
			"--checkpoint-hash", spMagiConfig.CheckpointHash,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("HTTP server started"),
			wait.ForListeningPort(GethPortHTTP),
			wait.ForListeningPort(GethPortTCP),
		),
	}
	instance, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
		Reuse:            ReuseContainers,
	})
	if err != nil {
		return nil, err
	}

	return &SpMagiServer{
		log:      log,
		Instance: instance,
	}, nil
}
