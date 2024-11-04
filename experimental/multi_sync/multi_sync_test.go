package multi_sync

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
)

const interval = 10000
const withDelay = true
const withSymetricDelay = true
const withNoise = false

func TestSimulateSync(t *testing.T) {
	const N = 10
	const T = 100
	runSync(N, T, t)
}

func runSync(N, T int, t *testing.T) {

	times := []string{}
	estOffsetError := []string{}
	estDelayError := []string{}

	// Create offsets for N nodes (secret).
	offsets := make([]int, N)
	for i := 0; i < N; i++ {
		offsets[i] = rand.Intn(interval)
	}
	for i, offset := range offsets {
		fmt.Printf("N%d: %d\n", i, offset)
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
					// TODO: define a maximum delay as a constant/parameter
					delay[i][j] = rand.Intn(interval/10) + 1 // up to 10% of interval
				}
			}
		}
	}

	if withSymetricDelay {
		for i := range N {
			for j := range N {
				delay[i][j] = delay[j][i]
			}
		}
	}

	builder := &strings.Builder{}
	builder.WriteString("Iteration,")
	for i := range N {
		builder.WriteString(fmt.Sprintf("N%d,", i))
	}
	builder.WriteString("\n")
	times = append(times, builder.String())

	builder = &strings.Builder{}
	builder.WriteString("Iteration,")
	for i := range N {
		builder.WriteString(fmt.Sprintf("N%d,", i))
	}
	builder.WriteString("\n")
	estOffsetError = append(estOffsetError, builder.String())

	builder = &strings.Builder{}
	builder.WriteString("Iteration,")
	for i := range N {
		for j := range N {
			builder.WriteString(fmt.Sprintf("%d-%d,", i, j))
		}
	}
	builder.WriteString("\n")
	estDelayError = append(estDelayError, builder.String())

	// Run simulation.
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

		{
			builder := &strings.Builder{}
			builder.WriteString(fmt.Sprintf("%d,", step))
			for i := range N {
				builder.WriteString(fmt.Sprintf("%d,", emitTimes[i]+interval))
			}
			builder.WriteString("\n")
			times = append(times, builder.String())
		}

		// Print the offset estimate errors.
		{
			builder := &strings.Builder{}
			builder.WriteString(fmt.Sprintf("%d,", step))
			for i := range N {
				is := offsets[i]
				est := nodes[i].GetEstimatedGlobalOffset()
				diff := est - is
				//builder.WriteString(fmt.Sprintf("%d-%d = %d,", est, is, diff))
				builder.WriteString(fmt.Sprintf("%d,", diff))
			}
			builder.WriteString("\n")
			estOffsetError = append(estOffsetError, builder.String())
		}

		// Print the delay estimate errors.
		{
			builder := &strings.Builder{}
			builder.WriteString(fmt.Sprintf("%d,", step))
			for i := range N {
				for j := range N {
					is := delay[i][j]
					est := nodes[j].delayEstimate[i]
					diff := est - is
					//fmt.Printf("%d,%d,%d,", is, est, diff)
					//fmt.Printf("%d-%d=%d,", is, est, diff)
					builder.WriteString(fmt.Sprintf("%d,", diff))
					//fmt.Printf("N%d->N%d: is %d, est %d\n", i, j, delay[i][j], nodes[j].delayEstimate[i])
				}
			}
			builder.WriteString("\n")
			estDelayError = append(estDelayError, builder.String())
		}
	}

	if true {
		for _, line := range times {
			fmt.Print(line)
		}
		fmt.Printf("\n")
	}
	if true {
		for _, line := range estOffsetError {
			fmt.Print(line)
		}
		fmt.Printf("\n")
	}
	if true {
		for _, line := range estDelayError {
			fmt.Print(line)
		}
		fmt.Printf("\n")
	}

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

func (n *Node) GetEstimatedGlobalOffset() int {
	// compute the offset this node should emit its message
	should := interval * n.id / n.N

	// estimated local offset to global time
	return n.lastEmitTime - should
}

func (n *Node) GetNextEmitTime() int {
	// estimated local offset to global time
	offset := n.GetEstimatedGlobalOffset()

	// refresh the delay estimate
	{
		// collect the observed delays of the last iteration
		lastDelays := make([]int, n.N)
		for i := range n.N {
			expectedSend := interval * i / n.N
			lastReceived := n.lastReceiveTimes[i]
			lastDelays[i] = lastReceived - expectedSend - offset
		}

		// align the observed delays to minimize the impact on the offset
		// Note: part of the offset could 'hide' in the observed delays
		min := math.MaxInt
		for i := range n.N {
			if i != n.id {
				if lastDelays[i] < min {
					min = lastDelays[i]
				}
			}
		}
		//fmt.Printf("N: %d, min: %d\n", n.id, min)
		if min < math.MaxInt {
			min -= 1 // minimum expected delay
			for i := range n.N {
				lastDelays[i] -= min
			}
		}

		// update the delay estimates as a running average
		for i := range n.N {
			n.delayEstimate[i] = (lastDelays[i]*1 + n.delayEstimate[i]*2) / 3
		}
	}

	// compute the mean time of all received messages
	sum := 0
	for i := range n.lastReceiveTimes {
		cur := n.lastReceiveTimes[i]
		cur = cur - n.delayEstimate[i]
		sum += cur
	}
	mean := sum / len(n.lastReceiveTimes)

	// compute the offset this node should emit its message
	should := interval * n.id / n.N

	// The spread between the first and last message should be (N-1)/N of the interval.
	spread := interval * (n.N - 1) / n.N
	target := mean - spread/2 + should

	delta := (target - n.lastEmitTime) / 2

	n.lastEmitTime += delta
	return n.lastEmitTime
}

func (n *Node) Receive(time int, msg Message) {
	n.lastReceiveTimes[msg.from] = time
}
