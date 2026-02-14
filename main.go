package main

import (
	"fmt"

	"github.com/farhansabbir/rbac/lib"
)

func main() {
	var resource lib.Resource = lib.NewEmptyRule("allow-user-read")
	rule := resource.(*lib.Rule)
	rule.UpdateName("hello-rule")

	resource = lib.NewProfile("Business User profile", "user profile")
	profile := resource.(*lib.Profile)

	rule.AddTargetResourceID(profile.ResourceType, fmt.Sprintf("%d", profile.ID))
	profile.AddRule(rule)
	fmt.Println("Profile ID:" + fmt.Sprintf("%d", profile.GetResourceID()))
	fmt.Println("Resource Type:" + fmt.Sprintf("%s", profile.ResourceType))
	fmt.Println("Resource Name:" + profile.GetResourceName())
	fmt.Println("'" + profile.GetRules()[0].GetRuleAsDSL() + "' is valid? " + fmt.Sprintf("%t", profile.GetRules()[0].IsValidRule()))
}
