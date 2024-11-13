package inter

import (
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func TestBlock_GetPrevRandao_IsNeverZero(t *testing.T) {
	blk := Block{}
	if h := blk.GetPrevRandao(); h == (common.Hash{}) {
		t.Error("prevrandao must never be zero")
	}
}
