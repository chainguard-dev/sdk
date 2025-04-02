// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.29.3
// source: registry_entitlements.platform.proto

package v1

import (
	_ "chainguard.dev/sdk/proto/annotations"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type Entitlement_Type int32

const (
	Entitlement_UNKNOWN    Entitlement_Type = 0
	Entitlement_TRIAL      Entitlement_Type = 1
	Entitlement_PRODUCTION Entitlement_Type = 2
)

// Enum value maps for Entitlement_Type.
var (
	Entitlement_Type_name = map[int32]string{
		0: "UNKNOWN",
		1: "TRIAL",
		2: "PRODUCTION",
	}
	Entitlement_Type_value = map[string]int32{
		"UNKNOWN":    0,
		"TRIAL":      1,
		"PRODUCTION": 2,
	}
)

func (x Entitlement_Type) Enum() *Entitlement_Type {
	p := new(Entitlement_Type)
	*p = x
	return p
}

func (x Entitlement_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Entitlement_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_registry_entitlements_platform_proto_enumTypes[0].Descriptor()
}

func (Entitlement_Type) Type() protoreflect.EnumType {
	return &file_registry_entitlements_platform_proto_enumTypes[0]
}

func (x Entitlement_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Entitlement_Type.Descriptor instead.
func (Entitlement_Type) EnumDescriptor() ([]byte, []int) {
	return file_registry_entitlements_platform_proto_rawDescGZIP(), []int{0, 0}
}

// Entitlement contains information about what an organization is entitled to.
type Entitlement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ExternalId     string                 `protobuf:"bytes,2,opt,name=external_id,json=externalId,proto3" json:"external_id,omitempty"`
	CreateTime     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	UpdateTime     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
	ExpirationTime *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=expiration_time,json=expirationTime,proto3" json:"expiration_time,omitempty"`
	Type           Entitlement_Type       `protobuf:"varint,6,opt,name=type,proto3,enum=chainguard.platform.registry.Entitlement_Type" json:"type,omitempty"`
	// Keys can't be enum types, but string should match CatalogTier.
	Quota map[string]*ImageQuota `protobuf:"bytes,7,rep,name=quota,proto3" json:"quota,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Entitlement) Reset() {
	*x = Entitlement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_entitlements_platform_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entitlement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entitlement) ProtoMessage() {}

func (x *Entitlement) ProtoReflect() protoreflect.Message {
	mi := &file_registry_entitlements_platform_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entitlement.ProtoReflect.Descriptor instead.
func (*Entitlement) Descriptor() ([]byte, []int) {
	return file_registry_entitlements_platform_proto_rawDescGZIP(), []int{0}
}

func (x *Entitlement) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Entitlement) GetExternalId() string {
	if x != nil {
		return x.ExternalId
	}
	return ""
}

func (x *Entitlement) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

func (x *Entitlement) GetUpdateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

func (x *Entitlement) GetExpirationTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpirationTime
	}
	return nil
}

func (x *Entitlement) GetType() Entitlement_Type {
	if x != nil {
		return x.Type
	}
	return Entitlement_UNKNOWN
}

func (x *Entitlement) GetQuota() map[string]*ImageQuota {
	if x != nil {
		return x.Quota
	}
	return nil
}

type ImageQuota struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Current int32 `protobuf:"varint,1,opt,name=current,proto3" json:"current,omitempty"`
	Max     int32 `protobuf:"varint,2,opt,name=max,proto3" json:"max,omitempty"`
}

func (x *ImageQuota) Reset() {
	*x = ImageQuota{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_entitlements_platform_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImageQuota) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImageQuota) ProtoMessage() {}

func (x *ImageQuota) ProtoReflect() protoreflect.Message {
	mi := &file_registry_entitlements_platform_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImageQuota.ProtoReflect.Descriptor instead.
func (*ImageQuota) Descriptor() ([]byte, []int) {
	return file_registry_entitlements_platform_proto_rawDescGZIP(), []int{1}
}

func (x *ImageQuota) GetCurrent() int32 {
	if x != nil {
		return x.Current
	}
	return 0
}

func (x *ImageQuota) GetMax() int32 {
	if x != nil {
		return x.Max
	}
	return 0
}

type EntitlementFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Parent string `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
}

func (x *EntitlementFilter) Reset() {
	*x = EntitlementFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_entitlements_platform_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitlementFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitlementFilter) ProtoMessage() {}

func (x *EntitlementFilter) ProtoReflect() protoreflect.Message {
	mi := &file_registry_entitlements_platform_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitlementFilter.ProtoReflect.Descriptor instead.
func (*EntitlementFilter) Descriptor() ([]byte, []int) {
	return file_registry_entitlements_platform_proto_rawDescGZIP(), []int{2}
}

func (x *EntitlementFilter) GetParent() string {
	if x != nil {
		return x.Parent
	}
	return ""
}

type EntitlementList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*Entitlement `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *EntitlementList) Reset() {
	*x = EntitlementList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_entitlements_platform_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitlementList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitlementList) ProtoMessage() {}

func (x *EntitlementList) ProtoReflect() protoreflect.Message {
	mi := &file_registry_entitlements_platform_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitlementList.ProtoReflect.Descriptor instead.
func (*EntitlementList) Descriptor() ([]byte, []int) {
	return file_registry_entitlements_platform_proto_rawDescGZIP(), []int{3}
}

func (x *EntitlementList) GetItems() []*Entitlement {
	if x != nil {
		return x.Items
	}
	return nil
}

type EntitlementImage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The image repository UID.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The image repository catalog tier the image is associated with.
	Tier CatalogTier `protobuf:"varint,2,opt,name=tier,proto3,enum=chainguard.platform.registry.CatalogTier" json:"tier,omitempty"`
	// Human-readable image name corresponding to id.
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *EntitlementImage) Reset() {
	*x = EntitlementImage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_entitlements_platform_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitlementImage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitlementImage) ProtoMessage() {}

func (x *EntitlementImage) ProtoReflect() protoreflect.Message {
	mi := &file_registry_entitlements_platform_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitlementImage.ProtoReflect.Descriptor instead.
func (*EntitlementImage) Descriptor() ([]byte, []int) {
	return file_registry_entitlements_platform_proto_rawDescGZIP(), []int{4}
}

func (x *EntitlementImage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *EntitlementImage) GetTier() CatalogTier {
	if x != nil {
		return x.Tier
	}
	return CatalogTier_UNKNOWN
}

func (x *EntitlementImage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type EntitlementImagesFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Parent string `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
}

func (x *EntitlementImagesFilter) Reset() {
	*x = EntitlementImagesFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_entitlements_platform_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitlementImagesFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitlementImagesFilter) ProtoMessage() {}

func (x *EntitlementImagesFilter) ProtoReflect() protoreflect.Message {
	mi := &file_registry_entitlements_platform_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitlementImagesFilter.ProtoReflect.Descriptor instead.
func (*EntitlementImagesFilter) Descriptor() ([]byte, []int) {
	return file_registry_entitlements_platform_proto_rawDescGZIP(), []int{5}
}

func (x *EntitlementImagesFilter) GetParent() string {
	if x != nil {
		return x.Parent
	}
	return ""
}

type EntitlementImagesList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Images []*EntitlementImage `protobuf:"bytes,1,rep,name=images,proto3" json:"images,omitempty"`
}

func (x *EntitlementImagesList) Reset() {
	*x = EntitlementImagesList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_entitlements_platform_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitlementImagesList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitlementImagesList) ProtoMessage() {}

func (x *EntitlementImagesList) ProtoReflect() protoreflect.Message {
	mi := &file_registry_entitlements_platform_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitlementImagesList.ProtoReflect.Descriptor instead.
func (*EntitlementImagesList) Descriptor() ([]byte, []int) {
	return file_registry_entitlements_platform_proto_rawDescGZIP(), []int{6}
}

func (x *EntitlementImagesList) GetImages() []*EntitlementImage {
	if x != nil {
		return x.Images
	}
	return nil
}

type EntitlementSummaryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Parent string `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
}

func (x *EntitlementSummaryRequest) Reset() {
	*x = EntitlementSummaryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_entitlements_platform_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitlementSummaryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitlementSummaryRequest) ProtoMessage() {}

func (x *EntitlementSummaryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_registry_entitlements_platform_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitlementSummaryRequest.ProtoReflect.Descriptor instead.
func (*EntitlementSummaryRequest) Descriptor() ([]byte, []int) {
	return file_registry_entitlements_platform_proto_rawDescGZIP(), []int{7}
}

func (x *EntitlementSummaryRequest) GetParent() string {
	if x != nil {
		return x.Parent
	}
	return ""
}

type EntitlementSummaryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Keys can't be enum types, but string should match CatalogTier.
	Quota map[string]*ImageQuota `protobuf:"bytes,1,rep,name=quota,proto3" json:"quota,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *EntitlementSummaryResponse) Reset() {
	*x = EntitlementSummaryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_entitlements_platform_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitlementSummaryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitlementSummaryResponse) ProtoMessage() {}

func (x *EntitlementSummaryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_registry_entitlements_platform_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitlementSummaryResponse.ProtoReflect.Descriptor instead.
func (*EntitlementSummaryResponse) Descriptor() ([]byte, []int) {
	return file_registry_entitlements_platform_proto_rawDescGZIP(), []int{8}
}

func (x *EntitlementSummaryResponse) GetQuota() map[string]*ImageQuota {
	if x != nil {
		return x.Quota
	}
	return nil
}

var File_registry_entitlements_platform_proto protoreflect.FileDescriptor

var file_registry_entitlements_platform_proto_rawDesc = []byte{
	0x0a, 0x24, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x5f, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1c, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61,
	0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x79, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62,
	0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa1, 0x04, 0x0a, 0x0b, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x78, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x3b, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x43, 0x0a, 0x0f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0e, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x42, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x2e, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64,
	0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x79, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x4a, 0x0a, 0x05, 0x71, 0x75, 0x6f,
	0x74, 0x61, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e,
	0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x51, 0x75, 0x6f, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05,
	0x71, 0x75, 0x6f, 0x74, 0x61, 0x1a, 0x62, 0x0a, 0x0a, 0x51, 0x75, 0x6f, 0x74, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x3e, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72,
	0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x79, 0x2e, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x51, 0x75, 0x6f, 0x74, 0x61, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x2e, 0x0a, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x09,
	0x0a, 0x05, 0x54, 0x52, 0x49, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x52, 0x4f,
	0x44, 0x55, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x02, 0x22, 0x38, 0x0a, 0x0a, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x51, 0x75, 0x6f, 0x74, 0x61, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x6d, 0x61, 0x78, 0x22, 0x33, 0x0a, 0x11, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x65,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0x90, 0xaf, 0xa8, 0xd2, 0x05, 0x01,
	0x52, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x22, 0x52, 0x0a, 0x0f, 0x45, 0x6e, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x3f, 0x0a, 0x05, 0x69,
	0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x7b, 0x0a, 0x10,
	0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x3d, 0x0a, 0x04, 0x74, 0x69, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x29,
	0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x43, 0x61,
	0x74, 0x61, 0x6c, 0x6f, 0x67, 0x54, 0x69, 0x65, 0x72, 0x52, 0x04, 0x74, 0x69, 0x65, 0x72, 0x12,
	0x18, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xe2,
	0x41, 0x01, 0x03, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x39, 0x0a, 0x17, 0x45, 0x6e, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x46, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0x90, 0xaf, 0xa8, 0xd2, 0x05, 0x01, 0x52, 0x06, 0x70, 0x61,
	0x72, 0x65, 0x6e, 0x74, 0x22, 0x5f, 0x0a, 0x15, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x46, 0x0a,
	0x06, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e,
	0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x45, 0x6e, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x06, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x73, 0x22, 0x3b, 0x0a, 0x19, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1e, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x06, 0x90, 0xaf, 0xa8, 0xd2, 0x05, 0x01, 0x52, 0x06, 0x70, 0x61, 0x72, 0x65,
	0x6e, 0x74, 0x22, 0xdb, 0x01, 0x0a, 0x1a, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x59, 0x0a, 0x05, 0x71, 0x75, 0x6f, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x43, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c,
	0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e,
	0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x75, 0x6d, 0x6d, 0x61,
	0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x51, 0x75, 0x6f, 0x74, 0x61,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x71, 0x75, 0x6f, 0x74, 0x61, 0x1a, 0x62, 0x0a, 0x0a,
	0x51, 0x75, 0x6f, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x3e, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x51, 0x75, 0x6f, 0x74, 0x61, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x32, 0xc8, 0x04, 0x0a, 0x0c, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x12, 0xad, 0x01, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x2f, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75,
	0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x1a, 0x2d, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67,
	0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x39, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x27, 0x12, 0x25,
	0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x7b, 0x70, 0x61, 0x72, 0x65, 0x6e,
	0x74, 0x3d, 0x2a, 0x2a, 0x7d, 0x8a, 0xaf, 0xa8, 0xd2, 0x05, 0x06, 0x12, 0x04, 0x0a, 0x02, 0x90,
	0x0d, 0x12, 0xc5, 0x01, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x12, 0x35, 0x2e, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x1a, 0x33, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x40, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2e, 0x12,
	0x2c, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x7b, 0x70, 0x61, 0x72, 0x65,
	0x6e, 0x74, 0x3d, 0x2a, 0x2a, 0x7d, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x8a, 0xaf, 0xa8,
	0xd2, 0x05, 0x06, 0x12, 0x04, 0x0a, 0x02, 0x90, 0x0d, 0x12, 0xbf, 0x01, 0x0a, 0x07, 0x53, 0x75,
	0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x37, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61,
	0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x79, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x38,
	0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x41, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2f,
	0x12, 0x2d, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x7b, 0x70, 0x61, 0x72,
	0x65, 0x6e, 0x74, 0x3d, 0x2a, 0x2a, 0x7d, 0x3a, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x8a,
	0xaf, 0xa8, 0xd2, 0x05, 0x06, 0x12, 0x04, 0x0a, 0x02, 0x90, 0x0d, 0x42, 0x71, 0x0a, 0x27, 0x64,
	0x65, 0x76, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x73, 0x64,
	0x6b, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x42, 0x15, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x2d, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x64, 0x65, 0x76, 0x2f,
	0x73, 0x64, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_registry_entitlements_platform_proto_rawDescOnce sync.Once
	file_registry_entitlements_platform_proto_rawDescData = file_registry_entitlements_platform_proto_rawDesc
)

func file_registry_entitlements_platform_proto_rawDescGZIP() []byte {
	file_registry_entitlements_platform_proto_rawDescOnce.Do(func() {
		file_registry_entitlements_platform_proto_rawDescData = protoimpl.X.CompressGZIP(file_registry_entitlements_platform_proto_rawDescData)
	})
	return file_registry_entitlements_platform_proto_rawDescData
}

var file_registry_entitlements_platform_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_registry_entitlements_platform_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_registry_entitlements_platform_proto_goTypes = []any{
	(Entitlement_Type)(0),              // 0: chainguard.platform.registry.Entitlement.Type
	(*Entitlement)(nil),                // 1: chainguard.platform.registry.Entitlement
	(*ImageQuota)(nil),                 // 2: chainguard.platform.registry.ImageQuota
	(*EntitlementFilter)(nil),          // 3: chainguard.platform.registry.EntitlementFilter
	(*EntitlementList)(nil),            // 4: chainguard.platform.registry.EntitlementList
	(*EntitlementImage)(nil),           // 5: chainguard.platform.registry.EntitlementImage
	(*EntitlementImagesFilter)(nil),    // 6: chainguard.platform.registry.EntitlementImagesFilter
	(*EntitlementImagesList)(nil),      // 7: chainguard.platform.registry.EntitlementImagesList
	(*EntitlementSummaryRequest)(nil),  // 8: chainguard.platform.registry.EntitlementSummaryRequest
	(*EntitlementSummaryResponse)(nil), // 9: chainguard.platform.registry.EntitlementSummaryResponse
	nil,                                // 10: chainguard.platform.registry.Entitlement.QuotaEntry
	nil,                                // 11: chainguard.platform.registry.EntitlementSummaryResponse.QuotaEntry
	(*timestamppb.Timestamp)(nil),      // 12: google.protobuf.Timestamp
	(CatalogTier)(0),                   // 13: chainguard.platform.registry.CatalogTier
}
var file_registry_entitlements_platform_proto_depIdxs = []int32{
	12, // 0: chainguard.platform.registry.Entitlement.create_time:type_name -> google.protobuf.Timestamp
	12, // 1: chainguard.platform.registry.Entitlement.update_time:type_name -> google.protobuf.Timestamp
	12, // 2: chainguard.platform.registry.Entitlement.expiration_time:type_name -> google.protobuf.Timestamp
	0,  // 3: chainguard.platform.registry.Entitlement.type:type_name -> chainguard.platform.registry.Entitlement.Type
	10, // 4: chainguard.platform.registry.Entitlement.quota:type_name -> chainguard.platform.registry.Entitlement.QuotaEntry
	1,  // 5: chainguard.platform.registry.EntitlementList.items:type_name -> chainguard.platform.registry.Entitlement
	13, // 6: chainguard.platform.registry.EntitlementImage.tier:type_name -> chainguard.platform.registry.CatalogTier
	5,  // 7: chainguard.platform.registry.EntitlementImagesList.images:type_name -> chainguard.platform.registry.EntitlementImage
	11, // 8: chainguard.platform.registry.EntitlementSummaryResponse.quota:type_name -> chainguard.platform.registry.EntitlementSummaryResponse.QuotaEntry
	2,  // 9: chainguard.platform.registry.Entitlement.QuotaEntry.value:type_name -> chainguard.platform.registry.ImageQuota
	2,  // 10: chainguard.platform.registry.EntitlementSummaryResponse.QuotaEntry.value:type_name -> chainguard.platform.registry.ImageQuota
	3,  // 11: chainguard.platform.registry.Entitlements.ListEntitlements:input_type -> chainguard.platform.registry.EntitlementFilter
	6,  // 12: chainguard.platform.registry.Entitlements.ListEntitlementImages:input_type -> chainguard.platform.registry.EntitlementImagesFilter
	8,  // 13: chainguard.platform.registry.Entitlements.Summary:input_type -> chainguard.platform.registry.EntitlementSummaryRequest
	4,  // 14: chainguard.platform.registry.Entitlements.ListEntitlements:output_type -> chainguard.platform.registry.EntitlementList
	7,  // 15: chainguard.platform.registry.Entitlements.ListEntitlementImages:output_type -> chainguard.platform.registry.EntitlementImagesList
	9,  // 16: chainguard.platform.registry.Entitlements.Summary:output_type -> chainguard.platform.registry.EntitlementSummaryResponse
	14, // [14:17] is the sub-list for method output_type
	11, // [11:14] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_registry_entitlements_platform_proto_init() }
func file_registry_entitlements_platform_proto_init() {
	if File_registry_entitlements_platform_proto != nil {
		return
	}
	file_registry_platform_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_registry_entitlements_platform_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Entitlement); i {
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
		file_registry_entitlements_platform_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ImageQuota); i {
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
		file_registry_entitlements_platform_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*EntitlementFilter); i {
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
		file_registry_entitlements_platform_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*EntitlementList); i {
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
		file_registry_entitlements_platform_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*EntitlementImage); i {
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
		file_registry_entitlements_platform_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*EntitlementImagesFilter); i {
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
		file_registry_entitlements_platform_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*EntitlementImagesList); i {
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
		file_registry_entitlements_platform_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*EntitlementSummaryRequest); i {
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
		file_registry_entitlements_platform_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*EntitlementSummaryResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_registry_entitlements_platform_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_registry_entitlements_platform_proto_goTypes,
		DependencyIndexes: file_registry_entitlements_platform_proto_depIdxs,
		EnumInfos:         file_registry_entitlements_platform_proto_enumTypes,
		MessageInfos:      file_registry_entitlements_platform_proto_msgTypes,
	}.Build()
	File_registry_entitlements_platform_proto = out.File
	file_registry_entitlements_platform_proto_rawDesc = nil
	file_registry_entitlements_platform_proto_goTypes = nil
	file_registry_entitlements_platform_proto_depIdxs = nil
}
