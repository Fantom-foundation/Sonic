package launcher

import (
	"path"

	"github.com/Fantom-foundation/go-opera/integration"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/ethereum/go-ethereum/cmd/utils"
)

func makeDBsProducer(cfg *config) kvdb.FullDBProducer {
	if err := integration.CheckStateInitialized(path.Join(cfg.Node.DataDir, "chaindata"), cfg.DBs); err != nil {
		utils.Fatalf(err.Error())
	}
	producer, err := integration.GetDbProducer(path.Join(cfg.Node.DataDir, "chaindata"), cfg.DBs.RuntimeCache)
	if err != nil {
		utils.Fatalf("Failed to initialize DB producer: %v", err)
	}
	return producer
}
