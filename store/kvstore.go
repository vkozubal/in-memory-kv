package store

import (
	"fmt"
	"sync"
)

type Registry struct {
	reg sync.Map
}

func (r *Registry) Set(key string, value string) error {
	if len(value) > 1024 {
		return fmt.Errorf("key too long")
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
	if r.Exists(key) {

	}
	r.reg.Delete(key)
	return nil
}

func (r *Registry) Exists(key string) bool {
	_, found := r.reg.Load(key)
	return found
}
