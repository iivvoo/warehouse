package warehouse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWarehouse(t *testing.T) {
	t.Run("Test get/set with hit", func(t *testing.T) {
		cache := New[string, string]()

		defer cache.Stop()

		cache.Set("hello", "world")

		assert.Equal(t, "world", cache.Get("hello"))
		assert.True(t, cache.HasKey("hello"))
	})
	t.Run("Test get/set with miss", func(t *testing.T) {
		cache := New[string, string]()

		defer cache.Stop()

		assert.Equal(t, "", cache.Get("hello")) // "" is the Zero value for strings
		assert.False(t, cache.HasKey("hello"))
	})
	t.Run("Test expiration", func(t *testing.T) {
		cache := New[string, string]()

		defer cache.Stop()

		cache.SetWithExpiration("hello", "world", -time.Second)

		assert.Equal(t, "", cache.Get("hello"))
	})
	t.Run("Test GetSet with hit", func(t *testing.T) {
		cache := New[string, string]()

		defer cache.Stop()

		cache.GetSet("hello", func(k string) string { return "world" })

		assert.Equal(t, "world", cache.Get("hello"))
	})
}
