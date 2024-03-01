package balancer

import (
	"sync"
)

type smart struct {
	feed  chan entry
	table table
}

type table struct {
	sync.RWMutex
	entrys []entry
}

type entry struct {
	address string
	weight  int64
}

// newRound creates smart
func newSmart(addrs []string) (s *smart) {
	entrys := make([]entry, 0, len(addrs))
	for weight, addr := range addrs {
		entrys = append(entrys, entry{address: addr, weight: int64(weight + 1)})
	}

	return &smart{
		feed:  make(chan entry, 1),
		table: table{entrys: entrys},
	}
}

// close finds similar backend weight or one closest to it
// and returns its index
//
// Note: Method might be called only in push method
func (s *smart) close(weight int64) (index int) {
	mid := -1
	left := 0
	right := len(s.table.entrys) - 1

	for left <= right {
		mid = int((left + right) / 2)
		wt := s.table.entrys[mid].weight

		if wt == weight {
			return mid
		}

		if wt < weight {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return mid
}

// push returns best result from table, deletes it
// and inserts new result without changing order
func (s *smart) push(ent entry) (out entry, free bool) {
	if ok := s.table.TryLock(); ok {
		defer s.table.Unlock()
	} else {
		return ent, false
	}

	i := s.close(ent.weight)
	if i < 0 {
		return ent, false
	}

	out = s.table.entrys[0]
	weight := s.table.entrys[i].weight

	switch {
	case i == 0:
		s.table.entrys[i] = ent
	case weight < ent.weight:
		buf := make([]entry, len(s.table.entrys))
		copy(buf, s.table.entrys[1:])

		buf[i] = ent
		i++

		copy(buf[i:], s.table.entrys[i:])
		s.table.entrys = buf
	case ent.weight < weight:
		buf := make([]entry, len(s.table.entrys))
		copy(buf, s.table.entrys[1:])

		buf[i-1] = ent

		copy(buf[i:], s.table.entrys[i:])
		s.table.entrys = buf
	case weight == ent.weight:
		buf := make([]entry, len(s.table.entrys))
		copy(buf, s.table.entrys[1:])

		buf[i] = ent
		for i++; i < len(buf); i++ {
			e := s.table.entrys[i-1]
			e.weight++

			buf[i] = e
		}

		s.table.entrys = buf
	}

	return out, true
}

// add adds new backend to results table
func (s *smart) add(addr string) {
	defer s.table.Unlock()
	s.table.Lock()

	ent := entry{addr, 0}
	buf := make([]entry, len(s.table.entrys)+1)
	buf[0] = ent
	for i, e := range s.table.entrys {
		e.weight++
		buf[i+1] = e
	}

	s.table.entrys = buf
}

// cut shortens the results table
//
// Note: No need to worry about deleting every record of same backend,
// all records are deleted automatically.
// Therefore, here we have simple reduction method
func (s *smart) cut() {
	defer s.table.Unlock()
	s.table.Lock()

	if len(s.table.entrys) > 0 {
		s.table.entrys = s.table.entrys[:len(s.table.entrys)-1]
	}
}
