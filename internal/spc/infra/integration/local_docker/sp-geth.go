package local_docker

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
)

type SpGethServer struct {
}

func (s SpGethServer) Port() (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s SpGethServer) Close() error {
	//TODO implement me
	panic("implement me")
}

func (s SpGethServer) Address() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SpGethServer) Prep() error {
	//TODO implement me
	panic("implement me")
}

func NewSpGethServer(
	_ context.Context,
	_ *logrus.Logger,
	_ *testcontainers.DockerNetwork,
) (*SpGethServer, error) {
	//TODO implement me
	panic("implement me")
}
