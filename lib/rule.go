package lib

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cespare/xxhash/v2"
)

type Verb uint8

const (
	VerbRead Verb = 1 << iota
	VerbCreate
	VerbUpdate
	VerbDelete
	VerbList
	VerbExecute
)

func (v Verb) String() string {
	switch v {
	case VerbRead:
		return "read"
	case VerbCreate:
		return "create"
	case VerbUpdate:
		return "update"
	case VerbDelete:
		return "delete"
	case VerbList:
		return "list"
	case VerbExecute:
		return "execute"
	default:
		return "*"
	}
}

type Action uint8

const (
	ActionAllow Action = 1 << iota
	ActionDeny
	ActionAllowAndForwardToNextRule
)

func (a Action) String() string {
	switch a {
	case ActionAllow:
		return "allow"
	case ActionDeny:
		return "deny"
	case ActionAllowAndForwardToNextRule:
		return "allow_and_forward_to_next_rule"
	default:
		return "deny"
	}
}

type ActionOption struct {
	Action     Action `json:"action"`
	NextRuleID uint64 `json:"next_rule_id"`
}

type Rule struct {
	ruleID                 uint64
	ruleName               string
	ruleDescription        string
	ruleResourceType       ResourceType
	ruleCreatedAt          time.Time
	ruleUpdatedAt          time.Time
	ruleDeletedAt          time.Time
	ruleTargetResourceType ResourceType
	ruleTargetResourceID   string
	ruleVerb               Verb
	ruleAction             Action
	ruleForwardRuleID      uint64
}

func (r *Rule) String() string {
	return fmt.Sprintf("rule rule_id=%d rule_name=%s rule_description=%s rule_resource_type=%d rule_created_at=%s rule_updated_at=%s rule_deleted_at=%s rule_target_resource_type=%d rule_target_resource_id=%s rule_verb=%s rule_action=%s rule_forward_rule_id=%d", r.ruleID, r.ruleName, r.ruleDescription, r.ruleResourceType, r.ruleCreatedAt, r.ruleUpdatedAt, r.ruleDeletedAt, r.ruleTargetResourceType, r.ruleTargetResourceID, r.ruleVerb, r.ruleAction, r.ruleForwardRuleID)
}

func (r *Rule) JSON() string {
	js, _ := json.Marshal(r)
	return string(js)
}

func (r *Rule) MarshalJSON() ([]byte, error) {
	// We map private fields to a public-facing map or anonymous struct
	return json.Marshal(struct {
		ID                 uint64       `json:"id"`
		Name               string       `json:"name"`
		TargetResourceType ResourceType `json:"target_resource_type"`
		TargetResourceID   string       `json:"target_resource_id"`
		Verb               string       `json:"verb"`
		Action             string       `json:"action"`
		ForwardRuleID      uint64       `json:"forward_rule_id,omitempty"`
	}{
		ID:                 r.ruleID,
		Name:               r.ruleName,
		TargetResourceType: r.ruleTargetResourceType,
		TargetResourceID:   r.ruleTargetResourceID,
		Verb:               r.ruleVerb.String(),   // Good chance to use the string representation
		Action:             r.ruleAction.String(), // for better JSON readability
		ForwardRuleID:      r.ruleForwardRuleID,
	})
}

func (r *Rule) GetResourceID() uint64 {
	return r.ruleID
}

func (r *Rule) GetResourceType() ResourceType {
	return r.ruleResourceType
}

func (r *Rule) GetResourceName() string {
	return r.ruleName
}

func (r *Rule) GetResourceDescription() string {
	return r.ruleDescription
}

func (r *Rule) GetResourceCreatedAt() time.Time {
	return r.ruleCreatedAt
}

func (r *Rule) GetResourceUpdatedAt() time.Time {
	return r.ruleUpdatedAt
}

func (r *Rule) GetResourceDeletedAt() time.Time {
	return r.ruleDeletedAt
}

func (r *Rule) IsActive() bool {
	return r.ruleDeletedAt.IsZero()
}

func NewRule(name string, description string, targetResourceID string, verb Verb, action Action) *Rule {
	rule := &Rule{
		ruleID:               xxhash.Sum64String(fmt.Sprint(ResourceTypeRule) + name + description),
		ruleName:             name,
		ruleDescription:      description,
		ruleResourceType:     ResourceTypeRule,
		ruleCreatedAt:        time.Now(),
		ruleUpdatedAt:        time.Now(),
		ruleDeletedAt:        time.Time{},
		ruleTargetResourceID: targetResourceID,
		ruleVerb:             verb,
		ruleAction:           action,
	}
	return rule
}

func NewEmptyRule(name string) *Rule {
	rule := &Rule{
		ruleID:                 xxhash.Sum64String(name),
		ruleName:               name,
		ruleDescription:        "",
		ruleResourceType:       ResourceTypeRule,
		ruleCreatedAt:          time.Now(),
		ruleUpdatedAt:          time.Now(),
		ruleDeletedAt:          time.Time{},
		ruleTargetResourceID:   "",
		ruleTargetResourceType: ResourceTypeNone,
		ruleVerb:               0,
		ruleAction:             ActionDeny,
	}
	return rule
}

// func (r *Rule) String() string {
// 	return fmt.Sprintf("rule %d:%s:%s:%s:%s", r.ruleID, r.ruleTargetResourceType, r.ruleTargetResourceID, r.ruleVerb, r.ruleAction)
// }

func (r *Rule) Update(name string, description string, targetResourceID string, verb Verb, action Action) *Rule {
	r.ruleName = name
	r.ruleDescription = description
	r.ruleTargetResourceID = targetResourceID
	r.ruleVerb = verb
	r.ruleAction = action
	r.ruleUpdatedAt = time.Now()
	return r
}

func (r *Rule) UpdateName(name string) *Rule {
	r.ruleName = name
	r.ruleUpdatedAt = time.Now()
	return r
}

func (r *Rule) UpdateDescription(description string) *Rule {
	r.ruleDescription = description
	r.ruleUpdatedAt = time.Now()
	return r
}

func (r *Rule) GetVerb() Verb {
	return r.ruleVerb
}

func (r *Rule) UpdateVerb(verb Verb) *Rule {
	r.ruleVerb = verb
	r.ruleUpdatedAt = time.Now()
	return r
}

func (r *Rule) RemoveVerb(verb Verb) *Rule {
	if r.ruleVerb == verb {
		r.ruleVerb = 0
		r.ruleUpdatedAt = time.Now()
		return r
	}
	return r
}

func (r *Rule) UpdateAction(actionOption ActionOption) (*Rule, error) {
	if actionOption.Action == ActionAllowAndForwardToNextRule {
		if actionOption.NextRuleID == 0 {
			return nil, fmt.Errorf("Invalid NextRuleID for ActionAllowAndForwardToNextRule")
		}
		r.ruleForwardRuleID = actionOption.NextRuleID
	}
	r.ruleAction = actionOption.Action
	r.ruleUpdatedAt = time.Now()
	return r, nil
}

func (r *Rule) AddTargetResourceID(targetResourceType ResourceType, targetResourceID string) *Rule {
	r.ruleTargetResourceType = targetResourceType
	r.ruleTargetResourceID = targetResourceID
	r.ruleUpdatedAt = time.Now()
	return r
}

func (r *Rule) RemoveTargetResourceID(deleteID string) *Rule {
	if r.ruleTargetResourceID == deleteID {
		r.ruleTargetResourceID = ""
		r.ruleUpdatedAt = time.Now()
		return r
	}
	return r
}

func (r *Rule) SoftDelete() *Rule {
	r.ruleDeletedAt = time.Now()
	return r
}

func (r *Rule) Restore() *Rule {
	r.ruleDeletedAt = time.Time{}
	return r
}

func (r *Rule) GetRuleName() string {
	return r.ruleName
}

func (r *Rule) GetRuleDescription() string {
	return r.ruleDescription
}

func (r *Rule) GetRuleAction() Action {
	return r.ruleAction
}

func (r *Rule) GetTargetResourceType() ResourceType {
	return r.ruleTargetResourceType
}

func (r *Rule) GetRuleTargetResourceIDs() string {
	return r.ruleTargetResourceID
}

func (r *Rule) GetRuleAsDSL() string {
	// ruleid:targetresourcetype:targetresourceID:verb:action
	return fmt.Sprintf("rule %d:%s:%s:%s:%s", r.ruleID, r.ruleTargetResourceType, r.ruleTargetResourceID, r.ruleVerb, r.ruleAction)
}

func (r *Rule) IsValidRuleSyntax() (bool, error) {
	if r.ruleResourceType == ResourceTypeRule { // only valid if this is a rule resourcetype, false otherwise
		if r.ruleTargetResourceType == ResourceTypeAll {
			if r.ruleTargetResourceID != "" {
				return false, fmt.Errorf("TargetResourceID must be empty for ResourceTypeAll")
			}
		}
		if r.ruleTargetResourceID != "" {
			if r.ruleTargetResourceType == ResourceTypeNone {
				return false, fmt.Errorf("TargetResourceType cannot be ResourceTypeNone when TargetResourceID is set")
			}
		}
		if r.ruleAction == ActionAllowAndForwardToNextRule {
			if r.ruleForwardRuleID == 0 {
				return false, fmt.Errorf("ForwardRuleID must be set for ActionAllowAndForwardToNextRule")
			}
		}
		if r.ruleForwardRuleID != 0 {
			if r.ruleAction != ActionAllowAndForwardToNextRule {
				return false, fmt.Errorf("ForwardRuleID must be set for ActionAllowAndForwardToNextRule")
			}
		}

		return true, nil
	} else {
		return false, fmt.Errorf("Not a %s", ResourceTypeRule.String())
	}
}
