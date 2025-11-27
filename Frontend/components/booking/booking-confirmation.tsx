"use client"

import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { CheckCircle, MapPin, User, Bed } from "lucide-react"
import type { Room, Bed as BedType, Booking } from "@/lib/types"

interface BookingConfirmationProps {
  readonly isOpen: boolean
  readonly onClose: () => void
  // Either provide a booking (for post-booking confirmation) or room+bed details (pre-booking)
  readonly booking?: Booking | null
  readonly room?: Room | null
  readonly bed?: BedType | null
  readonly buildingName?: string
  // onConfirm is optional; if provided the dialog will show a Confirm button
  readonly onConfirm?: (() => void) | null
}

export function BookingConfirmation({ isOpen, onClose, booking, room, bed, buildingName, onConfirm }: BookingConfirmationProps) {
  // If neither booking nor room+bed are provided, nothing to show
  if (!booking && (!room || !bed)) return null

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <CheckCircle className="h-5 w-5 text-green-600" />
            <span>Confirm Booking</span>
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-6">
          <div className="bg-slate-50 rounded-lg p-4 space-y-3">
            <div className="flex items-center space-x-2">
              <MapPin className="h-4 w-4 text-slate-500" />
              <span className="font-medium">{booking?.buildingName ?? buildingName}</span>
            </div>
            <div className="flex items-center space-x-2">
              <Bed className="h-4 w-4 text-slate-500" />
              <span>
                Room {booking?.roomNumber ?? room?.number} - Bed {booking?.bedNumber ?? bed?.number}
              </span>
            </div>
            <div className="flex items-center space-x-2">
              <User className="h-4 w-4 text-slate-500" />
              <span>{room?.type ?? booking?.status ?? ""}</span>
              {room?.price !== undefined && (
                <Badge variant="secondary" className="ml-2">
                  ₹{room.price}/month
                </Badge>
              )}
            </div>
          </div>

          <div className="text-sm text-slate-600">
            <p>Are you sure you want to book this bed? This action cannot be undone.</p>
          </div>

          <div className="flex space-x-3">
            <Button variant="outline" onClick={onClose} className="flex-1 bg-transparent">
              Close
            </Button>
            {onConfirm && (
              <Button onClick={onConfirm} className="flex-1">
                Confirm Booking
              </Button>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
