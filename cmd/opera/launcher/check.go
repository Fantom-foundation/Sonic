package launcher

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/gossip"
	"time"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"
)

func checkEvm(ctx *cli.Context) error {
	if len(ctx.Args()) != 0 {
		utils.Fatalf("This command doesn't require an argument.")
	}

	cfg := makeAllConfigs(ctx)

	rawDbs := makeDBsProducer(cfg)
	gdb, err := gossip.NewStore(rawDbs, cfg.OperaStore)
	if err != nil {
		return fmt.Errorf("failed to create gossip store: %w", err)
	}
	defer gdb.Close()

	start := time.Now()

	lastBlockIdx := gdb.GetLatestBlockIndex()
	lastBlock := gdb.GetBlock(lastBlockIdx)
	if lastBlock == nil {
		log.Crit("Verification of the database failed - unable to get the last block")
	}

	err = gdb.EvmStore().VerifyWorldState(uint64(lastBlockIdx), common.Hash(lastBlock.Root))
	if err != nil {
		log.Crit("Verification of the Fantom World State failed", "err", err)
	}
	log.Info("EVM storage is verified", "last", lastBlockIdx, "elapsed", common.PrettyDuration(time.Since(start)))
	return nil
}
