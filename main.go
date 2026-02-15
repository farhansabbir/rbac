package main

import (
	"encoding/json"
	"fmt"

	"github.com/farhansabbir/rbac/lib"
)

func main() {
	var users []*lib.User
	var profiles []*lib.Profile
	var rules []*lib.Rule

	// gk := lib.NewGatekeeper()
	rule := lib.NewEmptyRule("readonly-rule")
	rule.UpdateVerb(lib.VerbList)
	rules = append(rules, rule)
	rule2 := lib.NewEmptyRule("writeonly-rule")
	rules = append(rules, rule2)

	profile := lib.NewProfile("Business User profile", "user profile")
	rule.AddTargetResourceID(profile.GetResourceType(), fmt.Sprintf("%d", profile.GetResourceID()))
	profile.AddRule(rule)
	profile.AddRule(rule2)
	profiles = append(profiles, profile)

	user := lib.NewUser("John Doe", "Administrator", "john.doe@example.com")
	users = append(users, user)
	user1 := lib.NewUser("Jane Doe", "Regular User", "jane.doe@example.com")
	users = append(users, user1)
	user.AddProfile(profile)
	js, _ := json.Marshal(users)
	// rctx, err := lib.NewRequestContext(user.GetResourceID(), user.GetProfiles(), lib.ResourceTypeURL, profile.GetResourceID(), lib.VerbList, nil)
	// if err != nil {
	// 	fmt.Println("Error creating request context:", err)
	// 	return
	// }
	// fmt.Println(gk.IsRequestAllowed(rctx))
	fmt.Println(string(js))

}
