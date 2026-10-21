package properties

import (
	"context"
	"testing"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/stretchr/testify/require"
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
		changes []*propertypb.PropertyChange
		allowed func(string) bool
		wantErr error
	}{
		{"set", []*propertypb.PropertyChange{{Key: "user.note", Value: stringValue}}, allowed, nil},
		{"clear", []*propertypb.PropertyChange{{Key: "user.note"}}, allowed, nil},
		{
			"duplicate",
			[]*propertypb.PropertyChange{
				{Key: "user.note", Value: stringValue},
				{Key: "user.note"},
			},
			allowed,
			ErrInvalidPropertyChange,
		},
		{
			"invalid key",
			[]*propertypb.PropertyChange{{Key: "user.", Value: stringValue}},
			allowed,
			ErrInvalidPropertyChange,
		},
		{
			"denied",
			[]*propertypb.PropertyChange{{Key: "user.note", Value: stringValue}},
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
