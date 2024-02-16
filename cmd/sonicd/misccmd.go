package main

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/config"
	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/urfave/cli.v1"
	"os"
	"runtime"

	"github.com/Fantom-foundation/go-opera/gossip"
)

var (
	versionCommand = cli.Command{
		Action:    version,
		Name:      "version",
		Usage:     "Print version numbers",
		ArgsUsage: " ",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
The output of this command is supposed to be machine-readable.
`,
	}
)

func version(ctx *cli.Context) error {
	fmt.Println(config.ClientIdentifier)
	fmt.Println("Version:", params.VersionWithMeta())
	if config.GitCommit != "" {
		fmt.Println("Git Commit:", config.GitCommit)
	}
	if config.GitDate != "" {
		fmt.Println("Git Commit Date:", config.GitDate)
	}
	fmt.Println("Architecture:", runtime.GOARCH)
	fmt.Println("Protocol Versions:", []uint{gossip.ProtocolVersion})
	fmt.Println("Go Version:", runtime.Version())
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Printf("GOPATH=%s\n", os.Getenv("GOPATH"))
	fmt.Printf("GOROOT=%s\n", runtime.GOROOT())
	return nil
}
