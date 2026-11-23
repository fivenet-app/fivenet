package properties

import (
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
)

// ResolveUserProperties attaches the visible generic and derived system
// properties to users. Callers must pass definitions already resolved for the
// requesting user's property grants.
func ResolveUserProperties(
	usersList []*users.User,
	valuesByUserID map[int32][]*propertypb.Property,
	propertyList *permissionsattributes.PropertyList,
	definitions map[string]*propertypb.PropertyDefinition,
	groups []*propertypb.PropertyGroup,
) {
	if propertyList == nil {
		return
	}

	for _, user := range usersList {
		if user == nil {
			continue
		}

		values := valuesByUserID[user.GetUserId()]
		visibleValues := make([]*propertypb.Property, 0, len(values)+6)
		for _, value := range values {
			if value != nil && propertyList.Contains(value.GetKey()) &&
				!isDerivedUserPropertyKey(value.GetKey()) {
				visibleValues = append(visibleValues, value)
			}
		}
		visibleValues = append(visibleValues, derivedUserProperties(user, propertyList)...)
		user.Properties = ResolvePropertiesResourceWithGroups(visibleValues, definitions, groups)
	}
}

func derivedUserProperties(
	user *users.User,
	propertyList *permissionsattributes.PropertyList,
) []*propertypb.Property {
	if user == nil || propertyList == nil {
		return nil
	}

	result := make([]*propertypb.Property, 0, 6)
	if propertyList.Contains(UserPropertyKeyDateOfBirth) && user.GetDateofbirth() != "" {
		value := &propertypb.PropertyValue{}
		value.SetStringValue(user.GetDateofbirth())
		result = append(result, &propertypb.Property{Key: UserPropertyKeyDateOfBirth, Value: value})
	}
	if propertyList.Contains(UserPropertyKeySex) && user.HasSex() {
		value := &propertypb.PropertyValue{}
		value.SetStringValue(user.GetSex())
		result = append(result, &propertypb.Property{Key: UserPropertyKeySex, Value: value})
	}
	if propertyList.Contains(UserPropertyKeyHeight) && user.HasHeight() {
		value := &propertypb.PropertyValue{}
		value.SetNumberValue(float64(user.GetHeight()))
		result = append(result, &propertypb.Property{Key: UserPropertyKeyHeight, Value: value})
	}
	if propertyList.Contains(UserPropertyKeyPhoneNumber) && user.HasPhoneNumber() {
		value := &propertypb.PropertyValue{}
		value.SetStringValue(user.GetPhoneNumber())
		result = append(result, &propertypb.Property{Key: UserPropertyKeyPhoneNumber, Value: value})
	}
	if propertyList.Contains(UserPropertyKeyEmail) && user.GetProps().GetEmail() != "" {
		value := &propertypb.PropertyValue{}
		value.SetStringValue(user.GetProps().GetEmail())
		result = append(result, &propertypb.Property{Key: UserPropertyKeyEmail, Value: value})
	}
	if propertyList.Contains(UserPropertyKeyVisum) && user.HasVisum() {
		value := &propertypb.PropertyValue{}
		value.SetIntegerValue(int64(user.GetVisum()))
		result = append(result, &propertypb.Property{Key: UserPropertyKeyVisum, Value: value})
	}
	return result
}

func isDerivedUserPropertyKey(key string) bool {
	switch key {
	case UserPropertyKeyDateOfBirth,
		UserPropertyKeySex,
		UserPropertyKeyHeight,
		UserPropertyKeyPhoneNumber,
		UserPropertyKeyEmail,
		UserPropertyKeyVisum:
		return true
	default:
		return false
	}
}
