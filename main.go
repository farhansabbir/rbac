package main

import (
	"fmt"

	"github.com/farhansabbir/rbac/lib"
)

func main() {
	var resource lib.Resource = lib.NewRule("allow-user-read", "user read", []string{}, []lib.Verbs{lib.List}, lib.Allow)
	fmt.Println(resource.GetResourceID())
	fmt.Println(resource.GetResourceType())
	fmt.Println(resource.IsActive())
	fmt.Println(resource.GetResourceCreatedAt())
	rule := resource.(*lib.Rule)

	fmt.Println(rule.TargetResourceIDs)

}
