package concurrency

import (
	"sync"
)

type SafeStringListMap struct {
	mu      sync.Mutex
	listMap map[string][]string
}

func NewSafeStringListMap() *SafeStringListMap {
	lm := &SafeStringListMap{
		listMap: make(map[string][]string),
	}
	return lm
}

func (lm *SafeStringListMap) Add(key, newValue string) {
	lm.mu.Lock()
	v, _ := lm.listMap[key]
	lm.listMap[key] = append(v, newValue)
	lm.mu.Unlock()
}

func (lm *SafeStringListMap) Clear() map[string][]string {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	listMap := lm.listMap
	lm.listMap = make(map[string][]string)
	return listMap
}

func (lm *SafeStringListMap) Get() map[string][]string {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	return lm.listMap
}
