package genesis

import (
	"bytes"
	"fmt"
	"github.com/Fantom-foundation/go-opera/opera/genesis"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"os"
	"sort"
)

type SectionHash struct {
	Name string
	Hash hash.Hash
}

type Metadata struct {
	Header genesis.Header
	Hashes []SectionHash
}

func CalculateHashFromGenesis(header genesis.Header, genesisHashes genesis.Hashes) (common.Hash, error) {
	var metadata Metadata
	metadata.Header = header

	// add section hashes in deterministic order
	sectionNames := make(sort.StringSlice, 0, len(genesisHashes))
	for sectionName := range genesisHashes {
		if sectionName == "signature" {
			continue
		}
		sectionNames = append(sectionNames, sectionName)
	}
	sectionNames.Sort()

	for _, sectionName := range sectionNames {
		metadata.Hashes = append(metadata.Hashes, SectionHash{
			Name: sectionName,
			Hash: genesisHashes[sectionName],
		})
	}

	encodedMetadata, err := rlp.EncodeToBytes(metadata)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to RLP encode genesis metadata: %w", err)
	}
	return crypto.Keccak256Hash(encodedMetadata), nil
}

func CheckGenesisSignature(hash common.Hash, signature []byte) error {
	recoveredPubKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return err
	}
	for _, pubkey := range allowedPubkeys {
		if bytes.Equal(recoveredPubKey, pubkey) {
			return nil
		}
	}
	return fmt.Errorf("genesis signature does not match any trusted pubkey (pubkey: %x)", recoveredPubKey)
}

func WriteSignatureIntoGenesisFile(header genesis.Header, signature []byte, file string) error {
	out, err := os.OpenFile(file, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = out.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	defer out.Close()

	tmpDir, err := os.MkdirTemp("", "signing-genesis-tmp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	writer := newUnitWriter(out)
	if err := writer.Start(header, "signature", tmpDir); err != nil {
		return err
	}
	_, err = writer.Write(signature)
	if err != nil {
		return err
	}
	_, err = writer.Flush()
	return err
}
