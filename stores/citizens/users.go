package citizensstore

import (
	"context"
	"errors"

	citizenslicenses "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/licenses"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type ListCitizensOptions struct {
	IncludePhoneNumber             bool
	IncludeWanted                  bool
	IncludeJob                     bool
	IncludeTrafficInfractionPoints bool
	IncludeOpenFines               bool
	IncludeBloodType               bool
	IncludeMugshot                 bool
	IncludeEmail                   bool
}

type GetUserOptions struct {
	IncludePhoneNumber             bool
	IncludeWanted                  bool
	IncludeJob                     bool
	IncludeTrafficInfractionPoints bool
	IncludeOpenFines               bool
	IncludeBloodType               bool
	IncludeMugshot                 bool
	IncludeEmail                   bool
	IncludePropsUpdated            bool
	IncludeLicenses                bool
}

type UserActivityOptions struct {
	UserID int32
	Types  []usersactivity.UserActivityType
}

type ListUserActivityOptions struct {
	UserActivityOptions

	Sort   *database.Sort
	Offset int64
	Limit  int64
}

type CountUserActivityOptions struct {
	UserActivityOptions
}

func (s *Store) ListCitizens(
	ctx context.Context,
	req *pbcitizens.ListCitizensRequest,
	opts ListCitizensOptions,
) (*pbcitizens.ListCitizensResponse, error) {
	tUser := table.FivenetUser.AS("user")
	tUserProps := table.FivenetUserProps
	tFiles := table.FivenetFiles.AS("mugshot")

	selectors := dbutils.Columns{
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
		tUser.Dateofbirth,
		tUser.Sex,
		tUser.Height,
		tUserProps.UserID,
		s.customDB.Columns.User.GetVisum(tUser.Alias()),
	}
	condition := s.customDB.Conditions.User.GetFilter(tUser.Alias())
	orderBys := []mysql.OrderByClause{}

	if opts.IncludePhoneNumber {
		selectors = append(selectors, tUser.PhoneNumber)
		if req.GetPhoneNumber() != "" {
			phoneNumber := dbutils.PrepareForLikeSearch(req.GetPhoneNumber())
			condition = condition.AND(tUser.PhoneNumber.LIKE(mysql.String(phoneNumber)))
		}
	}
	if opts.IncludeWanted {
		selectors = append(selectors, tUserProps.Wanted)
		if req.Wanted != nil && req.GetWanted() {
			condition = condition.AND(tUserProps.Wanted.IS_TRUE())
			orderBys = append(orderBys, tUserProps.UpdatedAt.DESC())
		}
	}
	if opts.IncludeJob {
		selectors = append(selectors, tUserProps.Job, tUserProps.JobGrade)
	}
	if opts.IncludeTrafficInfractionPoints {
		selectors = append(selectors, tUserProps.TrafficInfractionPoints)
		if req.TrafficInfractionPoints != nil && req.GetTrafficInfractionPoints() > 0 {
			condition = condition.AND(
				tUserProps.TrafficInfractionPoints.GT_EQ(
					mysql.Uint32(req.GetTrafficInfractionPoints()),
				),
			)
		}
	}
	if opts.IncludeOpenFines {
		selectors = append(selectors, tUserProps.OpenFines)
		if req.OpenFines != nil && req.GetOpenFines() > 0 {
			condition = condition.AND(tUserProps.OpenFines.GT_EQ(mysql.Int64(req.GetOpenFines())))
		}
	}
	if opts.IncludeBloodType {
		selectors = append(selectors, tUserProps.BloodType)
	}
	if opts.IncludeMugshot {
		selectors = append(selectors, tUserProps.MugshotFileID, tFiles.ID, tFiles.FilePath)
	}
	if opts.IncludeEmail {
		selectors = append(selectors, tUserProps.Email)
	}

	if search := dbutils.PrepareForLikeSearch(req.GetSearch()); search != "" {
		condition = condition.AND(
			mysql.CONCAT(tUser.Firstname, mysql.String(" "), tUser.Lastname).
				LIKE(mysql.String(search)),
		)
	}
	if req.GetDateofbirth() != "" {
		condition = condition.AND(
			tUser.Dateofbirth.LIKE(
				mysql.String(dbutils.PrepareForLikeSearch(req.GetDateofbirth())),
			),
		)
	}
	if req.GetMinHeight() > 0 {
		condition = condition.AND(tUser.Height.GT_EQ(mysql.Float(float64(req.GetMinHeight()))))
	}
	if req.GetMaxHeight() > 0 {
		condition = condition.AND(tUser.Height.LT_EQ(mysql.Float(float64(req.GetMaxHeight()))))
	}

	countStmt := tUser.
		SELECT(mysql.COUNT(tUser.ID).AS("data_count.total")).
		OPTIMIZER_HINTS(mysql.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(tUser.
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUser.ID),
			),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	pag, limit := req.GetPagination().GetResponse(count.Total)
	resp := &pbcitizens.ListCitizensResponse{Pagination: pag}
	if count.Total <= 0 {
		return resp, nil
	}

	if req.GetSort() != nil && len(req.GetSort().GetColumns()) > 0 {
		for _, sc := range req.GetSort().GetColumns() {
			var column mysql.Column
			switch sc.GetId() {
			case "trafficInfractionPoints":
				if opts.IncludeTrafficInfractionPoints {
					column = tUserProps.TrafficInfractionPoints
				}
			case "openFines":
				if opts.IncludeOpenFines {
					column = tUserProps.OpenFines
				}
			case "name":
				fallthrough
			default:
				column = tUser.Firstname
			}
			if column == nil {
				column = tUser.Firstname
			}

			if sc.GetDesc() {
				orderBys = append(orderBys, column.DESC(), tUser.Lastname.DESC())
			} else {
				orderBys = append(orderBys, column.ASC(), tUser.Lastname.ASC())
			}
		}
	} else {
		orderBys = append(orderBys, tUser.Firstname.ASC(), tUser.Lastname.ASC())
	}

	stmt := tUser.
		SELECT(tUser.ID, selectors.Get()...).
		OPTIMIZER_HINTS(mysql.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(tUser.
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUser.ID),
			).LEFT_JOIN(tFiles,
			tFiles.ID.EQ(tUserProps.MugshotFileID),
		),
		).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Users); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Store) GetUser(
	ctx context.Context,
	req *pbcitizens.GetUserRequest,
	opts GetUserOptions,
) (*pbcitizens.GetUserResponse, error) {
	tUser := table.FivenetUser.AS("user")
	tUserProps := table.FivenetUserProps
	tFiles := table.FivenetFiles.AS("mugshot")

	selectors := dbutils.Columns{
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
		tUser.Dateofbirth,
		tUser.Sex,
		tUser.Height,
		tUserProps.UserID,
		s.customDB.Columns.User.GetVisum(tUser.Alias()),
	}
	if opts.IncludePropsUpdated {
		selectors = append(selectors, tUserProps.UpdatedAt)
	}
	if opts.IncludePhoneNumber {
		selectors = append(selectors, tUser.PhoneNumber)
	}
	if opts.IncludeWanted {
		selectors = append(selectors, tUserProps.Wanted)
	}
	if opts.IncludeJob {
		selectors = append(selectors, tUserProps.Job, tUserProps.JobGrade)
	}
	if opts.IncludeTrafficInfractionPoints {
		selectors = append(selectors, tUserProps.TrafficInfractionPoints)
	}
	if opts.IncludeOpenFines {
		selectors = append(selectors, tUserProps.OpenFines)
	}
	if opts.IncludeBloodType {
		selectors = append(selectors, tUserProps.BloodType)
	}
	if opts.IncludeMugshot {
		selectors = append(selectors, tUserProps.MugshotFileID, tFiles.ID, tFiles.FilePath)
	}
	if opts.IncludeEmail {
		selectors = append(selectors, tUserProps.Email)
	}

	resp := &pbcitizens.GetUserResponse{User: &users.User{}}
	stmt := tUser.
		SELECT(tUser.ID, selectors.Get()...).
		FROM(tUser.
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUser.ID),
			).LEFT_JOIN(tFiles,
			tFiles.ID.EQ(tUserProps.MugshotFileID),
		),
		).
		WHERE(tUser.ID.EQ(mysql.Int32(req.GetUserId()))).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, resp.GetUser()); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if resp.GetUser() == nil || resp.GetUser().GetUserId() <= 0 {
		return nil, nil
	}

	if resp.GetUser().GetProps() == nil {
		resp.User.Props = &usersprops.UserProps{UserId: resp.GetUser().GetUserId()}
	}

	if opts.IncludeLicenses {
		licenses, err := s.GetUserLicenses(ctx, req.GetUserId())
		if err != nil {
			return nil, err
		}
		resp.User.Licenses = licenses
	}

	return resp, nil
}

func (s *Store) GetUserLicenses(
	ctx context.Context,
	userId int32,
) ([]*citizenslicenses.License, error) {
	tCitizenLicenses := table.FivenetUserLicenses
	tLicenses := table.FivenetLicenses

	stmt := tCitizenLicenses.
		SELECT(
			tLicenses.Type.AS("license.type"),
			tLicenses.Label.AS("license.label"),
		).
		FROM(
			tCitizenLicenses.
				LEFT_JOIN(tLicenses,
					tCitizenLicenses.Type.EQ(tLicenses.Type)),
		).
		WHERE(tCitizenLicenses.UserID.EQ(mysql.Int32(userId))).
		LIMIT(15)

	var licenses []*citizenslicenses.License
	if err := stmt.QueryContext(ctx, s.db, &licenses); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return licenses, nil
}
