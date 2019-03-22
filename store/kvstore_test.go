package store

import "testing"

func TestCanSetGetKey(t *testing.T) {
	r := Registry{}
	err := r.Set("key", "v")
	if err != nil {
		t.Errorf("Should be able to Set key.")
	}
	val, err := r.Get("key")
	if err != nil {
		t.Errorf("Should be able to Get key.")
	}
	if val != "v" {
		t.Errorf("Should fetch the same value")
	}
}

func TestSetFailOnEmptyKey(t *testing.T) {
	r := Registry{}
	err := r.Set("", "valid value")
	if err == nil || err.Error() != "empty value not allowed" {
		t.Errorf("Should fail on empty key")
	}
}
