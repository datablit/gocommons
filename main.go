package main

import (
	"fmt"
	"time"

	"github.com/datablit/gocommons/syncache"
)

func main() {
	cache := syncache.New[string](2 * time.Minute)

	// 1. GetOrLoad (auto-caches result)
	val, _ := cache.GetOrLoad("greeting", func() (string, error) {
		fmt.Println("Loading greeting...")
		return "Hello, world!", nil
	})
	fmt.Println(val) // "Hello, world!"

	// 2. Set manually
	cache.Set("message", "Hi again!")
	msg, _ := cache.GetOrLoad("message", func() (string, error) {
		return "This won't be called", nil
	})
	fmt.Println(msg) // "Hi again!"

	// 3. Clear example
	cache.Delete("message")
}
