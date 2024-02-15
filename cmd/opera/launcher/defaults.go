package launcher

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/opera/launcher/utils"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	DefaultHTTPPort = 8545 // Default TCP port for the HTTP RPC server
	DefaultWSPort   = 8546 // Default TCP port for the websocket RPC server
)

// NodeDefaultConfig contains reasonable default settings.
var NodeDefaultConfig = node.Config{
	HTTPPort:            DefaultHTTPPort,
	HTTPModules:         []string{},
	HTTPVirtualHosts:    []string{"localhost"},
	HTTPTimeouts:        rpc.DefaultHTTPTimeouts,
	WSPort:              DefaultWSPort,
	WSModules:           []string{},
	GraphQLVirtualHosts: []string{"localhost"},
	P2P: p2p.Config{
		NoDiscovery: false, // enable discovery v4 by default
		DiscoveryV5: true,  // enable discovery v5 by default
		ListenAddr:  fmt.Sprintf(":%d", utils.ListenPortFlag.Value),
		MaxPeers:    50,
		NAT:         nat.Any(),
	},
}
