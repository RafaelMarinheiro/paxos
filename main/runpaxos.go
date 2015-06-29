package main

import "fmt"
import "time"
import "math/rand"
import "strconv"
import "github.com/RafaelMarinheiro/paxos"

func main() {
	rand.Seed(42)
	n := 29
	fnet := paxos.CreateFakeNetwork(n)
	go fnet.Start()

	interfaces := make([]paxos.FakeNetworkInterface, 0)
	for i := 0; i < n; i++ {
		net := paxos.CreateFakeNetworkInterface(fnet, i+1)

		interfaces = append(interfaces, net)

		go net.Run()
	}

	nodes := make([]*paxos.PaxosNode, 0)
	for i := 0; i < n; i++ {
		node := paxos.CreatePaxosNode()
		nodes = append(nodes, node)

		go node.Start(interfaces[i])
	}

	for iteration := 0; iteration < 2000; iteration++ {
		key := "leader" + strconv.Itoa(iteration)
		perm := rand.Perm(n)
		for _, no := range perm {
			time.Sleep(1 * time.Millisecond)
			go nodes[no].Propose(key, strconv.Itoa(no))
		}
		learned := nodes[0].Wait(key)

		fmt.Println("Learned: " + key + " = " + learned)
	}

	fnet.Stop()
}
