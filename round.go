package balancer

import (
	"sync"
)

type round struct {
	sync.RWMutex
	index     uint64
	addresses []string
}

// newRound creates round
func newRound(addrs []string) (r *round) {
	return &round{
		addresses: addrs,
	}
}

// add adds a new address to address list
func (r *round) add(addr string) {
	defer r.Unlock()
	r.Lock()

	exist := false
	for _, a := range r.addresses {
		if a == addr {
			exist = true
			break
		}
	}

	if !exist {
		r.addresses = append(r.addresses, addr)
	}
}

// remove removes address from address list
func (r *round) remove(addr string) {
	defer r.Unlock()
	r.Lock()

	for i, a := range r.addresses {
		if a == addr {
			buf := make([]string, len(r.addresses)-1)
			copy(buf, r.addresses[:i])
			copy(buf[i:], r.addresses[i:])

			r.addresses = buf
			break
		}
	}
}
