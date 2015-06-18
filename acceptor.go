package paxos

type Acceptor struct{
	node nodeinfo
	highestPrepare int64
	highestAccepted int64
	valueAccepted int64
}

func (Acceptor *) Run(net Network){
	
} 