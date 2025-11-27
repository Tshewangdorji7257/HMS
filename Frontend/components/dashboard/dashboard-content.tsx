"use client"

import { useCallback } from "react"
import { QuickStats } from "./quick-stats"
import { QuickBookWidget } from "@/components/booking/quick-book-widget"
import { FilterDialog } from "@/components/search/filter-dialog"
import { SearchResults } from "@/components/search/search-results"
import { ActiveFilters } from "@/components/search/active-filters"
import { useHostelData } from "@/hooks/use-hostel-data"
import { useSearchFilters } from "@/hooks/use-search-filters"
import { Input } from "@/components/ui/input"
import { Search, Sparkles, MapPin } from "lucide-react"
import { Skeleton } from "@/components/ui/skeleton"
import { useRouter } from "next/navigation"

export function DashboardContent() {
  const { data, loading } = useHostelData()
  const router = useRouter()

  const {
    searchQuery,
    setSearchQuery,
    filters,
    setFilters,
    filteredBuildings,
    availableAmenities,
    removeFilter,
    clearAllFilters,
  } = useSearchFilters(data.buildings)

  const handleViewDetails = useCallback((buildingId: string) => {
    router.push(`/building/${buildingId}`)
  }, [router])

  if (loading) {
    return (
      <div className="space-y-12">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} className="h-36 rounded-xl" />
          ))}
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-[420px] rounded-xl" />
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-16">
      {/* Enhanced stats section */}
      <div className="animate-fade-in">
        <div className="text-center mb-8">
          <h2 className="text-2xl font-serif font-light text-slate-800 mb-2">Campus Overview</h2>
          <p className="text-slate-600">Real-time statistics across all hostel buildings</p>
          <div className="w-16 h-0.5 bg-gradient-to-r from-rose-500 to-pink-500 mx-auto mt-4 rounded-full"></div>
        </div>
        <QuickStats />
      </div>

      {/* Enhanced search and booking section */}
      <div className="grid grid-cols-1 lg:grid-cols-5 gap-8 animate-slide-up">
        <div className="lg:col-span-2">
          <QuickBookWidget />
        </div>
        <div className="lg:col-span-3 space-y-6">
          <div className="bg-white/80 backdrop-blur-sm rounded-2xl p-8 border border-slate-200/50 shadow-lg">
            <div className="flex items-center gap-4 mb-6">
              <div className="p-3 rounded-xl bg-gradient-to-br from-blue-100 to-indigo-100">
                <Search className="h-6 w-6 text-blue-600" />
              </div>
              <div>
                <h3 className="text-xl font-serif font-light text-slate-800">Find Your Perfect Room</h3>
                <p className="text-slate-500 text-sm">Search across all buildings and amenities</p>
              </div>
            </div>

            <div className="flex flex-col sm:flex-row gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 text-slate-400 h-5 w-5" />
                <Input
                  placeholder="Search buildings, amenities, locations..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-12 h-12 rounded-xl border-slate-200 focus:border-rose-300 focus:ring-rose-200 transition-all duration-300 bg-white/70"
                />
              </div>
              <div className="flex items-center gap-2">
                <FilterDialog filters={filters} onFiltersChange={setFilters} availableAmenities={availableAmenities} />
                <div className="hidden sm:flex items-center text-sm text-slate-500 bg-slate-50 rounded-lg px-3 py-2">
                  <MapPin className="h-4 w-4 mr-1" />
                  Campus Wide
                </div>
              </div>
            </div>

            {/* Active Filters */}
            <div className="mt-6">
              <ActiveFilters filters={filters} onRemoveFilter={removeFilter} onClearAll={clearAllFilters} />
            </div>
          </div>
        </div>
      </div>

      {/* Enhanced results section */}
      <div className="animate-slide-up">
        <div className="flex items-center justify-between mb-10">
          <div className="flex items-center gap-4">
            <div className="p-3 rounded-xl bg-gradient-to-br from-amber-100 to-orange-100">
              <Sparkles className="h-6 w-6 text-amber-600" />
            </div>
            <div>
              <h2 className="text-3xl font-serif font-light text-slate-800 tracking-tight">
                {searchQuery || filters.roomTypes.length > 0 || filters.amenities.length > 0
                  ? "Search Results"
                  : "Buildings"}
              </h2>
              <p className="text-slate-600 mt-1">
                {filteredBuildings.length} building{filteredBuildings.length !== 1 ? "s" : ""} available •
                {filteredBuildings.reduce((sum, b) => sum + b.availableBeds, 0)} beds ready
              </p>
            </div>
          </div>

          {/* Results summary */}
          <div className="hidden md:flex items-center space-x-4 text-sm text-slate-500">
            <div className="flex items-center space-x-2 bg-slate-50 rounded-lg px-3 py-2">
              <div className="w-2 h-2 bg-green-500 rounded-full"></div>
              <span>Available</span>
            </div>
            <div className="flex items-center space-x-2 bg-slate-50 rounded-lg px-3 py-2">
              <div className="w-2 h-2 bg-amber-500 rounded-full"></div>
              <span>Limited</span>
            </div>
            <div className="flex items-center space-x-2 bg-slate-50 rounded-lg px-3 py-2">
              <div className="w-2 h-2 bg-rose-500 rounded-full"></div>
              <span>Full</span>
            </div>
          </div>
        </div>

        <SearchResults buildings={filteredBuildings} searchQuery={searchQuery} onViewDetails={handleViewDetails} />
      </div>
    </div>
  )
}
