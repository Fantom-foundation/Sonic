package main

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/genesis"
	ogenesis "github.com/Fantom-foundation/go-opera/opera/genesis"
	"github.com/Fantom-foundation/go-opera/opera/genesisstore"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/term"
	"gopkg.in/urfave/cli.v1"
	"os"
	"syscall"
)

var (
	KeystoreFlag = &cli.StringFlag{
		Name:  "keystore",
		Usage: "Directory for the keystore",
	}
)

func signGenesis(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		return fmt.Errorf("this command requires an argument - the genesis file to import")
	}

	header, genesisHashes, err := getGenesisHeaderHashes(ctx.Args().First())
	if err != nil {
		return err
	}

	for sectionName, sectionHash := range genesisHashes {
		log.Info("Section", "name", sectionName, "hash", hexutil.Encode(sectionHash.Bytes()))
	}

	hash, err := genesis.CalculateHashFromGenesis(header, genesisHashes)
	if err != nil {
		return err
	}
	log.Info("Hash to sign", "hash", hexutil.Encode(hash.Bytes()))

	keystoreFilename := ctx.String(KeystoreFlag.Name)
	if keystoreFilename == "" {
		return fmt.Errorf("please specify the --%s flag", KeystoreFlag.Name)
	}
	keystoreJson, err := os.ReadFile(keystoreFilename)
	if err != nil {
		return err
	}

	fmt.Print("Keystore passphrase: ")
	passphrase, err := term.ReadPassword(syscall.Stdin)
	fmt.Printf("\n")
	if err != nil {
		return err
	}
	key, err := keystore.DecryptKey(keystoreJson, string(passphrase))
	if err != nil {
		return err
	}

	publicKeyBytes := crypto.FromECDSAPub(&key.PrivateKey.PublicKey)
	log.Info("Signing key opened", "pubkey", hexutil.Encode(publicKeyBytes))

	signature, err := crypto.Sign(hash.Bytes(), key.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to sign metadata: %w", err)
	}

	log.Info("Signed", "signature", hexutil.Encode(signature))

	if err := genesis.CheckGenesisSignature(hash, signature); err != nil {
		return err
	}
	if err = genesis.WriteSignatureIntoGenesisFile(header, signature, ctx.Args().First()); err != nil {
		return fmt.Errorf("failed to write signature into genesis file: %w", err)
	}
	log.Info("Signature successfully written into genesis file")
	return nil
}

func getGenesisHeaderHashes(genesisFile string) (ogenesis.Header, ogenesis.Hashes, error) {
	genesisReader, err := os.Open(genesisFile)
	if err != nil {
		return ogenesis.Header{}, nil, fmt.Errorf("failed to open the genesis file: %w", err)
	}
	defer genesisReader.Close()

	genesisStore, genesisHashes, err := genesisstore.OpenGenesisStore(genesisReader)
	if err != nil {
		return ogenesis.Header{}, nil, fmt.Errorf("failed to read genesis file: %w", err)
	}
	defer genesisStore.Close()

	return genesisStore.Header(), genesisHashes, nil
}
