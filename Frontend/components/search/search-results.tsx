"use client"

import { memo } from "react"
import { BuildingCard } from "@/components/dashboard/building-card"
import type { Building } from "@/lib/types"
import { Card, CardContent } from "@/components/ui/card"
import { Search, Filter, MapPin, Building2 } from "lucide-react"

interface SearchResultsProps {
  buildings: Building[]
  searchQuery: string
  onViewDetails: (buildingId: string) => void
}

export const SearchResults = memo(function SearchResults({ buildings, searchQuery, onViewDetails }: SearchResultsProps) {
  if (buildings.length === 0) {
    return (
      <Card className="border-0 bg-white/60 backdrop-blur-sm shadow-lg">
        <CardContent className="py-16 text-center">
          <div className="relative mb-8">
            <div className="w-24 h-24 bg-gradient-to-br from-slate-100 to-slate-200 rounded-full flex items-center justify-center mx-auto">
              <Search className="h-10 w-10 text-slate-400" />
            </div>
            <div className="absolute -bottom-2 -right-2 w-8 h-8 bg-gradient-to-br from-rose-100 to-pink-100 rounded-full flex items-center justify-center">
              <Filter className="h-4 w-4 text-rose-500" />
            </div>
          </div>

          <h3 className="text-2xl font-serif font-light text-slate-800 mb-3">No buildings found</h3>
          <p className="text-slate-600 mb-2 max-w-md mx-auto leading-relaxed">
            {searchQuery
              ? `We couldn't find any buildings matching "${searchQuery}"`
              : "No buildings match your current filters"}
          </p>
          <p className="text-sm text-slate-500">Try adjusting your search terms or removing some filters</p>

          <div className="mt-8 flex items-center justify-center space-x-6 text-sm text-slate-400">
            <div className="flex items-center space-x-2">
              <Building2 className="h-4 w-4" />
              <span>10 Total Buildings</span>
            </div>
            <div className="flex items-center space-x-2">
              <MapPin className="h-4 w-4" />
              <span>Campus Wide</span>
            </div>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 md:gap-8">
      {buildings.map((building, index) => (
        <div key={building.id} className="animate-fade-in" style={{ animationDelay: `${index * 90}ms` }}>
          <BuildingCard building={building} onViewDetails={onViewDetails} />
        </div>
      ))}
    </div>
  )
})
