"use client"

import { useState, useMemo } from "react"
import type { Building } from "@/lib/types"
import type { FilterOptions } from "@/components/search/filter-dialog"
import { useDebounce } from "./use-debounce"

export function useSearchFilters(buildings: Building[]) {
  const [searchQuery, setSearchQuery] = useState("")
  const debouncedSearchQuery = useDebounce(searchQuery, 300)
  const [filters, setFilters] = useState<FilterOptions>({
    roomTypes: [],
    amenities: [],
    availabilityRange: [0, 100],
    sortBy: "name",
  })

  // Get all available amenities from buildings
  const availableAmenities = useMemo(() => {
    const amenitySet = new Set<string>()
    if (Array.isArray(buildings)) {
      buildings.forEach((building) => {
        building.amenities?.forEach((amenity) => amenitySet.add(amenity))
        building.rooms?.forEach((room) => {
          room.amenities?.forEach((amenity) => amenitySet.add(amenity))
        })
      })
    }
    return Array.from(amenitySet).sort()
  }, [buildings])

  // Filter and sort buildings based on search query and filters
  const filteredBuildings = useMemo(() => {
    let filtered = Array.isArray(buildings) ? buildings : []

    // Text search with debounced query
    if (debouncedSearchQuery.trim()) {
      const query = debouncedSearchQuery.toLowerCase()
      filtered = filtered.filter(
        (building) =>
          building.name.toLowerCase().includes(query) ||
          building.description.toLowerCase().includes(query) ||
          building.amenities.some((amenity) => amenity.toLowerCase().includes(query)),
      )
    }

    // Room type filter
    if (filters.roomTypes.length > 0) {
      filtered = filtered.filter((building) => building.rooms.some((room) => filters.roomTypes.includes(room.type)))
    }

    // Amenities filter
    if (filters.amenities.length > 0) {
      filtered = filtered.filter((building) =>
        filters.amenities.every(
          (amenity) =>
            building.amenities.includes(amenity) || building.rooms.some((room) => room.amenities.includes(amenity)),
        ),
      )
    }

    // Availability filter
    filtered = filtered.filter(
      (building) =>
        building.availableBeds >= filters.availabilityRange[0] &&
        (filters.availabilityRange[1] === 100 || building.availableBeds <= filters.availabilityRange[1]),
    )

    // Sort buildings
    switch (filters.sortBy) {
      case "availability":
        filtered.sort((a, b) => b.availableBeds - a.availableBeds)
        break
      case "occupancy":
        filtered.sort((a, b) => {
          const aOccupancy = (a.totalBeds - a.availableBeds) / a.totalBeds
          const bOccupancy = (b.totalBeds - b.availableBeds) / b.totalBeds
          return aOccupancy - bOccupancy
        })
        break
      case "name":
      default:
        filtered.sort((a, b) => a.name.localeCompare(b.name))
        break
    }

    return filtered
  }, [buildings, debouncedSearchQuery, filters])

  const removeFilter = (filterType: string, value?: string) => {
    switch (filterType) {
      case "roomType":
        setFilters({
          ...filters,
          roomTypes: filters.roomTypes.filter((type) => type !== value),
        })
        break
      case "amenity":
        setFilters({
          ...filters,
          amenities: filters.amenities.filter((amenity) => amenity !== value),
        })
        break
      case "availability":
        setFilters({
          ...filters,
          availabilityRange: [0, 100],
        })
        break
      case "sortBy":
        setFilters({
          ...filters,
          sortBy: "name",
        })
        break
    }
  }

  const clearAllFilters = () => {
    setFilters({
      roomTypes: [],
      amenities: [],
      availabilityRange: [0, 100],
      sortBy: "name",
    })
    setSearchQuery("")
  }

  return {
    searchQuery,
    setSearchQuery,
    filters,
    setFilters,
    filteredBuildings,
    availableAmenities,
    removeFilter,
    clearAllFilters,
  }
}
