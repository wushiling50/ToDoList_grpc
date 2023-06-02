DIR = $(shell pwd)/cmd
CONFIG_PATH = $(shell pwd)/config
PB_PATH = $(shell pwd)/service

.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install google.golang.org/grpc@latest
	go install go.etcd.io/etcd/client/v3@latest

.PHONY: proto
proto:
	@for file in $(PB_PATH)/*.proto; do \
		protoc --go_out=$(PB_PATH) --go-grpc_out=$(PB_PATH); \
	done

.PHONY: build-up
build-up:
	docker-compose up --build

.PHONY: down
down:
	docker-compose down