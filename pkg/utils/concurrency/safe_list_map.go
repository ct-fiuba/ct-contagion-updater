package concurrency

import (
	"sync"
)

type SafeStringListMap struct {
	mu      sync.Mutex
	listMap map[string][]string
	count   int
}

func NewSafeStringListMap() *SafeStringListMap {
	lm := &SafeStringListMap{
		listMap: make(map[string][]string),
		count:   0,
	}
	return lm
}

func (lm *SafeStringListMap) Add(key, newValue string) {
	lm.mu.Lock()
	v, _ := lm.listMap[key]
	lm.listMap[key] = append(v, newValue)
	lm.count++
	lm.mu.Unlock()
}

func (lm *SafeStringListMap) Clear() map[string][]string {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	listMap := lm.listMap
	lm.listMap = make(map[string][]string)
	lm.count = 0
	return listMap
}

func (lm *SafeStringListMap) Get() map[string][]string {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	return lm.listMap
}

func (lm *SafeStringListMap) Count() int {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	return lm.count
}
