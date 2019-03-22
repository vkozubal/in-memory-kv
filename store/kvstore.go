package store

import (
	"fmt"
	"sync"
)

func NewRegistry() (Registry, error) {
	return &registry{}, nil
}

type Registry interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) bool
}

type registry struct {
	reg sync.Map
}

func (r *registry) Set(key string, value string) error {
	if err := checkKey(key); err != nil {
		return err
	}
	if len(value) > 512 {
		return fmt.Errorf("value too long")
	}
	r.reg.Store(key, value)
	return nil
}

func (r *registry) Get(key string) (string, error) {
	if val, found := r.reg.Load(key); found {
		val := val.(string)
		return val, nil
	}
	return "", fmt.Errorf("not found")
}

func (r *registry) Delete(key string) error {
	if !r.Exists(key) {
		return fmt.Errorf("not found")
	}
	r.reg.Delete(key)
	return nil
}

func (r *registry) Exists(key string) bool {
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
