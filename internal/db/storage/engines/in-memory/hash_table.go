package inmemory

import "sync"

type Hasheable interface {
	Get(k string) (string, bool)
	Set(k string, v string)
	Del(k string)
}

type HashTable struct {
	data map[string]string
	mu   *sync.Mutex // todo bench for mutex or rwmutex
}

func (h *HashTable) Get(k string) (string, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	value, exists := h.data[k]

	return value, exists
}

func (h *HashTable) Set(k, v string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.data[k] = v
}

func (h *HashTable) Del(k string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.data, k)
}

func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]string),
		mu:   &sync.Mutex{},
	}
}
