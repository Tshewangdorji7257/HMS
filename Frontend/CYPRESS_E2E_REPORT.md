# Cypress End-to-End Testing Report
## Hostel Management System - Frontend & Backend Integration

**Report Date:** December 5, 2025  
**Testing Type:** End-to-End Testing  
**Framework:** Cypress  
**Status:** ⏳ PLANNED

---

## Executive Summary

⚠️ **Status:** PLANNED - Tests Not Yet Implemented

End-to-end testing will validate complete user workflows from frontend UI through backend APIs, ensuring the entire HMS application functions correctly from a user perspective.

---

## Test Plan

### Objectives
- Test complete user journeys
- Validate UI/API integration
- Test cross-browser compatibility
- Verify data persistence
- Test error handling flows

### Technology Stack
- **Cypress** v13.6.0 (E2E testing)
- **TypeScript** (Test scripts)
- **Mochawesome** (Reporting)

### Test Environment
```
Frontend: http://localhost:3000
Backend: http://localhost:8000
Database: PostgreSQL (test instance)
```

---

## Planned Test Suites

### 1. Authentication & Authorization

#### User Registration Flow
```typescript
describe('User Registration', () => {
  it('should allow new user signup', () => {
    cy.visit('/auth')
    cy.get('[data-cy=signup-tab]').click()
    cy.get('[data-cy=email-input]').type('newuser@test.com')
    cy.get('[data-cy=password-input]').type('SecurePass123!')
    cy.get('[data-cy=name-input]').type('Test User')
    cy.get('[data-cy=signup-button]').click()
    
    cy.url().should('include', '/dashboard')
    cy.contains('Welcome, Test User').should('be.visible')
  })
  
  it('should reject duplicate email', () => {
    // Test email uniqueness constraint
  })
  
  it('should validate password strength', () => {
    // Test password validation
  })
})
```

#### User Login Flow
```typescript
describe('User Login', () => {
  it('should login with valid credentials', () => {
    cy.visit('/auth')
    cy.get('[data-cy=email-input]').type('test@example.com')
    cy.get('[data-cy=password-input]').type('password123')
    cy.get('[data-cy=login-button]').click()
    
    cy.url().should('include', '/dashboard')
    cy.getCookie('token').should('exist')
  })
  
  it('should reject invalid credentials', () => {
    cy.visit('/auth')
    cy.get('[data-cy=email-input]').type('wrong@example.com')
    cy.get('[data-cy=password-input]').type('wrongpass')
    cy.get('[data-cy=login-button]').click()
    
    cy.contains('Invalid credentials').should('be.visible')
    cy.url().should('include', '/auth')
  })
  
  it('should redirect unauthenticated users', () => {
    cy.visit('/dashboard')
    cy.url().should('include', '/auth')
  })
})
```

#### Admin Access Control
```typescript
describe('Admin Authorization', () => {
  it('should allow admin access to admin panel', () => {
    cy.loginAsAdmin()
    cy.visit('/admin')
    cy.contains('Admin Dashboard').should('be.visible')
  })
  
  it('should deny regular user access to admin', () => {
    cy.loginAsUser()
    cy.visit('/admin')
    cy.contains('Access Denied').should('be.visible')
  })
})
```

---

### 2. Building Search & Discovery

#### Building List View
```typescript
describe('Building Listings', () => {
  beforeEach(() => {
    cy.loginAsUser()
  })
  
  it('should display all available buildings', () => {
    cy.visit('/dashboard')
    cy.get('[data-cy=building-card]').should('have.length.greaterThan', 0)
  })
  
  it('should show building details', () => {
    cy.get('[data-cy=building-card]').first().within(() => {
      cy.get('[data-cy=building-name]').should('be.visible')
      cy.get('[data-cy=building-location]').should('be.visible')
      cy.get('[data-cy=available-rooms]').should('be.visible')
    })
  })
})
```

#### Building Search
```typescript
describe('Building Search', () => {
  it('should filter by location', () => {
    cy.visit('/dashboard')
    cy.get('[data-cy=location-filter]').type('North Campus')
    cy.get('[data-cy=search-button]').click()
    
    cy.get('[data-cy=building-card]').each(($card) => {
      cy.wrap($card).should('contain', 'North Campus')
    })
  })
  
  it('should filter by gender', () => {
    cy.get('[data-cy=gender-filter]').select('male')
    cy.get('[data-cy=building-card]').each(($card) => {
      cy.wrap($card).should('contain', 'Male')
    })
  })
  
  it('should filter by price range', () => {
    cy.get('[data-cy=min-price]').type('5000')
    cy.get('[data-cy=max-price]').type('10000')
    cy.get('[data-cy=search-button]').click()
    
    // Verify results within price range
  })
})
```

#### Building Details Page
```typescript
describe('Building Details', () => {
  it('should navigate to building details', () => {
    cy.visit('/dashboard')
    cy.get('[data-cy=building-card]').first().click()
    
    cy.url().should('match', /\/building\/[\w-]+/)
    cy.get('[data-cy=building-details]').should('be.visible')
  })
  
  it('should display room grid', () => {
    cy.visit('/building/bldg-1')
    cy.get('[data-cy=room-card]').should('have.length.greaterThan', 0)
  })
  
  it('should show bed availability', () => {
    cy.get('[data-cy=room-card]').first().within(() => {
      cy.get('[data-cy=bed-grid]').should('be.visible')
      cy.get('[data-cy=vacant-bed]').should('exist')
    })
  })
})
```

---

### 3. Booking Creation Flow

#### Complete Booking Process
```typescript
describe('Booking Creation', () => {
  beforeEach(() => {
    cy.loginAsUser()
  })
  
  it('should create a booking successfully', () => {
    // Navigate to building
    cy.visit('/building/bldg-1')
    
    // Select room
    cy.get('[data-cy=room-card]').first().click()
    
    // Select bed
    cy.get('[data-cy=vacant-bed]').first().click()
    
    // Fill booking form
    cy.get('[data-cy=check-in-date]').type('2025-06-01')
    cy.get('[data-cy=check-out-date]').type('2025-12-31')
    
    // Submit booking
    cy.get('[data-cy=confirm-booking]').click()
    
    // Verify confirmation
    cy.contains('Booking Confirmed').should('be.visible')
    cy.get('[data-cy=booking-id]').should('be.visible')
    
    // Verify redirection
    cy.url().should('include', '/bookings')
  })
  
  it('should validate booking dates', () => {
    cy.visit('/building/bldg-1')
    cy.get('[data-cy=room-card]').first().click()
    cy.get('[data-cy=vacant-bed]').first().click()
    
    // Past date should be rejected
    cy.get('[data-cy=check-in-date]').type('2024-01-01')
    cy.get('[data-cy=confirm-booking]').click()
    
    cy.contains('Check-in date cannot be in the past').should('be.visible')
  })
  
  it('should prevent double booking', () => {
    // Book a bed
    cy.createBooking('bed-1', '2025-06-01', '2025-12-31')
    
    // Try to book same bed
    cy.visit('/building/bldg-1')
    cy.get('[data-cy=occupied-bed]').first().should('have.class', 'disabled')
    cy.get('[data-cy=occupied-bed]').first().click()
    
    cy.contains('This bed is already occupied').should('be.visible')
  })
})
```

#### Booking Validation
```typescript
describe('Booking Validation', () => {
  it('should require check-in date', () => {
    cy.visit('/building/bldg-1')
    cy.get('[data-cy=room-card]').first().click()
    cy.get('[data-cy=vacant-bed]').first().click()
    cy.get('[data-cy=check-out-date]').type('2025-12-31')
    cy.get('[data-cy=confirm-booking]').click()
    
    cy.contains('Check-in date is required').should('be.visible')
  })
  
  it('should validate date range', () => {
    // Check-out before check-in should fail
    cy.get('[data-cy=check-in-date]').type('2025-12-31')
    cy.get('[data-cy=check-out-date]').type('2025-06-01')
    cy.get('[data-cy=confirm-booking]').click()
    
    cy.contains('Check-out must be after check-in').should('be.visible')
  })
})
```

---

### 4. My Bookings Management

#### View Bookings
```typescript
describe('My Bookings', () => {
  beforeEach(() => {
    cy.loginAsUser()
    cy.visit('/bookings')
  })
  
  it('should display user bookings', () => {
    cy.get('[data-cy=booking-card]').should('have.length.greaterThan', 0)
  })
  
  it('should show booking details', () => {
    cy.get('[data-cy=booking-card]').first().within(() => {
      cy.get('[data-cy=building-name]').should('be.visible')
      cy.get('[data-cy=room-number]').should('be.visible')
      cy.get('[data-cy=bed-number]').should('be.visible')
      cy.get('[data-cy=check-in-date]').should('be.visible')
      cy.get('[data-cy=booking-status]').should('be.visible')
    })
  })
  
  it('should filter by status', () => {
    cy.get('[data-cy=status-filter]').select('confirmed')
    cy.get('[data-cy=booking-card]').each(($card) => {
      cy.wrap($card).find('[data-cy=booking-status]').should('contain', 'Confirmed')
    })
  })
})
```

#### Cancel Booking
```typescript
describe('Booking Cancellation', () => {
  it('should cancel a booking', () => {
    cy.visit('/bookings')
    cy.get('[data-cy=booking-card]').first().within(() => {
      cy.get('[data-cy=cancel-button]').click()
    })
    
    // Confirm cancellation
    cy.get('[data-cy=confirm-cancel]').click()
    
    // Verify status changed
    cy.contains('Booking Cancelled').should('be.visible')
    cy.get('[data-cy=booking-status]').should('contain', 'Cancelled')
  })
  
  it('should show cancellation confirmation', () => {
    cy.get('[data-cy=cancel-button]').click()
    cy.contains('Are you sure you want to cancel?').should('be.visible')
    cy.get('[data-cy=cancel-modal]').should('be.visible')
  })
  
  it('should allow cancellation dismissal', () => {
    cy.get('[data-cy=cancel-button]').click()
    cy.get('[data-cy=dismiss-cancel]').click()
    cy.get('[data-cy=cancel-modal]').should('not.exist')
  })
})
```

---

### 5. Admin Dashboard

#### Admin Building Management
```typescript
describe('Admin Building Management', () => {
  beforeEach(() => {
    cy.loginAsAdmin()
    cy.visit('/admin')
  })
  
  it('should display all bookings', () => {
    cy.get('[data-cy=all-bookings]').should('be.visible')
    cy.get('[data-cy=booking-row]').should('have.length.greaterThan', 0)
  })
  
  it('should allow booking status update', () => {
    cy.get('[data-cy=booking-row]').first().within(() => {
      cy.get('[data-cy=status-select]').select('confirmed')
      cy.get('[data-cy=update-status]').click()
    })
    
    cy.contains('Status updated').should('be.visible')
  })
  
  it('should show statistics', () => {
    cy.get('[data-cy=total-bookings]').should('be.visible')
    cy.get('[data-cy=total-revenue]').should('be.visible')
    cy.get('[data-cy=occupancy-rate]').should('be.visible')
  })
})
```

---

## Custom Cypress Commands

### Authentication Commands
```typescript
// cypress/support/commands.ts

Cypress.Commands.add('loginAsUser', () => {
  cy.session('user-session', () => {
    cy.request({
      method: 'POST',
      url: 'http://localhost:8000/api/auth/login',
      body: {
        email: 'test@example.com',
        password: 'password123'
      }
    }).then((response) => {
      window.localStorage.setItem('token', response.body.token)
    })
  })
})

Cypress.Commands.add('loginAsAdmin', () => {
  cy.session('admin-session', () => {
    cy.request({
      method: 'POST',
      url: 'http://localhost:8000/api/auth/admin/login',
      body: {
        email: 'admin@example.com',
        password: 'admin123'
      }
    }).then((response) => {
      window.localStorage.setItem('token', response.body.token)
    })
  })
})
```

### Booking Commands
```typescript
Cypress.Commands.add('createBooking', (bedId, checkIn, checkOut) => {
  cy.request({
    method: 'POST',
    url: 'http://localhost:8000/api/bookings',
    headers: {
      'Authorization': `Bearer ${window.localStorage.getItem('token')}`
    },
    body: {
      bed_id: bedId,
      check_in_date: checkIn,
      check_out_date: checkOut
    }
  })
})
```

---

## Test Data Strategy

### Fixtures
```typescript
// cypress/fixtures/users.json
{
  "validUser": {
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  },
  "admin": {
    "email": "admin@example.com",
    "password": "admin123"
  }
}

// cypress/fixtures/buildings.json
{
  "sampleBuilding": {
    "name": "North Campus Building A",
    "location": "North Campus",
    "gender": "male",
    "total_rooms": 50,
    "price_per_semester": 8000
  }
}
```

### Database Seeding
```bash
# Before tests
npm run db:seed:test

# After tests
npm run db:clean:test
```

---

## CI/CD Integration

### GitHub Actions Workflow
```yaml
name: Cypress E2E Tests

on: [push, pull_request]

jobs:
  cypress:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Start services
        run: docker-compose up -d
      
      - name: Wait for services
        run: npm run wait-on:services
      
      - name: Run Cypress tests
        uses: cypress-io/github-action@v6
        with:
          browser: chrome
          record: true
          parallel: true
      
      - name: Upload screenshots
        if: failure()
        uses: actions/upload-artifact@v3
        with:
          name: cypress-screenshots
          path: cypress/screenshots
```

---

## Success Criteria

### Functional Requirements
✅ All critical user journeys pass
✅ Authentication flows work correctly
✅ Booking creation succeeds
✅ Data persists correctly
✅ Error messages display properly

### Non-Functional Requirements
✅ Tests run in < 10 minutes
✅ Zero flaky tests
✅ Screenshots on failure
✅ Video recordings available

---

## Timeline

| Phase | Tasks | Duration | Status |
|-------|-------|----------|--------|
| **Setup** | Install Cypress, configure | 2 days | ⏳ PENDING |
| **Auth Tests** | Login, signup, authorization | 3 days | ⏳ PENDING |
| **Search Tests** | Building search & filters | 2 days | ⏳ PENDING |
| **Booking Tests** | Complete booking flow | 4 days | ⏳ PENDING |
| **Admin Tests** | Admin dashboard tests | 2 days | ⏳ PENDING |
| **CI/CD** | GitHub Actions integration | 2 days | ⏳ PENDING |

**Total:** 3 weeks

---

## Conclusion

Comprehensive E2E testing suite is planned to validate all critical user workflows. Tests will cover authentication, building search, booking creation, and admin management with focus on real user interactions.

**Next Steps:**
1. Set up Cypress in Frontend project
2. Create test fixtures and commands
3. Implement authentication test suite
4. Add booking flow tests
5. Integrate with CI/CD pipeline

**Estimated Completion:** 3 weeks

---

**Report Prepared By:** GitHub Copilot  
**Framework:** Cypress v13.6.0  
**Status:** PLANNED - Ready for Implementation
