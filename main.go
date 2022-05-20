package main

import (
	"fmt"

	"github.com/iivvoo/warehouse/warehouse"
)

func main() {
	fmt.Println("Hello World")

	cache := warehouse.New[string, string]()
	cache.Set("foo", "hello")
	fmt.Println(cache, cache.Get("foo"))
	cache.Cleanup()
}
