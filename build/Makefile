ROOT = ..

build:
	go build -o $(ROOT)/bin/ $(ROOT)/cmd/ova-checklist-api

generate:
	go generate -v -x $(ROOT)/...
	protoc --proto_path $(ROOT)/pkg/proto \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--go_out=$(ROOT)/internal/server/generated/service \
		--go-grpc_out=$(ROOT)/internal/server/generated/service \
		service.proto
