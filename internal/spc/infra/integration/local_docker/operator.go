package local_docker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/SpecularL2/specular-cli/internal/service/config"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/samber/lo"
	"github.com/testcontainers/testcontainers-go"

	"github.com/SpecularL2/specular-cli/internal/spc/infra/integration"
)

var (
	_, _file, _, _  = runtime.Caller(0)                           // nolint:gochecknoglobals
	ProjectRoot     = filepath.Join(filepath.Dir(_file), "../..") // nolint:gochecknoglobals
	ReuseContainers = os.Getenv("REUSE_CONTAINERS") == "1"        // nolint:gochecknoglobals
)

const (
	ContainerContextTimeout = 1 * time.Minute
)

func CreateDockerNetwork(ctx context.Context, name *string) (*testcontainers.DockerNetwork, error) {
	if name == nil {
		if ReuseContainers {
			name = lo.ToPtr("myapp-test")
		} else {
			name = lo.ToPtr(uuid.NewString())
		}
	}
	// t.Logf("docker network name: %s", lo.FromPtr(name))

	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{Name: *name},
	})
	if err != nil {
		return nil, err
	}

	return network.(*testcontainers.DockerNetwork), nil
}

type SpcDockerCluster struct {
	ctx     context.Context
	log     *logrus.Logger
	network *testcontainers.DockerNetwork

	Geth    integration.ServerInstance
	SpGeth  integration.ServerInstance
	Magi    integration.ServerInstance
	Sidecar integration.ServerInstance

	Teardown func()
}

func (sdc *SpcDockerCluster) Workspace() *config.Workspace {
	return &config.Workspace{
		L1GethURL: fmt.Sprintf("http://%s", lo.Must(sdc.Geth.Address())),
	}
}

func (sdc *SpcDockerCluster) Close() error {
	if sdc.Geth != nil {
		if err := sdc.Geth.Close(); err != nil {
			return err
		}
	}

	if sdc.SpGeth != nil {
		if err := sdc.SpGeth.Close(); err != nil {
			return err
		}
	}

	if sdc.Magi != nil {
		if err := sdc.Magi.Close(); err != nil {
			return err
		}
	}

	if sdc.Sidecar != nil {
		if err := sdc.Sidecar.Close(); err != nil {
			return err
		}
	}

	if sdc.network != nil {
		if err := sdc.network.Remove(context.Background()); err != nil {
			return err
		}
	}

	return nil
}

func WithAll() func(*SpcDockerCluster) error {
	return func(sdc *SpcDockerCluster) error {
		var err error
		sdc.Geth, err = NewGethServer(sdc.ctx, sdc.log, sdc.network)
		if err != nil {
			return err
		}

		//sdc.SpGeth, err = NewSpGethServer(sdc.ctx, sdc.log, sdc.network)
		//if err != nil {
		//	return err
		//}
		return nil
	}
}

func WithGeth() func(*SpcDockerCluster) error {
	return func(sdc *SpcDockerCluster) error {
		var err error
		sdc.Geth, err = NewGethServer(sdc.ctx, sdc.log, sdc.network)
		if err != nil {
			return err
		}
		return nil
	}
}

func WithSpGeth() func(*SpcDockerCluster) error {
	return func(sdc *SpcDockerCluster) error {
		var err error
		sdc.SpGeth, err = NewSpGethServer(sdc.ctx, sdc.log, sdc.network)
		if err != nil {
			return err
		}
		return nil
	}
}

func NewSpcDockerCluster(ctx context.Context, log *logrus.Logger, options ...func(*SpcDockerCluster) error) (*SpcDockerCluster, error) {
	network, err := CreateDockerNetwork(ctx, nil)
	if err != nil {
		return nil, err
	}

	tc := &SpcDockerCluster{
		ctx:     ctx,
		log:     log,
		network: network,
	}

	for _, o := range options {
		if err := o(tc); err != nil {
			return nil, err
		}
	}

	return tc, nil
}
