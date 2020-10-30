// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package generated

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// AppsServiceClient is the client API for AppsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AppsServiceClient interface {
	Search(ctx context.Context, in *SearchAppsRequest, opts ...grpc.CallOption) (*AppsElasticResponse, error)
	Apps(ctx context.Context, in *ListAppsRequest, opts ...grpc.CallOption) (*AppsMongoResponse, error)
}

type appsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAppsServiceClient(cc grpc.ClientConnInterface) AppsServiceClient {
	return &appsServiceClient{cc}
}

func (c *appsServiceClient) Search(ctx context.Context, in *SearchAppsRequest, opts ...grpc.CallOption) (*AppsElasticResponse, error) {
	out := new(AppsElasticResponse)
	err := c.cc.Invoke(ctx, "/generated.AppsService/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appsServiceClient) Apps(ctx context.Context, in *ListAppsRequest, opts ...grpc.CallOption) (*AppsMongoResponse, error) {
	out := new(AppsMongoResponse)
	err := c.cc.Invoke(ctx, "/generated.AppsService/Apps", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AppsServiceServer is the server API for AppsService service.
// All implementations must embed UnimplementedAppsServiceServer
// for forward compatibility
type AppsServiceServer interface {
	Search(context.Context, *SearchAppsRequest) (*AppsElasticResponse, error)
	Apps(context.Context, *ListAppsRequest) (*AppsMongoResponse, error)
	mustEmbedUnimplementedAppsServiceServer()
}

// UnimplementedAppsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAppsServiceServer struct {
}

func (UnimplementedAppsServiceServer) Search(context.Context, *SearchAppsRequest) (*AppsElasticResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedAppsServiceServer) Apps(context.Context, *ListAppsRequest) (*AppsMongoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Apps not implemented")
}
func (UnimplementedAppsServiceServer) mustEmbedUnimplementedAppsServiceServer() {}

// UnsafeAppsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AppsServiceServer will
// result in compilation errors.
type UnsafeAppsServiceServer interface {
	mustEmbedUnimplementedAppsServiceServer()
}

func RegisterAppsServiceServer(s grpc.ServiceRegistrar, srv AppsServiceServer) {
	s.RegisterService(&_AppsService_serviceDesc, srv)
}

func _AppsService_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchAppsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppsServiceServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.AppsService/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppsServiceServer).Search(ctx, req.(*SearchAppsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AppsService_Apps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAppsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppsServiceServer).Apps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.AppsService/Apps",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppsServiceServer).Apps(ctx, req.(*ListAppsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AppsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "generated.AppsService",
	HandlerType: (*AppsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _AppsService_Search_Handler,
		},
		{
			MethodName: "Apps",
			Handler:    _AppsService_Apps_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apps.proto",
}
