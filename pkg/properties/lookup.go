package properties

import propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"

// Find returns the property identified by key, or nil if the collection does
// not contain it. It is deliberately nil-safe so it can be exposed to
// document templates without requiring callers to guard absent collections.
func Find(collection *propertypb.Properties, key string) *propertypb.Property {
	if collection == nil {
		return nil
	}
	for _, property := range collection.GetProperties() {
		if property != nil && property.GetKey() == key {
			return property
		}
	}
	return nil
}

// FindValue returns the scalar value identified by key, or nil when the
// property or its value is absent.
func FindValue(collection *propertypb.Properties, key string) *propertypb.PropertyValue {
	return Find(collection, key).GetValue()
}
