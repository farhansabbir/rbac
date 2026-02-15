package controllers

import (
	"context"
	"fmt"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/farhansabbir/rbac/lib"
)

var (
	globalController *Controller
	initOnce         sync.Once
	wg               sync.WaitGroup
)

// Controller is the main entry point (Singleton)
type Controller struct {
	ucinstance *UserController
	ctx        context.Context
	cancel     context.CancelFunc
}

// UserController manages user state and events
type UserController struct {
	id     uint64
	mux    sync.RWMutex
	users  map[uint64]*lib.User
	events chan string
}

// GetController initializes the system once and returns the singleton
func GetController() *Controller {
	initOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())

		globalController = &Controller{
			ctx:    ctx,
			cancel: cancel,
			ucinstance: &UserController{
				id:     xxhash.Sum64String("user_controller_singleton"),
				users:  make(map[uint64]*lib.User),
				events: make(chan string, 100), // Buffered channel
			},
		}

		// Start background processes
		globalController.StartUserControllerEventLoop()
		fmt.Println("System Controller initialized")
	})
	return globalController
}

// GetUserController returns the sub-controller
func (c *Controller) GetUserController() *UserController {
	return c.ucinstance
}

// StartUserControllerEventLoop runs in the background
func (c *Controller) StartUserControllerEventLoop() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("UC event loop running...")

		for {
			select {
			case <-c.ctx.Done():
				fmt.Println("UC event loop received shutdown signal")
				return
			case msg, ok := <-c.ucinstance.events:
				if !ok {
					return
				}
				fmt.Printf("[EVENT LOG]: %s\n", msg)
			}
		}
	}()
}

// Stop safely shuts down all background loops
func (c *Controller) Stop() {
	c.cancel() // Trigger context cancellation
	close(c.ucinstance.events)
	wg.Wait() // Wait for goroutines to finish
	fmt.Println("All systems stopped.")
}

// --- UserController Methods ---

func (uc *UserController) CreateUser(name, description, email string) *lib.User {
	u := lib.NewUser(name, description, email)

	uc.mux.Lock()
	uc.users[u.GetResourceID()] = u
	uc.mux.Unlock()

	uc.events <- fmt.Sprintf("User Created: %s (ID: %d)", u.GetResourceName(), u.GetResourceID())
	return u
}

func (uc *UserController) GetUser(id uint64) *lib.User {
	uc.mux.RLock()
	defer uc.mux.RUnlock()
	return uc.users[id]
}

func (uc *UserController) DeleteUser(id uint64) bool {
	uc.mux.Lock()
	defer uc.mux.Unlock()

	if user, ok := uc.users[id]; ok {
		user.SoftDelete()
		delete(uc.users, id)
		uc.events <- fmt.Sprintf("User Deleted: %d", id)
		return true
	}
	return false
}

func (uc *UserController) ListUsers() []*lib.User {
	uc.mux.RLock()
	defer uc.mux.RUnlock()

	list := make([]*lib.User, 0, len(uc.users))
	for _, u := range uc.users {
		list = append(list, u)
	}
	return list
}
