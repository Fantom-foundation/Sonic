package genesis

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/Fantom-foundation/go-opera/opera/genesis"
	"github.com/Fantom-foundation/go-opera/opera/genesisstore"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"os"
)

type Metadata struct {
	Header genesis.Header
	Hashes []hash.Hash
}

func GetGenesisMetadataFromGenesisStore(genesisStore *genesisstore.Store, genesisHashes genesis.Hashes) *Metadata {
	g := genesisStore.Genesis()
	var metadata Metadata
	metadata.Header = genesis.Header{
		GenesisID:   g.GenesisID,
		NetworkID:   g.NetworkID,
		NetworkName: g.NetworkName,
	}
	for _, sectionHash := range genesisHashes {
		metadata.Hashes = append(metadata.Hashes, sectionHash)
	}
	return &metadata
}

func SignMetadata(genesisHashes *Metadata, privateKey *ecdsa.PrivateKey) (*genesis.SignedMetadata, error) {
	encodedTemplate, err := rlp.EncodeToBytes(genesisHashes)
	if err != nil {
		return nil, fmt.Errorf("failed to RLP encode genesis metadata: %w", err)
	}
	hashesHash := crypto.Keccak256Hash(encodedTemplate)

	signature, err := crypto.Sign(hashesHash.Bytes(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed sign genesis metadata: %w", err)
	}
	return &genesis.SignedMetadata{
		Signature: signature,
		Hashes:    encodedTemplate,
	}, nil
}

func WriteSignedMetadataIntoGenesisFile(header genesis.Header, signedMetadata *genesis.SignedMetadata, out *os.File) error {
	tmpDir, err := os.MkdirTemp("", "signing-genesis-tmp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	writer := newUnitWriter(out)
	if err := writer.Start(header, "signature", tmpDir); err != nil {
		return err
	}
	if err := rlp.Encode(writer, signedMetadata); err != nil {
		return err
	}
	_, err = writer.Flush()
	return err
}
