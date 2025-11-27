"use client"

import { useState, useEffect, createContext, useContext, type ReactNode } from "react"
import { adminAuthService, type AdminAuthState } from "@/lib/admin"

interface AdminAuthContextType extends AdminAuthState {
  login: (email: string, password: string) => Promise<{ success: boolean; error?: string }>
  signup: (email: string, password: string, name: string) => Promise<{ success: boolean; error?: string }>
  logout: () => void
  loading: boolean
}

const AdminAuthContext = createContext<AdminAuthContextType | undefined>(undefined)

export function AdminAuthProvider({ children }: { children: ReactNode }) {
  const [authState, setAuthState] = useState<AdminAuthState>({ user: null, isAuthenticated: false, token: null })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const state = adminAuthService.getAuthState()
    setAuthState(state)
    setLoading(false)
  }, [])

  const login = async (email: string, password: string) => {
    setLoading(true)
    const result = await adminAuthService.login(email, password)
    if (result.success) {
      setAuthState(adminAuthService.getAuthState())
    }
    setLoading(false)
    return result
  }

  const signup = async (email: string, password: string, name: string) => {
    setLoading(true)
    const result = await adminAuthService.signup(email, password, name)
    if (result.success) {
      setAuthState(adminAuthService.getAuthState())
    }
    setLoading(false)
    return result
  }

  const logout = () => {
    adminAuthService.logout()
    setAuthState({ user: null, isAuthenticated: false, token: null })
  }

  return (
    <AdminAuthContext.Provider value={{ ...authState, login, signup, logout, loading }}>
      {children}
    </AdminAuthContext.Provider>
  )
}

export function useAdminAuth() {
  const context = useContext(AdminAuthContext)
  if (context === undefined) {
    throw new Error("useAdminAuth must be used within an AdminAuthProvider")
  }
  return context
}
