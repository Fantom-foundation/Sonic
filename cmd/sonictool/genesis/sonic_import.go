package genesis

import (
	"archive/tar"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/db"
	"github.com/Fantom-foundation/go-opera/opera/genesisstore/filelog"
	"github.com/ethereum/go-ethereum/log"
	"github.com/klauspost/pgzip"
	"io"
	"os"
	"path/filepath"
	"time"
)

func SonicImport(dataDir string, genesisFile *os.File) error {
	info, err := genesisFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get genesis file stats: %w", err)
	}

	if err := db.RemoveDatabase(dataDir); err != nil {
		return fmt.Errorf("failed to remove existing data from the datadir: %w", err)
	}

	chaindataDir := filepath.Join(dataDir, "chaindata")
	err = os.MkdirAll(chaindataDir, 0700)
	if err != nil {
		return fmt.Errorf("failed to create chaindataDir directory: %w", err)
	}
	setGenesisProcessing(chaindataDir)

	log.Info("Unpacking Sonic genesis")
	start := time.Now()
	hasher := sha256.New()
	reader := filelog.Wrap(genesisFile, "sonic-genesis", uint64(info.Size()), time.Minute)
	teeReader := io.TeeReader(reader, hasher)
	uncompressedStream, err := pgzip.NewReaderN(teeReader, 2621440, 64) // 30% faster than native gzip
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
	log.Info("Unpacking finished", "elapsed", time.Since(start))

	name, ok := allowedSonicGenesisHashes[hash]
	if !ok {
		_ = db.RemoveDatabase(dataDir)
		return fmt.Errorf("hash of the genesis file does not match any allowed value: %s", hash)
	}

	setGenesisComplete(chaindataDir)
	log.Info("Successfully imported Sonic genesis", "name", name)
	return nil
}
