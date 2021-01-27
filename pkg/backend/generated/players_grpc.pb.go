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
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PlayersServiceClient is the client API for PlayersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlayersServiceClient interface {
	Search(ctx context.Context, in *SearchPlayersRequest, opts ...grpc.CallOption) (*PlayersElasticResponse, error)
	List(ctx context.Context, in *ListPlayersRequest, opts ...grpc.CallOption) (*PlayersMongoResponse, error)
}

type playersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPlayersServiceClient(cc grpc.ClientConnInterface) PlayersServiceClient {
	return &playersServiceClient{cc}
}

func (c *playersServiceClient) Search(ctx context.Context, in *SearchPlayersRequest, opts ...grpc.CallOption) (*PlayersElasticResponse, error) {
	out := new(PlayersElasticResponse)
	err := c.cc.Invoke(ctx, "/generated.PlayersService/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playersServiceClient) List(ctx context.Context, in *ListPlayersRequest, opts ...grpc.CallOption) (*PlayersMongoResponse, error) {
	out := new(PlayersMongoResponse)
	err := c.cc.Invoke(ctx, "/generated.PlayersService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlayersServiceServer is the server API for PlayersService service.
// All implementations must embed UnimplementedPlayersServiceServer
// for forward compatibility
type PlayersServiceServer interface {
	Search(context.Context, *SearchPlayersRequest) (*PlayersElasticResponse, error)
	List(context.Context, *ListPlayersRequest) (*PlayersMongoResponse, error)
	mustEmbedUnimplementedPlayersServiceServer()
}

// UnimplementedPlayersServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPlayersServiceServer struct {
}

func (UnimplementedPlayersServiceServer) Search(context.Context, *SearchPlayersRequest) (*PlayersElasticResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedPlayersServiceServer) List(context.Context, *ListPlayersRequest) (*PlayersMongoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedPlayersServiceServer) mustEmbedUnimplementedPlayersServiceServer() {}

// UnsafePlayersServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlayersServiceServer will
// result in compilation errors.
type UnsafePlayersServiceServer interface {
	mustEmbedUnimplementedPlayersServiceServer()
}

func RegisterPlayersServiceServer(s grpc.ServiceRegistrar, srv PlayersServiceServer) {
	s.RegisterService(&PlayersService_ServiceDesc, srv)
}

func _PlayersService_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchPlayersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayersServiceServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.PlayersService/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayersServiceServer).Search(ctx, req.(*SearchPlayersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlayersService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPlayersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayersServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.PlayersService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayersServiceServer).List(ctx, req.(*ListPlayersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PlayersService_ServiceDesc is the grpc.ServiceDesc for PlayersService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PlayersService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "generated.PlayersService",
	HandlerType: (*PlayersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _PlayersService_Search_Handler,
		},
		{
			MethodName: "List",
			Handler:    _PlayersService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "players.proto",
}
