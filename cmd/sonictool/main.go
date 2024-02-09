package main

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/flags"
	_ "github.com/Fantom-foundation/go-opera/version"
	"gopkg.in/urfave/cli.v1"
	"os"
	"sort"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags).
	gitCommit = ""
	gitDate   = ""
)

var (
	DataDirFlag = cli.StringFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
	}
)

func main() {
	app := flags.NewApp(gitCommit, gitDate, "the Sonic management tool")
	app.Commands = []cli.Command{
		{
			Name:     "genesis",
			Usage:    "Download or import genesis files",
			Description: "TBD",
			Subcommands: []cli.Command{
				{
					Name:   "sonic",
					Usage:  "Initialize the database from a tar.gz genesis file",
					Action: sonicGenesisImport,
					Flags: []cli.Flag{
						DataDirFlag,
						GenesisFlag,
					},
					Description: "TBD",
					CustomHelpTemplate: AppHelpTemplate,
				},
				{
					Name:   "legacy",
					Usage:  "Initialize the database from a legacy genesis file",
					Action: legacyGenesisImport,
					Flags: []cli.Flag{
						DataDirFlag,
						GenesisFlag,
						ExperimentalFlag,
						CacheFlag,
					},
					Description: "TBD",
					CustomHelpTemplate: AppHelpTemplate,
				},
				{
					Name:   "json",
					Usage:  "Initialize the database from a testing JSON genesis file",
					Action: jsonGenesisImport,
					Flags: []cli.Flag{
						DataDirFlag,
						GenesisFlag,
						ExperimentalFlag,
						CacheFlag,
					},
					Description: "TBD",
					CustomHelpTemplate: AppHelpTemplate,
				},
				{
					Name:   "fake",
					Usage:  "Initialize the database for a fakenet testing network",
					Action: fakeGenesisImport,
					Flags: []cli.Flag{
						DataDirFlag,
						FakeNetFlag,
						CacheFlag,
					},
					Description: "TBD",
					CustomHelpTemplate: AppHelpTemplate,
				},
			},
		},
		{
			Name:     "check",
			Usage:    "Check EVM database consistency",
			Description: "TBD",
			Subcommands: []cli.Command{
				{
					Name:   "live",
					Usage:  "Check EVM live state database",
					Action: checkLive,
					Flags: []cli.Flag{
						DataDirFlag,
						CacheFlag,
					},
					Description: "TBD",
					CustomHelpTemplate: AppHelpTemplate,
				},
				{
					Name:   "archive",
					Usage:  "Check EVM archive states database",
					Action: checkArchive,
					Flags: []cli.Flag{
						DataDirFlag,
						CacheFlag,
					},
					Description: "TBD",
					CustomHelpTemplate: AppHelpTemplate,
				},
			},
		},
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
