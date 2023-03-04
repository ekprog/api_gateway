protoc -I ./proto/test_service \
--go_out ./pkg/test_service \
--go_opt paths=source_relative \
--go-grpc_out ./pkg/test_service \
--go-grpc_opt paths=source_relative \
--grpc-gateway_out ./pkg/test_service \
--grpc-gateway_opt paths=source_relative \
./proto/test_service/api/**/*.proto
