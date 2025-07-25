// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.29.3
// source: auth.proto

package annotations

import (
	capabilities "chainguard.dev/sdk/proto/capabilities"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
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

type IAM struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Mode:
	//
	//	*IAM_Disabled
	//	*IAM_Enabled
	Mode isIAM_Mode `protobuf_oneof:"mode"`
}

func (x *IAM) Reset() {
	*x = IAM{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IAM) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IAM) ProtoMessage() {}

func (x *IAM) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IAM.ProtoReflect.Descriptor instead.
func (*IAM) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{0}
}

func (m *IAM) GetMode() isIAM_Mode {
	if m != nil {
		return m.Mode
	}
	return nil
}

func (x *IAM) GetDisabled() *emptypb.Empty {
	if x, ok := x.GetMode().(*IAM_Disabled); ok {
		return x.Disabled
	}
	return nil
}

func (x *IAM) GetEnabled() *IAM_Rules {
	if x, ok := x.GetMode().(*IAM_Enabled); ok {
		return x.Enabled
	}
	return nil
}

type isIAM_Mode interface {
	isIAM_Mode()
}

type IAM_Disabled struct {
	Disabled *emptypb.Empty `protobuf:"bytes,1,opt,name=disabled,proto3,oneof"`
}

type IAM_Enabled struct {
	Enabled *IAM_Rules `protobuf:"bytes,2,opt,name=enabled,proto3,oneof"`
}

func (*IAM_Disabled) isIAM_Mode() {}

func (*IAM_Enabled) isIAM_Mode() {}

type IAM_Rules struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A list of capabilities required by a particular API.
	// This field is either scoped or unscoped, as determined
	// by the field below.
	//   - When it is "scoped", this field is combined with the
	//     field designated by "(iam_scope) = true" (see below)
	//     on the request message to indicate what capabilities
	//     the caller needs at what scope in order to authorize
	//     the action they are performing.
	//   - When is it "unscoped", this field is used to determine
	//     the set of scopes the caller has the appropriate access
	//     to so that the RPC itself can scope down the results
	//     it returns.
	Capabilities []capabilities.Capability `protobuf:"varint,1,rep,packed,name=capabilities,proto3,enum=chainguard.capabilities.Capability" json:"capabilities,omitempty"`
	// Unscoped is set on APIs where the request itself doesn't
	// carry a field with "iam_scope", and instead scopes itself
	// to the set of groups to which the caller has access
	// according to their OIDC token.
	Unscoped bool `protobuf:"varint,2,opt,name=unscoped,proto3" json:"unscoped,omitempty"`
}

func (x *IAM_Rules) Reset() {
	*x = IAM_Rules{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IAM_Rules) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IAM_Rules) ProtoMessage() {}

func (x *IAM_Rules) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IAM_Rules.ProtoReflect.Descriptor instead.
func (*IAM_Rules) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{0, 0}
}

func (x *IAM_Rules) GetCapabilities() []capabilities.Capability {
	if x != nil {
		return x.Capabilities
	}
	return nil
}

func (x *IAM_Rules) GetUnscoped() bool {
	if x != nil {
		return x.Unscoped
	}
	return false
}

var file_auth_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*IAM)(nil),
		Field:         189350641,
		Name:          "chainguard.annotations.iam",
		Tag:           "bytes,189350641,opt,name=iam",
		Filename:      "auth.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         189350642,
		Name:          "chainguard.annotations.iam_scope",
		Tag:           "varint,189350642,opt,name=iam_scope",
		Filename:      "auth.proto",
	},
}

// Extension fields to descriptorpb.MethodOptions.
var (
	// optional chainguard.annotations.IAM iam = 189350641;
	E_Iam = &file_auth_proto_extTypes[0] // randomly chosen
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional bool iam_scope = 189350642;
	E_IamScope = &file_auth_proto_extTypes[1] // one more than above
)

var File_auth_proto protoreflect.FileDescriptor

var file_auth_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x63, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65,
	0x73, 0x2f, 0x63, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf0, 0x01, 0x0a, 0x03, 0x49, 0x41, 0x4d, 0x12, 0x34, 0x0a, 0x08,
	0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x48, 0x00, 0x52, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c,
	0x65, 0x64, 0x12, 0x3d, 0x0a, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64,
	0x2e, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x49, 0x41, 0x4d,
	0x2e, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x48, 0x00, 0x52, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65,
	0x64, 0x1a, 0x6c, 0x0a, 0x05, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x47, 0x0a, 0x0c, 0x63, 0x61,
	0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0e,
	0x32, 0x23, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x63, 0x61,
	0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x70, 0x61, 0x62,
	0x69, 0x6c, 0x69, 0x74, 0x79, 0x52, 0x0c, 0x63, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74,
	0x69, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x6e, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x75, 0x6e, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x64, 0x42,
	0x06, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x3a, 0x50, 0x0a, 0x03, 0x69, 0x61, 0x6d, 0x12, 0x1e,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xf1,
	0x85, 0xa5, 0x5a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x67,
	0x75, 0x61, 0x72, 0x64, 0x2e, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x49, 0x41, 0x4d, 0x52, 0x03, 0x69, 0x61, 0x6d, 0x3a, 0x3d, 0x0a, 0x09, 0x69, 0x61, 0x6d,
	0x5f, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xf2, 0x85, 0xa5, 0x5a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08,
	0x69, 0x61, 0x6d, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x42, 0x26, 0x5a, 0x24, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x2e, 0x64, 0x65, 0x76, 0x2f, 0x73, 0x64, 0x6b, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_proto_rawDescOnce sync.Once
	file_auth_proto_rawDescData = file_auth_proto_rawDesc
)

func file_auth_proto_rawDescGZIP() []byte {
	file_auth_proto_rawDescOnce.Do(func() {
		file_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_proto_rawDescData)
	})
	return file_auth_proto_rawDescData
}

var file_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_auth_proto_goTypes = []any{
	(*IAM)(nil),                        // 0: chainguard.annotations.IAM
	(*IAM_Rules)(nil),                  // 1: chainguard.annotations.IAM.Rules
	(*emptypb.Empty)(nil),              // 2: google.protobuf.Empty
	(capabilities.Capability)(0),       // 3: chainguard.capabilities.Capability
	(*descriptorpb.MethodOptions)(nil), // 4: google.protobuf.MethodOptions
	(*descriptorpb.FieldOptions)(nil),  // 5: google.protobuf.FieldOptions
}
var file_auth_proto_depIdxs = []int32{
	2, // 0: chainguard.annotations.IAM.disabled:type_name -> google.protobuf.Empty
	1, // 1: chainguard.annotations.IAM.enabled:type_name -> chainguard.annotations.IAM.Rules
	3, // 2: chainguard.annotations.IAM.Rules.capabilities:type_name -> chainguard.capabilities.Capability
	4, // 3: chainguard.annotations.iam:extendee -> google.protobuf.MethodOptions
	5, // 4: chainguard.annotations.iam_scope:extendee -> google.protobuf.FieldOptions
	0, // 5: chainguard.annotations.iam:type_name -> chainguard.annotations.IAM
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	5, // [5:6] is the sub-list for extension type_name
	3, // [3:5] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_auth_proto_init() }
func file_auth_proto_init() {
	if File_auth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_auth_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*IAM); i {
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
		file_auth_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*IAM_Rules); i {
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
	file_auth_proto_msgTypes[0].OneofWrappers = []any{
		(*IAM_Disabled)(nil),
		(*IAM_Enabled)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_auth_proto_goTypes,
		DependencyIndexes: file_auth_proto_depIdxs,
		MessageInfos:      file_auth_proto_msgTypes,
		ExtensionInfos:    file_auth_proto_extTypes,
	}.Build()
	File_auth_proto = out.File
	file_auth_proto_rawDesc = nil
	file_auth_proto_goTypes = nil
	file_auth_proto_depIdxs = nil
}
