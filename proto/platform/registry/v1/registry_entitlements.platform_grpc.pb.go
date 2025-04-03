// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: registry_entitlements.platform.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Entitlements_ListEntitlements_FullMethodName      = "/chainguard.platform.registry.Entitlements/ListEntitlements"
	Entitlements_ListEntitlementImages_FullMethodName = "/chainguard.platform.registry.Entitlements/ListEntitlementImages"
	Entitlements_Summary_FullMethodName               = "/chainguard.platform.registry.Entitlements/Summary"
)

// EntitlementsClient is the client API for Entitlements service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Entitlements is a service for viewing configuration and features enabled on a registry.
// NOTE: This API is EARLY ACCESS and is subject to change without notice.
type EntitlementsClient interface {
	ListEntitlements(ctx context.Context, in *EntitlementFilter, opts ...grpc.CallOption) (*EntitlementList, error)
	ListEntitlementImages(ctx context.Context, in *EntitlementImagesFilter, opts ...grpc.CallOption) (*EntitlementImagesList, error)
	// Summary provides a group-level summary of entitlements.
	Summary(ctx context.Context, in *EntitlementSummaryRequest, opts ...grpc.CallOption) (*EntitlementSummaryResponse, error)
}

type entitlementsClient struct {
	cc grpc.ClientConnInterface
}

func NewEntitlementsClient(cc grpc.ClientConnInterface) EntitlementsClient {
	return &entitlementsClient{cc}
}

func (c *entitlementsClient) ListEntitlements(ctx context.Context, in *EntitlementFilter, opts ...grpc.CallOption) (*EntitlementList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EntitlementList)
	err := c.cc.Invoke(ctx, Entitlements_ListEntitlements_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *entitlementsClient) ListEntitlementImages(ctx context.Context, in *EntitlementImagesFilter, opts ...grpc.CallOption) (*EntitlementImagesList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EntitlementImagesList)
	err := c.cc.Invoke(ctx, Entitlements_ListEntitlementImages_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *entitlementsClient) Summary(ctx context.Context, in *EntitlementSummaryRequest, opts ...grpc.CallOption) (*EntitlementSummaryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EntitlementSummaryResponse)
	err := c.cc.Invoke(ctx, Entitlements_Summary_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EntitlementsServer is the server API for Entitlements service.
// All implementations must embed UnimplementedEntitlementsServer
// for forward compatibility.
//
// Entitlements is a service for viewing configuration and features enabled on a registry.
// NOTE: This API is EARLY ACCESS and is subject to change without notice.
type EntitlementsServer interface {
	ListEntitlements(context.Context, *EntitlementFilter) (*EntitlementList, error)
	ListEntitlementImages(context.Context, *EntitlementImagesFilter) (*EntitlementImagesList, error)
	// Summary provides a group-level summary of entitlements.
	Summary(context.Context, *EntitlementSummaryRequest) (*EntitlementSummaryResponse, error)
	mustEmbedUnimplementedEntitlementsServer()
}

// UnimplementedEntitlementsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEntitlementsServer struct{}

func (UnimplementedEntitlementsServer) ListEntitlements(context.Context, *EntitlementFilter) (*EntitlementList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEntitlements not implemented")
}
func (UnimplementedEntitlementsServer) ListEntitlementImages(context.Context, *EntitlementImagesFilter) (*EntitlementImagesList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEntitlementImages not implemented")
}
func (UnimplementedEntitlementsServer) Summary(context.Context, *EntitlementSummaryRequest) (*EntitlementSummaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Summary not implemented")
}
func (UnimplementedEntitlementsServer) mustEmbedUnimplementedEntitlementsServer() {}
func (UnimplementedEntitlementsServer) testEmbeddedByValue()                      {}

// UnsafeEntitlementsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EntitlementsServer will
// result in compilation errors.
type UnsafeEntitlementsServer interface {
	mustEmbedUnimplementedEntitlementsServer()
}

func RegisterEntitlementsServer(s grpc.ServiceRegistrar, srv EntitlementsServer) {
	// If the following call pancis, it indicates UnimplementedEntitlementsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Entitlements_ServiceDesc, srv)
}

func _Entitlements_ListEntitlements_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntitlementFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntitlementsServer).ListEntitlements(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Entitlements_ListEntitlements_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntitlementsServer).ListEntitlements(ctx, req.(*EntitlementFilter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Entitlements_ListEntitlementImages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntitlementImagesFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntitlementsServer).ListEntitlementImages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Entitlements_ListEntitlementImages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntitlementsServer).ListEntitlementImages(ctx, req.(*EntitlementImagesFilter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Entitlements_Summary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntitlementSummaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntitlementsServer).Summary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Entitlements_Summary_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntitlementsServer).Summary(ctx, req.(*EntitlementSummaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Entitlements_ServiceDesc is the grpc.ServiceDesc for Entitlements service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Entitlements_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chainguard.platform.registry.Entitlements",
	HandlerType: (*EntitlementsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListEntitlements",
			Handler:    _Entitlements_ListEntitlements_Handler,
		},
		{
			MethodName: "ListEntitlementImages",
			Handler:    _Entitlements_ListEntitlementImages_Handler,
		},
		{
			MethodName: "Summary",
			Handler:    _Entitlements_Summary_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "registry_entitlements.platform.proto",
}
