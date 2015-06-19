package paxos

type Acceptor struct{
	data map[string]Proposal
}

func (acceptor * Acceptor) AcceptMessage(message Message) bool{
	switch message.mtype{
	case prepare:
		return true
	case accept:
		return true;
	default:
		return false
	}
}

func (acceptor * Acceptor) Run(in <-chan Message, out chan<- Message){
	for m := range in {
		proposal := m.proposal

		if _, exists := acceptor.data[proposal.key]; !exists{
			acceptor.data[proposal.key] = nullProposal(proposal.key)
		}

		current := acceptor.data[proposal.key];

		switch m.mtype{
		case prepare:
			if proposal.number > current.promise{
				current.promise = proposal.number
				//Save
				acceptor.data[proposal.key] = current;
				
				//Reply message
			}

		case accept:
			if proposal.number >= current.promise && proposal.number != current.number{
				current.promise = proposal.number
				current.number = proposal.number
				current.value = proposal.value

				//Save
				acceptor.data[proposal.key] = current;

				//Reply and send to learners
			}
		}
	}
	close(out)
}