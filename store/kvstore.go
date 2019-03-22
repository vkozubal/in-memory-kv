package store

import (
	"fmt"
	"github.com/stefankopieczek/gossip/log"
	"sync"
	"sync/atomic"
)

func NewRegistry(size uint64) (Registry, error) {
	return &registry{maxSize: size, currentSize: 0}, nil
}

type Registry interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) bool
}

type registry struct {
	reg         sync.Map
	maxSize     uint64
	currentSize uint64
}

func (r *registry) Set(key string, value string) error {
	if count := r.getCurrentSize(); count > r.maxSize-1 {
		log.Debug("Reached maximum store capacity `%d`.", r.maxSize)
		return fmt.Errorf("Reached store max size %d.", r.maxSize)
	}

	if err := checkKey(key); err != nil {
		return err
	}
	if len(value) > 512 {
		return fmt.Errorf("value too long")
	}
	r.reg.Store(key, value)
	r.increaseSize()
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
	r.decreaseSize()
	return nil
}

func (r *registry) decreaseSize() {
	atomic.AddUint64(&r.currentSize, ^uint64(1-1))
}

func (r *registry) increaseSize() uint64 {
	return atomic.AddUint64(&r.currentSize, 1)
}

func (r *registry) getCurrentSize() uint64 {
	return atomic.LoadUint64(&r.currentSize)
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
