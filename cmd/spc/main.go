package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/SpecularL2/specular-cli/internal/service/di"
)

func main() {
	application, _, err := di.SetupApplication()
	exitCode := 0
	if err != nil {
		log.Fatalf("failed to setup application #{err}")
	}
	if err := application.Run(); err != nil {
		if !errors.Is(err, context.Canceled) {
			application.GetLogger().WithError(err).Log(
				logrus.FatalLevel,
				"application failed",
			)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
