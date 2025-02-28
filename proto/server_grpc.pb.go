// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.13.0
// source: proto/server.proto

package managerpb

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

// ShardServiceClient is the client API for ShardService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShardServiceClient interface {
	GetRecordsForNode(ctx context.Context, in *NodeRequest, opts ...grpc.CallOption) (*RecordsResponse, error)
	SubscribeRebalance(ctx context.Context, in *NodeRequest, opts ...grpc.CallOption) (ShardService_SubscribeRebalanceClient, error)
	RegisterNode(ctx context.Context, in *NodeRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	RemoveRecord(ctx context.Context, in *RemoveRecordRequest, opts ...grpc.CallOption) (*RemoveRecordResponse, error)
}

type shardServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewShardServiceClient(cc grpc.ClientConnInterface) ShardServiceClient {
	return &shardServiceClient{cc}
}

func (c *shardServiceClient) GetRecordsForNode(ctx context.Context, in *NodeRequest, opts ...grpc.CallOption) (*RecordsResponse, error) {
	out := new(RecordsResponse)
	err := c.cc.Invoke(ctx, "/shard.ShardService/GetRecordsForNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shardServiceClient) SubscribeRebalance(ctx context.Context, in *NodeRequest, opts ...grpc.CallOption) (ShardService_SubscribeRebalanceClient, error) {
	stream, err := c.cc.NewStream(ctx, &ShardService_ServiceDesc.Streams[0], "/shard.ShardService/SubscribeRebalance", opts...)
	if err != nil {
		return nil, err
	}
	x := &shardServiceSubscribeRebalanceClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ShardService_SubscribeRebalanceClient interface {
	Recv() (*RebalanceEvent, error)
	grpc.ClientStream
}

type shardServiceSubscribeRebalanceClient struct {
	grpc.ClientStream
}

func (x *shardServiceSubscribeRebalanceClient) Recv() (*RebalanceEvent, error) {
	m := new(RebalanceEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *shardServiceClient) RegisterNode(ctx context.Context, in *NodeRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/shard.ShardService/RegisterNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shardServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/shard.ShardService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shardServiceClient) RemoveRecord(ctx context.Context, in *RemoveRecordRequest, opts ...grpc.CallOption) (*RemoveRecordResponse, error) {
	out := new(RemoveRecordResponse)
	err := c.cc.Invoke(ctx, "/shard.ShardService/RemoveRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShardServiceServer is the server API for ShardService service.
// All implementations must embed UnimplementedShardServiceServer
// for forward compatibility
type ShardServiceServer interface {
	GetRecordsForNode(context.Context, *NodeRequest) (*RecordsResponse, error)
	SubscribeRebalance(*NodeRequest, ShardService_SubscribeRebalanceServer) error
	RegisterNode(context.Context, *NodeRequest) (*RegisterResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	RemoveRecord(context.Context, *RemoveRecordRequest) (*RemoveRecordResponse, error)
	mustEmbedUnimplementedShardServiceServer()
}

// UnimplementedShardServiceServer must be embedded to have forward compatible implementations.
type UnimplementedShardServiceServer struct {
}

func (UnimplementedShardServiceServer) GetRecordsForNode(context.Context, *NodeRequest) (*RecordsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecordsForNode not implemented")
}
func (UnimplementedShardServiceServer) SubscribeRebalance(*NodeRequest, ShardService_SubscribeRebalanceServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeRebalance not implemented")
}
func (UnimplementedShardServiceServer) RegisterNode(context.Context, *NodeRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterNode not implemented")
}
func (UnimplementedShardServiceServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedShardServiceServer) RemoveRecord(context.Context, *RemoveRecordRequest) (*RemoveRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveRecord not implemented")
}
func (UnimplementedShardServiceServer) mustEmbedUnimplementedShardServiceServer() {}

// UnsafeShardServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShardServiceServer will
// result in compilation errors.
type UnsafeShardServiceServer interface {
	mustEmbedUnimplementedShardServiceServer()
}

func RegisterShardServiceServer(s grpc.ServiceRegistrar, srv ShardServiceServer) {
	s.RegisterService(&ShardService_ServiceDesc, srv)
}

func _ShardService_GetRecordsForNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShardServiceServer).GetRecordsForNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shard.ShardService/GetRecordsForNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShardServiceServer).GetRecordsForNode(ctx, req.(*NodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShardService_SubscribeRebalance_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(NodeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ShardServiceServer).SubscribeRebalance(m, &shardServiceSubscribeRebalanceServer{stream})
}

type ShardService_SubscribeRebalanceServer interface {
	Send(*RebalanceEvent) error
	grpc.ServerStream
}

type shardServiceSubscribeRebalanceServer struct {
	grpc.ServerStream
}

func (x *shardServiceSubscribeRebalanceServer) Send(m *RebalanceEvent) error {
	return x.ServerStream.SendMsg(m)
}

func _ShardService_RegisterNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShardServiceServer).RegisterNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shard.ShardService/RegisterNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShardServiceServer).RegisterNode(ctx, req.(*NodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShardService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShardServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shard.ShardService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShardServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShardService_RemoveRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShardServiceServer).RemoveRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shard.ShardService/RemoveRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShardServiceServer).RemoveRecord(ctx, req.(*RemoveRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ShardService_ServiceDesc is the grpc.ServiceDesc for ShardService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShardService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shard.ShardService",
	HandlerType: (*ShardServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRecordsForNode",
			Handler:    _ShardService_GetRecordsForNode_Handler,
		},
		{
			MethodName: "RegisterNode",
			Handler:    _ShardService_RegisterNode_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _ShardService_Ping_Handler,
		},
		{
			MethodName: "RemoveRecord",
			Handler:    _ShardService_RemoveRecord_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeRebalance",
			Handler:       _ShardService_SubscribeRebalance_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/server.proto",
}
