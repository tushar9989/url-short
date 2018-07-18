package cache

import (
	"sync"
	"time"
)

type item struct {
	value      map[string]string
	lastAccess int64
}

type TTLMap struct {
	m map[string]*item
	l sync.RWMutex
}

func New(maxTTL int) (m *TTLMap) {
	m = &TTLMap{m: make(map[string]*item, 0)}
	go func() {
		for now := range time.Tick(time.Minute) {
			m.l.Lock()
			for k, v := range m.m {
				if now.Unix()-v.lastAccess > int64(maxTTL) {
					delete(m.m, k)
				}
			}
			m.l.Unlock()
		}
	}()
	return
}

func (m *TTLMap) Len() int {
	return len(m.m)
}

func (m *TTLMap) Put(k string, v map[string]string) {
	m.l.Lock()
	it, ok := m.m[k]
	if !ok {
		it = &item{value: v}
		m.m[k] = it
	}
	it.lastAccess = time.Now().Unix()
	m.l.Unlock()
}

func (m *TTLMap) Get(k string) (v map[string]string, found bool) {
	m.l.RLock()
	if it, ok := m.m[k]; ok {
		v = it.value
		found = true
	}
	m.l.RUnlock()
	return
}
