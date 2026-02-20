package lib

import (
	"testing"

	"github.com/farhansabbir/rbac/core"
)

// Helper to reset global state between tests
func resetGlobals() {
	Users = []*core.User{}
	Profiles = []*core.Profile{}
	Rules = []*core.Rule{}
}

func TestGatekeeper_IsRequestAllowed_BasicAllow(t *testing.T) {
	resetGlobals()
	gk := NewGatekeeper()

	// 1. Setup Rule: Allow "Read" on "Profile" resources
	rule := core.NewEmptyRule("allow-read-profiles")
	rule.UpdateVerb(core.VerbRead)
	rule.SetTargetResourceTypeAndID(core.ResourceTypeProfile, core.ResourceIDAll)
	rule.UpdateAction(core.ActionOption{Action: core.ActionAllow})

	// 2. Setup Profile & User
	profile := core.NewProfile("basic-profile", "test profile")
	profile.AddRule(rule)

	user := core.NewUser("John", "User", "john@example.com")
	user.AddProfile(profile)

	// Register in globals (simulating DB)
	Users = append(Users, user)

	// 3. Create Request: Can John Read a Profile?
	ctx, err := NewRequestContext(user.GetResourceID(), core.ResourceTypeProfile, 12345, core.VerbRead, nil)
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
	allowRule := core.NewEmptyRule("allow-read")
	allowRule.UpdateVerb(core.VerbRead)
	allowRule.SetTargetResourceTypeAndID(core.ResourceTypeProfile, core.ResourceIDAll)
	allowRule.UpdateAction(core.ActionOption{Action: core.ActionAllow})

	// Rule 2: Deny Read (The "Strict" Rule)
	denyRule := core.NewEmptyRule("deny-read")
	denyRule.UpdateVerb(core.VerbRead)
	denyRule.SetTargetResourceTypeAndID(core.ResourceTypeProfile, core.ResourceIDAll)
	denyRule.UpdateAction(core.ActionOption{Action: core.ActionDeny})

	// Profile has BOTH rules
	profile := core.NewProfile("mixed-profile", "mixed")
	profile.AddRule(allowRule)
	profile.AddRule(denyRule)

	user := core.NewUser("Jane", "User", "jane@example.com")
	user.AddProfile(profile)
	Users = append(Users, user)

	// Request
	ctx, _ := NewRequestContext(user.GetResourceID(), core.ResourceTypeProfile, 999, core.VerbRead, nil)

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
	rule := core.NewEmptyRule("allow-url-read")
	rule.UpdateVerb(core.VerbRead)
	rule.SetTargetResourceTypeAndID(core.ResourceTypeURL, core.ResourceIDAll)
	rule.UpdateAction(core.ActionOption{
		Action: core.ActionAllow,
	})

	profile := core.NewProfile("url-profile", "urls")
	profile.AddRule(rule)

	user := core.NewUser("Bob", "User", "bob@example.com")
	user.AddProfile(profile)
	Users = append(Users, user)

	// Request: Try to Read a PROFILE (Mismatch!)
	ctx, _ := NewRequestContext(user.GetResourceID(), core.ResourceTypeProfile, 123, core.VerbRead, nil)

	allowed, _ := gk.IsRequestAllowed(ctx)
	if allowed {
		t.Errorf("Expected DENY (Resource Type Mismatch), but got ALLOW")
	}
}

func TestGatekeeper_VerbBitmaskMatching(t *testing.T) {
	resetGlobals()
	gk := NewGatekeeper()

	// Rule: Allow Read OR List (Bitmask: 0001 | 0010 = 0011)
	rule := core.NewEmptyRule("read-list-rule")
	rule.UpdateVerb(core.VerbRead | core.VerbList)
	rule.SetTargetResourceTypeAndID(core.ResourceTypeProfile, core.ResourceIDAll)
	rule.UpdateAction(core.ActionOption{
		Action: core.ActionAllow,
	})

	profile := core.NewProfile("reader-profile", "reader")
	profile.AddRule(rule)

	user := core.NewUser("Alice", "User", "alice@example.com")
	user.AddProfile(profile)
	Users = append(Users, user)

	// Request 1: Ask for Read (Should Match)
	ctxRead, _ := NewRequestContext(user.GetResourceID(), core.ResourceTypeProfile, 555, core.VerbRead, nil)
	if allowed, _ := gk.IsRequestAllowed(ctxRead); !allowed {
		t.Errorf("Expected ALLOW for Read request, got DENY")
	}

	// Request 2: Ask for List (Should Match)
	ctxList, _ := NewRequestContext(user.GetResourceID(), core.ResourceTypeProfile, 555, core.VerbList, nil)
	if allowed, _ := gk.IsRequestAllowed(ctxList); !allowed {
		t.Errorf("Expected ALLOW for List request, got DENY")
	}

	// Request 3: Ask for Delete (Should NOT Match)
	ctxDelete, _ := NewRequestContext(user.GetResourceID(), core.ResourceTypeProfile, 555, core.VerbDelete, nil)
	if allowed, _ := gk.IsRequestAllowed(ctxDelete); allowed {
		t.Errorf("Expected DENY for Delete request, got ALLOW")
	}
}

func TestGatekeeper_ForwardingRule(t *testing.T) {
	resetGlobals()
	gk := NewGatekeeper()

	// Rule: Forwarding Action
	rule := core.NewEmptyRule("forwarding-rule")
	rule.UpdateVerb(core.VerbExecute)
	rule.SetTargetResourceTypeAndID(core.ResourceTypeProfile, core.ResourceIDAll)

	// Set Action to Forward with a dummy NextID
	rule.UpdateAction(core.ActionOption{
		Action:     core.ActionAllowAndForwardToNextRule,
		NextRuleID: 9999,
	})

	profile := core.NewProfile("forward-profile", "forward")
	profile.AddRule(rule)

	user := core.NewUser("Dave", "User", "dave@example.com")
	user.AddProfile(profile)
	Users = append(Users, user)

	// Request
	ctx, _ := NewRequestContext(user.GetResourceID(), core.ResourceTypeProfile, 101, core.VerbExecute, nil)

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
	rule := core.NewEmptyRule("specific-id-rule")
	rule.UpdateVerb(core.VerbRead)
	rule.SetTargetResourceTypeAndID(core.ResourceTypeProfile, "100") // String ID
	rule.UpdateAction(core.ActionOption{
		Action: core.ActionAllow,
	})

	profile := core.NewProfile("strict-profile", "strict")
	profile.AddRule(rule)

	user := core.NewUser("Eve", "User", "eve@example.com")
	user.AddProfile(profile)
	Users = append(Users, user)

	// Request 1: ID 100 (Should Match)
	ctxMatch, _ := NewRequestContext(user.GetResourceID(), core.ResourceTypeProfile, 100, core.VerbRead, nil)
	if allowed, _ := gk.IsRequestAllowed(ctxMatch); !allowed {
		t.Errorf("Expected ALLOW for ID 100, got DENY")
	}

	// Request 2: ID 101 (Should Fail)
	ctxMismatch, _ := NewRequestContext(user.GetResourceID(), core.ResourceTypeProfile, 101, core.VerbRead, nil)
	if allowed, _ := gk.IsRequestAllowed(ctxMismatch); allowed {
		t.Errorf("Expected DENY for ID 101, got ALLOW")
	}
}
