package authclaims

import "github.com/golang-jwt/jwt/v5"

type AccountInfoClaims struct {
	jwt.RegisteredClaims

	AccID          int64    `json:"aid"`
	Username       string   `json:"usr"`
	Groups         []string `json:"grps"`
	CanBeSuperuser bool     `json:"wheel,omitempty"`
}

type UserInfoClaims struct {
	jwt.RegisteredClaims

	AccID    int64   `json:"aid"`
	UserID   int32   `json:"uid"`
	Job      *string `json:"jb,omitempty"`
	JobGrade *int32  `json:"jbg,omitempty"`

	Superuser   *bool        `json:"su,omitempty"`
	OriginalJob *UserJobInfo `json:"og,omitempty"`
}

type UserJobInfo struct {
	Job      string `json:"jb"`
	JobGrade int32  `json:"jbg"`
}

// CombinedClaims combines both AccountInfoClaims and UserInfoClaims into a single struct.
// Fields are flattened to avoid nested JSON structure issue caused by the included `jwt.RegisteredClaims`.
type CombinedClaims struct {
	jwt.RegisteredClaims

	// AccountInfoClaims fields
	AccID          int64    `json:"aid"`
	Username       string   `json:"usr"`
	Groups         []string `json:"grps"`
	CanBeSuperuser bool     `json:"wheel"`

	// UserInfoClaims fields
	UserID   int32   `json:"uid"`
	Job      *string `json:"jb,omitempty"`
	JobGrade *int32  `json:"jbg,omitempty"`

	Superuser   *bool        `json:"su,omitempty"`
	OriginalJob *UserJobInfo `json:"og,omitempty"`
}

func (c *CombinedClaims) GetAccountInfoClaims() *AccountInfoClaims {
	return &AccountInfoClaims{
		RegisteredClaims: c.RegisteredClaims,
		AccID:            c.AccID,
		Username:         c.Username,
		Groups:           c.Groups,
		CanBeSuperuser:   c.CanBeSuperuser,
	}
}

func (c *CombinedClaims) GetUserInfoClaims() *UserInfoClaims {
	return &UserInfoClaims{
		RegisteredClaims: c.RegisteredClaims,
		UserID:           c.UserID,
		Job:              c.Job,
		JobGrade:         c.JobGrade,
		Superuser:        c.Superuser,
		OriginalJob:      c.OriginalJob,
	}
}
