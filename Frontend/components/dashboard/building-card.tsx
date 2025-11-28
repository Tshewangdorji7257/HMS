"use client"

import { memo, useMemo, useCallback } from "react"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import type { Building } from "@/lib/types"
import { Users, Bed, MapPin, Star } from "lucide-react"
import Image from "next/image"

interface BuildingCardProps {
  building: Building
  onViewDetails: (buildingId: string) => void
}

export const BuildingCard = memo(function BuildingCard(props: Readonly<BuildingCardProps>) {
  const { building, onViewDetails } = props

  const occupancyRate = useMemo(() => 
    Math.abs(Math.round(((building.totalBeds - building.availableBeds) / building.totalBeds) * 100)),
    [building.totalBeds, building.availableBeds]
  )

  const handleClick = useCallback(() => {
    onViewDetails(building.id)
  }, [onViewDetails, building.id])

  const getOccupancyColor = (rate: number) => {
    if (rate < 50) return "text-emerald-700 bg-emerald-50 border-emerald-200"
    if (rate < 80) return "text-amber-700 bg-amber-50 border-amber-200"
    return "text-rose-700 bg-rose-50 border-rose-200"
  }

  // Amenities are intentionally hidden on the building card per request.

  return (
    <Card className="group hover:shadow-2xl transition-all duration-500 hover:-translate-y-2 border-0 bg-white/80 backdrop-blur-sm overflow-hidden flex flex-col h-full">
      <div className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-black/20 to-transparent z-10" />
        <Image
          src={
            building.image ||
            (() => {
              const query = encodeURIComponent(`Modern ${building.name} hostel building exterior`)
              return `/placeholder.svg?height=240&width=400&query=${query}`
            })()
          }
          alt={building.name}
          width={400}
          height={240}
          className="w-full h-60 object-cover group-hover:scale-110 transition-transform duration-700"
        />

        {/* Floating badges */}
        <div className="absolute top-4 right-4 z-20 flex flex-col gap-2">
          <Badge className={`${getOccupancyColor(occupancyRate)} border font-medium shadow-lg backdrop-blur-sm`}>
            {occupancyRate}% occupied
          </Badge>
          {building.availableBeds > 0 && (
            <Badge className="bg-green-500/90 text-white border-0 shadow-lg backdrop-blur-sm">
              {building.availableBeds} beds available
            </Badge>
          )}
        </div>

        {/* Building name overlay */}
        <div className="absolute bottom-4 left-4 right-4 z-20">
          <h3 className="text-2xl font-serif font-light text-white mb-1 drop-shadow-lg">{building.name}</h3>
          <div className="flex items-center text-white/90 text-sm">
            <MapPin className="h-4 w-4 mr-1" />
            <span>Campus District</span>
          </div>
        </div>
      </div>

      <CardContent className="p-6 flex flex-col flex-1">
        <p className="text-slate-600 leading-relaxed line-clamp-2 mb-6 min-h-[44px]">{building.description}</p>
        <div className="grid grid-cols-2 gap-4 mb-6">
          <div className="bg-slate-50 rounded-xl p-4 text-center flex flex-col items-center justify-between">
            <div className="flex items-center justify-center mb-2">
              <div className="p-2 rounded-lg bg-blue-100">
                <Users className="h-5 w-5 text-blue-600" />
              </div>
            </div>
            <div className="text-2xl font-serif font-light text-slate-800">{building.totalRooms}</div>
            <div className="text-sm text-slate-500">Total Rooms</div>
          </div>
          <div className="bg-slate-50 rounded-xl p-4 text-center flex flex-col items-center justify-between">
            <div className="flex items-center justify-center mb-2">
              <div className="p-2 rounded-lg bg-green-100">
                <Bed className="h-5 w-5 text-green-600" />
              </div>
            </div>
            <div className="text-2xl font-serif font-light text-slate-800">{building.availableBeds}</div>
            <div className="text-sm text-slate-500">Available</div>
          </div>
        </div>
        <div className="mt-auto">
          <Button
            onClick={handleClick}
            className="w-full h-12 bg-gradient-to-r from-rose-500 to-pink-500 hover:from-rose-600 hover:to-pink-600 text-white font-medium rounded-xl shadow-lg hover:shadow-xl transition-all duration-300 transform hover:scale-[1.02]"
            disabled={building.availableBeds === 0}
          >
            {building.availableBeds === 0 ? (
              <>
                <Bed className="h-4 w-4 mr-2" />
                Fully Occupied
              </>
            ) : (
              <>
                <Star className="h-4 w-4 mr-2" />
                Explore Rooms
              </>
            )}
          </Button>
        </div>
      </CardContent>
    </Card>
  )
})
