package lib

import (
	"fmt"
	"time"
)

type RequestContext struct {
	PrincipalID         uint64         `json:"principal_id"`
	PrincipalProfiles   []uint64       `json:"principal_profiles"`
	RequestResourceType ResourceType   `json:"request_resource_type"`
	RequestResourceID   uint64         `json:"request_resource_id"` // Changed to string for matching
	RequestVerb         Verb           `json:"request_verb"`
	ContextDT           time.Time      `json:"context_dt"`
	Attributes          map[string]any `json:"attributes"`
}

func (ctx *RequestContext) String() string {
	return fmt.Sprintf("Principal: %d, Target: %s:%d, Verb: %s",
		ctx.PrincipalID, ctx.RequestResourceType, ctx.RequestResourceID, ctx.RequestVerb)
}

func NewRequestContext(principalID uint64, resType ResourceType, resID uint64, verb Verb, attrs map[string]any) (*RequestContext, error) {
	// 1. Basic Field Validation
	if principalID == 0 || resType == ResourceTypeNone || resID == 0 {
		return nil, fmt.Errorf("missing core context fields")
	}

	// 2. Verb Validation (Simplified check)
	isValidVerb := (verb & (VerbRead | VerbCreate | VerbUpdate | VerbDelete | VerbList | VerbExecute)) != 0
	if !isValidVerb {
		return nil, fmt.Errorf("invalid request verb: %s", verb)
	}

	var principalProfiles []uint64
	profs, _ := GetUserProfilesFromUserID(principalID)
	for _, prof := range profs {
		principalProfiles = append(principalProfiles, prof.GetResourceID())
	}

	return &RequestContext{
		PrincipalID:         principalID,
		PrincipalProfiles:   principalProfiles,
		RequestResourceType: resType,
		RequestResourceID:   resID,
		RequestVerb:         verb,
		Attributes:          attrs,
		ContextDT:           time.Now(),
	}, nil
}
