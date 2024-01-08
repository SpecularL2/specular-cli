package local_docker

import (
	"context"
	"fmt"
	"time"

	"github.com/SpecularL2/specular-cli/internal/service/config"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	GethContainerName  = "geth"
	GethHost           = "geth"
	GethPortHTTP       = "8545"
	GethPortTCP        = "8546"
	GethPortP2P_TCPUDP = "30303"
	GethPortP2P_UDP    = "30304"
	GethImage          = "ethereum/client-go:stable"
	NetworkId          = 1
)

type GethServer struct {
	log      *logrus.Logger
	Instance testcontainers.Container
	config   *config.GethConfig
}

func (g GethServer) Port() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	port, err := g.Instance.MappedPort(ctx, GethPortTCP)
	if err != nil {
		return 0, err
	}
	return port.Int(), nil
}

func (g GethServer) Close() error {
	if !ReuseContainers {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		if err := g.Instance.Terminate(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (g GethServer) Address() (string, error) {
	port, err := g.Port()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("127.0.0.1:%d", port), nil
}

func (g GethServer) Prep() error {
	address, err := g.Address()
	if err != nil {
		return err
	}
	g.log.Infof("geth: %s", address)
	return nil
}

func NewGethServer(
	ctx context.Context,
	log *logrus.Logger,
	dockerNetwork *testcontainers.DockerNetwork,
	// gc *config.GethConfig,
) (*GethServer, error) {
	ctx, cancel := context.WithTimeout(ctx, ContainerContextTimeout)
	defer cancel()

	// TODO: refactor that this is a constructor's argument
	gc := &config.GethConfig{
		GethPortHTTP: GethPortHTTP,
		GethPortTCP:  GethPortTCP,
	}

	request := testcontainers.ContainerRequest{
		Name:     lo.Ternary(ReuseContainers, GethContainerName, ""),
		Hostname: GethHost,
		Image:    GethImage,
		ExposedPorts: []string{
			fmt.Sprintf("%s/tcp", GethPortHTTP),
			fmt.Sprintf("%s/tcp", GethPortTCP),
			fmt.Sprintf("%s/tcp", GethPortP2P_TCPUDP),
			fmt.Sprintf("%s/udp", GethPortP2P_TCPUDP),
			fmt.Sprintf("%s/udp", GethPortP2P_UDP),
		},
		AutoRemove:     false,
		SkipReaper:     ReuseContainers,
		Env:            map[string]string{},
		Networks:       []string{dockerNetwork.Name},
		NetworkAliases: map[string][]string{dockerNetwork.Name: {GethHost}},
		Entrypoint:     append([]string{"geth"}, gc.Args()...),
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

	return &GethServer{
		log:      log,
		Instance: instance,
		config:   gc,
	}, nil
}
