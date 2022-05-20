package main

import (
	"fmt"

	"github.com/iivvoo/warehouse/warehouse"
)

func main() {
	cache := warehouse.New[string, string]()
	cache.Set("foo", "hello")
	fmt.Println(cache.Get("foo"))
}
