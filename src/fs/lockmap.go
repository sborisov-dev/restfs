package fs

import "sync"

type lockmap struct {
	items map[string]*sync.RWMutex
	sync.Mutex
}

func newLockmap() *lockmap {
	return &lockmap{
		items: make(map[string]*sync.RWMutex),
	}
}

func (lm *lockmap) getMutex(name string) *sync.RWMutex {
	m, ok := lm.items[name]
	if !ok {
		m = &sync.RWMutex{}
		lm.items[name] = m
	}

	return m
}
