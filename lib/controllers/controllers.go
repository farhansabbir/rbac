package controllers

import (
	"context"
	"fmt"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/farhansabbir/rbac/core"
)

var (
	globalController *Controller
	initOnce         sync.Once
	wg               sync.WaitGroup
)

// Controller is the main entry point (Singleton)
type Controller struct {
	ucinstance *UserController
	pcinstance *ProfileController
	rcinstance *RuleController
	ctx        context.Context
	cancel     context.CancelFunc
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
				users:  make(map[uint64]*core.User),
				events: make(chan string, 100), // Buffered channel
			},
			pcinstance: &ProfileController{
				id:       xxhash.Sum64String("profile_controller_singleton"),
				profiles: make(map[uint64]*core.Profile),
				events:   make(chan string, 100), // Buffered channel
			},
			rcinstance: &RuleController{
				id:     xxhash.Sum64String("rule_controller_singleton"),
				rules:  make(map[uint64]*core.Rule),
				events: make(chan string, 100), // Buffered channel
			},
		}

		// Start background processes
		globalController.startEventLoop()
		fmt.Println("System Controller initialized")
	})
	return globalController
}

// GetUserController returns the sub-controller
func (c *Controller) GetUserController() *UserController {
	return c.ucinstance
}

// GetProfileController returns the sub-controller
func (c *Controller) GetProfileController() *ProfileController {
	return c.pcinstance
}

// GetRuleController returns the sub-controller
func (c *Controller) GetRuleController() *RuleController {
	return c.rcinstance
}

// StartEventLoop runs in the background
func (c *Controller) startEventLoop() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Controller event loop running...")

		for {
			select {
			case <-c.ctx.Done():
				fmt.Println("Controller event loop received shutdown signal")
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
