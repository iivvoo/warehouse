package warehouse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWarehouse(t *testing.T) {
	assert.True(t, false)

	cache := New[string, any]()
	cache.Set("x", 1)
}
