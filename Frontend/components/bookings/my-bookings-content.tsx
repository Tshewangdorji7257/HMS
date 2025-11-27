"use client"

import { useState } from "react"
import { useHostelData } from "@/hooks/use-hostel-data"
import { useAuth } from "@/hooks/use-auth"
import { BookingCard } from "./booking-card"
import { BookingStats } from "./booking-stats"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Calendar, Home } from "lucide-react"
import { Skeleton } from "@/components/ui/skeleton"
import { useRouter } from "next/navigation"

export function MyBookingsContent() {
  const { data, loading } = useHostelData()
  const { user } = useAuth()
  const router = useRouter()
  const [activeTab, setActiveTab] = useState("active")

  const allBookings = data.bookings.filter((booking) => booking.userId === user?.id)
  const activeBookings = allBookings.filter((booking) => booking.status === "active")
  const cancelledBookings = allBookings.filter((booking) => booking.status === "cancelled")
  const activeBooking = allBookings.find((booking) => booking.status === "active") || null

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-24" />
          ))}
        </div>
        <Skeleton className="h-96" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Booking Stats */}
      <BookingStats bookings={allBookings} activeBooking={activeBooking} />

      {/* Quick Actions */}
      <div className="flex flex-col sm:flex-row gap-4">
        <Button onClick={() => router.push("/")} className="flex items-center space-x-2">
          <Home className="h-4 w-4" />
          <span>Browse Buildings</span>
        </Button>
      </div>

      {/* Bookings List */}
      {allBookings.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <Calendar className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-foreground mb-2">No bookings yet</h3>
            <p className="text-muted-foreground mb-4">You haven't made any room reservations.</p>
            <Button onClick={() => router.push("/")} className="flex items-center space-x-2">
              <Home className="h-4 w-4" />
              <span>Browse Available Rooms</span>
            </Button>
          </CardContent>
        </Card>
      ) : (
        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="active">Active ({activeBookings.length})</TabsTrigger>
            <TabsTrigger value="history">History ({cancelledBookings.length})</TabsTrigger>
          </TabsList>

          <TabsContent value="active" className="space-y-4">
            {activeBookings.length === 0 ? (
              <Card>
                <CardContent className="py-8 text-center">
                  <p className="text-muted-foreground">No active bookings</p>
                </CardContent>
              </Card>
            ) : (
              activeBookings.map((booking) => <BookingCard key={booking.id} booking={booking} />)
            )}
          </TabsContent>

          <TabsContent value="history" className="space-y-4">
            {cancelledBookings.length === 0 ? (
              <Card>
                <CardContent className="py-8 text-center">
                  <p className="text-muted-foreground">No booking history</p>
                </CardContent>
              </Card>
            ) : (
              cancelledBookings.map((booking) => <BookingCard key={booking.id} booking={booking} />)
            )}
          </TabsContent>
        </Tabs>
      )}
    </div>
  )
}
