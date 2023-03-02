package model

func (u *User) GetLicenseFromIdentifier() string {
	return "license:" + u.Identifier[6:]
}
