package launcher

import (
	"gopkg.in/urfave/cli.v1"
)

var (
	GenesisExportSections = cli.StringFlag{
		Name:  "export.sections",
		Usage: `Genesis sections to export separated by comma (e.g. "brs-1" or "ers" or "evm-2" or "fws")`,
		Value: "brs,ers,fws",
	}
	importCommand = cli.Command{
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
	}
	exportCommand = cli.Command{
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
    opera export events

Requires a first argument of the file to write to.
Optional second and third arguments control the first and
last epoch to write. If the file ends with .gz, the output will
be gzipped
`,
			},
			{
				Name:      "genesis",
				Usage:     "Export current state into a genesis file",
				ArgsUsage: "<filename or dry-run> [<epochFrom> <epochTo>] [--export.sections=brs,ers,fws]",
				Action:    exportGenesis,
				Flags: []cli.Flag{
					DataDirFlag,
					GenesisExportSections,
				},
				Description: `
    opera export genesis

Export current state into a genesis file.
Requires a first argument of the file to write to.
Optional second and third arguments control the first and
last epoch to write.
Pass dry-run instead of filename for calculation of hashes without exporting data.
`,
			},
		},
	}
)
