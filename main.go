package main

import (
	"fmt"

	"github.com/farhansabbir/rbac/lib"
)

func main() {
	gk := lib.NewGatekeeper()
	rule := lib.NewEmptyRule("readonly-rule")
	rule.UpdateVerb(lib.VerbList)

	rule2 := lib.NewEmptyRule("writeonly-rule")

	profile := lib.NewProfile("Business User profile", "user profile")
	rule.AddTargetResourceID(profile.GetResourceType(), fmt.Sprintf("%d", profile.GetResourceID()))
	profile.AddRule(rule)
	profile.AddRule(rule2)
	user := lib.NewUser("John Doe", "Administrator", "john.doe@example.com")
	user.AddProfile(profile)
	fmt.Println(user.GetProfiles())

	rctx, err := lib.NewRequestContext(user.GetResourceID(), user.GetProfiles(), lib.ResourceTypeURL, profile.GetResourceID(), lib.VerbList, nil)
	if err != nil {
		fmt.Println("Error creating request context:", err)
		return
	}
	fmt.Println(gk.IsRequestAllowed(rctx))
}
