"use client"

import { memo, useMemo, useCallback } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import type { Room } from "@/lib/types"
import { Users, Bed, Wifi, Bath } from "lucide-react"

interface RoomCardProps {
  room: Room
  buildingName: string
  onViewRoom: (roomId: string) => void
}

export const RoomCard = memo(function RoomCard({ room, buildingName, onViewRoom }: RoomCardProps) {
  const roomTypeColor = useMemo(() => {
    switch (room.type) {
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
  }, [room.type])

  const getAmenityIcon = useCallback((amenity: string) => {
    switch (amenity.toLowerCase()) {
      case "wi-fi":
        return <Wifi className="h-3 w-3" />
      case "private bathroom":
      case "shared bathroom":
        return <Bath className="h-3 w-3" />
      default:
        return null
    }
  }, [])

  const handleClick = useCallback(() => {
    onViewRoom(room.id)
  }, [onViewRoom, room.id])

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">Room {room.number}</CardTitle>
          <Badge className={roomTypeColor}>{room.type}</Badge>
        </div>
        <p className="text-sm text-muted-foreground">{buildingName}</p>
      </CardHeader>

      <CardContent className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div className="flex items-center space-x-2 text-sm text-muted-foreground">
            <Users className="h-4 w-4" />
            <span>{room.totalBeds} beds</span>
          </div>
          <div className="flex items-center space-x-2 text-sm text-muted-foreground">
            <Bed className="h-4 w-4 text-green-600" />
            <span>{room.availableBeds} available</span>
          </div>
        </div>

        <div className="flex flex-wrap gap-1">
          {room.amenities.slice(0, 3).map((amenity) => (
            <div key={amenity} className="flex items-center space-x-1 text-xs text-muted-foreground">
              {getAmenityIcon(amenity)}
              <span>{amenity}</span>
            </div>
          ))}
          {room.amenities.length > 3 && (
            <span className="text-xs text-muted-foreground">+{room.amenities.length - 3} more</span>
          )}
        </div>

        <Button onClick={handleClick} className="w-full" size="sm">
          View Room Details
        </Button>
      </CardContent>
    </Card>
  )
})
