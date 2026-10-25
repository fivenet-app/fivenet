package citizens

import (
	context "context"
	"slices"

	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListUserActivity(
	ctx context.Context,
	req *pbcitizens.ListUserActivityRequest,
) (*pbcitizens.ListUserActivityResponse, error) {
	logging.InjectFields(ctx, logging.Fields{citizenIDLogFieldKey, req.GetUserId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbcitizens.ListUserActivityResponse{
		Activity: []*usersactivity.UserActivity{},
	}

	// User can't see their own activities, unless they have "Own" perm attribute, or are a superuser
	fields, err := permscitizens.CitizensService.ListUserActivity.FieldsTyped.Get(s.perms, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if userInfo.GetUserId() == req.GetUserId() {
		// If isn't superuser or doesn't have 'Own' activity feed access
		if !userInfo.GetJobAdmin() &&
			!fields.Contains(permscitizens.CitizensServiceListUserActivityFieldsPermValueOwn) {
			return resp, nil
		}
	}
	propertyList, err := s.perms.AttrPropertyListContext(ctx,
		userInfo,
		permscitizens.CitizensService.ListCitizens.Properties,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	visiblePropertyKeys, err := s.effectiveActivityPropertyKeys(
		ctx,
		propertyList.GetKeys(),
		userInfo.GetJobAdmin(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	canViewFiles := false
	if slices.Contains(visiblePropertyKeys, citizenFilesPropertyKey) {
		canViewFiles, err = s.canViewCitizenFiles(ctx, userInfo, req.GetUserId())
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}
	activityTypes := effectiveUserActivityTypes(req.GetTypes(), visiblePropertyKeys, canViewFiles)
	queryOpts := citizensstore.CountUserActivityOptions{
		UserActivityOptions: citizensstore.UserActivityOptions{
			UserID:              req.GetUserId(),
			Types:               activityTypes,
			PropertyKeys:        req.GetPropertyKeys(),
			VisiblePropertyKeys: visiblePropertyKeys,
			From:                req.GetFrom(),
			To:                  req.GetTo(),
		},
	}
	count, err := s.store.CountUserActivity(ctx, queryOpts)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count, 20)
	resp.Pagination = pag
	if count <= 0 {
		return resp, nil
	}

	activities, err := s.store.ListUserActivity(ctx, citizensstore.ListUserActivityOptions{
		UserActivityOptions: queryOpts.UserActivityOptions,
		Sort:                req.GetSort(),
		Offset:              req.GetPagination().GetOffset(),
		Limit:               limit,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	resp.Activity = activities

	targets := make([]citizenshydrator.BasicTarget, 0, len(resp.GetActivity())*2)
	for i := range resp.GetActivity() {
		activity := resp.GetActivity()[i]
		targets = append(targets, citizenshydrator.BasicTarget{
			UserID: activity.GetTargetUserId(),
			Set: func(user *usershort.UserShort) {
				resp.Activity[i].TargetUser = user
			},
		})
		targets = append(targets, citizenshydrator.BasicTarget{
			UserID: activity.GetSourceUserId(),
			Set: func(user *usershort.UserShort) {
				resp.Activity[i].SourceUser = user
			},
		})
	}

	if err := s.hydrator.HydrateBasicTargetsSafeFunc(userInfo)(
		ctx,
		nil,
		targets,
	); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	return resp, nil
}

func effectiveUserActivityTypes(
	requested []usersactivity.UserActivityType,
	visiblePropertyKeys []string,
	canViewFiles bool,
) []usersactivity.UserActivityType {
	allTypes := []usersactivity.UserActivityType{
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_NAME,
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_LICENSES,
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_WANTED,
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_TRAFFIC_INFRACTION_POINTS,
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_MUGSHOT,
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_LABELS,
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_JOB,
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_DOCUMENT,
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_JAIL,
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_FINE,
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_FILE,
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_PROPERTY,
	}
	if len(requested) > 0 {
		allTypes = requested
	}

	allowed := make([]usersactivity.UserActivityType, 0, len(allTypes))
	for _, activityType := range allTypes {
		if activityType == usersactivity.UserActivityType_USER_ACTIVITY_TYPE_FILE && !canViewFiles {
			continue
		}
		if propertyKey := userActivityPropertyKey(activityType); propertyKey != "" &&
			!slices.Contains(visiblePropertyKeys, propertyKey) {
			continue
		}
		allowed = append(allowed, activityType)
	}
	if len(allowed) == 0 {
		return []usersactivity.UserActivityType{
			usersactivity.UserActivityType_USER_ACTIVITY_TYPE_UNSPECIFIED,
		}
	}
	return allowed
}

func userActivityPropertyKey(activityType usersactivity.UserActivityType) string {
	switch activityType {
	case usersactivity.UserActivityType_USER_ACTIVITY_TYPE_LICENSES:
		return "user.licenses"
	case usersactivity.UserActivityType_USER_ACTIVITY_TYPE_WANTED:
		return "user.wanted"
	case usersactivity.UserActivityType_USER_ACTIVITY_TYPE_TRAFFIC_INFRACTION_POINTS:
		return "user.traffic_infraction_points"
	case usersactivity.UserActivityType_USER_ACTIVITY_TYPE_MUGSHOT:
		return "user.mugshot"
	case usersactivity.UserActivityType_USER_ACTIVITY_TYPE_FILE:
		return "user.files"
	case usersactivity.UserActivityType_USER_ACTIVITY_TYPE_LABELS:
		return "user.labels"
	case usersactivity.UserActivityType_USER_ACTIVITY_TYPE_JOB:
		return "user.job"
	case usersactivity.UserActivityType_USER_ACTIVITY_TYPE_FINE:
		return "user.open_fines"
	default:
		return ""
	}
}

func (s *Server) effectiveActivityPropertyKeys(
	ctx context.Context,
	allowedKeys []string,
	propertyAdmin bool,
) ([]string, error) {
	keys := append([]string(nil), allowedKeys...)
	if propertyAdmin {
		for _, definition := range propertyregistry.SystemPropertyDefinitionsForScope(
			properties.PropertyScope_PROPERTY_SCOPE_USER,
		) {
			keys = append(keys, definition.GetKey())
		}
		customKeys, err := s.store.ListCustomPropertyKeys(ctx)
		if err != nil {
			return nil, err
		}
		keys = append(keys, customKeys...)
	}
	slices.Sort(keys)
	keys = slices.Compact(keys)

	definitions, err := s.store.ResolvePropertyDefinitions(ctx, keys)
	if err != nil {
		return nil, err
	}
	return propertyregistry.EnabledPropertyKeys(keys, definitions), nil
}
