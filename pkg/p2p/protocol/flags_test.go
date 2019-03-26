package protocol

import "testing"

func TestMessageFlag_Has(t *testing.T) {
	var f MessageFlag = System

	if f.Has(System) == false {
		t.Fail()
	}
}
