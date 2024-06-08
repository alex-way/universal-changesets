proto:
    protoc --go_out=./pkg/plugin --go_opt=paths=source_relative --go-grpc_out=./pkg/plugin --go-grpc_opt=paths=source_relative ./pkg/plugin/plugin.proto
run:
    go run main.go
