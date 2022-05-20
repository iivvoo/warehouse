package warehouse

import (
	"time"

	"github.com/iivvoo/warehouse/genx"
)

/*


   c := warehouse[string].New()

   // Should key also be generic? We can do that and optionally define a StringKeyCache

*/

type entry[T any] struct {
	exp   time.Time
	value T
}

func (e *entry[T]) Expired() bool {
	return !e.exp.IsZero() || e.exp.Before(time.Now())
}

type warehouse[K comparable, T any] struct {
	cache      map[K]*entry[T]
	expiration time.Duration
}

func New[K comparable, T any]() *warehouse[K, T] {
	return &warehouse[K, T]{
		cache:      make(map[K]*entry[T]),
		expiration: 0,
	}
}

func (w *warehouse[K, T]) SetWithExpiration(k K, v T, expiration time.Duration) {
	expires := time.Time{} // never
	if expiration != 0 {
		expires = time.Now().Add(expiration)
	}
	w.cache[k] = &entry[T]{exp: expires, value: v}
}

func (w *warehouse[K, T]) Set(k K, v T) {
	w.SetWithExpiration(k, v, w.expiration)
}

func (w *warehouse[K, T]) Get(k K) T {
	e := w.cache[k]

	if e == nil {
		return genx.Zero[T]()
	}

	// check if it expired
	if e.Expired() {
		return e.value
	}
	return genx.Zero[T]()
}

func (w *warehouse[K, T]) GetSetWithExpiration(k K, callable func(k K) T, expiration time.Duration) T {
	if e := w.cache[k]; e != nil && !e.Expired() {
		return e.value
	}

	v := callable(k)

	w.SetWithExpiration(k, v, expiration)
	return v
}

func (w *warehouse[K, T]) GetSet(k K, callable func(k K) T) T {
	return w.GetSetWithExpiration(k, callable, w.expiration)
}

func (w *warehouse[K, T]) Cleanup() {
	// This touches all entries during each run which is not very efficient for huge caches.
	for k, v := range w.cache {
		if v.Expired() {
			delete(w.cache, k)
		}
	}
}
