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

test-ut:
	go test $(ROOT)/...

test-integration:
	go test --tags=integration $(ROOT)/test/...

docker-compose-up:
	docker compose -f docker-compose.yml up --build

docker-compose-up-integration-test:
	docker compose -f docker-compose-integration-test.yml up --build

docker-compose-down:
	docker compose -f docker-compose.yml down --rmi all --volumes
