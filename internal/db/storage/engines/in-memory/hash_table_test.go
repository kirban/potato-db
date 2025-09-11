package inmemory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// helper to check key existence
func existsInHashTable(ht *HashTable, key string) bool {
	_, ok := ht.data[key]

	return ok
}

func TestHashTable_Get(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		keyArg         string
		expectedVal    string
		expectedExists bool
	}{
		"get existing key": {
			keyArg:         "key",
			expectedVal:    "value",
			expectedExists: true,
		},
		"get non existing key": {
			keyArg:         "nkey",
			expectedVal:    "",
			expectedExists: false,
		},
	}

	ht := &HashTable{
		data: map[string]string{
			"key": "value",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v, exists := ht.Get(tc.keyArg)

			assert.Equal(t, tc.expectedVal, v)
			assert.Equal(t, tc.expectedExists, exists)
		})
	}
}

func TestHashTable_Set(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		keyArg string
		value  string
	}{
		"set updates existing key": {
			keyArg: "newKey",
			value:  "newValue",
		},
		"set adds key": {
			keyArg: "key123",
			value:  "123",
		},
	}

	ht := &HashTable{
		data: map[string]string{
			"key": "value",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ht.Set(tc.keyArg, tc.value)
			assert.True(t, existsInHashTable(ht, tc.keyArg))
			assert.Equal(t, tc.value, ht.data[tc.keyArg])
		})
	}
}

func TestHashTable_Del(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		keyArg string
	}{
		"removes existing key": {
			keyArg: "key",
		},
		"removes non existing key": {
			keyArg: "nkey",
		},
	}

	ht := &HashTable{
		data: map[string]string{
			"key": "value",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ht.Del(tc.keyArg)
			assert.False(t, existsInHashTable(ht, tc.keyArg))
			assert.Empty(t, ht.data[tc.keyArg])
		})
	}
}
