package lib

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cespare/xxhash/v2"
)

type Profile struct {
	// aka policy
	profID           uint64
	profName         string
	profDescription  string
	profResourceType ResourceType
	profCreatedAt    time.Time
	profUpdatedAt    time.Time
	profDeletedAt    time.Time
	profRuleMap      map[uint64][]*Rule
}

func (p *Profile) GetResourceID() uint64 {
	return p.profID
}

func (p *Profile) GetResourceName() string {
	return p.profName
}

func (p *Profile) GetResourceDescription() string {
	return p.profDescription
}

func (p *Profile) GetRuleMap() map[uint64][]*Rule {
	return p.profRuleMap
}

func (p *Profile) GetResourceType() ResourceType {
	return p.profResourceType
}

func (p *Profile) GetResourceCreatedAt() time.Time {
	return p.profCreatedAt
}

func (p *Profile) GetResourceUpdatedAt() time.Time {
	return p.profUpdatedAt
}

func (p *Profile) GetResourceDeletedAt() time.Time {
	return p.profDeletedAt
}

func (p *Profile) IsActive() bool {
	return p.profDeletedAt.IsZero()
}

func NewProfile(name string, description string) *Profile {
	return &Profile{
		profID:           xxhash.Sum64String(name + description),
		profName:         name,
		profDescription:  description,
		profResourceType: ResourceTypeProfile,
		profCreatedAt:    time.Now(),
		profUpdatedAt:    time.Now(),
		profRuleMap:      make(map[uint64][]*Rule),
	}
}

func (p *Profile) UpdateName(name string) *Profile {
	p.profName = name
	return p
}

func (p *Profile) UpdateDescription(description string) *Profile {
	p.profDescription = description
	return p
}

func (p *Profile) AddRule(rule *Rule) *Profile {
	valid, _ := rule.IsValidRuleSyntax()
	if valid {
		p.profRuleMap[uint64(rule.GetTargetResourceType())] = append(p.profRuleMap[uint64(rule.GetTargetResourceType())], rule)
		return p
	}
	return p
}

func (p *Profile) RemoveRule(ruleID uint64) *Profile {
	for targetresourcetype, rule := range p.profRuleMap {
		fmt.Print(targetresourcetype)
		fmt.Println(rule)
	}
	return p
}

func (p *Profile) GetProfileAsJSON() string {
	jsonBytes, err := json.Marshal(p.GetProfileAsMap())
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (p *Profile) GetProfileAsMap() map[string]any {
	return map[string]any{
		"id":           p.profID,
		"name":         p.profName,
		"description":  p.profDescription,
		"resourceType": p.profResourceType,
		"createdAt":    p.profCreatedAt,
		"updatedAt":    p.profUpdatedAt,
		"deletedAt":    p.profDeletedAt,
		"ruleMap":      p.profRuleMap,
	}
}
