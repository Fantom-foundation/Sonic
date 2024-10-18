package app

import (
	"github.com/Fantom-foundation/go-opera/config"
	"github.com/Fantom-foundation/go-opera/version"
	"github.com/ethereum/go-ethereum/crypto"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Fantom-foundation/go-opera/integration/makefakegenesis"
	"github.com/Fantom-foundation/go-opera/inter/validatorpk"
)

const (
	ipcAPIs  = "abft:1.0 admin:1.0 dag:1.0 debug:1.0 ftm:1.0 net:1.0 personal:1.0 rpc:1.0 trace:1.0 txpool:1.0 web3:1.0"
)

func TestFakeNetFlag_NonValidator(t *testing.T) {
	// Start an opera console, make sure it's cleaned up and terminate the console
	dataDir := tmpdir(t)
	initFakenetDatadir(dataDir, 3)
	cli := exec(t,
		"--fakenet", "0/3",
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none", "--datadir", dataDir,
		"--exitwhensynced.epoch", "1")

	// Gather all the infos the welcome message needs to contain
	cli.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	cli.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	cli.SetTemplateFunc("gover", runtime.Version)
	cli.SetTemplateFunc("version", func() string { return version.VersionWithCommit("", "") })
	cli.SetTemplateFunc("niltime", genesisStart)
	cli.SetTemplateFunc("apis", func() string { return ipcAPIs })
	cli.ExpectExit()

	wantMessages := []string{
		"Unlocked fake validator",
	}
	for _, m := range wantMessages {
		if strings.Contains(cli.StderrText(), m) {
			t.Errorf("stderr text contains %q", m)
		}
	}
}

func TestFakeNetFlag_Validator(t *testing.T) {
	// Start an opera console, make sure it's cleaned up and terminate the console
	dataDir := tmpdir(t)
	initFakenetDatadir(dataDir, 3)
	cli := exec(t,
		"--fakenet", "3/3",
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none", "--datadir", dataDir,
		"--exitwhensynced.epoch", "1")

	// Gather all the infos the welcome message needs to contain
	va := readFakeValidator("3/3")
	cli.Coinbase = "0x0000000000000000000000000000000000000000"
	cli.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	cli.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	cli.SetTemplateFunc("gover", runtime.Version)
	cli.SetTemplateFunc("version", func() string { return version.VersionWithCommit("", "") })
	cli.SetTemplateFunc("niltime", genesisStart)
	cli.SetTemplateFunc("apis", func() string { return ipcAPIs })
	cli.ExpectExit()

	wantMessages := []string{
		"Unlocked validator key",
		"pubkey=" + va.String(),
	}
	for _, m := range wantMessages {
		if !strings.Contains(cli.StderrText(), m) {
			t.Errorf("stderr text does not contain %q", m)
		}
	}
}

func readFakeValidator(fakenet string) *validatorpk.PubKey {
	n, _, err := config.ParseFakeGen(fakenet)
	if err != nil {
		panic(err)
	}

	if n < 1 {
		return nil
	}

	return &validatorpk.PubKey{
		Raw:  crypto.FromECDSAPub(&makefakegenesis.FakeKey(n).PublicKey),
		Type: validatorpk.Types.Secp256k1,
	}
}

func genesisStart() string {
	return time.Unix(int64(makefakegenesis.FakeGenesisTime.Unix()), 0).Format("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)")
}
