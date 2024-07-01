package vecclock

import (
	"github.com/Fantom-foundation/go-opera/vecclock/highestbefore"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
)

type allVecs struct {
	after  []idx.Event
	before []highestbefore.Type
}
