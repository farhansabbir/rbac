package controllers

import (
	"sync"

	"github.com/farhansabbir/rbac/lib"
)

type RuleController struct {
	id     uint64
	mux    sync.RWMutex
	rules  map[uint64]*lib.Rule
	events chan string
}
