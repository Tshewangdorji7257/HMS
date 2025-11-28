"use client"

import { useState, useEffect, createContext, useContext, type ReactNode } from "react"
import { authService, type AuthState } from "@/lib/auth"

interface AuthContextType extends AuthState {
  login: (email: string, password: string) => Promise<{ success: boolean; error?: string }>
  signup: (email: string, password: string, name: string) => Promise<{ success: boolean; error?: string }>
  logout: () => void
  loading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [authState, setAuthState] = useState<AuthState>({ user: null, isAuthenticated: false, token: null })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Initialize auth state from localStorage
    const initAuth = () => {
      const state = authService.getAuthState()
      setAuthState(state)
      setLoading(false)
    }

    initAuth()
  }, [])

  const login = async (email: string, password: string) => {
    setLoading(true)
    const result = await authService.login(email, password)

    if (result.success) {
      const newState = authService.getAuthState()
      setAuthState(newState)
    }

    setLoading(false)
    return result
  }

  const signup = async (email: string, password: string, name: string) => {
    setLoading(true)
    const result = await authService.signup(email, password, name)

    if (result.success) {
      const newState = authService.getAuthState()
      setAuthState(newState)
    }

    setLoading(false)
    return result
  }

  const logout = () => {
    authService.logout()
    setAuthState({ user: null, isAuthenticated: false, token: null })
  }

  return (
    <AuthContext.Provider
      value={{
        ...authState,
        login,
        signup,
        logout,
        loading,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}
