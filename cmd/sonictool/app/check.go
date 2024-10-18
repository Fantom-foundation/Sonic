package app

import (
	"context"
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/check"
	"github.com/Fantom-foundation/go-opera/config/flags"
	"gopkg.in/urfave/cli.v1"
	"os/signal"
	"syscall"
)

func checkLive(ctx *cli.Context) error {
	dataDir := ctx.GlobalString(flags.DataDirFlag.Name)
	if dataDir == "" {
		return fmt.Errorf("--%s need to be set", flags.DataDirFlag.Name)
	}
	cacheRatio, err := cacheScaler(ctx)
	if err != nil {
		return err
	}

	cancelCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return check.CheckLiveStateDb(cancelCtx, dataDir, cacheRatio)
}

func checkArchive(ctx *cli.Context) error {
	dataDir := ctx.GlobalString(flags.DataDirFlag.Name)
	if dataDir == "" {
		return fmt.Errorf("--%s need to be set", flags.DataDirFlag.Name)
	}
	cacheRatio, err := cacheScaler(ctx)
	if err != nil {
		return err
	}

	cancelCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return check.CheckArchiveStateDb(cancelCtx, dataDir, cacheRatio)
}
