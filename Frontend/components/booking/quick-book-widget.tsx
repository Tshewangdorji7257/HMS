"use client"

import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { useHostelData } from "@/hooks/use-hostel-data"
import { useRouter } from "next/navigation"
import { MapPin, Zap, Building2, Bed } from "lucide-react"

export function QuickBookWidget() {
  const { data } = useHostelData()
  const [selectedBuilding, setSelectedBuilding] = useState("")
  const router = useRouter()

  const availableBuildings = data.buildings.filter((building) => building.availableBeds > 0)

  const handleQuickSearch = () => {
    if (selectedBuilding) {
      router.push(`/building/${selectedBuilding}`)
    }
  }

  const selectedBuildingData = availableBuildings.find((b) => b.id === selectedBuilding)

  return (
    <Card className="bg-gradient-to-br from-white via-rose-50/30 to-pink-50/20 border-0 shadow-xl overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-br from-rose-500/5 to-pink-500/5" />
      <CardHeader className="relative pb-4">
        <CardTitle className="flex items-center space-x-3">
          <div className="p-2 rounded-xl bg-gradient-to-br from-rose-100 to-pink-100">
            <Zap className="h-5 w-5 text-rose-600" />
          </div>
          <div>
            <span className="text-xl font-serif font-light text-slate-800">Quick Book</span>
            <p className="text-sm text-slate-500 font-normal mt-1">Find your room instantly</p>
          </div>
        </CardTitle>
      </CardHeader>

      <CardContent className="space-y-6 relative">
        <div className="space-y-3">
          <label className="text-sm font-medium text-slate-700 flex items-center">
            <Building2 className="h-4 w-4 mr-2 text-slate-500" />
            Select Building
          </label>
          <Select value={selectedBuilding} onValueChange={setSelectedBuilding}>
            <SelectTrigger className="h-12 rounded-xl border-slate-200 bg-white/80 backdrop-blur-sm hover:bg-white transition-all duration-300">
              <SelectValue placeholder="Choose your preferred building..." />
            </SelectTrigger>
            <SelectContent className="rounded-xl border-slate-200">
              {availableBuildings.map((building) => (
                <SelectItem key={building.id} value={building.id} className="rounded-lg">
                  <div className="flex items-center justify-between w-full">
                    <div className="flex items-center space-x-3">
                      <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                      <span className="font-medium">{building.name}</span>
                    </div>
                    <div className="flex items-center space-x-2 text-xs text-slate-500 ml-4">
                      <Bed className="h-3 w-3" />
                      <span>{building.availableBeds} available</span>
                    </div>
                  </div>
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {selectedBuildingData && (
          <div className="bg-white/60 backdrop-blur-sm rounded-xl p-4 border border-slate-200/50">
            <div className="grid grid-cols-2 gap-4 text-center">
              <div>
                <div className="text-2xl font-serif font-light text-slate-800">{selectedBuildingData.totalRooms}</div>
                <div className="text-xs text-slate-500">Total Rooms</div>
              </div>
              <div>
                <div className="text-2xl font-serif font-light text-green-600">
                  {selectedBuildingData.availableBeds}
                </div>
                <div className="text-xs text-slate-500">Available Beds</div>
              </div>
            </div>
          </div>
        )}

        <Button
          onClick={handleQuickSearch}
          disabled={!selectedBuilding}
          className="w-full h-12 bg-gradient-to-r from-rose-500 to-pink-500 hover:from-rose-600 hover:to-pink-600 text-white font-medium rounded-xl shadow-lg hover:shadow-xl transition-all duration-300 transform hover:scale-[1.02]"
        >
          <MapPin className="mr-2 h-4 w-4" />
          Explore Rooms
        </Button>

        <div className="text-center space-y-2">
          <div className="flex items-center justify-center space-x-2 text-sm text-slate-600">
            <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
            <span>{availableBuildings.length} buildings with available beds</span>
          </div>
          <p className="text-xs text-slate-500">
            {availableBuildings.reduce((sum, b) => sum + b.availableBeds, 0)} total beds ready for booking
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
