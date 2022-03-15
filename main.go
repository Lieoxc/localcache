package main

import (
	"fmt"
	"goCache/cache"
)

func main() {
	fmt.Println("hello world")
	lru := cache.NewLRU(10)
	lru.Add("1", "qqqq")
	if val, ok := lru.Get("1"); ok {
		fmt.Println(val)
	}

}
