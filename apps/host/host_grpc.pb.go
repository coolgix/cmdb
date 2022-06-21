// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: apps/host/pb/host.proto

package host

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

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	// 同步云商的主机资源
	SyncHost(ctx context.Context, in *Host, opts ...grpc.CallOption) (*Host, error)
	// 查询本地同步后的主机资源列表
	QueryHost(ctx context.Context, in *QueryHostRequest, opts ...grpc.CallOption) (*HostSet, error)
	// 查询主机详情信息,查询单个
	DescribeHost(ctx context.Context, in *DescribeHostRequest, opts ...grpc.CallOption) (*Host, error)
	// 更新主机信息, 同步更新云商资源信息
	UpdateHost(ctx context.Context, in *UpdateHostRequest, opts ...grpc.CallOption) (*Host, error)
	// 释放主机, 按计划释放后, 信息会保留一段时间
	ReleaseHost(ctx context.Context, in *ReleaseHostRequest, opts ...grpc.CallOption) (*Host, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) SyncHost(ctx context.Context, in *Host, opts ...grpc.CallOption) (*Host, error) {
	out := new(Host)
	err := c.cc.Invoke(ctx, "/course.cmdb.host.Service/SyncHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) QueryHost(ctx context.Context, in *QueryHostRequest, opts ...grpc.CallOption) (*HostSet, error) {
	out := new(HostSet)
	err := c.cc.Invoke(ctx, "/course.cmdb.host.Service/QueryHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DescribeHost(ctx context.Context, in *DescribeHostRequest, opts ...grpc.CallOption) (*Host, error) {
	out := new(Host)
	err := c.cc.Invoke(ctx, "/course.cmdb.host.Service/DescribeHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) UpdateHost(ctx context.Context, in *UpdateHostRequest, opts ...grpc.CallOption) (*Host, error) {
	out := new(Host)
	err := c.cc.Invoke(ctx, "/course.cmdb.host.Service/UpdateHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ReleaseHost(ctx context.Context, in *ReleaseHostRequest, opts ...grpc.CallOption) (*Host, error) {
	out := new(Host)
	err := c.cc.Invoke(ctx, "/course.cmdb.host.Service/ReleaseHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	// 同步云商的主机资源
	SyncHost(context.Context, *Host) (*Host, error)
	// 查询本地同步后的主机资源列表
	QueryHost(context.Context, *QueryHostRequest) (*HostSet, error)
	// 查询主机详情信息,查询单个
	DescribeHost(context.Context, *DescribeHostRequest) (*Host, error)
	// 更新主机信息, 同步更新云商资源信息
	UpdateHost(context.Context, *UpdateHostRequest) (*Host, error)
	// 释放主机, 按计划释放后, 信息会保留一段时间
	ReleaseHost(context.Context, *ReleaseHostRequest) (*Host, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) SyncHost(context.Context, *Host) (*Host, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncHost not implemented")
}
func (UnimplementedServiceServer) QueryHost(context.Context, *QueryHostRequest) (*HostSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryHost not implemented")
}
func (UnimplementedServiceServer) DescribeHost(context.Context, *DescribeHostRequest) (*Host, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeHost not implemented")
}
func (UnimplementedServiceServer) UpdateHost(context.Context, *UpdateHostRequest) (*Host, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateHost not implemented")
}
func (UnimplementedServiceServer) ReleaseHost(context.Context, *ReleaseHostRequest) (*Host, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseHost not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_SyncHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Host)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).SyncHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/course.cmdb.host.Service/SyncHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).SyncHost(ctx, req.(*Host))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_QueryHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).QueryHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/course.cmdb.host.Service/QueryHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).QueryHost(ctx, req.(*QueryHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DescribeHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DescribeHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/course.cmdb.host.Service/DescribeHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DescribeHost(ctx, req.(*DescribeHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_UpdateHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).UpdateHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/course.cmdb.host.Service/UpdateHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).UpdateHost(ctx, req.(*UpdateHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ReleaseHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ReleaseHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/course.cmdb.host.Service/ReleaseHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ReleaseHost(ctx, req.(*ReleaseHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "course.cmdb.host.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SyncHost",
			Handler:    _Service_SyncHost_Handler,
		},
		{
			MethodName: "QueryHost",
			Handler:    _Service_QueryHost_Handler,
		},
		{
			MethodName: "DescribeHost",
			Handler:    _Service_DescribeHost_Handler,
		},
		{
			MethodName: "UpdateHost",
			Handler:    _Service_UpdateHost_Handler,
		},
		{
			MethodName: "ReleaseHost",
			Handler:    _Service_ReleaseHost_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apps/host/pb/host.proto",
}
