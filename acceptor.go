package paxos

// import "fmt"

type acceptorData struct {
	key              string
	value            string
	promisedProposal int64
	acceptedProposal int64
	accepted         bool
}

type Acceptor struct {
	data map[string]acceptorData
}

func nullAcceptorData(key string) acceptorData {
	return acceptorData{key: key, promisedProposal: -1, acceptedProposal: -1}
}

func (acceptor *Acceptor) AcceptMessage(message Message) bool {
	switch message.mtype {
	case prepare:
		return true
	case acceptRequest:
		return true
	default:
		return false
	}
}

func (acceptor *Acceptor) Run(in <-chan Message, out chan<- Message) {
	for m := range in {

		proposal := m.proposal

		if _, exists := acceptor.data[proposal.key]; !exists {
			acceptor.data[proposal.key] = nullAcceptorData(proposal.key)
		}

		current := acceptor.data[proposal.key]

		switch m.mtype {
		case prepare:
			if proposal.number > current.promisedProposal {
				current.promisedProposal = proposal.number
				//Save
				acceptor.data[proposal.key] = current

				//Reply message

				out <- Message{
					To:    m.From,
					mtype: promise,
					proposal: Proposal{
						key:     current.key,
						number:  current.promisedProposal,
						promise: current.acceptedProposal,
						value:   current.value}}
			}

		case acceptRequest:
			if proposal.number >= current.promisedProposal && proposal.number != current.acceptedProposal {
				current.promisedProposal = proposal.number
				current.acceptedProposal = proposal.number
				current.value = proposal.value

				//Save
				acceptor.data[proposal.key] = current

				//Reply and send to learners
				out <- Message{
					To:    Addr(0), //Broadcast
					mtype: accept,
					proposal: Proposal{
						key:     current.key,
						number:  current.promisedProposal,
						promise: current.acceptedProposal,
						value:   current.value}}
			}
		}
	}
	// close(out)
}
