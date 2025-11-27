# Protocol Buffer Definitions

This directory contains the Protocol Buffer (protobuf) definitions for gRPC services.

## Prerequisites

Install the Protocol Buffer compiler and Go plugins:

```bash
# Install protoc compiler
# Download from: https://github.com/protocolbuffers/protobuf/releases

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Add to PATH (if not already)
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Windows Installation

1. Download protoc from https://github.com/protocolbuffers/protobuf/releases
2. Extract and add to PATH
3. Install Go plugins:
```powershell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Generating Code

### Linux/Mac:
```bash
chmod +x generate.sh
./generate.sh
```

### Windows:
```powershell
.\generate.ps1
```

## Proto Files

- `auth.proto` - Authentication service definitions
- `building.proto` - Building/Room/Bed service definitions

## Generated Files

Generated files will be placed in:
- `../auth-service/proto/auth/` - Auth service gRPC code
- `../building-service/proto/building/` - Building service gRPC code
