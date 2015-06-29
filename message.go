package paxos

import "fmt"

type Addr uint

type MessageType int

const (
	propose       MessageType = iota // Ask the proposer to propose something
	prepare                          // Proposer asks Acceptor to prepare
	promise                          // Acceptor promises Proposer
	refusePromise                    // Acceptor denies Proposal
	acceptRequest                    // Proposer request Acceptor to learn avalue
	accept                           // Acceptor accepts value and broadcast the message
	refuseAccept                     // Acceptor denies Accept
)

func (m MessageType) String() string {
	switch m {
	case propose:
		return "propose"
	case prepare:
		return "prepare"
	case promise:
		return "promise"
	case acceptRequest:
		return "acceptRequest"
	case accept:
		return "accept"
	default:
		return "???"
	}
}

type Proposal struct {
	promise int64
	number  int64
	key     string
	value   string
}

func nullProposal(key string) Proposal {
	return Proposal{promise: -1, number: -1, key: key}
}

type Message struct {
	From     Addr
	To       Addr
	mtype    MessageType
	proposal Proposal
}

func (p Proposal) String() string {
	return fmt.Sprintf("Proposal(%d, %d, %s, %s)", p.promise, p.number, p.key, p.value)
}

func (m Message) String() string {
	return fmt.Sprintf("Message(From: %d, To: %d\n\ttype: %s\n\t%s", m.From, m.To, m.mtype.String(), m.proposal.String())
}
