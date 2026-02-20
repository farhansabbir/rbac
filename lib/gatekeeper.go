package lib

import (
	"fmt"
	"sync/atomic"

	"github.com/farhansabbir/rbac/core"
)

var (
	Users    []*core.User
	Profiles []*core.Profile
	Rules    []*core.Rule
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
	// 1. Basic Validation
	if requestcontext.RequestResourceType == core.ResourceTypeNone {
		g.incrementRequestsRejected()
		return false, fmt.Errorf("RequestResourceType cannot be ResourceTypeNone")
	}

	// 2. Resolve User
	user, err := GetUserByID(requestcontext.PrincipalID)
	if err != nil {
		g.incrementRequestsRejected()
		return false, err
	}
	if !user.IsActive() {
		g.incrementRequestsRejected()
		return false, fmt.Errorf("User %d is not active", user.GetResourceID())
	}

	// 3. Get Active Profiles
	profiles, err := GetActiveProfilesByUserID(user.GetResourceID())
	if err != nil {
		// No active profiles = Implicit Deny
		g.incrementRequestsRejected()
		return false, err
	}

	// We assume "Implicit Deny" by default.
	// We only switch this to true if we find an explicit Allow.
	requestAllowed := false

	// 4. Evaluate Profiles
	for _, prof := range profiles {
		// OPTIMIZATION: Only fetch rules that match the Requested Resource Type OR Global Rules.
		// This replaces GetActiveRulesByProfileID which was inefficient.
		relevantRules := prof.GetAssociatedRules(requestcontext.RequestResourceType)
		globalRules := prof.GetAssociatedRules(core.ResourceTypeAll)

		allRulesToCheck := append(relevantRules, globalRules...)

		for _, rule := range allRulesToCheck {
			// Skip inactive rules
			if !rule.IsActive() {
				continue
			}

			// Check match (returns bool now, clearer logic)
			if RuleMatches(rule, requestcontext) {
				// LOGIC: Deny-Overrides-Allow
				switch rule.GetRuleAction() {
				case core.ActionDeny:
					// CRITICAL FIX: Return immediately on Deny.
					// Do NOT continue checking other rules.
					fmt.Printf("Explicit DENY by rule ID: %d\n", rule.GetResourceID())
					g.incrementRequestsRejected()
					return false, nil

				case core.ActionAllow:
					// Mark as allowed, but KEEP CHECKING in case a later rule Denies it.
					fmt.Printf("Matched ALLOW rule ID: %d\n", rule.GetResourceID())
					requestAllowed = true

				case core.ActionAllowAndForwardToNextRule:
					// Treat as Allow for now (forwarding logic would go here)
					fmt.Printf("Matched ALLOWandForwardToNextRule rule ID: %d, forward to rule ID: %d\n", rule.GetResourceID(), rule.GetResourceForwardRuleID())
					requestAllowed = true
				}
			}
		}
	}

	// 5. Final Decision
	if requestAllowed {
		g.incrementRequestsAccepted()
		return true, nil
	}

	// Implicit Deny
	g.incrementRequestsRejected()
	return false, nil
}

// RuleMatches returns true if the rule APPLIES to the request.
// It uses pointers (*Rule) to avoid copying the struct.
func RuleMatches(rule *core.Rule, ctx *RequestContext) bool {
	// 1. Check Resource Type
	// (Already filtered by optimization, but safety check)
	if rule.GetTargetResourceType() != core.ResourceTypeAll &&
		rule.GetTargetResourceType() != ctx.RequestResourceType {
		return false
	}

	// 2. Check Resource ID (String match or Wildcard)
	targetID := rule.GetTargetResourceID()
	reqIDStr := fmt.Sprint(ctx.RequestResourceID)

	if targetID != core.ResourceIDAll && targetID != reqIDStr {
		return false
	}

	// 3. Check Verb (Bitwise Match)
	ruleVerb := rule.GetVerb()

	// CRITICAL FIX: Use Bitwise AND (&)
	// If rule is VerbAll, it matches everything.
	// Otherwise, check if the requested bit is set in the rule.
	if ruleVerb != core.VerbAll && (ruleVerb&ctx.RequestVerb) == 0 {
		return false
	}

	return true
}

// --- Helper Functions (No changes needed, kept for context) ---

func (g *Gatekeeper) GetGKStats() (uint64, uint64) {
	return g.requestsRejected, g.requestsAccepted
}

func GetUserByID(id uint64) (*core.User, error) {
	for _, user := range Users {
		if user.GetResourceID() == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("User with ID %d not found", id)
}

func GetActiveProfilesByUserID(userid uint64) ([]core.Profile, error) {
	var profiles []core.Profile
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

func GetUserProfilesFromUserID(userid uint64) ([]core.Profile, error) {
	var profiles []core.Profile
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

func GetProfileByID(id uint64) (*core.Profile, error) {
	for _, profile := range Profiles {
		if profile.GetResourceID() == id {
			return profile, nil
		}
	}
	return nil, fmt.Errorf("Profile with ID %d not found", id)
}
