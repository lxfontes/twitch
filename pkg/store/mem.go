package store

import (
	"encoding/json"
	"sync"
	"time"
)

var _ Store = &memStore{}
var _ Namespace = storage{}

type storage map[string][]byte

func (s storage) FindByID(id string, out interface{}) error {
	rawItem, ok := s[id]
	if !ok {
		return ErrItemNotFound
	}

	return json.Unmarshal(rawItem, out)
}

func (s storage) Delete(id string) error {
	delete(s, id)
	return nil
}

func (s storage) Save(id string, item interface{}, expiration time.Time) error {
	rw, err := json.Marshal(item)
	if err != nil {
		return err
	}

	s[id] = rw
	return nil
}

type memStore struct {
	things map[string]storage
	mtx    sync.Mutex
}

func NewMemoryStore() *memStore {
	return &memStore{
		things: map[string]storage{},
	}
}

func (ms *memStore) Namespace(name string) Namespace {
	namespace, ok := ms.things[name]
	if !ok {
		namespace = storage{}
		ms.things[name] = namespace
	}

	return namespace
}
