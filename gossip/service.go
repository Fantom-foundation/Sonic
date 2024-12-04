package gossip

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Fantom-foundation/lachesis-base/lachesis"
	"github.com/Fantom-foundation/lachesis-base/ltypes"
	"github.com/Fantom-foundation/lachesis-base/utils/workers"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	notify "github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/dnsdisc"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/Fantom-foundation/go-opera/ethapi"
	"github.com/Fantom-foundation/go-opera/eventcheck"
	"github.com/Fantom-foundation/go-opera/eventcheck/basiccheck"
	"github.com/Fantom-foundation/go-opera/eventcheck/epochcheck"
	"github.com/Fantom-foundation/go-opera/eventcheck/gaspowercheck"
	"github.com/Fantom-foundation/go-opera/eventcheck/heavycheck"
	"github.com/Fantom-foundation/go-opera/eventcheck/parentscheck"
	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc/drivermodule"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc/eventmodule"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc/evmmodule"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc/sealmodule"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc/verwatcher"
	"github.com/Fantom-foundation/go-opera/gossip/emitter"
	"github.com/Fantom-foundation/go-opera/gossip/filters"
	"github.com/Fantom-foundation/go-opera/gossip/gasprice"
	"github.com/Fantom-foundation/go-opera/gossip/proclogger"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/logger"
	"github.com/Fantom-foundation/go-opera/utils/signers/gsignercache"
	"github.com/Fantom-foundation/go-opera/utils/txtime"
	"github.com/Fantom-foundation/go-opera/utils/wgmutex"
	"github.com/Fantom-foundation/go-opera/valkeystore"
	"github.com/Fantom-foundation/go-opera/vecmt"
)

type ServiceFeed struct {
	scope notify.SubscriptionScope

	newEpoch        notify.Feed
	newEmittedEvent notify.Feed
	newBlock        notify.Feed
	newLogs         notify.Feed
}

func (f *ServiceFeed) SubscribeNewEpoch(ch chan<- ltypes.EpochID) notify.Subscription {
	return f.scope.Track(f.newEpoch.Subscribe(ch))
}

func (f *ServiceFeed) SubscribeNewEmitted(ch chan<- *inter.EventPayload) notify.Subscription {
	return f.scope.Track(f.newEmittedEvent.Subscribe(ch))
}

func (f *ServiceFeed) SubscribeNewBlock(ch chan<- evmcore.ChainHeadNotify) notify.Subscription {
	return f.scope.Track(f.newBlock.Subscribe(ch))
}

func (f *ServiceFeed) SubscribeNewLogs(ch chan<- []*types.Log) notify.Subscription {
	return f.scope.Track(f.newLogs.Subscribe(ch))
}

type BlockProc struct {
	SealerModule     blockproc.SealerModule
	TxListenerModule blockproc.TxListenerModule
	PreTxTransactor  blockproc.TxTransactor
	PostTxTransactor blockproc.TxTransactor
	EventsModule     blockproc.ConfirmedEventsModule
	EVMModule        blockproc.EVM
}

func DefaultBlockProc() BlockProc {
	return BlockProc{
		SealerModule:     sealmodule.New(),
		TxListenerModule: drivermodule.NewDriverTxListenerModule(),
		PreTxTransactor:  drivermodule.NewDriverTxPreTransactor(),
		PostTxTransactor: drivermodule.NewDriverTxTransactor(),
		EventsModule:     eventmodule.New(),
		EVMModule:        evmmodule.New(),
	}
}

// Service implements go-ethereum/node.Service interface.
type Service struct {
	config Config

	// server
	p2pServer *p2p.Server
	Name      string

	accountManager *accounts.Manager

	// application
	store               *Store
	engine              lachesis.Consensus
	dagIndexer          *vecmt.Index
	engineMu            *sync.RWMutex
	emitters            []*emitter.Emitter
	txpool              TxPool
	heavyCheckReader    HeavyCheckReader
	gasPowerCheckReader GasPowerCheckReader
	checkers            *eventcheck.Checkers
	uniqueEventIDs      uniqueID

	// version watcher
	verWatcher *verwatcher.VerWarcher

	blockProcWg        sync.WaitGroup
	blockProcTasks     *workers.Workers
	blockProcTasksDone chan struct{}
	blockProcModules   BlockProc

	blockBusyFlag uint32
	eventBusyFlag uint32

	feed ServiceFeed

	gpo *gasprice.Oracle

	// application protocol
	handler *handler

	operaDialCandidates enode.Iterator

	EthAPI        *EthAPIBackend
	netRPCService *ethapi.PublicNetAPI

	procLogger *proclogger.Logger

	stopped   bool
	haltCheck func(oldEpoch, newEpoch ltypes.EpochID, time time.Time) bool

	tflusher PeriodicFlusher

	bootstrapping bool

	logger.Instance
}

func NewService(stack *node.Node, config Config, store *Store, blockProc BlockProc,
	engine lachesis.Consensus, dagIndexer *vecmt.Index, newTxPool func(evmcore.StateReader) TxPool,
	haltCheck func(oldEpoch, newEpoch ltypes.EpochID, age time.Time) bool) (*Service, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	localNodeId := enode.PubkeyToIDV4(&stack.Server().PrivateKey.PublicKey)
	svc, err := newService(config, store, blockProc, engine, dagIndexer, newTxPool, localNodeId)
	if err != nil {
		return nil, err
	}

	svc.p2pServer = stack.Server()
	svc.accountManager = stack.AccountManager()
	svc.EthAPI.SetExtRPCEnabled(stack.Config().ExtRPCEnabled())
	// Create the net API service
	svc.netRPCService = ethapi.NewPublicNetAPI(svc.p2pServer, store.GetRules().NetworkID)
	svc.haltCheck = haltCheck

	return svc, nil
}

func newService(config Config, store *Store, blockProc BlockProc, engine lachesis.Consensus, dagIndexer *vecmt.Index, newTxPool func(evmcore.StateReader) TxPool, localId enode.ID) (*Service, error) {
	svc := &Service{
		config:             config,
		blockProcTasksDone: make(chan struct{}),
		Name:               fmt.Sprintf("Node-%d", rand.Int()),
		store:              store,
		engine:             engine,
		blockProcModules:   blockProc,
		dagIndexer:         dagIndexer,
		engineMu:           new(sync.RWMutex),
		uniqueEventIDs:     uniqueID{new(big.Int)},
		procLogger:         proclogger.NewLogger(),
		Instance:           logger.New("gossip-service"),
	}

	svc.blockProcTasks = workers.New(new(sync.WaitGroup), svc.blockProcTasksDone, 1)

	// load epoch DB
	svc.store.loadEpochStore(svc.store.GetEpoch())
	es := svc.store.getEpochStore(svc.store.GetEpoch())
	svc.dagIndexer.Reset(svc.store.GetValidators(), es.table.DagIndex, func(id ltypes.EventHash) ltypes.Event {
		return svc.store.GetEvent(id)
	})

	// load caches for mutable values to avoid race condition
	svc.store.GetBlockEpochState()
	svc.store.GetHighestLamport()
	svc.store.GetUpgradeHeights()
	svc.store.GetGenesisID()
	netVerStore := verwatcher.NewStore(store.table.NetworkVersion)
	netVerStore.GetNetworkVersion()
	netVerStore.GetMissedVersion()

	// create checkers
	net := store.GetRules()
	txSigner := gsignercache.Wrap(types.LatestSignerForChainID(new(big.Int).SetUint64(net.NetworkID)))
	svc.heavyCheckReader.Store = store
	svc.heavyCheckReader.Pubkeys.Store(readEpochPubKeys(svc.store, svc.store.GetEpoch()))                                          // read pub keys of current epoch from DB
	svc.gasPowerCheckReader.Ctx.Store(NewGasPowerContext(svc.store, svc.store.GetValidators(), svc.store.GetEpoch(), net.Economy)) // read gaspower check data from DB
	svc.checkers = makeCheckers(config.HeavyCheck, txSigner, &svc.heavyCheckReader, &svc.gasPowerCheckReader, svc.store)

	// create GPO
	svc.gpo = gasprice.NewOracle(svc.config.GPO, nil)

	// create tx pool
	stateReader := &EvmStateReader{
		ServiceFeed: &svc.feed,
		store:       svc.store,
		gpo:         svc.gpo,
	}
	svc.txpool = newTxPool(stateReader)
	svc.gpo.SetReader(&GPOBackend{svc.store, svc.txpool})

	// init dialCandidates
	dnsclient := dnsdisc.NewClient(dnsdisc.Config{})
	var err error
	svc.operaDialCandidates, err = dnsclient.NewIterator(config.OperaDiscoveryURLs...)
	if err != nil {
		return nil, err
	}

	// create protocol manager
	svc.handler, err = newHandler(handlerConfig{
		config:   config,
		notifier: &svc.feed,
		txpool:   svc.txpool,
		engineMu: svc.engineMu,
		checkers: svc.checkers,
		s:        store,
		localId:  localId,
		process: processCallback{
			Event: func(event *inter.EventPayload) error {
				done := svc.procLogger.EventConnectionStarted(event, false)
				defer done()
				return svc.processEvent(event)
			},
			SwitchEpochTo: svc.SwitchEpochTo,
		},
		localEndPointSource: localEndPointSource{svc},
	})
	if err != nil {
		return nil, err
	}

	rpc.SetExecutionTimeLimit(config.RPCTimeout)

	// create API backend
	svc.EthAPI = &EthAPIBackend{false, svc, stateReader, txSigner, config.AllowUnprotectedTxs}

	svc.verWatcher = verwatcher.New(netVerStore)
	svc.tflusher = svc.makePeriodicFlusher()

	return svc, nil
}

type localEndPointSource struct {
	service *Service
}

func (s localEndPointSource) GetLocalEndPoint() *enode.Node {
	return s.service.p2pServer.LocalNode().Node()
}

// makeCheckers builds event checkers
func makeCheckers(heavyCheckCfg heavycheck.Config, txSigner types.Signer, heavyCheckReader *HeavyCheckReader, gasPowerCheckReader *GasPowerCheckReader, store *Store) *eventcheck.Checkers {
	// create signatures checker
	heavyCheck := heavycheck.New(heavyCheckCfg, heavyCheckReader, txSigner)

	// create gaspower checker
	gaspowerCheck := gaspowercheck.New(gasPowerCheckReader)

	return &eventcheck.Checkers{
		Basiccheck:    basiccheck.New(),
		Epochcheck:    epochcheck.New(store),
		Parentscheck:  parentscheck.New(),
		Heavycheck:    heavyCheck,
		Gaspowercheck: gaspowerCheck,
	}
}

// makePeriodicFlusher makes PeriodicFlusher
func (s *Service) makePeriodicFlusher() PeriodicFlusher {
	// Normally the diffs are committed by message processing. Yet, since we have async EVM snapshots generation,
	// we need to periodically commit its data
	return PeriodicFlusher{
		period: 10 * time.Millisecond,
		callback: PeriodicFlusherCallaback{
			busy: func() bool {
				// try to lock engineMu/blockProcWg pair as rarely as possible to not hurt
				// events/blocks pipeline concurrency
				return atomic.LoadUint32(&s.eventBusyFlag) != 0 || atomic.LoadUint32(&s.blockBusyFlag) != 0
			},
			commitNeeded: func() bool {
				// use slightly higher size threshold to avoid locking the mutex/wg pair and hurting events/blocks concurrency
				// PeriodicFlusher should mostly commit only data generated by async EVM snapshots generation
				return s.store.isCommitNeeded(1200, 1000)
			},
			commit: func() {
				s.engineMu.Lock()
				defer s.engineMu.Unlock()
				// Note: blockProcWg.Wait() is already called by s.commit
				if s.store.isCommitNeeded(1200, 1000) {
					s.commit(false)
				}
			},
		},
		wg:   sync.WaitGroup{},
		quit: make(chan struct{}),
	}
}

func (s *Service) EmitterWorld(signer valkeystore.SignerI) emitter.World {
	return emitter.World{
		External: &emitterWorld{
			emitterWorldProc: emitterWorldProc{s},
			emitterWorldRead: emitterWorldRead{s.store},
			WgMutex:          wgmutex.New(s.engineMu, &s.blockProcWg),
		},
		TxPool:   s.txpool,
		Signer:   signer,
		TxSigner: s.EthAPI.signer,
	}
}

// RegisterEmitter must be called before service is started
func (s *Service) RegisterEmitter(em *emitter.Emitter) {
	txtime.Enabled = true // enable tracking of tx times
	s.emitters = append(s.emitters, em)
}

// MakeProtocols constructs the P2P protocol definitions for `opera`.
func MakeProtocols(svc *Service, backend *handler, disc enode.Iterator) []p2p.Protocol {
	nodeIter := enode.NewFairMix(time.Second)
	nodeIter.AddSource(disc)
	nodeIter.AddSource(backend.GetSuggestedPeerIterator())

	protocols := make([]p2p.Protocol, len(ProtocolVersions))
	for i, version := range ProtocolVersions {
		version := version // Closure

		protocols[i] = p2p.Protocol{
			Name:    ProtocolName,
			Version: version,
			Length:  protocolLengths[version],
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				// wait until handler has started
				backend.started.Wait()
				peer := newPeer(version, p, rw, backend.config.Protocol.PeerCache)
				defer peer.Close()

				select {
				case <-backend.quitSync:
					return p2p.DiscQuitting
				default:
					backend.wg.Add(1)
					defer backend.wg.Done()
					return backend.handle(peer)
				}
			},
			NodeInfo: func() interface{} {
				return backend.NodeInfo()
			},
			PeerInfo: func(id enode.ID) interface{} {
				if p := backend.peers.Peer(id.String()); p != nil {
					return p.Info()
				}
				return nil
			},
			Attributes:     []enr.Entry{currentENREntry(svc, 0 /* time */)},
			DialCandidates: nodeIter,
		}
	}
	return protocols
}

// Protocols returns protocols the service can communicate on.
func (s *Service) Protocols() []p2p.Protocol {
	return MakeProtocols(s, s.handler, s.operaDialCandidates)
}

// APIs returns api methods the service wants to expose on rpc channels.
func (s *Service) APIs() []rpc.API {
	apis := ethapi.GetAPIs(s.EthAPI)

	apis = append(apis, []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicEthereumAPI(s),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.EthAPI, s.config.FilterAPI),
			Public:    true,
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   ethapi.NewPublicDebugAPI(s.EthAPI, s.config.MaxResponseSize, s.config.StructLogLimit),
			Public:    true,
		}, {
			Namespace: "trace",
			Version:   "1.0",
			Service:   ethapi.NewPublicTxTraceAPI(s.EthAPI, s.config.MaxResponseSize),
			Public:    true,
		},
	}...)

	// eth-namespace is doubled as ftm-namespace for branding purpose
	for _, api := range apis {
		if api.Namespace == "eth" {
			apis = append(apis, rpc.API{
				Namespace: "ftm",
				Version:   api.Version,
				Service:   api.Service,
				Public:    api.Public,
			})
		}
	}

	return apis
}

// Start method invoked when the node is ready to start the service.
func (s *Service) Start() error {
	s.gpo.Start()
	// start tflusher before starting snapshots generation
	s.tflusher.Start()
	blockState := s.store.GetBlockState()
	if s.store.evm.CheckLiveStateHash(blockState.LastBlock.Idx, blockState.FinalizedStateRoot) != nil {
		return errors.New("fullsync isn't possible because state root is missing")
	}

	// start blocks processor
	s.blockProcTasks.Start(1)

	// start p2p
	StartENRUpdater(s, s.p2pServer.LocalNode())
	s.handler.Start(s.p2pServer.MaxPeers)

	// start emitters
	for _, em := range s.emitters {
		em.Start()
	}

	s.verWatcher.Start()

	if s.haltCheck != nil && s.haltCheck(s.store.GetEpoch(), s.store.GetEpoch(), s.store.GetBlockState().LastBlock.Time.Time()) {
		// halt syncing
		s.stopped = true
	}

	return nil
}

// WaitBlockEnd waits until parallel block processing is complete (if any)
func (s *Service) WaitBlockEnd() {
	s.blockProcWg.Wait()
}

// Stop method invoked when the node terminates the service.
func (s *Service) Stop() error {
	defer log.Info("Fantom service stopped")
	s.verWatcher.Stop()
	for _, em := range s.emitters {
		em.Stop()
	}

	// Stop all the peer-related stuff first.
	s.operaDialCandidates.Close()

	s.handler.Stop()
	s.feed.scope.Close()
	s.gpo.Stop()
	// it's safe to stop tflusher only before locking engineMu
	s.tflusher.Stop()

	// flush the state at exit, after all the routines stopped
	s.engineMu.Lock()
	defer s.engineMu.Unlock()
	s.stopped = true

	s.blockProcWg.Wait()
	close(s.blockProcTasksDone)

	err := s.dagIndexer.Close()
	if err != nil {
		return err
	}

	return s.store.Commit()
}

// AccountManager return node's account manager
func (s *Service) AccountManager() *accounts.Manager {
	return s.accountManager
}
