package lib

import "time"

type Profile struct {
	// aka policy
	ID           uint64       `json:"id"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Rules        []*Rule      `json:"rules"`
	ResourceType ResourceType `json:"resource_type"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    time.Time    `json:"deleted_at"`
}
