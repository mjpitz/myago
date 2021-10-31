package paxos

import (
	"context"
	"fmt"
	"sync"

	"github.com/jonboulle/clockwork"
	"github.com/mjpitz/myago/cluster"
	"golang.org/x/sync/errgroup"
)

// Config contains configurable elements of Paxos.
type Config struct {
	Clock          clockwork.Clock
	IDGenerator    IDGenerator
	PromiseLog     Log
	AcceptedLog    Log
	RecordedLog    Log
	AcceptorDialer func(ctx context.Context, member string) (AcceptorClient, error)
	ObserverDialer func(ctx context.Context, member string) (ObserverClient, error)
}

// Validate ensures the configuration is valid.
func (c *Config) Validate() error {
	if c.Clock == nil {
		c.Clock = clockwork.NewRealClock()
	}

	switch {
	case c.IDGenerator == nil:
		return fmt.Errorf("missing id generator")
	case c.PromiseLog == nil:
		return fmt.Errorf("missing promise log")
	case c.AcceptedLog == nil:
		return fmt.Errorf("missing promise log")
	case c.RecordedLog == nil:
		return fmt.Errorf("missing recorded log")
	case c.AcceptorDialer == nil:
		return fmt.Errorf("missing acceptor dialer")
	case c.ObserverDialer == nil:
		return fmt.Errorf("missing observer dialer")
	}
	return nil
}

// New constructs a new instance of paxos given the provided configuration. It returns an error should the provided
// configuration be invalid.
func New(cfg *Config) (*Paxos, error) {
	acceptor, err := NewAcceptor(cfg.PromiseLog, cfg.AcceptedLog)
	if err != nil {
		return nil, err
	}

	return &Paxos{
		Proposer: Proposer{
			Clock:       cfg.Clock,
			IDGenerator: cfg.IDGenerator,
			Acceptor: &MultiAcceptorClient{
				Dialer: cfg.AcceptorDialer,
				cache:  &sync.Map{},
			},
		},
		Observer: Observer{
			Dialer: cfg.ObserverDialer,
			Log:    cfg.RecordedLog,
		},
		Acceptor: acceptor,
	}, nil
}

// Paxos defines the core elements of a paxos participant.
type Paxos struct {
	// Proposer contains the logic required to propose changes to the paxos state machine. Any member in paxos can act
	// as a proposer. Proposers communicate with all acceptor to propose changes to the log.
	Proposer

	// Observer contains the logic required to be an observer of the paxos protocol. Every member in paxos _must_ be an
	// observer. Observers watch all acceptor to learn about the records they've accepted.
	Observer

	// Acceptor must implement the functionality of an AcceptorServer and an ObserverServer. The ObserverServer is how
	// other members of the cluster learn about changes.
	Acceptor
}

func (p *Paxos) Start(ctx context.Context, membership *cluster.Membership) error {
	acceptor, ok := p.Proposer.Acceptor.(*MultiAcceptorClient)
	if !ok {
		return p.Observer.Start(ctx, membership)
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return acceptor.Start(ctx, membership)
	})

	group.Go(func() error {
		return p.Observer.Start(ctx, membership)
	})

	return group.Wait()
}

var _ cluster.Discovery = &Paxos{}
