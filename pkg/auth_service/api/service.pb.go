// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: api/service.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_api_service_proto protoreflect.FileDescriptor

var file_api_service_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x12, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x65, 0x6c,
	0x69, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x39, 0x0a, 0x0d, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x50,
	0x69, 0x6e, 0x67, 0x12, 0x0d, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x1a, 0x0f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x9b, 0x02, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x31, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x12, 0x10, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x28, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x12, 0x0d, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0e, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x2b, 0x0a, 0x06, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x12, 0x0e, 0x2e, 0x52,
	0x65, 0x76, 0x6f, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x52,
	0x65, 0x76, 0x6f, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x2b, 0x0a, 0x06, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x12, 0x0e, 0x2e, 0x56, 0x65, 0x72, 0x69,
	0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x56, 0x65, 0x72, 0x69,
	0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x2e, 0x0a, 0x07,
	0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x12, 0x0f, 0x2e, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73,
	0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x52, 0x65, 0x66, 0x72, 0x65,
	0x73, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x25, 0x0a, 0x04,
	0x4c, 0x69, 0x73, 0x74, 0x12, 0x0c, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x08, 0x5a, 0x06, 0x70, 0x62, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_api_service_proto_goTypes = []interface{}{
	(*EmptyMessage)(nil),     // 0: EmptyMessage
	(*RegisterRequest)(nil),  // 1: RegisterRequest
	(*LoginRequest)(nil),     // 2: LoginRequest
	(*RevokeRequest)(nil),    // 3: RevokeRequest
	(*VerifyRequest)(nil),    // 4: VerifyRequest
	(*RefreshRequest)(nil),   // 5: RefreshRequest
	(*ListRequest)(nil),      // 6: ListRequest
	(*StatusResponse)(nil),   // 7: StatusResponse
	(*RegisterResponse)(nil), // 8: RegisterResponse
	(*LoginResponse)(nil),    // 9: LoginResponse
	(*RevokeResponse)(nil),   // 10: RevokeResponse
	(*VerifyResponse)(nil),   // 11: VerifyResponse
	(*RefreshResponse)(nil),  // 12: RefreshResponse
	(*ListResponse)(nil),     // 13: ListResponse
}
var file_api_service_proto_depIdxs = []int32{
	0,  // 0: pb.StatusService.Ping:input_type -> EmptyMessage
	1,  // 1: pb.AuthService.Register:input_type -> RegisterRequest
	2,  // 2: pb.AuthService.Login:input_type -> LoginRequest
	3,  // 3: pb.AuthService.Revoke:input_type -> RevokeRequest
	4,  // 4: pb.AuthService.Verify:input_type -> VerifyRequest
	5,  // 5: pb.AuthService.Refresh:input_type -> RefreshRequest
	6,  // 6: pb.AuthService.List:input_type -> ListRequest
	7,  // 7: pb.StatusService.Ping:output_type -> StatusResponse
	8,  // 8: pb.AuthService.Register:output_type -> RegisterResponse
	9,  // 9: pb.AuthService.Login:output_type -> LoginResponse
	10, // 10: pb.AuthService.Revoke:output_type -> RevokeResponse
	11, // 11: pb.AuthService.Verify:output_type -> VerifyResponse
	12, // 12: pb.AuthService.Refresh:output_type -> RefreshResponse
	13, // 13: pb.AuthService.List:output_type -> ListResponse
	7,  // [7:14] is the sub-list for method output_type
	0,  // [0:7] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_api_service_proto_init() }
func file_api_service_proto_init() {
	if File_api_service_proto != nil {
		return
	}
	file_api_delivery_proto_init()
	file_api_app_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_api_service_proto_goTypes,
		DependencyIndexes: file_api_service_proto_depIdxs,
	}.Build()
	File_api_service_proto = out.File
	file_api_service_proto_rawDesc = nil
	file_api_service_proto_goTypes = nil
	file_api_service_proto_depIdxs = nil
}
