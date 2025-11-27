// Backend admin authentication using role-based access
import { apiFetch, handleApiError } from './api-config'

export interface AdminUser {
  id: string
  email: string
  name: string
  role: 'admin'
  created_at?: string
}

export interface AdminAuthState {
  user: AdminUser | null
  isAuthenticated: boolean
  token: string | null
}

export interface AdminAuthResponse {
  success: boolean
  message: string
  token: string
  user: AdminUser
}

class AdminAuthService {
  private readonly storageKey = "hostel-admin-auth"
  private readonly tokenKey = "hostel-admin-auth-token"

  getAuthState(): AdminAuthState {
    if (typeof globalThis.window === "undefined") {
      return { user: null, isAuthenticated: false, token: null }
    }

    try {
      const stored = localStorage.getItem(this.storageKey)
      const token = localStorage.getItem(this.tokenKey)
      
      if (stored && token) {
        const user = JSON.parse(stored)
        return { user, isAuthenticated: true, token }
      }
    } catch (error) {
      console.error("Error reading admin auth state:", error)
    }

    return { user: null, isAuthenticated: false, token: null }
  }

  async login(email: string, password: string): Promise<{ success: boolean; error?: string }> {
    try {
      if (!email || !password) {
        return { success: false, error: "Email and password are required" }
      }

      if (password.length < 6) {
        return { success: false, error: "Password must be at least 6 characters" }
      }

      // Call backend API with admin role
      const response = await apiFetch<AdminAuthResponse>('/api/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      })

      if (response.success && response.token && response.user) {
        // Verify user has admin role
        if (response.user.role !== 'admin') {
          return { success: false, error: "Access denied. Admin privileges required." }
        }

        // Store user and token
        localStorage.setItem(this.storageKey, JSON.stringify(response.user))
        localStorage.setItem(this.tokenKey, response.token)
        return { success: true }
      }

      return { success: false, error: response.message || 'Login failed' }
    } catch (error: any) {
      console.error('Admin login error:', error)
      return { success: false, error: handleApiError(error) }
    }
  }

  async signup(email: string, password: string, name: string): Promise<{ success: boolean; error?: string }> {
    try {
      if (!email || !password || !name) {
        return { success: false, error: "All fields are required" }
      }

      if (password.length < 6) {
        return { success: false, error: "Password must be at least 6 characters" }
      }

      // Call backend API with admin role
      const response = await apiFetch<AdminAuthResponse>('/api/auth/signup', {
        method: 'POST',
        body: JSON.stringify({ email, password, name, role: 'admin' }),
      })

      if (response.success && response.token && response.user) {
        // Store user and token
        localStorage.setItem(this.storageKey, JSON.stringify(response.user))
        localStorage.setItem(this.tokenKey, response.token)
        return { success: true }
      }

      return { success: false, error: response.message || 'Signup failed' }
    } catch (error: any) {
      console.error('Admin signup error:', error)
      return { success: false, error: handleApiError(error) }
    }
  }

  async getProfile(token?: string): Promise<AdminUser | null> {
    try {
      const authToken = token || localStorage.getItem(this.tokenKey)
      if (!authToken) {
        return null
      }

      const response = await apiFetch<{ success: boolean; user: AdminUser }>('/api/auth/profile', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${authToken}`,
        },
      })

      if (response.success && response.user && response.user.role === 'admin') {
        // Update stored user
        localStorage.setItem(this.storageKey, JSON.stringify(response.user))
        return response.user
      }

      return null
    } catch (error) {
      console.error('Get admin profile error:', error)
      return null
    }
  }

  logout(): void {
    localStorage.removeItem(this.storageKey)
    localStorage.removeItem(this.tokenKey)
  }

  getToken(): string | null {
    return typeof globalThis.window !== 'undefined' ? localStorage.getItem(this.tokenKey) : null
  }
}

export const adminAuthService = new AdminAuthService()
