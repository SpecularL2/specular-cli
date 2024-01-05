// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/SpecularL2/specular-cli/internal/service/config"
	"github.com/SpecularL2/specular-cli/internal/spc/executor"
	"github.com/SpecularL2/specular-cli/internal/spc/workspace"
)

// Injectors from inject.go:

func SetupApplication() (*Application, func(), error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, nil, err
	}
	logger := config.NewLogger(configConfig)
	cancelChannel := config.NewCancelChannel()
	context := config.NewContext(logger, cancelChannel)
	workspaceHandler := workspace.NewWorkspaceHandler(configConfig, logger)
	runHandler := executor.NewRunHandler(configConfig, logger)
	application := &Application{
		ctx:       context,
		log:       logger,
		config:    configConfig,
		workspace: workspaceHandler,
		executor:  runHandler,
	}
	return application, func() {
	}, nil
}

func SetupApplicationForIntegrationTests(cfg *config.Config) (*TestApplication, func(), error) {
	logger := config.NewLogger(cfg)
	cancelChannel := config.NewCancelChannel()
	context := config.NewContext(logger, cancelChannel)
	workspaceHandler := workspace.NewWorkspaceHandler(cfg, logger)
	runHandler := executor.NewRunHandler(cfg, logger)
	application := &Application{
		ctx:       context,
		log:       logger,
		config:    cfg,
		workspace: workspaceHandler,
		executor:  runHandler,
	}
	testApplication := &TestApplication{
		Application: application,
		Ctx:         context,
		Log:         logger,
		Config:      cfg,
		workspace:   workspaceHandler,
	}
	return testApplication, func() {
	}, nil
}
