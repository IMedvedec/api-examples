protoc -I ./api/grpc/proto \
    --go_out ./api/grpc/build --go_opt paths=source_relative \
    --go-grpc_out ./api/grpc/build --go-grpc_opt paths=source_relative \
    greeting.proto

protoc -I ./api/grpc/proto \
	--grpc-gateway_out ./api/grpc/build \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    greeting.proto

