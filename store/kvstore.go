package store

import (
	"fmt"
	"sync"
)

type Registry struct {
	reg sync.Map
}

func (r *Registry) Set(key string, value string) error {
	err := checkKey(key)
	if err != nil {
		return err
	}
	if len(value) > 512 {
		return fmt.Errorf("value too long")
	}
	r.reg.Store(key, value)
	return nil
}

func (r *Registry) Get(key string) (string, error) {
	if val, found := r.reg.Load(key); found {
		val := val.(string)
		return val, nil
	}
	return "", fmt.Errorf("not found")
}

func (r *Registry) Delete(key string) error {
	if !r.Exists(key) {
		return fmt.Errorf("not found")
	}
	r.reg.Delete(key)
	return nil
}

func (r *Registry) Exists(key string) bool {
	_, found := r.reg.Load(key)
	return found
}

func checkKey(key string) error {
	if len(key) == 0 {
		return fmt.Errorf("empty value not allowed")
	}

	if len(key) > 16 {
		return fmt.Errorf("key too long")
	}
	return nil
}
