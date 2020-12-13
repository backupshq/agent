package utils

import (
	"sync"
)

type SignalMap struct {
	syncMap sync.Map
}

func (m *SignalMap) Load(key string) (chan bool, bool) {
	val, ok := m.syncMap.Load(key)
	if ok {
		return val.(chan bool), true
	}
	return nil, false
}

func (m *SignalMap) Store(key string, value chan bool) {
	m.syncMap.Store(key, value)
}

func (m *SignalMap) Delete(key string) {
	m.syncMap.Delete(key)
}
