package lib

import "time"

type ResourceType uint32

const (
	ResourceTypeNone ResourceType = iota
	ResourceTypeUser
	ResourceTypeProfile
	ResourceTypeURL
	ResourceTypeOrganization
	ResourceTypeProject
	ResourceTypeRole
	ResourceTypePermission
	ResourceTypeRule
	ResourceTypeAll
)

func (resourceType ResourceType) String() string {
	switch resourceType {
	case ResourceTypeUser:
		return "User"
	case ResourceTypeProfile:
		return "Profile"
	case ResourceTypeURL:
		return "URL"
	case ResourceTypeOrganization:
		return "Organization"
	case ResourceTypeProject:
		return "Project"
	case ResourceTypeRole:
		return "Role"
	case ResourceTypePermission:
		return "Permission"
	case ResourceTypeRule:
		return "Rule"
	case ResourceTypeNone:
		return ""
	default:
		return "*"
	}
}

type Resource interface {
	GetResourceID() uint64
	GetResourceName() string
	GetResourceDescription() string
	GetResourceType() ResourceType
	GetResourceCreatedAt() time.Time
	GetResourceUpdatedAt() time.Time
	GetResourceDeletedAt() time.Time
	IsActive() bool
}
