package controllers

import (
	"sync"

	"github.com/farhansabbir/rbac/core"
)

type ProfileController struct {
	id       uint64
	mux      sync.RWMutex
	profiles map[uint64]*core.Profile
	events   chan string
}
