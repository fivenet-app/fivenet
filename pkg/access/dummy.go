//nolint:protogetter // Package access provides dummy implementations of access-related types for testing or placeholder purposes. Direct proto field access is needed.
package access

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// DummyAccessLevel is a dummy enum type for access levels.
type DummyAccessLevel int32

// DummyAccessLevel_DUMMY_ACCESS_LEVEL_UNSPECIFIED is the default dummy value.
const DummyAccessLevel_DUMMY_ACCESS_LEVEL_UNSPECIFIED DummyAccessLevel = 0

// Number returns the dummy enum number (always 0).
func (x DummyAccessLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(0)
}

// DummyJobAccess is a dummy implementation of job access for testing or placeholder purposes.
type DummyJobAccess[V protoutils.ProtoEnum] struct {
	// Access holds the access level or permission value.
	Access V
}

// GetId returns a dummy ID (always 0).
func (x *DummyJobAccess[V]) GetId() uint64 { return 0 }

// GetTargetId returns a dummy target ID (always 0).
func (x *DummyJobAccess[V]) GetTargetId() uint64 { return 0 }

// GetJob returns a dummy job string (always empty).
func (x *DummyJobAccess[V]) GetJob() string { return "" }

// GetMinimumGrade returns a dummy grade (always 0).
func (x *DummyJobAccess[V]) GetMinimumGrade() int32 { return 0 }

// SetMinimumGrade is a no-op for the dummy implementation.
func (x *DummyJobAccess[V]) SetMinimumGrade(int32) {}

// GetAccess returns the stored access value.
func (x *DummyJobAccess[V]) GetAccess() V { return x.Access }

// SetAccess is a no-op for the dummy implementation.
func (x *DummyJobAccess[V]) SetAccess(V) {}

// ProtoReflect returns nil for the dummy implementation.
func (x *DummyJobAccess[V]) ProtoReflect() protoreflect.Message { return nil }

// DummyUserAccess is a dummy implementation of user access for testing or placeholder purposes.
type DummyUserAccess[V protoutils.ProtoEnum] struct {
	// Access holds the access level or permission value.
	Access V
}

// GetId returns a dummy ID (always 0).
func (x *DummyUserAccess[V]) GetId() uint64 { return 0 }

// GetTargetId returns a dummy target ID (always 0).
func (x *DummyUserAccess[V]) GetTargetId() uint64 { return 0 }

// GetUserId returns a dummy user ID (always 0).
func (x *DummyUserAccess[V]) GetUserId() int32 { return 0 }

// SetUserId is a no-op for the dummy implementation.
func (x *DummyUserAccess[V]) SetUserId(int32) {}

// GetAccess returns the stored access value.
func (x *DummyUserAccess[V]) GetAccess() V { return x.Access }

// SetAccess is a no-op for the dummy implementation.
func (x *DummyUserAccess[V]) SetAccess(V) {}

// ProtoReflect returns nil for the dummy implementation.
func (x *DummyUserAccess[V]) ProtoReflect() protoreflect.Message { return nil }

// DummyQualificationAccess is a dummy implementation of qualification access for testing or placeholder purposes.
type DummyQualificationAccess[V protoutils.ProtoEnum] struct {
	// Access holds the access level or permission value.
	Access V
}

// GetId returns a dummy ID (always 0).
func (x *DummyQualificationAccess[V]) GetId() uint64 { return 0 }

// GetTargetId returns a dummy target ID (always 0).
func (x *DummyQualificationAccess[V]) GetTargetId() uint64 { return 0 }

// GetQualificationId returns a dummy qualification ID (always 0).
func (x *DummyQualificationAccess[V]) GetQualificationId() uint64 { return 0 }

// SetQualificationId is a no-op for the dummy implementation.
func (x *DummyQualificationAccess[V]) SetQualificationId(uint64) {}

// GetAccess returns the stored access value.
func (x *DummyQualificationAccess[V]) GetAccess() V { return x.Access }

// SetAccess is a no-op for the dummy implementation.
func (x *DummyQualificationAccess[V]) SetAccess(V) {}

// ProtoReflect returns nil for the dummy implementation.
func (x *DummyQualificationAccess[V]) ProtoReflect() protoreflect.Message { return nil }
