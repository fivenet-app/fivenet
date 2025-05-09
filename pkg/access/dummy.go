package access

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type DummyAccessLevel int32

const DummyAccessLevel_DUMMY_ACCESS_LEVEL_UNSPECIFIED DummyAccessLevel = 0

func (x DummyAccessLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(0)
}

type DummyJobAccess[V protoutils.ProtoEnum] struct {
	Access V
}

func (x *DummyJobAccess[V]) GetId() uint64 {
	return 0
}

func (x *DummyJobAccess[V]) GetTargetId() uint64 {
	return 0
}

func (x *DummyJobAccess[V]) GetJob() string {
	return ""
}

func (x *DummyJobAccess[V]) GetMinimumGrade() int32 {
	return 0
}

func (x *DummyJobAccess[V]) SetMinimumGrade(int32) {
}

func (x *DummyJobAccess[V]) GetAccess() V {
	return x.Access
}

func (x *DummyJobAccess[V]) SetAccess(V) {}

func (x *DummyJobAccess[V]) ProtoReflect() protoreflect.Message {
	return nil
}

type DummyUserAccess[V protoutils.ProtoEnum] struct {
	Access V
}

func (x *DummyUserAccess[V]) GetId() uint64 {
	return 0
}

func (x *DummyUserAccess[V]) GetTargetId() uint64 {
	return 0
}

func (x *DummyUserAccess[V]) GetUserId() int32 {
	return 0
}

func (x *DummyUserAccess[V]) SetUserId(int32) {
}

func (x *DummyUserAccess[V]) GetAccess() V {
	return x.Access
}
func (x *DummyUserAccess[V]) SetAccess(V) {}

func (x *DummyUserAccess[V]) ProtoReflect() protoreflect.Message {
	return nil
}

type DummyQualificationAccess[V protoutils.ProtoEnum] struct {
	Access V
}

func (x *DummyQualificationAccess[V]) GetId() uint64 {
	return 0
}

func (x *DummyQualificationAccess[V]) GetTargetId() uint64 {
	return 0
}

func (x *DummyQualificationAccess[V]) GetQualificationId() uint64 {
	return 0
}

func (x *DummyQualificationAccess[V]) SetQualificationId(uint64) {}

func (x *DummyQualificationAccess[V]) GetAccess() V {
	return x.Access
}

func (x *DummyQualificationAccess[V]) SetAccess(V) {}

func (x *DummyQualificationAccess[V]) ProtoReflect() protoreflect.Message {
	return nil
}
