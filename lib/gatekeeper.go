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
	requestAccepted := false
	// If requested resource type is ResourceTypeNone, reject the request
	if requestcontext.RequestResourceType == ResourceTypeNone {
		g.incrementRequestsRejected()
		return false, fmt.Errorf("RequestResourceType cannot be ResourceTypeNone")
	}
	// If user is not found, reject the request
	user, err := GetUserByID(requestcontext.PrincipalID)
	if err != nil {
		g.incrementRequestsRejected()
		return false, err
	}
	// If user is not active, reject the request
	if !user.IsActive() {
		g.incrementRequestsRejected()
		return false, fmt.Errorf("User %d is not active", user.GetResourceID())
	}
	// user found, now check if user has any active profiles
	if profiles, err := GetActiveProfilesByUserID(user.GetResourceID()); err != nil {
		g.incrementRequestsRejected()
		return false, err
	} else {
		// user has active profiles
		for _, prof := range profiles {
			// now check if user has any active rules in that active profiles
			if rules, err := GetActiveRulesByProfileID(prof.GetResourceID()); err != nil {
				continue
			} else {
				// user has active rules in active profiles for active user
				// now browse the rules and match
				for _, rule := range rules {
					// fmt.Print("Username: " + user.GetResourceName())
					err := RuleMatcher(rule, requestcontext)
					if err != nil {
						continue
					} else {
						fmt.Println("Matched rule ID:", rule.JSON())
						if rule.GetRuleAction() == ActionDeny {
							// g.incrementRequestsRejected()
							requestAccepted = false
							continue
						} else if rule.GetRuleAction() == ActionAllowAndForwardToNextRule {
							g.incrementRequestsAccepted()
							fmt.Println("Chained rule for next processing")
							return true, nil
						} else if rule.GetRuleAction() == ActionAllow {
							g.incrementRequestsAccepted()
							fmt.Println("Allowed")
							return true, nil
						}
					}
				}
			}
		}
	}
	g.incrementRequestsRejected()
	return false, nil
}

func RuleMatcher(rule Rule, requestcontext *RequestContext) error {
	// check if requested resource matches rule target resource type
	if uint64(rule.GetTargetResourceType()) != uint64(requestcontext.RequestResourceType) { //requested resource type and rule resource type are not same?
		if rule.GetTargetResourceType() != ResourceTypeAll { // check if rule target resource type is ResourceTypeAll, this means that the rule applies to all resource types
			// and also rule resource type is not ResourceTypeAll? then we need to fail
			return fmt.Errorf("Resource type mismatch. Req: %s, Rule: %s", requestcontext.RequestResourceType, rule.GetTargetResourceType())
		} else { // this means that the rule applies to all resource types, we need to ignore matching resourceID and move forward
			// ignore matching resourceID and move forward to verb matching
			//
		}
	} else { // this means req resourcetype matches with rule resource type, we need to match req resourceID with rule target resource ID
		if fmt.Sprint(requestcontext.RequestResourceID) != rule.GetTargetResourceID() &&
			rule.GetTargetResourceID() != ResourceIDAll { //requested resource id and rule resource id are not same?
			return fmt.Errorf("Resource ID mismatch. Req: %d, Rule: %s", requestcontext.RequestResourceID, rule.GetTargetResourceID())
		}
	}
	// at this point we know that the resource type and resource id match with the rule
	// proceed to verb matching
	if requestcontext.RequestVerb != rule.GetVerb() &&
		rule.GetVerb() != VerbAll {
		return fmt.Errorf("Verb mismatch. Req: %s, Rule: %s", requestcontext.RequestVerb, rule.GetVerb())
	}
	// if fmt.Sprint(requestcontext.RequestResourceID) != rule.GetTargetResourceID() &&
	// 	rule.GetTargetResourceID() != ResourceIDAll { //requested resource id and rule resource id are not same?
	// 	return false, fmt.Errorf("Resource ID mismatch. Req: %d, Rule: %s", requestcontext.RequestResourceID, rule.GetTargetResourceID())
	// }
	return nil
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

func GetActiveProfilesByUserID(userid uint64) ([]Profile, error) {
	var profiles []Profile
	found := false
	userprofiles, err := GetUserProfilesFromUserID(userid)
	if err != nil {
		return nil, err
	}
	for _, profile := range userprofiles {
		if profile.IsActive() {
			profiles = append(profiles, profile)
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("User with ID %d does not have active profiles", userid)
	}
	return profiles, nil
}

func GetActiveRulesByProfileID(profileID uint64) ([]Rule, error) {
	var rules []Rule
	profile, err := GetProfileByID(profileID)
	if err != nil {
		return nil, err
	}
	found := false
	for _, allrules := range profile.GetRuleMap() {
		for _, rule := range allrules {
			if rule.IsActive() {
				rules = append(rules, *rule)
				found = true
			}
		}
	}
	if !found {
		return nil, fmt.Errorf("Profile with ID %d does not have active rules", profileID)
	}
	return rules, nil
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
	found := false
	for _, profile := range profiles {
		if resourcerules, ok := profile.GetRuleMap()[uint32(resourcetype)]; ok {
			found = true
			rules = append(rules, resourcerules...)
		} else {
			continue
		}
	}
	if !found {
		return nil, fmt.Errorf("No rules found for user '%d' and resource type '%s'", userID, resourcetype)
	}
	return rules, nil
}
