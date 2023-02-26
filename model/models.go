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

type ContentType string

const (
	PlaintextContentType ContentType = "plaintext"
	MarkdownContentType  ContentType = "markdown"
)
