package main

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/db"
	"github.com/Fantom-foundation/go-opera/cmdhelper"
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

// TODO: reuse flags from opera/launcher/flags
var (
	DataDirFlag = cli.StringFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
	}
	KeyStoreDirFlag = cli.StringFlag{
		Name:  "keystore",
		Usage: "Directory for the keystore (default = inside the datadir)",
	}
	PasswordFileFlag = cli.StringFlag{
		Name:  "password",
		Usage: "Password file to use for non-interactive password input",
		Value: "",
	}
	LightKDFFlag = cli.BoolFlag{
		Name:  "lightkdf",
		Usage: "Reduce key-derivation RAM & CPU usage at some expense of KDF strength",
	}
	CacheFlag = cli.IntFlag{
		Name:  "cache",
		Usage: "Megabytes of memory allocated to internal pebble caching",
		Value: db.DefaultCacheSize,
	}
)

func main() {
	app := cmdhelper.NewApp(gitCommit, gitDate, "the Sonic management tool")
	app.Flags = []cli.Flag{
		DataDirFlag,
		CacheFlag,
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
					ArgsUsage: "<filename>",
					Action: sonicGenesisImport,
					Description: "TBD",
				},
				{
					Name:   "legacy",
					Usage:  "Initialize the database from a legacy genesis file",
					ArgsUsage: "<filename>",
					Action: legacyGenesisImport,
					Flags: []cli.Flag{
						ExperimentalFlag,
						ModeFlag,
					},
					Description: "TBD",
				},
				{
					Name:   "json",
					Usage:  "Initialize the database from a testing JSON genesis file",
					ArgsUsage: "<filename>",
					Action: jsonGenesisImport,
					Flags: []cli.Flag{
						ExperimentalFlag,
						ModeFlag,
					},
					Description: "TBD",
				},
				{
					Name:   "fake",
					Usage:  "Initialize the database for a fakenet testing network",
					ArgsUsage: "<validators>",
					Action: fakeGenesisImport,
					Flags: []cli.Flag{
						ModeFlag,
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
					Description: "TBD",
				},
				{
					Name:   "archive",
					Usage:  "Check EVM archive states database",
					Action: checkArchive,
					Description: "TBD",
				},
			},
		},
		{
			Name:     "compact",
			Usage:    "Compact all pebble databases",
			Action: compactDbs,
			Description: "TBD",
		},
		{
			Name:     "cli",
			Usage:    "Start an interactive JavaScript environment, attach to a node",
			ArgsUsage: "[endpoint]",
			Action: remoteConsole,
			Flags: []cli.Flag{
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
		{
			Action:      checkConfig,
			Name:        "checkconfig",
			Usage:       "Checks configuration file",
			ArgsUsage:   "<filename>",
			Category:    "MISCELLANEOUS COMMANDS",
			Description: `The checkconfig checks configuration file.`,
		},
		{
			Action:      dumpConfig,
			Name:        "dumpconfig",
			Usage:       "Show configuration values",
			ArgsUsage:   "<filename>",
			Category:    "MISCELLANEOUS COMMANDS",
			Description: `The dumpconfig command shows configuration values.`,
		},

		{
			Name:     "account",
			Usage:    "Manage accounts",
			Category: "ACCOUNT COMMANDS",
			Description: `

Manage accounts, list all existing accounts, import a private key into a new
account, create a new account or update an existing account.

It supports interactive mode, when you are prompted for password as well as
non-interactive mode where passwords are supplied via a given password file.
Non-interactive mode is only meant for scripted use on test networks or known
safe environments.

Make sure you remember the password you gave when creating a new account (with
either new or import). Without it you are not able to unlock your account.

Note that exporting your key in unencrypted format is NOT supported.

Keys are stored under <DATADIR>/keystore.
It is safe to transfer the entire directory or the individual keys therein
between ethereum nodes by simply copying.

Make sure you backup your keys regularly.`,
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Usage:  "Print summary of existing accounts",
					Action: accountList,
					Flags: []cli.Flag{
						DataDirFlag,
						KeyStoreDirFlag,
					},
					Description: `
Print a short summary of all accounts`,
				},
				{
					Name:   "new",
					Usage:  "Create a new account",
					Action: accountCreate,
					Flags: []cli.Flag{
						DataDirFlag,
						KeyStoreDirFlag,
						PasswordFileFlag,
						LightKDFFlag,
					},
					Description: `
    sonictool account new

Creates a new account and prints the address.

The account is saved in encrypted format, you are prompted for a passphrase.

You must remember this passphrase to unlock your account in the future.

For non-interactive use the passphrase can be specified with the --password flag:

Note, this is meant to be used for testing only, it is a bad idea to save your
password to file or expose in any other way.
`,
				},
				{
					Name:      "update",
					Usage:     "Update an existing account",
					Action:    accountUpdate,
					ArgsUsage: "<address>",
					Flags: []cli.Flag{
						DataDirFlag,
						KeyStoreDirFlag,
						LightKDFFlag,
					},
					Description: `
    sonictool account update <address>

Update an existing account.

The account is saved in the newest version in encrypted format, you are prompted
for a passphrase to unlock the account and another to save the updated file.

This same command can therefore be used to migrate an account of a deprecated
format to the newest format or change the password for an account.

For non-interactive use the passphrase can be specified with the --password flag:

    sonictool account update [options] <address>

Since only one password can be given, only format update can be performed,
changing your password is only possible interactively.
`,
				},
				{
					Name:   "import",
					Usage:  "Import a private key into a new account",
					Action: accountImport,
					Flags: []cli.Flag{
						DataDirFlag,
						KeyStoreDirFlag,
						PasswordFileFlag,
						LightKDFFlag,
					},
					ArgsUsage: "<keyFile>",
					Description: `
    sonictool account import <keyfile>

Imports an unencrypted private key from <keyfile> and creates a new account.
Prints the address.

The keyfile is assumed to contain an unencrypted private key in hexadecimal format.

The account is saved in encrypted format, you are prompted for a passphrase.

You must remember this passphrase to unlock your account in the future.

For non-interactive use the passphrase can be specified with the -password flag:

    sonictool account import [options] <keyfile>

Note:
As you can directly copy your encrypted accounts to another ethereum instance,
this import mechanism is not needed when you transfer an account between
nodes.
`,
				},
			},
		},

		{
			Name:     "validator",
			Usage:    "Manage validators",
			Category: "VALIDATOR COMMANDS",
			Description: `

Create a new validator private key.

It supports interactive mode, when you are prompted for password as well as
non-interactive mode where passwords are supplied via a given password file.
Non-interactive mode is only meant for scripted use on test networks or known
safe environments.

Make sure you remember the password you gave when creating a new validator key.
Without it you are not able to unlock your validator key.

Note that exporting your key in unencrypted format is NOT supported.

Keys are stored under <DATADIR>/keystore/validator.
It is safe to transfer the entire directory or the individual keys therein
between Opera nodes by simply copying.

Make sure you backup your keys regularly.`,
			Subcommands: []cli.Command{
				{
					Name:   "new",
					Usage:  "Create a new validator key",
					Action: validatorKeyCreate,
					Flags: []cli.Flag{
						DataDirFlag,
						KeyStoreDirFlag,
						PasswordFileFlag,
					},
					Description: `
    sonictool validator new

Creates a new validator private key and prints the public key.

The key is saved in encrypted format, you are prompted for a passphrase.

You must remember this passphrase to unlock your key in the future.

For non-interactive use the passphrase can be specified with the --validator.password flag:

Note, this is meant to be used for testing only, it is a bad idea to save your
password to file or expose in any other way.
`,
				},
				{
					Name:   "convert",
					Usage:  "Convert an account key to a validator key",
					Action: validatorKeyConvert,
					Flags: []cli.Flag{
						DataDirFlag,
						KeyStoreDirFlag,
					},
					ArgsUsage: "<account address> <validator pubkey>",
					Description: `
    sonictool validator convert

Converts an account private key to a validator private key and saves in the validator keystore.
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
