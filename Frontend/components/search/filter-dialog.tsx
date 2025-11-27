"use client"

import { useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import { Checkbox } from "@/components/ui/checkbox"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Slider } from "@/components/ui/slider"
import { Filter, X } from "lucide-react"

export interface FilterOptions {
  roomTypes: string[]
  amenities: string[]
  availabilityRange: [number, number]
  sortBy: string
}

interface FilterDialogProps {
  filters: FilterOptions
  onFiltersChange: (filters: FilterOptions) => void
  availableAmenities: string[]
}

export function FilterDialog({ filters, onFiltersChange, availableAmenities }: FilterDialogProps) {
  const [open, setOpen] = useState(false)
  const [localFilters, setLocalFilters] = useState<FilterOptions>(filters)

  const roomTypes = ["single", "double", "triple", "quad"]
  const sortOptions = [
    { value: "name", label: "Building Name" },
    { value: "availability", label: "Most Available" },
    { value: "occupancy", label: "Least Occupied" },
  ]

  const handleRoomTypeChange = (roomType: string, checked: boolean) => {
    const newRoomTypes = checked
      ? [...localFilters.roomTypes, roomType]
      : localFilters.roomTypes.filter((type) => type !== roomType)

    setLocalFilters({ ...localFilters, roomTypes: newRoomTypes })
  }

  const handleAmenityChange = (amenity: string, checked: boolean) => {
    const newAmenities = checked
      ? [...localFilters.amenities, amenity]
      : localFilters.amenities.filter((a) => a !== amenity)

    setLocalFilters({ ...localFilters, amenities: newAmenities })
  }

  const handleApplyFilters = () => {
    onFiltersChange(localFilters)
    setOpen(false)
  }

  const handleClearFilters = () => {
    const clearedFilters: FilterOptions = {
      roomTypes: [],
      amenities: [],
      availabilityRange: [0, 100],
      sortBy: "name",
    }
    setLocalFilters(clearedFilters)
    onFiltersChange(clearedFilters)
    setOpen(false)
  }

  const hasActiveFilters =
    filters.roomTypes.length > 0 ||
    filters.amenities.length > 0 ||
    filters.availabilityRange[0] > 0 ||
    filters.availabilityRange[1] < 100 ||
    filters.sortBy !== "name"

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" className="flex items-center gap-2 bg-transparent relative">
          <Filter className="h-4 w-4" />
          Filters
          {hasActiveFilters && <div className="absolute -top-1 -right-1 h-2 w-2 bg-primary rounded-full" />}
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Filter Buildings</DialogTitle>
          <DialogDescription>Customize your search to find the perfect room</DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {/* Room Types */}
          <div className="space-y-3">
            <Label className="text-sm font-medium">Room Types</Label>
            <div className="grid grid-cols-2 gap-3">
              {roomTypes.map((roomType) => (
                <div key={roomType} className="flex items-center space-x-2">
                  <Checkbox
                    id={roomType}
                    checked={localFilters.roomTypes.includes(roomType)}
                    onCheckedChange={(checked) => handleRoomTypeChange(roomType, checked as boolean)}
                  />
                  <Label htmlFor={roomType} className="text-sm capitalize">
                    {roomType}
                  </Label>
                </div>
              ))}
            </div>
          </div>

          {/* Amenities */}
          <div className="space-y-3">
            <Label className="text-sm font-medium">Amenities</Label>
            <div className="grid grid-cols-1 gap-2 max-h-32 overflow-y-auto">
              {availableAmenities.slice(0, 8).map((amenity) => (
                <div key={amenity} className="flex items-center space-x-2">
                  <Checkbox
                    id={amenity}
                    checked={localFilters.amenities.includes(amenity)}
                    onCheckedChange={(checked) => handleAmenityChange(amenity, checked as boolean)}
                  />
                  <Label htmlFor={amenity} className="text-sm">
                    {amenity}
                  </Label>
                </div>
              ))}
            </div>
          </div>

          {/* Availability Range */}
          <div className="space-y-3">
            <Label className="text-sm font-medium">Minimum Available Beds: {localFilters.availabilityRange[0]}</Label>
            <Slider
              value={localFilters.availabilityRange}
              onValueChange={(value) =>
                setLocalFilters({ ...localFilters, availabilityRange: value as [number, number] })
              }
              max={50}
              min={0}
              step={1}
              className="w-full"
            />
          </div>

          {/* Sort By */}
          <div className="space-y-3">
            <Label className="text-sm font-medium">Sort By</Label>
            <Select
              value={localFilters.sortBy}
              onValueChange={(value) => setLocalFilters({ ...localFilters, sortBy: value })}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {sortOptions.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        <DialogFooter className="flex-col sm:flex-row gap-2">
          <Button variant="outline" onClick={handleClearFilters} className="w-full sm:w-auto bg-transparent">
            <X className="mr-2 h-4 w-4" />
            Clear All
          </Button>
          <Button onClick={handleApplyFilters} className="w-full sm:w-auto">
            Apply Filters
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
