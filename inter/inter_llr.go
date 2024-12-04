package inter

import (
	"crypto/sha256"

	"github.com/Fantom-foundation/lachesis-base/common/bigendian"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
)

type LlrBlockVotes struct {
	Start idx.BlockID
	Epoch idx.EpochID
	Votes []hash.Hash
}

func (bvs LlrBlockVotes) LastBlock() idx.BlockID {
	return bvs.Start + idx.BlockID(len(bvs.Votes)) - 1
}

type LlrEpochVote struct {
	Epoch idx.EpochID
	Vote  hash.Hash
}

type LlrSignedBlockVotes struct {
	Signed                       SignedEventLocator
	TxsAndMisbehaviourProofsHash hash.Hash
	EpochVoteHash                hash.Hash
	Val                          LlrBlockVotes
}

type LlrSignedEpochVote struct {
	Signed                       SignedEventLocator
	TxsAndMisbehaviourProofsHash hash.Hash
	BlockVotesHash               hash.Hash
	Val                          LlrEpochVote
}

func (r SignedEventLocator) Size() uint64 {
	return uint64(len(r.Sig)) + 3*32 + 4*4
}

func (bvs LlrSignedBlockVotes) Size() uint64 {
	return bvs.Signed.Size() + uint64(len(bvs.Val.Votes))*32 + 32*2 + 8 + 4
}

func (ers LlrEpochVote) Hash() hash.Hash {
	hasher := sha256.New()
	hasher.Write(ers.Epoch.Bytes())
	hasher.Write(ers.Vote.Bytes())
	return hash.BytesToHash(hasher.Sum(nil))
}

func (bvs LlrBlockVotes) Hash() hash.Hash {
	hasher := sha256.New()
	hasher.Write(bvs.Start.Bytes())
	hasher.Write(bvs.Epoch.Bytes())
	hasher.Write(bigendian.Uint32ToBytes(uint32(len(bvs.Votes))))
	for _, bv := range bvs.Votes {
		hasher.Write(bv.Bytes())
	}
	return hash.BytesToHash(hasher.Sum(nil))
}

func (bvs LlrSignedBlockVotes) CalcPayloadHash() hash.Hash {
	return hash.Of(bvs.TxsAndMisbehaviourProofsHash.Bytes(), hash.Of(bvs.EpochVoteHash.Bytes(), bvs.Val.Hash().Bytes()).Bytes())
}

func (ev LlrSignedEpochVote) CalcPayloadHash() hash.Hash {
	return hash.Of(ev.TxsAndMisbehaviourProofsHash.Bytes(), hash.Of(ev.Val.Hash().Bytes(), ev.BlockVotesHash.Bytes()).Bytes())
}

func (ev LlrSignedEpochVote) Size() uint64 {
	return ev.Signed.Size() + 32 + 32*2 + 4 + 4
}
