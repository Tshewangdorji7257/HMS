"use client"

import { useParams, useRouter } from "next/navigation"
import { useState, useMemo, useCallback } from "react"
import { useHostelData } from "@/hooks/use-hostel-data"
import { Header } from "@/components/layout/header"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { ArrowLeft, MapPin, Users, Bed, Wifi, Car, Search } from "lucide-react"
import { Skeleton } from "@/components/ui/skeleton"
import Image from "next/image"
import { RoomCard } from "@/components/building/room-card"

export default function BuildingPage() {
  const params = useParams()
  const router = useRouter()
  const { data, loading } = useHostelData()
  const [roomSearch, setRoomSearch] = useState("")

  const buildingId = params.buildingId as string
  const building = useMemo(() => data.buildings.find((b) => b.id === buildingId), [data.buildings, buildingId])

  // Filter rooms based on search with memoization
  const filteredRooms = useMemo(() => {
    if (!building) return []
    
    return building.rooms.filter((room) => {
      if (!roomSearch.trim()) return room.availableBeds > 0
      const query = roomSearch.toLowerCase()
      return (
        room.availableBeds > 0 &&
        (room.number.toLowerCase().includes(query) ||
          room.type.toLowerCase().includes(query) ||
          room.amenities.some((amenity) => amenity.toLowerCase().includes(query)))
      )
    })
  }, [building, roomSearch])

  const handleBack = useCallback(() => {
    router.push("/")
  }, [router])

  const handleViewRoom = useCallback((roomId: string) => {
    router.push(`/building/${buildingId}/room/${roomId}`)
  }, [router, buildingId])

  if (loading) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <main className="container mx-auto px-4 py-8">
          <Skeleton className="h-8 w-64 mb-6" />
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2">
              <Skeleton className="h-64 w-full mb-6" />
              <Skeleton className="h-32 w-full" />
            </div>
            <div>
              <Skeleton className="h-48 w-full" />
            </div>
          </div>
        </main>
      </div>
    )
  }

  if (!building) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <main className="container mx-auto px-4 py-8">
            <div className="text-center py-12">
            <h1 className="text-2xl font-bold text-foreground mb-2">Building Not Found</h1>
            <p className="text-muted-foreground mb-4">The building you're looking for doesn't exist.</p>
            <Button onClick={handleBack} variant="outline">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Dashboard
            </Button>
          </div>
        </main>
      </div>
    )
  }

  const occupancyRate = Math.abs(Math.round(((building.totalBeds - building.availableBeds) / building.totalBeds) * 100))

  const getAmenityIcon = (amenity: string) => {
    switch (amenity.toLowerCase()) {
      case "wi-fi":
        return <Wifi className="h-4 w-4" />
      case "parking":
        return <Car className="h-4 w-4" />
      default:
        return <MapPin className="h-4 w-4" />
    }
  }

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="container mx-auto px-4 py-8">
        {/* Back Button */}
        <Button onClick={handleBack} variant="ghost" className="mb-6">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Dashboard
        </Button>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Content */}
          <div className="lg:col-span-2 space-y-6">
            {/* Building Header */}
            <div>
              <div className="relative overflow-hidden rounded-lg mb-4">
                <Image
                  src={building.image || "/placeholder.svg"}
                  alt={building.name}
                  width={800}
                  height={400}
                  className="w-full h-64 object-cover"
                />
                <div className="absolute top-4 right-4">
                  <Badge className="bg-primary text-primary-foreground">{occupancyRate}% occupied</Badge>
                </div>
              </div>

              <div className="flex items-start justify-between">
                <div>
                  <h1 className="text-3xl font-bold text-foreground mb-2">{building.name}</h1>
                  <p className="text-muted-foreground text-lg">{building.description}</p>
                </div>
              </div>
            </div>

            {/* Building Stats */}
            <div className="grid grid-cols-3 gap-4">
              <Card>
                <CardContent className="p-4 text-center">
                  <Users className="h-6 w-6 text-primary mx-auto mb-2" />
                  <p className="text-2xl font-bold text-foreground">{building.totalRooms}</p>
                  <p className="text-sm text-muted-foreground">Total Rooms</p>
                </CardContent>
              </Card>
              <Card>
                <CardContent className="p-4 text-center">
                  <Bed className="h-6 w-6 text-primary mx-auto mb-2" />
                  <p className="text-2xl font-bold text-foreground">{building.totalBeds}</p>
                  <p className="text-sm text-muted-foreground">Total Beds</p>
                </CardContent>
              </Card>
              <Card>
                <CardContent className="p-4 text-center">
                  <Bed className="h-6 w-6 text-green-600 mx-auto mb-2" />
                  <p className="text-2xl font-bold text-foreground">{building.availableBeds}</p>
                  <p className="text-sm text-muted-foreground">Available</p>
                </CardContent>
              </Card>
            </div>

            {/* Room Search */}
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <h2 className="text-2xl font-bold text-foreground">Available Rooms</h2>
              </div>

              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  placeholder="Search rooms by number, type, or amenities..."
                  value={roomSearch}
                  onChange={(e) => setRoomSearch(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>

            {/* Rooms Grid */}
            <div>
              {filteredRooms.length === 0 ? (
                <div className="text-center py-8">
                  <p className="text-muted-foreground">
                    {roomSearch.trim()
                      ? `No rooms found matching "${roomSearch}"`
                      : "No available rooms in this building."}
                  </p>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {filteredRooms.map((room) => (
                    <RoomCard
                      key={room.id}
                      room={room}
                      buildingName={building.name}
                      onViewRoom={handleViewRoom}
                    />
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Amenities */}
            <Card>
              <CardHeader>
                <CardTitle>Building Amenities</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {building.amenities.map((amenity) => (
                    <div key={amenity} className="flex items-center space-x-3">
                      {getAmenityIcon(amenity)}
                      <span className="text-sm text-foreground">{amenity}</span>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Quick Actions */}
            <Card>
              <CardHeader>
                <CardTitle>Quick Actions</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <Button className="w-full" disabled={building.availableBeds === 0}>
                  {building.availableBeds === 0 ? "Fully Occupied" : "Find Available Room"}
                </Button>
                <Button variant="outline" className="w-full bg-transparent">
                  View Floor Plan
                </Button>
                <Button variant="outline" className="w-full bg-transparent">
                  Contact Manager
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>
      </main>
    </div>
  )
}
