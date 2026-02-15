package lib

import (
	"fmt"
	"time"
)

type RequestContext struct {
	PrincipalID         uint64         `json:"principal_id"`
	PrincipalProfiles   []Profile      `json:"principal_profile_id"`
	RequestResourceType ResourceType   `json:"request_resource_type"`
	RequestResourceID   uint64         `json:"request_resource_id"`
	RequestVerb         Verb           `json:"request_verb"`
	ContextDT           time.Time      `json:"context_dt"`
	Attributes          map[string]any `json:"attributes"`
}

func (ctx *RequestContext) String() string {
	return fmt.Sprintf("PrincipalID: %d, PrincipalProfileID: %d, TargetResourceType: %s, TargetResourceID: %d, RequestVerb: %s, Attributes: %v",
		ctx.PrincipalID, ctx.PrincipalID, ctx.RequestResourceType, ctx.RequestResourceID, ctx.RequestVerb, ctx.Attributes)
}

func NewRequestContext(principalID uint64, principalProfiles []Profile, requestResourceType ResourceType, requestResourceID uint64, requestVerb Verb, requestAttributes map[string]any) (*RequestContext, error) {
	if principalID == 0 ||
		principalProfiles == nil ||
		requestResourceType == 0 || requestResourceType == ResourceTypeNone ||
		requestResourceID == 0 ||
		(requestVerb != VerbCreate &&
			requestVerb != VerbDelete &&
			requestVerb != VerbUpdate &&
			requestVerb != VerbRead &&
			requestVerb != VerbList &&
			requestVerb != VerbExecute) {
		return nil, fmt.Errorf("invalid request context, principalID: %d, principalProfiles: %v, requestResourceType: %s, requestResourceID: %d, requestVerb: %s, requestAttributes: %v", principalID, principalProfiles, requestResourceType, requestResourceID, requestVerb, requestAttributes)
	}
	return &RequestContext{
		PrincipalID:         principalID,
		PrincipalProfiles:   principalProfiles,
		RequestResourceType: requestResourceType,
		RequestResourceID:   requestResourceID,
		RequestVerb:         requestVerb,
		Attributes:          requestAttributes,
		ContextDT:           time.Now(),
	}, nil
}
