package paxos

// import "fmt"

type learnerData struct {
	key        string
	round      int64
	votes      map[string]int
	hasVoted   map[Addr]bool
	hasLearned bool
}

type Learner struct {
	numNodes  int
	learnChan map[string]chan string
	data      map[string]learnerData
}

func nullLearnerData(key string) learnerData {
	return learnerData{key: key, round: -1}
}

func (learner *Learner) AcceptMessage(message Message) bool {
	switch message.mtype {
	case accept:
		return true
	default:
		return false
	}
}

func (learner *Learner) Run(in <-chan Message, out chan<- Message) {
	for m := range in {
		proposal := m.proposal

		if _, exists := learner.data[proposal.key]; !exists {
			learner.data[proposal.key] = nullLearnerData(proposal.key)
		}

		current := learner.data[proposal.key]

		switch m.mtype {
		case accept:
			// fmt.Println("Aqui!: "+ m.String())

			if proposal.number > current.round {
				current.votes = make(map[string]int)
				current.hasVoted = make(map[Addr]bool)
				current.hasLearned = false
				current.round = proposal.number
			}

			hasLearned := false

			if !current.hasLearned {
				if proposal.number == current.round {
					if _, exists := current.hasVoted[m.From]; !exists {
						current.hasVoted[m.From] = true

						if _, valueExists := current.votes[proposal.value]; !valueExists {
							current.votes[proposal.value] = 0
						}
						current.votes[proposal.value]++
					}

					if current.votes[proposal.value] >= quorum(learner.numNodes) {
						current.hasLearned = true
						hasLearned = true
					}
				}
			}

			learner.data[proposal.key] = current

			if hasLearned {
				if _, exists := learner.learnChan[proposal.key]; !exists {
					learner.learnChan[proposal.key] = make(chan string, 10)
				}
				learner.learnChan[proposal.key] <- proposal.value
			}

		}
	}
	// close(out)
}
