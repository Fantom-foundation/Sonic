package main

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/check"
	"gopkg.in/urfave/cli.v1"
)

func checkLive(ctx *cli.Context) error {
	dataDir := ctx.GlobalString(DataDirFlag.Name)
	if dataDir == "" {
		return fmt.Errorf("--%s need to be set", DataDirFlag.Name)
	}
	cacheRatio, err := cacheScaler(ctx)
	if err != nil {
		return err
	}

	return check.CheckLiveStateDb(dataDir, cacheRatio)
}

func checkArchive(ctx *cli.Context) error {
	dataDir := ctx.GlobalString(DataDirFlag.Name)
	if dataDir == "" {
		return fmt.Errorf("--%s need to be set", DataDirFlag.Name)
	}
	cacheRatio, err := cacheScaler(ctx)
	if err != nil {
		return err
	}

	return check.CheckArchiveStateDb(dataDir, cacheRatio)
}
