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
	app.Flags = []cli.Flag{
		DataDirFlag,
	}
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
						ModeFlag,
					},
					Description: "TBD",
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
						ModeFlag,
					},
					Description: "TBD",
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
				},
			},
		},
		{
			Name:     "compact",
			Usage:    "Compact all pebble databases",
			Action: compactDbs,
			Flags: []cli.Flag{
				DataDirFlag,
				CacheFlag,
			},
			Description: "TBD",
		},
		{
			Name:     "cli",
			Usage:    "Start an interactive JavaScript environment, attach to a node",
			ArgsUsage: "[endpoint]",
			Action: remoteConsole,
			Flags: []cli.Flag{
				DataDirFlag,
				JSpathFlag,
				PreloadJSFlag,
				ExecFlag,
			},
			Description: `
The Sonic console is an interactive shell for the JavaScript runtime environment
which exposes a node admin interface as well as the Dapp JavaScript API.
See https://github.com/ethereum/go-ethereum/wiki/JavaScript-Console.
This command allows to open a console attached to a running Sonic node.`,
		},
		{
			Name:      "import",
			Usage:     "Import a blockchain file",
			ArgsUsage: "<filename> (<filename 2> ... <filename N>)",
			Category:  "MISCELLANEOUS COMMANDS",
			Description: `
    opera import events

The import command imports events from an RLP-encoded files.
Events are fully verified.`,

			Subcommands: []cli.Command{
				{
					Action:    importEvents,
					Name:      "events",
					Usage:     "Import blockchain events",
					ArgsUsage: "<filename> (<filename 2> ... <filename N>)",
					Flags: []cli.Flag{
						DataDirFlag,
					},
					Description: `
The import command imports events from RLP-encoded files.
Events are fully verified.`,
				},
			},
		},
		{
			Name:     "export",
			Usage:    "Export blockchain",
			Category: "MISCELLANEOUS COMMANDS",

			Subcommands: []cli.Command{
				{
					Name:      "events",
					Usage:     "Export blockchain events",
					ArgsUsage: "<filename> [<epochFrom> <epochTo>]",
					Action:    exportEvents,
					Flags: []cli.Flag{
						DataDirFlag,
					},
					Description: `
    sonictool export events

Requires a first argument of the file to write to.
Optional second and third arguments control the first and
last epoch to write. If the file ends with .gz, the output will
be gzipped.
`,
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
