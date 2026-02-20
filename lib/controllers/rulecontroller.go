package controllers

import (
	"sync"

	"github.com/farhansabbir/rbac/core"
)

type RuleController struct {
	id     uint64
	mux    sync.RWMutex
	rules  map[uint64]*core.Rule
	events chan string
}
