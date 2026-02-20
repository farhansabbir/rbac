package controllers

import (
	"fmt"
	"sync"

	"github.com/farhansabbir/rbac/core"
)

// UserController manages user state and events
type UserController struct {
	id     uint64
	mux    sync.RWMutex
	users  map[uint64]*core.User
	events chan string
}

// --- UserController Methods ---

func (uc *UserController) CreateUser(name, description, email string) *core.User {
	u := core.NewUser(name, description, email)

	uc.mux.Lock()
	uc.users[u.GetResourceID()] = u
	uc.mux.Unlock()

	uc.events <- fmt.Sprintf("User Created: %s (ID: %d)", u.GetResourceName(), u.GetResourceID())
	return u
}

func (uc *UserController) GetUser(id uint64) *core.User {
	uc.mux.RLock()
	defer uc.mux.RUnlock()
	return uc.users[id]
}

func (uc *UserController) DeleteUser(id uint64) bool {
	uc.mux.Lock()
	defer uc.mux.Unlock()

	if user, ok := uc.users[id]; ok {
		user.SoftDelete()
		uc.events <- fmt.Sprintf("User Deleted: %d", id)
		return true
	}
	return false
}

func (uc *UserController) ListUsers() []*core.User {
	uc.mux.RLock()
	defer uc.mux.RUnlock()

	list := make([]*core.User, 0, len(uc.users))
	for _, u := range uc.users {
		list = append(list, u)
	}
	return list
}

func (uc *UserController) ListActiveUsers() []*core.User {
	uc.mux.RLock()
	defer uc.mux.RUnlock()

	list := make([]*core.User, 0, len(uc.users))
	for _, u := range uc.users {
		if u.IsActive() {
			list = append(list, u)
		}
	}
	return list
}
