"use client"

import { useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import type { Bed, Room, Building } from "@/lib/types"
import { Loader2, MapPin, Users, Calendar, CheckCircle } from "lucide-react"
import { useToast } from "@/hooks/use-toast"

interface BookingModalProps {
  isOpen: boolean
  onClose: () => void
  bed: Bed | null
  room: Room
  building: Building
  onConfirm: () => Promise<{ success: boolean; error?: string }>
}

export function BookingModal({ isOpen, onClose, bed, room, building, onConfirm }: BookingModalProps) {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  const { toast } = useToast()

  if (!bed) return null

  const handleConfirm = async () => {
    setLoading(true)
    setError("")

    try {
      const result = await onConfirm()
      if (result.success) {
        toast({
          title: "Booking Confirmed!",
          description: `You've successfully booked Bed ${bed.number} in Room ${room.number}, ${building.name}.`,
          variant: "default",
        })
        onClose()
      } else {
        setError(result.error || "Booking failed")
      }
    } catch (err) {
      setError("An unexpected error occurred")
    } finally {
      setLoading(false)
    }
  }

  const getRoomTypeColor = (type: string) => {
    switch (type) {
      case "single":
        return "bg-blue-100 text-blue-800"
      case "double":
        return "bg-green-100 text-green-800"
      case "triple":
        return "bg-yellow-100 text-yellow-800"
      case "quad":
        return "bg-purple-100 text-purple-800"
      default:
        return "bg-gray-100 text-gray-800"
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <CheckCircle className="h-5 w-5 text-primary" />
            <span>Confirm Booking</span>
          </DialogTitle>
          <DialogDescription>Please review your booking details before confirming.</DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {/* Booking Details */}
          <div className="bg-muted/50 rounded-lg p-4 space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Building</span>
              <span className="font-medium">{building.name}</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Room</span>
              <div className="flex items-center space-x-2">
                <span className="font-medium">Room {room.number}</span>
                <Badge className={getRoomTypeColor(room.type)}>{room.type}</Badge>
              </div>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Bed</span>
              <span className="font-medium">Bed {bed.number}</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Booking Date</span>
              <span className="font-medium">{new Date().toLocaleDateString()}</span>
            </div>
          </div>

          {/* Room Info */}
          <div className="grid grid-cols-3 gap-4 text-center">
            <div className="flex flex-col items-center space-y-1">
              <MapPin className="h-4 w-4 text-muted-foreground" />
              <span className="text-xs text-muted-foreground">Building</span>
              <span className="text-sm font-medium">{building.name}</span>
            </div>
            <div className="flex flex-col items-center space-y-1">
              <Users className="h-4 w-4 text-muted-foreground" />
              <span className="text-xs text-muted-foreground">Room Type</span>
              <span className="text-sm font-medium capitalize">{room.type}</span>
            </div>
            <div className="flex flex-col items-center space-y-1">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              <span className="text-xs text-muted-foreground">Available</span>
              <span className="text-sm font-medium">Now</span>
            </div>
          </div>

          {error && <div className="text-sm text-destructive bg-destructive/10 p-3 rounded-md">{error}</div>}

          <div className="text-xs text-muted-foreground bg-muted/30 p-3 rounded-md">
            <p className="font-medium mb-1">Important:</p>
            <ul className="space-y-1">
              <li>• You can only have one active booking at a time</li>
              <li>• You can cancel your booking anytime from "My Bookings"</li>
              <li>• This booking is free of charge</li>
            </ul>
          </div>
        </div>

        <DialogFooter className="flex-col sm:flex-row gap-2">
          <Button variant="outline" onClick={onClose} disabled={loading} className="w-full sm:w-auto bg-transparent">
            Cancel
          </Button>
          <Button onClick={handleConfirm} disabled={loading} className="w-full sm:w-auto">
            {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Confirm Booking
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
