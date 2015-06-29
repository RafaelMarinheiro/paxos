package paxos

// import "fmt"
import "math/rand"

type NetworkInterface interface {
	NumNodes() int64
	Address() Addr
	InputChannel() <-chan Message
	OutputChannel() chan<- Message
	Run()
}

type FakeNetwork struct {
	numNodes int
	inputs   []chan Message
	output   chan Message
}

func CreateFakeNetwork(numNodes int) FakeNetwork {
	in := make([]chan Message, numNodes)
	for i := 0; i < numNodes; i++ {
		in[i] = make(chan Message, 100)
	}
	out := make(chan Message, 100)

	return FakeNetwork{numNodes: numNodes, inputs: in, output: out}
}

func (fnet FakeNetwork) Stop() {
	close(fnet.output)
}

func (net FakeNetwork) Start() {
	defer func() {
		for _, c := range net.inputs {
			close(c)
		}
	}()

	sendMessage := func(m Message) {
		if m.To == Addr(0) { //Broadcast
			for _, c := range net.inputs {
				// m.To = Addr(i+1)
				c <- m
			}
		} else {
			net.inputs[int(m.To)-1] <- m
		}
	}

	for m := range net.output {
		select {
		case m2 := <-net.output:
			chance := rand.Float32()

			if chance < 0.1 {
				sendMessage(m2)
				sendMessage(m)
			} else {
				sendMessage(m)
				sendMessage(m2)
			}
		default:
			sendMessage(m)
		}
		// fmt.Println(m)

	}
}

type FakeNetworkInterface struct {
	fnet *FakeNetwork
	me   int
}

func (net FakeNetworkInterface) NumNodes() int64 {
	return int64(net.fnet.numNodes)
}

func (net FakeNetworkInterface) Address() Addr {
	return Addr(net.me)
}

func (net FakeNetworkInterface) InputChannel() <-chan Message {
	return net.fnet.inputs[net.me-1]
}

func (net FakeNetworkInterface) OutputChannel() chan<- Message {
	return net.fnet.output
}

func (net FakeNetworkInterface) Run() {

}

func CreateFakeNetworkInterface(fnet FakeNetwork, me int) FakeNetworkInterface {
	return FakeNetworkInterface{fnet: &fnet, me: me}
}
