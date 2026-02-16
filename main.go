package main

import (
	"fmt"

	"github.com/farhansabbir/rbac/lib"
)

func main() {

	gk := lib.NewGatekeeper()
	rule := lib.NewEmptyRule("readonly-rule")
	rule.UpdateVerb(lib.VerbList)
	lib.Rules = append(lib.Rules, rule)
	rule2 := lib.NewEmptyRule("writeonly-rule")
	lib.Rules = append(lib.Rules, rule2)

	profile := lib.NewProfile("Business User profile", "user profile")
	rule.AddTargetResourceID(profile.GetResourceType(), fmt.Sprintf("%d", profile.GetResourceID()))
	profile.AddRule(rule)
	profile.AddRule(rule2)
	lib.Profiles = append(lib.Profiles, profile)

	user := lib.NewUser("John Doe", "Administrator", "john.doe@example.com")
	lib.Users = append(lib.Users, user)
	user1 := lib.NewUser("Jane Doe", "Regular User", "jane.doe@example.com")
	lib.Users = append(lib.Users, user1)
	user.AddProfile(profile)
	rctx, err := lib.NewRequestContext(user.GetResourceID(), lib.ResourceTypeURL, profile.GetResourceID(), lib.VerbList, nil)
	if err != nil {
		fmt.Println("Error creating request context:", err)
		return
	}
	// fmt.Println(rctx)
	fmt.Println(gk.IsRequestAllowed(rctx))

}
