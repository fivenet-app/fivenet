package properties

import (
	"testing"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/stretchr/testify/require"
)

func TestFindPropertyAndValue(t *testing.T) {
	t.Parallel()

	value := &propertypb.PropertyValue{}
	value.SetStringValue("A+")
	collection := &propertypb.Properties{Properties: []*propertypb.Property{{
		Key:   "user.blood_type",
		Value: value,
	}}}

	require.Same(t, collection.Properties[0], Find(collection, "user.blood_type"))
	require.Same(t, value, FindValue(collection, "user.blood_type"))
	require.Nil(t, Find(nil, "user.blood_type"))
	require.Nil(t, FindValue(collection, "user.unknown"))
}
