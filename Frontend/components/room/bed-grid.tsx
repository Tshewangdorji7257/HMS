"use client"

import { useState, useCallback } from "react"
import { Button } from "@/components/ui/button"
import { useHostelData } from "@/hooks/use-hostel-data"
import { useAuth } from "@/hooks/use-auth"
import type { Room, Building, Bed, Booking } from "@/lib/types"
import { BookingModal } from "./booking-modal"
import { BookingConfirmation } from "../booking/booking-confirmation"
import { cn } from "@/lib/utils"

interface BedGridProps {
  readonly room: Room
  readonly building: Building
}

export function BedGrid({ room, building }: BedGridProps) {
  const [selectedBed, setSelectedBed] = useState<Bed | null>(null)
  const [showBookingModal, setShowBookingModal] = useState(false)
  const [showConfirmation, setShowConfirmation] = useState(false)
  const [confirmedBooking, setConfirmedBooking] = useState<Booking | null>(null)
  const { bookBed } = useHostelData()
  const { user } = useAuth()

  const handleBedClick = useCallback((bed: Bed) => {
    if (bed.isOccupied) return
    setSelectedBed(bed)
    setShowBookingModal(true)
  }, [])

  const handleBooking = useCallback(async () => {
    if (!selectedBed || !user) return { success: false, error: "Invalid booking data" }

    const result = await bookBed(building.id, room.id, selectedBed.id)
    if (result.success && result.booking) {
      setConfirmedBooking(result.booking)
      setShowBookingModal(false)
      setSelectedBed(null)
      setShowConfirmation(true)
    }
    return result
  }, [selectedBed, user, bookBed, building.id, room.id])

  const getBedStatus = (bed: Bed) => {
    if (!bed.isOccupied) return "available"
    if (bed.occupiedBy === user?.id) return "mine"
    return "occupied"
  }

  const getBedColor = (bed: Bed) => {
    const status = getBedStatus(bed)
    switch (status) {
      case "available":
        return "bg-green-100 border-green-300 hover:bg-green-200 text-green-800 hover:shadow-md"
      case "mine":
        return "bg-blue-100 border-blue-300 text-blue-800"
      case "occupied":
        return "bg-red-100 border-red-300 text-red-800"
      default:
        return "bg-gray-100 border-gray-300 text-gray-800"
    }
  }

  const getBedLabel = (bed: Bed) => {
    const status = getBedStatus(bed)
    switch (status) {
      case "available":
        return "Available"
      case "mine":
        return "Your Bed"
      case "occupied":
        return bed.occupiedByName || "Occupied"
      default:
        return "Unknown"
    }
  }

  // Create a grid layout based on room type
  const getGridLayout = () => {
    switch (room.type) {
      case "single":
        return "grid-cols-1"
      case "double":
        return "grid-cols-2"
      case "triple":
        return "grid-cols-3"
      case "quad":
        return "grid-cols-2"
      default:
        return "grid-cols-2"
    }
  }

  return (
    <>
      <div className="space-y-6">
        {/* Legend */}
        <div className="flex flex-wrap gap-4 text-sm">
          <div className="flex items-center space-x-2">
            <div className="w-4 h-4 bg-green-100 border border-green-300 rounded" />
            <span>Available</span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="w-4 h-4 bg-blue-100 border border-blue-300 rounded" />
            <span>Your Bed</span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="w-4 h-4 bg-red-100 border border-red-300 rounded" />
            <span>Occupied</span>
          </div>
        </div>

        {/* Room Layout Visualization */}
        <div className="bg-muted/20 rounded-lg p-6">
          <div className="text-center mb-4">
            <h3 className="font-medium text-foreground">Room Layout</h3>
            <p className="text-sm text-muted-foreground">Click on an available bed to book it</p>
          </div>

          {/* Bed Grid */}
          <div className={cn("grid gap-4 max-w-md mx-auto", getGridLayout())}>
            {room.beds.map((bed) => (
              <Button
                key={bed.id}
                variant="outline"
                className={cn(
                  "h-24 flex flex-col items-center justify-center space-y-1 transition-all duration-200 border-2",
                  getBedColor(bed),
                  bed.isOccupied ? "cursor-not-allowed opacity-75" : "cursor-pointer transform hover:scale-105",
                )}
                onClick={() => handleBedClick(bed)}
                disabled={bed.isOccupied}
              >
                <span className="font-medium">Bed {bed.number}</span>
                <span className="text-xs">{getBedLabel(bed)}</span>
              </Button>
            ))}
          </div>

          {/* Room Features */}
          <div className="mt-6 flex justify-center">
            <div className="text-xs text-muted-foreground bg-background/80 px-3 py-1 rounded-full">
              {room.type.charAt(0).toUpperCase() + room.type.slice(1)} Room • {room.totalBeds} Beds
            </div>
          </div>
        </div>

        {room.availableBeds === 0 && (
          <div className="text-center py-6 bg-muted/10 rounded-lg">
            <p className="text-muted-foreground">All beds in this room are currently occupied.</p>
            <p className="text-sm text-muted-foreground mt-1">Check back later or browse other rooms.</p>
          </div>
        )}
      </div>

      <BookingModal
        isOpen={showBookingModal}
        onClose={() => setShowBookingModal(false)}
        bed={selectedBed}
        room={room}
        building={building}
        onConfirm={handleBooking}
      />

      <BookingConfirmation
        isOpen={showConfirmation}
        onClose={() => setShowConfirmation(false)}
        booking={confirmedBooking}
      />
    </>
  )
}
