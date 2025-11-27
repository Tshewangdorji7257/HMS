# Backend Integration Guide

This document explains how the Next.js frontend is connected to the Golang microservices backend.

## 🏗️ Architecture

```
Frontend (Next.js) ──→ API Gateway (Port 8000) ──┬──→ Auth Service (Port 8001)
                                                  ├──→ Building Service (Port 8002)
                                                  └──→ Booking Service (Port 8003)
```

## 🚀 Getting Started

### Prerequisites

1. **Backend Services Running**
   - Clone the backend repository: https://github.com/Tshewangdorji7257/micro_services-backend
   - Start all services using Docker Compose:
     ```powershell
     cd backend
     docker-compose up -d
     ```

2. **Verify Backend Health**
   ```powershell
   # Check API Gateway
   Invoke-RestMethod -Uri "http://localhost:8000/health"
   
   # Check available buildings
   Invoke-RestMethod -Uri "http://localhost:8000/api/buildings"
   ```

### Frontend Setup

1. **Install Dependencies**
   ```powershell
   pnpm install
   ```

2. **Environment Configuration**
   
   The `.env.local` file is already configured:
   ```env
   NEXT_PUBLIC_API_URL=http://localhost:8000
   NEXT_PUBLIC_API_BASE_URL=http://localhost:8000/api
   ```

3. **Start Development Server**
   ```powershell
   pnpm dev
   ```

4. **Access Application**
   - Frontend: http://localhost:3000
   - Backend API Gateway: http://localhost:8000

## 📁 Key Files Updated

### 1. API Configuration (`lib/api-config.ts`)
- Centralized API configuration
- Helper functions for authenticated requests
- Error handling utilities

### 2. Authentication (`lib/auth.ts`)
- User authentication with JWT tokens
- Login and signup with backend API
- Token management and validation
- Profile retrieval

### 3. Admin Authentication (`lib/admin.ts`)
- Admin-specific authentication
- Role-based access control
- Separate token storage

### 4. Data Service (`lib/data.ts`)
- Building and room data fetching
- Booking creation and management
- Real-time data synchronization
- Caching for offline support

### 5. Next.js Config (`next.config.mjs`)
- API proxy configuration
- Environment variable setup
- CORS handling

## 🔐 Authentication Flow

### Student/User Authentication

1. **Signup**
   ```typescript
   import { authService } from '@/lib/auth'
   
   const result = await authService.signup(
     'student@example.com',
     'password123',
     'John Doe',
     'student'
   )
   ```

2. **Login**
   ```typescript
   const result = await authService.login(
     'student@example.com',
     'password123'
   )
   ```

3. **Get Auth State**
   ```typescript
   const { user, isAuthenticated, token } = authService.getAuthState()
   ```

### Admin Authentication

```typescript
import { adminAuthService } from '@/lib/admin'

// Admin signup with 'admin' role
const result = await adminAuthService.signup(
  'admin@example.com',
  'password123',
  'Admin User'
)

// Admin login
const result = await adminAuthService.login(
  'admin@example.com',
  'password123'
)
```

## 🏢 Data Operations

### Fetch Buildings

```typescript
import { hostelDataService } from '@/lib/data'

// Get all buildings
const buildings = await hostelDataService.getBuildings()

// Get specific building
const building = await hostelDataService.getBuildingById('bldg-1')

// Search buildings
const results = await hostelDataService.searchBuildings('RK')
```

### Manage Bookings

```typescript
// Create booking
const result = await hostelDataService.createBooking({
  buildingId: 'bldg-1',
  buildingName: 'RK A',
  roomId: 'bldg-1-room-001',
  roomNumber: '001',
  bedId: 'bldg-1-room-001-bed-1',
  bedNumber: 1
})

// Get user bookings
const bookings = await hostelDataService.getBookings()

// Cancel booking
const result = await hostelDataService.cancelBooking('booking-id')
```

## 🔄 Data Flow

### 1. User Signup/Login
```
Component → authService.login() → POST /api/auth/login → Backend
                                                       ↓
Component ← Token + User Data ← Response ← Auth Service
```

### 2. Fetch Buildings
```
Component → hostelDataService.getBuildings() → GET /api/buildings → Backend
                                                                   ↓
Component ← Buildings Array ← Response ← Building Service
```

### 3. Create Booking
```
Component → hostelDataService.createBooking() → POST /api/bookings → Backend
                                              (with JWT token)      ↓
Component ← Booking Data ← Response ← Booking Service
```

## 🛠️ Backend API Endpoints

### Authentication
- `POST /api/auth/signup` - Create new user
- `POST /api/auth/login` - User login
- `POST /api/auth/validate` - Validate token
- `GET /api/auth/profile` - Get user profile (requires auth)

### Buildings
- `GET /api/buildings` - Get all buildings
- `GET /api/buildings/:id` - Get specific building
- `GET /api/buildings/:buildingId/rooms/:roomId` - Get specific room
- `GET /api/buildings/search?q=query` - Search buildings

### Bookings
- `POST /api/bookings` - Create booking (requires auth)
- `GET /api/bookings/users/:userId` - Get user bookings (requires auth)
- `GET /api/bookings` - Get all bookings (admin only)
- `PUT /api/bookings/:id/cancel` - Cancel booking (requires auth)

## 🧪 Testing the Integration

### 1. Health Check
```powershell
# Test backend connectivity
Invoke-RestMethod -Uri "http://localhost:8000/health"
```

### 2. Create Test User
```powershell
$body = @{
    email = "test@example.com"
    password = "password123"
    name = "Test User"
    role = "student"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8000/api/auth/signup" `
    -Method POST `
    -ContentType "application/json" `
    -Body $body

Write-Host "Token: $($response.token)"
```

### 3. Fetch Buildings
```powershell
$buildings = Invoke-RestMethod -Uri "http://localhost:8000/api/buildings"
$buildings.buildings | Select-Object name, total_rooms, available_beds
```

### 4. Create Booking
```powershell
$token = "your-jwt-token-here"
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

Invoke-RestMethod -Uri "http://localhost:8000/api/bookings" `
    -Method POST `
    -ContentType "application/json" `
    -Headers @{Authorization = "Bearer $token"} `
    -Body $bookingBody
```

## 🐛 Troubleshooting

### Backend Not Running
**Error**: `fetch failed` or connection refused

**Solution**:
```powershell
# Check if backend services are running
cd backend
docker-compose ps

# Start services if not running
docker-compose up -d

# Check logs
docker-compose logs -f
```

### CORS Issues
**Error**: CORS policy error in browser console

**Solution**: The API Gateway handles CORS. Ensure:
1. Backend is running on `localhost:8000`
2. Frontend is running on `localhost:3000`
3. API Gateway CORS settings allow `localhost:3000`

### Authentication Errors
**Error**: 401 Unauthorized

**Solution**:
1. Check if token is stored: `localStorage.getItem('hostel-auth-token')`
2. Verify token is being sent in Authorization header
3. Token may have expired - login again

### Data Not Loading
**Error**: Empty data or cached data

**Solution**:
```typescript
// Clear cache and fetch fresh data
localStorage.removeItem('hostel-buildings-cache')
localStorage.removeItem('hostel-bookings-cache')

// Reload the page
window.location.reload()
```

## 📝 Development Tips

### 1. Debug API Calls
Add console logs in `lib/api-config.ts`:
```typescript
export async function apiFetch<T>(endpoint: string, options: RequestInit = {}) {
  console.log('API Request:', endpoint, options)
  // ... existing code
  console.log('API Response:', data)
  return data
}
```

### 2. Monitor Backend Logs
```powershell
# Watch all service logs
docker-compose logs -f

# Watch specific service
docker-compose logs -f api-gateway
docker-compose logs -f auth-service
```

### 3. Test with cURL
```powershell
# Test signup
curl -X POST http://localhost:8000/api/auth/signup `
  -H "Content-Type: application/json" `
  -d '{\"email\":\"test@example.com\",\"password\":\"test123\",\"name\":\"Test User\",\"role\":\"student\"}'

# Test buildings
curl http://localhost:8000/api/buildings
```

## 🚀 Production Deployment

### Environment Variables
Update `.env.local` for production:
```env
NEXT_PUBLIC_API_URL=https://your-backend-api.com
NEXT_PUBLIC_API_BASE_URL=https://your-backend-api.com/api
```

### Build for Production
```powershell
pnpm build
pnpm start
```

## 📚 Additional Resources

- Backend Repository: https://github.com/Tshewangdorji7257/micro_services-backend
- Backend API Documentation: See `backend/API_TESTING.md`
- Backend Setup Guide: See `backend/SETUP.md`
- Database Schema: See `backend/DATABASE_SCHEMA.md`

## ✅ Integration Checklist

- [x] Backend services running (Docker Compose)
- [x] API configuration files created
- [x] Authentication services updated
- [x] Data service updated
- [x] Environment variables configured
- [x] Next.js config updated with API proxy
- [ ] Backend health check passing
- [ ] User signup/login working
- [ ] Buildings data loading
- [ ] Booking creation working
- [ ] Admin panel functional

## 🎉 Success!

Your Next.js frontend is now fully integrated with the Golang microservices backend!

Test the complete flow:
1. Start backend: `docker-compose up -d`
2. Start frontend: `pnpm dev`
3. Visit: http://localhost:3000
4. Sign up → Browse buildings → Make a booking

---

**Need Help?** Check the troubleshooting section or review the backend documentation.
