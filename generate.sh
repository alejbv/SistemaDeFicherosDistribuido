#!/bin/zsh
protoc tagFileSystempb/tagFileSystem.proto --go_out=. --go-grpc_out=.
protoc server/proto/chord.proto --go_out=. --go-grpc_out=.