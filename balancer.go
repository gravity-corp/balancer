package balancer

import (
	"errors"
	"sync/atomic"
)

type Balancer struct {
	round round
	smart *smart
}

// NewBalancer creates and returns a load balancer
func NewBalancer(addrs []string) (bal *Balancer) {
	table := make([]entry, 0, len(addrs))
	var k int64 = 1
	for _, addr := range addrs {
		table = append(table, entry{address: addr, weight: k})
		k++
	}

	bal = &Balancer{
		round: round{index: 0, addresses: addrs},
		smart: &smart{
			to:      make(chan entry, 1),
			from:    make(chan entry, 1),
			enable:  make(chan string, 1),
			disable: make(chan struct{}, 1),
			table:   table,
		},
	}

	go bal.analyze()

	return bal
}

// Next returns the next backend address
func (b *Balancer) Next() (addr string, err error) {
	select {
	case entry := <-b.smart.from:
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

// Score writes the result to a table for future analysis
func (b *Balancer) Score(addr string, weight int64) {
	select {
	case b.smart.to <- entry{address: addr, weight: weight}:
	}
}

// Disable removes an address from the address list
// and reduces the results table
func (b *Balancer) Disable(addr string) {
	select {
	case b.smart.disable <- struct{}{}:
	}

	b.round.remove(addr)
}

// Enable adds an address to the address list
// and results table
func (b *Balancer) Enable(addr string) {
	select {
	case b.smart.enable <- addr:
	}

	b.round.add(addr)
}

// Analize calculates the most unloaded backend
// and returns its address
func (b *Balancer) analyze() {
	for {
		select {
		case entry := <-b.smart.to:
			i := b.smart.close(entry.weight)
			if i > -1 {
				entry = b.smart.push(i, entry)
				b.smart.from <- entry
			}
		case <-b.smart.disable:
			b.smart.cut()
		case addr := <-b.smart.enable:
			b.smart.enter(addr)
		}
	}
}
