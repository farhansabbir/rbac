package main

import (
	"fmt"

	"github.com/farhansabbir/rbac/lib"
)

func main() {
	rule := lib.NewEmptyRule("readonly-rule")
	rule.UpdateVerb(lib.VerbList)
	rule2 := lib.NewEmptyRule("writeonly-rule")

	profile := lib.NewProfile("Business User profile", "user profile")
	rule.AddTargetResourceID(profile.ResourceType, fmt.Sprintf("%d", profile.ID))
	profile.AddRule(rule)
	profile.AddRule(rule2)
	fmt.Println(profile.GetProfileAsJSON())
}
