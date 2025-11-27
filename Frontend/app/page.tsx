"use client"

import { useAuth } from "@/hooks/use-auth"
import { Header } from "@/components/layout/header"
import { DashboardContent } from "@/components/dashboard/dashboard-content"
import { Loader2, Sparkles } from "lucide-react"
import AuthPage from "@/components/auth/auth-page"

export default function HomePage() {
  const { isAuthenticated, loading } = useAuth()

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 via-rose-50/30 to-pink-50/20">
        <div className="text-center animate-fade-in">
          <div className="relative mb-8">
            <Loader2 className="h-16 w-16 animate-spin text-rose-500 mx-auto" />
            <div className="absolute inset-0 flex items-center justify-center">
              <Sparkles className="h-6 w-6 text-rose-400 animate-pulse" />
            </div>
          </div>
          <p className="text-slate-600 text-xl font-medium">Loading your experience...</p>
          <p className="text-slate-500 text-sm mt-2">Preparing your premium hostel dashboard</p>
        </div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return <AuthPage />
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-rose-50/20 to-pink-50/10">
      <Header />
      <main className="container mx-auto px-4 py-8">
        <div className="mb-12 text-center animate-fade-in">
          <h1 className="text-5xl font-serif font-light text-slate-800 mb-6 text-balance tracking-tight">
            Welcome to <span className="text-rose-600">Campus Hostel</span>
          </h1>
          <p className="text-slate-600 text-xl max-w-3xl mx-auto text-pretty leading-relaxed">
            Discover and book your perfect room across our premium campus buildings. Experience comfort, community, and
            convenience like never before.
          </p>
          <div className="w-24 h-1 bg-gradient-to-r from-rose-500 to-pink-500 rounded-full mx-auto mt-8"></div>
        </div>
        <div className="animate-slide-up">
          <DashboardContent />
        </div>
      </main>
    </div>
  )
}
