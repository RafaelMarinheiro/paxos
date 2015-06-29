package paxos

type PaxosNode struct {
	shouldStop  bool
	proposeChan chan Proposal
	learned     map[string]chan string
}

func CreatePaxosNode() *PaxosNode {
	return &PaxosNode{proposeChan: make(chan Proposal), learned: make(map[string]chan string)}
}

func (n *PaxosNode) Start(net NetworkInterface) {
	n.run(net)
}

func (n *PaxosNode) run(net NetworkInterface) {
	chans := make([]chan Message, 4)
	for i := 0; i < 4; i++ {
		chans[i] = make(chan Message)
	}

	defer func() {
		for _, c := range chans {
			close(c)
		}
	}()

	roles := make([]Role, 0)

	proposer := &Proposer{numNodes: int(net.NumNodes()),
		baseNumber:      int64(net.Address()),
		incrementNumber: net.NumNodes(),
		data:            make(map[string]proposerData)}

	acceptor := &Acceptor{data: make(map[string]acceptorData)}

	learner := &Learner{numNodes: int(net.NumNodes()),
		data:      make(map[string]learnerData),
		learnChan: n.learned}

	roles = append(roles, proposer)
	roles = append(roles, acceptor)
	roles = append(roles, learner)

	for i, _ := range roles {
		go roles[i].Run(chans[i], chans[3])
	}

	go func() {
		for m := range chans[3] {
			m.From = net.Address()
			net.OutputChannel() <- m
		}
	}()

	for {
		if n.shouldStop {
			return
		}

		var m Message
		select {
		case prop := <-n.proposeChan:
			m = Message{mtype: propose,
				proposal: prop}
		case m = <-net.InputChannel():
		}

		for i, _ := range roles {
			m.To = net.Address()
			if roles[i].AcceptMessage(m) {
				chans[i] <- m
			}
		}
	}
}

func (n *PaxosNode) Propose(key string, value string) {
	if _, exists := n.learned[key]; !exists {
		n.learned[key] = make(chan string, 10)
	}
	n.proposeChan <- Proposal{key: key, value: value}
}

func (n *PaxosNode) Status(key string) bool {
	return true
}

func (n *PaxosNode) Wait(key string) string {
	if _, exists := n.learned[key]; !exists {
		n.learned[key] = make(chan string, 10)
	}
	return <-n.learned[key]
}
