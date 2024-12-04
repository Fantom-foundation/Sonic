package inter

import (
	"crypto/sha256"

	"github.com/Fantom-foundation/lachesis-base/common/bigendian"
	"github.com/Fantom-foundation/lachesis-base/ltypes"
)

type LlrBlockVotes struct {
	Start ltypes.BlockID
	Epoch ltypes.EpochID
	Votes []ltypes.Hash
}

func (bvs LlrBlockVotes) LastBlock() ltypes.BlockID {
	return bvs.Start + ltypes.BlockID(len(bvs.Votes)) - 1
}

type LlrEpochVote struct {
	Epoch ltypes.EpochID
	Vote  ltypes.Hash
}

type LlrSignedBlockVotes struct {
	Signed                       SignedEventLocator
	TxsAndMisbehaviourProofsHash ltypes.Hash
	EpochVoteHash                ltypes.Hash
	Val                          LlrBlockVotes
}

type LlrSignedEpochVote struct {
	Signed                       SignedEventLocator
	TxsAndMisbehaviourProofsHash ltypes.Hash
	BlockVotesHash               ltypes.Hash
	Val                          LlrEpochVote
}

func (r SignedEventLocator) Size() uint64 {
	return uint64(len(r.Sig)) + 3*32 + 4*4
}

func (bvs LlrSignedBlockVotes) Size() uint64 {
	return bvs.Signed.Size() + uint64(len(bvs.Val.Votes))*32 + 32*2 + 8 + 4
}

func (ers LlrEpochVote) Hash() ltypes.Hash {
	hasher := sha256.New()
	hasher.Write(ers.Epoch.Bytes())
	hasher.Write(ers.Vote.Bytes())
	return ltypes.BytesToHash(hasher.Sum(nil))
}

func (bvs LlrBlockVotes) Hash() ltypes.Hash {
	hasher := sha256.New()
	hasher.Write(bvs.Start.Bytes())
	hasher.Write(bvs.Epoch.Bytes())
	hasher.Write(bigendian.Uint32ToBytes(uint32(len(bvs.Votes))))
	for _, bv := range bvs.Votes {
		hasher.Write(bv.Bytes())
	}
	return ltypes.BytesToHash(hasher.Sum(nil))
}

func (bvs LlrSignedBlockVotes) CalcPayloadHash() ltypes.Hash {
	return ltypes.Of(bvs.TxsAndMisbehaviourProofsHash.Bytes(), ltypes.Of(bvs.EpochVoteHash.Bytes(), bvs.Val.Hash().Bytes()).Bytes())
}

func (ev LlrSignedEpochVote) CalcPayloadHash() ltypes.Hash {
	return ltypes.Of(ev.TxsAndMisbehaviourProofsHash.Bytes(), ltypes.Of(ev.Val.Hash().Bytes(), ev.BlockVotesHash.Bytes()).Bytes())
}

func (ev LlrSignedEpochVote) Size() uint64 {
	return ev.Signed.Size() + 32 + 32*2 + 4 + 4
}
