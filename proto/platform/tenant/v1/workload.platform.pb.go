// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.21.12
// source: workload.platform.proto

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

type Workload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id, The Workload UIDP at which this Workload resides.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// name of the Workload.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// a short description of this Workload.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// remote_id is the remote ID of this Workload.
	RemoteId   string `protobuf:"bytes,4,opt,name=remote_id,json=remoteId,proto3" json:"remote_id,omitempty"`
	Labels     string `protobuf:"bytes,5,opt,name=labels,proto3" json:"labels,omitempty"`
	ApiVersion string `protobuf:"bytes,6,opt,name=api_version,json=apiVersion,proto3" json:"api_version,omitempty"`
	Kind       string `protobuf:"bytes,7,opt,name=kind,proto3" json:"kind,omitempty"`
	// last_seen tracks the timestamp at which this workload was last seen.
	LastSeen *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=last_seen,json=lastSeen,proto3" json:"last_seen,omitempty"`
	// owner_id is the remote_id of the Workload that is referenced via a
	// "controller" owner reference by this workload.
	OwnerId string `protobuf:"bytes,9,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
}

func (x *Workload) Reset() {
	*x = Workload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workload_platform_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Workload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Workload) ProtoMessage() {}

func (x *Workload) ProtoReflect() protoreflect.Message {
	mi := &file_workload_platform_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Workload.ProtoReflect.Descriptor instead.
func (*Workload) Descriptor() ([]byte, []int) {
	return file_workload_platform_proto_rawDescGZIP(), []int{0}
}

func (x *Workload) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Workload) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Workload) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Workload) GetRemoteId() string {
	if x != nil {
		return x.RemoteId
	}
	return ""
}

func (x *Workload) GetLabels() string {
	if x != nil {
		return x.Labels
	}
	return ""
}

func (x *Workload) GetApiVersion() string {
	if x != nil {
		return x.ApiVersion
	}
	return ""
}

func (x *Workload) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *Workload) GetLastSeen() *timestamppb.Timestamp {
	if x != nil {
		return x.LastSeen
	}
	return nil
}

func (x *Workload) GetOwnerId() string {
	if x != nil {
		return x.OwnerId
	}
	return ""
}

type WorkloadList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*Workload `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *WorkloadList) Reset() {
	*x = WorkloadList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workload_platform_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkloadList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkloadList) ProtoMessage() {}

func (x *WorkloadList) ProtoReflect() protoreflect.Message {
	mi := &file_workload_platform_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkloadList.ProtoReflect.Descriptor instead.
func (*WorkloadList) Descriptor() ([]byte, []int) {
	return file_workload_platform_proto_rawDescGZIP(), []int{1}
}

func (x *WorkloadList) GetItems() []*Workload {
	if x != nil {
		return x.Items
	}
	return nil
}

type WorkloadFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// remote_id is the remote ID of this Workload.
	RemoteId string `protobuf:"bytes,2,opt,name=remote_id,json=remoteId,proto3" json:"remote_id,omitempty"`
	// active_since is the timestamp after which returned workloads
	// should have been active (last seen).
	ActiveSince *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=active_since,json=activeSince,proto3" json:"active_since,omitempty"`
	// name filters on the resource name.
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// Return only the workloads owned by a particular remote_id.
	OwnerId string         `protobuf:"bytes,5,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	Uidp    *v1.UIDPFilter `protobuf:"bytes,100,opt,name=uidp,proto3" json:"uidp,omitempty"`
}

func (x *WorkloadFilter) Reset() {
	*x = WorkloadFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workload_platform_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkloadFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkloadFilter) ProtoMessage() {}

func (x *WorkloadFilter) ProtoReflect() protoreflect.Message {
	mi := &file_workload_platform_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkloadFilter.ProtoReflect.Descriptor instead.
func (*WorkloadFilter) Descriptor() ([]byte, []int) {
	return file_workload_platform_proto_rawDescGZIP(), []int{2}
}

func (x *WorkloadFilter) GetRemoteId() string {
	if x != nil {
		return x.RemoteId
	}
	return ""
}

func (x *WorkloadFilter) GetActiveSince() *timestamppb.Timestamp {
	if x != nil {
		return x.ActiveSince
	}
	return nil
}

func (x *WorkloadFilter) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *WorkloadFilter) GetOwnerId() string {
	if x != nil {
		return x.OwnerId
	}
	return ""
}

func (x *WorkloadFilter) GetUidp() *v1.UIDPFilter {
	if x != nil {
		return x.Uidp
	}
	return nil
}

var File_workload_platform_proto protoreflect.FileDescriptor

var file_workload_platform_proto_rawDesc = []byte{
	0x0a, 0x17, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x63, 0x68, 0x61, 0x69, 0x6e,
	0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x74,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x26,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f,
	0x76, 0x31, 0x2f, 0x75, 0x69, 0x64, 0x70, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8e, 0x02, 0x0a, 0x08, 0x57, 0x6f, 0x72, 0x6b, 0x6c,
	0x6f, 0x61, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x6d,
	0x6f, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65,
	0x6d, 0x6f, 0x74, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x1f,
	0x0a, 0x0b, 0x61, 0x70, 0x69, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b,
	0x69, 0x6e, 0x64, 0x12, 0x37, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x73, 0x65, 0x65, 0x6e,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x53, 0x65, 0x65, 0x6e, 0x12, 0x19, 0x0a, 0x08,
	0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x22, 0x4a, 0x0a, 0x0c, 0x57, 0x6f, 0x72, 0x6b, 0x6c,
	0x6f, 0x61, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x3a, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75,
	0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x74, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x05, 0x69, 0x74,
	0x65, 0x6d, 0x73, 0x22, 0xd7, 0x01, 0x0a, 0x0e, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64,
	0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x6d, 0x6f, 0x74,
	0x65, 0x49, 0x64, 0x12, 0x3d, 0x0a, 0x0c, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x73, 0x69,
	0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x69, 0x6e,
	0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x3a, 0x0a, 0x04, 0x75, 0x69, 0x64, 0x70, 0x18, 0x64, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x26, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61,
	0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x55, 0x49, 0x44,
	0x50, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x04, 0x75, 0x69, 0x64, 0x70, 0x32, 0x79, 0x0a,
	0x09, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x73, 0x12, 0x6c, 0x0a, 0x04, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x2a, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e,
	0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x1a, 0x28,
	0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x57, 0x6f, 0x72, 0x6b,
	0x6c, 0x6f, 0x61, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x0e, 0x8a, 0xaf, 0xa8, 0xd2, 0x05, 0x08,
	0x12, 0x06, 0x0a, 0x02, 0xb3, 0x09, 0x10, 0x01, 0x42, 0x73, 0x0a, 0x25, 0x64, 0x65, 0x76, 0x2e,
	0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x70,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x76,
	0x31, 0x42, 0x1b, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x54, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x2b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x64, 0x65, 0x76,
	0x2f, 0x73, 0x64, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x2f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_workload_platform_proto_rawDescOnce sync.Once
	file_workload_platform_proto_rawDescData = file_workload_platform_proto_rawDesc
)

func file_workload_platform_proto_rawDescGZIP() []byte {
	file_workload_platform_proto_rawDescOnce.Do(func() {
		file_workload_platform_proto_rawDescData = protoimpl.X.CompressGZIP(file_workload_platform_proto_rawDescData)
	})
	return file_workload_platform_proto_rawDescData
}

var file_workload_platform_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_workload_platform_proto_goTypes = []interface{}{
	(*Workload)(nil),              // 0: chainguard.platform.tenant.Workload
	(*WorkloadList)(nil),          // 1: chainguard.platform.tenant.WorkloadList
	(*WorkloadFilter)(nil),        // 2: chainguard.platform.tenant.WorkloadFilter
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
	(*v1.UIDPFilter)(nil),         // 4: chainguard.platform.common.UIDPFilter
}
var file_workload_platform_proto_depIdxs = []int32{
	3, // 0: chainguard.platform.tenant.Workload.last_seen:type_name -> google.protobuf.Timestamp
	0, // 1: chainguard.platform.tenant.WorkloadList.items:type_name -> chainguard.platform.tenant.Workload
	3, // 2: chainguard.platform.tenant.WorkloadFilter.active_since:type_name -> google.protobuf.Timestamp
	4, // 3: chainguard.platform.tenant.WorkloadFilter.uidp:type_name -> chainguard.platform.common.UIDPFilter
	2, // 4: chainguard.platform.tenant.Workloads.List:input_type -> chainguard.platform.tenant.WorkloadFilter
	1, // 5: chainguard.platform.tenant.Workloads.List:output_type -> chainguard.platform.tenant.WorkloadList
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_workload_platform_proto_init() }
func file_workload_platform_proto_init() {
	if File_workload_platform_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_workload_platform_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Workload); i {
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
		file_workload_platform_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkloadList); i {
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
		file_workload_platform_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkloadFilter); i {
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
			RawDescriptor: file_workload_platform_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_workload_platform_proto_goTypes,
		DependencyIndexes: file_workload_platform_proto_depIdxs,
		MessageInfos:      file_workload_platform_proto_msgTypes,
	}.Build()
	File_workload_platform_proto = out.File
	file_workload_platform_proto_rawDesc = nil
	file_workload_platform_proto_goTypes = nil
	file_workload_platform_proto_depIdxs = nil
}
