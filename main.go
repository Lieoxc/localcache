package main

import (
	"fmt"
	"goCache/cache"
)

func main() {
	fmt.Println("hello world")
	// lru := cache.NewLRU(10)
	// lru.Add("1", "qqqq")
	// if val, ok := lru.Get("1"); ok {
	// 	fmt.Println(val)
	// }
	lfu := cache.NewLFU(10)
	lfu.Add("1", "qqqq")
	lfu.Add("2", "AAA")
	if val, ok := lfu.Get("2"); ok {
		fmt.Println(val)
	}
}
