package storage

import (
	"fmt"
	"sync"

	"github.com/tormoder/chat/common"
	pb "github.com/tormoder/chat/proto"
)

type UserStorage interface {
	AddUser(user User) error
	GetUser(nick string) (User, bool)
	UpdateUser(user User) error
	DeleteUser(nick string) error
	GetAllUsers() []User
	GetAllOnlineUsers() []User
	GetAllOnlineUsersDTO() []*pb.User
	CheckCredentials(*pb.Credentials) (User, error)
}

type InMemoryUserStorage struct {
	users map[string]User
	mu    sync.RWMutex
}

func NewInMemoryUserStorage() *InMemoryUserStorage {
	return &InMemoryUserStorage{
		users: make(map[string]User),
	}
}

func (us *InMemoryUserStorage) AddUser(user User) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	u, found := us.users[user.Nick]
	if found {
		return fmt.Errorf("user %q already exists", u.Nick)
	}
	us.users[user.Nick] = user
	return nil
}

func (us *InMemoryUserStorage) GetUser(nick string) (User, bool) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	u, found := us.users[nick]
	return u, found
}

func (us *InMemoryUserStorage) UpdateUser(user User) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.users[user.Nick] = user
	return nil
}

func (us *InMemoryUserStorage) DeleteUser(nick string) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	delete(us.users, nick)
	return nil
}

func (us *InMemoryUserStorage) GetAllUsers() []User {
	return us.getAllUsers(true)
}

func (us *InMemoryUserStorage) GetAllOnlineUsers() []User {
	return us.getAllUsers(false)
}

func (us *InMemoryUserStorage) getAllUsers(offline bool) []User {
	us.mu.RLock()
	defer us.mu.RUnlock()
	var users []User
	for _, user := range us.users {
		if !offline && !user.Online {
			continue
		}
		users = append(users, user)
	}
	return users
}

func (us *InMemoryUserStorage) GetAllOnlineUsersDTO() []*pb.User {
	us.mu.RLock()
	defer us.mu.RUnlock()
	var users []*pb.User
	for _, user := range us.users {
		if !user.Online {
			continue
		}
		u := user.User
		users = append(users, &u)
	}
	return users
}

func (us *InMemoryUserStorage) CheckCredentials(creds *pb.Credentials) (User, error) {
	user, found := us.GetUser(creds.Nick)
	if !found {
		return User{}, common.AuthenticationError("user not found")
	}
	if !user.Online {
		return User{}, common.AuthenticationError("user not logged-in")
	}

	// TODO: Authentication

	return user, nil
}
