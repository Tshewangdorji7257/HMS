# Backend Connectivity Test Script
# Run this script to verify the backend is running and accessible

Write-Host "🔍 Testing Backend Connectivity..." -ForegroundColor Cyan
Write-Host ""

$baseUrl = "http://localhost:8000"
$allTestsPassed = $true

# Test 1: API Gateway Health Check
Write-Host "1. Testing API Gateway Health..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "$baseUrl/health" -Method GET -ErrorAction Stop
    Write-Host "   ✓ API Gateway is healthy" -ForegroundColor Green
    Write-Host "   Status: $($health.status)" -ForegroundColor Gray
} catch {
    Write-Host "   ✗ API Gateway health check failed" -ForegroundColor Red
    Write-Host "   Error: $_" -ForegroundColor Red
    $allTestsPassed = $false
}
Write-Host ""

# Test 2: Get Buildings
Write-Host "2. Testing Buildings API..." -ForegroundColor Yellow
try {
    $buildings = Invoke-RestMethod -Uri "$baseUrl/api/buildings" -Method GET -ErrorAction Stop
    if ($buildings.success -and $buildings.buildings) {
        Write-Host "   ✓ Buildings API working" -ForegroundColor Green
        Write-Host "   Found $($buildings.buildings.Count) buildings" -ForegroundColor Gray
        
        # Display first few buildings
        if ($buildings.buildings.Count -gt 0) {
            Write-Host "   Sample buildings:" -ForegroundColor Gray
            $buildings.buildings | Select-Object -First 3 | ForEach-Object {
                Write-Host "     - $($_.name): $($_.available_beds) beds available" -ForegroundColor Gray
            }
        }
    } else {
        Write-Host "   ⚠ Buildings API returned unexpected format" -ForegroundColor Yellow
        $allTestsPassed = $false
    }
} catch {
    Write-Host "   ✗ Buildings API failed" -ForegroundColor Red
    Write-Host "   Error: $_" -ForegroundColor Red
    $allTestsPassed = $false
}
Write-Host ""

# Test 3: Test User Signup
Write-Host "3. Testing User Signup..." -ForegroundColor Yellow
$testEmail = "test-$(Get-Random)@example.com"
$signupBody = @{
    email = $testEmail
    password = "test123456"
    name = "Test User"
    role = "student"
} | ConvertTo-Json

try {
    $signup = Invoke-RestMethod -Uri "$baseUrl/api/auth/signup" `
        -Method POST `
        -ContentType "application/json" `
        -Body $signupBody `
        -ErrorAction Stop
    
    if ($signup.success -and $signup.token) {
        Write-Host "   ✓ User signup successful" -ForegroundColor Green
        Write-Host "   User ID: $($signup.user.id)" -ForegroundColor Gray
        Write-Host "   Email: $($signup.user.email)" -ForegroundColor Gray
        $testToken = $signup.token
        $testUserId = $signup.user.id
    } else {
        Write-Host "   ⚠ Signup returned unexpected format" -ForegroundColor Yellow
        $allTestsPassed = $false
    }
} catch {
    Write-Host "   ✗ User signup failed" -ForegroundColor Red
    Write-Host "   Error: $_" -ForegroundColor Red
    $allTestsPassed = $false
}
Write-Host ""

# Test 4: Test User Login
Write-Host "4. Testing User Login..." -ForegroundColor Yellow
$loginBody = @{
    email = $testEmail
    password = "test123456"
} | ConvertTo-Json

try {
    $login = Invoke-RestMethod -Uri "$baseUrl/api/auth/login" `
        -Method POST `
        -ContentType "application/json" `
        -Body $loginBody `
        -ErrorAction Stop
    
    if ($login.success -and $login.token) {
        Write-Host "   ✓ User login successful" -ForegroundColor Green
        Write-Host "   Token received: $(if($login.token.Length -gt 20) { $login.token.Substring(0,20) + '...' } else { $login.token })" -ForegroundColor Gray
    } else {
        Write-Host "   ⚠ Login returned unexpected format" -ForegroundColor Yellow
        $allTestsPassed = $false
    }
} catch {
    Write-Host "   ✗ User login failed" -ForegroundColor Red
    Write-Host "   Error: $_" -ForegroundColor Red
    $allTestsPassed = $false
}
Write-Host ""

# Test 5: Test Token Validation
if ($testToken) {
    Write-Host "5. Testing Token Validation..." -ForegroundColor Yellow
    try {
        $validate = Invoke-RestMethod -Uri "$baseUrl/api/auth/validate" `
            -Method POST `
            -Headers @{Authorization = "Bearer $testToken"} `
            -ErrorAction Stop
        
        if ($validate.success) {
            Write-Host "   ✓ Token validation successful" -ForegroundColor Green
        } else {
            Write-Host "   ⚠ Token validation returned unexpected format" -ForegroundColor Yellow
            $allTestsPassed = $false
        }
    } catch {
        Write-Host "   ✗ Token validation failed" -ForegroundColor Red
        Write-Host "   Error: $_" -ForegroundColor Red
        $allTestsPassed = $false
    }
    Write-Host ""
}

# Test 6: Test Get User Profile
if ($testToken) {
    Write-Host "6. Testing Get User Profile..." -ForegroundColor Yellow
    try {
        $profile = Invoke-RestMethod -Uri "$baseUrl/api/auth/profile" `
            -Method GET `
            -Headers @{Authorization = "Bearer $testToken"} `
            -ErrorAction Stop
        
        if ($profile.success -and $profile.user) {
            Write-Host "   ✓ Get profile successful" -ForegroundColor Green
            Write-Host "   Name: $($profile.user.name)" -ForegroundColor Gray
            Write-Host "   Email: $($profile.user.email)" -ForegroundColor Gray
            Write-Host "   Role: $($profile.user.role)" -ForegroundColor Gray
        } else {
            Write-Host "   ⚠ Get profile returned unexpected format" -ForegroundColor Yellow
            $allTestsPassed = $false
        }
    } catch {
        Write-Host "   ✗ Get profile failed" -ForegroundColor Red
        Write-Host "   Error: $_" -ForegroundColor Red
        $allTestsPassed = $false
    }
    Write-Host ""
}

# Test 7: Test Create Booking
if ($testToken -and $testUserId -and $buildings.buildings.Count -gt 0) {
    Write-Host "7. Testing Create Booking..." -ForegroundColor Yellow
    
    # Find first available bed
    $building = $buildings.buildings[0]
    $room = $building.rooms | Where-Object { $_.available_beds -gt 0 } | Select-Object -First 1
    
    if ($room) {
        $bed = $room.beds | Where-Object { -not $_.is_occupied } | Select-Object -First 1
        
        if ($bed) {
            $bookingBody = @{
                user_id = $testUserId
                user_name = "Test User"
                building_id = $building.id
                building_name = $building.name
                room_id = $room.id
                room_number = $room.number
                bed_id = $bed.id
                bed_number = $bed.number
            } | ConvertTo-Json
            
            try {
                $booking = Invoke-RestMethod -Uri "$baseUrl/api/bookings" `
                    -Method POST `
                    -ContentType "application/json" `
                    -Headers @{Authorization = "Bearer $testToken"} `
                    -Body $bookingBody `
                    -ErrorAction Stop
                
                if ($booking.success -and $booking.booking) {
                    Write-Host "   ✓ Booking created successfully" -ForegroundColor Green
                    Write-Host "   Booking ID: $($booking.booking.id)" -ForegroundColor Gray
                    Write-Host "   Building: $($booking.booking.building_name)" -ForegroundColor Gray
                    Write-Host "   Room: $($booking.booking.room_number)" -ForegroundColor Gray
                    Write-Host "   Bed: $($booking.booking.bed_number)" -ForegroundColor Gray
                    $testBookingId = $booking.booking.id
                } else {
                    Write-Host "   ⚠ Booking returned unexpected format" -ForegroundColor Yellow
                    $allTestsPassed = $false
                }
            } catch {
                Write-Host "   ✗ Booking creation failed" -ForegroundColor Red
                Write-Host "   Error: $_" -ForegroundColor Red
                $allTestsPassed = $false
            }
        } else {
            Write-Host "   ⚠ No available beds found" -ForegroundColor Yellow
        }
    } else {
        Write-Host "   ⚠ No available rooms found" -ForegroundColor Yellow
    }
    Write-Host ""
}

# Test 8: Test Get User Bookings
if ($testToken -and $testUserId) {
    Write-Host "8. Testing Get User Bookings..." -ForegroundColor Yellow
    try {
        $bookings = Invoke-RestMethod -Uri "$baseUrl/api/bookings/users/$testUserId" `
            -Method GET `
            -Headers @{Authorization = "Bearer $testToken"} `
            -ErrorAction Stop
        
        if ($bookings.success) {
            Write-Host "   ✓ Get bookings successful" -ForegroundColor Green
            Write-Host "   Found $($bookings.bookings.Count) booking(s)" -ForegroundColor Gray
        } else {
            Write-Host "   ⚠ Get bookings returned unexpected format" -ForegroundColor Yellow
            $allTestsPassed = $false
        }
    } catch {
        Write-Host "   ✗ Get bookings failed" -ForegroundColor Red
        Write-Host "   Error: $_" -ForegroundColor Red
        $allTestsPassed = $false
    }
    Write-Host ""
}

# Test 9: Test Cancel Booking
if ($testToken -and $testBookingId) {
    Write-Host "9. Testing Cancel Booking..." -ForegroundColor Yellow
    try {
        $cancel = Invoke-RestMethod -Uri "$baseUrl/api/bookings/$testBookingId/cancel" `
            -Method PUT `
            -Headers @{Authorization = "Bearer $testToken"} `
            -ErrorAction Stop
        
        if ($cancel.success) {
            Write-Host "   ✓ Booking cancelled successfully" -ForegroundColor Green
        } else {
            Write-Host "   ⚠ Cancel booking returned unexpected format" -ForegroundColor Yellow
            $allTestsPassed = $false
        }
    } catch {
        Write-Host "   ✗ Cancel booking failed" -ForegroundColor Red
        Write-Host "   Error: $_" -ForegroundColor Red
        $allTestsPassed = $false
    }
    Write-Host ""
}

# Summary
Write-Host "════════════════════════════════════════" -ForegroundColor Cyan
if ($allTestsPassed) {
    Write-Host "✅ All tests passed!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Backend is properly configured and working." -ForegroundColor Green
    Write-Host "You can now start the frontend application:" -ForegroundColor Green
    Write-Host "  pnpm dev" -ForegroundColor Yellow
} else {
    Write-Host "⚠️  Some tests failed!" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Please check the errors above and ensure:" -ForegroundColor Yellow
    Write-Host "  1. Backend services are running (docker-compose up -d)" -ForegroundColor Yellow
    Write-Host "  2. All services are healthy" -ForegroundColor Yellow
    Write-Host "  3. Ports are not blocked by firewall" -ForegroundColor Yellow
}
Write-Host "════════════════════════════════════════" -ForegroundColor Cyan
