package balancer

import (
	"testing"
)

func TestClose(t *testing.T) {
	bal := NewBalancer([]string{"one", "two", "three"})

	i := bal.smart.close(3)
	if i != 2 {
		t.Errorf("expected index '%d' instead of index '%d'", 2, i)
	}

	bal = NewBalancer([]string{"one", "two", "three"})
	bal.smart.table[2].weight++

	i = bal.smart.close(3)
	if i != 2 {
		t.Errorf("expected index '%d' instead of index '%d'", 2, i)
	}
}

func TestPush(t *testing.T) {
	bal := NewBalancer([]string{"one"})
	ent := entry{"one", 0}
	bal.smart.push(0, ent)

	if bal.smart.table[0] != ent {
		t.Errorf("expected a null element with parameters '%s' and '%d'", ent.address, ent.weight)
	}

	bal = NewBalancer([]string{"one", "two"})
	ent = entry{"one", 3}
	bal.smart.push(1, ent)

	if bal.smart.table[1] != ent {
		t.Errorf("expected a last element with parameters '%s' and '%d'", ent.address, ent.weight)
	}

	bal = NewBalancer([]string{"one", "two"})
	bal.smart.table[1].weight++
	ent = entry{"one", 2}
	bal.smart.push(1, ent)

	if bal.smart.table[0] != ent {
		t.Errorf("expected a last element with parameters '%s' and '%d'", ent.address, ent.weight)
	}

	bal = NewBalancer([]string{"one", "two"})
	bal.smart.table[1].weight++
	ent = entry{"one", 2}
	bal.smart.push(1, ent)

	if bal.smart.table[0] != ent {
		t.Errorf("expected a last element with parameters '%s' and '%d'", ent.address, ent.weight)
	}
}

func TestEnter(t *testing.T) {
	bal := NewBalancer([]string{"one", "two"})
	bal.smart.table[1].weight++

	bal.smart.enter("three")

	if bal.smart.table[0].address != "three" {
		t.Errorf("expected a null element with address '%s'", "three")
	}

	if bal.smart.table[1].weight-bal.smart.table[0].weight == 1 {
		t.Error("the table should remain in ascending order")
	}
}

func TestCut(t *testing.T) {
	bal := NewBalancer([]string{"one"})
	l := len(bal.smart.table)
	bal.smart.cut()
	if l == len(bal.smart.table) {
		t.Error("table reduction was expected")
	}
}
