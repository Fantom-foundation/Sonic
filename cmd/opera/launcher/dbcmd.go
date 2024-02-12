package launcher

import (
	"fmt"
	"path"

	"github.com/Fantom-foundation/go-opera/integration"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
)

func makeDBsProducer(cfg *config) (kvdb.FullDBProducer, error) {
	if err := integration.CheckStateInitialized(path.Join(cfg.Node.DataDir, "chaindata"), cfg.DBs); err != nil {
		return nil, err
	}
	producer, err := integration.GetDbProducer(path.Join(cfg.Node.DataDir, "chaindata"), cfg.DBs.RuntimeCache)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize DB producer: %w", err)
	}
	return producer, nil
}
