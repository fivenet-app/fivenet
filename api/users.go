package api

import (
	"errors"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/proto/common"
	"github.com/galexrt/arpanet/query"
	"gorm.io/hints"
)

const (
	DefaultPageLimit = 25
)

var (
	Users = &users{}
)

type users struct {
}

func (a *users) SearchUsersByNamePages(firstname string, lastname string, offset int64, orderBys ...*common.OrderBy) ([]*model.User, int64, error) {
	return a.SearchUsersByNamePagesWithLimit(firstname, lastname, offset, DefaultPageLimit, orderBys...)
}

func (a *users) SearchUsersByNamePagesWithLimit(firstname string, lastname string, offset int64, limit int, orderBys ...*common.OrderBy) ([]*model.User, int64, error) {
	u := query.User
	q := u.Clauses(hints.UseIndex("users_firstname_lastname_IDX")).
		Preload(u.UserLicenses.RelationField)

	if firstname != "" {
		q = q.Where(u.Firstname.Like("%" + firstname + "%"))
	}
	if lastname != "" {
		q = q.Where(u.Lastname.Like("%" + lastname + "%"))
	}

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

func (a *users) GetUserByIdentifier(identifier string) (*model.User, error) {
	u := query.User
	return u.Where(u.Identifier.Eq(identifier)).Limit(1).First()
}

func ConvertModelUserToCommonCharacter(user *model.User) *common.Character {
	return &common.Character{
		Identifier:  user.Identifier,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Dateofbirth: user.Dateofbirth,
		Sex:         string(user.Sex),
		Height:      user.Height,
		Job:         user.Job,
		JobGrade:    int32(user.JobGrade),
		Visum:       int64(user.Visum),
		Playtime:    int64(user.Playtime),
	}
}
