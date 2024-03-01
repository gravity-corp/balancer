package balancer

import (
	"errors"
	"sync/atomic"
)

type RoundBalancer struct {
	round *round
}

// NewBalancer creates round balancer
func NewRoundBalancer(addrs []string) (b *RoundBalancer) {
	return &RoundBalancer{
		round: newRound(addrs),
	}
}

// Next returns next backend address
func (b *RoundBalancer) Next() (addr string, err error) {
	if len(b.round.addresses) > 0 {
		i := atomic.AddUint64(&b.round.index, 1) % uint64(len(b.round.addresses))
		return b.round.addresses[i], nil
	} else {
		return "", errors.New("no addresses are tracked")
	}
}

// Disable removes address from address list
// and reduces results table
func (b *RoundBalancer) Disable(addr string) {
	b.round.remove(addr)
}

// Enable adds address to address list
// and results table
func (b *RoundBalancer) Enable(addr string) {
	b.round.add(addr)
}
