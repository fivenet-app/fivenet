package accounts

import (
	"database/sql/driver"
	"encoding/json"
	"slices"

	timestamp "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
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
	var enabled bool
	if a.Enabled != nil {
		enabled = *a.Enabled
	}

	groups := &AccountGroups{}
	if a.Groups != nil {
		_ = groups.Scan(a.Groups)
	}

	return &Account{
		Id:        a.ID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Enabled:   enabled,
		Username:  *a.Username,
		License:   a.License,
		Groups:    groups,
		LastChar:  a.LastChar,
	}
}

func (ag *AccountGroups) ContainsAnyGroup(groups []string) bool {
	if ag == nil {
		return false
	}

	agg := ag.GetGroups()
	for _, g := range groups {
		if slices.Contains(agg, g) {
			return true
		}
	}

	return false
}

func (ag *AccountGroups) AddGroup(group string) {
	if ag == nil {
		ag = &AccountGroups{
			Groups: []string{},
		}
	}

	if ag.ContainsAnyGroup([]string{group}) {
		return
	}

	ag.Groups = append(ag.Groups, group)
}

func (ag *AccountGroups) Equal(in *AccountGroups) bool {
	if ag == nil || in == nil {
		return false
	}

	aGroups := ag.GetGroups()
	bGroups := in.GetGroups()

	if len(aGroups) != len(bGroups) {
		return false
	}

	aGroupSet := make(map[string]struct{}, len(aGroups))
	for _, g := range aGroups {
		aGroupSet[g] = struct{}{}
	}

	for _, g := range bGroups {
		if _, exists := aGroupSet[g]; !exists {
			return false
		}
	}

	return true
}

// Scan implements driver.Valuer for protobuf AccountGroups.
func (x *AccountGroups) Scan(value any) error {
	switch t := value.(type) {
	case string:
		var dest []string
		if err := json.Unmarshal([]byte(t), &dest); err != nil {
			return err
		}
		for _, group := range dest {
			x.Groups = append(x.Groups, group)
		}
		return nil
	case *string:
		if t == nil {
			return nil
		}
		var dest []string
		if err := json.Unmarshal([]byte(*t), &dest); err != nil {
			return err
		}
		for _, group := range dest {
			x.Groups = append(x.Groups, group)
		}
		return nil
	case []byte:
		var dest []string
		if err := json.Unmarshal(t, &dest); err != nil {
			return err
		}
		for _, group := range dest {
			x.Groups = append(x.Groups, group)
		}
		return nil
	}
	return nil
}

// Value marshals the AccountGroups value into driver.Valuer.
func (x *AccountGroups) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := json.Marshal(x.GetGroups())
	if err != nil {
		return nil, err
	}
	return string(out), err
}
