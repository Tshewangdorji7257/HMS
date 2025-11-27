"use client"

import { useState } from "react"
import { AdminAuthProvider, useAdminAuth } from "@/hooks/use-admin-auth"
import { AdminLoginForm } from "@/components/auth/admin-login-form"
import { AdminSignupForm } from "@/components/auth/admin-signup-form"
import { AdminBookingsTable } from "@/components/booking/admin-bookings-table"
import { Button } from "@/components/ui/button"
import { Building2, Sparkles } from "lucide-react"

function AdminContent() {
  const [mode, setMode] = useState<"login" | "signup">("login")
  const { user, isAuthenticated, logout } = useAdminAuth()

  if (!isAuthenticated) {
    return (
      <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center py-12 px-4 relative overflow-hidden bg-gradient-to-br from-slate-50 via-rose-50/40 to-pink-50/30">
        <div className="absolute inset-0 opacity-40 pointer-events-none">
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_25%_25%,rgba(244,114,182,0.18),transparent_60%)]" />
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_75%_75%,rgba(244,63,94,0.18),transparent_60%)]" />
        </div>
        <div className="absolute -top-24 -left-24 w-72 h-72 bg-gradient-to-br from-rose-300/30 to-pink-300/30 rounded-full blur-3xl" />
        <div className="absolute -bottom-24 -right-24 w-80 h-80 bg-gradient-to-br from-pink-200/30 to-rose-200/30 rounded-full blur-3xl" />

        <div className="relative w-full max-w-6xl grid md:grid-cols-2 gap-10 items-center">
          <div className="hidden md:flex flex-col gap-8 pl-4">
            <div className="space-y-6 animate-in fade-in slide-in-from-left-4 duration-500">
              <div className="flex items-center gap-4">
                <div className="relative">
                  <div className="p-6 rounded-2xl bg-gradient-to-br from-rose-500/20 to-pink-500/20 backdrop-blur-sm border border-rose-200/50 shadow-2xl">
                    <Building2 className="h-12 w-12 text-rose-600" />
                  </div>
                  <Sparkles className="h-6 w-6 text-rose-400 absolute -top-2 -right-2 animate-pulse" />
                </div>
                <h1 className="text-5xl font-serif font-light tracking-tight text-slate-800">
                  Admin <span className="text-rose-600 font-medium">Console</span>
                </h1>
              </div>
              <p className="text-slate-600 text-lg leading-relaxed max-w-md">
                Manage bookings, rooms, and user activities with a refined and secure administrative interface.
              </p>
              <div className="grid grid-cols-2 gap-4 max-w-md">
                <div className="p-4 rounded-xl bg-white/70 backdrop-blur-sm border border-white/60 shadow hover:shadow-md transition-all">
                  <p className="text-sm font-medium text-slate-700 mb-1">Real-time Data</p>
                  <p className="text-xs text-slate-500">Stay updated with live booking information</p>
                </div>
                <div className="p-4 rounded-xl bg-white/70 backdrop-blur-sm border border-white/60 shadow hover:shadow-md transition-all">
                  <p className="text-sm font-medium text-slate-700 mb-1">Secure Access</p>
                  <p className="text-xs text-slate-500">Multi-level authentication & roles</p>
                </div>
                <div className="p-4 rounded-xl bg-white/70 backdrop-blur-sm border border-white/60 shadow hover:shadow-md transition-all">
                  <p className="text-sm font-medium text-slate-700 mb-1">Insights</p>
                  <p className="text-xs text-slate-500">Track occupancy & performance metrics</p>
                </div>
                <div className="p-4 rounded-xl bg-white/70 backdrop-blur-sm border border-white/60 shadow hover:shadow-md transition-all">
                  <p className="text-sm font-medium text-slate-700 mb-1">Scalable</p>
                  <p className="text-xs text-slate-500">Architected for growth and resilience</p>
                </div>
              </div>
            </div>
          </div>

          <div className="flex justify-center animate-in fade-in slide-in-from-right-4 duration-500">
            {mode === "login" ? (
              <AdminLoginForm onToggleMode={() => setMode("signup")} />
            ) : (
              <AdminSignupForm onToggleMode={() => setMode("login")} />
            )}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      <div className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div className="space-y-2">
          <h1 className="text-4xl font-serif font-light tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-rose-600 to-pink-500">Admin Dashboard</h1>
          <p className="text-sm text-muted-foreground max-w-prose">Monitor bookings, manage capacity and ensure a seamless hostel experience.</p>
        </div>
        <div className="flex items-center gap-4 bg-white/70 backdrop-blur rounded-xl border border-slate-200/60 px-4 py-2 shadow-sm">
          <div className="text-xs text-slate-600 leading-tight">
            <div className="font-medium text-slate-800">{user?.name}</div>
            <div className="text-[10px] uppercase tracking-wide text-slate-500">Administrator</div>
          </div>
          <Button size="sm" variant="ghost" onClick={() => logout()} className="hover:text-rose-600">Sign out</Button>
        </div>
      </div>

      <div className="grid gap-8">
        <AdminBookingsTable />
      </div>
    </div>
  )
}

export default function AdminPage() {
  return (
    <AdminAuthProvider>
      <div className="container mx-auto p-6">
        <AdminContent />
      </div>
    </AdminAuthProvider>
  )
}
