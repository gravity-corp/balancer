package balancer

import (
	"testing"
	"time"
)

func TestNext(t *testing.T) {
	bal := NewBalancer([]string{"one"})
	_, err := bal.Next()
	if err != nil {
		t.Error(err)
	}

	bal = NewBalancer([]string{})
	_, err = bal.Next()
	if err == nil {
		t.Error("next should return an error if more than one address is not being tracked")
	}
}

func TestScore(t *testing.T) {
	bal := NewBalancer([]string{"one"})

	bal.Score("two", 2)
	time.Sleep(time.Second / 60)

	if bal.smart.table[0].address != "two" {
		t.Errorf("the embedding of the record in the table was expected")
	}
}

func TestEnable(t *testing.T) {
	bal := NewBalancer([]string{"one"})

	bal.Enable("two")
	time.Sleep(time.Second / 60)

	if bal.smart.table[0].address != "two" {
		t.Errorf("expected a null element with address '%s'", "two")
	}

	if bal.round.addresses[len(bal.round.addresses)-1] != "two" {
		t.Errorf("expected a last element with address '%s'", "two")
	}
}

func TestDisable(t *testing.T) {
	bal := NewBalancer([]string{"one"})

	bal.Disable("one")
	time.Sleep(time.Second / 60)

	if len(bal.round.addresses) != 0 {
		t.Errorf("the expected length of the list is '%d'", 0)
	}

	if len(bal.smart.table) != 0 {
		t.Errorf("the expected length of the table is '%d'", 0)
	}
}
