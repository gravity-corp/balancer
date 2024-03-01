package balancer

import "testing"

func TestAdd(t *testing.T) {
	bal := NewRoundBalancer([]string{"one"})
	bal.round.add("two")

	if len(bal.round.addresses) != 2 {
		t.Error("the list of addresses was expected to expand")
	}
}

func TestRemove(t *testing.T) {
	bal := NewRoundBalancer([]string{"one"})
	bal.round.remove("one")

	if len(bal.round.addresses) != 0 {
		t.Error("expected reduction of the list of addresses")
	}
}
