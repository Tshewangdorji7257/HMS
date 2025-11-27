# API Testing Guide

Complete guide for testing all API endpoints using cURL and PowerShell.

## 🔐 Authentication Flow

### 1. Sign Up (Create New User)

**cURL (Windows PowerShell)**:
```powershell
curl -X POST http://localhost:8000/api/auth/signup `
  -H "Content-Type: application/json" `
  -d '{\"email\":\"student@example.com\",\"password\":\"password123\",\"name\":\"John Doe\",\"role\":\"student\"}'
```

**PowerShell**:
```powershell
$body = @{
    email = "student@example.com"
    password = "password123"
    name = "John Doe"
    role = "student"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8000/api/auth/signup" `
    -Method POST `
    -ContentType "application/json" `
    -Body $body
```

**Expected Response**:
```json
{
  "success": true,
  "message": "User created successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "student@example.com",
    "name": "John Doe",
    "role": "student",
    "created_at": "2025-11-25T10:30:00Z"
  }
}
```

### 2. Login

**cURL**:
```powershell
curl -X POST http://localhost:8000/api/auth/login `
  -H "Content-Type: application/json" `
  -d '{\"email\":\"student@example.com\",\"password\":\"password123\"}'
```

**PowerShell**:
```powershell
$loginBody = @{
    email = "student@example.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8000/api/auth/login" `
    -Method POST `
    -ContentType "application/json" `
    -Body $loginBody

# Save token for later use
$token = $response.token
```

### 3. Validate Token

**cURL**:
```powershell
curl -X POST http://localhost:8000/api/auth/validate `
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**PowerShell**:
```powershell
Invoke-RestMethod -Uri "http://localhost:8000/api/auth/validate" `
    -Method POST `
    -Headers @{Authorization = "Bearer $token"}
```

### 4. Get User Profile

**cURL**:
```powershell
curl http://localhost:8000/api/auth/profile `
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**PowerShell**:
```powershell
Invoke-RestMethod -Uri "http://localhost:8000/api/auth/profile" `
    -Method GET `
    -Headers @{Authorization = "Bearer $token"}
```

## 🏢 Building & Room Management

### 1. Get All Buildings

**cURL**:
```powershell
curl http://localhost:8000/api/buildings
```

**PowerShell**:
```powershell
$buildings = Invoke-RestMethod -Uri "http://localhost:8000/api/buildings" -Method GET
$buildings.buildings | Format-Table name, total_rooms, available_beds
```

**Expected Response**:
```json
{
  "success": true,
  "buildings": [
    {
      "id": "bldg-1",
      "name": "RK A",
      "description": "Modern residence with excellent facilities",
      "total_rooms": 25,
      "total_beds": 50,
      "available_beds": 30,
      "amenities": ["Wi-Fi", "Gym", "Laundry", "Common Room"],
      "image": "/images/rk-a.jpg",
      "rooms": [...]
    }
  ]
}
```

### 2. Get Specific Building

**cURL**:
```powershell
curl http://localhost:8000/api/buildings/bldg-1
```

**PowerShell**:
```powershell
$buildingId = "bldg-1"
$building = Invoke-RestMethod -Uri "http://localhost:8000/api/buildings/$buildingId"
Write-Host "Building: $($building.building.name)"
Write-Host "Total Rooms: $($building.building.total_rooms)"
Write-Host "Available Beds: $($building.building.available_beds)"
```

### 3. Get Room Details

**cURL**:
```powershell
curl http://localhost:8000/api/buildings/bldg-1/rooms/bldg-1-room-001
```

**PowerShell**:
```powershell
$roomId = "bldg-1-room-001"
Invoke-RestMethod -Uri "http://localhost:8000/api/buildings/bldg-1/rooms/$roomId"
```

### 4. Search Buildings

**cURL**:
```powershell
curl http://localhost:8000/api/buildings/search?q=RK
```

**PowerShell**:
```powershell
$searchQuery = "RK"
Invoke-RestMethod -Uri "http://localhost:8000/api/buildings/search?q=$searchQuery"
```

## 📅 Booking Management

### 1. Create Booking

**cURL**:
```powershell
curl -X POST http://localhost:8000/api/bookings `
  -H "Content-Type: application/json" `
  -H "Authorization: Bearer YOUR_TOKEN" `
  -d '{\"user_id\":\"user-uuid\",\"user_name\":\"John Doe\",\"building_id\":\"bldg-1\",\"building_name\":\"RK A\",\"room_id\":\"bldg-1-room-001\",\"room_number\":\"001\",\"bed_id\":\"bldg-1-room-001-bed-1\",\"bed_number\":1}'
```

**PowerShell**:
```powershell
# First get an available bed
$buildings = Invoke-RestMethod -Uri "http://localhost:8000/api/buildings"
$firstBuilding = $buildings.buildings[0]
$firstRoom = $firstBuilding.rooms[0]
$availableBed = $firstRoom.beds | Where-Object { -not $_.is_occupied } | Select-Object -First 1

# Create booking
$bookingBody = @{
    user_id = "your-user-id"
    user_name = "John Doe"
    building_id = $firstBuilding.id
    building_name = $firstBuilding.name
    room_id = $firstRoom.id
    room_number = $firstRoom.number
    bed_id = $availableBed.id
    bed_number = $availableBed.number
} | ConvertTo-Json

$booking = Invoke-RestMethod -Uri "http://localhost:8000/api/bookings" `
    -Method POST `
    -ContentType "application/json" `
    -Headers @{Authorization = "Bearer $token"} `
    -Body $bookingBody
```

**Expected Response**:
```json
{
  "success": true,
  "message": "Booking created successfully",
  "booking": {
    "id": "booking-uuid",
    "user_id": "user-uuid",
    "user_name": "John Doe",
    "building_id": "bldg-1",
    "building_name": "RK A",
    "room_id": "bldg-1-room-001",
    "room_number": "001",
    "bed_id": "bldg-1-room-001-bed-1",
    "bed_number": 1,
    "booking_date": "2025-11-25T10:30:00Z",
    "status": "active",
    "created_at": "2025-11-25T10:30:00Z"
  }
}
```

### 2. Get All Bookings (Admin Only)

**cURL**:
```powershell
curl http://localhost:8000/api/bookings `
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

**PowerShell**:
```powershell
$allBookings = Invoke-RestMethod -Uri "http://localhost:8000/api/bookings" `
    -Headers @{Authorization = "Bearer $adminToken"}
$allBookings.bookings | Format-Table user_name, building_name, room_number, bed_number, status
```

### 3. Get User's Bookings

**cURL**:
```powershell
curl http://localhost:8000/api/bookings/users/user-uuid `
  -H "Authorization: Bearer YOUR_TOKEN"
```

**PowerShell**:
```powershell
$userId = "your-user-id"
$userBookings = Invoke-RestMethod -Uri "http://localhost:8000/api/bookings/users/$userId" `
    -Headers @{Authorization = "Bearer $token"}
```

### 4. Get Specific Booking

**cURL**:
```powershell
curl http://localhost:8000/api/bookings/booking-uuid `
  -H "Authorization: Bearer YOUR_TOKEN"
```

**PowerShell**:
```powershell
$bookingId = "booking-uuid"
Invoke-RestMethod -Uri "http://localhost:8000/api/bookings/$bookingId" `
    -Headers @{Authorization = "Bearer $token"}
```

### 5. Cancel Booking

**cURL**:
```powershell
curl -X PUT http://localhost:8000/api/bookings/booking-uuid/cancel `
  -H "Authorization: Bearer YOUR_TOKEN"
```

**PowerShell**:
```powershell
$bookingId = "booking-uuid"
Invoke-RestMethod -Uri "http://localhost:8000/api/bookings/$bookingId/cancel" `
    -Method PUT `
    -Headers @{Authorization = "Bearer $token"}
```

## 🧪 Complete Test Scenarios

### Scenario 1: Student Booking Flow

```powershell
# 1. Sign up as student
$signupBody = @{
    email = "student1@example.com"
    password = "student123"
    name = "Alice Student"
    role = "student"
} | ConvertTo-Json

$user = Invoke-RestMethod -Uri "http://localhost:8000/api/auth/signup" `
    -Method POST -ContentType "application/json" -Body $signupBody
$token = $user.token

# 2. Browse available buildings
$buildings = Invoke-RestMethod -Uri "http://localhost:8000/api/buildings"
Write-Host "Available buildings:"
$buildings.buildings | ForEach-Object {
    Write-Host "- $($_.name): $($_.available_beds) beds available"
}

# 3. Select a building and room
$selectedBuilding = $buildings.buildings[0]
$selectedRoom = $selectedBuilding.rooms | Where-Object { $_.available_beds -gt 0 } | Select-Object -First 1
$selectedBed = $selectedRoom.beds | Where-Object { -not $_.is_occupied } | Select-Object -First 1

Write-Host "`nSelected:"
Write-Host "Building: $($selectedBuilding.name)"
Write-Host "Room: $($selectedRoom.number) ($($selectedRoom.type))"
Write-Host "Bed: $($selectedBed.number)"

# 4. Create booking
$bookingBody = @{
    user_id = $user.user.id
    user_name = $user.user.name
    building_id = $selectedBuilding.id
    building_name = $selectedBuilding.name
    room_id = $selectedRoom.id
    room_number = $selectedRoom.number
    bed_id = $selectedBed.id
    bed_number = $selectedBed.number
} | ConvertTo-Json

$booking = Invoke-RestMethod -Uri "http://localhost:8000/api/bookings" `
    -Method POST -ContentType "application/json" `
    -Headers @{Authorization = "Bearer $token"} -Body $bookingBody

Write-Host "`n✓ Booking created: $($booking.booking.id)"

# 5. View my bookings
$myBookings = Invoke-RestMethod -Uri "http://localhost:8000/api/bookings/users/$($user.user.id)" `
    -Headers @{Authorization = "Bearer $token"}
Write-Host "`nMy bookings:"
$myBookings.bookings | Format-Table building_name, room_number, bed_number, status -AutoSize

# 6. Cancel booking
Invoke-RestMethod -Uri "http://localhost:8000/api/bookings/$($booking.booking.id)/cancel" `
    -Method PUT -Headers @{Authorization = "Bearer $token"}
Write-Host "`n✓ Booking cancelled"
```

### Scenario 2: Admin View All Bookings

```powershell
# 1. Sign up as admin
$adminBody = @{
    email = "admin@example.com"
    password = "admin123"
    name = "Admin User"
    role = "admin"
} | ConvertTo-Json

$admin = Invoke-RestMethod -Uri "http://localhost:8000/api/auth/signup" `
    -Method POST -ContentType "application/json" -Body $adminBody
$adminToken = $admin.token

# 2. View all bookings
$allBookings = Invoke-RestMethod -Uri "http://localhost:8000/api/bookings" `
    -Headers @{Authorization = "Bearer $adminToken"}

Write-Host "All Bookings:"
$allBookings.bookings | Format-Table user_name, building_name, room_number, status, booking_date -AutoSize

# 3. View buildings with occupancy
$buildings = Invoke-RestMethod -Uri "http://localhost:8000/api/buildings"
Write-Host "`nBuilding Occupancy:"
$buildings.buildings | ForEach-Object {
    $occupancyRate = (($_.total_beds - $_.available_beds) / $_.total_beds) * 100
    Write-Host "$($_.name): $($_.total_beds - $_.available_beds)/$($_.total_beds) beds occupied ($([math]::Round($occupancyRate, 1))%)"
}
```

### Scenario 3: Search and Filter

```powershell
# 1. Search for specific buildings
$searchResults = Invoke-RestMethod -Uri "http://localhost:8000/api/buildings/search?q=RK"
Write-Host "Buildings matching 'RK':"
$searchResults.buildings | Format-Table name, total_rooms, available_beds

# 2. Filter rooms by type
$buildings = Invoke-RestMethod -Uri "http://localhost:8000/api/buildings"
Write-Host "`nSingle Rooms:"
foreach ($building in $buildings.buildings) {
    $singleRooms = $building.rooms | Where-Object { $_.type -eq "single" }
    if ($singleRooms) {
        Write-Host "$($building.name):"
        $singleRooms | Format-Table number, price, available_beds -AutoSize
    }
}

# 3. Find cheapest rooms
$allRooms = @()
foreach ($building in $buildings.buildings) {
    foreach ($room in $building.rooms) {
        $allRooms += [PSCustomObject]@{
            Building = $building.name
            Room = $room.number
            Type = $room.type
            Price = $room.price
            Available = $room.available_beds
        }
    }
}
Write-Host "`nCheapest Rooms:"
$allRooms | Where-Object { $_.Available -gt 0 } | Sort-Object Price | Select-Object -First 10 | Format-Table
```

## 📊 Data Validation Tests

### Test 1: Duplicate Booking Prevention
```powershell
# Try to book twice with same user
$bookingBody = @{
    user_id = "user-id"
    user_name = "Test User"
    building_id = "bldg-1"
    building_name = "RK A"
    room_id = "bldg-1-room-001"
    room_number = "001"
    bed_id = "bldg-1-room-001-bed-1"
    bed_number = 1
} | ConvertTo-Json

# First booking should succeed
$booking1 = Invoke-RestMethod -Uri "http://localhost:8000/api/bookings" `
    -Method POST -ContentType "application/json" `
    -Headers @{Authorization = "Bearer $token"} -Body $bookingBody

# Second booking should fail
try {
    $booking2 = Invoke-RestMethod -Uri "http://localhost:8000/api/bookings" `
        -Method POST -ContentType "application/json" `
        -Headers @{Authorization = "Bearer $token"} -Body $bookingBody
    Write-Host "✗ Test Failed: Should not allow duplicate booking" -ForegroundColor Red
} catch {
    Write-Host "✓ Test Passed: Duplicate booking prevented" -ForegroundColor Green
}
```

### Test 2: Authentication Required
```powershell
# Try to create booking without token
try {
    Invoke-RestMethod -Uri "http://localhost:8000/api/bookings" `
        -Method POST -ContentType "application/json" `
        -Body $bookingBody
    Write-Host "✗ Test Failed: Should require authentication" -ForegroundColor Red
} catch {
    Write-Host "✓ Test Passed: Authentication required" -ForegroundColor Green
}
```

### Test 3: Invalid Email Format
```powershell
$invalidBody = @{
    email = "invalid-email"
    password = "test123"
    name = "Test User"
} | ConvertTo-Json

try {
    Invoke-RestMethod -Uri "http://localhost:8000/api/auth/signup" `
        -Method POST -ContentType "application/json" -Body $invalidBody
    Write-Host "✗ Test Failed: Should validate email format" -ForegroundColor Red
} catch {
    Write-Host "✓ Test Passed: Email validation works" -ForegroundColor Green
}
```

## 🎯 Performance Testing

### Load Test: Multiple Concurrent Signups
```powershell
$jobs = @()
for ($i = 1; $i -le 10; $i++) {
    $jobs += Start-Job -ScriptBlock {
        param($index)
        $body = @{
            email = "user$index@example.com"
            password = "test123"
            name = "User $index"
            role = "student"
        } | ConvertTo-Json
        
        Invoke-RestMethod -Uri "http://localhost:8000/api/auth/signup" `
            -Method POST -ContentType "application/json" -Body $body
    } -ArgumentList $i
}

$results = $jobs | Wait-Job | Receive-Job
Write-Host "Created $($results.Count) users concurrently"
$jobs | Remove-Job
```

## 📝 Export Test Data

### Save All Buildings to JSON
```powershell
$buildings = Invoke-RestMethod -Uri "http://localhost:8000/api/buildings"
$buildings | ConvertTo-Json -Depth 10 | Out-File "buildings-export.json"
Write-Host "Buildings exported to buildings-export.json"
```

### Generate Report
```powershell
$buildings = Invoke-RestMethod -Uri "http://localhost:8000/api/buildings"

$report = @"
# Hostel Management Report
Generated: $(Get-Date)

## Building Summary
Total Buildings: $($buildings.buildings.Count)
Total Rooms: $(($buildings.buildings | Measure-Object -Property total_rooms -Sum).Sum)
Total Beds: $(($buildings.buildings | Measure-Object -Property total_beds -Sum).Sum)
Available Beds: $(($buildings.buildings | Measure-Object -Property available_beds -Sum).Sum)

## Buildings:
"@

foreach ($building in $buildings.buildings) {
    $report += "`n### $($building.name)"
    $report += "`n- Rooms: $($building.total_rooms)"
    $report += "`n- Beds: $($building.total_beds)"
    $report += "`n- Available: $($building.available_beds)"
    $report += "`n- Amenities: $($building.amenities -join ', ')"
}

$report | Out-File "hostel-report.md"
Write-Host "Report generated: hostel-report.md"
```

---

**Happy Testing! 🧪**
