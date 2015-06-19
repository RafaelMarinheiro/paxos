package paxos

type Addr uint;

type MessageType int;

const (
	prepare MessageType = iota
	accept
)

type Proposal struct{
	promise int64
	number int64
	key string
	value string
}

func nullProposal(key string) Proposal{
	return Proposal{promise: -1, number: -1, key: key}
}

type Message struct{
	From Addr
	To Addr
	mtype MessageType
	proposal Proposal
}