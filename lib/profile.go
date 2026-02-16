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
	profRuleMap      map[uint32][]*Rule
}

func (p *Profile) String() string {
	return fmt.Sprintf("profile profile_id=%d profile_name=%s profile_description=%s profile_resource_type=%d profile_created_at=%s profile_updated_at=%s profile_deleted_at=%s profile_rule_map=%v", p.profID, p.profName, p.profDescription, p.profResourceType, p.profCreatedAt, p.profUpdatedAt, p.profDeletedAt, p.profRuleMap)
}

func (p *Profile) JSON() string {
	js, _ := json.Marshal(p)
	return string(js)
}

func (p *Profile) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID           uint64             `json:"profile_id"`
		Name         string             `json:"profile_name"`
		Description  string             `json:"profile_description"`
		ResourceType ResourceType       `json:"profile_resource_type"`
		CreatedAt    time.Time          `json:"profile_created_at"`
		UpdatedAt    time.Time          `json:"profile_updated_at"`
		DeletedAt    time.Time          `json:"profile_deleted_at"`
		RuleMap      map[uint32][]*Rule `json:"profile_rule_map"`
	}{
		ID:           p.profID,
		Name:         p.profName,
		Description:  p.profDescription,
		ResourceType: p.profResourceType,
		RuleMap:      p.profRuleMap,
		CreatedAt:    p.profCreatedAt,
		UpdatedAt:    p.profUpdatedAt,
		DeletedAt:    p.profDeletedAt,
	})
}

func (p *Profile) UnmarshalJSON(data []byte) error {
	var profile struct {
		ID           uint64             `json:"profile_id"`
		Name         string             `json:"profile_name"`
		Description  string             `json:"profile_description"`
		ResourceType ResourceType       `json:"profile_resource_type"`
		CreatedAt    time.Time          `json:"profile_created_at"`
		UpdatedAt    time.Time          `json:"profile_updated_at"`
		DeletedAt    time.Time          `json:"profile_deleted_at"`
		RuleMap      map[uint32][]*Rule `json:"profile_rule_map"`
	}

	if err := json.Unmarshal(data, &profile); err != nil {
		return err
	}

	p.profID = profile.ID
	p.profName = profile.Name
	p.profDescription = profile.Description
	p.profResourceType = profile.ResourceType
	p.profCreatedAt = profile.CreatedAt
	p.profUpdatedAt = profile.UpdatedAt
	p.profDeletedAt = profile.DeletedAt
	p.profRuleMap = profile.RuleMap

	return nil
}

func (p *Profile) GetAssociatedRules(resourceType ResourceType) []*Rule {
	// rules := []*Rule{}
	rules := p.profRuleMap[uint32(resourceType)]
	return rules
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

func (p *Profile) GetRuleMap() map[uint32][]*Rule {
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
		profID:           xxhash.Sum64String(fmt.Sprint(ResourceTypeProfile) + name + description),
		profName:         name,
		profDescription:  description,
		profResourceType: ResourceTypeProfile,
		profCreatedAt:    time.Now(),
		profUpdatedAt:    time.Now(),
		profRuleMap:      make(map[uint32][]*Rule),
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
		p.profRuleMap[uint32(rule.GetTargetResourceType())] = append(p.profRuleMap[uint32(rule.GetTargetResourceType())], rule)
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
