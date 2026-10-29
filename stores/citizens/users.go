package citizensstore

import (
	"context"
	"errors"

	citizenslicenses "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/licenses"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type ListCitizensOptions struct {
	UserIDs                        []int32
	IncludeDateOfBirth             bool
	IncludeSex                     bool
	IncludeHeight                  bool
	IncludePhoneNumber             bool
	IncludeVisum                   bool
	IncludeWanted                  bool
	IncludeJob                     bool
	IncludeTrafficInfractionPoints bool
	IncludeOpenFines               bool
	IncludeBloodType               bool
	IncludeMugshot                 bool
	IncludeEmail                   bool
	PropertyFilters                []*properties.PropertyFilter
	CanFilterProperties            bool
}

type GetUserOptions struct {
	IncludeDateOfBirth             bool
	IncludeSex                     bool
	IncludeHeight                  bool
	IncludePhoneNumber             bool
	IncludeVisum                   bool
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

// ListCitizensOptionsForProperties derives query field selection from the
// caller's allowed properties and their effective definitions. This keeps
// direct store consumers aligned with the citizens hydrator, including
// instance-level system-property enablement overrides.
func ListCitizensOptionsForProperties(
	propertyList *permissionsattributes.PropertyList,
	definitions map[string]*properties.PropertyDefinition,
) ListCitizensOptions {
	access := propertyregistry.NewEffectivePropertyAccess(
		properties.PropertyScope_PROPERTY_SCOPE_USER,
		propertyList,
		definitions,
		false,
	)
	canUse := access.CanView
	return ListCitizensOptions{
		IncludeDateOfBirth:             canUse("user.dateofbirth"),
		IncludeSex:                     canUse("user.sex"),
		IncludeHeight:                  canUse("user.height"),
		IncludePhoneNumber:             canUse("user.phone_number"),
		IncludeVisum:                   canUse("user.visum"),
		IncludeWanted:                  canUse("user.wanted"),
		IncludeJob:                     canUse("user.job"),
		IncludeTrafficInfractionPoints: canUse("user.traffic_infraction_points"),
		IncludeOpenFines:               canUse("user.open_fines"),
		IncludeBloodType:               canUse("user.blood_type"),
		IncludeMugshot:                 canUse("user.mugshot"),
		IncludeEmail:                   canUse("user.email"),
	}
}

// GetUserOptionsForProperties derives detail-query field selection from the
// same effective property set used for citizen listings.
func GetUserOptionsForProperties(
	propertyList *permissionsattributes.PropertyList,
	definitions map[string]*properties.PropertyDefinition,
	infoOnly bool,
) GetUserOptions {
	access := propertyregistry.NewEffectivePropertyAccess(
		properties.PropertyScope_PROPERTY_SCOPE_USER,
		propertyList,
		definitions,
		false,
	)
	listOptions := ListCitizensOptionsForProperties(propertyList, definitions)
	return GetUserOptions{
		IncludeDateOfBirth:             listOptions.IncludeDateOfBirth,
		IncludeSex:                     listOptions.IncludeSex,
		IncludeHeight:                  listOptions.IncludeHeight,
		IncludePhoneNumber:             listOptions.IncludePhoneNumber,
		IncludeVisum:                   listOptions.IncludeVisum,
		IncludeWanted:                  listOptions.IncludeWanted,
		IncludeJob:                     listOptions.IncludeJob,
		IncludeTrafficInfractionPoints: listOptions.IncludeTrafficInfractionPoints,
		IncludeOpenFines:               listOptions.IncludeOpenFines,
		IncludeBloodType:               listOptions.IncludeBloodType,
		IncludeMugshot:                 listOptions.IncludeMugshot,
		IncludeEmail:                   listOptions.IncludeEmail,
		IncludePropsUpdated:            propertyList != nil && propertyList.Len() > 0,
		IncludeLicenses:                !infoOnly && access.CanView("user.licenses"),
	}
}

type UserActivityOptions struct {
	UserID              int32
	Types               []usersactivity.UserActivityType
	PropertyKeys        []string
	VisiblePropertyKeys []string
	From                *timestamp.Timestamp
	To                  *timestamp.Timestamp
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
	tWanted := table.FivenetUserPropertyValues.AS("wanted_value")
	tTrafficPoints := table.FivenetUserPropertyValues.AS("traffic_points_value")
	tOpenFines := table.FivenetUserPropertyValues.AS("open_fines_value")

	selectors := dbutils.Columns{
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
	}
	if opts.IncludeDateOfBirth {
		selectors = append(selectors, tUser.Dateofbirth)
	}
	if opts.IncludeSex {
		selectors = append(selectors, tUser.Sex)
	}
	if opts.IncludeHeight {
		selectors = append(selectors, tUser.Height)
	}
	condition := s.customDB.Conditions.User.GetFilter(tUser.Alias())
	if len(opts.UserIDs) > 0 {
		ids := make([]mysql.Expression, 0, len(opts.UserIDs))
		for _, id := range opts.UserIDs {
			ids = append(ids, mysql.Int32(id))
		}
		condition = condition.AND(tUser.ID.IN(ids...))
	}
	orderBys := []mysql.OrderByClause{}

	if opts.IncludePhoneNumber {
		selectors = append(selectors, tUser.PhoneNumber)
		if req.GetPhoneNumber() != "" {
			phoneNumber := dbutils.PrepareForLikeSearch(req.GetPhoneNumber())
			condition = condition.AND(tUser.PhoneNumber.LIKE(mysql.String(phoneNumber)))
		}
	}
	if opts.IncludeVisum {
		selectors = append(selectors, s.customDB.Columns.User.GetVisum(tUser.Alias()))
	}
	if opts.CanFilterProperties {
		for i, filter := range opts.PropertyFilters {
			propertyCondition, err := userPropertyFilterCondition(filter, tUser.Alias(), i)
			if err != nil {
				return nil, err
			}
			condition = condition.AND(propertyCondition)
			if filter.GetKey() == "user.wanted" && filter.GetValue().GetBoolValue() {
				orderBys = append(orderBys, tWanted.UpdatedAt.DESC())
			}
		}
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
		FROM(tUser).
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
					column = tTrafficPoints.ValueInteger
				}
			case "openFines":
				if opts.IncludeOpenFines {
					column = tOpenFines.ValueInteger
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
			LEFT_JOIN(tWanted, mysql.AND(
				tWanted.UserID.EQ(tUser.ID),
				tWanted.PropertyKey.EQ(mysql.String("user.wanted")),
				tWanted.ValueType.EQ(mysql.Uint8(uint8(userPropertyValueTypeBool))),
			)).
			LEFT_JOIN(tTrafficPoints, mysql.AND(
				tTrafficPoints.UserID.EQ(tUser.ID),
				tTrafficPoints.PropertyKey.EQ(mysql.String("user.traffic_infraction_points")),
				tTrafficPoints.ValueType.EQ(mysql.Uint8(uint8(userPropertyValueTypeInteger))),
			)).
			LEFT_JOIN(tOpenFines, mysql.AND(
				tOpenFines.UserID.EQ(tUser.ID),
				tOpenFines.PropertyKey.EQ(mysql.String("user.open_fines")),
				tOpenFines.ValueType.EQ(mysql.Uint8(uint8(userPropertyValueTypeInteger))),
			)),
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

	selectors := dbutils.Columns{
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
	}
	if opts.IncludeDateOfBirth {
		selectors = append(selectors, tUser.Dateofbirth)
	}
	if opts.IncludeSex {
		selectors = append(selectors, tUser.Sex)
	}
	if opts.IncludeHeight {
		selectors = append(selectors, tUser.Height)
	}
	if opts.IncludePhoneNumber {
		selectors = append(selectors, tUser.PhoneNumber)
	}
	if opts.IncludeVisum {
		selectors = append(selectors, s.customDB.Columns.User.GetVisum(tUser.Alias()))
	}

	resp := &pbcitizens.GetUserResponse{User: &users.User{}}
	stmt := tUser.
		SELECT(tUser.ID, selectors.Get()...).
		FROM(tUser).
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
