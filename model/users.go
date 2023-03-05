package model

type Sex string

const (
	MaleSex   Sex = "m"
	FemaleSex Sex = "f"
)

type MoneyAccounts struct {
	BlackMoney int `json:"-"`
	Bank       int `json:"bank"`
	Cash       int `json:"-"`
}

func (u *User) GetLicenseFromIdentifier() string {
	return "license:" + u.Identifier[6:]
}
