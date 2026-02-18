package lib

import (
	"testing"
)

// Helper to reset global state between tests
func resetGlobals() {
	Users = []*User{}
	Profiles = []*Profile{}
	Rules = []*Rule{}
}

func TestGatekeeper_IsRequestAllowed_BasicAllow(t *testing.T) {
	resetGlobals()
	gk := NewGatekeeper()

	// 1. Setup Rule: Allow "Read" on "Profile" resources
	rule := NewEmptyRule("allow-read-profiles")
	rule.UpdateVerb(VerbRead)
	rule.SetTargetResourceTypeAndID(ResourceTypeProfile, ResourceIDAll)
	rule.ruleAction = ActionAllow

	// 2. Setup Profile & User
	profile := NewProfile("basic-profile", "test profile")
	profile.AddRule(rule)

	user := NewUser("John", "User", "john@example.com")
	user.AddProfile(profile)

	// Register in globals (simulating DB)
	Users = append(Users, user)

	// 3. Create Request: Can John Read a Profile?
	ctx, err := NewRequestContext(user.GetResourceID(), ResourceTypeProfile, 12345, VerbRead, nil)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	// 4. Assert
	allowed, err := gk.IsRequestAllowed(ctx)
	if err != nil {
		t.Fatalf("Gatekeeper error: %v", err)
	}
	if !allowed {
		t.Errorf("Expected ALLOW, got DENY")
	}
}

func TestGatekeeper_DenyOverridesAllow(t *testing.T) {
	resetGlobals()
	gk := NewGatekeeper()

	// Rule 1: Allow Read (The "Nice" Rule)
	allowRule := NewEmptyRule("allow-read")
	allowRule.UpdateVerb(VerbRead)
	allowRule.SetTargetResourceTypeAndID(ResourceTypeProfile, ResourceIDAll)
	allowRule.ruleAction = ActionAllow

	// Rule 2: Deny Read (The "Strict" Rule)
	denyRule := NewEmptyRule("deny-read")
	denyRule.UpdateVerb(VerbRead)
	denyRule.SetTargetResourceTypeAndID(ResourceTypeProfile, ResourceIDAll)
	denyRule.ruleAction = ActionDeny

	// Profile has BOTH rules
	profile := NewProfile("mixed-profile", "mixed")
	profile.AddRule(allowRule)
	profile.AddRule(denyRule)

	user := NewUser("Jane", "User", "jane@example.com")
	user.AddProfile(profile)
	Users = append(Users, user)

	// Request
	ctx, _ := NewRequestContext(user.GetResourceID(), ResourceTypeProfile, 999, VerbRead, nil)

	// Assert: Deny should win
	allowed, _ := gk.IsRequestAllowed(ctx)
	if allowed {
		t.Errorf("Expected DENY (due to Deny-Overrides-Allow), but got ALLOW")
	}
}

func TestGatekeeper_ResourceMismatch(t *testing.T) {
	resetGlobals()
	gk := NewGatekeeper()

	// Rule: Allow Read on URLS only
	rule := NewEmptyRule("allow-url-read")
	rule.UpdateVerb(VerbRead)
	rule.SetTargetResourceTypeAndID(ResourceTypeURL, ResourceIDAll)
	rule.ruleAction = ActionAllow

	profile := NewProfile("url-profile", "urls")
	profile.AddRule(rule)

	user := NewUser("Bob", "User", "bob@example.com")
	user.AddProfile(profile)
	Users = append(Users, user)

	// Request: Try to Read a PROFILE (Mismatch!)
	ctx, _ := NewRequestContext(user.GetResourceID(), ResourceTypeProfile, 123, VerbRead, nil)

	allowed, _ := gk.IsRequestAllowed(ctx)
	if allowed {
		t.Errorf("Expected DENY (Resource Type Mismatch), but got ALLOW")
	}
}

func TestGatekeeper_VerbBitmaskMatching(t *testing.T) {
	resetGlobals()
	gk := NewGatekeeper()

	// Rule: Allow Read OR List (Bitmask: 0001 | 0010 = 0011)
	rule := NewEmptyRule("read-list-rule")
	rule.UpdateVerb(VerbRead | VerbList)
	rule.SetTargetResourceTypeAndID(ResourceTypeProfile, ResourceIDAll)
	rule.ruleAction = ActionAllow

	profile := NewProfile("reader-profile", "reader")
	profile.AddRule(rule)

	user := NewUser("Alice", "User", "alice@example.com")
	user.AddProfile(profile)
	Users = append(Users, user)

	// Request 1: Ask for Read (Should Match)
	ctxRead, _ := NewRequestContext(user.GetResourceID(), ResourceTypeProfile, 555, VerbRead, nil)
	if allowed, _ := gk.IsRequestAllowed(ctxRead); !allowed {
		t.Errorf("Expected ALLOW for Read request, got DENY")
	}

	// Request 2: Ask for List (Should Match)
	ctxList, _ := NewRequestContext(user.GetResourceID(), ResourceTypeProfile, 555, VerbList, nil)
	if allowed, _ := gk.IsRequestAllowed(ctxList); !allowed {
		t.Errorf("Expected ALLOW for List request, got DENY")
	}

	// Request 3: Ask for Delete (Should NOT Match)
	ctxDelete, _ := NewRequestContext(user.GetResourceID(), ResourceTypeProfile, 555, VerbDelete, nil)
	if allowed, _ := gk.IsRequestAllowed(ctxDelete); allowed {
		t.Errorf("Expected DENY for Delete request, got ALLOW")
	}
}

func TestGatekeeper_ForwardingRule(t *testing.T) {
	resetGlobals()
	gk := NewGatekeeper()

	// Rule: Forwarding Action
	rule := NewEmptyRule("forwarding-rule")
	rule.UpdateVerb(VerbExecute)
	rule.SetTargetResourceTypeAndID(ResourceTypeProfile, ResourceIDAll)

	// Set Action to Forward with a dummy NextID
	rule.UpdateAction(ActionOption{
		Action:     ActionAllowAndForwardToNextRule,
		NextRuleID: 9999,
	})

	profile := NewProfile("forward-profile", "forward")
	profile.AddRule(rule)

	user := NewUser("Dave", "User", "dave@example.com")
	user.AddProfile(profile)
	Users = append(Users, user)

	// Request
	ctx, _ := NewRequestContext(user.GetResourceID(), ResourceTypeProfile, 101, VerbExecute, nil)

	// Assert
	allowed, _ := gk.IsRequestAllowed(ctx)
	if !allowed {
		t.Errorf("Expected ALLOW (via Forwarding Logic), got DENY")
	}

	// Check Stats
	rejected, accepted := gk.GetGKStats()
	if accepted != 1 {
		t.Errorf("Expected 1 accepted request, got %d", accepted)
	}
	if rejected != 0 {
		t.Errorf("Expected 0 rejected requests, got %d", rejected)
	}
}

func TestGatekeeper_WildcardIDMatching(t *testing.T) {
	resetGlobals()
	gk := NewGatekeeper()

	// Rule: Allow Read on Specific ID "100" only
	rule := NewEmptyRule("specific-id-rule")
	rule.UpdateVerb(VerbRead)
	rule.SetTargetResourceTypeAndID(ResourceTypeProfile, "100") // String ID
	rule.ruleAction = ActionAllow

	profile := NewProfile("strict-profile", "strict")
	profile.AddRule(rule)

	user := NewUser("Eve", "User", "eve@example.com")
	user.AddProfile(profile)
	Users = append(Users, user)

	// Request 1: ID 100 (Should Match)
	ctxMatch, _ := NewRequestContext(user.GetResourceID(), ResourceTypeProfile, 100, VerbRead, nil)
	if allowed, _ := gk.IsRequestAllowed(ctxMatch); !allowed {
		t.Errorf("Expected ALLOW for ID 100, got DENY")
	}

	// Request 2: ID 101 (Should Fail)
	ctxMismatch, _ := NewRequestContext(user.GetResourceID(), ResourceTypeProfile, 101, VerbRead, nil)
	if allowed, _ := gk.IsRequestAllowed(ctxMismatch); allowed {
		t.Errorf("Expected DENY for ID 101, got ALLOW")
	}
}
