package pair_sync

import (
	"fmt"
	"testing"

	"pgregory.net/rand"
)

const interval = 1000
const iterations = 100

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
	afterValid      bool
}

func (n *Node) WantEmit(time int) bool {
	passed := time - n.lastEmitTime
	limit := interval
	if n.afterValid {
		delta := interval / 100
		diff := n.afterDelta - n.beforeDelta
		diff = diff * 3 / 4
		if diff > delta {
			diff = delta
		}
		if diff < -delta {
			diff = -delta
		}
		limit += diff
		//fmt.Printf("N%d: Emitted at %d, passed %d, limit %d\n", n.id, time, passed, limit)
	}
	res := passed >= limit
	if res {
		n.beforeDelta = time - n.lastMessageTime
		n.lastEmitTime = time
	}
	return res
}

func (n *Node) Receive(time int, msg Message) {
	if msg.from == n.id {
		n.afterValid = false
	} else if !n.afterValid {
		n.afterDelta = time - n.lastEmitTime
		n.afterValid = true
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
			if i == 0 && time < duration/3 {
				continue
			}
			if i == 1 && time > duration*2/3 {
				continue
			}
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
	RunSim(4, iterations*interval)
	t.Fail()
}
