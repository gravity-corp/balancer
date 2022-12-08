package balancer

import (
	"testing"
)

func TestNewBalancer(t *testing.T) {
	_, err := NewBalancer([]string{"one"})
	if err == nil {
		t.Error("at least 2 addresses were expected")
	}
}

func TestClose(t *testing.T) {
	bal, _ := NewBalancer([]string{"one", "two", "three", "four", "five"})

	i := bal.smart.close(3)
	if i != 3 {
		t.Errorf("index '%d' was not expected", i)
	}
}

func TestPush(t *testing.T) {
	bal, _ := NewBalancer([]string{"one", "two", "three"})
	ent := entry{"one", 0}
	bal.smart.push(0, ent)

	if bal.smart.table[0].weight != 0 {
		t.Error("the weight of the zero element was expected to be 0")
	}

	bal, _ = NewBalancer([]string{"one", "two", "three"})
	ent = entry{"one", 4}
	bal.smart.push(2, ent)

	if bal.smart.table[2].weight != 4 {
		t.Error("the weight of the last element was expected to be 4")
	}

	bal, _ = NewBalancer([]string{"one", "two", "three"})
	ent = entry{"one", 1}
	bal.smart.push(1, ent)

	if bal.smart.table[1].weight != 1 {
		t.Error("the weight of the first element was expected to be 4")
	}

	bal, _ = NewBalancer([]string{"one", "two", "three", "four"})
	bal.smart.table[2].weight++
	bal.smart.table[3].weight++
	ent = entry{"one", 2}
	bal.smart.push(1, ent)

	if bal.smart.table[1].weight != 2 {
		t.Error("the weight of the first element was expected to be 2")
	}

	bal, _ = NewBalancer([]string{"one", "two", "three", "four"})
	bal.smart.table[2].weight++
	bal.smart.table[3].weight++
	ent = entry{"one", 2}
	bal.smart.push(2, ent)

	if bal.smart.table[1].weight != 2 {
		t.Error("the weight of the first element was expected to be 2")
	}
}
