package fakemultidb

import (
	"strings"

	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/kvdb/table"
)

type Producer struct {
	kvdb.FullDBProducer
}

// NewProducer of a combined producer for multiple types of DBs.
func NewProducer(producer kvdb.FullDBProducer) (*Producer, error) {
	return &Producer{producer}, nil
}

// OpenDB or create db with name.
// Example of names: "genesis", "gossip/S", "lachesis-123"
func (p *Producer) OpenDB(req string) (kvdb.Store, error) {
	slashPos := strings.LastIndexByte(req, '/')
	name := req
	if slashPos >= 0 {
		name = req[:slashPos]
	}
	db, err := p.FullDBProducer.OpenDB(name)
	if err != nil {
		return nil, err
	}
	cdb := &closableTable{
		Store:      db,
		underlying: db,
	}
	if slashPos >= 0 {
		cdb.Store = table.New(db, []byte(req[slashPos+1:]))
	}
	return cdb, nil
}

func (p *Producer) Close() error {
	return p.FullDBProducer.Close()
}

type closableTable struct {
	kvdb.Store
	underlying kvdb.Store
}

// Close leaves underlying database.
func (s *closableTable) Close() error {
	return s.underlying.Close()
}

// Drop whole database.
func (s *closableTable) Drop() {
	s.underlying.Drop()
}
