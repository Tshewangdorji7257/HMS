# Testing Guide - HMS (Hostel Management System)

Complete guide to run tests locally and verify CI/CD pipeline.

---

## 🚀 Quick Test Commands

### Run All Tests Locally (One Command)
```powershell
# From HMS root directory
cd Backend\auth-service; go test -v ./...; cd ..\..;
cd Backend\booking-service; go test -v ./...; cd ..\..;
cd Backend\building-service; go test -v ./...; cd ..\..;
cd Backend\api-gateway; go test -v ./...; cd ..\..;
cd Frontend; pnpm install; pnpm exec tsc --noEmit; pnpm lint; pnpm build
```

---

## 📋 Backend Tests

### Test All Backend Services
```powershell
# Test auth-service
cd Backend\auth-service
go test -v ./...
go test -v -race -coverprofile=coverage.out ./...
cd ..\..

# Test booking-service
cd Backend\booking-service
go test -v ./...
go test -v -race -coverprofile=coverage.out ./...
cd ..\..

# Test building-service
cd Backend\building-service
go test -v ./...
go test -v -race -coverprofile=coverage.out ./...
cd ..\..

# Test api-gateway
cd Backend\api-gateway
go test -v ./...
cd ..\..
```

### Check Backend Build
```powershell
# Build all services
cd Backend\auth-service; go build -v .; cd ..\..
cd Backend\booking-service; go build -v .; cd ..\..
cd Backend\building-service; go build -v .; cd ..\..
cd Backend\api-gateway; go build -v .; cd ..\..
```

### Run Backend Linting
```powershell
cd Backend\auth-service; go vet ./...; cd ..\..
cd Backend\booking-service; go vet ./...; cd ..\..
cd Backend\building-service; go vet ./...; cd ..\..
cd Backend\api-gateway; go vet ./...; cd ..\..
```

---

## 🎨 Frontend Tests

### Setup Frontend
```powershell
cd Frontend
pnpm install --no-frozen-lockfile
```

### Run TypeScript Check
```powershell
cd Frontend
pnpm exec tsc --noEmit
```

### Run Linting
```powershell
cd Frontend
pnpm lint
```

### Build Frontend
```powershell
cd Frontend
pnpm build
```

### Run Frontend Dev Server
```powershell
cd Frontend
pnpm dev
# Open http://localhost:3000
```

---

## 🐳 Docker Tests

### Build All Docker Images
```powershell
# Build auth-service
docker build -t hms-auth-service ./Backend/auth-service

# Build booking-service
docker build -t hms-booking-service ./Backend/booking-service

# Build building-service
docker build -t hms-building-service ./Backend/building-service

# Build api-gateway
docker build -t hms-api-gateway ./Backend/api-gateway

# Build frontend
docker build -t hms-frontend ./Frontend
```

### Run with Docker Compose
```powershell
cd Backend
docker-compose up -d
docker-compose ps
docker-compose logs
docker-compose down
```

---

## ✅ CI/CD Pipeline Tests

### Trigger GitHub Actions Manually
```powershell
# Commit and push to trigger CI
git add .
git commit -m "Test CI pipeline"
git push origin main
```

### Check CI Status
```powershell
# View workflow runs
# Go to: https://github.com/YOUR_USERNAME/HMS/actions
```

### CI Workflows Available
- **Backend CI**: Tests all 4 microservices
- **Frontend CI**: TypeScript, linting, and build
- **Full Stack CI**: Combined backend + frontend + integration tests
- **Docker Push**: Build and push images to GitHub Container Registry

---

## 🧪 Integration Tests

### Run Integration Tests Locally
```powershell
# Start backend services
cd Backend
docker-compose up -d

# Wait for services to start
Start-Sleep -Seconds 10

# Run frontend against backend
cd ..\Frontend
$env:NEXT_PUBLIC_API_URL="http://localhost:8080"
pnpm build
pnpm start

# Test endpoints
curl http://localhost:8080/health
curl http://localhost:3000

# Cleanup
cd ..\Backend
docker-compose down
```

---

## 📊 View Test Coverage

### Backend Coverage
```powershell
cd Backend\auth-service
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
cd ..\..
```

### Frontend Coverage (Future)
```powershell
cd Frontend
# pnpm test --coverage
# pnpm exec jest --coverage
```

---

## 🔍 Debug Failing Tests

### Check Backend Service Logs
```powershell
cd Backend
docker-compose logs auth-service
docker-compose logs booking-service
docker-compose logs building-service
docker-compose logs api-gateway
```

### Check Frontend Build Issues
```powershell
cd Frontend
pnpm exec tsc --noEmit --listFiles
pnpm build --debug
```

---

## 🎯 Quick Demo Test Script

Copy and run this complete test script:

```powershell
# HMS Complete Test Script
Write-Host "🚀 Starting HMS Test Suite..." -ForegroundColor Cyan

# Test Backend Services
Write-Host "`n📦 Testing Backend Services..." -ForegroundColor Yellow
$services = @("auth-service", "booking-service", "building-service", "api-gateway")
foreach ($service in $services) {
    Write-Host "Testing $service..." -ForegroundColor Green
    cd "Backend\$service"
    go test -v ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ $service tests failed!" -ForegroundColor Red
        cd ..\..
        exit 1
    }
    cd ..\..
}

# Test Frontend
Write-Host "`n🎨 Testing Frontend..." -ForegroundColor Yellow
cd Frontend
pnpm install --no-frozen-lockfile
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Frontend install failed!" -ForegroundColor Red
    exit 1
}

Write-Host "Running TypeScript check..." -ForegroundColor Green
pnpm exec tsc --noEmit
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ TypeScript check failed!" -ForegroundColor Red
    exit 1
}

Write-Host "Running ESLint..." -ForegroundColor Green
pnpm lint
if ($LASTEXITCODE -ne 0) {
    Write-Host "⚠️ Linting issues found (non-blocking)" -ForegroundColor Yellow
}

Write-Host "Building application..." -ForegroundColor Green
pnpm build
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Frontend build failed!" -ForegroundColor Red
    exit 1
}

cd ..

Write-Host "`n✅ All tests passed successfully!" -ForegroundColor Green
Write-Host "🎉 HMS is ready for deployment!" -ForegroundColor Cyan
```

---

## 📝 Pre-commit Checklist

Before pushing code, run:

```powershell
# 1. Format code
cd Backend\auth-service; go fmt ./...; cd ..\..

# 2. Run tests
cd Backend\auth-service; go test ./...; cd ..\..

# 3. Check TypeScript
cd Frontend; pnpm exec tsc --noEmit; cd ..

# 4. Build frontend
cd Frontend; pnpm build; cd ..

# 5. Commit and push
git add .
git commit -m "Your commit message"
git push origin main
```

---

## 🔗 Useful Links

- **GitHub Actions**: `https://github.com/YOUR_USERNAME/HMS/actions`
- **Docker Hub/GHCR**: `https://github.com/YOUR_USERNAME/HMS/pkgs/container/hms-frontend`
- **Local Frontend**: `http://localhost:3000`
- **Local API Gateway**: `http://localhost:8080`
- **Auth Service**: `http://localhost:8081`
- **Booking Service**: `http://localhost:8082`
- **Building Service**: `http://localhost:8083`

---

## 🆘 Common Issues

### Issue: "pnpm not found"
```powershell
npm install -g pnpm
```

### Issue: "go: command not found"
Download and install Go 1.24 from https://go.dev/dl/

### Issue: "docker: command not found"
Install Docker Desktop from https://www.docker.com/products/docker-desktop/

### Issue: TypeScript errors
```powershell
cd Frontend
rm -rf node_modules
pnpm install --force
pnpm exec tsc --noEmit
```

### Issue: Port already in use
```powershell
# Windows
netstat -ano | findstr :3000
taskkill /PID <PID> /F

# Or use different port
$env:PORT=3001
pnpm dev
```

---

## 🎓 CI/CD Pipeline Details

### What Gets Tested?

**Backend CI**:
- ✅ Go module verification
- ✅ `go vet` linting
- ✅ Unit tests with race detector
- ✅ Code coverage reports
- ✅ Service builds

**Frontend CI**:
- ✅ TypeScript type checking
- ✅ ESLint code quality
- ✅ Next.js production build
- ✅ Build artifact upload

**Integration Tests**:
- ✅ Verifies all services work together
- ✅ End-to-end workflow validation

**Docker Push**:
- ✅ Multi-platform builds (amd64, arm64)
- ✅ Push to GitHub Container Registry
- ✅ Security scanning with Snyk (optional)

---

**Last Updated**: November 27, 2025  
**Version**: 1.0  
**Maintainer**: HMS Development Team
