package proto

//go:generate protoc --go_out=../model/ --go_opt=paths=source_relative checklist.proto
//go:generate protoc --go_out=../service/ --go-grpc_out=../service/ --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative service.proto
