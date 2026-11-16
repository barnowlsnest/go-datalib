package lru

import (
	"errors"
	"sync"

	"github.com/barnowlsnest/go-datalib/v5/pkg/list"
	"github.com/barnowlsnest/go-datalib/v5/pkg/node"
)

// defaultLRUCapacity is used when a non-positive capacity is supplied to New.
const defaultLRUCapacity = 1024

// ErrCacheMiss is returned by Get when the requested key is not cached.
var ErrCacheMiss = errors.New("lru: cache miss")

type (
	// entry pairs a cached value with the list node tracking its recency.
	entry[T any] struct {
		node  *node.Node
		value *T
	}

	// LRU is a generic, thread-safe least-recently-used cache keyed by uint64.
	LRU[T any] struct {
		capacity int
		cache    map[uint64]*entry[T]
		usage    *list.LinkedList
		mux      sync.Mutex
	}
)

// New creates an LRU cache with the given capacity.
//
// A capacity less than or equal to zero falls back to defaultLRUCapacity.
func New[T any](capacity int) *LRU[T] {
	if capacity <= 0 {
		capacity = defaultLRUCapacity
	}

	return &LRU[T]{
		capacity: capacity,
		cache:    make(map[uint64]*entry[T]),
		usage:    list.New(),
	}
}

// Put inserts or updates the value for key, marking it most recently used.
//
// When inserting a new key into a full cache, the least recently used entry is
// evicted first.
func (lru *LRU[T]) Put(key uint64, value *T) {
	lru.mux.Lock()
	defer lru.mux.Unlock()

	if existing, ok := lru.cache[key]; ok {
		existing.value = value
		lru.usage.MoveToHead(existing.node)

		return
	}

	if len(lru.cache) >= lru.capacity {
		lru.evict()
	}

	n := node.New(key, nil, nil)
	lru.usage.Unshift(n)
	lru.cache[key] = &entry[T]{node: n, value: value}
}

// Get returns the value for key and marks it most recently used.
//
// It returns ErrCacheMiss if the key is not present.
func (lru *LRU[T]) Get(key uint64) (*T, error) {
	lru.mux.Lock()
	defer lru.mux.Unlock()

	cached, ok := lru.cache[key]
	if !ok {
		return nil, ErrCacheMiss
	}

	lru.usage.MoveToHead(cached.node)

	return cached.value, nil
}

// GetRange appends the values for the consecutive keys from through to, in key
// order, to dst and returns the extended slice. Each returned key is marked most
// recently used, as Get would mark it.
//
// It stops at the first key that is not cached and returns the values found
// before it together with ErrCacheMiss, so that key is from+len(added). A range
// with from above to appends nothing.
//
// The whole range is read under one lock acquisition. Calling Get once per key
// takes the lock once per key, which made the lock the bottleneck when many
// readers fetched the same recent range at once.
func (lru *LRU[T]) GetRange(from, to uint64, dst []*T) ([]*T, error) {
	if from > to {
		return dst, nil
	}

	lru.mux.Lock()
	defer lru.mux.Unlock()

	for key := from; ; key++ {
		cached, ok := lru.cache[key]
		if !ok {
			return dst, ErrCacheMiss
		}

		lru.usage.MoveToHead(cached.node)
		dst = append(dst, cached.value)

		// Checked here rather than in the loop condition so a range ending at
		// the largest key stops instead of wrapping around to zero.
		if key == to {
			return dst, nil
		}
	}
}

// Len returns the number of entries currently cached.
func (lru *LRU[T]) Len() int {
	lru.mux.Lock()
	defer lru.mux.Unlock()

	return len(lru.cache)
}

// Delete removes the entry for key from the cache, if present.
func (lru *LRU[T]) Delete(keys ...uint64) {
	lru.mux.Lock()
	defer lru.mux.Unlock()

	for _, key := range keys {
		cached, ok := lru.cache[key]
		if !ok {
			continue
		}

		lru.usage.MoveToTail(cached.node)
		lru.evict()
	}
}

// Clear removes all entries from the cache, leaving it empty while preserving
// its capacity.
func (lru *LRU[T]) Clear() {
	lru.mux.Lock()
	defer lru.mux.Unlock()

	lru.cache = make(map[uint64]*entry[T])
	lru.usage = list.New()
}

// evict removes the least recently used entry (the list tail).
//
// Callers must hold lru.mux.
func (lru *LRU[T]) evict() {
	evicted := lru.usage.Pop()
	if evicted == nil {
		return
	}

	delete(lru.cache, evicted.ID())
}
