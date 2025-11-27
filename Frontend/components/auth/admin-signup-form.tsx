"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardHeader } from "@/components/ui/card"
import { useAdminAuth } from "@/hooks/use-admin-auth"
import { Loader2, Mail, Lock, Eye, EyeOff, User, ShieldPlus, Sparkles } from "lucide-react"

interface Props {
  onToggleMode: () => void
}

export function AdminSignupForm({ onToggleMode }: Props) {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [name, setName] = useState("")
  const [error, setError] = useState("")
  const [showPassword, setShowPassword] = useState(false)
  const { signup, loading } = useAdminAuth()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    const res = await signup(email, password, name)
    if (!res.success) setError(res.error || "Signup failed")
  }

  return (
    <Card className="relative w-full max-w-md overflow-hidden border border-rose-100/60 bg-white/80 backdrop-blur-xl shadow-2xl rounded-3xl">
      <div className="absolute inset-0 bg-gradient-to-br from-white/70 via-rose-50/50 to-pink-50/40" />
      <CardHeader className="relative space-y-6 pb-2 pt-8">
        <div className="flex flex-col items-center gap-4">
          <div className="relative">
            <div className="p-5 rounded-2xl bg-gradient-to-br from-rose-500/20 to-pink-500/20 border border-rose-200/60 shadow-xl">
              <ShieldPlus className="h-10 w-10 text-rose-600" />
            </div>
            <Sparkles className="h-5 w-5 text-rose-400 absolute -top-1 -right-1 animate-pulse" />
          </div>
          <div className="text-center space-y-2">
            <h2 className="text-3xl font-serif font-light tracking-tight text-slate-800">Create Admin</h2>
            <p className="text-sm text-slate-500">Provision a new administrative account</p>
          </div>
        </div>
      </CardHeader>
      <CardContent className="relative p-8 pt-2">
        <form onSubmit={handleSubmit} className="space-y-6 animate-in fade-in slide-in-from-bottom-2 duration-300">
          <div className="space-y-2">
            <Label htmlFor="admin-name" className="text-slate-700 font-medium">Full Name</Label>
            <div className="relative group">
              <User className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-slate-400 group-focus-within:text-rose-500 transition-colors" />
              <Input
                id="admin-name"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
                placeholder="Jane Doe"
                className="pl-10 h-12 rounded-xl border-slate-200 focus:border-rose-300 focus:ring-rose-200 transition-all"
              />
            </div>
          </div>
          <div className="space-y-2">
            <Label htmlFor="admin-email" className="text-slate-700 font-medium">Email Address</Label>
            <div className="relative group">
              <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-slate-400 group-focus-within:text-rose-500 transition-colors" />
              <Input
                id="admin-email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                placeholder="admin@example.com"
                className="pl-10 h-12 rounded-xl border-slate-200 focus:border-rose-300 focus:ring-rose-200 transition-all"
              />
            </div>
          </div>
          <div className="space-y-2">
            <Label htmlFor="admin-password" className="text-slate-700 font-medium">Password</Label>
            <div className="relative group">
              <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-slate-400 group-focus-within:text-rose-500 transition-colors" />
              <Input
                id="admin-password"
                type={showPassword ? "text" : "password"}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                placeholder="Create a secure password"
                className="pl-10 pr-10 h-12 rounded-xl border-slate-200 focus:border-rose-300 focus:ring-rose-200 transition-all"
              />
              <button
                type="button"
                onClick={() => setShowPassword(v => !v)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                tabIndex={-1}
              >
                {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
              </button>
            </div>
          </div>
          {error && (
            <div className="text-sm text-destructive bg-destructive/10 p-3 rounded-lg border border-destructive/20 animate-in fade-in">
              {error}
            </div>
          )}
          <Button
            type="submit"
            className="w-full h-12 bg-gradient-to-r from-rose-500 to-pink-500 hover:from-rose-600 hover:to-pink-600 text-white font-medium rounded-xl shadow-lg hover:shadow-xl transition-all disabled:opacity-60"
            disabled={loading}
          >
            {loading ? (<><Loader2 className="mr-2 h-4 w-4 animate-spin" /> Creating...</>) : (<>Create Admin</>)}
          </Button>
        </form>
        <div className="mt-8 text-center text-sm text-slate-600">
          <span className="mr-1">Already have admin credentials?</span>
          <button
            onClick={onToggleMode}
            className="font-medium text-rose-600 hover:text-rose-700 hover:underline underline-offset-4 transition-colors"
          >
            Sign in
          </button>
        </div>
      </CardContent>
      <div className="absolute -bottom-16 -right-16 w-48 h-48 bg-gradient-to-br from-rose-300/30 to-pink-300/30 rounded-full blur-3xl" />
      <div className="absolute -top-10 -left-10 w-36 h-36 bg-gradient-to-br from-pink-200/40 to-rose-200/40 rounded-full blur-2xl" />
    </Card>
  )
}
