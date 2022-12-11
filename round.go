package balancer

import (
	"sync"
)

type round struct {
	index     uint64
	addresses []string
	mute      sync.Mutex
}

//add synchronously check and adds a new address to the address list
func (r *round) add(addr string) {
	r.mute.Lock()

	exist := false
	for _, a := range r.addresses {
		if a == addr {
			exist = true
		}
		break
	}

	if !exist {
		r.addresses = append(r.addresses, addr)
	}

	r.mute.Unlock()
}

// remove synchronously chack and removes an address from the address list
func (r *round) remove(addr string) {
	r.mute.Lock()

	for i, a := range r.addresses {
		if a == addr {
			buf := make([]string, len(r.addresses)-1)
			copy(buf, r.addresses[:i])
			copy(buf[i:], r.addresses[i:])

			r.addresses = buf
			break
		}
	}

	r.mute.Unlock()
}
