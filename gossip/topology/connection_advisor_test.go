package topology

import (
	"encoding/binary"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
)

func TestConnectionAdvisor_LinearTopologyTurnsIntoMesh(t *testing.T) {

	const N = 10

	nodes := []node{}
	for i := range N {
		nodes = append(nodes, node{
			id: createENode(i),
		})
	}

	// Create a linear topology
	for i := range nodes {
		if i > 0 {
			nodes[i].peers = append(nodes[i].peers, nodes[i-1].id)
		}
		if i < len(nodes)-1 {
			nodes[i].peers = append(nodes[i].peers, nodes[i+1].id)
		}
	}

	print := func() {
		for i := range nodes {
			fmt.Printf("Node %d has %d peers: ", i, len(nodes[i].peers))
			ids := []int{}
			for _, peer := range nodes[i].peers {
				ids = append(ids, toId(peer))
			}
			slices.Sort(ids)
			idStrings := []string{}
			for _, id := range ids {
				idStrings = append(idStrings, fmt.Sprintf("%d", id))
			}
			fmt.Printf("%s\n", strings.Join(idStrings, ", "))
		}
		fmt.Println()
	}

	print()

	for range 2 * N {
		newPeers := make([]*enode.Node, len(nodes))
		remPeers := make([]*enode.ID, len(nodes))
		for i, node := range nodes {
			self := node.id.ID()
			advisor := NewConnectionAdvisor(self)
			for _, peer := range node.peers {
				advisor.UpdatePeers(peer.ID(), nodes[toId(peer)].peers)
			}
			newPeer := advisor.GetNewPeerSuggestion()
			if newPeer != nil {
				fmt.Printf("Node %d should connect to %d\n", toId(node.id), toId(newPeer))
			} else {
				fmt.Printf("Node %d should not connect to any new peers\n", toId(node.id))
			}
			newPeers[i] = newPeer

			removePeer := advisor.GetRedundantPeerSuggestion()
			if removePeer != nil {
				fmt.Printf("Node %d should remove %d\n", toId(node.id), toId2(*removePeer))
			} else {
				fmt.Printf("Node %d should not remove any peers\n", toId(node.id))
			}
			remPeers[i] = removePeer
		}

		for i, newPeer := range newPeers {
			if newPeer == nil {
				continue
			}
			if slices.Contains(nodes[i].peers, newPeer) {
				continue
			}
			nodes[i].peers = append(nodes[i].peers, newPeer)
		}

		for i, remPeer := range remPeers {
			if remPeer == nil {
				continue
			}
			if len(nodes[i].peers) > 4 {
				nodes[i].peers = slices.DeleteFunc(nodes[i].peers, func(peer *enode.Node) bool {
					return peer.ID() == *remPeer
				})
			}
		}

		fmt.Println()

		print()
	}

	t.Fail()
}

type node struct {
	id    *enode.Node
	peers []*enode.Node
}

func createENode(i int) *enode.Node {
	id := enode.ID{}
	binary.BigEndian.PutUint64(id[24:], uint64(i))
	r := enr.Record{}
	r.Set(enr.WithEntry("nulladdr", &id))
	node, err := enode.New(enode.NullID{}, &r)
	if err != nil {
		panic(err)
	}
	return node
}

func toId(node *enode.Node) int {
	return toId2(node.ID())
}

func toId2(node enode.ID) int {
	return int(binary.BigEndian.Uint64(node[24:]))
}
