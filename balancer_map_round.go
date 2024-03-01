package balancer

import (
	"errors"
	"sync"
)

type MapRoundBalancer struct {
	sync.RWMutex
	table map[any]*RoundBalancer
}

// NewBalancer create and returns a load round balancer map
func NewMapRoundBalancer(addrs map[any][]string) (mb *MapRoundBalancer) {
	mb = &MapRoundBalancer{table: make(map[any]*RoundBalancer)}
	for key, value := range addrs {
		mb.table[key] = NewRoundBalancer(value)
	}

	return mb
}

// NextRange accepts handle func as an argument
// and renges next addresses by keys
func (mb *MapRoundBalancer) NextRange(f func(k any, addr string) bool) (err error) {
	defer mb.Unlock()
	mb.Lock()

	for key, b := range mb.table {
		addr, err := b.Next()
		if err != nil {
			break
		}

		if !f(key, addr) {
			break
		}
	}

	return err
}

// Next returns the next backend address by key
func (mb *MapRoundBalancer) Next(key any) (addr string, err error) {
	defer mb.Unlock()
	mb.Lock()

	if b, ok := mb.table[key]; ok {
		return b.Next()
	}

	return "", errors.New("no addresses are tracked")
}

// Disable removes an address from the address list
// by key and reduces the results table
func (mb *MapRoundBalancer) Disable(key any, addr string) {
	defer mb.Unlock()
	mb.Lock()

	if b, ok := mb.table[key]; ok {
		b.Disable(addr)
	}
}

// Enable adds an address to the address list
// by key and results table
func (mb *MapRoundBalancer) Enable(key any, addr string) {
	defer mb.Unlock()
	mb.Lock()

	if b, ok := mb.table[key]; ok {
		b.Enable(addr)
	}
}
