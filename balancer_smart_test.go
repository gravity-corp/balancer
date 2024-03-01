package balancer

import (
	"testing"
)

func TestSmartNext(t *testing.T) {
	bal := NewSmartBalancer([]string{"one"})
	_, err := bal.Next()
	if err != nil {
		t.Error(err)
	}

	bal = NewSmartBalancer([]string{})
	_, err = bal.Next()
	if err == nil {
		t.Error("next should return an error if no address is being tracked")
	}
}

func TestSmartScore(t *testing.T) {
	bal := NewSmartBalancer([]string{"one"})
	bal.Score("two", 2)

	if bal.smart.table.entrys[0].address != "two" {
		t.Errorf("the embedding of the record in the table was expected")
	}
}

func TestSmartEnable(t *testing.T) {
	bal := NewSmartBalancer([]string{"one"})
	bal.Enable("two")

	if bal.smart.table.entrys[0].address != "two" {
		t.Errorf("expected a null element with address '%s'", "two")
	}

	if bal.round.addresses[len(bal.round.addresses)-1] != "two" {
		t.Errorf("expected a last element with address '%s'", "two")
	}
}

func TestSmartDisable(t *testing.T) {
	bal := NewSmartBalancer([]string{"one"})
	bal.Disable("one")

	if len(bal.round.addresses) != 0 {
		t.Errorf("the expected length of the list is '%d'", 0)
	}

	if len(bal.smart.table.entrys) != 0 {
		t.Errorf("the expected length of the table is '%d'", 0)
	}
}
