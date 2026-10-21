package citizens

import (
	context "context"
	"slices"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var ZeroTrafficInfractionPoints uint32 = 0

func (s *Server) ListCitizens(
	ctx context.Context,
	req *pbcitizens.ListCitizensRequest,
) (*pbcitizens.ListCitizensResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Field Permission Check
	propertyList, err := s.perms.AttrPropertyListContext(ctx,
		userInfo,
		permscitizens.CitizensService.ListCitizens.Properties,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	wantedEnabled, err := s.store.IsSystemPropertyEnabled(ctx, "user.wanted")
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	propertyFilters := normalizeCitizenPropertyFilters(
		req,
		wantedEnabled && propertyList.Contains("user.wanted"),
	)
	if len(propertyFilters) > 0 {
		if err := requireCitizenPropertyKeys(
			userInfo,
			propertyList,
			propertyFilters,
		); err != nil {
			return nil, err
		}
	}
	resp, err := s.store.ListCitizens(ctx, req, citizensstore.ListCitizensOptions{
		IncludePhoneNumber:             propertyList.Contains("user.phone_number"),
		IncludeWanted:                  wantedEnabled && propertyList.Contains("user.wanted"),
		IncludeJob:                     propertyList.Contains("user.job"),
		IncludeTrafficInfractionPoints: propertyList.Contains("user.traffic_infraction_points"),
		IncludeOpenFines:               propertyList.Contains("user.open_fines"),
		IncludeBloodType:               propertyList.Contains("user.blood_type"),
		IncludeMugshot:                 propertyList.Contains("user.mugshot"),
		IncludeEmail:                   propertyList.Contains("user.email"),
		PropertyFilters:                propertyFilters,
		CanFilterProperties:            len(propertyFilters) > 0,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	userIDs := make([]int32, 0, len(resp.GetUsers()))
	for _, user := range resp.GetUsers() {
		if user.GetUserId() > 0 {
			userIDs = append(userIDs, user.GetUserId())
		}
	}
	usersByID, err := s.hydrator.HydrateByUserID(ctx, s.db, userInfo, userIDs)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	for i, listedUser := range resp.GetUsers() {
		if user, ok := usersByID[listedUser.GetUserId()]; ok {
			if !wantedEnabled && user.GetProps() != nil {
				user.Props.ClearWanted()
				user.Props.ClearWantedAt()
				user.Props.ClearWantedTill()
			}
			resp.Users[i] = user
		}
	}
	if propertyList.Len() > 0 && len(userIDs) > 0 {
		valuesByUserID, err := s.store.ListProperties(ctx, userIDs)
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
		catalog, err := s.store.ResolvePropertyCatalog(ctx, propertyList.GetKeys())
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
		resolveCitizenProperties(
			resp.GetUsers(),
			valuesByUserID,
			propertyList,
			catalog.Definitions,
			catalog.Groups,
		)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetUsers() {
		if resp.GetUsers()[i].GetProps() != nil && resp.Users[i].Props.JobName != nil {
			resp.Users[i].Job = resp.GetUsers()[i].GetProps().GetJobName()
			if resp.Users[i].Props.JobGradeNumber != nil {
				resp.Users[i].JobGrade = resp.GetUsers()[i].GetProps().GetJobGradeNumber()
			} else {
				resp.Users[i].JobGrade = 0
			}

			s.enricher.EnrichJobInfo(resp.GetUsers()[i])
		} else {
			jobInfoFn(resp.GetUsers()[i])
		}
	}

	return resp, nil
}

func resolveCitizenProperties(
	usersList []*users.User,
	valuesByUserID map[int32][]*properties.Property,
	propertyList *permissionsattributes.PropertyList,
	definitions map[string]*properties.PropertyDefinition,
	groupSets ...[]*properties.PropertyGroup,
) {
	groups := propertyregistry.ResolvePropertyGroups(definitions)
	if len(groupSets) > 0 {
		groups = groupSets[0]
	}
	for _, user := range usersList {
		values := valuesByUserID[user.GetUserId()]
		visibleValues := make([]*properties.Property, 0, len(values))
		for _, value := range values {
			if propertyList.Contains(value.GetKey()) {
				visibleValues = append(visibleValues, value)
			}
		}
		user.Properties = propertyregistry.ResolvePropertiesResourceWithGroups(
			visibleValues,
			definitions,
			groups,
		)
	}
}

func normalizeCitizenPropertyFilters(
	req *pbcitizens.ListCitizensRequest,
	includeWanted bool,
) []*properties.PropertyFilter {
	filters := append([]*properties.PropertyFilter(nil), req.GetPropertyFilters()...)
	if includeWanted && req.Wanted != nil && req.GetWanted() {
		value := &properties.PropertyValue{}
		value.SetBoolValue(true)
		filters = append(filters, &properties.PropertyFilter{
			Key:       "user.wanted",
			Operation: properties.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_EQUALS,
			Value:     value,
		})
	}
	if req.TrafficInfractionPoints != nil && req.GetTrafficInfractionPoints() > 0 {
		value := &properties.PropertyValue{}
		value.SetIntegerValue(int64(req.GetTrafficInfractionPoints()))
		filters = append(filters, &properties.PropertyFilter{
			Key:       "user.traffic_infraction_points",
			Operation: properties.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_GREATER_THAN_OR_EQUAL,
			Value:     value,
		})
	}
	if req.OpenFines != nil && req.GetOpenFines() > 0 {
		value := &properties.PropertyValue{}
		value.SetIntegerValue(req.GetOpenFines())
		filters = append(filters, &properties.PropertyFilter{
			Key:       "user.open_fines",
			Operation: properties.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_GREATER_THAN_OR_EQUAL,
			Value:     value,
		})
	}
	return filters
}

func requireCitizenPropertyKeys(
	userInfo *userinfo.UserInfo,
	allowed *permissionsattributes.PropertyList,
	filters []*properties.PropertyFilter,
) error {
	for _, filter := range filters {
		if !propertyregistry.IsSearchablePropertyKey(
			properties.PropertyScope_PROPERTY_SCOPE_USER,
			filter.GetKey(),
		) {
			return errorscitizens.ErrInvalidProperty
		}
		if !userInfo.GetJobAdmin() && !allowed.Contains(filter.GetKey()) {
			return errorscitizens.ErrPropsPropertiesDenied
		}
	}
	return nil
}

func (s *Server) GetUser(
	ctx context.Context,
	req *pbcitizens.GetUserRequest,
) (*pbcitizens.GetUserResponse, error) {
	logging.InjectFields(ctx, logging.Fields{citizenIDLogFieldKey, req.GetUserId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	grpc_audit.SetTargetUser(ctx, req.GetUserId(), "")

	infoOnly := req.InfoOnly != nil && req.GetInfoOnly()

	// Field Permission Check
	propertyList, err := s.perms.AttrPropertyListContext(ctx,
		userInfo,
		permscitizens.CitizensService.ListCitizens.Properties,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	wantedEnabled, err := s.store.IsSystemPropertyEnabled(ctx, "user.wanted")
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	resp, err := s.store.GetUser(ctx, req, citizensstore.GetUserOptions{
		IncludePropsUpdated:            propertyList.Len() > 0,
		IncludePhoneNumber:             propertyList.Contains("user.phone_number"),
		IncludeWanted:                  wantedEnabled && propertyList.Contains("user.wanted"),
		IncludeJob:                     propertyList.Contains("user.job"),
		IncludeTrafficInfractionPoints: propertyList.Contains("user.traffic_infraction_points"),
		IncludeOpenFines:               propertyList.Contains("user.open_fines"),
		IncludeBloodType:               propertyList.Contains("user.blood_type"),
		IncludeMugshot:                 propertyList.Contains("user.mugshot"),
		IncludeEmail:                   propertyList.Contains("user.email"),
		IncludeLicenses:                !infoOnly && propertyList.Contains("user.licenses"),
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if resp == nil || resp.GetUser() == nil || resp.GetUser().GetUserId() <= 0 {
		return nil, errorscitizens.ErrCitizenNotFound
	}

	grpc_audit.SetTargetUser(ctx, resp.GetUser().GetUserId(), resp.GetUser().GetJob())

	if slices.Contains(s.appCfg.Get().GetJobInfo().GetPublicJobs(), resp.GetUser().GetJob()) ||
		slices.Contains(s.appCfg.Get().GetJobInfo().GetHiddenJobs(), resp.GetUser().GetJob()) {
		// Make sure user has permission to see that grade
		check, err := s.checkIfUserCanAccess(
			userInfo,
			resp.GetUser().GetJob(),
			resp.GetUser().GetJobGrade(),
		)
		if err != nil {
			return nil, err
		}
		if !check {
			return nil, errorscitizens.ErrJobGradeNoPermission
		}
	}

	licenses := resp.GetUser().GetLicenses()
	hydrated, err := s.hydrator.GetByUserID(ctx, s.db, userInfo, resp.GetUser().GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if hydrated == nil {
		return nil, errorscitizens.ErrCitizenNotFound
	}
	hydrated.Licenses = licenses
	if !wantedEnabled && hydrated.GetProps() != nil {
		hydrated.Props.ClearWanted()
		hydrated.Props.ClearWantedAt()
		hydrated.Props.ClearWantedTill()
	}
	resp.User = hydrated
	values, err := s.store.GetProperties(ctx, resp.User.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	keys := append([]string(nil), propertyList.GetKeys()...)
	catalog, err := s.store.ResolvePropertyCatalog(ctx, keys)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	visibleValues := make([]*properties.Property, 0, len(values))
	for _, value := range values {
		if !propertyList.Contains(value.GetKey()) {
			continue
		}
		visibleValues = append(visibleValues, value)
	}
	resp.User.Properties = propertyregistry.ResolvePropertiesResourceWithGroups(
		visibleValues,
		catalog.Definitions,
		catalog.Groups,
	)

	// Only let user props override the job if the person isn't in a public job
	if resp.GetUser().GetProps() != nil && resp.User.Props.JobName != nil &&
		!slices.Contains(s.appCfg.Get().GetJobInfo().GetPublicJobs(), resp.GetUser().GetJob()) {
		resp.User.Job = resp.GetUser().GetProps().GetJobName()
		if resp.User.Props.JobGradeNumber != nil {
			resp.User.JobGrade = resp.GetUser().GetProps().GetJobGradeNumber()
		} else {
			resp.User.JobGrade = 0
		}

		s.enricher.EnrichJobInfo(resp.GetUser())
	} else {
		s.enricher.EnrichJobInfoSafe(userInfo, resp.GetUser())
	}

	if resp.GetUser().GetProps() == nil {
		resp.User.Props = &usersprops.UserProps{
			UserId: resp.GetUser().GetUserId(),
		}
	}

	if propertyList.Contains("user.labels") {
		labels, err := s.store.GetUserLabelsForUser(ctx, s.db, userInfo, req.GetUserId())
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
		resp.User.Props.Labels = labels
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)

	return resp, nil
}

func (s *Server) checkIfUserCanAccess(
	userInfo *userinfo.UserInfo,
	targetUserJob string,
	targetUserGrade int32,
) (bool, error) {
	// Skip if user is job unemployed
	unemployedJob := s.appCfg.Get().GetJobInfo().GetUnemployedJob()
	if unemployedJob.GetName() == targetUserJob {
		return true, nil
	}

	// If the user is not part of public or hidden jobs (e.g., police, medics), allow access
	if !slices.Contains(s.appCfg.Get().GetJobInfo().GetPublicJobs(), targetUserJob) &&
		!slices.Contains(s.appCfg.Get().GetJobInfo().GetHiddenJobs(), targetUserJob) {
		return true, nil
	}

	jobGrades, err := s.perms.AttrJobGradeList(
		userInfo,
		permscitizens.CitizensService.GetUser.Jobs,
	)
	if err != nil {
		return false, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if jobGrades.Len() == 0 && !userInfo.GetJobAdmin() {
		return false, errorscitizens.ErrJobGradeNoPermission
	}

	// Make sure user has permission to see that job's grade, otherwise deny access to the user
	if ok := jobGrades.HasJobGrade(targetUserJob, targetUserGrade); !ok &&
		!userInfo.GetJobAdmin() {
		return false, errorscitizens.ErrJobGradeNoPermission
	}

	return true, nil
}
