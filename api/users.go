package api

import (
	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/query"
	"gorm.io/hints"
)

const (
	DefaultPageLimit = 25
)

var Users = &users{}

type users struct {
}

func (a *users) SearchUsersByNamePages(firstname string, lastname string, offset int) ([]*model.User, int64, error) {
	return a.SearchUsersByNamePagesWithLimit(firstname, lastname, offset, DefaultPageLimit)
}

func (a *users) SearchUsersByNamePagesWithLimit(firstname string, lastname string, offset int, limit int) ([]*model.User, int64, error) {
	u := query.User
	return u.Clauses(hints.UseIndex("users_firstname_lastname_IDX")).
		Preload(u.UserLicenses).
		Where(u.Firstname.Like(firstname), u.Lastname.Like(lastname)).
		FindByPage(offset, limit)
}
