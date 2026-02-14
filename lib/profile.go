package lib

import (
	"time"

	"github.com/cespare/xxhash/v2"
)

type Profile struct {
	// aka policy
	ID           uint64       `json:"id"`            // interface Resource Implementer
	Name         string       `json:"name"`          // interface Resource Implementer
	Description  string       `json:"description"`   // interface Resource Implementer
	ResourceType ResourceType `json:"resource_type"` // interface Resource Implementer
	CreatedAt    time.Time    `json:"created_at"`    // interface Resource Implementer
	UpdatedAt    time.Time    `json:"updated_at"`    // interface Resource Implementer
	DeletedAt    time.Time    `json:"deleted_at"`    // interface Resource Implementer
	Rules        []*Rule      `json:"rules"`
}

func (p *Profile) GetResourceID() uint64 {
	return p.ID
}

func (p *Profile) GetResourceName() string {
	return p.Name
}

func (p *Profile) GetResourceDescription() string {
	return p.Description
}

func (p *Profile) GetRules() []*Rule {
	return p.Rules
}

func (p *Profile) GetResourceType() ResourceType {
	return p.ResourceType
}

func (p *Profile) GetResourceCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *Profile) GetResourceUpdatedAt() time.Time {
	return p.UpdatedAt
}

func (p *Profile) GetResourceDeletedAt() time.Time {
	return p.DeletedAt
}

func (p *Profile) IsActive() bool {
	return p.DeletedAt.IsZero()
}

func NewProfile(name string, description string) *Profile {
	return &Profile{
		ID:           xxhash.Sum64String(name + description),
		Name:         name,
		Description:  description,
		ResourceType: ResourceTypeProfile,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Rules:        []*Rule{},
	}
}

func (p *Profile) UpdateName(name string) *Profile {
	p.Name = name
	return p
}

func (p *Profile) UpdateDescription(description string) *Profile {
	p.Description = description
	return p
}

func (p *Profile) UpdateRules(rules []*Rule) *Profile {
	p.Rules = rules
	return p
}

func (p *Profile) AddRule(rule *Rule) *Profile {
	p.Rules = append(p.Rules, rule)
	return p
}

// func (p *Profile) GetRules() []*Rule {
// 	return p.Rules
// }

func (p *Profile) RemoveRule(ruleID uint64) *Profile {
	for i, rule := range p.Rules {
		if rule.ID == ruleID {
			p.Rules = append(p.Rules[:i], p.Rules[i+1:]...)
			return p
		}
	}
	return p
}
