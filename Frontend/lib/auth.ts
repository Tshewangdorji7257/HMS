// Backend authentication system
import { API_CONFIG, apiFetch, handleApiError } from './api-config'

export interface User {
  id: string
  email: string
  name: string
  role: 'student' | 'admin'
  created_at?: string
}

export interface AuthState {
  user: User | null
  isAuthenticated: boolean
  token: string | null
}

export interface AuthResponse {
  success: boolean
  message: string
  token: string
  user: User
}

class AuthService {
  private storageKey = "hostel-auth"
  private tokenKey = "hostel-auth-token"

  getAuthState(): AuthState {
    if (typeof window === "undefined") {
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
      console.error("Error reading auth state:", error)
    }

    return { user: null, isAuthenticated: false, token: null }
  }

  async login(email: string, password: string): Promise<{ success: boolean; error?: string }> {
    try {
      // Validate input
      if (!email || !password) {
        return { success: false, error: "Email and password are required" }
      }

      if (password.length < 6) {
        return { success: false, error: "Password must be at least 6 characters" }
      }

      // Call backend API
      const response = await apiFetch<AuthResponse>('/api/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      })

      if (response.success && response.token && response.user) {
        // Store user and token
        localStorage.setItem(this.storageKey, JSON.stringify(response.user))
        localStorage.setItem(this.tokenKey, response.token)
        return { success: true }
      }

      return { success: false, error: response.message || 'Login failed' }
    } catch (error: any) {
      console.error('Login error:', error)
      return { success: false, error: handleApiError(error) }
    }
  }

  async signup(
    email: string, 
    password: string, 
    name: string,
    role: 'student' | 'admin' = 'student'
  ): Promise<{ success: boolean; error?: string }> {
    try {
      // Validate input
      if (!email || !password || !name) {
        return { success: false, error: "All fields are required" }
      }

      if (password.length < 6) {
        return { success: false, error: "Password must be at least 6 characters" }
      }

      // Call backend API
      const response = await apiFetch<AuthResponse>('/api/auth/signup', {
        method: 'POST',
        body: JSON.stringify({ email, password, name, role }),
      })

      if (response.success && response.token && response.user) {
        // Store user and token
        localStorage.setItem(this.storageKey, JSON.stringify(response.user))
        localStorage.setItem(this.tokenKey, response.token)
        return { success: true }
      }

      return { success: false, error: response.message || 'Signup failed' }
    } catch (error: any) {
      console.error('Signup error:', error)
      return { success: false, error: handleApiError(error) }
    }
  }

  async getProfile(token?: string): Promise<User | null> {
    try {
      const authToken = token || localStorage.getItem(this.tokenKey)
      if (!authToken) {
        return null
      }

      const response = await apiFetch<{ success: boolean; user: User }>('/api/auth/profile', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${authToken}`,
        },
      })

      if (response.success && response.user) {
        // Update stored user
        localStorage.setItem(this.storageKey, JSON.stringify(response.user))
        return response.user
      }

      return null
    } catch (error) {
      console.error('Get profile error:', error)
      return null
    }
  }

  async validateToken(token?: string): Promise<boolean> {
    try {
      const authToken = token || localStorage.getItem(this.tokenKey)
      if (!authToken) {
        return false
      }

      const response = await apiFetch<{ success: boolean; valid: boolean }>('/api/auth/validate', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${authToken}`,
        },
      })

      return response.success && response.valid
    } catch (error) {
      console.error('Validate token error:', error)
      return false
    }
  }

  logout(): void {
    localStorage.removeItem(this.storageKey)
    localStorage.removeItem(this.tokenKey)
  }

  getToken(): string | null {
    return typeof window !== 'undefined' ? localStorage.getItem(this.tokenKey) : null
  }
}

export const authService = new AuthService()
