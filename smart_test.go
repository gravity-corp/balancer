package balancer

import (
	"testing"
)

func TestClose(t *testing.T) {
	bal := NewSmartBalancer([]string{"one", "two", "three"})
	i := bal.smart.close(3)
	if i != 2 {
		t.Errorf("expected index '%d' instead of index '%d'", 2, i)
	}

	bal = NewSmartBalancer([]string{"one", "two", "three"})
	bal.smart.table.entrys[2].weight++

	i = bal.smart.close(3)
	if i != 2 {
		t.Errorf("expected index '%d' instead of index '%d'", 2, i)
	}
}

func TestPush(t *testing.T) {
	bal := NewSmartBalancer([]string{"one"})
	ent := entry{"one", 0}
	bal.smart.push(ent)

	if bal.smart.table.entrys[0] != ent {
		t.Errorf("expected a null element with parameters '%s' and '%d'", ent.address, ent.weight)
	}

	bal = NewSmartBalancer([]string{"one", "two"})
	ent = entry{"one", 3}
	bal.smart.push(ent)

	if bal.smart.table.entrys[1] != ent {
		t.Errorf("expected a last element with parameters '%s' and '%d'", ent.address, ent.weight)
	}

	bal = NewSmartBalancer([]string{"one", "two"})
	bal.smart.table.entrys[1].weight++
	ent = entry{"one", 2}
	bal.smart.push(ent)

	if bal.smart.table.entrys[0] != ent {
		t.Errorf("expected a last element with parameters '%s' and '%d'", ent.address, ent.weight)
	}

	bal = NewSmartBalancer([]string{"one", "two"})
	bal.smart.table.entrys[1].weight++
	ent = entry{"one", 2}
	bal.smart.push(ent)

	if bal.smart.table.entrys[0] != ent {
		t.Errorf("expected a last element with parameters '%s' and '%d'", ent.address, ent.weight)
	}
}

func TestSmartAdd(t *testing.T) {
	bal := NewSmartBalancer([]string{"one", "two"})
	bal.smart.table.entrys[1].weight++
	bal.smart.add("three")

	if bal.smart.table.entrys[0].address != "three" {
		t.Errorf("expected a null element with address '%s'", "three")
	}

	if bal.smart.table.entrys[1].weight-bal.smart.table.entrys[0].weight == 1 {
		t.Error("the table should remain in ascending order")
	}
}

func TestCut(t *testing.T) {
	bal := NewSmartBalancer([]string{"one"})
	l := len(bal.smart.table.entrys)
	bal.smart.cut()

	if l == len(bal.smart.table.entrys) {
		t.Error("table reduction was expected")
	}
}
