"use client"

import { useParams, useRouter } from "next/navigation"
import { useHostelData } from "@/hooks/use-hostel-data"
import { Header } from "@/components/layout/header"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { ArrowLeft, Users, Bed } from "lucide-react"
import { Skeleton } from "@/components/ui/skeleton"
import { BedGrid } from "@/components/room/bed-grid"

export default function RoomPage() {
  const params = useParams()
  const router = useRouter()
  const { data, loading } = useHostelData()

  const buildingId = params.buildingId as string
  const roomId = params.roomId as string

  const building = data.buildings.find((b) => b.id === buildingId)
  const room = building?.rooms.find((r) => r.id === roomId)

  if (loading) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <main className="container mx-auto px-4 py-8">
          <Skeleton className="h-8 w-64 mb-6" />
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2">
              <Skeleton className="h-64 w-full mb-6" />
            </div>
            <div>
              <Skeleton className="h-48 w-full" />
            </div>
          </div>
        </main>
      </div>
    )
  }

  if (!building || !room) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <main className="container mx-auto px-4 py-8">
          <div className="text-center py-12">
            <h1 className="text-2xl font-bold text-foreground mb-2">Room Not Found</h1>
            <p className="text-muted-foreground mb-4">The room you're looking for doesn't exist.</p>
            <Button onClick={() => router.push("/")} variant="outline">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Dashboard
            </Button>
          </div>
        </main>
      </div>
    )
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
    <div className="min-h-screen bg-background">
      <Header />
      <main className="container mx-auto px-4 py-8">
        {/* Back Button */}
        <Button onClick={() => router.push(`/building/${buildingId}`)} variant="ghost" className="mb-6">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to {building.name}
        </Button>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Content */}
          <div className="lg:col-span-2 space-y-6">
            {/* Room Header */}
            <div>
              <div className="flex items-center justify-between mb-4">
                <div>
                  <h1 className="text-3xl font-bold text-foreground">Room {room.number}</h1>
                  <p className="text-muted-foreground text-lg">{building.name}</p>
                </div>
                <Badge className={getRoomTypeColor(room.type)}>{room.type} room</Badge>
              </div>
            </div>

            {/* Room Stats */}
            <div className="grid grid-cols-2 gap-4">
              <Card>
                <CardContent className="p-4 text-center">
                  <Users className="h-6 w-6 text-primary mx-auto mb-2" />
                  <p className="text-2xl font-bold text-foreground">{room.totalBeds}</p>
                  <p className="text-sm text-muted-foreground">Total Beds</p>
                </CardContent>
              </Card>
              <Card>
                <CardContent className="p-4 text-center">
                  <Bed className="h-6 w-6 text-green-600 mx-auto mb-2" />
                  <p className="text-2xl font-bold text-foreground">{room.availableBeds}</p>
                  <p className="text-sm text-muted-foreground">Available</p>
                </CardContent>
              </Card>
            </div>

            {/* Bed Layout */}
            <Card>
              <CardHeader>
                <CardTitle>Room Layout</CardTitle>
                <p className="text-sm text-muted-foreground">Click on an available bed to book it</p>
              </CardHeader>
              <CardContent>
                <BedGrid room={room} building={building} />
              </CardContent>
            </Card>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Room Amenities */}
            <Card>
              <CardHeader>
                <CardTitle>Room Amenities</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {room.amenities.map((amenity) => (
                    <div key={amenity} className="flex items-center space-x-2">
                      <div className="w-2 h-2 bg-primary rounded-full" />
                      <span className="text-sm text-foreground">{amenity}</span>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Booking Info */}
            <Card>
              <CardHeader>
                <CardTitle>Booking Information</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="text-sm">
                  <p className="text-muted-foreground">Room Type</p>
                  <p className="font-medium text-foreground capitalize">{room.type}</p>
                </div>
                <div className="text-sm">
                  <p className="text-muted-foreground">Available Beds</p>
                  <p className="font-medium text-foreground">
                    {room.availableBeds} out of {room.totalBeds}
                  </p>
                </div>
                <div className="text-sm">
                  <p className="text-muted-foreground">Building</p>
                  <p className="font-medium text-foreground">{building.name}</p>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </main>
    </div>
  )
}
