package paxos

// Vote is an internal structure used by multiple components to cast votes on behalf of the acceptor that they're
// communicating with.
type Vote struct {
	// Member contains which member of the cluster cast this vote.
	Member string
	// Payload contains the payload of the message we're voting on. This is usually a Promise or Proposal.
	Payload interface{}
}
