package properties

import (
	"errors"
	"time"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
)

// ManagedProperty describes a code-owned scalar/lifecycle property. Callers
// use it for routing; authorization remains scope-specific.
type ManagedProperty struct {
	Scope       propertypb.PropertyScope
	Key         string
	HandlerName string
	Apply       func(WantedProjection, *propertypb.PropertyMutation) error
}

var managedProperties = map[propertypb.PropertyScope]map[string]ManagedProperty{
	propertypb.PropertyScope_PROPERTY_SCOPE_USER: {
		"user.wanted": {
			Scope:       propertypb.PropertyScope_PROPERTY_SCOPE_USER,
			Key:         "user.wanted",
			HandlerName: "wanted",
			Apply:       ApplyWantedMutation,
		},
	},
	propertypb.PropertyScope_PROPERTY_SCOPE_VEHICLE: {
		"vehicle.wanted": {
			Scope:       propertypb.PropertyScope_PROPERTY_SCOPE_VEHICLE,
			Key:         "vehicle.wanted",
			HandlerName: "wanted",
			Apply:       ApplyWantedMutation,
		},
	},
}

// ManagedMutation looks up code-owned system handlers. Custom keys
// intentionally never enter managed routing.
func ManagedMutation(scope propertypb.PropertyScope, key string) (ManagedProperty, bool) {
	property, ok := managedProperties[scope][key]
	if !ok {
		return ManagedProperty{}, false
	}
	return property, ok
}

// WantedProjection is the legacy read projection updated by the managed
// wanted property. Both user and vehicle projections implement it.
type WantedProjection interface {
	SetWanted(bool)
	SetWantedTill(*timestamp.Timestamp)
}

// ApplyWantedMutation is the shared bridge between the public
// PropertyMutation contract and established wanted projections. It
// intentionally leaves state-only value fields unset.
func ApplyWantedMutation(managed WantedProjection, mutation *propertypb.PropertyMutation) error {
	if managed == nil {
		return errors.New("nil wanted projection")
	}
	if err := applyWantedValue(mutation, managed.SetWanted); err != nil {
		return err
	}
	if state := mutation.GetState(); state != nil {
		applyWantedExpiry(state, managed.SetWantedTill)
	}
	return nil
}

func applyWantedValue(mutation *propertypb.PropertyMutation, setWanted func(bool)) error {
	if mutation == nil {
		return errors.New("missing mutation")
	}
	if value := mutation.GetSetValue(); value != nil {
		if _, ok := value.GetValue().(*propertypb.PropertyValue_BoolValue); !ok {
			return errors.New("wanted value must be boolean")
		}
		setWanted(value.GetBoolValue())
	} else if mutation.GetClearValue() != nil {
		setWanted(false)
	}
	return nil
}

func applyWantedExpiry(
	state *propertypb.PropertyStatePatch,
	setWantedTill func(*timestamp.Timestamp),
) {
	if expiresAt := state.GetExpiresAt(); expiresAt != nil {
		setWantedTill(expiresAt)
	} else if state.GetClearExpiry() != nil {
		setWantedTill(timestamp.New(time.Time{}))
	}
}

// ClampWantedExpiry bounds an explicit managed wanted expiry. It deliberately
// leaves clear_expiry untouched: an explicit clear is distinct from a maximum
// duration, which the wanted housekeeper enforces from active_at.
func ClampWantedExpiry(
	scope propertypb.PropertyScope,
	mutations []*propertypb.PropertyMutation,
	maxExpiry time.Time,
) {
	for _, mutation := range mutations {
		managed, ok := ManagedMutation(scope, mutation.GetKey())
		if !ok || managed.HandlerName != "wanted" || mutation.GetState() == nil ||
			mutation.GetState().GetClearExpiry() != nil {
			continue
		}
		expiresAt := mutation.GetState().GetExpiresAt()
		if expiresAt != nil && expiresAt.AsTime().After(maxExpiry) {
			mutation.GetState().SetExpiresAt(timestamp.New(maxExpiry))
		}
	}
}

// LifecyclePatchForMutation exposes the transport-neutral lifecycle details
// needed by resource-specific activity payloads.
func LifecyclePatchForMutation(
	mutation *propertypb.PropertyMutation,
) (*timestamp.Timestamp, bool, bool) {
	if mutation == nil || mutation.GetState() == nil {
		return nil, false, false
	}
	state := mutation.GetState()
	return state.GetExpiresAt(), state.GetClearExpiry() != nil, true
}
