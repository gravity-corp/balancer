package balancer

type smart struct {
	from    chan entry
	to      chan entry
	enable  chan string
	disable chan struct{}
	table   []entry
}

type entry struct {
	address string
	weight  int64
}

// close finds a similar backend weight or the one closest to it
// and returns its index
func (s *smart) close(weight int64) (inx int) {
	first := 0
	last := len(s.table) - 1

	for first <= last {
		mid := first + (last-first)/2

		if s.table[mid].weight == weight {
			return mid
		}

		if first == mid {
			return mid
		}

		if s.table[mid].weight < weight {
			first = mid + 1
		} else {
			last = mid - 1
		}
	}

	return -1
}

// push returns the best result from the table, deletes it
// and inserts a new result without changing the order
func (s *smart) push(i int, ent entry) (out entry) {
	out = s.table[0]

	switch {
	case i == 0:
		s.table[i] = ent
	case s.table[i].weight < ent.weight:
		buf := make([]entry, len(s.table))
		copy(buf, s.table[1:])

		buf[i] = ent
		i++

		copy(buf[i:], s.table[i:])
		s.table = buf

	case ent.weight < s.table[i].weight:
		buf := make([]entry, len(s.table))
		copy(buf, s.table[1:])

		buf[i-1] = ent

		copy(buf[i:], s.table[i:])
		s.table = buf
	case s.table[i].weight == ent.weight:
		buf := make([]entry, len(s.table))
		copy(buf, s.table[1:])

		buf[i] = ent
		for i++; i < len(buf); i++ {
			e := s.table[i-1]
			e.weight++

			buf[i] = e
		}

		s.table = buf
	}

	return out
}

// add adds a new backend to the results table
func (s *smart) enter(addr string) {
	ent := entry{addr, 0}
	buf := make([]entry, len(s.table)+1)
	buf[0] = ent
	for i, e := range s.table {
		e.weight++
		buf[i+1] = e
	}

	s.table = buf
}

// cut shortens the results table
//
// Note: No need to worry about deleting every record of the same backend,
// all records are deleted automatically.
// Therefore, here we have a simple reduction method
func (s *smart) cut() {
	if len(s.table) > 0 {
		s.table = s.table[:len(s.table)-1]
	}
}
