package core

import (
	"encoding/json"
	"fmt"
	"sync"
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
	mux              sync.RWMutex
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID           uint64       `json:"user_id"`
		Name         string       `json:"user_name"`
		Description  string       `json:"user_description"`
		ResourceType ResourceType `json:"user_resource_type"`
		CreatedAt    time.Time    `json:"user_created_at"`
		UpdatedAt    time.Time    `json:"user_updated_at"`
		DeletedAt    time.Time    `json:"user_deleted_at"`
		Email        string       `json:"user_email"`
		Profiles     []*Profile   `json:"user_profiles"`
	}{
		ID:           u.userID,
		Name:         u.userName,
		Description:  u.userDescription,
		ResourceType: u.userResourceType,
		CreatedAt:    u.userCreatedAt,
		UpdatedAt:    u.userUpdatedAt,
		DeletedAt:    u.userDeletedAt,
		Email:        u.userEmail,
		Profiles:     u.userProfiles,
	})
}

func (u *User) GetResourceType() ResourceType {
	return u.userResourceType
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
		userID:           xxhash.Sum64String(fmt.Sprint(ResourceTypeUser) + name + description + email),
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

func (u *User) GetProfiles() []Profile {
	u.mux.RLock()
	defer u.mux.RUnlock()
	userProfiles := []Profile{}
	for _, profile := range u.userProfiles {
		userProfiles = append(userProfiles, *profile)
	}
	return userProfiles
}

func (u *User) Update(name string, description string, email string) *User {
	u.mux.Lock()
	defer u.mux.Unlock()
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
	u.mux.Lock()
	defer u.mux.Unlock()
	u.userDeletedAt = time.Now()
	return u
}

func (u *User) AddProfile(profile *Profile) *User {
	u.mux.Lock()
	defer u.mux.Unlock()
	u.userProfiles = append(u.userProfiles, profile)
	return u
}

func (u *User) RemoveProfile(profile *Profile) *User {
	u.mux.Lock()
	defer u.mux.Unlock()
	for i, p := range u.userProfiles {
		if p.GetResourceID() == profile.GetResourceID() {
			u.userProfiles = append(u.userProfiles[:i], u.userProfiles[i+1:]...)
			return u
		}
	}
	return u
}

func (u *User) GetEmail() string {
	return u.userEmail
}

func (u *User) String() string {

	return fmt.Sprintf("User: %s (%s)", u.userName, u.userEmail)
}
