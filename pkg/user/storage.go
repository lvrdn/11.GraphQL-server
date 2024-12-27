package user

import "sync"

type memoryStorage struct {
	Users []*User
	mu    *sync.RWMutex
}

func NewUserStorage() *memoryStorage {
	return &memoryStorage{
		Users: make([]*User, 0),
		mu:    &sync.RWMutex{},
	}
}

func (ms *memoryStorage) Add(user *User) (int, error) {

	ms.mu.Lock()
	id := len(ms.Users) + 1
	user.ID = id
	ms.Users = append(ms.Users, user)
	ms.mu.Unlock()

	return id, nil
}
