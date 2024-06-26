package types

import (
	"github.com/fivenet-app/fivenet/pkg/utils"
)

type Roles []*Role

type Role struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Color       *int   `yaml:"color,omitempty"`
	Permissions *int64 `yaml:"permissions,omitempty"`
	Job         string `yaml:"-"`

	Module string `yaml:"module,omitempty"`
}

type Users map[string]*User

func (u Users) Add(user *User) {
	existing, ok := u[user.ID]
	if !ok {
		u[user.ID] = user
	} else {
		existing.Merge(user)
	}
}

func (u Users) ToSlice() []*User {
	list := make([]*User, 0, len(u))

	for _, value := range u {
		list = append(list, value)
	}

	return list
}

type User struct {
	ID       string  `yaml:"userDiscordId"`
	Nickname *string `yaml:",omitempty"`
	Job      string  `yaml:"-"`

	Roles *UserRoles `yaml:"roles"`

	Kick       *bool  `yaml:"kick,omitempty"`
	KickReason string `yaml:"kickReason,omitempty"`
}

type UserRoles struct {
	Sum Roles `yaml:"-"`

	ToAdd    Roles `yaml:"toAdd,omitempty"`
	ToRemove Roles `yaml:"toRemove,omitempty"`
}

func (u *User) Merge(user *User) {
	if u.ID != user.ID {
		return
	}

	if user.Nickname != nil {
		u.Nickname = user.Nickname
	}

	if user.Kick != nil {
		u.Kick = user.Kick
	}
	if user.KickReason != "" {
		u.KickReason = user.KickReason
	}

	if u.Roles == nil {
		u.Roles = &UserRoles{}
	}
	if user.Roles != nil && len(user.Roles.Sum) > 0 {
		u.Roles.Sum = append(u.Roles.Sum, user.Roles.Sum...)
		u.Roles.Sum = utils.RemoveSliceDuplicates(u.Roles.Sum)
	}
}
