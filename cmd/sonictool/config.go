package main

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/opera/launcher"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func checkConfig(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		return fmt.Errorf("this command requires an argument - the config toml file")
	}
	configFile := ctx.Args().Get(0)
	_, err := launcher.MayMakeAllConfigs(ctx, configFile)
	return err
}

// dumpConfig is the dumpconfig command.
func dumpConfig(ctx *cli.Context) error {
	cfg := launcher.MakeAllConfigs(ctx)
	comment := ""

	out, err := launcher.TomlSettings.Marshal(&cfg)
	if err != nil {
		return err
	}

	dump := os.Stdout
	if ctx.NArg() > 0 {
		dump, err = os.OpenFile(ctx.Args().Get(0), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer dump.Close()
	}
	dump.WriteString(comment)
	dump.Write(out)

	return nil
}
