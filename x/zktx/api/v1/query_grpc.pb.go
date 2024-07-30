// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: hyle/zktx/v1/query.proto

package zktxv1

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
	Query_Contract_FullMethodName         = "/hyle.zktx.v1.Query/Contract"
	Query_ContractList_FullMethodName     = "/hyle.zktx.v1.Query/ContractList"
	Query_SettlementStatus_FullMethodName = "/hyle.zktx.v1.Query/SettlementStatus"
	Query_Params_FullMethodName           = "/hyle.zktx.v1.Query/Params"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Contract returns the current state of the contract.
	Contract(ctx context.Context, in *ContractRequest, opts ...grpc.CallOption) (*ContractResponse, error)
	// ContractList returns the list of all contracts with a given verifier and
	// program_id
	ContractList(ctx context.Context, in *ContractListRequest, opts ...grpc.CallOption) (*ContractListResponse, error)
	// SettlementStatus returns whether a TX has been settled or not
	SettlementStatus(ctx context.Context, in *SettlementStatusRequest, opts ...grpc.CallOption) (*SettlementStatusResponse, error)
	// Params returns the module parameters.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Contract(ctx context.Context, in *ContractRequest, opts ...grpc.CallOption) (*ContractResponse, error) {
	out := new(ContractResponse)
	err := c.cc.Invoke(ctx, Query_Contract_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ContractList(ctx context.Context, in *ContractListRequest, opts ...grpc.CallOption) (*ContractListResponse, error) {
	out := new(ContractListResponse)
	err := c.cc.Invoke(ctx, Query_ContractList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) SettlementStatus(ctx context.Context, in *SettlementStatusRequest, opts ...grpc.CallOption) (*SettlementStatusResponse, error) {
	out := new(SettlementStatusResponse)
	err := c.cc.Invoke(ctx, Query_SettlementStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Contract returns the current state of the contract.
	Contract(context.Context, *ContractRequest) (*ContractResponse, error)
	// ContractList returns the list of all contracts with a given verifier and
	// program_id
	ContractList(context.Context, *ContractListRequest) (*ContractListResponse, error)
	// SettlementStatus returns whether a TX has been settled or not
	SettlementStatus(context.Context, *SettlementStatusRequest) (*SettlementStatusResponse, error)
	// Params returns the module parameters.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Contract(context.Context, *ContractRequest) (*ContractResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Contract not implemented")
}
func (UnimplementedQueryServer) ContractList(context.Context, *ContractListRequest) (*ContractListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ContractList not implemented")
}
func (UnimplementedQueryServer) SettlementStatus(context.Context, *SettlementStatusRequest) (*SettlementStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SettlementStatus not implemented")
}
func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Contract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ContractRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Contract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Contract_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Contract(ctx, req.(*ContractRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ContractList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ContractListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ContractList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ContractList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ContractList(ctx, req.(*ContractListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_SettlementStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SettlementStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).SettlementStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_SettlementStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).SettlementStatus(ctx, req.(*SettlementStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Params_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hyle.zktx.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Contract",
			Handler:    _Query_Contract_Handler,
		},
		{
			MethodName: "ContractList",
			Handler:    _Query_ContractList_Handler,
		},
		{
			MethodName: "SettlementStatus",
			Handler:    _Query_SettlementStatus_Handler,
		},
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hyle/zktx/v1/query.proto",
}
