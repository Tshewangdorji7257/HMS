# Frontend Integration Guide

Complete guide for integrating the Next.js frontend with the Golang microservices backend.

## 📋 Overview

This guide will help you connect your existing Next.js/TypeScript Hostel Management System frontend to the newly deployed backend services.

**Frontend Repository**: https://github.com/Tshewangdorji7257/Hostel-Management-System

**Backend Base URL**: `http://localhost:8000`

## ⚠️ Important: CORS Configuration

**CORS has been properly configured** to prevent duplicate header errors:
- ✅ Only the API Gateway handles CORS headers
- ✅ Backend services do not add CORS headers
- ✅ Allowed origins: `http://localhost:3000` and `http://localhost:3001`
- ✅ Credentials are enabled for authentication

If you encounter CORS errors, ensure:
1. Your frontend is running on `localhost:3000` or `localhost:3001`
2. All services are running: `docker-compose ps`
3. Services were rebuilt after CORS fix: `docker-compose up -d --build`

---

## 🚀 Quick Start

### 1. Update Environment Variables

Create or update `.env.local` in your frontend project root:

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8000
NEXT_PUBLIC_API_TIMEOUT=30000
```

### 2. Create API Configuration File

Create `lib/api-config.ts`:

```typescript
export const API_CONFIG = {
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8000',
  timeout: Number(process.env.NEXT_PUBLIC_API_TIMEOUT) || 30000,
  headers: {
    'Content-Type': 'application/json',
  },
};

export const getAuthHeader = (token: string) => ({
  Authorization: `Bearer ${token}`,
});
```

---

## 🔐 Authentication Integration

### Update `lib/auth.ts`

Replace the localStorage-based authentication with backend API calls:

```typescript
import { API_CONFIG, getAuthHeader } from './api-config';

export interface User {
  id: string;
  email: string;
  name: string;
  role: 'student' | 'admin';
  created_at?: string;
}

export interface AuthResponse {
  success: boolean;
  message: string;
  token: string;
  user: User;
}

export const authService = {
  // Sign up new user
  async signup(
    email: string,
    password: string,
    name: string,
    role: 'student' | 'admin' = 'student'
  ): Promise<AuthResponse> {
    const response = await fetch(`${API_CONFIG.baseURL}/api/auth/signup`, {
      method: 'POST',
      headers: API_CONFIG.headers,
      body: JSON.stringify({ email, password, name, role }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Signup failed');
    }

    const data = await response.json();
    
    // Store token in localStorage
    if (data.token) {
      localStorage.setItem('token', data.token);
      localStorage.setItem('user', JSON.stringify(data.user));
    }

    return data;
  },

  // Login existing user
  async login(email: string, password: string): Promise<AuthResponse> {
    const response = await fetch(`${API_CONFIG.baseURL}/api/auth/login`, {
      method: 'POST',
      headers: API_CONFIG.headers,
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Login failed');
    }

    const data = await response.json();
    
    // Store token in localStorage
    if (data.token) {
      localStorage.setItem('token', data.token);
      localStorage.setItem('user', JSON.stringify(data.user));
    }

    return data;
  },

  // Logout user
  logout(): void {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  // Get current user from localStorage
  getCurrentUser(): User | null {
    const userJson = localStorage.getItem('user');
    return userJson ? JSON.parse(userJson) : null;
  },

  // Get stored token
  getToken(): string | null {
    return localStorage.getItem('token');
  },

  // Validate token with backend
  async validateToken(token: string): Promise<boolean> {
    try {
      const response = await fetch(`${API_CONFIG.baseURL}/api/auth/validate`, {
        method: 'POST',
        headers: {
          ...API_CONFIG.headers,
          ...getAuthHeader(token),
        },
      });

      return response.ok;
    } catch (error) {
      return false;
    }
  },

  // Get user profile from backend
  async getProfile(token: string): Promise<User> {
    const response = await fetch(`${API_CONFIG.baseURL}/api/auth/profile`, {
      method: 'GET',
      headers: {
        ...API_CONFIG.headers,
        ...getAuthHeader(token),
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch profile');
    }

    const data = await response.json();
    return data.user;
  },

  // Check if user is authenticated
  isAuthenticated(): boolean {
    const token = this.getToken();
    const user = this.getCurrentUser();
    return !!(token && user);
  },

  // Check if user is admin
  isAdmin(): boolean {
    const user = this.getCurrentUser();
    return user?.role === 'admin';
  },
};
```

---

## 🏢 Building/Room Data Integration

### Update `lib/data.ts`

Replace the mock data service with backend API calls:

```typescript
import { API_CONFIG, getAuthHeader } from './api-config';
import { authService } from './auth';

export interface Bed {
  id: string;
  number: number;
  is_occupied: boolean;
  occupied_by?: string;
  occupied_by_name?: string;
}

export interface Room {
  id: string;
  number: string;
  type: 'single' | 'double' | 'triple' | 'quad';
  total_beds: number;
  available_beds: number;
  amenities: string[];
  price: number;
  beds: Bed[];
}

export interface Building {
  id: string;
  name: string;
  description: string;
  total_rooms: number;
  total_beds: number;
  available_beds: number;
  amenities: string[];
  image: string;
  rooms: Room[];
}

export interface BuildingsResponse {
  success: boolean;
  buildings: Building[];
}

export interface BuildingResponse {
  success: boolean;
  building: Building;
}

export class HostelDataService {
  // Get all buildings
  async getAllBuildings(): Promise<Building[]> {
    const response = await fetch(`${API_CONFIG.baseURL}/api/buildings`, {
      method: 'GET',
      headers: API_CONFIG.headers,
    });

    if (!response.ok) {
      throw new Error('Failed to fetch buildings');
    }

    const data: BuildingsResponse = await response.json();
    return data.buildings;
  }

  // Get specific building by ID
  async getBuildingById(buildingId: string): Promise<Building> {
    const response = await fetch(
      `${API_CONFIG.baseURL}/api/buildings/${buildingId}`,
      {
        method: 'GET',
        headers: API_CONFIG.headers,
      }
    );

    if (!response.ok) {
      throw new Error('Failed to fetch building');
    }

    const data: BuildingResponse = await response.json();
    return data.building;
  }

  // Search buildings
  async searchBuildings(query: string): Promise<Building[]> {
    const response = await fetch(
      `${API_CONFIG.baseURL}/api/buildings/search?q=${encodeURIComponent(query)}`,
      {
        method: 'GET',
        headers: API_CONFIG.headers,
      }
    );

    if (!response.ok) {
      throw new Error('Failed to search buildings');
    }

    const data: BuildingsResponse = await response.json();
    return data.buildings;
  }

  // Get rooms by building ID
  async getRoomsByBuilding(buildingId: string): Promise<Room[]> {
    const building = await this.getBuildingById(buildingId);
    return building.rooms;
  }

  // Get specific room
  async getRoomById(buildingId: string, roomId: string): Promise<Room> {
    const response = await fetch(
      `${API_CONFIG.baseURL}/api/buildings/${buildingId}/rooms/${roomId}`,
      {
        method: 'GET',
        headers: API_CONFIG.headers,
      }
    );

    if (!response.ok) {
      throw new Error('Failed to fetch room');
    }

    const data = await response.json();
    return data.room;
  }

  // Get available rooms in a building
  async getAvailableRooms(buildingId: string): Promise<Room[]> {
    const building = await this.getBuildingById(buildingId);
    return building.rooms.filter((room) => room.available_beds > 0);
  }

  // Get building statistics
  async getBuildingStats(buildingId: string) {
    const building = await this.getBuildingById(buildingId);
    
    const totalOccupied = building.total_beds - building.available_beds;
    const occupancyRate = (totalOccupied / building.total_beds) * 100;

    return {
      totalRooms: building.total_rooms,
      totalBeds: building.total_beds,
      availableBeds: building.available_beds,
      occupiedBeds: totalOccupied,
      occupancyRate: Math.round(occupancyRate),
    };
  }
}

// Export singleton instance
export const hostelDataService = new HostelDataService();
```

---

## 📅 Booking Integration

### Create `lib/booking.ts`

```typescript
import { API_CONFIG, getAuthHeader } from './api-config';
import { authService } from './auth';

export interface Booking {
  id: string;
  user_id: string;
  user_name: string;
  building_id: string;
  building_name: string;
  room_id: string;
  room_number: string;
  bed_id: string;
  bed_number: number;
  booking_date: string;
  status: 'active' | 'cancelled';
  created_at: string;
  updated_at: string;
}

export interface BookingRequest {
  user_id: string;
  user_name: string;
  building_id: string;
  building_name: string;
  room_id: string;
  room_number: string;
  bed_id: string;
  bed_number: number;
}

export interface BookingsResponse {
  success: boolean;
  bookings: Booking[];
}

export interface BookingResponse {
  success: boolean;
  message: string;
  booking: Booking;
}

export class BookingService {
  // Create a new booking
  async createBooking(bookingData: BookingRequest): Promise<Booking> {
    const token = authService.getToken();
    
    if (!token) {
      throw new Error('Authentication required');
    }

    const response = await fetch(`${API_CONFIG.baseURL}/api/bookings`, {
      method: 'POST',
      headers: {
        ...API_CONFIG.headers,
        ...getAuthHeader(token),
      },
      body: JSON.stringify(bookingData),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Failed to create booking');
    }

    const data: BookingResponse = await response.json();
    return data.booking;
  }

  // Get user's bookings
  async getUserBookings(userId: string): Promise<Booking[]> {
    const token = authService.getToken();
    
    if (!token) {
      throw new Error('Authentication required');
    }

    const response = await fetch(
      `${API_CONFIG.baseURL}/api/bookings/users/${userId}`,
      {
        method: 'GET',
        headers: {
          ...API_CONFIG.headers,
          ...getAuthHeader(token),
        },
      }
    );

    if (!response.ok) {
      throw new Error('Failed to fetch bookings');
    }

    const data: BookingsResponse = await response.json();
    return data.bookings;
  }

  // Get all bookings (admin only)
  async getAllBookings(): Promise<Booking[]> {
    const token = authService.getToken();
    
    if (!token) {
      throw new Error('Authentication required');
    }

    const response = await fetch(`${API_CONFIG.baseURL}/api/bookings`, {
      method: 'GET',
      headers: {
        ...API_CONFIG.headers,
        ...getAuthHeader(token),
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch all bookings');
    }

    const data: BookingsResponse = await response.json();
    return data.bookings;
  }

  // Cancel a booking
  async cancelBooking(bookingId: string): Promise<void> {
    const token = authService.getToken();
    
    if (!token) {
      throw new Error('Authentication required');
    }

    const response = await fetch(
      `${API_CONFIG.baseURL}/api/bookings/${bookingId}/cancel`,
      {
        method: 'PUT',
        headers: {
          ...API_CONFIG.headers,
          ...getAuthHeader(token),
        },
      }
    );

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Failed to cancel booking');
    }
  }

  // Get active booking for user
  async getActiveBooking(userId: string): Promise<Booking | null> {
    const bookings = await this.getUserBookings(userId);
    const active = bookings.find((booking) => booking.status === 'active');
    return active || null;
  }

  // Check if user has active booking
  async hasActiveBooking(userId: string): Promise<boolean> {
    const active = await this.getActiveBooking(userId);
    return active !== null;
  }
}

// Export singleton instance
export const bookingService = new BookingService();
```

---

## 🎣 Update React Hooks

### Update `hooks/use-hostel-data.ts`

```typescript
import { useState, useEffect } from 'react';
import { hostelDataService, Building } from '@/lib/data';
import { bookingService, Booking, BookingRequest } from '@/lib/booking';
import { authService } from '@/lib/auth';
import { toast } from 'sonner';

export function useHostelData() {
  const [buildings, setBuildings] = useState<Building[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadBuildings();
  }, []);

  const loadBuildings = async () => {
    try {
      setLoading(true);
      const data = await hostelDataService.getAllBuildings();
      setBuildings(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load buildings');
      toast.error('Failed to load buildings');
    } finally {
      setLoading(false);
    }
  };

  const bookBed = async (
    buildingId: string,
    buildingName: string,
    roomId: string,
    roomNumber: string,
    bedId: string,
    bedNumber: number
  ) => {
    try {
      const user = authService.getCurrentUser();
      
      if (!user) {
        toast.error('Please login to book a bed');
        return false;
      }

      // Check if user already has an active booking
      const hasBooking = await bookingService.hasActiveBooking(user.id);
      if (hasBooking) {
        toast.error('You already have an active booking');
        return false;
      }

      const bookingData: BookingRequest = {
        user_id: user.id,
        user_name: user.name,
        building_id: buildingId,
        building_name: buildingName,
        room_id: roomId,
        room_number: roomNumber,
        bed_id: bedId,
        bed_number: bedNumber,
      };

      await bookingService.createBooking(bookingData);
      
      // Reload buildings to update availability
      await loadBuildings();
      
      toast.success('Bed booked successfully!');
      return true;
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to book bed';
      toast.error(message);
      return false;
    }
  };

  const cancelBooking = async (bookingId: string) => {
    try {
      await bookingService.cancelBooking(bookingId);
      
      // Reload buildings to update availability
      await loadBuildings();
      
      toast.success('Booking cancelled successfully!');
      return true;
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to cancel booking';
      toast.error(message);
      return false;
    }
  };

  return {
    buildings,
    loading,
    error,
    bookBed,
    cancelBooking,
    refresh: loadBuildings,
  };
}
```

### Update `hooks/use-auth.tsx`

```typescript
'use client';

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { authService, User } from '@/lib/auth';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<boolean>;
  signup: (email: string, password: string, name: string, role?: 'student' | 'admin') => Promise<boolean>;
  logout: () => void;
  isAuthenticated: boolean;
  isAdmin: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    // Check if user is already logged in
    const currentUser = authService.getCurrentUser();
    const token = authService.getToken();

    if (currentUser && token) {
      // Validate token with backend
      authService.validateToken(token).then((isValid) => {
        if (isValid) {
          setUser(currentUser);
        } else {
          // Token is invalid, clear auth data
          authService.logout();
        }
        setLoading(false);
      });
    } else {
      setLoading(false);
    }
  }, []);

  const login = async (email: string, password: string): Promise<boolean> => {
    try {
      const response = await authService.login(email, password);
      setUser(response.user);
      toast.success('Login successful!');
      return true;
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Login failed';
      toast.error(message);
      return false;
    }
  };

  const signup = async (
    email: string,
    password: string,
    name: string,
    role: 'student' | 'admin' = 'student'
  ): Promise<boolean> => {
    try {
      const response = await authService.signup(email, password, name, role);
      setUser(response.user);
      toast.success('Account created successfully!');
      return true;
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Signup failed';
      toast.error(message);
      return false;
    }
  };

  const logout = () => {
    authService.logout();
    setUser(null);
    toast.success('Logged out successfully');
    router.push('/login');
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        signup,
        logout,
        isAuthenticated: !!user,
        isAdmin: user?.role === 'admin',
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
```

---

## 🔄 Update Components

### Update `components/room/bed-grid.tsx`

```typescript
'use client';

import { Building, Room, Bed } from '@/lib/data';
import { useAuth } from '@/hooks/use-auth';
import { useHostelData } from '@/hooks/use-hostel-data';
import { Button } from '@/components/ui/button';
import { toast } from 'sonner';

interface BedGridProps {
  building: Building;
  room: Room;
}

export function BedGrid({ building, room }: BedGridProps) {
  const { isAuthenticated } = useAuth();
  const { bookBed } = useHostelData();

  const handleBookBed = async (bed: Bed) => {
    if (!isAuthenticated) {
      toast.error('Please login to book a bed');
      return;
    }

    if (bed.is_occupied) {
      toast.error('This bed is already occupied');
      return;
    }

    const success = await bookBed(
      building.id,
      building.name,
      room.id,
      room.number,
      bed.id,
      bed.number
    );

    if (success) {
      // Optionally refresh the page or navigate
      window.location.reload();
    }
  };

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {room.beds.map((bed) => (
        <div
          key={bed.id}
          className={`p-4 rounded-lg border-2 ${
            bed.is_occupied
              ? 'border-red-300 bg-red-50'
              : 'border-green-300 bg-green-50'
          }`}
        >
          <div className="flex flex-col items-center gap-2">
            <span className="text-lg font-semibold">Bed {bed.number}</span>
            <span
              className={`text-sm ${
                bed.is_occupied ? 'text-red-600' : 'text-green-600'
              }`}
            >
              {bed.is_occupied ? 'Occupied' : 'Available'}
            </span>
            {bed.is_occupied && bed.occupied_by_name && (
              <span className="text-xs text-gray-500">{bed.occupied_by_name}</span>
            )}
            {!bed.is_occupied && (
              <Button
                size="sm"
                onClick={() => handleBookBed(bed)}
                disabled={!isAuthenticated}
              >
                Book Now
              </Button>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
```

---

## 📄 Update Pages

### Update Login Page `app/login/page.tsx`

```typescript
'use client';

import { useState } from 'react';
import { useAuth } from '@/hooks/use-auth';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    const success = await login(email, password);
    
    if (success) {
      router.push('/');
    }
    
    setLoading(false);
  };

  return (
    <div className="min-h-screen flex items-center justify-center">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Login to Hostel Management</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="email" className="block text-sm font-medium mb-1">
                Email
              </label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="student@example.com"
                required
              />
            </div>
            <div>
              <label htmlFor="password" className="block text-sm font-medium mb-1">
                Password
              </label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                required
              />
            </div>
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? 'Logging in...' : 'Login'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
```

---

## 🧪 Testing Integration

### Test Script (`test-frontend-integration.ts`)

```typescript
import { authService } from '@/lib/auth';
import { hostelDataService } from '@/lib/data';
import { bookingService } from '@/lib/booking';

async function testIntegration() {
  console.log('🧪 Testing Frontend Integration...\n');

  try {
    // Test 1: Signup
    console.log('1️⃣ Testing Signup...');
    const signupResult = await authService.signup(
      `test${Date.now()}@example.com`,
      'test123',
      'Test User',
      'student'
    );
    console.log('✅ Signup successful:', signupResult.user.email);

    // Test 2: Get Buildings
    console.log('\n2️⃣ Testing Get Buildings...');
    const buildings = await hostelDataService.getAllBuildings();
    console.log(`✅ Found ${buildings.length} buildings`);

    // Test 3: Get Specific Building
    console.log('\n3️⃣ Testing Get Building Details...');
    const building = await hostelDataService.getBuildingById(buildings[0].id);
    console.log(`✅ Building: ${building.name} with ${building.total_rooms} rooms`);

    // Test 4: Create Booking
    console.log('\n4️⃣ Testing Create Booking...');
    const room = building.rooms[0];
    const bed = room.beds.find((b) => !b.is_occupied);

    if (bed) {
      const booking = await bookingService.createBooking({
        user_id: signupResult.user.id,
        user_name: signupResult.user.name,
        building_id: building.id,
        building_name: building.name,
        room_id: room.id,
        room_number: room.number,
        bed_id: bed.id,
        bed_number: bed.number,
      });
      console.log('✅ Booking created:', booking.id);

      // Test 5: Get User Bookings
      console.log('\n5️⃣ Testing Get User Bookings...');
      const bookings = await bookingService.getUserBookings(signupResult.user.id);
      console.log(`✅ Found ${bookings.length} booking(s)`);

      // Test 6: Cancel Booking
      console.log('\n6️⃣ Testing Cancel Booking...');
      await bookingService.cancelBooking(booking.id);
      console.log('✅ Booking cancelled successfully');
    }

    console.log('\n✅ All tests passed!');
  } catch (error) {
    console.error('❌ Test failed:', error);
  }
}

// Run tests
testIntegration();
```

---

## 🔒 Protected Routes

### Create Middleware `middleware.ts`

```typescript
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  const token = request.cookies.get('token')?.value;
  const { pathname } = request.nextUrl;

  // Protected routes
  const protectedRoutes = ['/booking', '/profile', '/admin'];
  const isProtectedRoute = protectedRoutes.some((route) =>
    pathname.startsWith(route)
  );

  if (isProtectedRoute && !token) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // Admin routes
  if (pathname.startsWith('/admin')) {
    // Additional check for admin role can be done here
    // by validating the token with the backend
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/booking/:path*', '/profile/:path*', '/admin/:path*'],
};
```

---

## 🎨 Error Handling

### Create Error Handler `lib/error-handler.ts`

```typescript
export class APIError extends Error {
  constructor(
    message: string,
    public statusCode?: number,
    public data?: any
  ) {
    super(message);
    this.name = 'APIError';
  }
}

export function handleAPIError(error: unknown): string {
  if (error instanceof APIError) {
    return error.message;
  }

  if (error instanceof Error) {
    return error.message;
  }

  return 'An unexpected error occurred';
}

export async function fetchWithErrorHandling(
  url: string,
  options?: RequestInit
): Promise<Response> {
  try {
    const response = await fetch(url, options);

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new APIError(
        errorData.message || `HTTP ${response.status}: ${response.statusText}`,
        response.status,
        errorData
      );
    }

    return response;
  } catch (error) {
    if (error instanceof APIError) {
      throw error;
    }

    throw new APIError('Network error: Unable to connect to server');
  }
}
```

---

## 📱 Example: Complete Booking Flow

```typescript
'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/hooks/use-auth';
import { hostelDataService, Building } from '@/lib/data';
import { bookingService } from '@/lib/booking';
import { useRouter } from 'next/navigation';

export default function BookingPage() {
  const { user, isAuthenticated } = useAuth();
  const [buildings, setBuildings] = useState<Building[]>([]);
  const [selectedBuilding, setSelectedBuilding] = useState<Building | null>(null);
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login');
      return;
    }

    loadBuildings();
  }, [isAuthenticated]);

  const loadBuildings = async () => {
    try {
      const data = await hostelDataService.getAllBuildings();
      setBuildings(data);
    } catch (error) {
      console.error('Failed to load buildings:', error);
    }
  };

  const handleBookBed = async (
    buildingId: string,
    roomId: string,
    bedId: string
  ) => {
    if (!user || loading) return;

    setLoading(true);
    try {
      const building = buildings.find((b) => b.id === buildingId);
      const room = building?.rooms.find((r) => r.id === roomId);
      const bed = room?.beds.find((b) => b.id === bedId);

      if (!building || !room || !bed) {
        throw new Error('Invalid selection');
      }

      await bookingService.createBooking({
        user_id: user.id,
        user_name: user.name,
        building_id: building.id,
        building_name: building.name,
        room_id: room.id,
        room_number: room.number,
        bed_id: bed.id,
        bed_number: bed.number,
      });

      // Refresh buildings to update availability
      await loadBuildings();
      
      alert('Booking successful!');
    } catch (error) {
      console.error('Booking failed:', error);
      alert('Failed to create booking');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto py-8">
      <h1 className="text-3xl font-bold mb-6">Book Your Bed</h1>
      
      {/* Building selection and booking UI */}
      {/* Implementation details here */}
    </div>
  );
}
```

---

## ✅ Checklist

Before deploying to production:

- [ ] Update `NEXT_PUBLIC_API_BASE_URL` in `.env.production`
- [ ] Enable CORS for your production domain in backend
- [ ] Test all authentication flows (signup, login, logout)
- [ ] Test booking creation and cancellation
- [ ] Test admin-specific features
- [ ] Implement proper error boundaries
- [ ] Add loading states for all async operations
- [ ] Test on different browsers and devices
- [ ] Implement rate limiting on frontend
- [ ] Add analytics tracking
- [ ] Test with production database

---

## 🚀 Deployment

### Update Production Environment

```env
# .env.production
NEXT_PUBLIC_API_BASE_URL=https://api.yourdomain.com
NEXT_PUBLIC_API_TIMEOUT=30000
```

### Build and Deploy

```bash
# Build frontend
npm run build

# Test production build locally
npm run start

# Deploy to Vercel/Netlify/etc
vercel deploy --prod
```

---

## 📞 Support

If you encounter issues:

1. Check browser console for errors
2. Verify backend is running: `docker-compose ps`
3. Test backend endpoints: See `API_TESTING.md`
4. Check CORS configuration in backend
5. Verify JWT token is being sent with requests

---

**Frontend integration complete! Your Next.js app is now connected to the Golang microservices backend. 🎉**
