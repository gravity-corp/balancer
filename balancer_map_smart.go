package balancer

import (
	"errors"
	"sync"
)

type MapSmartBalancer struct {
	sync.RWMutex
	table map[any]*SmartBalancer
}

// NewBalancer create and returns a load smart balancer map
func NewMapSmartBalancer(addrs map[any][]string) (mb *MapSmartBalancer) {
	mb = &MapSmartBalancer{table: make(map[any]*SmartBalancer)}
	for key, value := range addrs {
		mb.table[key] = NewSmartBalancer(value)
	}

	return mb
}

// NextRange accepts handle func as an argument
// and renges next addresses by keys
func (mb *MapSmartBalancer) NextRange(f func(k any, addr string) bool) (err error) {
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
func (mb *MapSmartBalancer) Next(key any) (addr string, err error) {
	defer mb.Unlock()
	mb.Lock()

	if b, ok := mb.table[key]; ok {
		return b.Next()
	}

	return "", errors.New("no addresses are tracked")
}

// Score writes the result to a table by key for future analysis
func (mb *MapSmartBalancer) Score(key any, addr string, weight int64) {
	defer mb.Unlock()
	mb.Lock()

	if b, ok := mb.table[key]; ok {
		b.Score(addr, weight)
	}
}

// Disable removes an address from the address list
// by key and reduces the results table
func (mb *MapSmartBalancer) Disable(key any, addr string) {
	defer mb.Unlock()
	mb.Lock()

	if b, ok := mb.table[key]; ok {
		b.Disable(addr)
	}
}

// Enable adds an address to the address list
// by key and results table
func (mb *MapSmartBalancer) Enable(key any, addr string) {
	defer mb.Unlock()
	mb.Lock()

	if b, ok := mb.table[key]; ok {
		b.Enable(addr)
	}
}
