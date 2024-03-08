// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.21.12
// source: context.platform.proto

package v1

import (
	_ "chainguard.dev/sdk/proto/annotations"
	v1 "chainguard.dev/sdk/proto/platform/common/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RecordContext_Ancestry_Role int32

const (
	RecordContext_Ancestry_UNKNOWN RecordContext_Ancestry_Role = 0
	RecordContext_Ancestry_BASE    RecordContext_Ancestry_Role = 1
	RecordContext_Ancestry_DERIVED RecordContext_Ancestry_Role = 2
)

// Enum value maps for RecordContext_Ancestry_Role.
var (
	RecordContext_Ancestry_Role_name = map[int32]string{
		0: "UNKNOWN",
		1: "BASE",
		2: "DERIVED",
	}
	RecordContext_Ancestry_Role_value = map[string]int32{
		"UNKNOWN": 0,
		"BASE":    1,
		"DERIVED": 2,
	}
)

func (x RecordContext_Ancestry_Role) Enum() *RecordContext_Ancestry_Role {
	p := new(RecordContext_Ancestry_Role)
	*p = x
	return p
}

func (x RecordContext_Ancestry_Role) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RecordContext_Ancestry_Role) Descriptor() protoreflect.EnumDescriptor {
	return file_context_platform_proto_enumTypes[0].Descriptor()
}

func (RecordContext_Ancestry_Role) Type() protoreflect.EnumType {
	return &file_context_platform_proto_enumTypes[0]
}

func (x RecordContext_Ancestry_Role) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RecordContext_Ancestry_Role.Descriptor instead.
func (RecordContext_Ancestry_Role) EnumDescriptor() ([]byte, []int) {
	return file_context_platform_proto_rawDescGZIP(), []int{0, 1, 0}
}

type RecordContext_Variant_Role int32

const (
	RecordContext_Variant_UNKNOWN RecordContext_Variant_Role = 0
	RecordContext_Variant_INDEX   RecordContext_Variant_Role = 1
	RecordContext_Variant_VARIANT RecordContext_Variant_Role = 2
)

// Enum value maps for RecordContext_Variant_Role.
var (
	RecordContext_Variant_Role_name = map[int32]string{
		0: "UNKNOWN",
		1: "INDEX",
		2: "VARIANT",
	}
	RecordContext_Variant_Role_value = map[string]int32{
		"UNKNOWN": 0,
		"INDEX":   1,
		"VARIANT": 2,
	}
)

func (x RecordContext_Variant_Role) Enum() *RecordContext_Variant_Role {
	p := new(RecordContext_Variant_Role)
	*p = x
	return p
}

func (x RecordContext_Variant_Role) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RecordContext_Variant_Role) Descriptor() protoreflect.EnumDescriptor {
	return file_context_platform_proto_enumTypes[1].Descriptor()
}

func (RecordContext_Variant_Role) Type() protoreflect.EnumType {
	return &file_context_platform_proto_enumTypes[1]
}

func (x RecordContext_Variant_Role) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RecordContext_Variant_Role.Descriptor instead.
func (RecordContext_Variant_Role) EnumDescriptor() ([]byte, []int) {
	return file_context_platform_proto_rawDescGZIP(), []int{0, 2, 0}
}

type RecordContext struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id, The RecordContext UIDP at which this RecordContext resides.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// name of the RecordContext.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// a short description of this RecordContext.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// last_seen tracks the timestamp at which this RecordContext was last seen.
	LastSeen *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=last_seen,json=lastSeen,proto3" json:"last_seen,omitempty"`
	// Types that are assignable to Kind:
	//
	//	*RecordContext_Workload_
	//	*RecordContext_Ancestry_
	//	*RecordContext_Variant_
	Kind isRecordContext_Kind `protobuf_oneof:"kind"`
}

func (x *RecordContext) Reset() {
	*x = RecordContext{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_platform_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordContext) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordContext) ProtoMessage() {}

func (x *RecordContext) ProtoReflect() protoreflect.Message {
	mi := &file_context_platform_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordContext.ProtoReflect.Descriptor instead.
func (*RecordContext) Descriptor() ([]byte, []int) {
	return file_context_platform_proto_rawDescGZIP(), []int{0}
}

func (x *RecordContext) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RecordContext) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RecordContext) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *RecordContext) GetLastSeen() *timestamppb.Timestamp {
	if x != nil {
		return x.LastSeen
	}
	return nil
}

func (m *RecordContext) GetKind() isRecordContext_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (x *RecordContext) GetWorkload() *RecordContext_Workload {
	if x, ok := x.GetKind().(*RecordContext_Workload_); ok {
		return x.Workload
	}
	return nil
}

func (x *RecordContext) GetAncestry() *RecordContext_Ancestry {
	if x, ok := x.GetKind().(*RecordContext_Ancestry_); ok {
		return x.Ancestry
	}
	return nil
}

func (x *RecordContext) GetVariant() *RecordContext_Variant {
	if x, ok := x.GetKind().(*RecordContext_Variant_); ok {
		return x.Variant
	}
	return nil
}

type isRecordContext_Kind interface {
	isRecordContext_Kind()
}

type RecordContext_Workload_ struct {
	Workload *RecordContext_Workload `protobuf:"bytes,100,opt,name=workload,proto3,oneof"`
}

type RecordContext_Ancestry_ struct {
	Ancestry *RecordContext_Ancestry `protobuf:"bytes,101,opt,name=ancestry,proto3,oneof"`
}

type RecordContext_Variant_ struct {
	Variant *RecordContext_Variant `protobuf:"bytes,102,opt,name=variant,proto3,oneof"`
}

func (*RecordContext_Workload_) isRecordContext_Kind() {}

func (*RecordContext_Ancestry_) isRecordContext_Kind() {}

func (*RecordContext_Variant_) isRecordContext_Kind() {}

type RecordContextList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*RecordContext `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *RecordContextList) Reset() {
	*x = RecordContextList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_platform_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordContextList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordContextList) ProtoMessage() {}

func (x *RecordContextList) ProtoReflect() protoreflect.Message {
	mi := &file_context_platform_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordContextList.ProtoReflect.Descriptor instead.
func (*RecordContextList) Descriptor() ([]byte, []int) {
	return file_context_platform_proto_rawDescGZIP(), []int{1}
}

func (x *RecordContextList) GetItems() []*RecordContext {
	if x != nil {
		return x.Items
	}
	return nil
}

type RecordContextFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uidp *v1.UIDPFilter `protobuf:"bytes,1,opt,name=uidp,proto3" json:"uidp,omitempty"`
	// active_since is the timestamp after which the records should
	// have last been observed in the returned context.
	ActiveSince *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=active_since,json=activeSince,proto3" json:"active_since,omitempty"`
	// Allow filtering results based on matching fields in the context
	// ranging from just a particular "kind" to the exact workload shape.
	// Only specified fields will be used as part of the match.
	//
	// Types that are assignable to Kind:
	//
	//	*RecordContextFilter_Workload
	//	*RecordContextFilter_Ancestry
	//	*RecordContextFilter_Variant
	Kind isRecordContextFilter_Kind `protobuf_oneof:"kind"`
}

func (x *RecordContextFilter) Reset() {
	*x = RecordContextFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_platform_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordContextFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordContextFilter) ProtoMessage() {}

func (x *RecordContextFilter) ProtoReflect() protoreflect.Message {
	mi := &file_context_platform_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordContextFilter.ProtoReflect.Descriptor instead.
func (*RecordContextFilter) Descriptor() ([]byte, []int) {
	return file_context_platform_proto_rawDescGZIP(), []int{2}
}

func (x *RecordContextFilter) GetUidp() *v1.UIDPFilter {
	if x != nil {
		return x.Uidp
	}
	return nil
}

func (x *RecordContextFilter) GetActiveSince() *timestamppb.Timestamp {
	if x != nil {
		return x.ActiveSince
	}
	return nil
}

func (m *RecordContextFilter) GetKind() isRecordContextFilter_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (x *RecordContextFilter) GetWorkload() *RecordContext_Workload {
	if x, ok := x.GetKind().(*RecordContextFilter_Workload); ok {
		return x.Workload
	}
	return nil
}

func (x *RecordContextFilter) GetAncestry() *RecordContext_Ancestry {
	if x, ok := x.GetKind().(*RecordContextFilter_Ancestry); ok {
		return x.Ancestry
	}
	return nil
}

func (x *RecordContextFilter) GetVariant() *RecordContext_Variant {
	if x, ok := x.GetKind().(*RecordContextFilter_Variant); ok {
		return x.Variant
	}
	return nil
}

type isRecordContextFilter_Kind interface {
	isRecordContextFilter_Kind()
}

type RecordContextFilter_Workload struct {
	Workload *RecordContext_Workload `protobuf:"bytes,100,opt,name=workload,proto3,oneof"`
}

type RecordContextFilter_Ancestry struct {
	Ancestry *RecordContext_Ancestry `protobuf:"bytes,101,opt,name=ancestry,proto3,oneof"`
}

type RecordContextFilter_Variant struct {
	Variant *RecordContext_Variant `protobuf:"bytes,102,opt,name=variant,proto3,oneof"`
}

func (*RecordContextFilter_Workload) isRecordContextFilter_Kind() {}

func (*RecordContextFilter_Ancestry) isRecordContextFilter_Kind() {}

func (*RecordContextFilter_Variant) isRecordContextFilter_Kind() {}

// Workload contexts are added to existence records that have been
// observed running on a cluster.
type RecordContext_Workload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// remote_id is the metadata.uid of the workload in which this
	// container was observed.
	RemoteId string `protobuf:"bytes,1,opt,name=remote_id,json=remoteId,proto3" json:"remote_id,omitempty"`
}

func (x *RecordContext_Workload) Reset() {
	*x = RecordContext_Workload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_platform_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordContext_Workload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordContext_Workload) ProtoMessage() {}

func (x *RecordContext_Workload) ProtoReflect() protoreflect.Message {
	mi := &file_context_platform_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordContext_Workload.ProtoReflect.Descriptor instead.
func (*RecordContext_Workload) Descriptor() ([]byte, []int) {
	return file_context_platform_proto_rawDescGZIP(), []int{0, 0}
}

func (x *RecordContext_Workload) GetRemoteId() string {
	if x != nil {
		return x.RemoteId
	}
	return ""
}

// Ancestry relationships are added to records when a "base image"
// relationship has been uncovered.  This context is added to BOTH
// records with their respective roles.  The base image will get
// the Role BASE, and the derivative image will get the Role DERIVED.
type RecordContext_Ancestry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role RecordContext_Ancestry_Role `protobuf:"varint,1,opt,name=role,proto3,enum=chainguard.platform.tenant.RecordContext_Ancestry_Role" json:"role,omitempty"`
	// image_id holds the digest of the related image, which can be used
	// to efficiently retrieve it's record.
	ImageId string `protobuf:"bytes,2,opt,name=image_id,json=imageId,proto3" json:"image_id,omitempty"`
}

func (x *RecordContext_Ancestry) Reset() {
	*x = RecordContext_Ancestry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_platform_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordContext_Ancestry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordContext_Ancestry) ProtoMessage() {}

func (x *RecordContext_Ancestry) ProtoReflect() protoreflect.Message {
	mi := &file_context_platform_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordContext_Ancestry.ProtoReflect.Descriptor instead.
func (*RecordContext_Ancestry) Descriptor() ([]byte, []int) {
	return file_context_platform_proto_rawDescGZIP(), []int{0, 1}
}

func (x *RecordContext_Ancestry) GetRole() RecordContext_Ancestry_Role {
	if x != nil {
		return x.Role
	}
	return RecordContext_Ancestry_UNKNOWN
}

func (x *RecordContext_Ancestry) GetImageId() string {
	if x != nil {
		return x.ImageId
	}
	return ""
}

// Variant relationships are added to records when we find an "index"
// containing multiple different variations (typically os/arch) of the
// same logical image.  These are referred to as "OCI Image Index",
// "Docker Manifest List", and occasionally "fat images".  This context
// is added to ALL records including the INDEX and all VARIANTs of that
// index.  The INDEX will typically contain N contexts carrying the Role
// INDEX, the id of the VARIANT's record, and the version information
// that discriminates that VARIANT from other VARIANTs.  The VARIANT
// will typically (but not always!) contain 1 context carrying the Role
// VARIANT, the id of the INDEX's record, and the version information
// that discriminates it among the other VARIANTs in the INDEX.
type RecordContext_Variant struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role RecordContext_Variant_Role `protobuf:"varint,1,opt,name=role,proto3,enum=chainguard.platform.tenant.RecordContext_Variant_Role" json:"role,omitempty"`
	// image_id holds the digest of the related image, which can be used
	// to efficiently retrieve it's record.
	ImageId string `protobuf:"bytes,2,opt,name=image_id,json=imageId,proto3" json:"image_id,omitempty"`
	// The version information distinguishing this variant
	// from other possible variants of the index.
	Version string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *RecordContext_Variant) Reset() {
	*x = RecordContext_Variant{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_platform_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordContext_Variant) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordContext_Variant) ProtoMessage() {}

func (x *RecordContext_Variant) ProtoReflect() protoreflect.Message {
	mi := &file_context_platform_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordContext_Variant.ProtoReflect.Descriptor instead.
func (*RecordContext_Variant) Descriptor() ([]byte, []int) {
	return file_context_platform_proto_rawDescGZIP(), []int{0, 2}
}

func (x *RecordContext_Variant) GetRole() RecordContext_Variant_Role {
	if x != nil {
		return x.Role
	}
	return RecordContext_Variant_UNKNOWN
}

func (x *RecordContext_Variant) GetImageId() string {
	if x != nil {
		return x.ImageId
	}
	return ""
}

func (x *RecordContext_Variant) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

var File_context_platform_proto protoreflect.FileDescriptor

var file_context_platform_proto_rawDesc = []byte{
	0x0a, 0x16, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67,
	0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x74, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x26, 0x70,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76,
	0x31, 0x2f, 0x75, 0x69, 0x64, 0x70, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8d, 0x06, 0x0a, 0x0d, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x37, 0x0a,
	0x09, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x73, 0x65, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x6c, 0x61,
	0x73, 0x74, 0x53, 0x65, 0x65, 0x6e, 0x12, 0x50, 0x0a, 0x08, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f,
	0x61, 0x64, 0x18, 0x64, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e,
	0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x74,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f, 0x6e, 0x74,
	0x65, 0x78, 0x74, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x48, 0x00, 0x52, 0x08,
	0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x50, 0x0a, 0x08, 0x61, 0x6e, 0x63, 0x65,
	0x73, 0x74, 0x72, 0x79, 0x18, 0x65, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x2e, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x41, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x72, 0x79, 0x48, 0x00,
	0x52, 0x08, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x72, 0x79, 0x12, 0x4d, 0x0a, 0x07, 0x76, 0x61,
	0x72, 0x69, 0x61, 0x6e, 0x74, 0x18, 0x66, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x2e, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x48, 0x00,
	0x52, 0x07, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x1a, 0x27, 0x0a, 0x08, 0x57, 0x6f, 0x72,
	0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65,
	0x49, 0x64, 0x1a, 0x9e, 0x01, 0x0a, 0x08, 0x41, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x72, 0x79, 0x12,
	0x4b, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x37, 0x2e,
	0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x2e, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x41, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x72,
	0x79, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x19, 0x0a, 0x08,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x49, 0x64, 0x22, 0x2a, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12,
	0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04,
	0x42, 0x41, 0x53, 0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x45, 0x52, 0x49, 0x56, 0x45,
	0x44, 0x10, 0x02, 0x1a, 0xb7, 0x01, 0x0a, 0x07, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x12,
	0x4a, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x36, 0x2e,
	0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x2e, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74,
	0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x22, 0x2b, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e,
	0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x10, 0x01,
	0x12, 0x0b, 0x0a, 0x07, 0x56, 0x41, 0x52, 0x49, 0x41, 0x4e, 0x54, 0x10, 0x02, 0x42, 0x06, 0x0a,
	0x04, 0x6b, 0x69, 0x6e, 0x64, 0x22, 0x54, 0x0a, 0x11, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x3f, 0x0a, 0x05, 0x69, 0x74,
	0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e,
	0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x78, 0x74, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x8b, 0x03, 0x0a, 0x13,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x46, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x12, 0x3a, 0x0a, 0x04, 0x75, 0x69, 0x64, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x26, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x55,
	0x49, 0x44, 0x50, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x04, 0x75, 0x69, 0x64, 0x70, 0x12,
	0x3d, 0x0a, 0x0c, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x73, 0x69, 0x6e, 0x63, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x0b, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x69, 0x6e, 0x63, 0x65, 0x12, 0x50,
	0x0a, 0x08, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x64, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x32, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c,
	0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x57, 0x6f, 0x72, 0x6b,
	0x6c, 0x6f, 0x61, 0x64, 0x48, 0x00, 0x52, 0x08, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64,
	0x12, 0x50, 0x0a, 0x08, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x72, 0x79, 0x18, 0x65, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x32, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x41, 0x6e,
	0x63, 0x65, 0x73, 0x74, 0x72, 0x79, 0x48, 0x00, 0x52, 0x08, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x74,
	0x72, 0x79, 0x12, 0x4d, 0x0a, 0x07, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x18, 0x66, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64,
	0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x56,
	0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x07, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e,
	0x74, 0x42, 0x06, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x32, 0x88, 0x01, 0x0a, 0x0e, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x73, 0x12, 0x76, 0x0a, 0x04,
	0x4c, 0x69, 0x73, 0x74, 0x12, 0x2f, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72,
	0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x74, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x46,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x1a, 0x2d, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61,
	0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x74, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74,
	0x4c, 0x69, 0x73, 0x74, 0x22, 0x0e, 0x8a, 0xaf, 0xa8, 0xd2, 0x05, 0x08, 0x12, 0x06, 0x0a, 0x02,
	0xe5, 0x04, 0x10, 0x01, 0x42, 0x78, 0x0a, 0x25, 0x64, 0x65, 0x76, 0x2e, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x2e, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x42, 0x20, 0x50,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x2b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x64, 0x65,
	0x76, 0x2f, 0x73, 0x64, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_context_platform_proto_rawDescOnce sync.Once
	file_context_platform_proto_rawDescData = file_context_platform_proto_rawDesc
)

func file_context_platform_proto_rawDescGZIP() []byte {
	file_context_platform_proto_rawDescOnce.Do(func() {
		file_context_platform_proto_rawDescData = protoimpl.X.CompressGZIP(file_context_platform_proto_rawDescData)
	})
	return file_context_platform_proto_rawDescData
}

var file_context_platform_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_context_platform_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_context_platform_proto_goTypes = []interface{}{
	(RecordContext_Ancestry_Role)(0), // 0: chainguard.platform.tenant.RecordContext.Ancestry.Role
	(RecordContext_Variant_Role)(0),  // 1: chainguard.platform.tenant.RecordContext.Variant.Role
	(*RecordContext)(nil),            // 2: chainguard.platform.tenant.RecordContext
	(*RecordContextList)(nil),        // 3: chainguard.platform.tenant.RecordContextList
	(*RecordContextFilter)(nil),      // 4: chainguard.platform.tenant.RecordContextFilter
	(*RecordContext_Workload)(nil),   // 5: chainguard.platform.tenant.RecordContext.Workload
	(*RecordContext_Ancestry)(nil),   // 6: chainguard.platform.tenant.RecordContext.Ancestry
	(*RecordContext_Variant)(nil),    // 7: chainguard.platform.tenant.RecordContext.Variant
	(*timestamppb.Timestamp)(nil),    // 8: google.protobuf.Timestamp
	(*v1.UIDPFilter)(nil),            // 9: chainguard.platform.common.UIDPFilter
}
var file_context_platform_proto_depIdxs = []int32{
	8,  // 0: chainguard.platform.tenant.RecordContext.last_seen:type_name -> google.protobuf.Timestamp
	5,  // 1: chainguard.platform.tenant.RecordContext.workload:type_name -> chainguard.platform.tenant.RecordContext.Workload
	6,  // 2: chainguard.platform.tenant.RecordContext.ancestry:type_name -> chainguard.platform.tenant.RecordContext.Ancestry
	7,  // 3: chainguard.platform.tenant.RecordContext.variant:type_name -> chainguard.platform.tenant.RecordContext.Variant
	2,  // 4: chainguard.platform.tenant.RecordContextList.items:type_name -> chainguard.platform.tenant.RecordContext
	9,  // 5: chainguard.platform.tenant.RecordContextFilter.uidp:type_name -> chainguard.platform.common.UIDPFilter
	8,  // 6: chainguard.platform.tenant.RecordContextFilter.active_since:type_name -> google.protobuf.Timestamp
	5,  // 7: chainguard.platform.tenant.RecordContextFilter.workload:type_name -> chainguard.platform.tenant.RecordContext.Workload
	6,  // 8: chainguard.platform.tenant.RecordContextFilter.ancestry:type_name -> chainguard.platform.tenant.RecordContext.Ancestry
	7,  // 9: chainguard.platform.tenant.RecordContextFilter.variant:type_name -> chainguard.platform.tenant.RecordContext.Variant
	0,  // 10: chainguard.platform.tenant.RecordContext.Ancestry.role:type_name -> chainguard.platform.tenant.RecordContext.Ancestry.Role
	1,  // 11: chainguard.platform.tenant.RecordContext.Variant.role:type_name -> chainguard.platform.tenant.RecordContext.Variant.Role
	4,  // 12: chainguard.platform.tenant.RecordContexts.List:input_type -> chainguard.platform.tenant.RecordContextFilter
	3,  // 13: chainguard.platform.tenant.RecordContexts.List:output_type -> chainguard.platform.tenant.RecordContextList
	13, // [13:14] is the sub-list for method output_type
	12, // [12:13] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_context_platform_proto_init() }
func file_context_platform_proto_init() {
	if File_context_platform_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_context_platform_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordContext); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_context_platform_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordContextList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_context_platform_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordContextFilter); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_context_platform_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordContext_Workload); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_context_platform_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordContext_Ancestry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_context_platform_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordContext_Variant); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_context_platform_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*RecordContext_Workload_)(nil),
		(*RecordContext_Ancestry_)(nil),
		(*RecordContext_Variant_)(nil),
	}
	file_context_platform_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*RecordContextFilter_Workload)(nil),
		(*RecordContextFilter_Ancestry)(nil),
		(*RecordContextFilter_Variant)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_context_platform_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_context_platform_proto_goTypes,
		DependencyIndexes: file_context_platform_proto_depIdxs,
		EnumInfos:         file_context_platform_proto_enumTypes,
		MessageInfos:      file_context_platform_proto_msgTypes,
	}.Build()
	File_context_platform_proto = out.File
	file_context_platform_proto_rawDesc = nil
	file_context_platform_proto_goTypes = nil
	file_context_platform_proto_depIdxs = nil
}
