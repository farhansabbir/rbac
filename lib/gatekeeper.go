package lib

import (
	"fmt"
)

var (
	Users    []*User
	Profiles []*Profile
	Rules    []*Rule
)

type Gatekeeper struct {
	requestsRejected uint64
	requestsAccepted uint64
}

func NewGatekeeper() *Gatekeeper {
	return &Gatekeeper{
		requestsRejected: 0,
		requestsAccepted: 0,
	}
}

func (g *Gatekeeper) incrementRequestsRejected() {
	g.requestsRejected++
}

func (g *Gatekeeper) incrementRequestsAccepted() {
	g.requestsAccepted++
}

func (g *Gatekeeper) IsRequestAllowed(requestcontext *RequestContext) (bool, error) {
	if requestcontext.RequestResourceType == ResourceTypeNone {
		g.incrementRequestsRejected()
		return false, fmt.Errorf("RequestResourceType cannot be ResourceTypeNone")
	}
	for _, profile := range requestcontext.PrincipalProfiles {
		fmt.Println(profile)
	}
	g.incrementRequestsAccepted()
	return true, nil
}

func (g *Gatekeeper) GetGKStats() (uint64, uint64) {
	return g.requestsRejected, g.requestsAccepted
}

func (g *Gatekeeper) GetUsers() []User {
	users := make([]User, len(Users))
	for i, user := range Users {
		users[i] = *user
	}
	return users
}

func (g *Gatekeeper) GetProfiles() []Profile {
	profiles := make([]Profile, len(Profiles))
	for i, profile := range Profiles {
		profiles[i] = *profile
	}
	return profiles
}

func (g *Gatekeeper) GetRules() []Rule {
	rules := make([]Rule, len(Rules))
	for i, rule := range Rules {
		rules[i] = *rule
	}
	return rules
}

func GetUserByID(id uint64) (*User, error) {
	for _, user := range Users {
		if user.GetResourceID() == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("User with ID %d not found", id)
}

func GetProfileByID(id uint64) (*Profile, error) {
	for _, profile := range Profiles {
		if profile.GetResourceID() == id {
			return profile, nil
		}
	}
	return nil, fmt.Errorf("Profile with ID %d not found", id)
}

func GetRuleByID(id uint64) (*Rule, error) {
	for _, rule := range Rules {
		if rule.GetResourceID() == id {
			return rule, nil
		}
	}
	return nil, fmt.Errorf("Profile with ID %d not found", id)
}

func GetUserProfilesByUserID(userid uint64) ([]Profile, error) {
	var profiles []Profile
	found := false
	for _, user := range Users {
		if user.GetResourceID() == userid {
			for _, profile := range user.GetProfiles() {
				profiles = append(profiles, profile)
			}
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("User with ID %d not found", userid)
	}
	return profiles, nil
}
