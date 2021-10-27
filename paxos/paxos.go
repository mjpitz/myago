package paxos

type Paxos struct {
	AcceptorServer
	ProposerServer
	ProposerClient
}

type LeaderPaxos struct {
	Paxos
}