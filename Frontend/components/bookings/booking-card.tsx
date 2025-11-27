"use client"

import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { useHostelData } from "@/hooks/use-hostel-data"
import { useToast } from "@/hooks/use-toast"
import type { Booking } from "@/lib/types"
import { MapPin, Calendar, Bed, X, AlertTriangle } from "lucide-react"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"

interface BookingCardProps {
  booking: Booking
}

export function BookingCard({ booking }: BookingCardProps) {
  const [loading, setLoading] = useState(false)
  const { cancelBooking } = useHostelData()
  const { toast } = useToast()

  const handleCancelBooking = async () => {
    setLoading(true)
    try {
      const result = await cancelBooking(booking.id)
      if (result.success) {
        toast({
          title: "Booking Cancelled",
          description: `Your booking for Bed ${booking.bedNumber} in Room ${booking.roomNumber} has been cancelled.`,
          variant: "default",
        })
      } else {
        toast({
          title: "Cancellation Failed",
          description: result.error || "Failed to cancel booking",
          variant: "destructive",
        })
      }
    } catch (error) {
      toast({
        title: "Error",
        description: "An unexpected error occurred",
        variant: "destructive",
      })
    } finally {
      setLoading(false)
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case "active":
        return "bg-green-100 text-green-800"
      case "cancelled":
        return "bg-red-100 text-red-800"
      default:
        return "bg-gray-100 text-gray-800"
    }
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    })
  }

  return (
    <Card className={`${booking.status === "active" ? "border-primary/20 bg-primary/5" : ""}`}>
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div>
            <CardTitle className="text-lg">{booking.buildingName}</CardTitle>
            <p className="text-sm text-muted-foreground">
              Room {booking.roomNumber} • Bed {booking.bedNumber}
            </p>
          </div>
          <Badge className={getStatusColor(booking.status)}>{booking.status}</Badge>
        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div className="flex items-center space-x-2 text-sm text-muted-foreground">
            <MapPin className="h-4 w-4" />
            <span>{booking.buildingName}</span>
          </div>
          <div className="flex items-center space-x-2 text-sm text-muted-foreground">
            <Bed className="h-4 w-4" />
            <span>Room {booking.roomNumber}</span>
          </div>
          <div className="flex items-center space-x-2 text-sm text-muted-foreground">
            <Calendar className="h-4 w-4" />
            <span>{formatDate(booking.bookingDate)}</span>
          </div>
        </div>

        {booking.status === "active" && (
          <div className="flex items-center justify-between pt-2 border-t">
            <div className="text-sm text-muted-foreground">
              <p>Booking ID: {booking.id.slice(-8)}</p>
            </div>
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button variant="destructive" size="sm" disabled={loading}>
                  <X className="h-4 w-4 mr-1" />
                  Cancel Booking
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle className="flex items-center space-x-2">
                    <AlertTriangle className="h-5 w-5 text-destructive" />
                    <span>Cancel Booking</span>
                  </AlertDialogTitle>
                  <AlertDialogDescription>
                    Are you sure you want to cancel your booking for Bed {booking.bedNumber} in Room{" "}
                    {booking.roomNumber} at {booking.buildingName}? This action cannot be undone.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Keep Booking</AlertDialogCancel>
                  <AlertDialogAction onClick={handleCancelBooking} className="bg-destructive hover:bg-destructive/90">
                    Cancel Booking
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </div>
        )}

        {booking.status === "cancelled" && (
          <div className="text-xs text-muted-foreground pt-2 border-t">
            <p>Cancelled on {formatDate(booking.bookingDate)}</p>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
