package api

import (
	"errors"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/query"
	"gorm.io/gorm"
)

var (
	Documents = &documentsAPI{}
)

type documentsAPI struct {
}

func (a *documentsAPI) FindDocuments(user *model.User) ([]*model.Document, error) {
	documents, err := a.prepareDocumentQuery(nil, user).
		Find()
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (a *documentsAPI) GetDocumentByID(user *model.User, documentID int) (*model.Document, error) {
	d := query.Document

	document, err := a.prepareDocumentQuery(d.Where(d.ID.Eq(uint(documentID))), user).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return document, nil
}

func (a *documentsAPI) prepareDocumentQuery(start query.IDocumentDo, user *model.User) query.IDocumentDo {
	d := query.Document
	dja := query.DocumentJobAccess
	dua := query.DocumentUserAccess

	if start == nil {
		start = d.Where()
	}
	return start.
		LeftJoin(dua, dua.DocumentID.EqCol(d.ID), dua.Identifier.Eq(user.Identifier)).
		LeftJoin(dja, dja.DocumentID.EqCol(d.ID), dja.Name.Eq(user.Job), dja.MinimumGrade.Lte(user.JobGrade)).
		Where(
			d.Where(
				d.Where(
					d.Public.Is(true)).
					Or(d.Creator.Eq(user.Identifier)),
			).
				Or(
					d.Where(
						d.Where(
							dua.Access.IsNotNull(),
							dua.Access.Neq(model.BlockedAccessRole),
						),
					).
						Or(
							dja.Where(
								dua.Access.IsNull(),
								dja.Access.IsNotNull(),
								dja.Access.Neq(model.BlockedAccessRole),
							),
						),
				),
		).
		Order(d.CreatedAt.Desc()).
		Preload(
			d.JobAccess.On(dja.Name.Eq(user.Job)),
			d.UserAccess.On(dua.Identifier.Eq(user.Identifier)),
		)
}
