package properties

import (
	"context"
	"errors"
	"math"
	"time"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
)

// MutationLifecycle is the transport-neutral lifecycle activity payload
// derived from a property mutation. Store/service layers convert it to their
// scope-specific activity protobuf.
type MutationLifecycle struct {
	ExpiresAt          *timestamp.Timestamp
	ExpiryCleared      bool
	StateReason        *string
	StateReasonCleared bool
}

func MutationLifecycleForMutation(mutation *propertypb.PropertyMutation) *MutationLifecycle {
	if mutation == nil {
		return nil
	}
	state := mutation.GetState()
	if state == nil && mutation.Reason == nil && mutation.ClearReason == nil {
		return nil
	}
	return &MutationLifecycle{
		ExpiresAt: state.GetExpiresAt(), ExpiryCleared: state.GetClearExpiry() != nil,
		StateReason: mutation.Reason, StateReasonCleared: mutation.ClearReason != nil,
	}
}

func SameOptionalString(left, right *string) bool {
	return left == nil && right == nil || left != nil && right != nil && *left == *right
}

func SameOptionalTime(left, right *time.Time) bool {
	return left == nil && right == nil || left != nil && right != nil && left.Equal(*right)
}

var (
	ErrInvalidPropertyChange = errors.New("invalid property change")
	ErrPropertyChangeDenied  = errors.New("property change denied")
)

// ScalarChange is an internal persistence representation derived from the
// public PropertyMutation contract. It is intentionally not a protobuf API.
type ScalarChange struct {
	Key   string
	Value *propertypb.PropertyValue
}

// ValidateMutationState checks that a lifecycle patch is permitted by the
// resolved definition. Persistence layers call this before touching their
// generic state table; managed handlers use the same policy gate.
func ValidateMutationState(
	definition *propertypb.PropertyDefinition,
	mutation *propertypb.PropertyMutation,
) error {
	if definition == nil || mutation == nil {
		return ErrInvalidPropertyChange
	}
	state := mutation.GetState()
	if state == nil {
		return nil
	}
	if !definition.GetEnabled() || (state.GetExpiresAt() == nil && state.GetClearExpiry() == nil) {
		return ErrInvalidPropertyChange
	}
	policy := definition.GetStatePolicy()
	if policy == nil || !policy.GetEnabled() {
		return ErrInvalidPropertyChange
	}
	if (state.GetExpiresAt() != nil || state.GetClearExpiry() != nil) &&
		!policy.GetSupportsExpiry() {
		return ErrInvalidPropertyChange
	}
	if expiresAt := state.GetExpiresAt(); expiresAt != nil &&
		(expiresAt.GetTimestamp() == nil || expiresAt.GetTimestamp().CheckValid() != nil) {
		return ErrInvalidPropertyChange
	}
	return nil
}

// ValidateMutationReason enforces the definition's mutation-scoped reason
// policy independently from expiry/state capabilities.
func ValidateMutationReason(
	definition *propertypb.PropertyDefinition,
	mutation *propertypb.PropertyMutation,
) error {
	if definition == nil || mutation == nil {
		return ErrInvalidPropertyChange
	}
	hasReason := mutation.Reason != nil
	hasClear := mutation.ClearReason != nil
	if hasClear && !mutation.GetClearReason() {
		return ErrInvalidPropertyChange
	}
	if hasReason && hasClear {
		return ErrInvalidPropertyChange
	}
	switch definition.GetReasonMode() {
	case propertypb.PropertyReasonMode_PROPERTY_REASON_MODE_NONE:
		if hasReason || hasClear {
			return ErrInvalidPropertyChange
		}
	case propertypb.PropertyReasonMode_PROPERTY_REASON_MODE_UNSPECIFIED,
		propertypb.PropertyReasonMode_PROPERTY_REASON_MODE_REQUIRED:
		// A clear removes the property state and therefore needs a concrete
		// rationale in its immutable activity, not a state-reason clear.
		if !hasReason {
			return ErrInvalidPropertyChange
		}
	}
	return nil
}

// ValidateMutationShape rejects ambiguous or no-op mutations before callers
// split scalar and lifecycle handling. A value clear already removes its state,
// so it cannot be combined with a state patch.
func ValidateMutationShape(mutation *propertypb.PropertyMutation) error {
	if mutation == nil || mutation.GetKey() == "" {
		return ErrInvalidPropertyChange
	}
	if mutation.GetSetValue() == nil && mutation.GetClearValue() == nil &&
		mutation.GetIntegerMath() == nil &&
		mutation.GetState() == nil &&
		mutation.Reason == nil &&
		mutation.ClearReason == nil {
		return ErrInvalidPropertyChange
	}
	if mutation.GetClearValue() != nil && mutation.GetState() != nil {
		return ErrInvalidPropertyChange
	}
	return nil
}

// ApplyIntegerMath calculates an integer mutation's effective value. Callers
// must load currentValue while holding the property's row lock. Missing values
// are represented by zero. The definition's integer bounds clamp the result.
func ApplyIntegerMath(
	definition *propertypb.PropertyDefinition,
	currentValue int64,
	operation *propertypb.IntegerMathOperation,
) (int64, error) {
	if definition == nil || !definition.GetEnabled() ||
		definition.GetValueType() != propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_INTEGER ||
		operation == nil || operation.GetOperand() <= 0 {
		return 0, ErrInvalidPropertyChange
	}

	operand := operation.GetOperand()
	var value int64
	switch operation.GetAction() {
	case propertypb.IntegerMathAction_INTEGER_MATH_ACTION_ADD:
		if (operand > 0 && currentValue > math.MaxInt64-operand) ||
			(operand < 0 && currentValue < math.MinInt64-operand) {
			return 0, ErrInvalidPropertyChange
		}
		value = currentValue + operand
	case propertypb.IntegerMathAction_INTEGER_MATH_ACTION_SUBTRACT:
		if (operand > 0 && currentValue < math.MinInt64+operand) ||
			(operand < 0 && currentValue > math.MaxInt64+operand) {
			return 0, ErrInvalidPropertyChange
		}
		value = currentValue - operand
	default:
		return 0, ErrInvalidPropertyChange
	}

	constraints := definition.GetValueConstraints()
	if constraints != nil {
		if constraints.IntegerMin != nil && value < constraints.GetIntegerMin() {
			value = constraints.GetIntegerMin()
		}
		if constraints.IntegerMax != nil && value > constraints.GetIntegerMax() {
			value = constraints.GetIntegerMax()
		}
	}
	return value, nil
}

// ValidateMutationKeys rejects duplicate keys and malformed mutations in one request.
func ValidateMutationKeys(mutations []*propertypb.PropertyMutation) error {
	keys := make(map[string]struct{}, len(mutations))
	for _, mutation := range mutations {
		if err := ValidateMutationShape(mutation); err != nil {
			return err
		}
		if _, exists := keys[mutation.GetKey()]; exists {
			return ErrInvalidPropertyChange
		}
		keys[mutation.GetKey()] = struct{}{}
	}
	return nil
}

type DefinitionResolver func(context.Context, []string) (map[string]*propertypb.PropertyDefinition, error)

// ValidateCustomPropertyChanges validates internal scalar changes derived from
// PropertyMutation before storage. An absent value clears a property.
func ValidateCustomPropertyChanges(
	ctx context.Context,
	scope propertypb.PropertyScope,
	changes []*ScalarChange,
	allowed func(string) bool,
	resolve DefinitionResolver,
) error {
	if len(changes) == 0 {
		return nil
	}
	if allowed == nil || resolve == nil {
		return ErrInvalidPropertyChange
	}

	keys := make([]string, 0, len(changes))
	seen := make(map[string]struct{}, len(changes))
	for _, change := range changes {
		if change == nil || !IsCustomPropertyKey(scope, change.Key) {
			return ErrInvalidPropertyChange
		}
		if !allowed(change.Key) {
			return ErrPropertyChangeDenied
		}
		if _, ok := seen[change.Key]; ok {
			return ErrInvalidPropertyChange
		}
		seen[change.Key] = struct{}{}
		keys = append(keys, change.Key)
	}

	definitions, err := resolve(ctx, keys)
	if err != nil {
		return err
	}
	for _, change := range changes {
		if change.Value == nil {
			continue
		}
		valueType, _, _, _, _, _, _, err := ValueColumns(change.Value)
		definition := definitions[change.Key]
		if err != nil || definition == nil || !definition.GetEnabled() ||
			definition.GetSource() != propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM ||
			definition.GetValueType() != valueType ||
			ValidateValueConstraints(definition, change.Value) != nil {
			return ErrInvalidPropertyChange
		}
	}

	return nil
}
