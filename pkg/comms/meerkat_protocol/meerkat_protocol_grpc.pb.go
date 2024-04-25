// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: meerkat_protocol.proto

package __

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

const (
	MeerkatGuide_SayHello_FullMethodName               = "/meerkat_protocol.MeerkatGuide/SayHello"
	MeerkatGuide_EchoText_FullMethodName               = "/meerkat_protocol.MeerkatGuide/EchoText"
	MeerkatGuide_JoinPoolProtocol_FullMethodName       = "/meerkat_protocol.MeerkatGuide/JoinPoolProtocol"
	MeerkatGuide_HandshakePoolProtocol_FullMethodName  = "/meerkat_protocol.MeerkatGuide/HandshakePoolProtocol"
	MeerkatGuide_DisconnectPoolProtocol_FullMethodName = "/meerkat_protocol.MeerkatGuide/DisconnectPoolProtocol"
	MeerkatGuide_DataModProtocol_FullMethodName        = "/meerkat_protocol.MeerkatGuide/DataModProtocol"
)

// MeerkatGuideClient is the client API for MeerkatGuide service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MeerkatGuideClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
	EchoText(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error)
	JoinPoolProtocol(ctx context.Context, in *PoolJoinRequest, opts ...grpc.CallOption) (MeerkatGuide_JoinPoolProtocolClient, error)
	HandshakePoolProtocol(ctx context.Context, in *PoolHandshakesRequest, opts ...grpc.CallOption) (*PoolHandshakeResponse, error)
	DisconnectPoolProtocol(ctx context.Context, in *PoolDisconnectRequest, opts ...grpc.CallOption) (*PoolDisconnectResponse, error)
	DataModProtocol(ctx context.Context, opts ...grpc.CallOption) (MeerkatGuide_DataModProtocolClient, error)
}

type meerkatGuideClient struct {
	cc grpc.ClientConnInterface
}

func NewMeerkatGuideClient(cc grpc.ClientConnInterface) MeerkatGuideClient {
	return &meerkatGuideClient{cc}
}

func (c *meerkatGuideClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, MeerkatGuide_SayHello_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *meerkatGuideClient) EchoText(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error) {
	out := new(EchoResponse)
	err := c.cc.Invoke(ctx, MeerkatGuide_EchoText_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *meerkatGuideClient) JoinPoolProtocol(ctx context.Context, in *PoolJoinRequest, opts ...grpc.CallOption) (MeerkatGuide_JoinPoolProtocolClient, error) {
	stream, err := c.cc.NewStream(ctx, &MeerkatGuide_ServiceDesc.Streams[0], MeerkatGuide_JoinPoolProtocol_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &meerkatGuideJoinPoolProtocolClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MeerkatGuide_JoinPoolProtocolClient interface {
	Recv() (*PoolJoinResponse, error)
	grpc.ClientStream
}

type meerkatGuideJoinPoolProtocolClient struct {
	grpc.ClientStream
}

func (x *meerkatGuideJoinPoolProtocolClient) Recv() (*PoolJoinResponse, error) {
	m := new(PoolJoinResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *meerkatGuideClient) HandshakePoolProtocol(ctx context.Context, in *PoolHandshakesRequest, opts ...grpc.CallOption) (*PoolHandshakeResponse, error) {
	out := new(PoolHandshakeResponse)
	err := c.cc.Invoke(ctx, MeerkatGuide_HandshakePoolProtocol_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *meerkatGuideClient) DisconnectPoolProtocol(ctx context.Context, in *PoolDisconnectRequest, opts ...grpc.CallOption) (*PoolDisconnectResponse, error) {
	out := new(PoolDisconnectResponse)
	err := c.cc.Invoke(ctx, MeerkatGuide_DisconnectPoolProtocol_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *meerkatGuideClient) DataModProtocol(ctx context.Context, opts ...grpc.CallOption) (MeerkatGuide_DataModProtocolClient, error) {
	stream, err := c.cc.NewStream(ctx, &MeerkatGuide_ServiceDesc.Streams[1], MeerkatGuide_DataModProtocol_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &meerkatGuideDataModProtocolClient{stream}
	return x, nil
}

type MeerkatGuide_DataModProtocolClient interface {
	Send(*DataModRequest) error
	CloseAndRecv() (*DataModResponse, error)
	grpc.ClientStream
}

type meerkatGuideDataModProtocolClient struct {
	grpc.ClientStream
}

func (x *meerkatGuideDataModProtocolClient) Send(m *DataModRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *meerkatGuideDataModProtocolClient) CloseAndRecv() (*DataModResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(DataModResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MeerkatGuideServer is the server API for MeerkatGuide service.
// All implementations must embed UnimplementedMeerkatGuideServer
// for forward compatibility
type MeerkatGuideServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloResponse, error)
	EchoText(context.Context, *EchoRequest) (*EchoResponse, error)
	JoinPoolProtocol(*PoolJoinRequest, MeerkatGuide_JoinPoolProtocolServer) error
	HandshakePoolProtocol(context.Context, *PoolHandshakesRequest) (*PoolHandshakeResponse, error)
	DisconnectPoolProtocol(context.Context, *PoolDisconnectRequest) (*PoolDisconnectResponse, error)
	DataModProtocol(MeerkatGuide_DataModProtocolServer) error
	mustEmbedUnimplementedMeerkatGuideServer()
}

// UnimplementedMeerkatGuideServer must be embedded to have forward compatible implementations.
type UnimplementedMeerkatGuideServer struct {
}

func (UnimplementedMeerkatGuideServer) SayHello(context.Context, *HelloRequest) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedMeerkatGuideServer) EchoText(context.Context, *EchoRequest) (*EchoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EchoText not implemented")
}
func (UnimplementedMeerkatGuideServer) JoinPoolProtocol(*PoolJoinRequest, MeerkatGuide_JoinPoolProtocolServer) error {
	return status.Errorf(codes.Unimplemented, "method JoinPoolProtocol not implemented")
}
func (UnimplementedMeerkatGuideServer) HandshakePoolProtocol(context.Context, *PoolHandshakesRequest) (*PoolHandshakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandshakePoolProtocol not implemented")
}
func (UnimplementedMeerkatGuideServer) DisconnectPoolProtocol(context.Context, *PoolDisconnectRequest) (*PoolDisconnectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisconnectPoolProtocol not implemented")
}
func (UnimplementedMeerkatGuideServer) DataModProtocol(MeerkatGuide_DataModProtocolServer) error {
	return status.Errorf(codes.Unimplemented, "method DataModProtocol not implemented")
}
func (UnimplementedMeerkatGuideServer) mustEmbedUnimplementedMeerkatGuideServer() {}

// UnsafeMeerkatGuideServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MeerkatGuideServer will
// result in compilation errors.
type UnsafeMeerkatGuideServer interface {
	mustEmbedUnimplementedMeerkatGuideServer()
}

func RegisterMeerkatGuideServer(s grpc.ServiceRegistrar, srv MeerkatGuideServer) {
	s.RegisterService(&MeerkatGuide_ServiceDesc, srv)
}

func _MeerkatGuide_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MeerkatGuideServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MeerkatGuide_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MeerkatGuideServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MeerkatGuide_EchoText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MeerkatGuideServer).EchoText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MeerkatGuide_EchoText_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MeerkatGuideServer).EchoText(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MeerkatGuide_JoinPoolProtocol_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PoolJoinRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MeerkatGuideServer).JoinPoolProtocol(m, &meerkatGuideJoinPoolProtocolServer{stream})
}

type MeerkatGuide_JoinPoolProtocolServer interface {
	Send(*PoolJoinResponse) error
	grpc.ServerStream
}

type meerkatGuideJoinPoolProtocolServer struct {
	grpc.ServerStream
}

func (x *meerkatGuideJoinPoolProtocolServer) Send(m *PoolJoinResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _MeerkatGuide_HandshakePoolProtocol_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PoolHandshakesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MeerkatGuideServer).HandshakePoolProtocol(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MeerkatGuide_HandshakePoolProtocol_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MeerkatGuideServer).HandshakePoolProtocol(ctx, req.(*PoolHandshakesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MeerkatGuide_DisconnectPoolProtocol_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PoolDisconnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MeerkatGuideServer).DisconnectPoolProtocol(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MeerkatGuide_DisconnectPoolProtocol_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MeerkatGuideServer).DisconnectPoolProtocol(ctx, req.(*PoolDisconnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MeerkatGuide_DataModProtocol_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MeerkatGuideServer).DataModProtocol(&meerkatGuideDataModProtocolServer{stream})
}

type MeerkatGuide_DataModProtocolServer interface {
	SendAndClose(*DataModResponse) error
	Recv() (*DataModRequest, error)
	grpc.ServerStream
}

type meerkatGuideDataModProtocolServer struct {
	grpc.ServerStream
}

func (x *meerkatGuideDataModProtocolServer) SendAndClose(m *DataModResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *meerkatGuideDataModProtocolServer) Recv() (*DataModRequest, error) {
	m := new(DataModRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MeerkatGuide_ServiceDesc is the grpc.ServiceDesc for MeerkatGuide service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MeerkatGuide_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "meerkat_protocol.MeerkatGuide",
	HandlerType: (*MeerkatGuideServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _MeerkatGuide_SayHello_Handler,
		},
		{
			MethodName: "EchoText",
			Handler:    _MeerkatGuide_EchoText_Handler,
		},
		{
			MethodName: "HandshakePoolProtocol",
			Handler:    _MeerkatGuide_HandshakePoolProtocol_Handler,
		},
		{
			MethodName: "DisconnectPoolProtocol",
			Handler:    _MeerkatGuide_DisconnectPoolProtocol_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "JoinPoolProtocol",
			Handler:       _MeerkatGuide_JoinPoolProtocol_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "DataModProtocol",
			Handler:       _MeerkatGuide_DataModProtocol_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "meerkat_protocol.proto",
}
