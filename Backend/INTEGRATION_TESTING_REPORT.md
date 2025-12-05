# Integration Testing Report
## Hostel Management System - Backend Services

**Report Date:** December 5, 2025  
**Project:** HMS Backend Microservices  
**Testing Type:** Integration Testing  
**Report Version:** 1.0

---

## Executive Summary

Integration testing validates the interaction between different components of the HMS backend microservices, including database operations, service-to-service communication, and external dependencies.

### Current Status

⚠️ **Status:** PLANNED - Not Yet Implemented

| Test Category | Planned | Implemented | Status |
|--------------|---------|-------------|--------|
| Database Integration | Yes | ❌ No | PENDING |
| Service-to-Service | Yes | ❌ No | PENDING |
| Consul Integration | Yes | ❌ No | PENDING |
| Email Service (SMTP) | Yes | ❌ No | PENDING |
| API Gateway Routing | Yes | ❌ No | PENDING |

---

## Test Plan

### 1. Objectives

- Validate database CRUD operations work correctly
- Test transaction handling and rollback scenarios
- Verify service-to-service communication via API Gateway
- Test Consul service discovery integration
- Validate email notification sending
- Ensure data consistency across services

### 2. Test Environment

#### Infrastructure Requirements:
```yaml
Services:
  - PostgreSQL 14+ (Test Database)
  - Consul 1.15+ (Service Discovery)
  - MailHog/SMTP Server (Email Testing)
  - Docker Compose (Container Orchestration)

Test Data:
  - Seed data for each service
  - Test fixtures for common scenarios
  - Mock external APIs where needed
```

#### Environment Configuration:
```bash
# Test Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=test_user
DB_PASSWORD=test_password
DB_NAME=hms_test

# Consul
CONSUL_ADDR=localhost:8500

# SMTP
SMTP_HOST=localhost
SMTP_PORT=1025
```

### 3. Test Scope

#### In Scope:
- ✅ Database connection and query execution
- ✅ CRUD operations for all entities
- ✅ Transaction management
- ✅ Foreign key constraints
- ✅ API Gateway proxying to services
- ✅ Consul service registration/discovery
- ✅ Email notification delivery
- ✅ Cross-service data consistency

#### Out of Scope:
- ❌ End-to-end user workflows (covered in E2E tests)
- ❌ Performance under load (covered in performance tests)
- ❌ Security penetration testing (covered in DAST)
- ❌ UI/Frontend integration (covered in Cypress)

---

## Planned Test Cases

### Auth Service Integration Tests

#### Database Tests:
```go
// Test File: auth-service/tests/integration/database_test.go

TestUserSignup_DatabaseIntegration
  ✓ Creates user record in database
  ✓ Password is hashed before storage
  ✓ Email uniqueness constraint enforced
  ✓ Returns user ID on success

TestUserLogin_DatabaseIntegration
  ✓ Retrieves user by email
  ✓ Validates password against hash
  ✓ Returns JWT token on success
  ✓ Returns error for invalid credentials

TestAdminSignup_DatabaseIntegration
  ✓ Creates admin user with role='admin'
  ✓ Admin-specific fields populated
  ✓ Transaction rolls back on validation error

TestGetUserById_DatabaseIntegration
  ✓ Retrieves user by ID
  ✓ Returns 404 for non-existent user
  ✓ Excludes password from response
```

#### Consul Integration:
```go
TestServiceRegistration_Consul
  ✓ Registers auth-service with Consul
  ✓ Health check endpoint accessible
  ✓ Service deregisters on shutdown

TestServiceDiscovery_Consul
  ✓ Other services can discover auth-service
  ✓ Load balancing works with multiple instances
```

---

### Building Service Integration Tests

#### Database Tests:
```go
// Test File: building-service/tests/integration/database_test.go

TestCreateBuilding_DatabaseIntegration
  ✓ Creates building with rooms and beds
  ✓ Nested transaction for related entities
  ✓ Rollback on partial failure

TestSearchBuildings_DatabaseIntegration
  ✓ Full-text search works correctly
  ✓ Filters by location, gender, amenities
  ✓ Pagination returns correct results

TestGetAvailableBeds_DatabaseIntegration
  ✓ Returns only vacant beds
  ✓ Filters by room type
  ✓ Respects booking date ranges

TestUpdateBedOccupancy_DatabaseIntegration
  ✓ Changes bed status to occupied
  ✓ Links bed to user_id
  ✓ Prevents double occupancy
```

#### Cross-Service Tests:
```go
TestBuildingBookingIntegration
  ✓ Booking service can query available beds
  ✓ Building service validates bed availability
  ✓ Occupancy updates when booking created
```

---

### Booking Service Integration Tests

#### Database Tests:
```go
// Test File: booking-service/tests/integration/database_test.go

TestCreateBooking_DatabaseIntegration
  ✓ Creates booking record
  ✓ Updates bed occupancy in building-service
  ✓ Transaction commits only if both succeed
  ✓ Generates booking confirmation ID

TestCancelBooking_DatabaseIntegration
  ✓ Changes booking status to 'cancelled'
  ✓ Releases bed in building-service
  ✓ Sends cancellation email
  ✓ Refund logic triggers correctly

TestGetUserBookings_DatabaseIntegration
  ✓ Retrieves all bookings for user
  ✓ Includes building and room details
  ✓ Orders by check-in date
```

#### Email Integration:
```go
TestBookingConfirmationEmail_SMTPIntegration
  ✓ Email sent on booking creation
  ✓ Contains correct booking details
  ✓ Uses proper template formatting
  ✓ Attachments included (if any)

TestCancellationEmail_SMTPIntegration
  ✓ Email sent on booking cancellation
  ✓ Contains cancellation reason
  ✓ Includes refund information
```

---

### API Gateway Integration Tests

#### Routing Tests:
```go
// Test File: api-gateway/tests/integration/routing_test.go

TestAuthServiceRouting_Integration
  ✓ POST /api/auth/signup routes to auth-service
  ✓ POST /api/auth/login routes to auth-service
  ✓ Authorization header forwarded
  ✓ Response returned to client

TestBuildingServiceRouting_Integration
  ✓ GET /api/buildings routes correctly
  ✓ Search queries preserved
  ✓ JSON response format maintained

TestBookingServiceRouting_Integration
  ✓ POST /api/bookings routes correctly
  ✓ Request body forwarded intact
  ✓ User authentication validated
```

#### Error Handling:
```go
TestServiceUnavailable_Integration
  ✓ Returns 502 when service down
  ✓ Error message indicates which service
  ✓ Client receives proper JSON error

TestTimeout_Integration
  ✓ Returns 504 on service timeout
  ✓ Timeout configurable per service
  ✓ Request cancelled properly
```

---

## Implementation Strategy

### Phase 1: Database Integration (Week 1)

**Setup:**
```bash
# Create Docker Compose for test database
docker-compose -f docker-compose.test.yml up -d postgres

# Run database migrations
cd auth-service
go run migrations/migrate.go up

# Seed test data
go run tests/fixtures/seed.go
```

**Test Structure:**
```
Backend/
├── auth-service/
│   └── tests/
│       ├── integration/
│       │   ├── database_test.go
│       │   ├── setup_test.go
│       │   └── teardown_test.go
│       └── fixtures/
│           ├── users.sql
│           └── seed.go
```

**Execution:**
```bash
# Run integration tests
go test ./tests/integration -tags=integration -v

# With coverage
go test ./tests/integration -tags=integration -coverprofile=integration_coverage.out
```

---

### Phase 2: Service-to-Service Integration (Week 2)

**Setup:**
```bash
# Start all services
docker-compose up -d

# Wait for services to be healthy
./scripts/wait-for-services.sh

# Run integration tests
go test ./tests/e2e -tags=integration
```

**Test Scenarios:**
```go
Scenario: Complete Booking Flow
  1. User signs up via auth-service ✓
  2. User searches buildings via building-service ✓
  3. User creates booking via booking-service ✓
  4. Booking-service updates bed occupancy ✓
  5. Email notification sent ✓
  6. All data consistent across services ✓
```

---

### Phase 3: External Service Integration (Week 3)

**Consul Integration:**
```bash
# Start Consul
docker run -d -p 8500:8500 consul agent -dev

# Verify service registration
curl http://localhost:8500/v1/catalog/services
```

**Email Integration:**
```bash
# Start MailHog for email testing
docker run -d -p 1025:1025 -p 8025:8025 mailhog/mailhog

# Access email UI
http://localhost:8025
```

---

## Test Data Strategy

### Database Fixtures

**Users:**
```sql
-- tests/fixtures/users.sql
INSERT INTO users (id, email, password, name, role, created_at)
VALUES 
  ('user-1', 'test@example.com', '$2a$10$hashed_password', 'Test User', 'user', NOW()),
  ('admin-1', 'admin@example.com', '$2a$10$hashed_password', 'Admin User', 'admin', NOW());
```

**Buildings:**
```sql
-- tests/fixtures/buildings.sql
INSERT INTO buildings (id, name, location, gender, total_rooms, available_rooms)
VALUES 
  ('bldg-1', 'Test Building A', 'Campus North', 'male', 50, 45);

INSERT INTO rooms (id, building_id, room_number, room_type, total_beds)
VALUES 
  ('room-1', 'bldg-1', '101', 'double', 2);

INSERT INTO beds (id, room_id, bed_number, is_occupied)
VALUES 
  ('bed-1', 'room-1', '1', false),
  ('bed-2', 'room-1', '2', false);
```

**Bookings:**
```sql
-- tests/fixtures/bookings.sql
INSERT INTO bookings (id, user_id, bed_id, check_in_date, check_out_date, status)
VALUES 
  ('book-1', 'user-1', 'bed-1', '2025-01-01', '2025-06-30', 'confirmed');
```

---

## Success Criteria

### Database Integration
- ✅ All CRUD operations work correctly
- ✅ Transactions commit/rollback properly
- ✅ Foreign key constraints enforced
- ✅ Query performance acceptable (<100ms)
- ✅ Connection pooling works under load

### Service Integration
- ✅ API Gateway routes to correct services
- ✅ Authentication propagates across services
- ✅ Error responses formatted consistently
- ✅ Cross-service data updates work atomically

### External Services
- ✅ Consul registers all services
- ✅ Health checks pass consistently
- ✅ Emails delivered within 5 seconds
- ✅ Email templates render correctly

---

## Current Gaps

### Why Integration Tests Not Implemented Yet

1. **Database Schema Evolution**
   - Schema changes frequently during development
   - Need stable schema before integration tests
   - Migrations not finalized

2. **Service Dependencies**
   - Some services not fully implemented
   - Need all services running for cross-service tests
   - Docker Compose setup needs refinement

3. **Test Infrastructure**
   - Test database not containerized
   - No CI/CD pipeline for integration tests
   - Lack of test data management strategy

---

## Recommendations

### Immediate Next Steps

1. **Create Test Docker Compose**
   ```yaml
   # docker-compose.test.yml
   version: '3.8'
   services:
     postgres-test:
       image: postgres:14
       environment:
         POSTGRES_DB: hms_test
         POSTGRES_USER: test_user
         POSTGRES_PASSWORD: test_password
       ports:
         - "5433:5432"
     
     consul-test:
       image: consul:latest
       ports:
         - "8500:8500"
     
     mailhog:
       image: mailhog/mailhog
       ports:
         - "1025:1025"
         - "8025:8025"
   ```

2. **Implement Test Helpers**
   ```go
   // tests/helpers/database.go
   func SetupTestDB() (*sql.DB, func()) {
       db := connectToTestDB()
       cleanup := func() {
           clearTestData(db)
           db.Close()
       }
       return db, cleanup
   }
   ```

3. **Create Integration Test Suite**
   - Start with auth-service database tests
   - Add building-service CRUD tests
   - Implement booking flow integration test
   - Add email notification tests

### Long-term Improvements

4. **Test Data Management**
   - Implement test data builder pattern
   - Create reusable fixtures
   - Add database seeding scripts

5. **CI/CD Integration**
   - Add integration tests to GitHub Actions
   - Run tests on every PR
   - Generate coverage reports

6. **Performance Baselines**
   - Establish acceptable response times
   - Monitor database query performance
   - Detect performance regressions

---

## Timeline

| Phase | Tasks | Duration | Status |
|-------|-------|----------|--------|
| **Phase 1** | Database integration tests | 1 week | ⏳ PENDING |
| **Phase 2** | Service-to-service tests | 1 week | ⏳ PENDING |
| **Phase 3** | External service tests | 1 week | ⏳ PENDING |
| **Phase 4** | CI/CD integration | 3 days | ⏳ PENDING |
| **Phase 5** | Documentation | 2 days | ⏳ PENDING |

**Total Estimated Time:** 3.5 weeks

---

## Conclusion

Integration testing is critical for validating the HMS backend microservices work correctly together. While unit tests verify individual components, integration tests ensure database operations, service communication, and external dependencies function as expected in a realistic environment.

The test plan is comprehensive and ready for implementation. Priority should be given to database integration tests first, followed by service-to-service communication tests, and finally external service integration.

**Current Status:** ⏳ **PLANNED** - Ready for implementation  
**Next Action:** Create test database Docker Compose and implement Phase 1

---

## Appendix A: Test Environment Setup

### Quick Start

```bash
# 1. Start test infrastructure
docker-compose -f docker-compose.test.yml up -d

# 2. Run database migrations
make migrate-test

# 3. Seed test data
make seed-test

# 4. Run integration tests
make test-integration

# 5. View coverage report
make coverage-integration

# 6. Cleanup
docker-compose -f docker-compose.test.yml down -v
```

### Environment Variables

```bash
# .env.test
NODE_ENV=test
GO_ENV=test

# Database
DB_HOST=localhost
DB_PORT=5433
DB_NAME=hms_test
DB_USER=test_user
DB_PASSWORD=test_password

# Services
AUTH_SERVICE_URL=http://localhost:8001
BUILDING_SERVICE_URL=http://localhost:8002
BOOKING_SERVICE_URL=http://localhost:8003
API_GATEWAY_URL=http://localhost:8000

# External Services
CONSUL_ADDR=localhost:8500
SMTP_HOST=localhost
SMTP_PORT=1025
```

---

**Report Prepared By:** GitHub Copilot  
**Review Date:** December 5, 2025  
**Next Review:** After Phase 1 implementation
