package paxos

// import "fmt"

type proposerData struct {
	key              string
	myValue          string
	myProposalNumber int64
	acceptedRequest  int64
	acceptedValue    string
	promised         map[Addr]bool
	npromised        int
	sentAccept       bool
}

type Proposer struct {
	numNodes        int
	baseNumber      int64
	incrementNumber int64
	data            map[string]proposerData
}

func quorum(k int) int {
	return k/2 + 1
}

func nullProposerData(key string, baseNumber int64) proposerData {
	return proposerData{key: key, myProposalNumber: baseNumber, acceptedRequest: -1}
}

func (proposer *Proposer) AcceptMessage(message Message) bool {
	switch message.mtype {
	case propose:
		return true
	case promise:
		return true
	case accept:
		return true
	default:
		return false
	}
}

func (proposer *Proposer) Run(in <-chan Message, out chan<- Message) {
	for m := range in {
		proposal := m.proposal

		if _, exists := proposer.data[proposal.key]; !exists {
			proposer.data[proposal.key] = nullProposerData(proposal.key, proposer.baseNumber)
		}

		current := proposer.data[proposal.key]

		switch m.mtype {
		case propose:
			current.myValue = proposal.value
			current.myProposalNumber += proposer.incrementNumber
			current.promised = make(map[Addr]bool)

			proposer.data[proposal.key] = current
			//Send proposal

			out <- Message{
				To:    Addr(0), //Broadcast
				mtype: prepare,
				proposal: Proposal{key: current.key,
					number: current.myProposalNumber}}
		case promise:
			if proposal.number == current.myProposalNumber {
				if proposal.promise > current.acceptedRequest {
					current.acceptedRequest = proposal.promise
					current.acceptedValue = proposal.value
				}

				if _, hasReplied := current.promised[m.From]; !hasReplied {
					current.promised[m.From] = true
					current.npromised++
					// fmt.Println(current.npromised)
				}

				willSend := false

				if current.npromised >= quorum(proposer.numNodes) && !current.sentAccept {
					willSend = true
					current.sentAccept = true
				}

				// save
				proposer.data[proposal.key] = current

				// send
				if willSend == true {
					daValue := current.myValue
					if current.acceptedRequest >= 0 {
						daValue = current.acceptedValue
					}

					out <- Message{
						To:    Addr(0), //Broadcast
						mtype: acceptRequest,
						proposal: Proposal{key: current.key,
							number: current.myProposalNumber,
							value:  daValue}}
				}

			}

		case accept:
			//nothing

		}
	}
	// close(out)
}
