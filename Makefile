# get makefile path
MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))

# get makefile dir
MAKEFILE_DIR := $(dir $(MAKEFILE_PATH))

# PROTO STUFF 
GOPATH_BIN = $(shell go env GOPATH)/bin
PROTOC_GEN_GO = $(GOPATH_BIN)/protoc-gen-go
PROTOC = $(MAKEFILE_DIR)protoc/bin/protoc
GO_OUT = $(MAKEFILE_DIR)src/generated
PROTOS_DIR = $(MAKEFILE_DIR)proto
PROTOS = $(wildcard $(MAKEFILE_DIR)proto/*.proto)

generate-proto:
#	$(PROTOC) --go_out=$(GO_OUT) -I $(PROTOS)
	$(PROTOC) --proto_path=${PROTOS_DIR} --plugin=protoc-gen-go=$(PROTOC_GEN_GO) --go_out=$(GO_OUT) --go_opt=paths=source_relative $(PROTOS)

generate:
	go generate ./...


