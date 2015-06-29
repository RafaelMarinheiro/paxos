package paxos

type Role interface {
	AcceptMessage(Message) bool
	Run(in <-chan Message, out chan<- Message)
}
