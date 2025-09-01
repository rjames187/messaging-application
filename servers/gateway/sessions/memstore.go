package sessions

import "errors"

type MemoryStore struct {
	store map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{store: map[string]string{}}
}

func (m *MemoryStore) Get(key string) (string, error) {
	value, ok := m.store[key]
	if !ok {
		return "", errors.New("no matches to given token in session store")
	}
	return value, nil
}

func (m *MemoryStore) Set(key string, value string) error {
	m.store[key] = value
	return nil
}

func (m *MemoryStore) Delete(key string) error {
	_, ok := m.store[key]
	if !ok {
		return errors.New("no matches to given token in session store")
	}
	delete(m.store, key)
	return nil
}
