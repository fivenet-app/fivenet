package citizens

import (
	context "context"
	"errors"
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
	definitionKeys := append([]string(nil), propertyList.GetKeys()...)
	definitionKeys = append(definitionKeys, requestedCitizenFilterPropertyKeys(req)...)
	definitions, err := s.store.ResolvePropertyDefinitions(ctx, definitionKeys)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	listOptions := citizensstore.ListCitizensOptionsForProperties(propertyList, definitions)
	propertyFilters := normalizeCitizenPropertyFilters(
		req,
		listOptions.IncludeWanted,
	)
	if len(propertyFilters) > 0 {
		if err := requireCitizenPropertyKeys(
			userInfo,
			propertyList,
			definitions,
			propertyFilters,
		); err != nil {
			return nil, err
		}
	}
	if err := requireCitizenLegacyFilterKeys(userInfo, propertyList, definitions, req); err != nil {
		return nil, err
	}
	listOptions.PropertyFilters = propertyFilters
	listOptions.CanFilterProperties = len(propertyFilters) > 0
	resp, err := s.store.ListCitizens(ctx, req, listOptions)
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
			if !listOptions.IncludeWanted && user.GetProps() != nil {
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
			if propertyList.Contains(value.GetKey()) &&
				!isDerivedCitizenPropertyKey(value.GetKey()) {
				visibleValues = append(visibleValues, value)
			}
		}
		visibleValues = append(visibleValues, derivedCitizenProperties(user, propertyList)...)
		user.Properties = propertyregistry.ResolvePropertiesResourceWithGroups(
			visibleValues,
			definitions,
			groups,
		)
	}
}

func derivedCitizenProperties(
	user *users.User,
	propertyList *permissionsattributes.PropertyList,
) []*properties.Property {
	if user == nil || propertyList == nil {
		return nil
	}

	result := make([]*properties.Property, 0, 6)
	if propertyList.Contains("user.dateofbirth") && user.GetDateofbirth() != "" {
		value := &properties.PropertyValue{}
		value.SetStringValue(user.GetDateofbirth())
		result = append(result, &properties.Property{Key: "user.dateofbirth", Value: value})
	}
	if propertyList.Contains("user.sex") && user.HasSex() {
		value := &properties.PropertyValue{}
		value.SetStringValue(user.GetSex())
		result = append(result, &properties.Property{Key: "user.sex", Value: value})
	}
	if propertyList.Contains("user.height") && user.HasHeight() {
		value := &properties.PropertyValue{}
		value.SetNumberValue(float64(user.GetHeight()))
		result = append(result, &properties.Property{Key: "user.height", Value: value})
	}
	if propertyList.Contains("user.phone_number") && user.HasPhoneNumber() {
		value := &properties.PropertyValue{}
		value.SetStringValue(user.GetPhoneNumber())
		result = append(result, &properties.Property{Key: "user.phone_number", Value: value})
	}
	if propertyList.Contains("user.email") && user.GetProps().GetEmail() != "" {
		value := &properties.PropertyValue{}
		value.SetStringValue(user.GetProps().GetEmail())
		result = append(result, &properties.Property{Key: "user.email", Value: value})
	}
	if propertyList.Contains("user.visum") && user.HasVisum() {
		value := &properties.PropertyValue{}
		value.SetIntegerValue(int64(user.GetVisum()))
		result = append(result, &properties.Property{Key: "user.visum", Value: value})
	}
	return result
}

func isDerivedCitizenPropertyKey(key string) bool {
	switch key {
	case "user.dateofbirth",
		"user.sex",
		"user.height",
		"user.phone_number",
		"user.email",
		"user.visum":
		return true
	default:
		return false
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
	definitions map[string]*properties.PropertyDefinition,
	filters []*properties.PropertyFilter,
) error {
	for _, filter := range filters {
		if err := requireCitizenPropertyAccess(
			userInfo,
			allowed,
			definitions,
			filter.GetKey(),
			true,
		); err != nil {
			return err
		}
	}
	return nil
}

func requireCitizenLegacyFilterKeys(
	userInfo *userinfo.UserInfo,
	allowed *permissionsattributes.PropertyList,
	definitions map[string]*properties.PropertyDefinition,
	req *pbcitizens.ListCitizensRequest,
) error {
	keys := requestedCitizenLegacyFilterPropertyKeys(req)
	for _, key := range keys {
		if err := requireCitizenPropertyAccess(
			userInfo,
			allowed,
			definitions,
			key,
			false,
		); err != nil {
			return err
		}
	}
	return nil
}

func requireCitizenPropertyAccess(
	userInfo *userinfo.UserInfo,
	allowed *permissionsattributes.PropertyList,
	definitions map[string]*properties.PropertyDefinition,
	key string,
	requireSearchable bool,
) error {
	access := propertyregistry.NewEffectivePropertyAccess(
		properties.PropertyScope_PROPERTY_SCOPE_USER,
		allowed,
		definitions,
		userInfo.GetJobAdmin(),
	)
	if !requireSearchable {
		if access.CanAccess(key) {
			return nil
		}
		return errorscitizens.ErrPropsPropertiesDenied
	}
	if errors.Is(access.ValidateFilter(key), propertyregistry.ErrPropertyFilterNotSearchable) {
		return errorscitizens.ErrInvalidProperty
	}
	if !access.CanAccess(key) {
		return errorscitizens.ErrPropsPropertiesDenied
	}
	return nil
}

func requestedCitizenFilterPropertyKeys(req *pbcitizens.ListCitizensRequest) []string {
	keys := make([]string, 0, len(req.GetPropertyFilters())+6)
	for _, filter := range req.GetPropertyFilters() {
		keys = append(keys, filter.GetKey())
	}
	if req.GetWanted() {
		keys = append(keys, "user.wanted")
	}
	if req.GetTrafficInfractionPoints() > 0 {
		keys = append(keys, "user.traffic_infraction_points")
	}
	if req.GetOpenFines() > 0 {
		keys = append(keys, "user.open_fines")
	}
	return append(keys, requestedCitizenLegacyFilterPropertyKeys(req)...)
}

func requestedCitizenLegacyFilterPropertyKeys(req *pbcitizens.ListCitizensRequest) []string {
	keys := make([]string, 0, 3)
	if req.GetPhoneNumber() != "" {
		keys = append(keys, "user.phone_number")
	}
	if req.GetDateofbirth() != "" {
		keys = append(keys, "user.dateofbirth")
	}
	if req.GetMinHeight() > 0 || req.GetMaxHeight() > 0 {
		keys = append(keys, "user.height")
	}
	return keys
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
	definitions, err := s.store.ResolvePropertyDefinitions(ctx, propertyList.GetKeys())
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	getOptions := citizensstore.GetUserOptionsForProperties(propertyList, definitions, infoOnly)
	resp, err := s.store.GetUser(ctx, req, getOptions)
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
	if !getOptions.IncludeWanted && hydrated.GetProps() != nil {
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
		if !propertyList.Contains(value.GetKey()) || isDerivedCitizenPropertyKey(value.GetKey()) {
			continue
		}
		visibleValues = append(visibleValues, value)
	}
	visibleValues = append(visibleValues, derivedCitizenProperties(resp.User, propertyList)...)
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
