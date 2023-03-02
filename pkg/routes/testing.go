package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/session"
	"github.com/galexrt/arpanet/query"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/hints"
)

type Testing struct {
}

func (r *Testing) Register(e *gin.Engine) {
	e.GET("/documents", r.DocumentsGET)
	e.GET("/documents/:id", r.DocumentsByIDGET)
	e.GET("/users", r.UsersGET)
}

func (r *Testing) DocumentsGET(c *gin.Context) {
	info := session.CitizenInfo{
		Identifier: "",
		Job:        "ambulance",
		JobGrade:   20,
	}

	d := query.Document
	dja := query.DocumentJobAccess
	dua := query.DocumentUserAccess
	documents, err := d.
		LeftJoin(dua, dua.DocumentID.EqCol(d.ID), dua.Identifier.Eq(info.Identifier)).
		LeftJoin(dja, dja.DocumentID.EqCol(d.ID), dja.Name.Eq(info.Job), dja.MinimumGrade.Lte(info.JobGrade)).
		Where(
			d.Where(
				d.Where(
					d.Public.Is(true)).
					Or(d.Creator.Eq(info.Identifier)),
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
			d.JobAccess.On(dja.Name.Eq(info.Job)),
			d.UserAccess.On(dua.Identifier.Eq(info.Identifier)),
		).
		Find()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSONP(http.StatusOK, documents)
}

func (r *Testing) DocumentsByIDGET(c *gin.Context) {
	info := session.CitizenInfo{
		Identifier: "",
		Job:        "ambulance",
		JobGrade:   20,
	}

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, "No document ID given!")
		return
	}
	documentID, _ := strconv.Atoi(id)

	d := query.Document
	dja := query.DocumentJobAccess
	dua := query.DocumentUserAccess
	document, err := d.Preload(d.Responses).
		LeftJoin(dua, dua.DocumentID.EqCol(d.ID), dua.Identifier.Eq(info.Identifier)).
		LeftJoin(dja, dja.DocumentID.EqCol(d.ID), dja.Name.Eq(info.Job), dja.MinimumGrade.Lte(info.JobGrade)).
		Where(
			d.Where(d.ID.Eq(uint(documentID))),
			d.Where(
				d.Where(
					d.Public.Is(true)).
					Or(d.Creator.Eq(info.Identifier)),
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
			d.JobAccess.On(dja.Name.Eq(info.Job)),
			d.UserAccess.On(dua.Identifier.Eq(info.Identifier)),
		).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSONP(http.StatusOK, document)
}

func (r *Testing) UsersGET(c *gin.Context) {
	firstname := c.Query("firstname")
	lastname := c.Query("lastname")

	offsetQuery := c.Query("offset")
	if offsetQuery == "" {
		offsetQuery = "0"
	}

	offset, _ := strconv.Atoi(offsetQuery)
	u := query.User
	users, count, err := u.Clauses(hints.UseIndex("users_firstname_lastname_IDX")).Preload(u.UserLicenses).Where(u.Firstname.Like(firstname), u.Lastname.Like(lastname)).FindByPage(offset, 25)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	_ = count

	c.JSONP(http.StatusOK, users)
}
