protoc --proto_path=api/proto --proto_path=third_party --go_out=plugins=grpc:pkg/api galasejahtera-service.proto
protoc --proto_path=api/proto --proto_path=third_party --grpc-gateway_out=logtostderr=true:pkg/api galasejahtera-service.proto
protoc --proto_path=api/proto --proto_path=third_party --swagger_out=logtostderr=true:api/swagger galasejahtera-service.proto
