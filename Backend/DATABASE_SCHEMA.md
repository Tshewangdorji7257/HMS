# Database Schema Documentation

Complete database schema for the Hostel Management System microservices.

## Overview

The system uses **PostgreSQL 15** with separate databases for each microservice to ensure data isolation and independent scalability.

### Databases:
1. **hostel_auth_db** - User authentication and profiles
2. **hostel_building_db** - Buildings, rooms, and beds inventory
3. **hostel_booking_db** - Booking records and reservations

## Database: hostel_auth_db

### Table: users

Stores user account information and authentication data.

```sql
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'student',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
```

#### Columns:

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | VARCHAR(255) | PRIMARY KEY | Unique user identifier (UUID) |
| email | VARCHAR(255) | UNIQUE, NOT NULL | User email address |
| name | VARCHAR(255) | NOT NULL | User's full name |
| password | TEXT | NOT NULL | Bcrypt hashed password |
| role | VARCHAR(50) | NOT NULL | User role: 'student' or 'admin' |
| created_at | TIMESTAMP | DEFAULT NOW | Account creation timestamp |
| updated_at | TIMESTAMP | DEFAULT NOW | Last update timestamp |

#### Sample Data:

```sql
INSERT INTO users (id, email, name, password, role) VALUES 
('550e8400-e29b-41d4-a716-446655440000', 'admin@example.com', 'Admin User', '$2a$10$...', 'admin'),
('550e8400-e29b-41d4-a716-446655440001', 'student@example.com', 'John Doe', '$2a$10$...', 'student');
```

#### Queries:

**Create User:**
```sql
INSERT INTO users (id, email, name, password, role, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
```

**Find by Email:**
```sql
SELECT * FROM users WHERE email = $1;
```

**Find by ID:**
```sql
SELECT * FROM users WHERE id = $1;
```

**Update User:**
```sql
UPDATE users SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2;
```

---

## Database: hostel_building_db

### Table: buildings

Stores hostel building information.

```sql
CREATE TABLE IF NOT EXISTS buildings (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    total_rooms INTEGER NOT NULL DEFAULT 0,
    total_beds INTEGER NOT NULL DEFAULT 0,
    available_beds INTEGER NOT NULL DEFAULT 0,
    amenities JSONB,
    image TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_buildings_name ON buildings(name);
CREATE INDEX idx_buildings_available_beds ON buildings(available_beds);
```

#### Columns:

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | VARCHAR(255) | PRIMARY KEY | Building identifier (e.g., 'bldg-1') |
| name | VARCHAR(255) | NOT NULL | Building name (e.g., 'RK A') |
| description | TEXT | | Building description |
| total_rooms | INTEGER | NOT NULL | Total number of rooms |
| total_beds | INTEGER | NOT NULL | Total bed capacity |
| available_beds | INTEGER | NOT NULL | Current available beds |
| amenities | JSONB | | JSON array of amenities |
| image | TEXT | | Building image URL |
| created_at | TIMESTAMP | DEFAULT NOW | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT NOW | Last update timestamp |

#### Amenities JSON Format:

```json
["Wi-Fi", "Gym", "Laundry", "Common Room", "Study Area", "Parking"]
```

### Table: rooms

Stores room information within buildings.

```sql
CREATE TABLE IF NOT EXISTS rooms (
    id VARCHAR(255) PRIMARY KEY,
    building_id VARCHAR(255) NOT NULL REFERENCES buildings(id),
    number VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    total_beds INTEGER NOT NULL,
    available_beds INTEGER NOT NULL DEFAULT 0,
    amenities JSONB,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(building_id, number)
);

CREATE INDEX idx_rooms_building ON rooms(building_id);
CREATE INDEX idx_rooms_type ON rooms(type);
CREATE INDEX idx_rooms_available_beds ON rooms(available_beds);
```

#### Columns:

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | VARCHAR(255) | PRIMARY KEY | Room identifier |
| building_id | VARCHAR(255) | FK, NOT NULL | Reference to buildings table |
| number | VARCHAR(50) | NOT NULL | Room number |
| type | VARCHAR(50) | NOT NULL | 'single', 'double', 'triple', 'quad' |
| total_beds | INTEGER | NOT NULL | Number of beds in room |
| available_beds | INTEGER | NOT NULL | Available beds |
| amenities | JSONB | | Room-specific amenities |
| price | DECIMAL(10,2) | NOT NULL | Price per bed per month |
| created_at | TIMESTAMP | DEFAULT NOW | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT NOW | Last update timestamp |

#### Room Types:

- **single**: 1 bed
- **double**: 2 beds
- **triple**: 3 beds
- **quad**: 4 beds

### Table: beds

Stores individual bed information and occupancy status.

```sql
CREATE TABLE IF NOT EXISTS beds (
    id VARCHAR(255) PRIMARY KEY,
    room_id VARCHAR(255) NOT NULL REFERENCES rooms(id),
    number INTEGER NOT NULL,
    is_occupied BOOLEAN DEFAULT FALSE,
    occupied_by VARCHAR(255),
    occupied_by_name VARCHAR(255),
    UNIQUE(room_id, number)
);

CREATE INDEX idx_beds_room ON beds(room_id);
CREATE INDEX idx_beds_occupied ON beds(is_occupied);
CREATE INDEX idx_beds_occupied_by ON beds(occupied_by);
```

#### Columns:

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | VARCHAR(255) | PRIMARY KEY | Bed identifier |
| room_id | VARCHAR(255) | FK, NOT NULL | Reference to rooms table |
| number | INTEGER | NOT NULL | Bed number within room |
| is_occupied | BOOLEAN | DEFAULT FALSE | Occupancy status |
| occupied_by | VARCHAR(255) | | User ID if occupied |
| occupied_by_name | VARCHAR(255) | | User name if occupied |

#### Sample Data:

```sql
-- Building
INSERT INTO buildings (id, name, description, total_rooms, total_beds, available_beds, amenities, image) 
VALUES (
    'bldg-1', 
    'RK A', 
    'Modern residence with excellent facilities',
    25, 
    50, 
    30,
    '["Wi-Fi", "Gym", "Laundry", "Common Room"]'::jsonb,
    '/images/rk-a.jpg'
);

-- Room
INSERT INTO rooms (id, building_id, number, type, total_beds, available_beds, amenities, price) 
VALUES (
    'bldg-1-room-001',
    'bldg-1',
    '001',
    'double',
    2,
    2,
    '["Window", "Desk", "Closet"]'::jsonb,
    5000.00
);

-- Beds
INSERT INTO beds (id, room_id, number, is_occupied) VALUES 
('bldg-1-room-001-bed-1', 'bldg-1-room-001', 1, FALSE),
('bldg-1-room-001-bed-2', 'bldg-1-room-001', 2, FALSE);
```

#### Queries:

**Get All Buildings with Rooms:**
```sql
SELECT 
    b.*,
    json_agg(
        json_build_object(
            'id', r.id,
            'number', r.number,
            'type', r.type,
            'total_beds', r.total_beds,
            'available_beds', r.available_beds,
            'price', r.price
        )
    ) as rooms
FROM buildings b
LEFT JOIN rooms r ON b.id = r.building_id
GROUP BY b.id;
```

**Update Bed Occupancy:**
```sql
UPDATE beds 
SET is_occupied = $1, occupied_by = $2, occupied_by_name = $3 
WHERE id = $4;
```

**Update Room Available Beds:**
```sql
UPDATE rooms 
SET available_beds = available_beds + $1 
WHERE id = $2;
```

---

## Database: hostel_booking_db

### Table: bookings

Stores all booking records.

```sql
CREATE TABLE IF NOT EXISTS bookings (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    building_id VARCHAR(255) NOT NULL,
    building_name VARCHAR(255) NOT NULL,
    room_id VARCHAR(255) NOT NULL,
    room_number VARCHAR(50) NOT NULL,
    bed_id VARCHAR(255) NOT NULL,
    bed_number INTEGER NOT NULL,
    booking_date TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bookings_user ON bookings(user_id);
CREATE INDEX idx_bookings_building ON bookings(building_id);
CREATE INDEX idx_bookings_room ON bookings(room_id);
CREATE INDEX idx_bookings_bed ON bookings(bed_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_date ON bookings(booking_date);
```

#### Columns:

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | VARCHAR(255) | PRIMARY KEY | Booking identifier (UUID) |
| user_id | VARCHAR(255) | NOT NULL | User who made the booking |
| user_name | VARCHAR(255) | NOT NULL | User's name (denormalized) |
| building_id | VARCHAR(255) | NOT NULL | Building reference |
| building_name | VARCHAR(255) | NOT NULL | Building name (denormalized) |
| room_id | VARCHAR(255) | NOT NULL | Room reference |
| room_number | VARCHAR(50) | NOT NULL | Room number (denormalized) |
| bed_id | VARCHAR(255) | NOT NULL | Bed reference |
| bed_number | INTEGER | NOT NULL | Bed number (denormalized) |
| booking_date | TIMESTAMP | NOT NULL | Date of booking |
| status | VARCHAR(50) | NOT NULL | 'active' or 'cancelled' |
| created_at | TIMESTAMP | DEFAULT NOW | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT NOW | Last update timestamp |

#### Status Values:

- **active**: Current active booking
- **cancelled**: Cancelled booking

#### Sample Data:

```sql
INSERT INTO bookings (
    id, user_id, user_name, 
    building_id, building_name, 
    room_id, room_number, 
    bed_id, bed_number, 
    booking_date, status
) VALUES (
    '650e8400-e29b-41d4-a716-446655440000',
    '550e8400-e29b-41d4-a716-446655440001',
    'John Doe',
    'bldg-1',
    'RK A',
    'bldg-1-room-001',
    '001',
    'bldg-1-room-001-bed-1',
    1,
    CURRENT_TIMESTAMP,
    'active'
);
```

#### Queries:

**Create Booking:**
```sql
INSERT INTO bookings (
    id, user_id, user_name, 
    building_id, building_name, 
    room_id, room_number, 
    bed_id, bed_number, 
    booking_date, status, 
    created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);
```

**Get User Active Bookings:**
```sql
SELECT * FROM bookings 
WHERE user_id = $1 AND status = 'active';
```

**Get All Bookings (Paginated):**
```sql
SELECT * FROM bookings 
ORDER BY created_at DESC 
LIMIT $1 OFFSET $2;
```

**Cancel Booking:**
```sql
UPDATE bookings 
SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP 
WHERE id = $1;
```

---

## Relationships

### Entity Relationship Diagram

```
┌─────────────┐
│    users    │ (auth_db)
│  - id (PK)  │
│  - email    │
│  - name     │
│  - password │
│  - role     │
└─────────────┘
       │
       │ (referenced by)
       ▼
┌──────────────┐
│   bookings   │ (booking_db)
│  - id (PK)   │
│  - user_id   │───┐
│  - bed_id    │   │
│  - status    │   │
└──────────────┘   │
       │           │
       │ (updates) │
       ▼           │
┌──────────────┐   │
│   buildings  │   │
│  - id (PK)   │   │
└──────┬───────┘   │
       │           │
       │ 1:N       │
       ▼           │
┌──────────────┐   │
│    rooms     │   │
│  - id (PK)   │   │
│  - bldg_id   │   │
└──────┬───────┘   │
       │           │
       │ 1:N       │
       ▼           │
┌──────────────┐   │
│     beds     │◄──┘
│  - id (PK)   │
│  - room_id   │
│  - occupied  │
└──────────────┘
```

### Cross-Database References

Since each service has its own database, references are maintained through:

1. **Application-Level Foreign Keys**: IDs stored but not enforced by DB
2. **Service-to-Service API Calls**: Validation via HTTP requests
3. **Eventual Consistency**: Updates propagated through service communication

Example: When a booking is created:
1. Booking service validates user via Auth service API
2. Booking service checks bed availability via Building service API
3. Booking service creates booking record
4. Booking service updates bed status via Building service API

---

## Database Initialization Scripts

### Initialize All Databases

Create this file as `init-databases.sql`:

```sql
-- Create databases
CREATE DATABASE hostel_auth_db;
CREATE DATABASE hostel_building_db;
CREATE DATABASE hostel_booking_db;

-- Connect to auth database
\c hostel_auth_db;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'student',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Connect to building database
\c hostel_building_db;

-- Create buildings table
CREATE TABLE IF NOT EXISTS buildings (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    total_rooms INTEGER NOT NULL DEFAULT 0,
    total_beds INTEGER NOT NULL DEFAULT 0,
    available_beds INTEGER NOT NULL DEFAULT 0,
    amenities JSONB,
    image TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create rooms table
CREATE TABLE IF NOT EXISTS rooms (
    id VARCHAR(255) PRIMARY KEY,
    building_id VARCHAR(255) NOT NULL REFERENCES buildings(id),
    number VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    total_beds INTEGER NOT NULL,
    available_beds INTEGER NOT NULL DEFAULT 0,
    amenities JSONB,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(building_id, number)
);

-- Create beds table
CREATE TABLE IF NOT EXISTS beds (
    id VARCHAR(255) PRIMARY KEY,
    room_id VARCHAR(255) NOT NULL REFERENCES rooms(id),
    number INTEGER NOT NULL,
    is_occupied BOOLEAN DEFAULT FALSE,
    occupied_by VARCHAR(255),
    occupied_by_name VARCHAR(255),
    UNIQUE(room_id, number)
);

-- Create indexes
CREATE INDEX idx_buildings_name ON buildings(name);
CREATE INDEX idx_rooms_building ON rooms(building_id);
CREATE INDEX idx_beds_room ON beds(room_id);
CREATE INDEX idx_beds_occupied ON beds(is_occupied);

-- Connect to booking database
\c hostel_booking_db;

-- Create bookings table
CREATE TABLE IF NOT EXISTS bookings (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    building_id VARCHAR(255) NOT NULL,
    building_name VARCHAR(255) NOT NULL,
    room_id VARCHAR(255) NOT NULL,
    room_number VARCHAR(50) NOT NULL,
    bed_id VARCHAR(255) NOT NULL,
    bed_number INTEGER NOT NULL,
    booking_date TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_bookings_user ON bookings(user_id);
CREATE INDEX idx_bookings_building ON bookings(building_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_date ON bookings(booking_date);
```

Run with:
```powershell
psql -U postgres -f init-databases.sql
```

---

## Backup and Restore

### Backup All Databases

```powershell
# Backup auth database
pg_dump -U postgres -d hostel_auth_db -F c -f auth_backup.dump

# Backup building database
pg_dump -U postgres -d hostel_building_db -F c -f building_backup.dump

# Backup booking database
pg_dump -U postgres -d hostel_booking_db -F c -f booking_backup.dump
```

### Restore Databases

```powershell
# Restore auth database
pg_restore -U postgres -d hostel_auth_db -c auth_backup.dump

# Restore building database
pg_restore -U postgres -d hostel_building_db -c building_backup.dump

# Restore booking database
pg_restore -U postgres -d hostel_booking_db -c booking_backup.dump
```

---

## Performance Optimization

### Recommended Indexes

Already included in schema, but for reference:

```sql
-- Auth DB
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Building DB
CREATE INDEX idx_buildings_available_beds ON buildings(available_beds);
CREATE INDEX idx_rooms_type ON rooms(type);
CREATE INDEX idx_rooms_available_beds ON rooms(available_beds);
CREATE INDEX idx_beds_occupied_by ON beds(occupied_by);

-- Booking DB
CREATE INDEX idx_bookings_user_status ON bookings(user_id, status);
CREATE INDEX idx_bookings_bed_status ON bookings(bed_id, status);
```

### Query Optimization Tips

1. **Use prepared statements** for repeated queries
2. **Limit result sets** with LIMIT and OFFSET for pagination
3. **Use EXPLAIN ANALYZE** to check query performance
4. **Create composite indexes** for frequently combined WHERE clauses
5. **Use JSONB operators** efficiently for amenity searches

---

## Maintenance

### Vacuum and Analyze

```sql
-- Run weekly
VACUUM ANALYZE users;
VACUUM ANALYZE buildings;
VACUUM ANALYZE rooms;
VACUUM ANALYZE beds;
VACUUM ANALYZE bookings;
```

### Check Database Size

```sql
SELECT 
    datname, 
    pg_size_pretty(pg_database_size(datname)) as size
FROM pg_database
WHERE datname IN ('hostel_auth_db', 'hostel_building_db', 'hostel_booking_db');
```

---

**Database schema complete and optimized for microservices architecture! 🗄️**
