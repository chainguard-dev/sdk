// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: registry_entitlements.platform.proto

package v1

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_Entitlements_ListEntitlements_0(ctx context.Context, marshaler runtime.Marshaler, client EntitlementsClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq EntitlementFilter
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["parent"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "parent")
	}

	protoReq.Parent, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "parent", err)
	}

	msg, err := client.ListEntitlements(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Entitlements_ListEntitlements_0(ctx context.Context, marshaler runtime.Marshaler, server EntitlementsServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq EntitlementFilter
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["parent"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "parent")
	}

	protoReq.Parent, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "parent", err)
	}

	msg, err := server.ListEntitlements(ctx, &protoReq)
	return msg, metadata, err

}

func request_Entitlements_ListEntitlementImages_0(ctx context.Context, marshaler runtime.Marshaler, client EntitlementsClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq EntitlementImagesFilter
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["parent"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "parent")
	}

	protoReq.Parent, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "parent", err)
	}

	msg, err := client.ListEntitlementImages(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Entitlements_ListEntitlementImages_0(ctx context.Context, marshaler runtime.Marshaler, server EntitlementsServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq EntitlementImagesFilter
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["parent"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "parent")
	}

	protoReq.Parent, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "parent", err)
	}

	msg, err := server.ListEntitlementImages(ctx, &protoReq)
	return msg, metadata, err

}

func request_Entitlements_Summary_0(ctx context.Context, marshaler runtime.Marshaler, client EntitlementsClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq EntitlementSummaryRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["parent"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "parent")
	}

	protoReq.Parent, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "parent", err)
	}

	msg, err := client.Summary(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Entitlements_Summary_0(ctx context.Context, marshaler runtime.Marshaler, server EntitlementsServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq EntitlementSummaryRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["parent"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "parent")
	}

	protoReq.Parent, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "parent", err)
	}

	msg, err := server.Summary(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterEntitlementsHandlerServer registers the http handlers for service Entitlements to "mux".
// UnaryRPC     :call EntitlementsServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterEntitlementsHandlerFromEndpoint instead.
// GRPC interceptors will not work for this type of registration. To use interceptors, you must use the "runtime.WithMiddlewares" option in the "runtime.NewServeMux" call.
func RegisterEntitlementsHandlerServer(ctx context.Context, mux *runtime.ServeMux, server EntitlementsServer) error {

	mux.Handle("GET", pattern_Entitlements_ListEntitlements_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/chainguard.platform.registry.Entitlements/ListEntitlements", runtime.WithHTTPPathPattern("/registry/v1/entitlements/{parent=**}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Entitlements_ListEntitlements_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Entitlements_ListEntitlements_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Entitlements_ListEntitlementImages_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/chainguard.platform.registry.Entitlements/ListEntitlementImages", runtime.WithHTTPPathPattern("/registry/v1/entitlements/{parent=**}/images"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Entitlements_ListEntitlementImages_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Entitlements_ListEntitlementImages_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Entitlements_Summary_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/chainguard.platform.registry.Entitlements/Summary", runtime.WithHTTPPathPattern("/registry/v1/entitlements/{parent=**}:summary"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Entitlements_Summary_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Entitlements_Summary_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterEntitlementsHandlerFromEndpoint is same as RegisterEntitlementsHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterEntitlementsHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.NewClient(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterEntitlementsHandler(ctx, mux, conn)
}

// RegisterEntitlementsHandler registers the http handlers for service Entitlements to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterEntitlementsHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterEntitlementsHandlerClient(ctx, mux, NewEntitlementsClient(conn))
}

// RegisterEntitlementsHandlerClient registers the http handlers for service Entitlements
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "EntitlementsClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "EntitlementsClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "EntitlementsClient" to call the correct interceptors. This client ignores the HTTP middlewares.
func RegisterEntitlementsHandlerClient(ctx context.Context, mux *runtime.ServeMux, client EntitlementsClient) error {

	mux.Handle("GET", pattern_Entitlements_ListEntitlements_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/chainguard.platform.registry.Entitlements/ListEntitlements", runtime.WithHTTPPathPattern("/registry/v1/entitlements/{parent=**}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Entitlements_ListEntitlements_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Entitlements_ListEntitlements_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Entitlements_ListEntitlementImages_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/chainguard.platform.registry.Entitlements/ListEntitlementImages", runtime.WithHTTPPathPattern("/registry/v1/entitlements/{parent=**}/images"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Entitlements_ListEntitlementImages_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Entitlements_ListEntitlementImages_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Entitlements_Summary_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/chainguard.platform.registry.Entitlements/Summary", runtime.WithHTTPPathPattern("/registry/v1/entitlements/{parent=**}:summary"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Entitlements_Summary_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Entitlements_Summary_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_Entitlements_ListEntitlements_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 3, 0, 4, 1, 5, 3}, []string{"registry", "v1", "entitlements", "parent"}, ""))

	pattern_Entitlements_ListEntitlementImages_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 3, 0, 4, 1, 5, 3, 2, 4}, []string{"registry", "v1", "entitlements", "parent", "images"}, ""))

	pattern_Entitlements_Summary_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 3, 0, 4, 1, 5, 3}, []string{"registry", "v1", "entitlements", "parent"}, "summary"))
)

var (
	forward_Entitlements_ListEntitlements_0 = runtime.ForwardResponseMessage

	forward_Entitlements_ListEntitlementImages_0 = runtime.ForwardResponseMessage

	forward_Entitlements_Summary_0 = runtime.ForwardResponseMessage
)
