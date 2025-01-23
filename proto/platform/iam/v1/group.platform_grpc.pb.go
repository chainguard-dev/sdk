// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: group.platform.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Groups_Create_FullMethodName = "/chainguard.platform.iam.Groups/Create"
	Groups_Update_FullMethodName = "/chainguard.platform.iam.Groups/Update"
	Groups_List_FullMethodName   = "/chainguard.platform.iam.Groups/List"
	Groups_Delete_FullMethodName = "/chainguard.platform.iam.Groups/Delete"
)

// GroupsClient is the client API for Groups service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GroupsClient interface {
	Create(ctx context.Context, in *CreateGroupRequest, opts ...grpc.CallOption) (*Group, error)
	Update(ctx context.Context, in *Group, opts ...grpc.CallOption) (*Group, error)
	List(ctx context.Context, in *GroupFilter, opts ...grpc.CallOption) (*GroupList, error)
	Delete(ctx context.Context, in *DeleteGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type groupsClient struct {
	cc grpc.ClientConnInterface
}

func NewGroupsClient(cc grpc.ClientConnInterface) GroupsClient {
	return &groupsClient{cc}
}

func (c *groupsClient) Create(ctx context.Context, in *CreateGroupRequest, opts ...grpc.CallOption) (*Group, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Group)
	err := c.cc.Invoke(ctx, Groups_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupsClient) Update(ctx context.Context, in *Group, opts ...grpc.CallOption) (*Group, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Group)
	err := c.cc.Invoke(ctx, Groups_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupsClient) List(ctx context.Context, in *GroupFilter, opts ...grpc.CallOption) (*GroupList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GroupList)
	err := c.cc.Invoke(ctx, Groups_List_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupsClient) Delete(ctx context.Context, in *DeleteGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Groups_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GroupsServer is the server API for Groups service.
// All implementations must embed UnimplementedGroupsServer
// for forward compatibility.
type GroupsServer interface {
	Create(context.Context, *CreateGroupRequest) (*Group, error)
	Update(context.Context, *Group) (*Group, error)
	List(context.Context, *GroupFilter) (*GroupList, error)
	Delete(context.Context, *DeleteGroupRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedGroupsServer()
}

// UnimplementedGroupsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGroupsServer struct{}

func (UnimplementedGroupsServer) Create(context.Context, *CreateGroupRequest) (*Group, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedGroupsServer) Update(context.Context, *Group) (*Group, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedGroupsServer) List(context.Context, *GroupFilter) (*GroupList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedGroupsServer) Delete(context.Context, *DeleteGroupRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedGroupsServer) mustEmbedUnimplementedGroupsServer() {}
func (UnimplementedGroupsServer) testEmbeddedByValue()                {}

// UnsafeGroupsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GroupsServer will
// result in compilation errors.
type UnsafeGroupsServer interface {
	mustEmbedUnimplementedGroupsServer()
}

func RegisterGroupsServer(s grpc.ServiceRegistrar, srv GroupsServer) {
	// If the following call pancis, it indicates UnimplementedGroupsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Groups_ServiceDesc, srv)
}

func _Groups_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupsServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Groups_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupsServer).Create(ctx, req.(*CreateGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Groups_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Group)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupsServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Groups_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupsServer).Update(ctx, req.(*Group))
	}
	return interceptor(ctx, in, info, handler)
}

func _Groups_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GroupFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupsServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Groups_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupsServer).List(ctx, req.(*GroupFilter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Groups_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupsServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Groups_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupsServer).Delete(ctx, req.(*DeleteGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Groups_ServiceDesc is the grpc.ServiceDesc for Groups service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Groups_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chainguard.platform.iam.Groups",
	HandlerType: (*GroupsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Groups_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Groups_Update_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Groups_List_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Groups_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "group.platform.proto",
}
