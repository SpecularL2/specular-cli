package local_docker

import (
	"context"
	"fmt"
	"time"

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
) (*GethServer, error) {
	ctx, cancel := context.WithTimeout(ctx, ContainerContextTimeout)
	defer cancel()

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
		Entrypoint: []string{
			"geth",
			"--dev",
			"--http",
			"--http.addr", "0.0.0.0",
			"--http.port", GethPortHTTP,
			"--http.api", "\"engine,personal,eth,net,web3,txpool,miner,debug\"",
			"--http.corsdomain", "\"*\"",
			"--http.vhosts", "\"*\"",
			"--ws",
			"--ws.addr", "0.0.0.0",
			"--ws.port", GethPortTCP,
			"--ws.api", "engine,personal,eth,net,web3,txpool,miner,debug",
			"--ws.origins", "\"*\"",
			"--authrpc.vhosts", "*",
			"--authrpc.addr", "0.0.0.0",
			// "--authrpc.port", "$AUTH_PORT",
			// "--authrpc.jwtsecret", "$JWT_SECRET_PATH",
			"--miner.recommit", "0",
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

	return &GethServer{
		log:      log,
		Instance: instance,
	}, nil
}
