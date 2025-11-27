"use client"

import { Button } from "@/components/ui/button"
import { useAuth } from "@/hooks/use-auth"
import { LogOut, Home, Calendar, Building2 } from "lucide-react"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { useRouter, usePathname } from "next/navigation"

export function Header() {
  const { user, logout } = useAuth()
  const router = useRouter()
  const pathname = usePathname()

  // Always render header; when no user show basic nav and admin access

  const isActive = (path: string) => pathname === path

  const handleLogout = () => {
    // clear auth and redirect to auth page
    logout()
    router.push('/auth')
  }

  return (
    <header className="border-b border-border/30 bg-card/80 backdrop-blur-xl sticky top-0 z-50 shadow-sm">
      <div className="container mx-auto px-6 py-5">
        <div className="flex items-center justify-between">
          {/* Left Side Logo */}
          <div className="flex items-center space-x-4 cursor-pointer group" onClick={() => router.push("/")}>
            <div className="p-2 rounded-xl bg-primary/10 group-hover:bg-primary/20 transition-colors duration-300">
              <Building2 className="h-6 w-6 text-primary" />
            </div>
            <h1 className="text-2xl font-serif font-light text-foreground tracking-tight">HMS</h1>
          </div>

          {/* Right Side Nav */}
          <div className="flex items-center space-x-4">
            <nav className="hidden md:flex items-center space-x-2">
              <Button
                variant="ghost"
                className={`px-6 py-2 rounded-full transition-all duration-300 ${
                  isActive("/")
                    ? "text-primary bg-primary/15 shadow-sm"
                    : "text-muted-foreground hover:text-foreground hover:bg-accent/50"
                }`}
                onClick={() => router.push("/")}
              >
                <Home className="h-4 w-4 mr-2" />
                Dashboard
              </Button>
              {!user && (
                <Button variant="ghost" className="px-6 py-2 rounded-full" onClick={() => router.push('/admin')}>
                  Admin
                </Button>
              )}
              <Button
                variant="ghost"
                className={`px-6 py-2 rounded-full transition-all duration-300 ${
                  isActive("/bookings")
                    ? "text-primary bg-primary/15 shadow-sm"
                    : "text-muted-foreground hover:text-foreground hover:bg-accent/50"
                }`}
                onClick={() => router.push("/bookings")}
              >
                <Calendar className="h-4 w-4 mr-2" />
                My Bookings
              </Button>
            </nav>

            {/* User Menu */}
            {user ? (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    className="relative h-12 w-12 rounded-full hover:bg-accent/50 transition-colors duration-300"
                  >
                    <Avatar className="h-10 w-10 ring-2 ring-primary/20">
                      <AvatarFallback className="bg-gradient-to-br from-primary to-primary/80 text-primary-foreground font-medium">
                        {user.name.charAt(0).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-64 p-2" align="end" forceMount>
                  <div className="flex items-center justify-start gap-3 p-3 rounded-lg bg-accent/30">
                    <Avatar className="h-10 w-10">
                      <AvatarFallback className="bg-gradient-to-br from-primary to-primary/80 text-primary-foreground">
                        {user.name.charAt(0).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                    <div className="flex flex-col space-y-1 leading-none">
                      <p className="font-medium text-foreground">{user.name}</p>
                      <p className="w-[180px] truncate text-sm text-muted-foreground">{user.email}</p>
                    </div>
                  </div>
                  <DropdownMenuSeparator className="my-2" />
                  <DropdownMenuItem onClick={() => router.push("/")} className="rounded-lg">
                    <Home className="mr-3 h-4 w-4" />
                    Dashboard
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => router.push("/bookings")} className="rounded-lg">
                    <Calendar className="mr-3 h-4 w-4" />
                    My Bookings
                  </DropdownMenuItem>
                  <DropdownMenuSeparator className="my-2" />
                  <DropdownMenuItem onClick={handleLogout} className="text-destructive rounded-lg">
                    <LogOut className="mr-3 h-4 w-4" />
                    Sign out
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            ) : (
              <div className="flex items-center gap-2">
                <Button variant="ghost" onClick={() => router.push('/auth')}>
                  Sign In
                </Button>
              </div>
            )}
            
          </div>
        </div>
      </div>
    </header>
  )
}
