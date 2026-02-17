package main

import (
	"fmt"

	"github.com/farhansabbir/rbac/lib"
)

func main() {

	gk := lib.NewGatekeeper()
	readallrule := lib.NewEmptyRule("readonly-rule")
	readallrule.UpdateVerb(lib.VerbAll)
	if _, err := readallrule.SetTargetResourceTypeAndID(lib.ResourceTypeProfile, lib.ResourceIDAll); err != nil {
		fmt.Println("Error setting target resource type and ID:", err)
		return
	}
	writerule := lib.NewEmptyRule("writeonly-rule")
	if _, err := writerule.SetTargetResourceType(lib.ResourceTypeAll); err != nil {
		fmt.Println("Error setting target resource type:", err)
		return
	}

	lib.Rules = append(lib.Rules, readallrule)
	lib.Rules = append(lib.Rules, writerule)

	mktprofile := lib.NewProfile("Business User profile", "Marketing user profile")

	mktprofile.AddRule(readallrule)
	mktprofile.AddRule(writerule)
	lib.Profiles = append(lib.Profiles, mktprofile)

	userjohn := lib.NewUser("John Doe", "Administrator", "john.doe@example.com")
	lib.Users = append(lib.Users, userjohn)
	userjane := lib.NewUser("Jane Doe", "Regular User", "jane.doe@example.com")
	lib.Users = append(lib.Users, userjane)
	userjohn.AddProfile(mktprofile)
	rctx, err := lib.NewRequestContext(userjohn.GetResourceID(), lib.ResourceTypeProfile, mktprofile.GetResourceID(), lib.VerbList, nil)
	if err != nil {
		fmt.Println("Error creating request context:", err)
		return
	}
	// fmt.Println(rctx)
	fmt.Println(gk.IsRequestAllowed(rctx))

}
