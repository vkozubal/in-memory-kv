package store

import (
	"testing"
)

func TestCanSetGetKeyPositive(t *testing.T) {
	r := registry{}
	err := r.Set("key", "v")
	if err != nil {
		t.Errorf("Should be able to Set key.")
	}
	val, err := r.Get("key")
	if err != nil {
		t.Errorf("Should be able to Get key.")
	}
	if val != "v" {
		t.Errorf("Should fetch the same value.")
	}
}

func TestSetFailOnEmptyKey(t *testing.T) {
	r := registry{}
	err := r.Set("", "valid value")
	if err == nil || err.Error() != "empty value not allowed" {
		t.Errorf("Should fail on empty key.")
	}
}

func TestSetKeyIsTooBig(t *testing.T) {
	r := registry{}
	err := r.Set("exceeding length k", "valid value")
	if err == nil || err.Error() != "key too long" {
		t.Errorf("Should fail on long key.")
	}
}

func TestSetValueIsTooLong(t *testing.T) {
	r := registry{}
	tooLongValue := "p2xqsqTpbXKCKT7i4l7ZiLRYIBg3rkD15BRJHXNqqisW0bZIjfoJeri3qMeb04wTr40t0qSQfMUNFBkkbpy1eCOXYHw4gwBzBp0kSe46Sof3ATHLdmGH4OUNiM4GJEUkNisCOgFEeRUepZSLfHc2hHpTwkjgeR0SKBGwfLQlZekZxjISvVgC3btBUnq8IueZV5XPhyJBmKV2Euh0VbIiNyrqe5tw32FbGegOL4SXNIfOeU0q7Vmdor2I8t8e5nViSo2hGqwa8MJwzISLq7weFk1BFC0bCUCkjki7TWXikoxvd4PZ8rjpG1fqPa5UKkJWOgcj6kNdJdKd8IL3fiE1xp9sEEUTJiKCOOeJb0DgpVu7rWWeN2LtjxlV2f7m5isbqbOJdyfZOqY4Maq70M4ODIKrK3r4WbQZvq8JlUDYbbwhqNMb7VDZVILp7xJbTHDU1P6fVkZRMWjPIUAyOoai0Y1f89XB6aecbnWTJkIv4myY5hokpghPcit5ER9EcPjSm"
	err := r.Set("max len key pam ", tooLongValue)
	if err == nil || err.Error() != "value too long" {
		t.Errorf("Should fail on long key")
	}
}

func TestDeleteNotFound(t *testing.T) {
	r := registry{}
	err := r.Delete("KEY")
	if err == nil || err.Error() != "not found" {
		t.Errorf("Should fail on non existent key.")
	}
}

func TestExists(t *testing.T) {
	r := registry{}
	if found := r.Exists("test_key"); found {
		t.Errorf("Test key shouldn't exist.")
	}

	if err := r.Set("test_key", "1"); err != nil {
		t.Errorf("Can't set key.")
	}

	if found := r.Exists("test_key"); !found {
		t.Errorf("Test key should exist.")
	}
}

func TestDeletePositive(t *testing.T) {
	r := registry{}
	if err := r.Set("KEY", "v"); err != nil {
		t.Fail()
	}

	if found := r.Exists("KEY"); !found {
		t.Errorf("Key should exist.")
	}

	if err := r.Delete("KEY"); err != nil {
		t.Errorf("Should fail on non existent key.")
	}

	if found := r.Exists("KEY"); found {
		t.Errorf("Key should not exist.")
	}
}

func TestStoreSizeIsMaintained(t *testing.T) {
	r := registry{maxSize: 1}
	if err := r.Set("1", "v1"); err != nil {
		t.Errorf("Should store. Capasity is not reached.")
	}

	if err := r.Set("2", "v2"); err == nil {
		t.Errorf("Capasity should be reached.")
	}

	// delete key
	if err := r.Delete("1"); err != nil {
		t.Errorf("Should be deleted.")
	}

	// insert again
	if err := r.Set("1", "v1"); err != nil {
		t.Errorf("Should store. Capasity is not reached.")
	}
}
