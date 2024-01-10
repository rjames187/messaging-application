package sessions

import "errors"

type MemoryStore struct {
	store map[string]int
}

func (m MemoryStore) New() MemoryStore {
	return MemoryStore{map[string]int{}}
}

func (m *MemoryStore) Get(sessionID string) (int, error) {
	userID := m.store[sessionID]
	if userID == 0 {
		return 0, errors.New("given session id is not in the session store")
	}
	return userID, nil
}

func (m *MemoryStore) Set(sessionID string, userID int) error {
	m.store[sessionID] = userID
	return nil
}

func (m *MemoryStore) Delete(sessionID string) error {
	userID := m.store[sessionID]
	if userID == 0 {
		return errors.New("given session id is not in the session store")
	}
	delete(m.store, sessionID)
	return nil
}