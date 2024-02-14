package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/opera/launcher"
	"github.com/ethereum/go-ethereum/node"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-opera/inter/validatorpk"
	"github.com/Fantom-foundation/go-opera/valkeystore"
	"github.com/Fantom-foundation/go-opera/valkeystore/encryption"
)

// validatorKeyCreate creates a new validator key into the keystore defined by the CLI flags.
func validatorKeyCreate(ctx *cli.Context) error {
	cfg := launcher.MakeAllConfigs(ctx)
	utils.SetNodeConfig(ctx, &cfg.Node)

	password := getPassPhrase("Your new validator key is locked with a password. Please give a password. Do not forget this password.", true, 0, utils.MakePasswordList(ctx))

	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		utils.Fatalf("Failed to create account: %v", err)
	}
	privateKey := crypto.FromECDSA(privateKeyECDSA)
	publicKey := validatorpk.PubKey{
		Raw:  crypto.FromECDSAPub(&privateKeyECDSA.PublicKey),
		Type: validatorpk.Types.Secp256k1,
	}

	valKeystore := valkeystore.NewDefaultFileRawKeystore(path.Join(getValKeystoreDir(cfg.Node), "validator"))
	err = valKeystore.Add(publicKey, privateKey, password)
	if err != nil {
		utils.Fatalf("Failed to create account: %v", err)
	}

	// Sanity check
	_, err = valKeystore.Get(publicKey, password)
	if err != nil {
		utils.Fatalf("Failed to decrypt the account: %v", err)
	}

	fmt.Printf("\nYour new key was generated\n\n")
	fmt.Printf("Public key:                  %s\n", publicKey.String())
	fmt.Printf("Public address of the key:   %s\n", crypto.PubkeyToAddress(privateKeyECDSA.PublicKey))
	fmt.Printf("Path of the secret key file: %s\n\n", valKeystore.PathOf(publicKey))
	fmt.Printf("- You can share your public key with anyone. Others need it to validate messages from you.\n")
	fmt.Printf("- You must NEVER share the secret key with anyone! The key controls access to your validator!\n")
	fmt.Printf("- You must BACKUP your key file! Without the key, it's impossible to operate the validator!\n")
	fmt.Printf("- You must REMEMBER your password! Without the password, it's impossible to decrypt the key!\n\n")
	return nil
}

// validatorKeyConvert converts account key to validator key.
func validatorKeyConvert(ctx *cli.Context) error {
	if len(ctx.Args()) < 2 {
		utils.Fatalf("This command requires 2 arguments.")
	}
	cfg := launcher.MakeAllConfigs(ctx)
	utils.SetNodeConfig(ctx, &cfg.Node)

	_, _, keydir, _ := cfg.Node.AccountConfig()

	pubkeyStr := ctx.Args().Get(1)
	pubkey, err := validatorpk.FromString(pubkeyStr)
	if err != nil {
		utils.Fatalf("Failed to decode the validator pubkey: %v", err)
	}

	var acckeypath string
	if strings.HasPrefix(ctx.Args().First(), "0x") {
		acckeypath, err = findAccountKeypath(common.HexToAddress(ctx.Args().First()), keydir)
		if err != nil {
			utils.Fatalf("Failed to find the account: %v", err)
		}
	} else {
		acckeypath = ctx.Args().First()
	}

	valkeypath := path.Join(keydir, "validator", common.Bytes2Hex(pubkey.Bytes()))
	err = encryption.MigrateAccountToValidatorKey(acckeypath, valkeypath, pubkey)
	if err != nil {
		utils.Fatalf("Failed to migrate the account key: %v", err)
	}
	fmt.Println("\nYour key was converted and saved to " + valkeypath)
	return nil
}

func findAccountKeypath(addr common.Address, keydir string) (keypath string, err error) {
	addrStr := strings.ToLower(addr.String())[2:]
	// find key path
	err = filepath.Walk(keydir, func(walk string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		_, filename := filepath.Split(walk)
		if strings.Contains(strings.ToLower(filename), addrStr) {
			keypath = walk
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return keypath, err
	}
	if len(keypath) == 0 {
		return keypath, errors.New("account not found")
	}
	return keypath, nil
}

func getValKeystoreDir(cfg node.Config) string {
	_, _, keydir, err := cfg.AccountConfig()
	if err != nil {
		utils.Fatalf("Failed to setup account config: %v", err)
	}
	return keydir
}
