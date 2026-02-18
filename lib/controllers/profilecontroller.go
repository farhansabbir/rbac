package controllers

import (
	"sync"

	"github.com/farhansabbir/rbac/lib"
)

type ProfileController struct {
	id       uint64
	mux      sync.RWMutex
	profiles map[uint64]*lib.Profile
	events   chan string
}
