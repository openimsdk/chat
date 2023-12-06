protoc --go_out=plugins=grpc:./common --go_opt=module=github.com/OpenIMSDK/chat/pkg/proto/common common/common.proto
protoc --go_out=plugins=grpc:./admin --go_opt=module=github.com/OpenIMSDK/chat/pkg/proto/admin admin/admin.proto
protoc --go_out=plugins=grpc:./chat --go_opt=module=github.com/OpenIMSDK/chat/pkg/proto/chat chat/chat.proto