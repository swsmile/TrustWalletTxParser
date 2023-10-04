package mem_store

import (
	"awesomeProject/entity"
	"awesomeProject/parser"
	"sync"
)

// inMemoryStore implements a simple in-memory store
type inMemoryStore struct {
	lock sync.RWMutex
	m    map[string][]entity.Transaction
}

func NewInMemoryStore() parser.Store {
	return &inMemoryStore{
		lock: sync.RWMutex{},
		m:    make(map[string][]entity.Transaction),
	}
}

func (s *inMemoryStore) Insert(address string, transaction entity.Transaction) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.m[address] = append(s.m[address], transaction)
}

func (s *inMemoryStore) Get(address string) []entity.Transaction {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.m[address]
}
