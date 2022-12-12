package events

import "sync"

type LocalStore struct {
	users map[string]*User
	mu    sync.RWMutex
}

func NewLocalStore() *LocalStore {
	return &LocalStore{
		users: make(map[string]*User),
	}
}
func (l *LocalStore) GetOrCreateUser(id string, con *Controller) (*User, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if u, ok := l.users[id]; ok {
		return u, nil
	}
	u := NewUser(id, con)
	l.users[id] = u
	return u, nil
}