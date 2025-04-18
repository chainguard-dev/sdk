// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: vulnerabilities.platform.proto

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
	Vulnerabilities_ListVulnReports_FullMethodName                = "/chainguard.platform.registry.Vulnerabilities/ListVulnReports"
	Vulnerabilities_GetRawVulnReport_FullMethodName               = "/chainguard.platform.registry.Vulnerabilities/GetRawVulnReport"
	Vulnerabilities_ListVulnCountReports_FullMethodName           = "/chainguard.platform.registry.Vulnerabilities/ListVulnCountReports"
	Vulnerabilities_ListCumulativeVulnCountReports_FullMethodName = "/chainguard.platform.registry.Vulnerabilities/ListCumulativeVulnCountReports"
)

// VulnerabilitiesClient is the client API for Vulnerabilities service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VulnerabilitiesClient interface {
	ListVulnReports(ctx context.Context, in *VulnReportFilter, opts ...grpc.CallOption) (*VulnReportList, error)
	GetRawVulnReport(ctx context.Context, in *GetRawVulnReportRequest, opts ...grpc.CallOption) (*RawVulnReport, error)
	ListVulnCountReports(ctx context.Context, in *VulnCountReportFilter, opts ...grpc.CallOption) (*VulnCountReportList, error)
	ListCumulativeVulnCountReports(ctx context.Context, in *VulnCountReportFilter, opts ...grpc.CallOption) (*VulnCountReportList, error)
}

type vulnerabilitiesClient struct {
	cc grpc.ClientConnInterface
}

func NewVulnerabilitiesClient(cc grpc.ClientConnInterface) VulnerabilitiesClient {
	return &vulnerabilitiesClient{cc}
}

func (c *vulnerabilitiesClient) ListVulnReports(ctx context.Context, in *VulnReportFilter, opts ...grpc.CallOption) (*VulnReportList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VulnReportList)
	err := c.cc.Invoke(ctx, Vulnerabilities_ListVulnReports_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vulnerabilitiesClient) GetRawVulnReport(ctx context.Context, in *GetRawVulnReportRequest, opts ...grpc.CallOption) (*RawVulnReport, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RawVulnReport)
	err := c.cc.Invoke(ctx, Vulnerabilities_GetRawVulnReport_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vulnerabilitiesClient) ListVulnCountReports(ctx context.Context, in *VulnCountReportFilter, opts ...grpc.CallOption) (*VulnCountReportList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VulnCountReportList)
	err := c.cc.Invoke(ctx, Vulnerabilities_ListVulnCountReports_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vulnerabilitiesClient) ListCumulativeVulnCountReports(ctx context.Context, in *VulnCountReportFilter, opts ...grpc.CallOption) (*VulnCountReportList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VulnCountReportList)
	err := c.cc.Invoke(ctx, Vulnerabilities_ListCumulativeVulnCountReports_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VulnerabilitiesServer is the server API for Vulnerabilities service.
// All implementations must embed UnimplementedVulnerabilitiesServer
// for forward compatibility.
type VulnerabilitiesServer interface {
	ListVulnReports(context.Context, *VulnReportFilter) (*VulnReportList, error)
	GetRawVulnReport(context.Context, *GetRawVulnReportRequest) (*RawVulnReport, error)
	ListVulnCountReports(context.Context, *VulnCountReportFilter) (*VulnCountReportList, error)
	ListCumulativeVulnCountReports(context.Context, *VulnCountReportFilter) (*VulnCountReportList, error)
	mustEmbedUnimplementedVulnerabilitiesServer()
}

// UnimplementedVulnerabilitiesServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedVulnerabilitiesServer struct{}

func (UnimplementedVulnerabilitiesServer) ListVulnReports(context.Context, *VulnReportFilter) (*VulnReportList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVulnReports not implemented")
}
func (UnimplementedVulnerabilitiesServer) GetRawVulnReport(context.Context, *GetRawVulnReportRequest) (*RawVulnReport, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRawVulnReport not implemented")
}
func (UnimplementedVulnerabilitiesServer) ListVulnCountReports(context.Context, *VulnCountReportFilter) (*VulnCountReportList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVulnCountReports not implemented")
}
func (UnimplementedVulnerabilitiesServer) ListCumulativeVulnCountReports(context.Context, *VulnCountReportFilter) (*VulnCountReportList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCumulativeVulnCountReports not implemented")
}
func (UnimplementedVulnerabilitiesServer) mustEmbedUnimplementedVulnerabilitiesServer() {}
func (UnimplementedVulnerabilitiesServer) testEmbeddedByValue()                         {}

// UnsafeVulnerabilitiesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VulnerabilitiesServer will
// result in compilation errors.
type UnsafeVulnerabilitiesServer interface {
	mustEmbedUnimplementedVulnerabilitiesServer()
}

func RegisterVulnerabilitiesServer(s grpc.ServiceRegistrar, srv VulnerabilitiesServer) {
	// If the following call pancis, it indicates UnimplementedVulnerabilitiesServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Vulnerabilities_ServiceDesc, srv)
}

func _Vulnerabilities_ListVulnReports_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VulnReportFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VulnerabilitiesServer).ListVulnReports(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vulnerabilities_ListVulnReports_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VulnerabilitiesServer).ListVulnReports(ctx, req.(*VulnReportFilter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vulnerabilities_GetRawVulnReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRawVulnReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VulnerabilitiesServer).GetRawVulnReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vulnerabilities_GetRawVulnReport_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VulnerabilitiesServer).GetRawVulnReport(ctx, req.(*GetRawVulnReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vulnerabilities_ListVulnCountReports_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VulnCountReportFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VulnerabilitiesServer).ListVulnCountReports(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vulnerabilities_ListVulnCountReports_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VulnerabilitiesServer).ListVulnCountReports(ctx, req.(*VulnCountReportFilter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vulnerabilities_ListCumulativeVulnCountReports_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VulnCountReportFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VulnerabilitiesServer).ListCumulativeVulnCountReports(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vulnerabilities_ListCumulativeVulnCountReports_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VulnerabilitiesServer).ListCumulativeVulnCountReports(ctx, req.(*VulnCountReportFilter))
	}
	return interceptor(ctx, in, info, handler)
}

// Vulnerabilities_ServiceDesc is the grpc.ServiceDesc for Vulnerabilities service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Vulnerabilities_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chainguard.platform.registry.Vulnerabilities",
	HandlerType: (*VulnerabilitiesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListVulnReports",
			Handler:    _Vulnerabilities_ListVulnReports_Handler,
		},
		{
			MethodName: "GetRawVulnReport",
			Handler:    _Vulnerabilities_GetRawVulnReport_Handler,
		},
		{
			MethodName: "ListVulnCountReports",
			Handler:    _Vulnerabilities_ListVulnCountReports_Handler,
		},
		{
			MethodName: "ListCumulativeVulnCountReports",
			Handler:    _Vulnerabilities_ListCumulativeVulnCountReports_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "vulnerabilities.platform.proto",
}
