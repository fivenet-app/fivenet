package model

func (u *User) GetLicenseFromIdentifier() string {
	return u.Identifier[6:]
}
