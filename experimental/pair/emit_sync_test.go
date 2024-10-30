package pair_sync

import (
	"fmt"
	"math"
	"testing"

	"pgregory.net/rand"
)

const interval = 100

type ID int

type Message struct {
	from ID
}

type Node struct {
	id ID

	lastEmitTime    int
	lastMessageTime int
	beforeDelta     int
	afterDelta      int
}

func (n *Node) WantEmit(time int) bool {
	passed := time - n.lastEmitTime
	limit := interval
	if n.afterDelta < math.MaxInt {
		diff := n.afterDelta - n.beforeDelta
		diff = diff * 3 / 4
		limit += diff
	}
	res := passed >= limit
	if res {
		n.lastEmitTime = time
	}
	return res
}

func (n *Node) Receive(time int, msg Message) {
	if msg.from == n.id {
		n.beforeDelta = time - n.lastMessageTime
		n.afterDelta = math.MaxInt
	} else if delta := time - n.lastEmitTime; delta < n.afterDelta {
		n.afterDelta = delta
	}
	n.lastMessageTime = time
}

func RunSim(numNodes, duration int) {
	nodes := make([]Node, numNodes)
	for i := 0; i < numNodes; i++ {
		nodes[i].id = ID(i)
	}

	offset := make([]int, numNodes)
	for i := range offset {
		offset[i] = rand.Intn(interval)
	}

	fmt.Printf("Time,Node,Iteration,delta\n")
	last := 0
	progress := make([]int, numNodes)
	for time := 0; time < duration; time++ {
		for i := 0; i < numNodes; i++ {
			if nodes[i].WantEmit(time + offset[i]) {
				progress[i]++
				fmt.Printf("%d,%d,%d,%d\n", time, i, progress[i], time-last)
				last = time
				msg := Message{from: nodes[i].id}
				// TODO: add delay
				for i := range nodes {
					nodes[i].Receive(time+offset[i], msg)
				}
			}
		}
	}
}

func TestSync_RunExample(t *testing.T) {
	RunSim(2, 100*interval)
	t.Fail()
}
