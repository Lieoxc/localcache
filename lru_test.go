package cache

import (
	"testing"
)

func TestGet(t *testing.T) {
	lru := NewLRU(10)
	str := "qqqq"
	lru.Add("1", str)
	if val, ok := lru.Get("1"); ok {
		if val.(string) != str {
			t.Errorf("Expect %s but got %v", str, val)
		}
	}
}
func TestGetALL(t *testing.T) {
	lru := NewLRU(3)
	lru.Add("1", "q")
	lru.Add("2", "w")
	lru.Add("3", "e")
	lru.Add("4", "r")

	data := lru.GetALL()
	for _, v := range data {
		if v.value.(string) == "q" {
			t.Errorf("your lru cannot delete Least  recently  use data")
		}
	}
}
