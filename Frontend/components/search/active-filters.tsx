"use client"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import type { FilterOptions } from "./filter-dialog"
import { X } from "lucide-react"

interface ActiveFiltersProps {
  filters: FilterOptions
  onRemoveFilter: (filterType: string, value?: string) => void
  onClearAll: () => void
}

export function ActiveFilters({ filters, onRemoveFilter, onClearAll }: ActiveFiltersProps) {
  const hasActiveFilters =
    filters.roomTypes.length > 0 ||
    filters.amenities.length > 0 ||
    filters.availabilityRange[0] > 0 ||
    filters.availabilityRange[1] < 100 ||
    filters.sortBy !== "name"

  if (!hasActiveFilters) return null

  return (
    <div className="flex flex-wrap items-center gap-2 p-4 bg-muted/30 rounded-lg">
      <span className="text-sm font-medium text-muted-foreground">Active filters:</span>

      {filters.roomTypes.map((roomType) => (
        <Badge key={roomType} variant="secondary" className="flex items-center gap-1">
          <span className="capitalize">{roomType}</span>
          <X
            className="h-3 w-3 cursor-pointer hover:text-destructive"
            onClick={() => onRemoveFilter("roomType", roomType)}
          />
        </Badge>
      ))}

      {filters.amenities.map((amenity) => (
        <Badge key={amenity} variant="secondary" className="flex items-center gap-1">
          <span>{amenity}</span>
          <X
            className="h-3 w-3 cursor-pointer hover:text-destructive"
            onClick={() => onRemoveFilter("amenity", amenity)}
          />
        </Badge>
      ))}

      {(filters.availabilityRange[0] > 0 || filters.availabilityRange[1] < 100) && (
        <Badge variant="secondary" className="flex items-center gap-1">
          <span>Min beds: {filters.availabilityRange[0]}</span>
          <X className="h-3 w-3 cursor-pointer hover:text-destructive" onClick={() => onRemoveFilter("availability")} />
        </Badge>
      )}

      {filters.sortBy !== "name" && (
        <Badge variant="secondary" className="flex items-center gap-1">
          <span>Sort: {filters.sortBy}</span>
          <X className="h-3 w-3 cursor-pointer hover:text-destructive" onClick={() => onRemoveFilter("sortBy")} />
        </Badge>
      )}

      <Button variant="ghost" size="sm" onClick={onClearAll} className="text-xs">
        Clear all
      </Button>
    </div>
  )
}
