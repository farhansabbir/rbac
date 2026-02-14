package lib

import (
	"encoding/json"
	"fmt"
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
	// Rules        []*Rule            `json:"-"`
	RuleMap map[uint64][]*Rule `json:"rule_map"`
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

func (p *Profile) GetRuleMap() map[uint64][]*Rule {
	return p.RuleMap
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
		RuleMap:      make(map[uint64][]*Rule),
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

func (p *Profile) AddRule(rule *Rule) *Profile {
	valid, _ := rule.IsValidRuleSyntax()
	if valid {
		p.RuleMap[uint64(rule.TargetResourceType)] = append(p.RuleMap[uint64(rule.TargetResourceType)], rule)
		return p
	}
	return p
}

func (p *Profile) RemoveRule(ruleID uint64) *Profile {
	for targetresourcetype, rule := range p.RuleMap {
		fmt.Print(targetresourcetype)
		fmt.Println(rule)
	}
	return p
}

func (p *Profile) GetProfileAsJSON() string {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (p *Profile) GetProfileAsMap() map[string]interface{} {
	return map[string]interface{}{
		"id":           p.ID,
		"name":         p.Name,
		"description":  p.Description,
		"resourceType": p.ResourceType,
		"createdAt":    p.CreatedAt,
		"updatedAt":    p.UpdatedAt,
		"deletedAt":    p.DeletedAt,
		"ruleMap":      p.RuleMap,
	}
}
