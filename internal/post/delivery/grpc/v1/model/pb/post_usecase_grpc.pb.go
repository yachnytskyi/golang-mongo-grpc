// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: post_usecase.proto

package model

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
	PostUseCase_GetPostById_FullMethodName    = "/model.PostUseCase/GetPostById"
	PostUseCase_GetAllPosts_FullMethodName    = "/model.PostUseCase/GetAllPosts"
	PostUseCase_CreatePost_FullMethodName     = "/model.PostUseCase/CreatePost"
	PostUseCase_UpdatePostById_FullMethodName = "/model.PostUseCase/UpdatePostById"
	PostUseCase_DeletePostById_FullMethodName = "/model.PostUseCase/DeletePostById"
)

// PostUseCaseClient is the client API for PostUseCase service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostUseCaseClient interface {
	GetPostById(ctx context.Context, in *PostById, opts ...grpc.CallOption) (*PostView, error)
	GetAllPosts(ctx context.Context, in *Posts, opts ...grpc.CallOption) (PostUseCase_GetAllPostsClient, error)
	CreatePost(ctx context.Context, in *PostCreate, opts ...grpc.CallOption) (*PostView, error)
	UpdatePostById(ctx context.Context, in *PostUpdate, opts ...grpc.CallOption) (*PostView, error)
	DeletePostById(ctx context.Context, in *PostById, opts ...grpc.CallOption) (*PostDeleteView, error)
}

type postUseCaseClient struct {
	cc grpc.ClientConnInterface
}

func NewPostUseCaseClient(cc grpc.ClientConnInterface) PostUseCaseClient {
	return &postUseCaseClient{cc}
}

func (c *postUseCaseClient) GetPostById(ctx context.Context, in *PostById, opts ...grpc.CallOption) (*PostView, error) {
	out := new(PostView)
	err := c.cc.Invoke(ctx, PostUseCase_GetPostById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postUseCaseClient) GetAllPosts(ctx context.Context, in *Posts, opts ...grpc.CallOption) (PostUseCase_GetAllPostsClient, error) {
	stream, err := c.cc.NewStream(ctx, &PostUseCase_ServiceDesc.Streams[0], PostUseCase_GetAllPosts_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &postUseCaseGetAllPostsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PostUseCase_GetAllPostsClient interface {
	Recv() (*Post, error)
	grpc.ClientStream
}

type postUseCaseGetAllPostsClient struct {
	grpc.ClientStream
}

func (x *postUseCaseGetAllPostsClient) Recv() (*Post, error) {
	m := new(Post)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *postUseCaseClient) CreatePost(ctx context.Context, in *PostCreate, opts ...grpc.CallOption) (*PostView, error) {
	out := new(PostView)
	err := c.cc.Invoke(ctx, PostUseCase_CreatePost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postUseCaseClient) UpdatePostById(ctx context.Context, in *PostUpdate, opts ...grpc.CallOption) (*PostView, error) {
	out := new(PostView)
	err := c.cc.Invoke(ctx, PostUseCase_UpdatePostById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postUseCaseClient) DeletePostById(ctx context.Context, in *PostById, opts ...grpc.CallOption) (*PostDeleteView, error) {
	out := new(PostDeleteView)
	err := c.cc.Invoke(ctx, PostUseCase_DeletePostById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostUseCaseServer is the server API for PostUseCase service.
// All implementations must embed UnimplementedPostUseCaseServer
// for forward compatibility
type PostUseCaseServer interface {
	GetPostById(context.Context, *PostById) (*PostView, error)
	GetAllPosts(*Posts, PostUseCase_GetAllPostsServer) error
	CreatePost(context.Context, *PostCreate) (*PostView, error)
	UpdatePostById(context.Context, *PostUpdate) (*PostView, error)
	DeletePostById(context.Context, *PostById) (*PostDeleteView, error)
	mustEmbedUnimplementedPostUseCaseServer()
}

// UnimplementedPostUseCaseServer must be embedded to have forward compatible implementations.
type UnimplementedPostUseCaseServer struct {
}

func (UnimplementedPostUseCaseServer) GetPostById(context.Context, *PostById) (*PostView, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPostById not implemented")
}
func (UnimplementedPostUseCaseServer) GetAllPosts(*Posts, PostUseCase_GetAllPostsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAllPosts not implemented")
}
func (UnimplementedPostUseCaseServer) CreatePost(context.Context, *PostCreate) (*PostView, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePost not implemented")
}
func (UnimplementedPostUseCaseServer) UpdatePostById(context.Context, *PostUpdate) (*PostView, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePostById not implemented")
}
func (UnimplementedPostUseCaseServer) DeletePostById(context.Context, *PostById) (*PostDeleteView, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePostById not implemented")
}
func (UnimplementedPostUseCaseServer) mustEmbedUnimplementedPostUseCaseServer() {}

// UnsafePostUseCaseServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostUseCaseServer will
// result in compilation errors.
type UnsafePostUseCaseServer interface {
	mustEmbedUnimplementedPostUseCaseServer()
}

func RegisterPostUseCaseServer(s grpc.ServiceRegistrar, srv PostUseCaseServer) {
	s.RegisterService(&PostUseCase_ServiceDesc, srv)
}

func _PostUseCase_GetPostById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostById)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostUseCaseServer).GetPostById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostUseCase_GetPostById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostUseCaseServer).GetPostById(ctx, req.(*PostById))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostUseCase_GetAllPosts_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Posts)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PostUseCaseServer).GetAllPosts(m, &postUseCaseGetAllPostsServer{stream})
}

type PostUseCase_GetAllPostsServer interface {
	Send(*Post) error
	grpc.ServerStream
}

type postUseCaseGetAllPostsServer struct {
	grpc.ServerStream
}

func (x *postUseCaseGetAllPostsServer) Send(m *Post) error {
	return x.ServerStream.SendMsg(m)
}

func _PostUseCase_CreatePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostCreate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostUseCaseServer).CreatePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostUseCase_CreatePost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostUseCaseServer).CreatePost(ctx, req.(*PostCreate))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostUseCase_UpdatePostById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostUpdate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostUseCaseServer).UpdatePostById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostUseCase_UpdatePostById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostUseCaseServer).UpdatePostById(ctx, req.(*PostUpdate))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostUseCase_DeletePostById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostById)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostUseCaseServer).DeletePostById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostUseCase_DeletePostById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostUseCaseServer).DeletePostById(ctx, req.(*PostById))
	}
	return interceptor(ctx, in, info, handler)
}

// PostUseCase_ServiceDesc is the grpc.ServiceDesc for PostUseCase service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PostUseCase_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "model.PostUseCase",
	HandlerType: (*PostUseCaseServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPostById",
			Handler:    _PostUseCase_GetPostById_Handler,
		},
		{
			MethodName: "CreatePost",
			Handler:    _PostUseCase_CreatePost_Handler,
		},
		{
			MethodName: "UpdatePostById",
			Handler:    _PostUseCase_UpdatePostById_Handler,
		},
		{
			MethodName: "DeletePostById",
			Handler:    _PostUseCase_DeletePostById_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAllPosts",
			Handler:       _PostUseCase_GetAllPosts_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "post_usecase.proto",
}