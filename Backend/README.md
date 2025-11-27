# Hostel Management System - Microservices Backend

A comprehensive microservices-based backend for a Hostel Management System built with **Golang**, featuring JWT authentication, role-based access control, and RESTful APIs.

## 🏗️ Architecture Overview

This backend follows a **microservices architecture** with three independent services:

```
┌─────────────────┐
│   Frontend      │
│  (Next.js)      │
└────────┬────────┘
         │
    ┌────▼─────┐
    │   API    │
    │ Gateway  │  Port 8000
    │  (Proxy) │
    └────┬─────┘
         │
    ┌────┴──────────────────────────────┐
    │                                   │
┌───▼────────┐  ┌────────────┐  ┌──────▼────┐
│   Auth     │  │  Building  │  │  Booking  │
│  Service   │  │  Service   │  │  Service  │
│  Port 8001 │  │ Port 8002  │  │ Port 8003 │
└─────┬──────┘  └─────┬──────┘  └─────┬─────┘
      │               │               │
┌─────▼──────┐  ┌─────▼──────┐  ┌─────▼─────┐
│ Auth DB    │  │Building DB │  │Booking DB │
│ PostgreSQL │  │ PostgreSQL │  │PostgreSQL │
│ Port 5432  │  │ Port 5433  │  │ Port 5434 │
└────────────┘  └────────────┘  └───────────┘
```

## 📦 Microservices

### 1. **Auth Service** (Port 8001)
- User authentication and authorization
- JWT token generation and validation
- Role-based access control (Student/Admin)
- User profile management

**Database**: `hostel_auth_db` (PostgreSQL)

**Key Features**:
- Secure password hashing with bcrypt
- JWT with configurable expiry
- Email validation
- User session management

### 2. **Building & Room Service** (Port 8002)
- Building management (CRUD operations)
- Room management with multiple types
- Bed inventory and availability tracking
- Real-time occupancy updates

**Database**: `hostel_building_db` (PostgreSQL)

**Key Features**:
- 10 pre-seeded buildings (RK A, RK B, H A-F, NK, Lhawang)
- Support for single, double, triple, and quad rooms
- Amenities tracking (Wi-Fi, Gym, Laundry, etc.)
- Bed occupancy status management

### 3. **Booking Service** (Port 8003)
- Room booking and reservations
- Booking cancellation
- User booking history
- Integration with Building Service for occupancy updates

**Database**: `hostel_booking_db` (PostgreSQL)

**Key Features**:
- Single active booking per user enforcement
- Automatic bed status synchronization
- Booking status tracking (active/cancelled)
- User-specific booking retrieval

### 4. **API Gateway** (Port 8000)
- Single entry point for all services
- Request routing and load balancing
- Service health monitoring
- CORS handling

## 🚀 Getting Started

### Prerequisites

- **Go 1.21** or higher
- **PostgreSQL 15** or higher
- **Docker & Docker Compose** (optional, for containerized deployment)
- **Git**

### Option 1: Docker Compose (Recommended)

1. **Clone the repository**
```powershell
git clone <repository-url>
cd backend
```

2. **Build and start all services**
```powershell
docker-compose up --build
```

3. **Services will be available at:**
- API Gateway: http://localhost:8000
- Auth Service: http://localhost:8001
- Building Service: http://localhost:8002
- Booking Service: http://localhost:8003

4. **Health check**
```powershell
curl http://localhost:8000/health
```

### Option 2: Manual Setup

#### Step 1: Setup PostgreSQL Databases

Create three separate databases:

```sql
CREATE DATABASE hostel_auth_db;
CREATE DATABASE hostel_building_db;
CREATE DATABASE hostel_booking_db;
```

#### Step 2: Configure Environment Variables

For each service, copy `.env.example` to `.env` and update:

**auth-service/.env**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=hostel_auth_db
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRY=24h
PORT=8001
```

**building-service/.env**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=hostel_building_db
PORT=8002
AUTH_SERVICE_URL=http://localhost:8001
```

**booking-service/.env**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=hostel_booking_db
PORT=8003
AUTH_SERVICE_URL=http://localhost:8001
BUILDING_SERVICE_URL=http://localhost:8002
```

**api-gateway/.env**
```env
PORT=8000
AUTH_SERVICE_URL=http://localhost:8001
BUILDING_SERVICE_URL=http://localhost:8002
BOOKING_SERVICE_URL=http://localhost:8003
```

#### Step 3: Install Dependencies and Run Services

For each service (auth-service, building-service, booking-service, api-gateway):

```powershell
cd <service-folder>
go mod download
go run main.go
```

## 📚 API Documentation

### Base URL
```
http://localhost:8000
```

### Authentication Endpoints

#### 1. **Sign Up**
```http
POST /api/auth/signup
Content-Type: application/json

{
  "email": "student@example.com",
  "password": "password123",
  "name": "John Doe",
  "role": "student"
}
```

**Response**:
```json
{
  "success": true,
  "message": "User created successfully",
  "token": "eyJhbGc...",
  "user": {
    "id": "uuid",
    "email": "student@example.com",
    "name": "John Doe",
    "role": "student"
  }
}
```

#### 2. **Login**
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "student@example.com",
  "password": "password123"
}
```

#### 3. **Validate Token**
```http
POST /api/auth/validate
Authorization: Bearer <token>
```

#### 4. **Get User Profile**
```http
GET /api/auth/profile
Authorization: Bearer <token>
```

### Building Endpoints

#### 1. **Get All Buildings**
```http
GET /api/buildings
```

**Response**:
```json
{
  "success": true,
  "buildings": [
    {
      "id": "bldg-1",
      "name": "RK A",
      "description": "Modern residence...",
      "total_rooms": 25,
      "total_beds": 50,
      "available_beds": 30,
      "amenities": ["Wi-Fi", "Gym", "Laundry"],
      "rooms": [...]
    }
  ]
}
```

#### 2. **Get Building by ID**
```http
GET /api/buildings/{buildingId}
```

#### 3. **Get Room by ID**
```http
GET /api/buildings/{buildingId}/rooms/{roomId}
```

#### 4. **Search Buildings**
```http
GET /api/buildings/search?q=RK
```

#### 5. **Update Bed Occupancy** (Internal use)
```http
PUT /api/buildings/beds/{bedId}/occupancy
Content-Type: application/json

{
  "is_occupied": true,
  "occupied_by": "user-id",
  "occupied_by_name": "John Doe"
}
```

### Booking Endpoints

#### 1. **Create Booking**
```http
POST /api/bookings
Content-Type: application/json
Authorization: Bearer <token>

{
  "user_id": "user-uuid",
  "user_name": "John Doe",
  "building_id": "bldg-1",
  "building_name": "RK A",
  "room_id": "bldg-1-room-001",
  "room_number": "001",
  "bed_id": "bldg-1-room-001-bed-1",
  "bed_number": 1
}
```

**Response**:
```json
{
  "success": true,
  "message": "Booking created successfully",
  "booking": {
    "id": "booking-uuid",
    "user_id": "user-uuid",
    "building_name": "RK A",
    "room_number": "001",
    "bed_number": 1,
    "status": "active",
    "booking_date": "2025-11-25T10:30:00Z"
  }
}
```

#### 2. **Get All Bookings**
```http
GET /api/bookings
Authorization: Bearer <token>
```

#### 3. **Get User Bookings**
```http
GET /api/bookings/users/{userId}
Authorization: Bearer <token>
```

#### 4. **Get Booking by ID**
```http
GET /api/bookings/{bookingId}
Authorization: Bearer <token>
```

#### 5. **Cancel Booking**
```http
PUT /api/bookings/{bookingId}/cancel
Authorization: Bearer <token>
```

## 🔐 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt for secure password storage
- **Role-Based Access Control**: Student and Admin roles
- **CORS Configuration**: Configurable cross-origin requests
- **Input Validation**: Request validation at all endpoints
- **SQL Injection Protection**: Prepared statements

## 🛠️ Technology Stack

- **Language**: Go 1.21
- **Database**: PostgreSQL 15
- **Libraries**:
  - `gorilla/mux` - HTTP router
  - `lib/pq` - PostgreSQL driver
  - `golang-jwt/jwt` - JWT implementation
  - `golang.org/x/crypto` - Password hashing
  - `rs/cors` - CORS middleware
  - `joho/godotenv` - Environment variables
  - `google/uuid` - UUID generation

## 📁 Project Structure

```
backend/
├── api-gateway/
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   └── .env.example
├── auth-service/
│   ├── main.go
│   ├── models/
│   │   └── user.go
│   ├── handlers/
│   │   └── auth.go
│   ├── middleware/
│   │   └── auth.go
│   ├── utils/
│   │   ├── jwt.go
│   │   └── password.go
│   ├── database/
│   │   └── postgres.go
│   ├── go.mod
│   ├── Dockerfile
│   └── .env.example
├── building-service/
│   ├── main.go
│   ├── models/
│   │   └── building.go
│   ├── handlers/
│   │   └── building.go
│   ├── database/
│   │   └── postgres.go
│   ├── go.mod
│   ├── Dockerfile
│   └── .env.example
├── booking-service/
│   ├── main.go
│   ├── models/
│   │   └── booking.go
│   ├── handlers/
│   │   └── booking.go
│   ├── utils/
│   │   └── config.go
│   ├── database/
│   │   └── postgres.go
│   ├── go.mod
│   ├── Dockerfile
│   └── .env.example
├── docker-compose.yml
└── README.md
```

## 🧪 Testing

### Manual Testing with cURL

**1. Create a user:**
```powershell
curl -X POST http://localhost:8000/api/auth/signup `
  -H "Content-Type: application/json" `
  -d '{\"email\":\"test@example.com\",\"password\":\"password123\",\"name\":\"Test User\"}'
```

**2. Get all buildings:**
```powershell
curl http://localhost:8000/api/buildings
```

**3. Create a booking:**
```powershell
curl -X POST http://localhost:8000/api/bookings `
  -H "Content-Type: application/json" `
  -H "Authorization: Bearer YOUR_TOKEN" `
  -d '{\"user_id\":\"user-id\",\"user_name\":\"Test User\",\"building_id\":\"bldg-1\",\"building_name\":\"RK A\",\"room_id\":\"bldg-1-room-001\",\"room_number\":\"001\",\"bed_id\":\"bldg-1-room-001-bed-1\",\"bed_number\":1}'
```

### Testing with Postman

Import the API endpoints into Postman:
1. Create a new collection
2. Add requests for each endpoint
3. Set up environment variables for base URL and tokens

## 🔄 Database Schema

### Auth Service

**users table**:
- `id` (VARCHAR, PK)
- `email` (VARCHAR, UNIQUE)
- `name` (VARCHAR)
- `password` (TEXT)
- `role` (VARCHAR) - 'student' or 'admin'
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Building Service

**buildings table**:
- `id` (VARCHAR, PK)
- `name` (VARCHAR)
- `description` (TEXT)
- `total_rooms` (INT)
- `total_beds` (INT)
- `available_beds` (INT)
- `amenities` (JSONB)
- `image` (TEXT)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

**rooms table**:
- `id` (VARCHAR, PK)
- `building_id` (VARCHAR, FK)
- `number` (VARCHAR)
- `type` (VARCHAR) - 'single', 'double', 'triple', 'quad'
- `total_beds` (INT)
- `available_beds` (INT)
- `amenities` (JSONB)
- `price` (DECIMAL)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

**beds table**:
- `id` (VARCHAR, PK)
- `room_id` (VARCHAR, FK)
- `number` (INT)
- `is_occupied` (BOOLEAN)
- `occupied_by` (VARCHAR, nullable)
- `occupied_by_name` (VARCHAR, nullable)

### Booking Service

**bookings table**:
- `id` (VARCHAR, PK)
- `user_id` (VARCHAR)
- `user_name` (VARCHAR)
- `building_id` (VARCHAR)
- `building_name` (VARCHAR)
- `room_id` (VARCHAR)
- `room_number` (VARCHAR)
- `bed_id` (VARCHAR)
- `bed_number` (INT)
- `booking_date` (TIMESTAMP)
- `status` (VARCHAR) - 'active' or 'cancelled'
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

## 🚢 Deployment

### Production Considerations

1. **Environment Variables**:
   - Use strong JWT secrets
   - Secure database credentials
   - Configure proper CORS origins

2. **Database**:
   - Use managed PostgreSQL services
   - Enable SSL connections
   - Set up regular backups
   - Configure connection pooling

3. **Security**:
   - Use HTTPS in production
   - Implement rate limiting
   - Add request logging
   - Set up monitoring and alerts

4. **Scaling**:
   - Use container orchestration (Kubernetes)
   - Implement service discovery
   - Add load balancers
   - Use caching (Redis)

## 🤝 Integration with Frontend

The backend is designed to work seamlessly with the Next.js frontend from the repository:
```
https://github.com/Tshewangdorji7257/Hostel-Management-System
```

**Configuration**:
Update the frontend API base URL to point to the API Gateway:
```typescript
// In frontend config
const API_BASE_URL = "http://localhost:8000";
```

## 📝 Development

### Adding a New Endpoint

1. Define the model in `models/`
2. Create handler in `handlers/`
3. Add route in `main.go`
4. Update API documentation

### Adding a New Service

1. Create new service directory
2. Implement service with proper structure
3. Add database configuration
4. Update `docker-compose.yml`
5. Update API Gateway routes

## 🐛 Troubleshooting

**Database connection errors**:
- Verify PostgreSQL is running
- Check connection credentials
- Ensure database exists

**Port already in use**:
```powershell
# Find and kill process using port 8000
Get-Process -Id (Get-NetTCPConnection -LocalPort 8000).OwningProcess | Stop-Process -Force
```

**Docker issues**:
```powershell
# Rebuild containers
docker-compose down -v
docker-compose up --build
```

## 📄 License

This project is provided as-is for educational purposes.

## 👥 Contributors

Built as a comprehensive microservices backend for Hostel Management System.

---

**Made with ❤️ using Go and PostgreSQL**
