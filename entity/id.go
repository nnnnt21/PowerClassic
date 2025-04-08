package entity

import (
	"fmt"
	"sync"
)

type EntityIdManager struct {
	mu     sync.Mutex
	nextId byte
	used   map[byte]bool
}

func NewEntityIdManager() *EntityIdManager {
	return &EntityIdManager{
		nextId: 0,
		used:   make(map[byte]bool),
	}
}

func (m *EntityIdManager) NextEntityId() byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	start := m.nextId
	for {
		candidate := m.nextId
		if candidate == 255 {
			m.nextId = 0
			candidate = m.nextId
		}
		if !m.used[candidate] {
			m.used[candidate] = true
			m.incrementNextId()
			return candidate
		}
		m.incrementNextId()
		if m.nextId == start {
			panic(fmt.Errorf("no available entity id"))
		}
	}
}

func (m *EntityIdManager) incrementNextId() {
	m.nextId++
	if m.nextId == 255 {
		m.nextId = 0
	}
}

func (m *EntityIdManager) UnregisterEntityId(id byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.used, id)
}
