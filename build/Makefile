ROOT = ..

build:
	go build -o $(ROOT)/bin/ $(ROOT)/cmd/ova-checklist-api

generate:
	go generate -v -x \
		$(ROOT)/pkg/proto \
		$(ROOT)/internal/generated/flusher \
		$(ROOT)/internal/generated/repo
