package lib

import (
	"time"

	"github.com/cespare/xxhash/v2"
)

type User struct {
	userID           uint64
	userName         string
	userDescription  string
	userResourceType ResourceType
	userCreatedAt    time.Time
	userUpdatedAt    time.Time
	userDeletedAt    time.Time
	userEmail        string
	userProfiles     []*Profile
}

func (u *User) GetResourceType() ResourceType {
	return ResourceTypeUser
}

func (u *User) GetResourceID() uint64 {
	return u.userID
}

func (u *User) GetResourceName() string {
	return u.userName
}

func (u *User) GetResourceDescription() string {
	return u.userDescription
}

func (u *User) GetResourceCreatedAt() time.Time {
	return u.userCreatedAt
}

func (u *User) GetResourceUpdatedAt() time.Time {
	return u.userUpdatedAt
}

func (u *User) GetResourceDeletedAt() time.Time {
	return u.userDeletedAt
}

func (u *User) IsActive() bool {
	return u.userDeletedAt.IsZero()
}

func NewUser(name string, description string, email string) *User {
	u := &User{
		userID:           xxhash.Sum64String(name + description + email),
		userName:         name,
		userResourceType: ResourceTypeUser,
		userDescription:  description,
		userCreatedAt:    time.Now(),
		userUpdatedAt:    time.Now(),
		userDeletedAt:    time.Time{},
		userProfiles:     []*Profile{},
		userEmail:        email,
	}
	return u
}

func (u *User) Update(name string, description string, email string) *User {
	u.userName = name
	u.userDescription = description
	u.userEmail = email
	u.userUpdatedAt = time.Now()
	return u
}

func (u *User) Restore() *User {
	u.userDeletedAt = time.Time{}
	return u
}

func (u *User) SoftDelete() *User {
	u.userDeletedAt = time.Now()
	return u
}

func (u *User) AddProfile(profile *Profile) *User {
	u.userProfiles = append(u.userProfiles, profile)
	return u
}

func (u *User) RemoveProfile(profile *Profile) *User {
	for i, p := range u.userProfiles {
		if p.GetResourceID() == profile.GetResourceID() {
			u.userProfiles = append(u.userProfiles[:i], u.userProfiles[i+1:]...)
			return u
		}
	}
	return u
}
