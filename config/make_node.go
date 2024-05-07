package config

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/config/flags"
	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/gossip"
	"github.com/Fantom-foundation/go-opera/gossip/emitter"
	"github.com/Fantom-foundation/go-opera/integration"
	"github.com/Fantom-foundation/go-opera/utils/errlock"
	"github.com/Fantom-foundation/go-opera/valkeystore"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"gopkg.in/urfave/cli.v1"
	"path"
	"time"
)

func MakeNode(ctx *cli.Context, cfg *Config) (*node.Node, *gossip.Service, func(), error) {
	var success bool
	var cleanup []func()
	defer func() { // if the function fails, clean-up in defer, otherwise return cleaning function
		if !success {
			for i := len(cleanup) - 1; i >= 0; i-- {
				cleanup[i]()
			}
		}
	}()

	// check errlock file
	errlock.SetDefaultDatadir(cfg.Node.DataDir)
	if err := errlock.Check(); err != nil {
		return nil, nil, nil, err
	}

	engine, dagIndex, gdb, cdb, blockProc, closeDBs, err := integration.MakeEngine(path.Join(cfg.Node.DataDir, "chaindata"), cfg.AppConfigs())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to make consensus engine: %w", err)
	}
	cleanup = append(cleanup, func() {
		if err := gdb.Close(); err != nil {
			log.Warn("Failed to close gossip store", "err", err)
		}
		if err := cdb.Close(); err != nil {
			log.Warn("Failed to close consensus database", "err", err)
		}
		if closeDBs != nil {
			if err := closeDBs(); err != nil {
				log.Warn("Failed to close databases", "err", err)
			}
		}
	})

	// substitute default bootnodes if requested
	networkName := ""
	if gdb.HasBlockEpochState() {
		networkName = gdb.GetRules().Name
	}
	if needDefaultBootnodes(cfg.Node.P2P.BootstrapNodes) {
		bootnodes := Bootnodes[networkName]
		if bootnodes == nil {
			bootnodes = []string{}
		}
		setBootnodes(ctx, bootnodes, &cfg.Node)
	}

	stack, err := makeNetworkStack(ctx, &cfg.Node)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unlock validator key: %w", err)
	}
	cleanup = append(cleanup, func() {
		if err := stack.Close(); err != nil && err != node.ErrNodeStopped {
			log.Warn("Failed to close stack", "err", err)
		}
	})

	_, _, keystoreDir, err := cfg.Node.AccountConfig()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to setup account config: %w", err)
	}
	valKeystore := valkeystore.NewDefaultFileKeystore(path.Join(keystoreDir, "validator"))
	valPubkey := cfg.Emitter.Validator.PubKey
	if key := getFakeValidatorKey(ctx); key != nil && cfg.Emitter.Validator.ID != 0 {
		if err := addFakeValidatorKey(ctx, key, valPubkey, valKeystore); err != nil {
			return nil, nil, nil, err
		}
		coinbase := integration.SetAccountKey(stack.AccountManager(), key, "fakepassword")
		log.Info("Unlocked fake validator account", "address", coinbase.Address.Hex())
	}

	// unlock validator key
	if !valPubkey.Empty() {
		err := unlockValidatorKey(ctx, valPubkey, valKeystore)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to unlock validator key: %w", err)
		}
	}
	signer := valkeystore.NewSigner(valKeystore)

	// Create and register a gossip network service.
	newTxPool := func(reader evmcore.StateReader) gossip.TxPool {
		if cfg.TxPool.Journal != "" {
			cfg.TxPool.Journal = path.Join(cfg.Node.DataDir, cfg.TxPool.Journal)
		}
		return evmcore.NewTxPool(cfg.TxPool, reader.Config(), reader)
	}
	haltCheck := func(oldEpoch, newEpoch idx.Epoch, age time.Time) bool {
		stop := ctx.GlobalIsSet(flags.ExitWhenAgeFlag.Name) && ctx.GlobalDuration(flags.ExitWhenAgeFlag.Name) >= time.Since(age)
		stop = stop || ctx.GlobalIsSet(flags.ExitWhenEpochFlag.Name) && idx.Epoch(ctx.GlobalUint64(flags.ExitWhenEpochFlag.Name)) <= newEpoch
		if stop {
			go func() {
				// do it in a separate thread to avoid deadlock
				_ = stack.Close()
			}()
			return true
		}
		return false
	}
	svc, err := gossip.NewService(stack, cfg.Opera, gdb, blockProc, engine, dagIndex, newTxPool, haltCheck)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create the gossip service: %w", err)
	}
	err = engine.Bootstrap(svc.GetConsensusCallbacks())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to bootstrap the consensus engine: %w", err)
	}
	err = engine.Reset(gdb.GetEpoch(), gdb.GetValidators())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to reset the consensus engine: %w", err)
	}
	svc.ReprocessEpochEvents()

	// commit dbs to avoid dirty state when the rest of the startup fails
	if err := gdb.Commit(); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to commit gossip store: %w", err)
	}

	if cfg.Emitter.Validator.ID != 0 {
		svc.RegisterEmitter(emitter.NewEmitter(cfg.Emitter, svc.EmitterWorld(signer)))
	}

	stack.RegisterAPIs(svc.APIs())
	stack.RegisterProtocols(svc.Protocols())
	stack.RegisterLifecycle(svc)

	success = true // skip cleanup in defer - keep it for the returned cleanup function
	return stack, svc, func() {
		for i := len(cleanup) - 1; i >= 0; i-- {
			cleanup[i]()
		}
	}, nil
}

func makeNetworkStack(ctx *cli.Context, cfg *node.Config) (*node.Node, error) {
	stack, err := node.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create the protocol stack: %w", err)
	}
	return stack, nil
}
