// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.29.3
// source: role.platform.proto

package v1

import (
	_ "chainguard.dev/sdk/proto/annotations"
	v1 "chainguard.dev/sdk/proto/platform/common/v1"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Role struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id, The Group path under which this Role resides.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// name, human readable name of group.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// description, human readable description of group.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// capabilities, human readable list of capabilities supported by the group.
	Capabilities []string `protobuf:"bytes,4,rep,name=capabilities,proto3" json:"capabilities,omitempty"`
}

func (x *Role) Reset() {
	*x = Role{}
	if protoimpl.UnsafeEnabled {
		mi := &file_role_platform_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Role) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Role) ProtoMessage() {}

func (x *Role) ProtoReflect() protoreflect.Message {
	mi := &file_role_platform_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Role.ProtoReflect.Descriptor instead.
func (*Role) Descriptor() ([]byte, []int) {
	return file_role_platform_proto_rawDescGZIP(), []int{0}
}

func (x *Role) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Role) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Role) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Role) GetCapabilities() []string {
	if x != nil {
		return x.Capabilities
	}
	return nil
}

type RoleList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*Role `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *RoleList) Reset() {
	*x = RoleList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_role_platform_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoleList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoleList) ProtoMessage() {}

func (x *RoleList) ProtoReflect() protoreflect.Message {
	mi := &file_role_platform_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoleList.ProtoReflect.Descriptor instead.
func (*RoleList) Descriptor() ([]byte, []int) {
	return file_role_platform_proto_rawDescGZIP(), []int{1}
}

func (x *RoleList) GetItems() []*Role {
	if x != nil {
		return x.Items
	}
	return nil
}

type RoleFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id is the exact UIDP of the record.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// name is the exact name of the record
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// parent is the exact UIDP of the parent, or / for root
	Parent string `protobuf:"bytes,3,opt,name=parent,proto3" json:"parent,omitempty"`
	// uidp filters records based on their position in the group hierarchy.
	Uidp *v1.UIDPFilter `protobuf:"bytes,4,opt,name=uidp,proto3" json:"uidp,omitempty"`
}

func (x *RoleFilter) Reset() {
	*x = RoleFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_role_platform_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoleFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoleFilter) ProtoMessage() {}

func (x *RoleFilter) ProtoReflect() protoreflect.Message {
	mi := &file_role_platform_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoleFilter.ProtoReflect.Descriptor instead.
func (*RoleFilter) Descriptor() ([]byte, []int) {
	return file_role_platform_proto_rawDescGZIP(), []int{2}
}

func (x *RoleFilter) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RoleFilter) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RoleFilter) GetParent() string {
	if x != nil {
		return x.Parent
	}
	return ""
}

func (x *RoleFilter) GetUidp() *v1.UIDPFilter {
	if x != nil {
		return x.Uidp
	}
	return nil
}

type CreateRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// parent_id, The Group UIDP path under which the new Role resides.
	ParentId string `protobuf:"bytes,1,opt,name=parent_id,json=parentId,proto3" json:"parent_id,omitempty"`
	// Role to create.
	Role *Role `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *CreateRoleRequest) Reset() {
	*x = CreateRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_role_platform_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRoleRequest) ProtoMessage() {}

func (x *CreateRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_role_platform_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRoleRequest.ProtoReflect.Descriptor instead.
func (*CreateRoleRequest) Descriptor() ([]byte, []int) {
	return file_role_platform_proto_rawDescGZIP(), []int{3}
}

func (x *CreateRoleRequest) GetParentId() string {
	if x != nil {
		return x.ParentId
	}
	return ""
}

func (x *CreateRoleRequest) GetRole() *Role {
	if x != nil {
		return x.Role
	}
	return nil
}

type DeleteRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id is the exact UIDP of the record.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteRoleRequest) Reset() {
	*x = DeleteRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_role_platform_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRoleRequest) ProtoMessage() {}

func (x *DeleteRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_role_platform_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRoleRequest.ProtoReflect.Descriptor instead.
func (*DeleteRoleRequest) Descriptor() ([]byte, []int) {
	return file_role_platform_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteRoleRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_role_platform_proto protoreflect.FileDescriptor

var file_role_platform_proto_rawDesc = []byte{
	0x0a, 0x13, 0x72, 0x6f, 0x6c, 0x65, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72,
	0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x26, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x69, 0x64, 0x70, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x78, 0x0a, 0x04, 0x52, 0x6f, 0x6c,
	0x65, 0x12, 0x16, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0x90,
	0xaf, 0xa8, 0xd2, 0x05, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x22, 0x0a, 0x0c, 0x63, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74,
	0x69, 0x65, 0x73, 0x22, 0x3f, 0x0a, 0x08, 0x52, 0x6f, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12,
	0x33, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d,
	0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x05, 0x69,
	0x74, 0x65, 0x6d, 0x73, 0x22, 0x84, 0x01, 0x0a, 0x0a, 0x52, 0x6f, 0x6c, 0x65, 0x46, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x12,
	0x3a, 0x0a, 0x04, 0x75, 0x69, 0x64, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e,
	0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x55, 0x49, 0x44, 0x50, 0x46,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x04, 0x75, 0x69, 0x64, 0x70, 0x22, 0x6b, 0x0a, 0x11, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x23, 0x0a, 0x09, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x06, 0x90, 0xaf, 0xa8, 0xd2, 0x05, 0x01, 0x52, 0x08, 0x70, 0x61, 0x72,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64,
	0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x52, 0x6f,
	0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0x2b, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0x90, 0xaf, 0xa8, 0xd2, 0x05,
	0x01, 0x52, 0x02, 0x69, 0x64, 0x32, 0xf9, 0x03, 0x0a, 0x05, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x12,
	0x8b, 0x01, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x2a, 0x2e, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75,
	0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61, 0x6d,
	0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x22, 0x36, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x24, 0x3a, 0x04, 0x72,
	0x6f, 0x6c, 0x65, 0x22, 0x1c, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x6f, 0x6c,
	0x65, 0x73, 0x2f, 0x7b, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x3d, 0x2a, 0x2a,
	0x7d, 0x8a, 0xaf, 0xa8, 0xd2, 0x05, 0x06, 0x12, 0x04, 0x0a, 0x02, 0xad, 0x02, 0x12, 0x74, 0x0a,
	0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67,
	0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61,
	0x6d, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x1a, 0x1d, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75,
	0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61, 0x6d,
	0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x22, 0x2c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x3a, 0x01, 0x2a,
	0x1a, 0x15, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2f,
	0x7b, 0x69, 0x64, 0x3d, 0x2a, 0x2a, 0x7d, 0x8a, 0xaf, 0xa8, 0xd2, 0x05, 0x06, 0x12, 0x04, 0x0a,
	0x02, 0xae, 0x02, 0x12, 0x73, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x23, 0x2e, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x1a, 0x21, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c,
	0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x4c,
	0x69, 0x73, 0x74, 0x22, 0x23, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x12, 0x0d, 0x2f, 0x69, 0x61,
	0x6d, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x8a, 0xaf, 0xa8, 0xd2, 0x05, 0x08,
	0x12, 0x06, 0x0a, 0x02, 0xaf, 0x02, 0x10, 0x01, 0x12, 0x77, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x12, 0x2a, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x29, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x2a, 0x15,
	0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2f, 0x7b, 0x69,
	0x64, 0x3d, 0x2a, 0x2a, 0x7d, 0x8a, 0xaf, 0xa8, 0xd2, 0x05, 0x06, 0x12, 0x04, 0x0a, 0x02, 0xb0,
	0x02, 0x42, 0x66, 0x0a, 0x22, 0x64, 0x65, 0x76, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75,
	0x61, 0x72, 0x64, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x42, 0x14, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x49, 0x41, 0x4d, 0x52, 0x6f, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x28, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x64, 0x65, 0x76, 0x2f,
	0x73, 0x64, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_role_platform_proto_rawDescOnce sync.Once
	file_role_platform_proto_rawDescData = file_role_platform_proto_rawDesc
)

func file_role_platform_proto_rawDescGZIP() []byte {
	file_role_platform_proto_rawDescOnce.Do(func() {
		file_role_platform_proto_rawDescData = protoimpl.X.CompressGZIP(file_role_platform_proto_rawDescData)
	})
	return file_role_platform_proto_rawDescData
}

var file_role_platform_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_role_platform_proto_goTypes = []any{
	(*Role)(nil),              // 0: chainguard.platform.iam.Role
	(*RoleList)(nil),          // 1: chainguard.platform.iam.RoleList
	(*RoleFilter)(nil),        // 2: chainguard.platform.iam.RoleFilter
	(*CreateRoleRequest)(nil), // 3: chainguard.platform.iam.CreateRoleRequest
	(*DeleteRoleRequest)(nil), // 4: chainguard.platform.iam.DeleteRoleRequest
	(*v1.UIDPFilter)(nil),     // 5: chainguard.platform.common.UIDPFilter
	(*emptypb.Empty)(nil),     // 6: google.protobuf.Empty
}
var file_role_platform_proto_depIdxs = []int32{
	0, // 0: chainguard.platform.iam.RoleList.items:type_name -> chainguard.platform.iam.Role
	5, // 1: chainguard.platform.iam.RoleFilter.uidp:type_name -> chainguard.platform.common.UIDPFilter
	0, // 2: chainguard.platform.iam.CreateRoleRequest.role:type_name -> chainguard.platform.iam.Role
	3, // 3: chainguard.platform.iam.Roles.Create:input_type -> chainguard.platform.iam.CreateRoleRequest
	0, // 4: chainguard.platform.iam.Roles.Update:input_type -> chainguard.platform.iam.Role
	2, // 5: chainguard.platform.iam.Roles.List:input_type -> chainguard.platform.iam.RoleFilter
	4, // 6: chainguard.platform.iam.Roles.Delete:input_type -> chainguard.platform.iam.DeleteRoleRequest
	0, // 7: chainguard.platform.iam.Roles.Create:output_type -> chainguard.platform.iam.Role
	0, // 8: chainguard.platform.iam.Roles.Update:output_type -> chainguard.platform.iam.Role
	1, // 9: chainguard.platform.iam.Roles.List:output_type -> chainguard.platform.iam.RoleList
	6, // 10: chainguard.platform.iam.Roles.Delete:output_type -> google.protobuf.Empty
	7, // [7:11] is the sub-list for method output_type
	3, // [3:7] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_role_platform_proto_init() }
func file_role_platform_proto_init() {
	if File_role_platform_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_role_platform_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Role); i {
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
		file_role_platform_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*RoleList); i {
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
		file_role_platform_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*RoleFilter); i {
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
		file_role_platform_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*CreateRoleRequest); i {
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
		file_role_platform_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteRoleRequest); i {
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
			RawDescriptor: file_role_platform_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_role_platform_proto_goTypes,
		DependencyIndexes: file_role_platform_proto_depIdxs,
		MessageInfos:      file_role_platform_proto_msgTypes,
	}.Build()
	File_role_platform_proto = out.File
	file_role_platform_proto_rawDesc = nil
	file_role_platform_proto_goTypes = nil
	file_role_platform_proto_depIdxs = nil
}
