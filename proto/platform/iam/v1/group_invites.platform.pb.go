// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.29.3
// source: group_invites.platform.proto

package v1

import (
	_ "chainguard.dev/sdk/proto/annotations"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

type RegistrationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Kind:
	//
	//	*RegistrationRequest_Human_
	//	*RegistrationRequest_Cluster_
	Kind isRegistrationRequest_Kind `protobuf_oneof:"kind"`
}

func (x *RegistrationRequest) Reset() {
	*x = RegistrationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_invites_platform_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistrationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrationRequest) ProtoMessage() {}

func (x *RegistrationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_group_invites_platform_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistrationRequest.ProtoReflect.Descriptor instead.
func (*RegistrationRequest) Descriptor() ([]byte, []int) {
	return file_group_invites_platform_proto_rawDescGZIP(), []int{0}
}

func (m *RegistrationRequest) GetKind() isRegistrationRequest_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (x *RegistrationRequest) GetHuman() *RegistrationRequest_Human {
	if x, ok := x.GetKind().(*RegistrationRequest_Human_); ok {
		return x.Human
	}
	return nil
}

func (x *RegistrationRequest) GetCluster() *RegistrationRequest_Cluster {
	if x, ok := x.GetKind().(*RegistrationRequest_Cluster_); ok {
		return x.Cluster
	}
	return nil
}

type isRegistrationRequest_Kind interface {
	isRegistrationRequest_Kind()
}

type RegistrationRequest_Human_ struct {
	Human *RegistrationRequest_Human `protobuf:"bytes,1,opt,name=human,proto3,oneof"`
}

type RegistrationRequest_Cluster_ struct {
	Cluster *RegistrationRequest_Cluster `protobuf:"bytes,2,opt,name=cluster,proto3,oneof"`
}

func (*RegistrationRequest_Human_) isRegistrationRequest_Kind() {}

func (*RegistrationRequest_Cluster_) isRegistrationRequest_Kind() {}

type GroupInvite struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id, The group UIDP under which this invite resides.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// expiration, timestamp this invite becomes no longer valid.
	Expiration *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=expiration,proto3" json:"expiration,omitempty"`
	// key_id is used to identify the verification key for this code.
	KeyId string `protobuf:"bytes,3,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
	// role is the role the invited identity will be role-bound to the group with.
	Role *Role `protobuf:"bytes,4,opt,name=role,proto3" json:"role,omitempty"`
	// code is the json-encoded authentication code.
	Code string `protobuf:"bytes,5,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *GroupInvite) Reset() {
	*x = GroupInvite{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_invites_platform_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupInvite) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupInvite) ProtoMessage() {}

func (x *GroupInvite) ProtoReflect() protoreflect.Message {
	mi := &file_group_invites_platform_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupInvite.ProtoReflect.Descriptor instead.
func (*GroupInvite) Descriptor() ([]byte, []int) {
	return file_group_invites_platform_proto_rawDescGZIP(), []int{1}
}

func (x *GroupInvite) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GroupInvite) GetExpiration() *timestamppb.Timestamp {
	if x != nil {
		return x.Expiration
	}
	return nil
}

func (x *GroupInvite) GetKeyId() string {
	if x != nil {
		return x.KeyId
	}
	return ""
}

func (x *GroupInvite) GetRole() *Role {
	if x != nil {
		return x.Role
	}
	return nil
}

func (x *GroupInvite) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type StoredGroupInvite struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id, The group UIDP under which this invite resides.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// expiration, timestamp this invite becomes no longer valid.
	Expiration *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=expiration,proto3" json:"expiration,omitempty"`
	// key_id is used to identify the verification key for this code.
	KeyId string `protobuf:"bytes,3,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
	// role is the role the invited identity will be role-bound to the group with.
	Role *Role `protobuf:"bytes,4,opt,name=role,proto3" json:"role,omitempty"`
	// email is the email address that is allowed to accept this invite code. If blank,
	// anyone with the invite code an accept.
	Email string `protobuf:"bytes,5,opt,name=email,proto3" json:"email,omitempty"`
	// created_at is the timestamp for when the invite was created.
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// single_use indicates whether or not the invite will be deleted after a user joins the group.
	SingleUse bool `protobuf:"varint,7,opt,name=single_use,json=singleUse,proto3" json:"single_use,omitempty"`
}

func (x *StoredGroupInvite) Reset() {
	*x = StoredGroupInvite{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_invites_platform_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoredGroupInvite) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoredGroupInvite) ProtoMessage() {}

func (x *StoredGroupInvite) ProtoReflect() protoreflect.Message {
	mi := &file_group_invites_platform_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoredGroupInvite.ProtoReflect.Descriptor instead.
func (*StoredGroupInvite) Descriptor() ([]byte, []int) {
	return file_group_invites_platform_proto_rawDescGZIP(), []int{2}
}

func (x *StoredGroupInvite) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *StoredGroupInvite) GetExpiration() *timestamppb.Timestamp {
	if x != nil {
		return x.Expiration
	}
	return nil
}

func (x *StoredGroupInvite) GetKeyId() string {
	if x != nil {
		return x.KeyId
	}
	return ""
}

func (x *StoredGroupInvite) GetRole() *Role {
	if x != nil {
		return x.Role
	}
	return nil
}

func (x *StoredGroupInvite) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *StoredGroupInvite) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *StoredGroupInvite) GetSingleUse() bool {
	if x != nil {
		return x.SingleUse
	}
	return false
}

type GroupInviteList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*StoredGroupInvite `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *GroupInviteList) Reset() {
	*x = GroupInviteList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_invites_platform_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupInviteList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupInviteList) ProtoMessage() {}

func (x *GroupInviteList) ProtoReflect() protoreflect.Message {
	mi := &file_group_invites_platform_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupInviteList.ProtoReflect.Descriptor instead.
func (*GroupInviteList) Descriptor() ([]byte, []int) {
	return file_group_invites_platform_proto_rawDescGZIP(), []int{3}
}

func (x *GroupInviteList) GetItems() []*StoredGroupInvite {
	if x != nil {
		return x.Items
	}
	return nil
}

type GroupInviteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// group, The Group UIDP path under which the new group Invite targets.
	Group string `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
	// expiration, timestamp this invite becomes no longer valid.
	Ttl *durationpb.Duration `protobuf:"bytes,2,opt,name=ttl,proto3" json:"ttl,omitempty"`
	// role is the Role UIDP the invited identity will be role-bound to the group with.
	Role string `protobuf:"bytes,3,opt,name=role,proto3" json:"role,omitempty"`
	// email is the exact email address that may accept this invite code, if specified.
	Email string `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	// if single_use is set to true, then the invite will be deleted after a user joins the group.
	SingleUse bool `protobuf:"varint,5,opt,name=single_use,json=singleUse,proto3" json:"single_use,omitempty"`
}

func (x *GroupInviteRequest) Reset() {
	*x = GroupInviteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_invites_platform_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupInviteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupInviteRequest) ProtoMessage() {}

func (x *GroupInviteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_group_invites_platform_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupInviteRequest.ProtoReflect.Descriptor instead.
func (*GroupInviteRequest) Descriptor() ([]byte, []int) {
	return file_group_invites_platform_proto_rawDescGZIP(), []int{4}
}

func (x *GroupInviteRequest) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

func (x *GroupInviteRequest) GetTtl() *durationpb.Duration {
	if x != nil {
		return x.Ttl
	}
	return nil
}

func (x *GroupInviteRequest) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *GroupInviteRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *GroupInviteRequest) GetSingleUse() bool {
	if x != nil {
		return x.SingleUse
	}
	return false
}

type DeleteGroupInviteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id is the exact UIDP of the record.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteGroupInviteRequest) Reset() {
	*x = DeleteGroupInviteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_invites_platform_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteGroupInviteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteGroupInviteRequest) ProtoMessage() {}

func (x *DeleteGroupInviteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_group_invites_platform_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteGroupInviteRequest.ProtoReflect.Descriptor instead.
func (*DeleteGroupInviteRequest) Descriptor() ([]byte, []int) {
	return file_group_invites_platform_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteGroupInviteRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GroupInviteFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// group is used to identify the group this record is rooted under.
	Group string `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
	// id is the exact UID of the record.
	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	// key_id is the identify the verification key for this code.
	KeyId string `protobuf:"bytes,3,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
}

func (x *GroupInviteFilter) Reset() {
	*x = GroupInviteFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_invites_platform_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupInviteFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupInviteFilter) ProtoMessage() {}

func (x *GroupInviteFilter) ProtoReflect() protoreflect.Message {
	mi := &file_group_invites_platform_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupInviteFilter.ProtoReflect.Descriptor instead.
func (*GroupInviteFilter) Descriptor() ([]byte, []int) {
	return file_group_invites_platform_proto_rawDescGZIP(), []int{6}
}

func (x *GroupInviteFilter) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

func (x *GroupInviteFilter) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GroupInviteFilter) GetKeyId() string {
	if x != nil {
		return x.KeyId
	}
	return ""
}

type RegistrationRequest_Human struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// code is the json-encoded authentication code.
	// +optional
	Code string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *RegistrationRequest_Human) Reset() {
	*x = RegistrationRequest_Human{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_invites_platform_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistrationRequest_Human) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrationRequest_Human) ProtoMessage() {}

func (x *RegistrationRequest_Human) ProtoReflect() protoreflect.Message {
	mi := &file_group_invites_platform_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistrationRequest_Human.ProtoReflect.Descriptor instead.
func (*RegistrationRequest_Human) Descriptor() ([]byte, []int) {
	return file_group_invites_platform_proto_rawDescGZIP(), []int{0, 0}
}

func (x *RegistrationRequest_Human) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type RegistrationRequest_Cluster struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// code is the json-encoded authentication code.
	Code string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	// cluster_id is an optional cluster id if registering a cluster.
	ClusterId string `protobuf:"bytes,2,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
}

func (x *RegistrationRequest_Cluster) Reset() {
	*x = RegistrationRequest_Cluster{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_invites_platform_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistrationRequest_Cluster) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrationRequest_Cluster) ProtoMessage() {}

func (x *RegistrationRequest_Cluster) ProtoReflect() protoreflect.Message {
	mi := &file_group_invites_platform_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistrationRequest_Cluster.ProtoReflect.Descriptor instead.
func (*RegistrationRequest_Cluster) Descriptor() ([]byte, []int) {
	return file_group_invites_platform_proto_rawDescGZIP(), []int{0, 1}
}

func (x *RegistrationRequest_Cluster) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *RegistrationRequest_Cluster) GetClusterId() string {
	if x != nil {
		return x.ClusterId
	}
	return ""
}

var File_group_invites_platform_proto protoreflect.FileDescriptor

var file_group_invites_platform_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x73, 0x2e,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17,
	0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x72, 0x6f, 0x6c, 0x65, 0x2e, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x96, 0x02, 0x0a, 0x13, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x4a, 0x0a, 0x05, 0x68, 0x75, 0x6d, 0x61, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x32, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x48, 0x75, 0x6d, 0x61, 0x6e, 0x48, 0x00, 0x52, 0x05, 0x68, 0x75, 0x6d, 0x61, 0x6e, 0x12, 0x50,
	0x0a, 0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x34, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61,
	0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x48, 0x00, 0x52, 0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x1a, 0x1b, 0x0a, 0x05, 0x48, 0x75, 0x6d, 0x61, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x1a, 0x3c, 0x0a,
	0x07, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x6b,
	0x69, 0x6e, 0x64, 0x22, 0xb7, 0x01, 0x0a, 0x0b, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76,
	0x69, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x3a, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x15, 0x0a, 0x06, 0x6b, 0x65, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x6b, 0x65, 0x79, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72,
	0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x52,
	0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x99, 0x02,
	0x0a, 0x11, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x64, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76,
	0x69, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x3a, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x15, 0x0a, 0x06, 0x6b, 0x65, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x6b, 0x65, 0x79, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72,
	0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x52,
	0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12,
	0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x69,
	0x6e, 0x67, 0x6c, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09,
	0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x55, 0x73, 0x65, 0x22, 0x53, 0x0a, 0x0f, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x40, 0x0a, 0x05,
	0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x64, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0xa8,
	0x01, 0x0a, 0x12, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0x90, 0xaf, 0xa8, 0xd2, 0x05, 0x01, 0x52, 0x05, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x12, 0x2b, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x03, 0x74, 0x74, 0x6c,
	0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x72, 0x6f, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x69,
	0x6e, 0x67, 0x6c, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09,
	0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x55, 0x73, 0x65, 0x22, 0x32, 0x0a, 0x18, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x06, 0x90, 0xaf, 0xa8, 0xd2, 0x05, 0x01, 0x52, 0x02, 0x69, 0x64, 0x22, 0x50, 0x0a,
	0x11, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x15, 0x0a, 0x06, 0x6b, 0x65, 0x79, 0x5f,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6b, 0x65, 0x79, 0x49, 0x64, 0x32,
	0xb2, 0x05, 0x0a, 0x0c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x73,
	0x12, 0xd5, 0x01, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x2b, 0x2e, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76, 0x69, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e,
	0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69,
	0x61, 0x6d, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x22, 0x78,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x25, 0x3a, 0x01, 0x2a, 0x22, 0x20, 0x2f, 0x69, 0x61, 0x6d, 0x2f,
	0x76, 0x31, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x73,
	0x2f, 0x7b, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x3d, 0x2a, 0x2a, 0x7d, 0x8a, 0xaf, 0xa8, 0xd2, 0x05,
	0x08, 0x12, 0x06, 0x0a, 0x04, 0xc9, 0x01, 0x91, 0x03, 0xc2, 0xf0, 0x8e, 0xfc, 0x0b, 0x39, 0x0a,
	0x2e, 0x64, 0x65, 0x76, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x6e,
	0x76, 0x69, 0x74, 0x65, 0x2e, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x76, 0x31, 0x12,
	0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x12, 0x76, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x57, 0x69, 0x74, 0x68, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x2b, 0x2e, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76, 0x69, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e,
	0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69,
	0x61, 0x6d, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x22, 0x10,
	0x8a, 0xaf, 0xa8, 0xd2, 0x05, 0x0a, 0x12, 0x08, 0x0a, 0x04, 0xc9, 0x01, 0x91, 0x03, 0x10, 0x01,
	0x12, 0x89, 0x01, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2a, 0x2e, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e,
	0x69, 0x61, 0x6d, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x46,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x1a, 0x28, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61,
	0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x2e,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x22,
	0x2b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x12, 0x15, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x76, 0x31,
	0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x73, 0x8a, 0xaf,
	0xa8, 0xd2, 0x05, 0x08, 0x12, 0x06, 0x0a, 0x02, 0xcb, 0x01, 0x10, 0x01, 0x12, 0xc5, 0x01, 0x0a,
	0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x31, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67,
	0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61,
	0x6d, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76,
	0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x70, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1f, 0x2a, 0x1d, 0x2f, 0x69, 0x61, 0x6d,
	0x2f, 0x76, 0x31, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65,
	0x73, 0x2f, 0x7b, 0x69, 0x64, 0x3d, 0x2a, 0x2a, 0x7d, 0x8a, 0xaf, 0xa8, 0xd2, 0x05, 0x06, 0x12,
	0x04, 0x0a, 0x02, 0xcc, 0x01, 0xc2, 0xf0, 0x8e, 0xfc, 0x0b, 0x39, 0x0a, 0x2e, 0x64, 0x65, 0x76,
	0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x69, 0x61, 0x6d, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65,
	0x2e, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x2e, 0x76, 0x31, 0x12, 0x05, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x18, 0x01, 0x42, 0x2a, 0x5a, 0x28, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61,
	0x72, 0x64, 0x2e, 0x64, 0x65, 0x76, 0x2f, 0x73, 0x64, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x76, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_group_invites_platform_proto_rawDescOnce sync.Once
	file_group_invites_platform_proto_rawDescData = file_group_invites_platform_proto_rawDesc
)

func file_group_invites_platform_proto_rawDescGZIP() []byte {
	file_group_invites_platform_proto_rawDescOnce.Do(func() {
		file_group_invites_platform_proto_rawDescData = protoimpl.X.CompressGZIP(file_group_invites_platform_proto_rawDescData)
	})
	return file_group_invites_platform_proto_rawDescData
}

var file_group_invites_platform_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_group_invites_platform_proto_goTypes = []any{
	(*RegistrationRequest)(nil),         // 0: chainguard.platform.iam.RegistrationRequest
	(*GroupInvite)(nil),                 // 1: chainguard.platform.iam.GroupInvite
	(*StoredGroupInvite)(nil),           // 2: chainguard.platform.iam.StoredGroupInvite
	(*GroupInviteList)(nil),             // 3: chainguard.platform.iam.GroupInviteList
	(*GroupInviteRequest)(nil),          // 4: chainguard.platform.iam.GroupInviteRequest
	(*DeleteGroupInviteRequest)(nil),    // 5: chainguard.platform.iam.DeleteGroupInviteRequest
	(*GroupInviteFilter)(nil),           // 6: chainguard.platform.iam.GroupInviteFilter
	(*RegistrationRequest_Human)(nil),   // 7: chainguard.platform.iam.RegistrationRequest.Human
	(*RegistrationRequest_Cluster)(nil), // 8: chainguard.platform.iam.RegistrationRequest.Cluster
	(*timestamppb.Timestamp)(nil),       // 9: google.protobuf.Timestamp
	(*Role)(nil),                        // 10: chainguard.platform.iam.Role
	(*durationpb.Duration)(nil),         // 11: google.protobuf.Duration
	(*emptypb.Empty)(nil),               // 12: google.protobuf.Empty
}
var file_group_invites_platform_proto_depIdxs = []int32{
	7,  // 0: chainguard.platform.iam.RegistrationRequest.human:type_name -> chainguard.platform.iam.RegistrationRequest.Human
	8,  // 1: chainguard.platform.iam.RegistrationRequest.cluster:type_name -> chainguard.platform.iam.RegistrationRequest.Cluster
	9,  // 2: chainguard.platform.iam.GroupInvite.expiration:type_name -> google.protobuf.Timestamp
	10, // 3: chainguard.platform.iam.GroupInvite.role:type_name -> chainguard.platform.iam.Role
	9,  // 4: chainguard.platform.iam.StoredGroupInvite.expiration:type_name -> google.protobuf.Timestamp
	10, // 5: chainguard.platform.iam.StoredGroupInvite.role:type_name -> chainguard.platform.iam.Role
	9,  // 6: chainguard.platform.iam.StoredGroupInvite.created_at:type_name -> google.protobuf.Timestamp
	2,  // 7: chainguard.platform.iam.GroupInviteList.items:type_name -> chainguard.platform.iam.StoredGroupInvite
	11, // 8: chainguard.platform.iam.GroupInviteRequest.ttl:type_name -> google.protobuf.Duration
	4,  // 9: chainguard.platform.iam.GroupInvites.Create:input_type -> chainguard.platform.iam.GroupInviteRequest
	4,  // 10: chainguard.platform.iam.GroupInvites.CreateWithGroup:input_type -> chainguard.platform.iam.GroupInviteRequest
	6,  // 11: chainguard.platform.iam.GroupInvites.List:input_type -> chainguard.platform.iam.GroupInviteFilter
	5,  // 12: chainguard.platform.iam.GroupInvites.Delete:input_type -> chainguard.platform.iam.DeleteGroupInviteRequest
	1,  // 13: chainguard.platform.iam.GroupInvites.Create:output_type -> chainguard.platform.iam.GroupInvite
	1,  // 14: chainguard.platform.iam.GroupInvites.CreateWithGroup:output_type -> chainguard.platform.iam.GroupInvite
	3,  // 15: chainguard.platform.iam.GroupInvites.List:output_type -> chainguard.platform.iam.GroupInviteList
	12, // 16: chainguard.platform.iam.GroupInvites.Delete:output_type -> google.protobuf.Empty
	13, // [13:17] is the sub-list for method output_type
	9,  // [9:13] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_group_invites_platform_proto_init() }
func file_group_invites_platform_proto_init() {
	if File_group_invites_platform_proto != nil {
		return
	}
	file_role_platform_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_group_invites_platform_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*RegistrationRequest); i {
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
		file_group_invites_platform_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GroupInvite); i {
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
		file_group_invites_platform_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*StoredGroupInvite); i {
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
		file_group_invites_platform_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*GroupInviteList); i {
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
		file_group_invites_platform_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*GroupInviteRequest); i {
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
		file_group_invites_platform_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteGroupInviteRequest); i {
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
		file_group_invites_platform_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*GroupInviteFilter); i {
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
		file_group_invites_platform_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*RegistrationRequest_Human); i {
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
		file_group_invites_platform_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*RegistrationRequest_Cluster); i {
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
	file_group_invites_platform_proto_msgTypes[0].OneofWrappers = []any{
		(*RegistrationRequest_Human_)(nil),
		(*RegistrationRequest_Cluster_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_group_invites_platform_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_group_invites_platform_proto_goTypes,
		DependencyIndexes: file_group_invites_platform_proto_depIdxs,
		MessageInfos:      file_group_invites_platform_proto_msgTypes,
	}.Build()
	File_group_invites_platform_proto = out.File
	file_group_invites_platform_proto_rawDesc = nil
	file_group_invites_platform_proto_goTypes = nil
	file_group_invites_platform_proto_depIdxs = nil
}
