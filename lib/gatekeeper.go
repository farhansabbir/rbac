package lib

import (
	"fmt"
	"sync/atomic"
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
	atomic.AddUint64(&g.requestsRejected, 1)
}

func (g *Gatekeeper) incrementRequestsAccepted() {
	atomic.AddUint64(&g.requestsAccepted, 1)
}

func (g *Gatekeeper) IsRequestAllowed(requestcontext *RequestContext) (bool, error) {
	if requestcontext.RequestResourceType == ResourceTypeNone {
		g.incrementRequestsRejected()
		return false, fmt.Errorf("RequestResourceType cannot be ResourceTypeNone")
	}
	user, err := GetUserByID(requestcontext.PrincipalID)
	if err != nil {
		return false, err
	}
	if !user.IsActive() {
		g.incrementRequestsRejected()
		return false, fmt.Errorf("User is not active")
	}
	profiles, err := GetUserProfilesFromUserID(user.GetResourceID())
	if err != nil {
		return false, err
	}
	fmt.Println("Total profiles: " + fmt.Sprint(len(profiles)))
	for _, profile := range profiles {
		fmt.Println(profile)
	}
	g.incrementRequestsAccepted()
	return true, nil
}

func (g *Gatekeeper) GetGKStats() (uint64, uint64) {
	return g.requestsRejected, g.requestsAccepted
}

// these will need to move to a separate file of corresponding controllers
func (g *Gatekeeper) GetAllUsers() []User {
	users := make([]User, len(Users))
	for i, user := range Users {
		users[i] = *user
	}
	return users
}

func GetUserByID(id uint64) (*User, error) {
	for _, user := range Users {
		if user.GetResourceID() == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("User with ID %d not found", id)
}

func GetUserProfilesFromUserID(userid uint64) ([]Profile, error) {
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

func (g *Gatekeeper) GetAllProfiles() []Profile {
	profiles := make([]Profile, len(Profiles))
	for i, profile := range Profiles {
		profiles[i] = *profile
	}
	return profiles
}

func GetProfileByID(id uint64) (*Profile, error) {
	for _, profile := range Profiles {
		if profile.GetResourceID() == id {
			return profile, nil
		}
	}
	return nil, fmt.Errorf("Profile with ID %d not found", id)
}

func (g *Gatekeeper) GetAllRules() []Rule {
	rules := make([]Rule, len(Rules))
	for i, rule := range Rules {
		rules[i] = *rule
	}
	return rules
}

func GetRuleByID(id uint64) (*Rule, error) {
	for _, rule := range Rules {
		if rule.GetResourceID() == id {
			return rule, nil
		}
	}
	return nil, fmt.Errorf("Profile with ID %d not found", id)
}

func GetRulesByUserID(userID uint64) ([]*Rule, error) {
	var rules []*Rule
	profiles, err := GetUserProfilesFromUserID(userID)
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		for _, rulelist := range profile.GetRuleMap() {
			rules = append(rules, rulelist...)
		}
	}
	return rules, nil
}

func GetRulesByUserIDAndResourceType(userID uint64, resourcetype ResourceType) ([]*Rule, error) {
	profiles, err := GetUserProfilesFromUserID(userID)
	if err != nil {
		return nil, err
	}
	var rules []*Rule
	for _, profile := range profiles {
		rules = append(rules, profile.GetRuleMap()[uint32(resourcetype)]...)
	}
	return rules, nil
}
