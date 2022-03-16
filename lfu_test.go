package localcache

import (
	"testing"
)

func TestLfuGet(t *testing.T) {
	lru := NewLFU(2)
	lru.Add("1", "1")
	lru.Add("2", "2")
	lru.Add("3", "3")
	if _, ok := lru.Get("1"); ok {
		t.Errorf("Expect Ok is false")
	}
	if val, ok := lru.Get("2"); ok {
		if val.(string) != "2" {
			t.Errorf("wan GET 2")
		}
	}
	if val, ok := lru.Get("3"); ok {
		if val.(string) != "3" {
			t.Errorf("wan GET 3")
		}
	}
	lru.Add("4", "4")
	if _, ok := lru.Get("2"); ok {
		t.Errorf("Expect Ok is false")
	}
	if val, ok := lru.Get("4"); ok {
		if val.(string) != "4" {
			t.Errorf("wan GET 4")
		}
	}
}
