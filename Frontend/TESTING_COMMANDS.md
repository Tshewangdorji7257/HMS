# Testing Commands

## Quick Tests

### Test Backend Health
```powershell
pnpm test:backend
```

### Test Buildings API
```powershell
pnpm test:buildings
```

## Manual API Testing

### 1. Health Check
```powershell
Invoke-RestMethod -Uri http://localhost:8000/health
```

### 2. Get All Buildings
```powershell
$buildings = Invoke-RestMethod -Uri http://localhost:8000/api/buildings
$buildings.buildings | Select-Object name, available_beds | Format-Table
```

### 3. Create Test User
```powershell
$body = @{
    email = "test@example.com"
    password = "password123"
    name = "Test User"
    role = "student"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri http://localhost:8000/api/auth/signup `
    -Method POST `
    -ContentType "application/json" `
    -Body $body

Write-Host "User ID: $($response.user.id)"
Write-Host "Token: $($response.token)"
```

### 4. Login User
```powershell
$body = @{
    email = "test@example.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri http://localhost:8000/api/auth/login `
    -Method POST `
    -ContentType "application/json" `
    -Body $body

$token = $response.token
Write-Host "Token: $token"
```

### 5. Get User Profile
```powershell
# Use token from previous step
$profile = Invoke-RestMethod -Uri http://localhost:8000/api/auth/profile `
    -Headers @{Authorization = "Bearer $token"}

Write-Host "Name: $($profile.user.name)"
Write-Host "Email: $($profile.user.email)"
Write-Host "Role: $($profile.user.role)"
```

### 6. Search Buildings
```powershell
$results = Invoke-RestMethod -Uri "http://localhost:8000/api/buildings/search?q=RK"
$results.buildings | Select-Object name, available_beds
```

### 7. Get Specific Building
```powershell
# Get first building ID
$buildings = Invoke-RestMethod -Uri http://localhost:8000/api/buildings
$buildingId = $buildings.buildings[0].id

# Get building details
$building = Invoke-RestMethod -Uri "http://localhost:8000/api/buildings/$buildingId"
Write-Host "Building: $($building.building.name)"
Write-Host "Rooms: $($building.building.rooms.Count)"
```

## Start Application

### Development Mode
```powershell
pnpm dev
```

### Build for Production
```powershell
pnpm build
pnpm start
```

## Check Backend Services

### View Running Containers
```powershell
docker ps
```

### View Backend Logs
```powershell
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api-gateway
docker-compose logs -f auth-service
docker-compose logs -f building-service
docker-compose logs -f booking-service
```

### Restart Backend Services
```powershell
docker-compose restart
```

### Stop Backend Services
```powershell
docker-compose down
```

### Start Backend Services
```powershell
docker-compose up -d
```

## Troubleshooting

### Clear Frontend Cache
```powershell
# In browser console (F12):
localStorage.clear()
location.reload()
```

### Check Port Usage
```powershell
# Check if port 3000 is in use
netstat -ano | findstr :3000

# Check if port 8000 is in use
netstat -ano | findstr :8000
```

### Reset Everything
```powershell
# Stop frontend (Ctrl+C in terminal)

# Stop backend
docker-compose down

# Start backend
docker-compose up -d

# Wait for services
Start-Sleep -Seconds 30

# Test backend
Invoke-RestMethod -Uri http://localhost:8000/health

# Start frontend
pnpm dev
```
