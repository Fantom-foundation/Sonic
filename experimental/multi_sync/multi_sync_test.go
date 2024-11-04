package multi_sync

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

const interval = 10000
const withDelay = false
const withNoise = false

func TestSimulateSync(t *testing.T) {
	const N = 10
	const T = 20
	runSync(N, T, t)
}

func runSync(N, T int, t *testing.T) {

	// Create offsets for N nodes (secret).
	offsets := make([]int, N)
	for i := 0; i < N; i++ {
		offsets[i] = rand.Intn(interval)
	}

	// Create N nodes.
	nodes := make([]Node, N)
	for i := range nodes {
		nodes[i] = newNode(N, i)
	}

	// create message delay matrix
	// delay[i][j] is the delay from i to j
	// delay[i][i] is 0
	// delay[i][j] is random
	delay := make([][]int, N)
	for i := range delay {
		delay[i] = make([]int, N)
		if withDelay {
			for j := range delay[i] {
				if i == j {
					delay[i][j] = 0
				} else {
					delay[i][j] = rand.Intn(interval/10) + 1 // 1 to 100 ms delay
				}
			}
		}
	}

	// Run simulation.
	output := true
	if output {
		fmt.Printf("Iteration")
		for i := range N {
			fmt.Printf(",N%d", i)
		}
		fmt.Printf("\n")
	}
	for step := range T {

		// Compute the next emit time.
		emitTimes := make([]int, N)
		for i := range N {
			emitTimes[i] = nodes[i].GetNextEmitTime() - offsets[i]
		}

		// Deliver all messages.
		for i := range N {
			for j := range N {
				noise := 0
				if withNoise && delay[i][j] > 1 {
					noise = rand.Intn(delay[i][j] / 2)
				}
				nodes[j].Receive(emitTimes[i]+delay[i][j]+noise+offsets[j], Message{from: i})
			}
		}

		if output {
			fmt.Printf("%d,", step)
			for i := range N {
				fmt.Printf("%d,", emitTimes[i]+interval)
			}
			fmt.Printf("\n")
		}

		// Print the delay estimate errors.
		if false {
			for i := range N {
				for j := range N {
					is := delay[i][j]
					est := nodes[j].delayEstimate[i]
					diff := is - est
					//fmt.Printf("%d,%d,%d,", is, est, diff)
					fmt.Printf("%d,", diff)
					//fmt.Printf("N%d->N%d: is %d, est %d\n", i, j, delay[i][j], nodes[j].delayEstimate[i])
				}
			}
			fmt.Printf("\n")
		}
	}

	/*
		for i := range N {
			for j := range N {
				fmt.Printf("N%d->N%d: is %d, est %d\n", i, j, delay[i][j], nodes[j].delayEstimate[i])
			}
		}
	*/

	t.Fail()
}

type Message struct {
	from int
}

type Node struct {
	N                int
	id               int
	lastEmitTime     int
	lastReceiveTimes []int
	delayEstimate    []int
}

func newNode(N, id int) Node {
	return Node{
		N:                N,
		id:               id,
		lastReceiveTimes: make([]int, N),
		delayEstimate:    make([]int, N),
	}
}

func (n *Node) GetNextEmitTime() int {
	// compute the offset this node should emit its message
	should := interval * n.id / n.N

	// estimated local offset to global time
	offset := n.lastEmitTime - should

	// refresh the delay estimate
	{
		for i := range n.N {
			expectedSend := interval*i/n.N + offset
			lastReceived := n.lastReceiveTimes[i]
			lastDelay := lastReceived - expectedSend

			//fmt.Printf("N%d: source %d, expected %d, received %d, delay %d\n", n.id, i, expectedSend, lastReceived, lastDelay)
			if lastDelay >= 0 {
				n.delayEstimate[i] = lastDelay*1/3 + n.delayEstimate[i]*2/3
			}
		}
	}

	// TODO: factor in the delay estimate

	// compute the mean time of all received messages
	sum := 0
	min := math.MaxInt
	max := math.MinInt
	for i := range n.lastReceiveTimes {
		cur := n.lastReceiveTimes[i]
		//cur = cur - n.delayEstimate[i]

		sum += cur
		if n.lastReceiveTimes[i] < min {
			min = cur
		}
		if n.lastReceiveTimes[i] > max {
			max = cur
		}
	}
	mean := sum / len(n.lastReceiveTimes)

	//fmt.Printf("N:%d, min:%d, max:%d, spread: %d, mean:%d\n", n.id, min, max, max-min, mean)

	//spread := (max-min)/4 + interval/4 // 50:50 mix of actual spread and intended interval

	// The spread between the first and last message should be (N-1)/N of the interval.
	spread := interval * (n.N - 1) / n.N
	target := mean - spread/2 + should

	delta := (target - n.lastEmitTime) / 2

	/*
		if n.id == 0 {
			fmt.Printf("N%d: mean %d, target %d, delta %d, last-emit: %d, next-emit: %d\n", n.id, mean, target, delta, n.lastEmitTime, n.lastEmitTime+delta)
		}
	*/

	n.lastEmitTime += delta
	return n.lastEmitTime
}

func (n *Node) Receive(time int, msg Message) {
	n.lastReceiveTimes[msg.from] = time
}
