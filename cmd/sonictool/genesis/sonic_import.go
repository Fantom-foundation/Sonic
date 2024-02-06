package genesis

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"io"
	"os"
	"path/filepath"
)

func SonicImport(dataDir string, genesisReader *os.File) error {
	if err := removeDatabase(dataDir); err != nil {
		return fmt.Errorf("failed to remove existing data from the datadir: %w", err)
	}

	chaindataDir := filepath.Join(dataDir, "chaindata")
	err := os.MkdirAll(chaindataDir, 0700)
	if err != nil {
		return fmt.Errorf("failed to create chaindataDir directory: %w", err)
	}
	setGenesisProcessing(chaindataDir)

	log.Info("Unpacking Sonic genesis")
	hasher := sha256.New()
	teeReader := io.TeeReader(genesisReader, hasher)
	uncompressedStream, err := gzip.NewReader(teeReader)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("untar failed: %w", err)
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath.Join(dataDir, header.Name), 0700); err != nil {
				return fmt.Errorf("mkdir failed: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(filepath.Join(dataDir, header.Name))
			if err != nil {
				return fmt.Errorf("create file failed: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("file write failed: %w", err)
			}
			outFile.Close()
		default:
			return fmt.Errorf("uknown tar content type: %x in %s", header.Typeflag, header.Name)
		}
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	if !allowedSonicGenesisHashes[hash] {
		_ = removeDatabase(dataDir)
		return fmt.Errorf("hash of the genesis file does not match any allowed value: %s", hash)
	}

	setGenesisComplete(chaindataDir)
	log.Info("Successfully imported Sonic genesis")
	return nil
}
