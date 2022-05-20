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
	if !e.exp.IsZero() || e.exp.After(time.Now()) {
		return genx.Zero[T]()
	}
	return e.value
}
