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
)

func (a Action) String() string {
	switch a {
	case ActionAllow:
		return "allow"
	case ActionDeny:
		return "deny"
	default:
		return "deny"
	}
}

type Rule struct {
	ID                 uint64       `json:"id"`                    // interface Resource Implementer
	Name               string       `json:"name,omitempty"`        // interface Resource Implementer
	Description        string       `json:"description,omitempty"` // interface Resource Implementer
	ResourceType       ResourceType `json:"resource_type"`         // interface Resource Implementer
	CreatedAt          time.Time    `json:"created_at"`            // interface Resource Implementer
	UpdatedAt          time.Time    `json:"updated_at"`            // interface Resource Implementer
	DeletedAt          time.Time    `json:"deleted_at"`            // interface Resource Implementer
	TargetResourceType ResourceType `json:"target_resource_type"`
	TargetResourceID   string       `json:"target_resource_id,omitempty"`
	Verb               Verb         `json:"verb"`
	Action             Action       `json:"action"`
}

func (r *Rule) GetResourceID() uint64 {
	return r.ID
}

func (r *Rule) GetResourceType() ResourceType {
	return r.ResourceType
}

func (r *Rule) GetResourceName() string {
	return r.Name
}

func (r *Rule) GetResourceDescription() string {
	return r.Description
}

func (r *Rule) GetResourceCreatedAt() time.Time {
	return r.CreatedAt
}

func (r *Rule) GetResourceUpdatedAt() time.Time {
	return r.UpdatedAt
}

func (r *Rule) GetResourceDeletedAt() time.Time {
	return r.DeletedAt
}

func (r *Rule) IsActive() bool {
	return r.DeletedAt.IsZero()
}

func NewRule(name string, description string, targetResourceID string, verb Verb, action Action) *Rule {
	rule := &Rule{
		ID:               xxhash.Sum64String(name + description),
		Name:             name,
		Description:      description,
		ResourceType:     ResourceTypeRule,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		DeletedAt:        time.Time{},
		TargetResourceID: targetResourceID,
		Verb:             verb,
		Action:           action,
	}
	return rule
}

func NewEmptyRule(name string) *Rule {
	rule := &Rule{
		ID:                 xxhash.Sum64String(name),
		Name:               name,
		Description:        "",
		ResourceType:       ResourceTypeRule,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		DeletedAt:          time.Time{},
		TargetResourceID:   "",
		TargetResourceType: ResourceTypeNone,
		Verb:               0,
		Action:             ActionDeny,
	}
	return rule
}

func (r *Rule) Update(name string, description string, targetResourceID string, verb Verb, action Action) *Rule {
	r.Name = name
	r.Description = description
	r.TargetResourceID = targetResourceID
	r.Verb = verb
	r.Action = action
	r.UpdatedAt = time.Now()
	return r
}

func (r *Rule) UpdateName(name string) *Rule {
	r.Name = name
	r.UpdatedAt = time.Now()
	return r
}

func (r *Rule) UpdateDescription(description string) *Rule {
	r.Description = description
	r.UpdatedAt = time.Now()
	return r
}

func (r *Rule) GetVerb() Verb {
	return r.Verb
}

func (r *Rule) UpdateVerb(verb Verb) *Rule {
	r.Verb = verb
	r.UpdatedAt = time.Now()
	return r
}

func (r *Rule) RemoveVerb(verb Verb) *Rule {
	if r.Verb == verb {
		r.Verb = 0
		r.UpdatedAt = time.Now()
		return r
	}
	return r
}

func (r *Rule) UpdateAction(action Action) *Rule {
	r.Action = action
	r.UpdatedAt = time.Now()
	return r
}

func (r *Rule) AddTargetResourceID(targetResourceType ResourceType, targetResourceID string) *Rule {
	r.TargetResourceType = targetResourceType
	r.TargetResourceID = targetResourceID
	r.UpdatedAt = time.Now()
	return r
}

func (r *Rule) RemoveTargetResourceID(deleteID string) *Rule {
	if r.TargetResourceID == deleteID {
		r.TargetResourceID = ""
		r.UpdatedAt = time.Now()
		return r
	}
	return r
}

func (r *Rule) SoftDelete() *Rule {
	r.DeletedAt = time.Now()
	return r
}

func (r *Rule) Restore() *Rule {
	r.DeletedAt = time.Time{}
	return r
}

func (r *Rule) GetRuleName() string {
	return r.Name
}

func (r *Rule) GetRuleDescription() string {
	return r.Description
}

func (r *Rule) GetRuleAction() Action {
	return r.Action
}

func (r *Rule) GetRuleTargetResourceIDs() string {
	return r.TargetResourceID
}

func (r *Rule) GetRuleAsJSON() string {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (r *Rule) GetRuleAsDSL() string {
	// ruleid:targetresourcetype:targetresourceID:verb:action
	return fmt.Sprintf("rule %d:%s:%s:%s:%s", r.ID, r.TargetResourceType, r.TargetResourceID, r.Verb, r.Action)
}

func (r *Rule) IsValidRuleSyntax() (bool, error) {
	if r.ResourceType == ResourceTypeRule { // only valid if this is a rule resourcetype, false otherwise
		if r.TargetResourceType == ResourceTypeAll {
			if r.TargetResourceID != "" {
				return false, fmt.Errorf("TargetResourceID must be empty for ResourceTypeAll")
			}
		}
		if r.TargetResourceID != "" {
			if r.TargetResourceType == ResourceTypeNone {
				return false, fmt.Errorf("TargetResourceType cannot be ResourceTypeNone when TargetResourceID is set")
			}
		}

		return true, nil
	} else {
		return false, fmt.Errorf("Not a %s", ResourceTypeRule.String())
	}
}
