package main

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/chain"
	"github.com/Fantom-foundation/go-opera/config/flags"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/log"
	gzip "github.com/klauspost/pgzip"
	"gopkg.in/urfave/cli.v1"
	"io"
	"os"
	"strconv"
	"strings"
)

func exportEvents(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		return fmt.Errorf("this command requires an argument - the output file")
	}

	fn := ctx.Args().First()

	dataDir := ctx.GlobalString(flags.DataDirFlag.Name)
	if dataDir == "" {
		return fmt.Errorf("--%s need to be set", flags.DataDirFlag.Name)
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
		return fmt.Errorf("export error: %w\n", err)
	}

	return nil
}

func importEvents(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		return fmt.Errorf("this command requires an argument - the input file")
	}

	err := chain.EventsImport(ctx, ctx.Args()...)
	if err != nil {
		return err
	}

	return nil
}
