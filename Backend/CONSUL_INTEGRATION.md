# Consul Service Discovery Integration

Consul has been successfully integrated into the HMS (Hostel Management System) for service discovery and health checking.

## 🎯 What is Consul?

Consul is a service mesh solution providing:
- **Service Discovery**: Automatically find services without hardcoded IPs/ports
- **Health Checking**: Monitor service health and remove unhealthy instances
- **Key/Value Store**: Configuration management
- **Multi-datacenter**: Support for distributed deployments

## 🏗️ Architecture

```
Consul Server (Port 8500)
    ↓
    ├── Auth Service (8001, gRPC: 9001)
    ├── Building Service (8002, gRPC: 9002)
    ├── Booking Service (8003)
    └── API Gateway (8000)
```

## ✅ What's Implemented

### 1. Consul Server
- **Port**: 8500 (HTTP API & Web UI)
- **DNS**: 8600/UDP
- **UI**: http://localhost:8500/ui
- **Mode**: Single-node development mode

### 2. Service Registration
Each service registers itself with:
- Service name (e.g., `auth-service`)
- Service ID (e.g., `auth-service-1`)
- HTTP port (8001-8003)
- gRPC port (9001-9002) in metadata
- Health check endpoint

### 3. Health Checks
- **Endpoint**: `/health`
- **Interval**: Every 10 seconds
- **Timeout**: 5 seconds
- **Auto-deregister**: After 30 seconds of critical state

### 4. Service Discovery Functions
```go
// Discover a service by name
address, err := consul.DiscoverService("auth-service")
// Returns: http://172.18.0.5:8001

// Get gRPC address
grpcAddr, err := consul.GetGRPCAddress("auth-service")
// Returns: 172.18.0.5:9001
```

## 🚀 Quick Start

### 1. Start All Services with Consul
```powershell
cd Backend
docker-compose up -d
```

### 2. Access Consul UI
```
http://localhost:8500/ui
```

### 3. View Registered Services
```powershell
curl http://localhost:8500/v1/agent/services
```

### 4. Check Service Health
```powershell
curl http://localhost:8500/v1/health/service/auth-service
```

## 📋 Environment Variables

Each service now has:

```yaml
environment:
  CONSUL_HOST: consul
  CONSUL_PORT: 8500
  SERVICE_NAME: auth-service
  SERVICE_ID: auth-service-1
  PORT: 8001
  GRPC_PORT: 9001
```

## 🔧 Integration in Go Services

### Step 1: Add Consul Package
```bash
go get github.com/hashicorp/consul/api
```

### Step 2: Initialize in main.go
```go
import "auth-service/consul"

func main() {
    // Initialize Consul client
    if err := consul.InitConsul(); err != nil {
        log.Fatalf("Failed to init consul: %v", err)
    }

    // Register service
    if err := consul.RegisterService(); err != nil {
        log.Fatalf("Failed to register service: %v", err)
    }
    defer consul.DeregisterService()

    // Your existing code...
}
```

### Step 3: Discover Services
```go
// Instead of hardcoded URLs:
authServiceURL := "http://auth-service:8001"  // ❌ Old way

// Use service discovery:
authServiceURL, err := consul.DiscoverService("auth-service")  // ✅ New way
if err != nil {
    log.Printf("Service discovery failed, using fallback: %v", err)
    authServiceURL = "http://auth-service:8001"
}
```

## 🎨 Consul UI Features

Access at `http://localhost:8500/ui`:

1. **Services**: View all registered services
2. **Nodes**: See all Consul agents
3. **Key/Value**: Store configuration
4. **Intentions**: Service-to-service authorization

## 📊 Health Check Status

| Status | Color | Meaning |
|--------|-------|---------|
| **Passing** | 🟢 Green | Service is healthy |
| **Warning** | 🟡 Yellow | Service degraded |
| **Critical** | 🔴 Red | Service unhealthy |

## 🔍 Useful Commands

### Check Consul Members
```powershell
docker exec hostel_consul consul members
```

### View Service Catalog
```powershell
docker exec hostel_consul consul catalog services
```

### Query DNS
```powershell
# Windows (nslookup)
nslookup auth-service.service.consul 127.0.0.1 -port=8600

# Linux/Mac (dig)
dig @127.0.0.1 -p 8600 auth-service.service.consul
```

### Deregister a Service
```powershell
curl -X PUT http://localhost:8500/v1/agent/service/deregister/auth-service-1
```

## 🆚 Before vs After

### Before (Hardcoded URLs):
```go
authServiceURL := "http://auth-service:8001"
buildingServiceURL := "http://building-service:8002"
```

**Problems:**
- ❌ Can't scale horizontally (no load balancing)
- ❌ Manual failover required
- ❌ No health awareness

### After (Consul Discovery):
```go
authServiceURL, _ := consul.DiscoverService("auth-service")
buildingServiceURL, _ := consul.DiscoverService("building-service")
```

**Benefits:**
- ✅ Auto-discovery of healthy instances
- ✅ Load balancing across multiple instances
- ✅ Automatic failover
- ✅ Health-aware routing

## 🔐 Production Considerations

For production deployment:

1. **Multi-node Cluster**: Run 3-5 Consul servers
2. **ACLs**: Enable authentication and authorization
3. **TLS**: Encrypt all communication
4. **Backup**: Regular snapshots of Consul data
5. **Monitoring**: Integrate with Prometheus/Grafana

### Multi-node Setup
```yaml
consul-server-1:
  command: agent -server -bootstrap-expect=3 -ui -client=0.0.0.0

consul-server-2:
  command: agent -server -retry-join=consul-server-1 -client=0.0.0.0

consul-server-3:
  command: agent -server -retry-join=consul-server-1 -client=0.0.0.0
```

## 🐛 Troubleshooting

### Service Not Appearing in Consul
```powershell
# Check if service can reach Consul
docker exec hostel_auth_service ping consul

# Check Consul logs
docker logs hostel_consul
```

### Health Check Failing
```powershell
# Test health endpoint directly
curl http://localhost:8001/health

# Check service logs
docker logs hostel_auth_service
```

### Discovery Returns Empty
```powershell
# Verify service is registered
curl http://localhost:8500/v1/agent/services | jq

# Check health status
curl http://localhost:8500/v1/health/service/auth-service?passing=true
```

## 📚 Additional Resources

- **Consul Docs**: https://www.consul.io/docs
- **API Reference**: https://www.consul.io/api-docs
- **Best Practices**: https://learn.hashicorp.com/consul

---

**Status**: ✅ Fully Integrated  
**Version**: Consul 1.18  
**Mode**: Development (Single-node)  
**Last Updated**: November 27, 2025
