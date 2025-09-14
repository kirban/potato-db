package inmemory

type Hasheable interface {
	Get(k string) (string, bool)
	Set(k string, v string)
	Del(k string)
}

type HashTable struct {
	data map[string]string
}

func (h *HashTable) Get(k string) (string, bool) {
	value, exists := h.data[k]

	return value, exists
}

func (h *HashTable) Set(k, v string) {
	h.data[k] = v
}

func (h *HashTable) Del(k string) {
	delete(h.data, k)
}

func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]string),
	}
}
