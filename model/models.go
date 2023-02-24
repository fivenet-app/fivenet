package model

type Sex string

const (
	MaleSex   Sex = "m"
	FemaleSex Sex = "f"
)

type Accounts struct {
	BlackMoney int `json:"black_money"`
	Bank       int `json:"bank"`
	Cash       int `json:"money"`
}

type ContentType string

const (
	PlaintextContentType ContentType = "plaintext"
	MarkdownContentType  ContentType = "markdown"
)
