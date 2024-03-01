package balancer

import (
	"testing"
)

func TestRoundNext(t *testing.T) {
	bal := NewRoundBalancer([]string{"one"})
	_, err := bal.Next()
	if err != nil {
		t.Error(err)
	}

	bal = NewRoundBalancer([]string{})
	_, err = bal.Next()
	if err == nil {
		t.Error("next should return an error if no address is being tracked")
	}
}

func TestRoundEnable(t *testing.T) {
	bal := NewRoundBalancer([]string{"one"})
	bal.Enable("two")

	if bal.round.addresses[len(bal.round.addresses)-1] != "two" {
		t.Errorf("expected a last element with address '%s'", "two")
	}
}

func TestRoundDisable(t *testing.T) {
	bal := NewRoundBalancer([]string{"one"})
	bal.Disable("one")

	if len(bal.round.addresses) != 0 {
		t.Errorf("the expected length of the list is '%d'", 0)
	}
}
