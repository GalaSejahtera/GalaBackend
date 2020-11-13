protoc --proto_path=api/proto --proto_path=third_party --go_out=plugins=grpc:pkg/api safeworkout-service.proto
protoc --proto_path=api/proto --proto_path=third_party --grpc-gateway_out=logtostderr=true:pkg/api safeworkout-service.proto
protoc --proto_path=api/proto --proto_path=third_party --swagger_out=logtostderr=true:api/swagger safeworkout-service.proto
