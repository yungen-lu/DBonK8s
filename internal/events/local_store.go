package events

import (
	"errors"
	"sync"
)

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

func (l *LocalStore) GetUser(id string) (*User, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if u, ok := l.users[id]; ok {
		return u, nil
	} else {
		return nil, errors.New("can't find user")
	}
}

func (l *LocalStore) GetAllUser() ([]*User, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if len(l.users) == 0 {
		return nil, errors.New("no user founded")
	}
	var users []*User
	for _, u := range l.users {
		users = append(users, u)
	}
	return users, nil
}
