package lib

import "time"

type ResourceType uint8

const (
	ResourceTypeUser ResourceType = iota
	ResourceTypeProfile
	ResourceTypeURL
	ResourceTypeOrganization
	ResourceTypeProject
	ResourceTypeRole
	ResourceTypePermission
	ResourceTypeRule
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
	default:
		return "Unknown"
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
