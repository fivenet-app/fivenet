// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.4
// source: resources/livemap/livemap.proto

package livemap

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	centrum "github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	jobs "github.com/fivenet-app/fivenet/gen/go/proto/resources/jobs"
	timestamp "github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MarkerType int32

const (
	MarkerType_MARKER_TYPE_UNSPECIFIED MarkerType = 0
	MarkerType_MARKER_TYPE_DOT         MarkerType = 1
	MarkerType_MARKER_TYPE_CIRCLE      MarkerType = 2
	MarkerType_MARKER_TYPE_ICON        MarkerType = 3
)

// Enum value maps for MarkerType.
var (
	MarkerType_name = map[int32]string{
		0: "MARKER_TYPE_UNSPECIFIED",
		1: "MARKER_TYPE_DOT",
		2: "MARKER_TYPE_CIRCLE",
		3: "MARKER_TYPE_ICON",
	}
	MarkerType_value = map[string]int32{
		"MARKER_TYPE_UNSPECIFIED": 0,
		"MARKER_TYPE_DOT":         1,
		"MARKER_TYPE_CIRCLE":      2,
		"MARKER_TYPE_ICON":        3,
	}
)

func (x MarkerType) Enum() *MarkerType {
	p := new(MarkerType)
	*p = x
	return p
}

func (x MarkerType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MarkerType) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_livemap_livemap_proto_enumTypes[0].Descriptor()
}

func (MarkerType) Type() protoreflect.EnumType {
	return &file_resources_livemap_livemap_proto_enumTypes[0]
}

func (x MarkerType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MarkerType.Descriptor instead.
func (MarkerType) EnumDescriptor() ([]byte, []int) {
	return file_resources_livemap_livemap_proto_rawDescGZIP(), []int{0}
}

type UserMarker struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	UserId    int32                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	X         float64                `protobuf:"fixed64,2,opt,name=x,proto3" json:"x,omitempty"`
	Y         float64                `protobuf:"fixed64,3,opt,name=y,proto3" json:"y,omitempty"`
	UpdatedAt *timestamp.Timestamp   `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	// @sanitize: method=StripTags
	Postal *string `protobuf:"bytes,5,opt,name=postal,proto3,oneof" json:"postal,omitempty"`
	// @sanitize: method=StripTags
	Color         *string         `protobuf:"bytes,6,opt,name=color,proto3,oneof" json:"color,omitempty"`
	Job           string          `protobuf:"bytes,7,opt,name=job,proto3" json:"job,omitempty"`
	JobLabel      string          `protobuf:"bytes,8,opt,name=job_label,json=jobLabel,proto3" json:"job_label,omitempty"`
	User          *jobs.Colleague `protobuf:"bytes,9,opt,name=user,proto3" json:"user,omitempty" alias:"user"` // @gotags: alias:"user"
	UnitId        *uint64         `protobuf:"varint,10,opt,name=unit_id,json=unitId,proto3,oneof" json:"unit_id,omitempty"`
	Unit          *centrum.Unit   `protobuf:"bytes,11,opt,name=unit,proto3,oneof" json:"unit,omitempty"`
	Hidden        bool            `protobuf:"varint,12,opt,name=hidden,proto3" json:"hidden,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserMarker) Reset() {
	*x = UserMarker{}
	mi := &file_resources_livemap_livemap_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserMarker) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserMarker) ProtoMessage() {}

func (x *UserMarker) ProtoReflect() protoreflect.Message {
	mi := &file_resources_livemap_livemap_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserMarker.ProtoReflect.Descriptor instead.
func (*UserMarker) Descriptor() ([]byte, []int) {
	return file_resources_livemap_livemap_proto_rawDescGZIP(), []int{0}
}

func (x *UserMarker) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UserMarker) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *UserMarker) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}

func (x *UserMarker) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *UserMarker) GetPostal() string {
	if x != nil && x.Postal != nil {
		return *x.Postal
	}
	return ""
}

func (x *UserMarker) GetColor() string {
	if x != nil && x.Color != nil {
		return *x.Color
	}
	return ""
}

func (x *UserMarker) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *UserMarker) GetJobLabel() string {
	if x != nil {
		return x.JobLabel
	}
	return ""
}

func (x *UserMarker) GetUser() *jobs.Colleague {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *UserMarker) GetUnitId() uint64 {
	if x != nil && x.UnitId != nil {
		return *x.UnitId
	}
	return 0
}

func (x *UserMarker) GetUnit() *centrum.Unit {
	if x != nil {
		return x.Unit
	}
	return nil
}

func (x *UserMarker) GetHidden() bool {
	if x != nil {
		return x.Hidden
	}
	return false
}

type MarkerMarker struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Id        uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	X         float64                `protobuf:"fixed64,2,opt,name=x,proto3" json:"x,omitempty"`
	Y         float64                `protobuf:"fixed64,3,opt,name=y,proto3" json:"y,omitempty"`
	CreatedAt *timestamp.Timestamp   `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp   `protobuf:"bytes,5,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	ExpiresAt *timestamp.Timestamp   `protobuf:"bytes,6,opt,name=expires_at,json=expiresAt,proto3,oneof" json:"expires_at,omitempty"`
	DeletedAt *timestamp.Timestamp   `protobuf:"bytes,7,opt,name=deleted_at,json=deletedAt,proto3,oneof" json:"deleted_at,omitempty"`
	// @sanitize
	Name string `protobuf:"bytes,8,opt,name=name,proto3" json:"name,omitempty"`
	// @sanitize
	Description *string `protobuf:"bytes,9,opt,name=description,proto3,oneof" json:"description,omitempty"`
	// @sanitize: method=StripTags
	Postal *string `protobuf:"bytes,10,opt,name=postal,proto3,oneof" json:"postal,omitempty"`
	// @sanitize: method=StripTags
	Color         *string          `protobuf:"bytes,11,opt,name=color,proto3,oneof" json:"color,omitempty"`
	Job           string           `protobuf:"bytes,12,opt,name=job,proto3" json:"job,omitempty"`
	JobLabel      string           `protobuf:"bytes,13,opt,name=job_label,json=jobLabel,proto3" json:"job_label,omitempty"`
	Type          MarkerType       `protobuf:"varint,14,opt,name=type,proto3,enum=resources.livemap.MarkerType" json:"type,omitempty" alias:"markerType"` // @gotags: alias:"markerType"
	Data          *MarkerData      `protobuf:"bytes,15,opt,name=data,proto3" json:"data,omitempty" alias:"markerData"`                                    // @gotags: alias:"markerData"
	CreatorId     *int32           `protobuf:"varint,16,opt,name=creator_id,json=creatorId,proto3,oneof" json:"creator_id,omitempty"`
	Creator       *users.UserShort `protobuf:"bytes,17,opt,name=creator,proto3,oneof" json:"creator,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MarkerMarker) Reset() {
	*x = MarkerMarker{}
	mi := &file_resources_livemap_livemap_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MarkerMarker) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkerMarker) ProtoMessage() {}

func (x *MarkerMarker) ProtoReflect() protoreflect.Message {
	mi := &file_resources_livemap_livemap_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkerMarker.ProtoReflect.Descriptor instead.
func (*MarkerMarker) Descriptor() ([]byte, []int) {
	return file_resources_livemap_livemap_proto_rawDescGZIP(), []int{1}
}

func (x *MarkerMarker) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MarkerMarker) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *MarkerMarker) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}

func (x *MarkerMarker) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *MarkerMarker) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *MarkerMarker) GetExpiresAt() *timestamp.Timestamp {
	if x != nil {
		return x.ExpiresAt
	}
	return nil
}

func (x *MarkerMarker) GetDeletedAt() *timestamp.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

func (x *MarkerMarker) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *MarkerMarker) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *MarkerMarker) GetPostal() string {
	if x != nil && x.Postal != nil {
		return *x.Postal
	}
	return ""
}

func (x *MarkerMarker) GetColor() string {
	if x != nil && x.Color != nil {
		return *x.Color
	}
	return ""
}

func (x *MarkerMarker) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *MarkerMarker) GetJobLabel() string {
	if x != nil {
		return x.JobLabel
	}
	return ""
}

func (x *MarkerMarker) GetType() MarkerType {
	if x != nil {
		return x.Type
	}
	return MarkerType_MARKER_TYPE_UNSPECIFIED
}

func (x *MarkerMarker) GetData() *MarkerData {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *MarkerMarker) GetCreatorId() int32 {
	if x != nil && x.CreatorId != nil {
		return *x.CreatorId
	}
	return 0
}

func (x *MarkerMarker) GetCreator() *users.UserShort {
	if x != nil {
		return x.Creator
	}
	return nil
}

// @dbscanner
type MarkerData struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Data:
	//
	//	*MarkerData_Circle
	//	*MarkerData_Icon
	Data          isMarkerData_Data `protobuf_oneof:"data"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MarkerData) Reset() {
	*x = MarkerData{}
	mi := &file_resources_livemap_livemap_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MarkerData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkerData) ProtoMessage() {}

func (x *MarkerData) ProtoReflect() protoreflect.Message {
	mi := &file_resources_livemap_livemap_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkerData.ProtoReflect.Descriptor instead.
func (*MarkerData) Descriptor() ([]byte, []int) {
	return file_resources_livemap_livemap_proto_rawDescGZIP(), []int{2}
}

func (x *MarkerData) GetData() isMarkerData_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *MarkerData) GetCircle() *CircleMarker {
	if x != nil {
		if x, ok := x.Data.(*MarkerData_Circle); ok {
			return x.Circle
		}
	}
	return nil
}

func (x *MarkerData) GetIcon() *IconMarker {
	if x != nil {
		if x, ok := x.Data.(*MarkerData_Icon); ok {
			return x.Icon
		}
	}
	return nil
}

type isMarkerData_Data interface {
	isMarkerData_Data()
}

type MarkerData_Circle struct {
	Circle *CircleMarker `protobuf:"bytes,3,opt,name=circle,proto3,oneof"`
}

type MarkerData_Icon struct {
	Icon *IconMarker `protobuf:"bytes,4,opt,name=icon,proto3,oneof"`
}

func (*MarkerData_Circle) isMarkerData_Data() {}

func (*MarkerData_Icon) isMarkerData_Data() {}

type CircleMarker struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Radius        int32                  `protobuf:"varint,1,opt,name=radius,proto3" json:"radius,omitempty"`
	Opacity       *float32               `protobuf:"fixed32,2,opt,name=opacity,proto3,oneof" json:"opacity,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CircleMarker) Reset() {
	*x = CircleMarker{}
	mi := &file_resources_livemap_livemap_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CircleMarker) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CircleMarker) ProtoMessage() {}

func (x *CircleMarker) ProtoReflect() protoreflect.Message {
	mi := &file_resources_livemap_livemap_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CircleMarker.ProtoReflect.Descriptor instead.
func (*CircleMarker) Descriptor() ([]byte, []int) {
	return file_resources_livemap_livemap_proto_rawDescGZIP(), []int{3}
}

func (x *CircleMarker) GetRadius() int32 {
	if x != nil {
		return x.Radius
	}
	return 0
}

func (x *CircleMarker) GetOpacity() float32 {
	if x != nil && x.Opacity != nil {
		return *x.Opacity
	}
	return 0
}

type IconMarker struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// @sanitize: method=StripTags
	Icon          string `protobuf:"bytes,1,opt,name=icon,proto3" json:"icon,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *IconMarker) Reset() {
	*x = IconMarker{}
	mi := &file_resources_livemap_livemap_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IconMarker) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IconMarker) ProtoMessage() {}

func (x *IconMarker) ProtoReflect() protoreflect.Message {
	mi := &file_resources_livemap_livemap_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IconMarker.ProtoReflect.Descriptor instead.
func (*IconMarker) Descriptor() ([]byte, []int) {
	return file_resources_livemap_livemap_proto_rawDescGZIP(), []int{4}
}

func (x *IconMarker) GetIcon() string {
	if x != nil {
		return x.Icon
	}
	return ""
}

type Coords struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	X             float64                `protobuf:"fixed64,1,opt,name=x,proto3" json:"x,omitempty"`
	Y             float64                `protobuf:"fixed64,2,opt,name=y,proto3" json:"y,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Coords) Reset() {
	*x = Coords{}
	mi := &file_resources_livemap_livemap_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Coords) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Coords) ProtoMessage() {}

func (x *Coords) ProtoReflect() protoreflect.Message {
	mi := &file_resources_livemap_livemap_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Coords.ProtoReflect.Descriptor instead.
func (*Coords) Descriptor() ([]byte, []int) {
	return file_resources_livemap_livemap_proto_rawDescGZIP(), []int{5}
}

func (x *Coords) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Coords) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}

var File_resources_livemap_livemap_proto protoreflect.FileDescriptor

const file_resources_livemap_livemap_proto_rawDesc = "" +
	"\n" +
	"\x1fresources/livemap/livemap.proto\x12\x11resources.livemap\x1a\x1dresources/centrum/units.proto\x1a\x1fresources/jobs/colleagues.proto\x1a#resources/timestamp/timestamp.proto\x1a\x1bresources/users/users.proto\x1a\x17validate/validate.proto\"\xf4\x03\n" +
	"\n" +
	"UserMarker\x12 \n" +
	"\auser_id\x18\x01 \x01(\x05B\a\xfaB\x04\x1a\x02 \x00R\x06userId\x12\f\n" +
	"\x01x\x18\x02 \x01(\x01R\x01x\x12\f\n" +
	"\x01y\x18\x03 \x01(\x01R\x01y\x12B\n" +
	"\n" +
	"updated_at\x18\x04 \x01(\v2\x1e.resources.timestamp.TimestampH\x00R\tupdatedAt\x88\x01\x01\x12$\n" +
	"\x06postal\x18\x05 \x01(\tB\a\xfaB\x04r\x02\x180H\x01R\x06postal\x88\x01\x01\x126\n" +
	"\x05color\x18\x06 \x01(\tB\x1b\xfaB\x18r\x162\x11^#[A-Fa-f0-9]{6}$\x98\x01\aH\x02R\x05color\x88\x01\x01\x12\x19\n" +
	"\x03job\x18\a \x01(\tB\a\xfaB\x04r\x02\x18\x14R\x03job\x12\x1b\n" +
	"\tjob_label\x18\b \x01(\tR\bjobLabel\x12-\n" +
	"\x04user\x18\t \x01(\v2\x19.resources.jobs.ColleagueR\x04user\x12\x1c\n" +
	"\aunit_id\x18\n" +
	" \x01(\x04H\x03R\x06unitId\x88\x01\x01\x120\n" +
	"\x04unit\x18\v \x01(\v2\x17.resources.centrum.UnitH\x04R\x04unit\x88\x01\x01\x12\x16\n" +
	"\x06hidden\x18\f \x01(\bR\x06hiddenB\r\n" +
	"\v_updated_atB\t\n" +
	"\a_postalB\b\n" +
	"\x06_colorB\n" +
	"\n" +
	"\b_unit_idB\a\n" +
	"\x05_unit\"\xf1\x06\n" +
	"\fMarkerMarker\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12\f\n" +
	"\x01x\x18\x02 \x01(\x01R\x01x\x12\f\n" +
	"\x01y\x18\x03 \x01(\x01R\x01y\x12B\n" +
	"\n" +
	"created_at\x18\x04 \x01(\v2\x1e.resources.timestamp.TimestampH\x00R\tcreatedAt\x88\x01\x01\x12B\n" +
	"\n" +
	"updated_at\x18\x05 \x01(\v2\x1e.resources.timestamp.TimestampH\x01R\tupdatedAt\x88\x01\x01\x12B\n" +
	"\n" +
	"expires_at\x18\x06 \x01(\v2\x1e.resources.timestamp.TimestampH\x02R\texpiresAt\x88\x01\x01\x12B\n" +
	"\n" +
	"deleted_at\x18\a \x01(\v2\x1e.resources.timestamp.TimestampH\x03R\tdeletedAt\x88\x01\x01\x12\x1e\n" +
	"\x04name\x18\b \x01(\tB\n" +
	"\xfaB\ar\x05\x10\x01\x18\xff\x01R\x04name\x12%\n" +
	"\vdescription\x18\t \x01(\tH\x04R\vdescription\x88\x01\x01\x12$\n" +
	"\x06postal\x18\n" +
	" \x01(\tB\a\xfaB\x04r\x02\x180H\x05R\x06postal\x88\x01\x01\x126\n" +
	"\x05color\x18\v \x01(\tB\x1b\xfaB\x18r\x162\x11^#[A-Fa-f0-9]{6}$\x98\x01\aH\x06R\x05color\x88\x01\x01\x12\x19\n" +
	"\x03job\x18\f \x01(\tB\a\xfaB\x04r\x02\x18\x14R\x03job\x12\x1b\n" +
	"\tjob_label\x18\r \x01(\tR\bjobLabel\x121\n" +
	"\x04type\x18\x0e \x01(\x0e2\x1d.resources.livemap.MarkerTypeR\x04type\x121\n" +
	"\x04data\x18\x0f \x01(\v2\x1d.resources.livemap.MarkerDataR\x04data\x12+\n" +
	"\n" +
	"creator_id\x18\x10 \x01(\x05B\a\xfaB\x04\x1a\x02 \x00H\aR\tcreatorId\x88\x01\x01\x129\n" +
	"\acreator\x18\x11 \x01(\v2\x1a.resources.users.UserShortH\bR\acreator\x88\x01\x01B\r\n" +
	"\v_created_atB\r\n" +
	"\v_updated_atB\r\n" +
	"\v_expires_atB\r\n" +
	"\v_deleted_atB\x0e\n" +
	"\f_descriptionB\t\n" +
	"\a_postalB\b\n" +
	"\x06_colorB\r\n" +
	"\v_creator_idB\n" +
	"\n" +
	"\b_creator\"\x89\x01\n" +
	"\n" +
	"MarkerData\x129\n" +
	"\x06circle\x18\x03 \x01(\v2\x1f.resources.livemap.CircleMarkerH\x00R\x06circle\x123\n" +
	"\x04icon\x18\x04 \x01(\v2\x1d.resources.livemap.IconMarkerH\x00R\x04iconB\v\n" +
	"\x04data\x12\x03\xf8B\x01\"b\n" +
	"\fCircleMarker\x12\x16\n" +
	"\x06radius\x18\x01 \x01(\x05R\x06radius\x12.\n" +
	"\aopacity\x18\x02 \x01(\x02B\x0f\xfaB\f\n" +
	"\n" +
	"\x1d\x00\x00\x96B-\x00\x00\x80?H\x00R\aopacity\x88\x01\x01B\n" +
	"\n" +
	"\b_opacity\"*\n" +
	"\n" +
	"IconMarker\x12\x1c\n" +
	"\x04icon\x18\x01 \x01(\tB\b\xfaB\x05r\x03\x18\x80\x01R\x04icon\"$\n" +
	"\x06Coords\x12\f\n" +
	"\x01x\x18\x01 \x01(\x01R\x01x\x12\f\n" +
	"\x01y\x18\x02 \x01(\x01R\x01y*l\n" +
	"\n" +
	"MarkerType\x12\x1b\n" +
	"\x17MARKER_TYPE_UNSPECIFIED\x10\x00\x12\x13\n" +
	"\x0fMARKER_TYPE_DOT\x10\x01\x12\x16\n" +
	"\x12MARKER_TYPE_CIRCLE\x10\x02\x12\x14\n" +
	"\x10MARKER_TYPE_ICON\x10\x03BGZEgithub.com/fivenet-app/fivenet/gen/go/proto/resources/livemap;livemapb\x06proto3"

var (
	file_resources_livemap_livemap_proto_rawDescOnce sync.Once
	file_resources_livemap_livemap_proto_rawDescData []byte
)

func file_resources_livemap_livemap_proto_rawDescGZIP() []byte {
	file_resources_livemap_livemap_proto_rawDescOnce.Do(func() {
		file_resources_livemap_livemap_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_livemap_livemap_proto_rawDesc), len(file_resources_livemap_livemap_proto_rawDesc)))
	})
	return file_resources_livemap_livemap_proto_rawDescData
}

var file_resources_livemap_livemap_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resources_livemap_livemap_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_resources_livemap_livemap_proto_goTypes = []any{
	(MarkerType)(0),             // 0: resources.livemap.MarkerType
	(*UserMarker)(nil),          // 1: resources.livemap.UserMarker
	(*MarkerMarker)(nil),        // 2: resources.livemap.MarkerMarker
	(*MarkerData)(nil),          // 3: resources.livemap.MarkerData
	(*CircleMarker)(nil),        // 4: resources.livemap.CircleMarker
	(*IconMarker)(nil),          // 5: resources.livemap.IconMarker
	(*Coords)(nil),              // 6: resources.livemap.Coords
	(*timestamp.Timestamp)(nil), // 7: resources.timestamp.Timestamp
	(*jobs.Colleague)(nil),      // 8: resources.jobs.Colleague
	(*centrum.Unit)(nil),        // 9: resources.centrum.Unit
	(*users.UserShort)(nil),     // 10: resources.users.UserShort
}
var file_resources_livemap_livemap_proto_depIdxs = []int32{
	7,  // 0: resources.livemap.UserMarker.updated_at:type_name -> resources.timestamp.Timestamp
	8,  // 1: resources.livemap.UserMarker.user:type_name -> resources.jobs.Colleague
	9,  // 2: resources.livemap.UserMarker.unit:type_name -> resources.centrum.Unit
	7,  // 3: resources.livemap.MarkerMarker.created_at:type_name -> resources.timestamp.Timestamp
	7,  // 4: resources.livemap.MarkerMarker.updated_at:type_name -> resources.timestamp.Timestamp
	7,  // 5: resources.livemap.MarkerMarker.expires_at:type_name -> resources.timestamp.Timestamp
	7,  // 6: resources.livemap.MarkerMarker.deleted_at:type_name -> resources.timestamp.Timestamp
	0,  // 7: resources.livemap.MarkerMarker.type:type_name -> resources.livemap.MarkerType
	3,  // 8: resources.livemap.MarkerMarker.data:type_name -> resources.livemap.MarkerData
	10, // 9: resources.livemap.MarkerMarker.creator:type_name -> resources.users.UserShort
	4,  // 10: resources.livemap.MarkerData.circle:type_name -> resources.livemap.CircleMarker
	5,  // 11: resources.livemap.MarkerData.icon:type_name -> resources.livemap.IconMarker
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_resources_livemap_livemap_proto_init() }
func file_resources_livemap_livemap_proto_init() {
	if File_resources_livemap_livemap_proto != nil {
		return
	}
	file_resources_livemap_livemap_proto_msgTypes[0].OneofWrappers = []any{}
	file_resources_livemap_livemap_proto_msgTypes[1].OneofWrappers = []any{}
	file_resources_livemap_livemap_proto_msgTypes[2].OneofWrappers = []any{
		(*MarkerData_Circle)(nil),
		(*MarkerData_Icon)(nil),
	}
	file_resources_livemap_livemap_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_livemap_livemap_proto_rawDesc), len(file_resources_livemap_livemap_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_livemap_livemap_proto_goTypes,
		DependencyIndexes: file_resources_livemap_livemap_proto_depIdxs,
		EnumInfos:         file_resources_livemap_livemap_proto_enumTypes,
		MessageInfos:      file_resources_livemap_livemap_proto_msgTypes,
	}.Build()
	File_resources_livemap_livemap_proto = out.File
	file_resources_livemap_livemap_proto_goTypes = nil
	file_resources_livemap_livemap_proto_depIdxs = nil
}
