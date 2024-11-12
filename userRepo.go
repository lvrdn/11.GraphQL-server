package main

type UserStorageMemory struct {
	Users []*User
}

func NewUserStorage() *UserStorageMemory {
	return &UserStorageMemory{
		Users: make([]*User, 0),
	}
}

func (us *UserStorageMemory) Add(user *User) (int, error) {
	id := len(us.Users) + 1

	user.ID = id

	us.Users = append(us.Users, user)

	return id, nil
}
