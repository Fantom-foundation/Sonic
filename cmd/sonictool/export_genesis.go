package main

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/chain"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/db"
	"github.com/Fantom-foundation/go-opera/config/flags"
	"github.com/Fantom-foundation/go-opera/integration"
	"github.com/Fantom-foundation/lachesis-base/utils/cachescale"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	GenesisExportSections = cli.StringFlag{
		Name:  "export.sections",
		Usage: `Genesis sections to export separated by comma (e.g. "brs-1" or "ers" or "fws")`,
		Value: "brs,ers,fws,fwa",
	}
)

func exportGenesis(ctx *cli.Context) error {
	dataDir := ctx.GlobalString(flags.DataDirFlag.Name)
	if dataDir == "" {
		return fmt.Errorf("--%s need to be set", flags.DataDirFlag.Name)
	}
	fileName := ctx.Args().First()
	if fileName == "" {
		return fmt.Errorf("the output file name must be provided as an argument")
	}

	sectionsStr := ctx.String(GenesisExportSections.Name)
	sections := map[string]string{}
	for _, str := range strings.Split(sectionsStr, ",") {
		before := len(sections)
		if strings.HasPrefix(str, "brs") {
			sections["brs"] = str
		} else if strings.HasPrefix(str, "ers") {
			sections["ers"] = str
		} else if strings.HasPrefix(str, "fws") {
			sections["fws"] = str
		} else if strings.HasPrefix(str, "fwa") {
			sections["fwa"] = str
		} else {
			return fmt.Errorf("unknown section '%s': has to start with either 'brs' or 'ers' or 'fws' or `fwa`", str)
		}
		if len(sections) == before {
			return fmt.Errorf("duplicate section: '%s'", str)
		}
	}

	cacheRatio, err := cacheScaler(ctx)
	if err != nil {
		return err
	}
	chaindataDir := filepath.Join(dataDir, "chaindata")
	dbs, err := integration.GetDbProducer(chaindataDir, integration.DBCacheConfig{
		Cache:   cacheRatio.U64(480 * opt.MiB),
		Fdlimit: 100,
	})
	if err != nil {
		return fmt.Errorf("failed to make DB producer: %v", err)
	}
	defer dbs.Close()

	gdb, err := db.MakeGossipDb(dbs, dataDir, false, cachescale.Identity)
	if err != nil {
		return err
	}
	defer gdb.Close()

	fh, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer fh.Close()

	tmpPath := path.Join(dataDir, "tmp-genesis-export")
	_ = os.RemoveAll(tmpPath)
	defer os.RemoveAll(tmpPath)

	return chain.ExportGenesis(gdb, sections, fh, tmpPath)
}
