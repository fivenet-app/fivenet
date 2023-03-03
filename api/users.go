package api

import (
	"errors"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/proto/common"
	"github.com/galexrt/arpanet/query"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

const (
	DefaultPageLimit = 25
)

var (
	Users = &usersAPI{}
)

type usersAPI struct {
}

func (a *usersAPI) SearchUsersByNamePages(firstname string, lastname string, offset int64, orderBys ...*common.OrderBy) ([]*model.User, int64, error) {
	return a.SearchUsersByNamePagesWithLimit(firstname, lastname, offset, DefaultPageLimit, orderBys...)
}

func (a *usersAPI) SearchUsersByNamePagesWithLimit(firstname string, lastname string, offset int64, limit int, orderBys ...*common.OrderBy) ([]*model.User, int64, error) {
	u := query.User
	q := u.Clauses(hints.UseIndex("idx_users_firstname_lastname")).Preload(u.UserProps.RelationField)

	if firstname != "" {
		q = q.Where(u.Firstname.Like("%" + firstname + "%"))
	}
	if lastname != "" {
		q = q.Where(u.Lastname.Like("%" + lastname + "%"))
	}

	// Convert our proto abstracted `common.OrderBy` to actual gorm order by instructions
	if len(orderBys) > 0 {
		for _, orderBy := range orderBys {
			field, ok := u.GetFieldByName(orderBy.Column)
			if !ok {
				return nil, 0, errors.New("orderBy column not found")
			}

			if orderBy.Desc {
				q = q.Order(field.Desc())
			} else {
				q = q.Order(field)
			}
		}
	} else {
		q = q.Order(u.Firstname)
	}

	return q.FindByPage(int(offset), limit)
}

func (a *usersAPI) GetUserByIdentifier(identifier string) (*model.User, error) {
	u := query.User
	q := u.Preload(u.UserLicenses.RelationField, u.UserProps.RelationField)
	return q.Where(u.Identifier.Eq(identifier)).Limit(1).First()
}

func (a *usersAPI) SetUserProps(userProps *model.UserProps) error {
	ups := query.UserProps
	return ups.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(userProps)
}
