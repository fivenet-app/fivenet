package properties

import (
	"context"
	"testing"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestValidateCustomPropertyChanges(t *testing.T) {
	t.Parallel()

	stringValue := &propertypb.PropertyValue{}
	stringValue.SetStringValue("watchlist")
	definition := &propertypb.PropertyDefinition{
		Key:       "user.note",
		Enabled:   true,
		Source:    propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM,
		ValueType: propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_STRING,
	}
	resolver := func(_ context.Context, keys []string) (map[string]*propertypb.PropertyDefinition, error) {
		return map[string]*propertypb.PropertyDefinition{"user.note": definition}, nil
	}
	allowed := func(key string) bool { return key == "user.note" }

	tests := []struct {
		name    string
		changes []*ScalarChange
		allowed func(string) bool
		wantErr error
	}{
		{"set", []*ScalarChange{{Key: "user.note", Value: stringValue}}, allowed, nil},
		{"clear", []*ScalarChange{{Key: "user.note"}}, allowed, nil},
		{
			"duplicate",
			[]*ScalarChange{
				{Key: "user.note", Value: stringValue},
				{Key: "user.note"},
			},
			allowed,
			ErrInvalidPropertyChange,
		},
		{
			"invalid key",
			[]*ScalarChange{{Key: "user.", Value: stringValue}},
			allowed,
			ErrInvalidPropertyChange,
		},
		{
			"denied",
			[]*ScalarChange{{Key: "user.note", Value: stringValue}},
			func(string) bool { return false },
			ErrPropertyChangeDenied,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateCustomPropertyChanges(
				t.Context(),
				propertypb.PropertyScope_PROPERTY_SCOPE_USER,
				test.changes,
				test.allowed,
				resolver,
			)
			if test.wantErr == nil {
				require.NoError(t, err)
				return
			}
			require.ErrorIs(t, err, test.wantErr)
		})
	}
}

func TestValidateCustomPropertyChangesAppliesDefinitionConstraints(t *testing.T) {
	t.Parallel()
	minimum := uint32(4)
	definition := &propertypb.PropertyDefinition{
		Scope:            propertypb.PropertyScope_PROPERTY_SCOPE_USER,
		Key:              "user.callsign",
		Source:           propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM,
		ValueType:        propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_STRING,
		Enabled:          true,
		ValueConstraints: &propertypb.PropertyValueConstraints{StringMinLength: &minimum},
	}
	value := &propertypb.PropertyValue{}
	value.SetStringValue("abc")
	err := ValidateCustomPropertyChanges(
		t.Context(),
		propertypb.PropertyScope_PROPERTY_SCOPE_USER,
		[]*ScalarChange{{Key: definition.GetKey(), Value: value}},
		func(string) bool { return true },
		func(context.Context, []string) (map[string]*propertypb.PropertyDefinition, error) {
			return map[string]*propertypb.PropertyDefinition{definition.GetKey(): definition}, nil
		},
	)
	require.ErrorIs(t, err, ErrInvalidPropertyChange)
}

func TestValidateCustomPropertyChangesHonorsAllowedValuesMode(t *testing.T) {
	t.Parallel()
	definition := &propertypb.PropertyDefinition{
		Scope:     propertypb.PropertyScope_PROPERTY_SCOPE_USER,
		Key:       "user.callsign",
		Source:    propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM,
		ValueType: propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_STRING,
		Enabled:   true,
		ValueConstraints: &propertypb.PropertyValueConstraints{
			StringAllowedValues:     []string{"Alpha"},
			StringAllowedValuesMode: propertypb.PropertyStringAllowedValuesMode_PROPERTY_STRING_ALLOWED_VALUES_MODE_SUGGEST,
		},
	}
	resolver := func(context.Context, []string) (map[string]*propertypb.PropertyDefinition, error) {
		return map[string]*propertypb.PropertyDefinition{definition.GetKey(): definition}, nil
	}
	change := func(value string) []*ScalarChange {
		propertyValue := &propertypb.PropertyValue{}
		propertyValue.SetStringValue(value)
		return []*ScalarChange{{Key: definition.GetKey(), Value: propertyValue}}
	}

	require.NoError(
		t,
		ValidateCustomPropertyChanges(
			t.Context(),
			definition.GetScope(),
			change("Bravo"),
			func(string) bool { return true },
			resolver,
		),
	)

	definition.ValueConstraints.StringAllowedValuesMode = propertypb.PropertyStringAllowedValuesMode_PROPERTY_STRING_ALLOWED_VALUES_MODE_STRICT
	require.ErrorIs(
		t,
		ValidateCustomPropertyChanges(
			t.Context(),
			definition.GetScope(),
			change("Bravo"),
			func(string) bool { return true },
			resolver,
		),
		ErrInvalidPropertyChange,
	)

	definition.ValueConstraints.StringAllowedValuesMode = propertypb.PropertyStringAllowedValuesMode(
		99,
	)
	require.ErrorIs(
		t,
		ValidateCustomPropertyChanges(
			t.Context(),
			definition.GetScope(),
			change("Alpha"),
			func(string) bool { return true },
			resolver,
		),
		ErrInvalidPropertyChange,
	)
}

func TestValidateMutationStateAndShape(t *testing.T) {
	t.Parallel()

	policy := &propertypb.PropertyDefinition{
		Enabled: true,
		StatePolicy: &propertypb.PropertyStatePolicy{
			Enabled: true, SupportsExpiry: true,
			ExpiryAction: propertypb.PropertyExpiryAction_PROPERTY_EXPIRY_ACTION_CLEAR_VALUE,
		},
	}
	state := &propertypb.PropertyStatePatch{}
	state.SetExpiresAt(timestamp.Now())
	valid := &propertypb.PropertyMutation{Key: "user.note", State: state}
	require.NoError(t, ValidateMutationState(policy, valid))
	require.NoError(t, ValidateMutationShape(valid))

	disabled := &propertypb.PropertyDefinition{Enabled: false, StatePolicy: policy.GetStatePolicy()}
	require.ErrorIs(t, ValidateMutationState(disabled, valid), ErrInvalidPropertyChange)

	reason := "ignored"
	clearPropMutation := &propertypb.PropertyMutation{
		Key:    "user.note",
		Reason: &reason,
	}
	clearPropMutation.SetClearValue(&emptypb.Empty{})
	require.NoError(t, ValidateMutationShape(clearPropMutation))

	invalidState := &propertypb.PropertyStatePatch{}
	invalidState.SetExpiresAt(&timestamp.Timestamp{})
	invalidTimestamp := &propertypb.PropertyMutation{Key: "user.note", State: invalidState}
	require.ErrorIs(t, ValidateMutationState(policy, invalidTimestamp), ErrInvalidPropertyChange)
}

func TestValidateMutationReason(t *testing.T) {
	t.Parallel()

	reason := "reviewed"
	mutation := func(reason *string, clearReason *bool) *propertypb.PropertyMutation {
		return &propertypb.PropertyMutation{
			Key:         "user.note",
			Reason:      reason,
			ClearReason: clearReason,
		}
	}
	tests := []struct {
		name     string
		mode     propertypb.PropertyReasonMode
		mutation *propertypb.PropertyMutation
		wantErr  bool
	}{
		{
			"required rejects omitted",
			propertypb.PropertyReasonMode_PROPERTY_REASON_MODE_REQUIRED,
			mutation(nil, nil),
			true,
		},
		{
			"required accepts supplied",
			propertypb.PropertyReasonMode_PROPERTY_REASON_MODE_REQUIRED,
			mutation(&reason, nil),
			false,
		},
		{
			"required rejects clear",
			propertypb.PropertyReasonMode_PROPERTY_REASON_MODE_REQUIRED,
			mutation(nil, new(true)),
			true,
		},
		{
			"optional accepts omitted",
			propertypb.PropertyReasonMode_PROPERTY_REASON_MODE_OPTIONAL,
			mutation(nil, nil),
			false,
		},
		{
			"optional accepts supplied",
			propertypb.PropertyReasonMode_PROPERTY_REASON_MODE_OPTIONAL,
			mutation(&reason, nil),
			false,
		},
		{
			"optional accepts clear",
			propertypb.PropertyReasonMode_PROPERTY_REASON_MODE_OPTIONAL,
			mutation(nil, new(true)),
			false,
		},
		{
			"none accepts omitted",
			propertypb.PropertyReasonMode_PROPERTY_REASON_MODE_NONE,
			mutation(nil, nil),
			false,
		},
		{
			"none rejects supplied",
			propertypb.PropertyReasonMode_PROPERTY_REASON_MODE_NONE,
			mutation(&reason, nil),
			true,
		},
		{
			"none rejects clear",
			propertypb.PropertyReasonMode_PROPERTY_REASON_MODE_NONE,
			mutation(nil, new(true)),
			true,
		},
		{
			"unspecified defaults required",
			propertypb.PropertyReasonMode_PROPERTY_REASON_MODE_UNSPECIFIED,
			mutation(nil, nil),
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			definition := &propertypb.PropertyDefinition{ReasonMode: test.mode}
			err := ValidateMutationReason(definition, test.mutation)
			if test.wantErr {
				require.ErrorIs(t, err, ErrInvalidPropertyChange)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateMutationKeysRejectsDuplicates(t *testing.T) {
	t.Parallel()
	value := &propertypb.PropertyValue{}
	value.SetStringValue("advanced")
	mutation := &propertypb.PropertyMutation{Key: "user.training"}
	mutation.SetSetValue(value)
	require.ErrorIs(
		t,
		ValidateMutationKeys([]*propertypb.PropertyMutation{mutation, mutation}),
		ErrInvalidPropertyChange,
	)
}
