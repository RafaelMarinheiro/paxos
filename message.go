package paxos

type MessageType int

const(
	_ , _ MessageType = iota, -iota
	prepare, rejectPrepare
	promise, rejectPromise
	accept, _
	accepted, notAccepted
	propose, rejectPropose
)

type nodeinfo struct{
	addr int
}

type Message struct{
	From int
	To int
	mtype MessageType
	proposeNumber int64
	proposeValue int64
}

type Network interface{
	//Send a message to a node
	Send(Message)

	//Receive a message from a noide
	Receive() Message
}