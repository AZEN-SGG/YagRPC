generate-proto:
	protoc --go_out=./pkg/api --go_opt=paths=source_relative --go-grpc_out=./pkg/api --go-grpc_opt=paths=source_relative 3_1/order.proto