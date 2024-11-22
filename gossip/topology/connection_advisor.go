package topology

import (
	"sync"

	"github.com/ethereum/go-ethereum/p2p/enode"
)

type ConnectionAdvisor interface {
	// GetNewPeerSuggestion returns a new peer that should be connected to.
	GetNewPeerSuggestion() *enode.Node

	// GetRedundantPeerSuggestion returns the ID of a peer that should be
	// disconnected in favor of another peer if needed.
	GetRedundantPeerSuggestion() *enode.ID

	// UpdatePeers updates information about the peers connected to a given peer.
	// This information is used to obtain an overview on the local node's
	// neighborhood from which decisions changes in the peer set can be made.
	// The provided peer is assumed to be a peer of the local node.
	UpdatePeers(peer enode.ID, peers []*enode.Node)
}

func NewConnectionAdvisor(localId enode.ID) ConnectionAdvisor {
	return &connectionAdvisor{
		neighborhood: make(map[enode.ID][]*enode.Node),
		localId:      localId,
	}
}

type connectionAdvisor struct {
	mu sync.Mutex

	// Keep track of the neighbors of each peer.
	// TODO: add support for forgetting old information.
	neighborhood map[enode.ID][]*enode.Node

	// The ID of the local node.
	localId enode.ID
}

func (c *connectionAdvisor) GetNewPeerSuggestion() *enode.Node {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Search for a peer of a peer that is not already connected to the local node.
	for _, peers := range c.neighborhood {
		for _, peer := range peers {
			if peer.ID() == c.localId {
				continue
			}
			if _, found := c.neighborhood[peer.ID()]; !found {
				//			fmt.Printf("Suggestion: %v\n", peer.ID())
				return peer
			}
		}
	}

	return nil
}

func (c *connectionAdvisor) GetRedundantPeerSuggestion() *enode.ID {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Count the number of indirect connections to all peers and
	// recommend the one that has the most indirect connections.
	count := map[enode.ID]int{}
	for _, peers := range c.neighborhood {
		for _, peer := range peers {
			if _, found := c.neighborhood[peer.ID()]; found {
				count[peer.ID()]++
			}
		}
	}
	delete(count, c.localId)

	var maxCount int
	var maxPeer enode.ID
	for peer, c := range count {
		if c > maxCount {
			maxCount = c
			maxPeer = peer
		}
	}

	if maxCount == 0 {
		return nil
	}
	return &maxPeer
}

func (c *connectionAdvisor) UpdatePeers(peer enode.ID, peers []*enode.Node) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.neighborhood[peer] = peers
}
