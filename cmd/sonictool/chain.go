package main

import (
	"compress/gzip"
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/chain"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	GenesisExportSections = cli.StringFlag{
		Name:  "export.sections",
		Usage: `Genesis sections to export separated by comma (e.g. "brs-1" or "ers" or "evm-2" or "fws")`,
		Value: "brs,ers,fws",
	}
)

func exportEvents(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	fn := ctx.Args().First()

	dataDir := ctx.String(DataDirFlag.Name)
	if dataDir == "" {
		return fmt.Errorf("--%s need to be set", DataDirFlag.Name)
	}

	// Open the file handle and potentially wrap with a gzip stream
	fh, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer fh.Close()

	var writer io.Writer = fh
	if strings.HasSuffix(fn, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}

	from := idx.Epoch(1)
	if len(ctx.Args()) > 1 {
		n, err := strconv.ParseUint(ctx.Args().Get(1), 10, 32)
		if err != nil {
			return err
		}
		from = idx.Epoch(n)
	}
	to := idx.Epoch(0)
	if len(ctx.Args()) > 2 {
		n, err := strconv.ParseUint(ctx.Args().Get(2), 10, 32)
		if err != nil {
			return err
		}
		to = idx.Epoch(n)
	}

	log.Info("Exporting events to file", "file", fn)
	err = chain.ExportEvents(writer, dataDir, from, to)
	if err != nil {
		utils.Fatalf("Export error: %v\n", err)
	}

	return nil
}

func importEvents(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	err := chain.EventsImport(ctx, ctx.Args()...)
	if err != nil {
		return err
	}

	return nil
}
