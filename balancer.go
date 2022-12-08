package balancer

import (
	"errors"
	"sync/atomic"
)

type Balancer struct {
	round round
	smart *smart
}

type round struct {
	index     uint64
	addresses []string
}

type smart struct {
	from  chan entry
	to    chan entry
	table []entry
}

type entry struct {
	address string
	weight  int64
}

// Returns new balancer
func NewBalancer(addrs []string) (bal *Balancer, err error) {
	if len(addrs) < 2 {
		return nil, errors.New("at least 2 addresses were expected")
	}

	table := make([]entry, 0, len(addrs))
	var k int64
	for _, addr := range addrs {
		table = append(table, entry{address: addr, weight: k})
		k++
	}

	bal = &Balancer{
		round: round{index: 0, addresses: addrs},
		smart: &smart{to: make(chan entry, 1), from: make(chan entry, 1), table: table},
	}

	go bal.analyze()

	return bal, nil
}

// Returns the next backend address
func (b *Balancer) Next() (addr string) {
	select {
	case entry := <-b.smart.from:
		return entry.address
	default:
		i := atomic.AddUint64(&b.round.index, 1) % uint64(len(b.round.addresses))
		return b.round.addresses[i]
	}
}

// Sends the result to the backend weight analyzer
func (b *Balancer) Score(addr string, weight int64) {
	select {
	case b.smart.to <- entry{address: addr, weight: weight}:
	}
}

// Analyzes the weight of the server and sends
// one with the lowest weight
func (b *Balancer) analyze() {
	for {
		select {
		case entry := <-b.smart.to:

			i := b.smart.close(entry.weight)
			entry = b.smart.push(i, entry)

			b.smart.from <- entry
		}
	}
}

// Looking for a backend with the same weight or close to it
func (s *smart) close(weight int64) (inx int) {
	first := 0
	last := len(s.table) - 1

	for first <= last {
		mid := first + (last-first)/2

		if first == mid {
			return first
		}

		if first == last {
			return first
		}

		if last == mid {
			return last
		}

		if s.table[mid].weight < weight {
			first = mid + 1
		} else {
			last = mid - 1
		}
	}

	return
}

// Insert the result into the table
func (s *smart) push(i int, ent entry) (out entry) {
	out = s.table[0]

	switch {
	case i == 0:
		s.table[i] = ent
	case i == len(s.table)-1:
		switch {
		case s.table[i].weight < ent.weight:
			buf := make([]entry, len(s.table))
			copy(buf, s.table[1:])
			buf[i] = ent

			s.table = buf
		case ent.weight < s.table[i].weight:
			buf := make([]entry, len(s.table))
			copy(buf, s.table[1:i])
			buf[i-1] = ent
			buf[i] = s.table[i]

			s.table = buf
		case ent.weight == s.table[i].weight:
			s.table[i] = ent
		}
	default:
		switch {
		case ent.weight < s.table[i].weight:
			switch {
			case s.table[i-1].weight < ent.weight:
				buf := make([]entry, len(s.table))
				copy(buf, s.table[1:i])

				buf[i-1] = ent

				copy(buf[i:], s.table[i:])

				s.table = buf

			case s.table[i-1].weight == ent.weight:
				buf := make([]entry, len(s.table))
				copy(buf, s.table[1:i])

				ent.weight++
				buf[i-1] = ent

				for ; i < len(buf); i++ {
					e := s.table[i]
					e.weight++

					buf[i] = e
				}

				s.table = buf
			}
		case s.table[i].weight < ent.weight:
			switch {
			case ent.weight < s.table[i+1].weight:
				buf := make([]entry, len(s.table))
				copy(buf, s.table[1:i+1])

				buf[i] = ent

				copy(buf[i+1:], s.table[i+1:])

				s.table = buf
			case s.table[i+1].weight == ent.weight:
				buf := make([]entry, len(s.table))
				copy(buf, s.table[1:i+1])

				ent.weight++
				buf[i] = ent

				for i++; i < len(buf); i++ {
					e := s.table[i]
					e.weight++

					buf[i] = e
				}

				s.table = buf
			}
		}
	}

	return out
}
