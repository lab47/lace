package main

import (
	"context"
	"os"

	"github.com/lab47/lace/pkg/build"
	"github.com/lab47/lsvd/logger"
)

func main() {
	log := logger.New(logger.Info)

	cwd, err := os.Getwd()
	if err != nil {
		log.Error("unable to calculate working directory", "error", err)
		os.Exit(1)
	}

	b, err := build.LoadBuilder(log, cwd)
	if err != nil {
		log.Error("unable to load build config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	err = b.Run(ctx)
	if err != nil {
		log.Error("error running build", "error", err)
		os.Exit(1)
	}
}
