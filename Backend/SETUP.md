# Quick Start Guide

This guide will help you get the Hostel Management System backend up and running quickly.

## ⚡ Quick Setup (5 minutes)

### Step 1: Start with Docker Compose

```powershell
# Navigate to backend directory
cd backend

# Start all services
docker-compose up -d

# Wait for services to start (about 30 seconds)
Start-Sleep -Seconds 30

# Verify services are running
docker-compose ps
```

Expected output:
```
NAME                    STATUS    PORTS
api-gateway             running   0.0.0.0:8000->8000/tcp
auth-service            running   0.0.0.0:8001->8001/tcp
building-service        running   0.0.0.0:8002->8002/tcp
booking-service         running   0.0.0.0:8003->8003/tcp
auth-db                 running   0.0.0.0:5433->5432/tcp
building-db             running   0.0.0.0:5434->5432/tcp
booking-db              running   0.0.0.0:5435->5432/tcp
```

### Step 2: Test the API

```powershell
# Health check
Invoke-RestMethod -Uri "http://localhost:8000/health" -Method GET

# Get all buildings (should return 10 buildings)
Invoke-RestMethod -Uri "http://localhost:8000/api/buildings" -Method GET
```

### Step 3: Create your first user

```powershell
$body = @{
    email = "admin@example.com"
    password = "admin123"
    name = "Admin User"
    role = "admin"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8000/api/auth/signup" `
    -Method POST `
    -ContentType "application/json" `
    -Body $body

# Save the token for future requests
$token = $response.token
Write-Host "Your token: $token"
```

### Step 4: Make a booking

```powershell
# First, get available buildings
$buildings = Invoke-RestMethod -Uri "http://localhost:8000/api/buildings" -Method GET
$firstBuilding = $buildings.buildings[0]
$firstRoom = $firstBuilding.rooms[0]
$firstBed = $firstRoom.beds[0]

# Create a booking
$bookingBody = @{
    user_id = $response.user.id
    user_name = $response.user.name
    building_id = $firstBuilding.id
    building_name = $firstBuilding.name
    room_id = $firstRoom.id
    room_number = $firstRoom.number
    bed_id = $firstBed.id
    bed_number = $firstBed.number
} | ConvertTo-Json

$booking = Invoke-RestMethod -Uri "http://localhost:8000/api/bookings" `
    -Method POST `
    -ContentType "application/json" `
    -Headers @{Authorization = "Bearer $token"} `
    -Body $bookingBody

Write-Host "Booking created: $($booking.booking.id)"
```

## 🔧 Manual Setup (Without Docker)

### Prerequisites Installation

#### 1. Install Go
```powershell
# Download from https://golang.org/dl/
# Or using Chocolatey
choco install golang
```

#### 2. Install PostgreSQL
```powershell
# Download from https://www.postgresql.org/download/windows/
# Or using Chocolatey
choco install postgresql
```

#### 3. Create Databases
```powershell
# Connect to PostgreSQL
psql -U postgres

# In psql shell:
CREATE DATABASE hostel_auth_db;
CREATE DATABASE hostel_building_db;
CREATE DATABASE hostel_booking_db;
\q
```

### Service Configuration

#### 1. Auth Service
```powershell
cd backend\auth-service

# Create .env file
@"
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=hostel_auth_db
JWT_SECRET=$(New-Guid)
JWT_EXPIRY=24h
PORT=8001
"@ | Out-File -FilePath .env -Encoding utf8

# Install dependencies
go mod download

# Run service
go run main.go
```

#### 2. Building Service
```powershell
# Open new terminal
cd backend\building-service

# Create .env file
@"
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=hostel_building_db
PORT=8002
AUTH_SERVICE_URL=http://localhost:8001
"@ | Out-File -FilePath .env -Encoding utf8

# Install dependencies
go mod download

# Run service
go run main.go
```

#### 3. Booking Service
```powershell
# Open new terminal
cd backend\booking-service

# Create .env file
@"
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=hostel_booking_db
PORT=8003
AUTH_SERVICE_URL=http://localhost:8001
BUILDING_SERVICE_URL=http://localhost:8002
"@ | Out-File -FilePath .env -Encoding utf8

# Install dependencies
go mod download

# Run service
go run main.go
```

#### 4. API Gateway
```powershell
# Open new terminal
cd backend\api-gateway

# Create .env file
@"
PORT=8000
AUTH_SERVICE_URL=http://localhost:8001
BUILDING_SERVICE_URL=http://localhost:8002
BOOKING_SERVICE_URL=http://localhost:8003
"@ | Out-File -FilePath .env -Encoding utf8

# Install dependencies
go mod download

# Run service
go run main.go
```

## 🧪 Testing the Setup

### Complete Test Script

Save this as `test-api.ps1`:

```powershell
# Test API endpoints
$baseUrl = "http://localhost:8000"

Write-Host "Testing Hostel Management API..." -ForegroundColor Green

# 1. Health Check
Write-Host "`n1. Testing health endpoint..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "$baseUrl/health" -Method GET
    Write-Host "✓ Health check passed" -ForegroundColor Green
} catch {
    Write-Host "✗ Health check failed: $_" -ForegroundColor Red
    exit 1
}

# 2. Sign Up
Write-Host "`n2. Creating new user..." -ForegroundColor Yellow
$signupBody = @{
    email = "test$(Get-Random)@example.com"
    password = "test123"
    name = "Test User"
    role = "student"
} | ConvertTo-Json

try {
    $signupResponse = Invoke-RestMethod -Uri "$baseUrl/api/auth/signup" `
        -Method POST `
        -ContentType "application/json" `
        -Body $signupBody
    
    $token = $signupResponse.token
    $userId = $signupResponse.user.id
    Write-Host "✓ User created: $($signupResponse.user.email)" -ForegroundColor Green
    Write-Host "  Token: $token" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Signup failed: $_" -ForegroundColor Red
    exit 1
}

# 3. Get Profile
Write-Host "`n3. Getting user profile..." -ForegroundColor Yellow
try {
    $profile = Invoke-RestMethod -Uri "$baseUrl/api/auth/profile" `
        -Method GET `
        -Headers @{Authorization = "Bearer $token"}
    
    Write-Host "✓ Profile retrieved: $($profile.user.name)" -ForegroundColor Green
} catch {
    Write-Host "✗ Get profile failed: $_" -ForegroundColor Red
}

# 4. Get All Buildings
Write-Host "`n4. Fetching all buildings..." -ForegroundColor Yellow
try {
    $buildings = Invoke-RestMethod -Uri "$baseUrl/api/buildings" -Method GET
    Write-Host "✓ Found $($buildings.buildings.Count) buildings" -ForegroundColor Green
    Write-Host "  Buildings: $($buildings.buildings.name -join ', ')" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Get buildings failed: $_" -ForegroundColor Red
    exit 1
}

# 5. Get Specific Building
Write-Host "`n5. Getting building details..." -ForegroundColor Yellow
$buildingId = $buildings.buildings[0].id
try {
    $building = Invoke-RestMethod -Uri "$baseUrl/api/buildings/$buildingId" -Method GET
    Write-Host "✓ Building: $($building.building.name)" -ForegroundColor Green
    Write-Host "  Total Rooms: $($building.building.total_rooms)" -ForegroundColor Cyan
    Write-Host "  Available Beds: $($building.building.available_beds)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Get building failed: $_" -ForegroundColor Red
}

# 6. Create Booking
Write-Host "`n6. Creating a booking..." -ForegroundColor Yellow
$firstRoom = $building.building.rooms[0]
$firstBed = $firstRoom.beds | Where-Object { -not $_.is_occupied } | Select-Object -First 1

if ($firstBed) {
    $bookingBody = @{
        user_id = $userId
        user_name = $signupResponse.user.name
        building_id = $building.building.id
        building_name = $building.building.name
        room_id = $firstRoom.id
        room_number = $firstRoom.number
        bed_id = $firstBed.id
        bed_number = $firstBed.number
    } | ConvertTo-Json

    try {
        $booking = Invoke-RestMethod -Uri "$baseUrl/api/bookings" `
            -Method POST `
            -ContentType "application/json" `
            -Headers @{Authorization = "Bearer $token"} `
            -Body $bookingBody
        
        $bookingId = $booking.booking.id
        Write-Host "✓ Booking created: $bookingId" -ForegroundColor Green
        Write-Host "  Room: $($booking.booking.room_number), Bed: $($booking.booking.bed_number)" -ForegroundColor Cyan
    } catch {
        Write-Host "✗ Create booking failed: $_" -ForegroundColor Red
    }
} else {
    Write-Host "⚠ No available beds found" -ForegroundColor Yellow
}

# 7. Get User Bookings
Write-Host "`n7. Fetching user bookings..." -ForegroundColor Yellow
try {
    $userBookings = Invoke-RestMethod -Uri "$baseUrl/api/bookings/users/$userId" `
        -Method GET `
        -Headers @{Authorization = "Bearer $token"}
    
    Write-Host "✓ Found $($userBookings.bookings.Count) booking(s)" -ForegroundColor Green
} catch {
    Write-Host "✗ Get bookings failed: $_" -ForegroundColor Red
}

# 8. Cancel Booking (if created)
if ($bookingId) {
    Write-Host "`n8. Cancelling booking..." -ForegroundColor Yellow
    try {
        $cancelled = Invoke-RestMethod -Uri "$baseUrl/api/bookings/$bookingId/cancel" `
            -Method PUT `
            -Headers @{Authorization = "Bearer $token"}
        
        Write-Host "✓ Booking cancelled successfully" -ForegroundColor Green
    } catch {
        Write-Host "✗ Cancel booking failed: $_" -ForegroundColor Red
    }
}

Write-Host "`n✓ All tests completed!" -ForegroundColor Green
```

Run the test script:
```powershell
.\test-api.ps1
```

## 🔍 Verification Checklist

- [ ] All 4 services are running (check with `docker-compose ps`)
- [ ] Can access http://localhost:8000/health
- [ ] Can create a new user via signup
- [ ] Can login with created user
- [ ] Can fetch all buildings (should return 10)
- [ ] Can create a booking
- [ ] Can view user's bookings
- [ ] Can cancel a booking

## 🐛 Common Issues

### Issue 1: Port Already in Use
```powershell
# Find process using port 8000
Get-NetTCPConnection -LocalPort 8000 -ErrorAction SilentlyContinue

# Kill the process
Get-Process -Id (Get-NetTCPConnection -LocalPort 8000).OwningProcess | Stop-Process -Force
```

### Issue 2: Database Connection Failed
```powershell
# Check if PostgreSQL is running
Get-Service -Name postgresql*

# Restart PostgreSQL
Restart-Service postgresql-x64-15
```

### Issue 3: Docker Containers Not Starting
```powershell
# View container logs
docker-compose logs auth-service
docker-compose logs building-service
docker-compose logs booking-service

# Rebuild containers
docker-compose down -v
docker-compose up --build
```

### Issue 4: JWT Token Invalid
- Ensure JWT_SECRET is the same across all services
- Check if token has expired (default: 24 hours)
- Verify Authorization header format: `Bearer <token>`

## 📊 Monitoring

### View Logs
```powershell
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f auth-service
docker-compose logs -f building-service
docker-compose logs -f booking-service
```

### Check Database Data
```powershell
# Connect to auth database
docker exec -it auth-db psql -U postgres -d hostel_auth_db

# View users
SELECT id, email, name, role FROM users;

# Connect to building database
docker exec -it building-db psql -U postgres -d hostel_building_db

# View buildings
SELECT id, name, total_rooms, available_beds FROM buildings;

# Connect to booking database
docker exec -it booking-db psql -U postgres -d hostel_booking_db

# View bookings
SELECT id, user_name, building_name, room_number, bed_number, status FROM bookings;
```

## 🚀 Next Steps

1. **Integrate with Frontend**:
   - Update frontend API base URL to `http://localhost:8000`
   - Test authentication flow
   - Test booking flow

2. **Customize Configuration**:
   - Change JWT expiry time
   - Add more buildings via seed data
   - Configure CORS for production domain

3. **Add Features**:
   - Email notifications
   - Payment integration
   - Admin dashboard
   - Booking history reports

## 📞 Support

If you encounter issues:
1. Check logs: `docker-compose logs`
2. Verify environment variables in `.env` files
3. Ensure all prerequisites are installed
4. Review API documentation in README.md

Happy coding! 🎉
