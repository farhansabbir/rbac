package lib

import (
	"time"

	"github.com/cespare/xxhash/v2"
)

type Verbs uint8

const (
	Read Verbs = iota
	Create
	Update
	Delete
	List
	Execute
)

func (v Verbs) String() string {
	switch v {
	case Read:
		return "read"
	case Create:
		return "create"
	case Update:
		return "update"
	case Delete:
		return "delete"
	case List:
		return "list"
	case Execute:
		return "execute"
	default:
		return "unknown"
	}
}

type Action uint8

const (
	Allow Action = iota
	Deny
)

func (a Action) String() string {
	switch a {
	case Allow:
		return "allow"
	case Deny:
		return "deny"
	default:
		return "unknown"
	}
}

type Rule struct {
	ID                uint64       `json:"id"`
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	TargetResourceIDs []string     `json:"target_resource_ids"`
	Verbs             []Verbs      `json:"verbs"`
	Action            Action       `json:"action"`
	ResourceType      ResourceType `json:"resource_type"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
	DeletedAt         time.Time    `json:"deleted_at"`
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

func NewRule(name string, description string, subjectResourceIDs []string, verbs []Verbs, action Action) *Rule {
	rule := &Rule{
		ID:                xxhash.Sum64String(name + description),
		Name:              name,
		Description:       description,
		TargetResourceIDs: subjectResourceIDs,
		Verbs:             verbs,
		Action:            action,
		ResourceType:      ResourceTypeRule,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		DeletedAt:         time.Time{},
	}
	return rule
}

func (r *Rule) Update(name string, description string, subjectResourceIDs []string, verbs []Verbs, action Action) *Rule {
	r.Name = name
	r.Description = description
	r.TargetResourceIDs = subjectResourceIDs
	r.Verbs = verbs
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

func (r *Rule) UpdateVerbs(verbs []Verbs) *Rule {
	r.Verbs = verbs
	r.UpdatedAt = time.Now()
	return r
}

func (r *Rule) RemoveVerb(verb Verbs) *Rule {
	for i, v := range r.Verbs {
		if v == verb {
			r.Verbs = append(r.Verbs[:i], r.Verbs[i+1:]...)
			r.UpdatedAt = time.Now()
			return r
		}
	}
	return r
}

func (r *Rule) RemoveVerbs(verbs []Verbs) *Rule {
	for _, verb := range verbs {
		r.RemoveVerb(verb)
	}
	return r
}

func (r *Rule) UpdateAction(action Action) *Rule {
	r.Action = action
	r.UpdatedAt = time.Now()
	return r
}

func (r *Rule) UpdateTargetResourceIDs(subjectResourceIDs []string) *Rule {
	r.TargetResourceIDs = subjectResourceIDs
	r.UpdatedAt = time.Now()
	return r
}

func (r *Rule) RemoveTargetResourceID(deleteID string) *Rule {
	for i, id := range r.TargetResourceIDs {
		if id == deleteID {
			r.TargetResourceIDs = append(r.TargetResourceIDs[:i], r.TargetResourceIDs[i+1:]...)
			r.UpdatedAt = time.Now()
			return r
		}
	}
	return r
}

func (r *Rule) RemoveTargetResourceIDs(deleteIDs []string) *Rule {
	for _, id := range deleteIDs {
		r.RemoveTargetResourceID(id)
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
