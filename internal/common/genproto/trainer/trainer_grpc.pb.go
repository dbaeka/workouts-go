// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: trainer.proto

package trainer

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

// TrainerServiceClient is the client API for TrainerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TrainerServiceClient interface {
	IsHourAvailable(ctx context.Context, in *IsHourAvailableRequest, opts ...grpc.CallOption) (*IsHourAvailableResponse, error)
	UpdateHour(ctx context.Context, in *UpdateHourRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
}

type trainerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTrainerServiceClient(cc grpc.ClientConnInterface) TrainerServiceClient {
	return &trainerServiceClient{cc}
}

func (c *trainerServiceClient) IsHourAvailable(ctx context.Context, in *IsHourAvailableRequest, opts ...grpc.CallOption) (*IsHourAvailableResponse, error) {
	out := new(IsHourAvailableResponse)
	err := c.cc.Invoke(ctx, "/trainer.TrainerService/IsHourAvailable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trainerServiceClient) UpdateHour(ctx context.Context, in *UpdateHourRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/trainer.TrainerService/UpdateHour", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TrainerServiceServer is the server API for TrainerService service.
// All implementations must embed UnimplementedTrainerServiceServer
// for forward compatibility
type TrainerServiceServer interface {
	IsHourAvailable(context.Context, *IsHourAvailableRequest) (*IsHourAvailableResponse, error)
	UpdateHour(context.Context, *UpdateHourRequest) (*EmptyResponse, error)
	mustEmbedUnimplementedTrainerServiceServer()
}

// UnimplementedTrainerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTrainerServiceServer struct {
}

func (UnimplementedTrainerServiceServer) IsHourAvailable(context.Context, *IsHourAvailableRequest) (*IsHourAvailableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsHourAvailable not implemented")
}
func (UnimplementedTrainerServiceServer) UpdateHour(context.Context, *UpdateHourRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateHour not implemented")
}
func (UnimplementedTrainerServiceServer) mustEmbedUnimplementedTrainerServiceServer() {}

// UnsafeTrainerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TrainerServiceServer will
// result in compilation errors.
type UnsafeTrainerServiceServer interface {
	mustEmbedUnimplementedTrainerServiceServer()
}

func RegisterTrainerServiceServer(s grpc.ServiceRegistrar, srv TrainerServiceServer) {
	s.RegisterService(&TrainerService_ServiceDesc, srv)
}

func _TrainerService_IsHourAvailable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsHourAvailableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrainerServiceServer).IsHourAvailable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trainer.TrainerService/IsHourAvailable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrainerServiceServer).IsHourAvailable(ctx, req.(*IsHourAvailableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrainerService_UpdateHour_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateHourRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrainerServiceServer).UpdateHour(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trainer.TrainerService/UpdateHour",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrainerServiceServer).UpdateHour(ctx, req.(*UpdateHourRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TrainerService_ServiceDesc is the grpc.ServiceDesc for TrainerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TrainerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "trainer.TrainerService",
	HandlerType: (*TrainerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsHourAvailable",
			Handler:    _TrainerService_IsHourAvailable_Handler,
		},
		{
			MethodName: "UpdateHour",
			Handler:    _TrainerService_UpdateHour_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "trainer.proto",
}
