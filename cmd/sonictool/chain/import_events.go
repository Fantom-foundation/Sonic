package chain

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Fantom-foundation/go-opera/config"
	"github.com/Fantom-foundation/go-opera/gossip"
	"github.com/Fantom-foundation/go-opera/gossip/emitter"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/utils/ioread"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	gzip "github.com/klauspost/pgzip"
	"github.com/status-im/keycard-go/hexutils"
	"gopkg.in/urfave/cli.v1"
	"io"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func EventsImport(ctx *cli.Context, files ...string) error {
	// avoid P2P interaction, API calls and events emitting
	cfg, err := config.MakeAllConfigs(ctx)
	if err != nil {
		return err
	}
	cfg.Opera.Protocol.EventsSemaphoreLimit.Size = math.MaxUint32
	cfg.Opera.Protocol.EventsSemaphoreLimit.Num = math.MaxUint32
	cfg.Emitter.Validator = emitter.ValidatorConfig{}
	cfg.TxPool.Journal = ""
	cfg.Node.IPCPath = ""
	cfg.Node.HTTPHost = ""
	cfg.Node.WSHost = ""
	cfg.Node.P2P.ListenAddr = ""
	cfg.Node.P2P.NoDiscovery = true
	cfg.Node.P2P.BootstrapNodes = nil
	cfg.Node.P2P.DiscoveryV5 = false
	cfg.Node.P2P.BootstrapNodesV5 = nil
	cfg.Node.P2P.StaticNodes = nil
	cfg.Node.P2P.TrustedNodes = nil

	node, svc, nodeClose, err := config.MakeNode(ctx, cfg)
	if err != nil {
		return err
	}
	defer nodeClose()

	if err := node.Start(); err != nil {
		return fmt.Errorf("error starting protocol stack: %w", err)
	}

	for _, fn := range files {
		log.Info("Importing events from file", "file", fn)
		if err := importEventsFile(svc, fn); err != nil {
			log.Error("Import error", "file", fn, "err", err)
			return err
		}
	}
	return nil
}

func checkEventsFileHeader(reader io.Reader) error {
	headerAndVersion := make([]byte, len(eventsFileHeader)+len(eventsFileVersion))
	err := ioread.ReadAll(reader, headerAndVersion)
	if err != nil {
		return err
	}
	if bytes.Compare(headerAndVersion[:len(eventsFileHeader)], eventsFileHeader) != 0 {
		return errors.New("expected an events file, mismatched file header")
	}
	if bytes.Compare(headerAndVersion[len(eventsFileHeader):], eventsFileVersion) != 0 {
		got := hexutils.BytesToHex(headerAndVersion[len(eventsFileHeader):])
		expected := hexutils.BytesToHex(eventsFileVersion)
		return errors.New(fmt.Sprintf("wrong version of events file, got=%s, expected=%s", got, expected))
	}
	return nil
}

func importEventsFile(srv *gossip.Service, fn string) error {
	// Watch for Ctrl-C while the import is running.
	// If a signal is received, the import will stop.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	// Open the file handle and potentially unwrap the gzip stream
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer fh.Close()

	var reader io.Reader = fh
	if strings.HasSuffix(fn, ".gz") {
		if reader, err = gzip.NewReader(reader); err != nil {
			return err
		}
		defer reader.(*gzip.Reader).Close()
	}

	// Check file version and header
	if err := checkEventsFileHeader(reader); err != nil {
		return err
	}

	stream := rlp.NewStream(reader, 0)

	start := time.Now()
	last := hash.Event{}

	batch := make(inter.EventPayloads, 0, 8*1024)
	batchSize := 0
	maxBatchSize := 8 * 1024 * 1024
	epoch := idx.Epoch(0)
	txs := 0
	events := 0

	processBatch := func() error {
		if batch.Len() == 0 {
			return nil
		}
		done := make(chan struct{})
		err := srv.DagProcessor().Enqueue("", batch.Bases(), true, nil, func() {
			done <- struct{}{}
		})
		if err != nil {
			return err
		}
		<-done
		last = batch[batch.Len()-1].ID()
		batch = batch[:0]
		batchSize = 0
		return nil
	}

	for {
		select {
		case <-interrupt:
			return fmt.Errorf("interrupted")
		default:
		}
		e := new(inter.EventPayload)
		err = stream.Decode(e)
		if err == io.EOF {
			err = processBatch()
			if err != nil {
				return err
			}
			break
		}
		if err != nil {
			return err
		}
		if e.Epoch() != epoch || batchSize >= maxBatchSize {
			err = processBatch()
			if err != nil {
				return err
			}
		}
		epoch = e.Epoch()
		batch = append(batch, e)
		batchSize += 1024 + e.Size()
		txs += e.Txs().Len()
		events++
	}
	srv.WaitBlockEnd()
	log.Info("Events import is finished", "file", fn, "last", last.String(), "imported", events, "txs", txs, "elapsed", common.PrettyDuration(time.Since(start)))

	return nil
}
