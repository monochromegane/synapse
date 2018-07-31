package synapse

import (
	"testing"
)

func TestProfileMerge(t *testing.T) {
	p1 := Profile{}
	p1["k1"] = "v1"
	p1["k2"] = "v2"

	p2 := Profile{}
	p2["k2"] = "v3"
	p2["k3"] = "v4"
	p1.Merge(p2)

	if len(p1) != 3 {
		t.Errorf("Profile should merge p1 and p2.")
	}

	for k, v := range p1 {
		if k == "k1" && v != "v1" {
			t.Errorf("Profile should merge %s:%s.", k, v)
		}
		if k == "k2" && v != "v3" {
			t.Errorf("Profile should merge %s:%s.", k, v)
		}
		if k == "k3" && v != "v4" {
			t.Errorf("Profile should merge %s:%s.", k, v)
		}
	}
}
