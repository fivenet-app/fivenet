package citizenstore

import (
	context "context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/dataenricher"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/proto/resources/common/database"
	users "github.com/galexrt/arpanet/proto/resources/users"
	"github.com/galexrt/arpanet/query/arpanet/model"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	user     = table.Users.AS("user")
	ul       = table.UserLicenses
	licenses = table.Licenses

	aup = table.ArpanetUserProps
	aua = table.ArpanetUserActivity
	adr = table.ArpanetDocumentsRelations.AS("document_relation")
	ad  = table.ArpanetDocuments.AS("document")
)

type Server struct {
	CitizenStoreServiceServer

	db *sql.DB
	p  perms.Permissions
	c  *dataenricher.Enricher
}

func NewServer(db *sql.DB, p perms.Permissions, c *dataenricher.Enricher) *Server {
	return &Server{
		db: db,
		p:  p,
		c:  c,
	}
}

func (s *Server) FindUsers(ctx context.Context, req *FindUsersRequest) (*FindUsersResponse, error) {
	userId, _, _ := auth.GetUserInfoFromContext(ctx)

	selectors := jet.ProjectionList{
		user.ID,
		user.Identifier,
		user.Job,
		user.JobGrade,
		user.Firstname,
		user.Lastname,
		user.Dateofbirth,
		user.Sex,
		user.Height,
		user.PhoneNumber,
		user.Visum,
		user.Playtime,
		aup.UserID,
	}
	// Field Permission Check
	if s.p.Can(userId, CitizenStoreServicePermKey, "FindUsers", "UserProps") {
		selectors = append(selectors, aup.Wanted)
	}

	req.SearchName = strings.ReplaceAll(req.SearchName, "%", "")

	condition := jet.Bool(true)
	if req.SearchName != "" {
		condition = condition.AND(jet.BoolExp(jet.Raw("MATCH(firstname,lastname) AGAINST ($search IN NATURAL LANGUAGE MODE)", jet.RawArgs{"$search": req.SearchName})))
	}
	if req.Wanted {
		condition = condition.AND(aup.Wanted.IS_TRUE())
	}

	// Get total count of values
	countStmt := user.
		SELECT(
			jet.COUNT(user.ID).AS("datacount.totalcount"),
		).
		FROM(
			user.
				LEFT_JOIN(aup, aup.UserID.EQ(user.ID)),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	resp := &FindUsersResponse{
		Pagination: database.EmptyPaginationResponse(req.Pagination.Offset),
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := user.
		SELECT(
			selectors[0], selectors[1:]...,
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname")).
		FROM(user.
			LEFT_JOIN(aup,
				aup.UserID.EQ(user.ID),
			),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		LIMIT(database.DefaultPageLimit)

	// Convert our proto abstracted `common.OrderBy` to actual gorm order by instructions
	orderBys := []jet.OrderByClause{}
	if len(req.OrderBy) > 0 {
		for _, orderBy := range req.OrderBy {
			var column jet.Column
			switch orderBy.Column {
			case "firstname":
				column = user.Firstname
			case "job":
				column = user.Job
			}

			if orderBy.Desc {
				orderBys = append(orderBys, column.DESC())
			} else {
				orderBys = append(orderBys, column.ASC())
			}
		}

		stmt = stmt.ORDER_BY(orderBys...)
	}

	if err := stmt.QueryContext(ctx, s.db, &resp.Users); err != nil {
		return nil, err
	}

	database.PaginationHelper(resp.Pagination,
		count.TotalCount,
		req.Pagination.Offset,
		len(resp.Users))

	for i := 0; i < len(resp.Users); i++ {
		s.c.EnrichJobInfo(resp.Users[i])
	}

	return resp, nil
}

func (s *Server) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	userId := auth.GetUserIDFromContext(ctx)

	selectors := jet.ProjectionList{
		user.ID,
		user.Identifier,
		user.Job,
		user.JobGrade,
		user.Firstname,
		user.Lastname,
		user.Dateofbirth,
		user.Sex,
		user.Height,
		user.PhoneNumber,
		user.Visum,
		user.Playtime,
		aup.UserID,
	}

	// Field Permission Check
	if s.p.Can(userId, CitizenStoreServicePermKey, "FindUsers", "UserProps") {
		selectors = append(selectors, aup.Wanted)
	}

	resp := &GetUserResponse{
		User: &users.User{},
	}
	stmt := user.
		SELECT(
			selectors[0], selectors[1:]...,
		).
		FROM(
			user.
				LEFT_JOIN(aup,
					aup.UserID.EQ(user.ID),
				),
		).
		WHERE(user.ID.EQ(jet.Int32(req.UserId))).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, resp.User); err != nil {
		return nil, err
	}

	// Check if user can see licenses and fetch them separately
	if s.p.Can(userId, CitizenStoreServicePermKey, "FindUsers", "Licenses") {
		stmt := user.
			SELECT(
				ul.Type.AS("license.type"),
				licenses.Label.AS("license.label"),
			).
			FROM(
				ul.
					INNER_JOIN(user,
						ul.Owner.EQ(user.Identifier),
					).
					LEFT_JOIN(licenses,
						licenses.Type.EQ(ul.Type)),
			).
			WHERE(user.ID.EQ(jet.Int32(req.UserId))).
			LIMIT(15)

		if err := stmt.QueryContext(ctx, s.db, &resp.User.Licenses); err != nil {
			if !errors.Is(qrm.ErrNoRows, err) {
				return nil, err
			}
		}
	}

	s.c.EnrichJobInfo(resp.User)

	return resp, nil
}

func (s *Server) GetUserActivity(ctx context.Context, req *GetUserActivityRequest) (*GetUserActivityResponse, error) {
	userId, _, _ := auth.GetUserInfoFromContext(ctx)

	resp := &GetUserActivityResponse{}
	// An user can never see their own activity on their own "profile"
	if userId == req.UserId {
		return resp, nil
	}

	uTarget := user.AS("target_user")
	uSource := user.AS("source_user")
	stmt := aua.
		SELECT(
			aua.AllColumns,
			uTarget.ID,
			uTarget.Identifier,
			uTarget.Job,
			uTarget.JobGrade,
			uTarget.Firstname,
			uTarget.Lastname,
			uSource.ID,
			uSource.Identifier,
			uSource.Job,
			uSource.JobGrade,
			uSource.Firstname,
			uSource.Lastname,
		).
		FROM(
			aua.
				LEFT_JOIN(uTarget,
					uTarget.ID.EQ(aua.TargetUserID),
				).
				LEFT_JOIN(uSource,
					uSource.ID.EQ(aua.SourceUserID),
				),
		).
		WHERE(
			aua.TargetUserID.EQ(jet.Int32(req.UserId)),
		).
		LIMIT(12)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	for i := 0; i < len(resp.Activity); i++ {
		s.c.EnrichJobInfo(resp.Activity[i].SourceUser)
		s.c.EnrichJobInfo(resp.Activity[i].TargetUser)
	}

	return resp, nil
}

func (s *Server) SetUserProps(ctx context.Context, req *SetUserPropsRequest) (*SetUserPropsResponse, error) {
	userId := auth.GetUserIDFromContext(ctx)

	// Field Permission Check
	if !s.p.Can(userId, CitizenStoreServicePermKey, "SetUserProps", "Wanted") {
		return nil, status.Error(codes.PermissionDenied, "You are not allowed to set user wanted status!")
	}

	stmt := aup.
		INSERT(
			aup.AllColumns,
		).
		VALUES(
			req.Props.UserId,
			req.Props.Wanted,
		).
		ON_DUPLICATE_KEY_UPDATE(
			aup.Wanted.SET(jet.Bool(req.Props.Wanted)),
		)
	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	// Create user activity
	key := "UserProps.Wanted"
	newValue := strconv.FormatBool(req.Props.Wanted)
	oldValue := strconv.FormatBool(!req.Props.Wanted)
	s.addUserAcitvity(ctx, &model.ArpanetUserActivity{
		SourceUserID: userId,
		TargetUserID: req.Props.UserId,
		Type:         int16(users.USER_ACTIVITY_TYPE_CHANGED),
		Key:          key,
		OldValue:     &oldValue,
		NewValue:     &newValue,
	})

	return &SetUserPropsResponse{}, nil
}

func (s *Server) addUserAcitvity(ctx context.Context, activity *model.ArpanetUserActivity) error {
	stmt := aua.
		INSERT(
			aua.SourceUserID,
			aua.TargetUserID,
			aua.Type,
			aua.Key,
			aua.OldValue,
			aua.NewValue,
		).
		MODEL(activity)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}

func (s *Server) GetUserDocuments(ctx context.Context, req *GetUserDocumentsRequest) (*GetUserDocumentsResponse, error) {
	resp := &GetUserDocumentsResponse{}

	uCreator := user.AS("creator")
	uSource := user.AS("source_user")
	uTarget := user.AS("target_user")
	stmt := adr.
		SELECT(
			adr.AllColumns,
			ad.ID,
			ad.CreatedAt,
			ad.UpdatedAt,
			ad.CategoryID,
			ad.CreatorID,
			ad.State,
			ad.Closed,
			ad.Title,
			uCreator.ID,
			uCreator.Identifier,
			uCreator.Job,
			uCreator.JobGrade,
			uCreator.Firstname,
			uCreator.Lastname,
			uSource.ID,
			uSource.Identifier,
			uSource.Job,
			uSource.JobGrade,
			uSource.Firstname,
			uSource.Lastname,
			uTarget.ID,
			uTarget.Identifier,
			uTarget.Job,
			uTarget.JobGrade,
			uTarget.Firstname,
			uTarget.Lastname,
		).
		FROM(
			adr.
				LEFT_JOIN(ad,
					adr.DocumentID.EQ(ad.ID),
				).
				LEFT_JOIN(uCreator,
					ad.CreatorID.EQ(uCreator.ID),
				).
				LEFT_JOIN(uSource,
					uSource.ID.EQ(adr.SourceUserID),
				).
				LEFT_JOIN(uTarget,
					uTarget.ID.EQ(adr.TargetUserID),
				),
		).
		WHERE(
			adr.TargetUserID.EQ(jet.Int32(req.UserId)),
		)

	if err := stmt.QueryContext(ctx, s.db, &resp.Relations); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	for i := 0; i < len(resp.Relations); i++ {
		s.c.EnrichJobInfo(resp.Relations[i].SourceUser)
		s.c.EnrichJobInfo(resp.Relations[i].TargetUser)
	}

	return resp, nil
}
