package main

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/check"
	"github.com/Fantom-foundation/go-opera/config/flags"
	"gopkg.in/urfave/cli.v1"
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

	return check.CheckLiveStateDb(dataDir, cacheRatio)
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

	return check.CheckArchiveStateDb(dataDir, cacheRatio)
}
