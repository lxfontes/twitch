package store

import "testing"

func TestMemory(t *testing.T) {
	performStoreTest(t, NewMemoryStore())
}
