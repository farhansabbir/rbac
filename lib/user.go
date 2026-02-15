package lib

import (
	"time"

	"github.com/cespare/xxhash/v2"
)

type User struct {
	ID           uint64       `json:"id"`            // interface Resource Implementer
	Name         string       `json:"name"`          // interface Resource Implementer
	Description  string       `json:"description"`   // interface Resource Implementer
	ResourceType ResourceType `json:"resource_type"` // interface Resource Implementer
	CreatedAt    time.Time    `json:"created_at"`    // interface Resource Implementer
	UpdatedAt    time.Time    `json:"updated_at"`    // interface Resource Implementer
	DeletedAt    time.Time    `json:"deleted_at"`    // interface Resource Implementer
	Email        string       `json:"email"`
	Profiles     []*Profile   `json:"profiles"`
}

func (u *User) GetResourceType() ResourceType {
	return ResourceTypeUser
}

func (u *User) GetResourceID() uint64 {
	return u.ID
}

func (u *User) GetResourceName() string {
	return u.Name
}

func (u *User) GetResourceDescription() string {
	return u.Description
}

func (u *User) GetResourceCreatedAt() time.Time {
	return u.CreatedAt
}

func (u *User) GetResourceUpdatedAt() time.Time {
	return u.UpdatedAt
}

func (u *User) GetResourceDeletedAt() time.Time {
	return u.DeletedAt
}

func (u *User) IsActive() bool {
	return u.DeletedAt.IsZero()
}

func NewUser(name string, description string, email string) *User {
	u := &User{
		ID:           xxhash.Sum64String(name + description + email),
		Name:         name,
		ResourceType: ResourceTypeUser,
		Description:  description,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		DeletedAt:    time.Time{},
		Profiles:     []*Profile{},
		Email:        email,
	}
	return u
}

func (u *User) Update(name string, description string, email string) *User {
	u.Name = name
	u.Description = description
	u.Email = email
	u.UpdatedAt = time.Now()
	return u
}

func (u *User) Restore() *User {
	u.DeletedAt = time.Time{}
	return u
}

func (u *User) SoftDelete() *User {
	u.DeletedAt = time.Now()
	return u
}

func (u *User) AddProfile(profile *Profile) *User {
	u.Profiles = append(u.Profiles, profile)
	return u
}

func (u *User) RemoveProfile(profile *Profile) *User {
	for i, p := range u.Profiles {
		if p.GetResourceID() == profile.GetResourceID() {
			u.Profiles = append(u.Profiles[:i], u.Profiles[i+1:]...)
			return u
		}
	}
	return u
}
