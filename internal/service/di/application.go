package di

import (
	"context"
	"os"
	"os/signal"

	"github.com/SpecularL2/specular-cli/internal/spc/executor"

	"github.com/SpecularL2/specular-cli/internal/spc/workspace"

	"github.com/sirupsen/logrus"

	"github.com/SpecularL2/specular-cli/internal/service/config"

	"golang.org/x/sync/errgroup"
)

type WaitGroup interface {
	Add(int)
	Done()
	Wait()
}

type Application struct {
	ctx    context.Context
	log    *logrus.Logger
	config *config.Config

	workspace *workspace.WorkspaceHandler
	executor  *executor.RunHandler
}

func (app *Application) Run() error {
	var _, cancel = context.WithCancel(app.ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	errGroup, _ := errgroup.WithContext(app.ctx)

	switch {
	case app.config.Args.Run != nil:
		return app.executor.Cmd()
	case app.config.Args.Workspace != nil:
		return app.workspace.Cmd()
	}

	err := errGroup.Wait()
	return err
}

func (app *Application) ShutdownAndCleanup() {
	app.log.Info("app shutting down")
}

func (app *Application) GetLogger() *logrus.Logger {
	return app.log
}

func (app *Application) GetContext() context.Context {
	return app.ctx
}

func (app *Application) GetConfig() *config.Config {
	return app.config
}

type TestApplication struct {
	*Application
	Ctx    context.Context
	Log    *logrus.Logger
	Config *config.Config

	workspace *workspace.WorkspaceHandler
}
