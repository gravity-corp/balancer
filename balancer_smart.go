package balancer

import (
	"errors"
	"sync/atomic"
)

type SmartBalancer struct {
	round *round
	smart *smart
}

// NewBalancer creates smart balancer
func NewSmartBalancer(addrs []string) (b *SmartBalancer) {
	return &SmartBalancer{
		round: newRound(addrs),
		smart: newSmart(addrs),
	}
}

// Next returns next backend address
func (b *SmartBalancer) Next() (addr string, err error) {
	select {
	case entry := <-b.smart.feed:
		return entry.address, nil
	default:
		if len(b.round.addresses) > 0 {
			i := atomic.AddUint64(&b.round.index, 1) % uint64(len(b.round.addresses))
			return b.round.addresses[i], nil
		} else {
			return "", errors.New("no addresses are tracked")
		}
	}
}

// Score writes result to table for future analysis
func (b *SmartBalancer) Score(addr string, weight int64) {
	entry := entry{address: addr, weight: weight}
	entry, free := b.smart.push(entry)
	if !free {
		return
	}

	select {
	case b.smart.feed <- entry:
	default:
		<-b.smart.feed
		b.smart.feed <- entry
	}
}

// Disable removes address from address list
// and reduces results table
func (b *SmartBalancer) Disable(addr string) {
	b.smart.cut()
	b.round.remove(addr)
}

// Enable adds address to address list
// and results table
func (b *SmartBalancer) Enable(addr string) {
	b.smart.add(addr)
	b.round.add(addr)
}
