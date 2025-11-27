"use client"

import type React from "react"

import { useState } from "react"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useAuth } from "@/hooks/use-auth"
import { Building2, Mail, Lock, User, Eye, EyeOff, Sparkles, Loader2 } from "lucide-react"
import { useToast } from "@/hooks/use-toast"
import { useRouter } from "next/navigation"

export default function AuthPage() {
  const [isLogin, setIsLogin] = useState(true)
  const [showPassword, setShowPassword] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    password: "",
  })
  const { login, signup } = useAuth()
  const { toast } = useToast()
  const router = useRouter()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)

    try {
      let result
      if (isLogin) {
        result = await login(formData.email, formData.password)
      } else {
        result = await signup(formData.email, formData.password, formData.name)
      }

      if (result.success) {
        toast({
          title: isLogin ? "Welcome back!" : "Account created!",
          description: isLogin ? "You have successfully signed in." : "Your account has been created successfully.",
        })
        // Redirect to dashboard after a short delay so the toast is visible
        setTimeout(() => {
          router.push('/')
        }, 600)
      } else {
        toast({
          title: "Authentication failed",
          description: result.error || "An error occurred. Please try again.",
          variant: "destructive",
        })
      }
    } catch (error) {
      toast({
        title: "Error",
        description: "An unexpected error occurred. Please try again.",
        variant: "destructive",
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleInputChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
  }

  return (
    <div className="h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 via-rose-50/30 to-pink-50/20 p-4 relative overflow-hidden">
      <div className="absolute inset-0 opacity-40">
        <div className="absolute inset-0 bg-gradient-to-br from-slate-100/50 to-transparent"></div>
        <div className="absolute top-0 left-0 w-full h-full bg-[radial-gradient(circle_at_50%_50%,rgba(241,245,249,0.4)_1px,transparent_1px)] bg-[length:60px_60px]"></div>
      </div>

      {/* Floating decorative elements */}
      <div className="absolute top-20 left-20 w-32 h-32 bg-gradient-to-br from-rose-200/20 to-pink-200/20 rounded-full blur-xl animate-pulse"></div>
      <div className="absolute bottom-20 right-20 w-40 h-40 bg-gradient-to-br from-pink-200/20 to-rose-200/20 rounded-full blur-xl animate-pulse delay-1000"></div>
      <div className="absolute top-1/2 left-10 w-24 h-24 bg-gradient-to-br from-slate-200/30 to-slate-300/30 rounded-full blur-xl animate-pulse delay-500"></div>

      <div className="w-full max-w-6xl relative z-10 grid lg:grid-cols-2 gap-8 items-center max-h-[95vh] overflow-hidden">
        {/* Left side - Branding and features */}
        <div className="hidden lg:block space-y-6 animate-fade-in">
          <div className="text-center lg:text-left">
            <h1 className="text-5xl font-serif font-light text-slate-800 mb-4 tracking-tight leading-tight">
              Hostel
              <span className="block text-rose-600 font-medium">Management System</span>
            </h1>
            <p className="text-slate-600 text-lg leading-relaxed mb-6">
              Experience student accommodation with world-class amenities and a vibrant community atmosphere.
            </p>
            <div className="w-20 h-1 bg-gradient-to-r from-rose-500 to-pink-500 rounded-full"></div>
          </div>

          {/* Feature highlights */}
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-white/60 backdrop-blur-sm rounded-xl p-4 border border-white/50 shadow-lg hover:shadow-xl transition-all duration-300">
              <div className="p-2 rounded-lg bg-gradient-to-br from-blue-500/20 to-blue-600/20 w-fit mb-2">
                <Sparkles className="h-5 w-5 text-blue-600" />
              </div>
              <h3 className="font-semibold text-slate-800 mb-1 text-sm">Community Living</h3>
              <p className="text-slate-600 text-xs leading-relaxed">
                Connect with fellow students in our vibrant community spaces
              </p>
            </div>
            <div className="bg-white/60 backdrop-blur-sm rounded-xl p-4 border border-white/50 shadow-lg hover:shadow-xl transition-all duration-300">
              <div className="p-2 rounded-lg bg-gradient-to-br from-green-500/20 to-green-600/20 w-fit mb-2">
                <Sparkles className="h-5 w-5 text-green-600" />
              </div>
              <h3 className="font-semibold text-slate-800 mb-1 text-sm">Safe & Secure</h3>
              <p className="text-slate-600 text-xs leading-relaxed">24/7 security with modern access control systems</p>
            </div>
            <div className="bg-white/60 backdrop-blur-sm rounded-xl p-4 border border-white/50 shadow-lg hover:shadow-xl transition-all duration-300">
              <div className="p-2 rounded-lg bg-gradient-to-br from-purple-500/20 to-purple-600/20 w-fit mb-2">
                <Sparkles className="h-5 w-5 text-purple-600" />
              </div>
              <h3 className="font-semibold text-slate-800 mb-1 text-sm">High-Speed WiFi</h3>
              <p className="text-slate-600 text-xs leading-relaxed">
                Blazing fast internet for all your academic needs
              </p>
            </div>
            <div className="bg-white/60 backdrop-blur-sm rounded-xl p-4 border border-white/50 shadow-lg hover:shadow-xl transition-all duration-300">
              <div className="p-2 rounded-lg bg-gradient-to-br from-amber-500/20 to-amber-600/20 w-fit mb-2">
                <Sparkles className="h-5 w-5 text-amber-600" />
              </div>
              <h3 className="font-semibold text-slate-800 mb-1 text-sm">Amenities</h3>
              <p className="text-slate-600 text-xs leading-relaxed">
                Gym, study rooms, laundry, and recreational facilities
              </p>
            </div>
          </div>
        </div>

        {/* Right side - Auth form */}
        <div className="w-full max-w-md mx-auto lg:mx-0">
          <div className="text-center mb-4 lg:hidden animate-fade-in">
            <h1 className="text-3xl font-serif font-light text-slate-800 mb-2 tracking-tight">Campus Hostel</h1>
            <p className="text-slate-600 text-base leading-relaxed">Your home away from home</p>
            <div className="w-16 h-0.5 bg-gradient-to-r from-rose-500 to-pink-500 mx-auto mt-3 rounded-full"></div>
          </div>

          <Card className="border-0 shadow-2xl bg-white/80 backdrop-blur-xl overflow-hidden">
            <div className="absolute inset-0 bg-gradient-to-br from-white/50 to-rose-50/30"></div>
            <CardContent className="p-6 relative">
              <div className="text-center mb-6">
                <div className="flex items-center justify-center space-x-2 mb-3">
                  <Sparkles className="h-4 w-4 text-rose-500" />
                  <h2 className="text-xl font-serif font-light text-slate-800">
                    {isLogin ? "Welcome Back" : "Join Campus Hostel"}
                  </h2>
                  <Sparkles className="h-4 w-4 text-rose-500" />
                </div>
                <p className="text-slate-600 text-xs">
                  {isLogin ? "Sign in to access your booking dashboard" : "Create your account to start booking rooms"}
                </p>
              </div>

              <form onSubmit={handleSubmit} className="space-y-4">
                {!isLogin && (
                  <div className="space-y-2">
                    <Label htmlFor="name" className="text-slate-700 font-medium">
                      Full Name
                    </Label>
                    <div className="relative">
                      <User className="absolute left-3 top-1/2 transform -translate-y-1/2 text-slate-400 h-5 w-5" />
                      <Input
                        id="name"
                        type="text"
                        placeholder="Enter your full name"
                        value={formData.name}
                        onChange={(e) => handleInputChange("name", e.target.value)}
                        className="pl-10 h-10 rounded-xl border-slate-200 focus:border-rose-300 focus:ring-rose-200 transition-all duration-300"
                        required
                        disabled={isSubmitting}
                      />
                    </div>
                  </div>
                )}

                <div className="space-y-2">
                  <Label htmlFor="email" className="text-slate-700 font-medium">
                    Email Address
                  </Label>
                  <div className="relative">
                    <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 text-slate-400 h-5 w-5" />
                    <Input
                      id="email"
                      type="email"
                      placeholder="Enter your email"
                      value={formData.email}
                      onChange={(e) => handleInputChange("email", e.target.value)}
                      className="pl-10 h-10 rounded-xl border-slate-200 focus:border-rose-300 focus:ring-rose-200 transition-all duration-300"
                      required
                      disabled={isSubmitting}
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="password" className="text-slate-700 font-medium">
                    Password
                  </Label>
                  <div className="relative">
                    <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 text-slate-400 h-5 w-5" />
                    <Input
                      id="password"
                      type={showPassword ? "text" : "password"}
                      placeholder="Enter your password"
                      value={formData.password}
                      onChange={(e) => handleInputChange("password", e.target.value)}
                      className="pl-10 pr-10 h-10 rounded-xl border-slate-200 focus:border-rose-300 focus:ring-rose-200 transition-all duration-300"
                      required
                      disabled={isSubmitting}
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword(!showPassword)}
                      className="absolute right-3 top-1/2 transform -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors"
                      disabled={isSubmitting}
                    >
                      {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                    </button>
                  </div>
                </div>

                <Button
                  type="submit"
                  disabled={isSubmitting}
                  className="w-full h-10 bg-gradient-to-r from-rose-500 to-pink-500 hover:from-rose-600 hover:to-pink-600 text-white font-medium rounded-xl shadow-lg hover:shadow-xl transition-all duration-300 transform hover:scale-[1.02] disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none">
                >
                  {isSubmitting ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      {isLogin ? "Signing In..." : "Creating Account..."}
                    </>
                  ) : (
                    <>{isLogin ? "Sign In" : "Create Account"}</>
                  )}
                </Button>
              </form>

              <div className="mt-6 text-center">
                <div className="relative">
                  <div className="absolute inset-0 flex items-center">
                    <div className="w-full border-t border-slate-200"></div>
                  </div>
                  <div className="relative flex justify-center text-sm">
                    <span className="px-4 bg-white text-slate-500">or</span>
                  </div>
                </div>

                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => setIsLogin(!isLogin)}
                  className="mt-4 text-slate-600 hover:text-rose-600 transition-colors"
                  disabled={isSubmitting}
                >
                  {isLogin ? "Don't have an account? Sign up" : "Already have an account? Sign in"}
                </Button>
                <div className="mt-4">
                  <Button variant="link" onClick={() => router.push('/admin')} className="text-xs text-muted-foreground">
                    Admin portal
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Enhanced footer */}
          <div className="text-center mt-4 text-xs text-slate-500">
            <p>Secure • Reliable • Student-Focused</p>
            <div className="flex items-center justify-center space-x-3 mt-1">
              <div className="w-1.5 h-1.5 bg-green-500 rounded-full animate-pulse"></div>
              <span>All systems operational</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
