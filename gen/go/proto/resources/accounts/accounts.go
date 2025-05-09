package accounts

import (
	timestamp "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
)

func ConvertFromModelAcc(a *model.FivenetAccounts) *Account {
	var createdAt *timestamp.Timestamp
	if a.CreatedAt != nil {
		createdAt = timestamp.New(*a.CreatedAt)
	}
	var updatedAt *timestamp.Timestamp
	if a.UpdatedAt != nil {
		updatedAt = timestamp.New(*a.UpdatedAt)
	}

	return &Account{
		Id:        a.ID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Username:  *a.Username,
		License:   a.License,
	}
}
