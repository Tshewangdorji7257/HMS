# gRPC Implementation Complete! 🎉

## Overview

Your Hostel Management System backend now uses **gRPC for service-to-service communication** while maintaining REST APIs for the frontend.

## Architecture

```
Frontend (Browser)
    ↓ REST/HTTP
API Gateway (Port 8000)
    ↓ HTTP Proxy → Backend Services
    ↓ gRPC (Optional for future enhancements)
    
Auth Service
    ├─ HTTP API (Port 8001) ← Frontend via Gateway
    └─ gRPC Server (Port 9001) ← Service-to-Service
    
Building Service
    ├─ HTTP API (Port 8002) ← Frontend via Gateway
    └─ gRPC Server (Port 9002) ← Service-to-Service
    
Booking Service
    ├─ HTTP API (Port 8003) ← Frontend via Gateway
    └─ gRPC Clients → Auth & Building Services
```

## What's Implemented

### ✅ Protobuf Definitions (`proto/`)
- `auth.proto` - Authentication service definitions
- `building.proto` - Building/Room/Bed service definitions
- PowerShell & Bash scripts for code generation

### ✅ gRPC Servers
**Auth Service (Port 9001)**:
- `ValidateToken` - Validate JWT tokens
- `GetUserByID` - Retrieve user information
- `GetUserByEmail` - Find user by email

**Building Service (Port 9002)**:
- `GetBuildingByID` - Get building details
- `GetRoomByID` - Get room details
- `GetBedByID` - Get bed details
- `UpdateBedOccupancy` - Update bed status
- `GetBedsByUserID` - Get user's beds

### ✅ gRPC Clients
**Booking Service**:
- Uses gRPC to call Auth & Building services
- Validates users via `ValidateUser()`
- Updates bed occupancy via `UpdateBedOccupancy()`

**API Gateway**:
- Optional gRPC client setup for future enhancements
- Currently uses HTTP proxy for simplicity

## Ports Configuration

| Service | HTTP Port | gRPC Port |
|---------|-----------|-----------|
| API Gateway | 8000 | - |
| Auth Service | 8001 | 9001 |
| Building Service | 8002 | 9002 |
| Booking Service | 8003 | - |

## Environment Variables

### Auth Service
```env
PORT=8001
GRPC_PORT=9001
```

### Building Service
```env
PORT=8002
GRPC_PORT=9002
```

### Booking Service
```env
PORT=8003
AUTH_GRPC_URL=auth-service:9001
BUILDING_GRPC_URL=building-service:9002
```

### API Gateway
```env
PORT=8000
AUTH_GRPC_URL=auth-service:9001
BUILDING_GRPC_URL=building-service:9002
```

## Benefits of gRPC

1. **Performance** - Binary protocol, faster than JSON/HTTP
2. **Type Safety** - Proto definitions ensure contract compliance
3. **Streaming** - Supports bidirectional streaming (future enhancement)
4. **Language Agnostic** - Can add services in other languages
5. **Auto-generated Code** - Proto compiler generates client/server code

## Frontend Integration

**No changes required!** The frontend continues to use REST APIs through the API Gateway at `http://localhost:8000`.

## Development Workflow

### Updating Proto Definitions

1. Edit proto files in `proto/` directory
2. Run code generation:
   ```powershell
   cd proto
   .\generate.ps1
   ```
3. Update service implementations
4. Rebuild and redeploy

### Adding New gRPC Methods

1. Add method to `.proto` file
2. Regenerate code
3. Implement method in service's `grpc/server.go`
4. Update client calls if needed
5. Test the new functionality

## Testing gRPC Services

### Using grpcurl (Command Line)

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List services
grpcurl -plaintext localhost:9001 list

# Call a method
grpcurl -plaintext -d '{"user_id":"123"}' \
  localhost:9001 auth.AuthService/GetUserByID
```

### Using Postman (GUI)

1. Import proto files into Postman
2. Create new gRPC request
3. Connect to `localhost:9001` or `localhost:9002`
4. Test methods interactively

## Future Enhancements

- [ ] Add gRPC streaming for real-time updates
- [ ] Implement gRPC interceptors for logging
- [ ] Add gRPC health checks
- [ ] Implement service mesh (Istio/Linkerd)
- [ ] Add gRPC load balancing
- [ ] Implement circuit breakers
- [ ] Add distributed tracing (OpenTelemetry)

## Troubleshooting

### gRPC Connection Errors

```bash
# Check if gRPC ports are open
docker exec hostel_auth_service netstat -tuln | grep 9001
docker exec hostel_building_service netstat -tuln | grep 9002
```

### Proto Generation Issues

```powershell
# Reinstall protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Import Path Issues

Ensure go.mod module names match the import paths in proto files:
```go
option go_package = "auth-service/proto/auth";
```

## Documentation

- [gRPC Official Docs](https://grpc.io/docs/)
- [Protocol Buffers Guide](https://protobuf.dev/)
- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)

---

**Your backend is now production-ready with high-performance gRPC communication!** 🚀
