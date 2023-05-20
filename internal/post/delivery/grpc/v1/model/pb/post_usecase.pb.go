// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: post_usecase.proto

package model

import (
	model "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/grpc/v1/model"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Posts struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page  *int64 `protobuf:"varint,1,opt,name=page,proto3,oneof" json:"page,omitempty"`
	Limit *int64 `protobuf:"varint,2,opt,name=limit,proto3,oneof" json:"limit,omitempty"`
}

func (x *Posts) Reset() {
	*x = Posts{}
	if protoimpl.UnsafeEnabled {
		mi := &file_post_usecase_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Posts) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Posts) ProtoMessage() {}

func (x *Posts) ProtoReflect() protoreflect.Message {
	mi := &file_post_usecase_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Posts.ProtoReflect.Descriptor instead.
func (*Posts) Descriptor() ([]byte, []int) {
	return file_post_usecase_proto_rawDescGZIP(), []int{0}
}

func (x *Posts) GetPage() int64 {
	if x != nil && x.Page != nil {
		return *x.Page
	}
	return 0
}

func (x *Posts) GetLimit() int64 {
	if x != nil && x.Limit != nil {
		return *x.Limit
	}
	return 0
}

type PostById struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PostID string `protobuf:"bytes,1,opt,name=PostID,proto3" json:"PostID,omitempty"`
}

func (x *PostById) Reset() {
	*x = PostById{}
	if protoimpl.UnsafeEnabled {
		mi := &file_post_usecase_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostById) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostById) ProtoMessage() {}

func (x *PostById) ProtoReflect() protoreflect.Message {
	mi := &file_post_usecase_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostById.ProtoReflect.Descriptor instead.
func (*PostById) Descriptor() ([]byte, []int) {
	return file_post_usecase_proto_rawDescGZIP(), []int{1}
}

func (x *PostById) GetPostID() string {
	if x != nil {
		return x.PostID
	}
	return ""
}

type PostDeleteView struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *PostDeleteView) Reset() {
	*x = PostDeleteView{}
	if protoimpl.UnsafeEnabled {
		mi := &file_post_usecase_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostDeleteView) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostDeleteView) ProtoMessage() {}

func (x *PostDeleteView) ProtoReflect() protoreflect.Message {
	mi := &file_post_usecase_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostDeleteView.ProtoReflect.Descriptor instead.
func (*PostDeleteView) Descriptor() ([]byte, []int) {
	return file_post_usecase_proto_rawDescGZIP(), []int{2}
}

func (x *PostDeleteView) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_post_usecase_proto protoreflect.FileDescriptor

var file_post_usecase_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x63, 0x61, 0x73, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x0a, 0x70, 0x6f, 0x73,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15,
	0x72, 0x70, 0x63, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4e, 0x0a, 0x05, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x12, 0x17,
	0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x01, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x88,
	0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x42, 0x08, 0x0a, 0x06, 0x5f,
	0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x22, 0x0a, 0x08, 0x50, 0x6f, 0x73, 0x74, 0x42, 0x79, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x6f, 0x73, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x50, 0x6f, 0x73, 0x74, 0x49, 0x44, 0x22, 0x2a, 0x0a, 0x0e, 0x50, 0x6f, 0x73,
	0x74, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x56, 0x69, 0x65, 0x77, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x32, 0x87, 0x02, 0x0a, 0x0b, 0x50, 0x6f, 0x73, 0x74, 0x55, 0x73,
	0x65, 0x43, 0x61, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x73, 0x74,
	0x12, 0x0f, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x42, 0x79, 0x49,
	0x64, 0x1a, 0x0f, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x56, 0x69,
	0x65, 0x77, 0x22, 0x00, 0x12, 0x29, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x73, 0x74, 0x73,
	0x12, 0x0c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x1a, 0x0b,
	0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x22, 0x00, 0x30, 0x01, 0x12,
	0x32, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x11, 0x2e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x1a, 0x0f, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x56, 0x69, 0x65,
	0x77, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x73,
	0x74, 0x12, 0x11, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x1a, 0x0f, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x50, 0x6f, 0x73,
	0x74, 0x56, 0x69, 0x65, 0x77, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x0f, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x50, 0x6f,
	0x73, 0x74, 0x42, 0x79, 0x49, 0x64, 0x1a, 0x15, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x50,
	0x6f, 0x73, 0x74, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x56, 0x69, 0x65, 0x77, 0x22, 0x00, 0x42,
	0x4f, 0x5a, 0x4d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x79, 0x61,
	0x63, 0x68, 0x6e, 0x79, 0x74, 0x73, 0x6b, 0x79, 0x69, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67,
	0x2d, 0x6d, 0x6f, 0x6e, 0x67, 0x6f, 0x2d, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65,
	0x72, 0x79, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_post_usecase_proto_rawDescOnce sync.Once
	file_post_usecase_proto_rawDescData = file_post_usecase_proto_rawDesc
)

func file_post_usecase_proto_rawDescGZIP() []byte {
	file_post_usecase_proto_rawDescOnce.Do(func() {
		file_post_usecase_proto_rawDescData = protoimpl.X.CompressGZIP(file_post_usecase_proto_rawDescData)
	})
	return file_post_usecase_proto_rawDescData
}

var file_post_usecase_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_post_usecase_proto_goTypes = []interface{}{
	(*Posts)(nil),            // 0: model.Posts
	(*PostById)(nil),         // 1: model.PostById
	(*PostDeleteView)(nil),   // 2: model.PostDeleteView
	(*model.PostCreate)(nil), // 3: model.PostCreate
	(*model.PostUpdate)(nil), // 4: model.PostUpdate
	(*model.PostView)(nil),   // 5: model.PostView
	(*model.Post)(nil),       // 6: model.Post
}
var file_post_usecase_proto_depIdxs = []int32{
	1, // 0: model.PostUseCase.GetPost:input_type -> model.PostById
	0, // 1: model.PostUseCase.GetPosts:input_type -> model.Posts
	3, // 2: model.PostUseCase.CreatePost:input_type -> model.PostCreate
	4, // 3: model.PostUseCase.UpdatePost:input_type -> model.PostUpdate
	1, // 4: model.PostUseCase.DeletePost:input_type -> model.PostById
	5, // 5: model.PostUseCase.GetPost:output_type -> model.PostView
	6, // 6: model.PostUseCase.GetPosts:output_type -> model.Post
	5, // 7: model.PostUseCase.CreatePost:output_type -> model.PostView
	5, // 8: model.PostUseCase.UpdatePost:output_type -> model.PostView
	2, // 9: model.PostUseCase.DeletePost:output_type -> model.PostDeleteView
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_post_usecase_proto_init() }
func file_post_usecase_proto_init() {
	if File_post_usecase_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_post_usecase_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Posts); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_post_usecase_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostById); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_post_usecase_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostDeleteView); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_post_usecase_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_post_usecase_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_post_usecase_proto_goTypes,
		DependencyIndexes: file_post_usecase_proto_depIdxs,
		MessageInfos:      file_post_usecase_proto_msgTypes,
	}.Build()
	File_post_usecase_proto = out.File
	file_post_usecase_proto_rawDesc = nil
	file_post_usecase_proto_goTypes = nil
	file_post_usecase_proto_depIdxs = nil
}
