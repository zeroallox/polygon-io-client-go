package pwsmodels

import (
    "sync"
)

// Pool is a generic wrapper around sync.Pool
type Pool[T any] struct {
    pool sync.Pool
}

// Release returns all items to the Pool.
func (this *Pool[T]) Release(item ...T) {
    for _, cItem := range item {
        this.pool.Put(cItem)
    }
}

// Acquire gets an item from the Pool.
func (this *Pool[T]) Acquire() T {
    return this.pool.Get().(T)
}
