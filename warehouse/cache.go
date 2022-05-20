package warehouse

import (
	"sync"
	"time"
)

type entry[T any] struct {
	exp   time.Time
	value T
}

func (e *entry[T]) Expired() bool {
	return !(e.exp.IsZero() || e.exp.After(time.Now()))
}

type warehouse[K comparable, T any] struct {
	cache      map[K]*entry[T]
	expiration time.Duration
	ticker     *time.Ticker
	mut        sync.RWMutex
}

const DefaultTimeout = 0 // never expires

func New[K comparable, T any]() *warehouse[K, T] {
	return NewWithExpiration[K, T](DefaultTimeout)
}

func NewWithExpiration[K comparable, T any](exp time.Duration) *warehouse[K, T] {
	w := &warehouse[K, T]{
		cache:      make(map[K]*entry[T]),
		expiration: exp,
		ticker:     time.NewTicker(time.Minute),
	}

	go w.Loop()

	return w
}

func (w *warehouse[K, T]) SetWithExpiration(k K, v T, expiration time.Duration) {
	w.mut.Lock()
	defer w.mut.Unlock()

	expires := time.Time{} // never
	if expiration != 0 {
		expires = time.Now().Add(expiration)
	}
	w.cache[k] = &entry[T]{exp: expires, value: v}
}

func (w *warehouse[K, T]) Set(k K, v T) {
	w.SetWithExpiration(k, v, w.expiration)
}

func (w *warehouse[K, T]) GetFound(k K) (T, bool) {
	var zero T

	w.mut.RLock()
	w.mut.RUnlock()

	e := w.cache[k]

	if e == nil || e.Expired() {
		return zero, false
	}

	return e.value, true
}

func (w *warehouse[K, T]) Get(k K) T {
	v, _ := w.GetFound(k)
	return v
}

// implement HasKey
func (w *warehouse[K, T]) HasKey(k K) bool {
	_, found := w.GetFound(k)
	return found
}

func (w *warehouse[K, T]) GetSetWithExpiration(k K, callable func(k K) T, expiration time.Duration) T {
	if v, found := w.GetFound(k); found {
		return v
	}

	v := callable(k)

	w.SetWithExpiration(k, v, expiration)
	return v
}

func (w *warehouse[K, T]) GetSet(k K, callable func(k K) T) T {
	return w.GetSetWithExpiration(k, callable, w.expiration)
}

func (w *warehouse[K, T]) Cleanup() {
	w.mut.Lock()
	defer w.mut.Unlock()
	// This touches all entries during each run which is not very efficient for huge caches.
	for k, v := range w.cache {
		if v.Expired() {
			delete(w.cache, k)
		}
	}
}

func (w *warehouse[K, T]) Loop() {
	for range w.ticker.C {
		w.Cleanup()
	}
}

func (w *warehouse[K, T]) Stop() {
	// terminates Cleanup Loop
	w.ticker.Stop()
}
