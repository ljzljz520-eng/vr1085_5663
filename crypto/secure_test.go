package crypto

import "testing"

func TestWrapUnwrap(t *testing.T) {
	w, e := WrapSecret("k", "secret", "p")
	if e != nil {
		t.Fatal(e)
	}
	s, e := UnwrapSecret("k", w, "p")
	if e != nil || s != "secret" {
		t.Fatalf("%v %s", e, s)
	}
}
func TestBadCiphertext(t *testing.T) {
	if _, e := UnwrapSecret("k", "!", "p"); e == nil {
		t.Fatal("expected error")
	}
}
