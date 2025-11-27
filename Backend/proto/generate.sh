#!/bin/bash

# Generate Go code from proto files

# Auth Service
protoc --go_out=../auth-service --go_opt=paths=source_relative \
    --go-grpc_out=../auth-service --go-grpc_opt=paths=source_relative \
    auth.proto

# Building Service
protoc --go_out=../building-service --go_opt=paths=source_relative \
    --go-grpc_out=../building-service --go-grpc_opt=paths=source_relative \
    building.proto

echo "Proto files generated successfully!"
