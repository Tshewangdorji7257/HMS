"use client"

import { Card, CardContent } from "@/components/ui/card"
import type { Booking } from "@/lib/types"
import { Calendar, CheckCircle, XCircle } from "lucide-react"

interface BookingStatsProps {
  bookings: Booking[]
  activeBooking: Booking | null
}

export function BookingStats({ bookings, activeBooking }: BookingStatsProps) {
  const totalBookings = bookings.length
  const activeBookings = bookings.filter((booking) => booking.status === "active").length
  const cancelledBookings = bookings.filter((booking) => booking.status === "cancelled").length

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center justify-between">
            <div className="space-y-2">
              <p className="text-sm font-medium text-muted-foreground">Total Bookings</p>
              <p className="text-2xl font-bold text-foreground">{totalBookings}</p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-primary/10 flex items-center justify-center">
              <Calendar className="h-6 w-6 text-primary" />
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-6">
          <div className="flex items-center justify-between">
            <div className="space-y-2">
              <p className="text-sm font-medium text-muted-foreground">Active Bookings</p>
              <p className="text-2xl font-bold text-foreground">{activeBookings}</p>
              {activeBooking && (
                <p className="text-xs text-muted-foreground">
                  {activeBooking.buildingName} • Room {activeBooking.roomNumber}
                </p>
              )}
            </div>
            <div className="h-12 w-12 rounded-lg bg-green-100 flex items-center justify-center">
              <CheckCircle className="h-6 w-6 text-green-600" />
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-6">
          <div className="flex items-center justify-between">
            <div className="space-y-2">
              <p className="text-sm font-medium text-muted-foreground">Cancelled</p>
              <p className="text-2xl font-bold text-foreground">{cancelledBookings}</p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-red-100 flex items-center justify-center">
              <XCircle className="h-6 w-6 text-red-600" />
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
