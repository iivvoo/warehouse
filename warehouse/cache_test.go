package warehouse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWarehouse(t *testing.T) {
	t.Run("Test get/set with hit", func(t *testing.T) {
		cache := New[string, string]()

		defer cache.Stop()

		cache.Set("hello", "world")

		assert.Equal(t, "world", cache.Get("hello"))
	})
	t.Run("Test get/set with miss", func(t *testing.T) {
		cache := New[string, string]()

		defer cache.Stop()

		assert.Equal(t, "", cache.Get("hello")) // "" is the Zero value for strings
	})
}
