package modelhelper

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

type LicenseType string

const (
	WeaponLicense = "weapon"

	DrivingPermitLicense    = "dmv"
	DriversLicense          = "drive"
	MotorcycleLicense       = "drive_bike"
	CommercialDriverLicense = "drive_truck"

	PilotLicense = "aircraft"
	BoatLicense  = "boat"
)

type UserActivityType string

const (
	ChangedActivityType   UserActivityType = "changed"
	CreatedActivityType   UserActivityType = "created"
	MentionedActivityType UserActivityType = "mentioned"
)
