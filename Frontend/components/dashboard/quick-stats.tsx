"use client"

import { useMemo } from "react"
import { StatsCard } from "./stats-card"
import { Building2, Bed, Users, Calendar } from "lucide-react"
import { useHostelData } from "@/hooks/use-hostel-data"

export function QuickStats() {
  const { data } = useHostelData()

  const stats = useMemo(() => {
    const totalBuildings = data.buildings.length
    const totalBeds = data.buildings.reduce((sum, building) => sum + building.totalBeds, 0)
    const availableBeds = data.buildings.reduce((sum, building) => sum + building.availableBeds, 0)
    const occupancyRate = totalBeds > 0 ? Math.abs(Math.round(((totalBeds - availableBeds) / totalBeds) * 100)) : 0
    const activeBookings = data.bookings.filter((booking) => booking.status === "active").length
    
    return { totalBuildings, totalBeds, availableBeds, occupancyRate, activeBookings }
  }, [data.buildings, data.bookings])

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <StatsCard title="Total Buildings" value={stats.totalBuildings} description="Across campus" icon={Building2} />
      <StatsCard title="Available Beds" value={stats.availableBeds} description={`Out of ${stats.totalBeds} total`} icon={Bed} />
      <StatsCard
        title="Occupancy Rate"
        value={`${stats.occupancyRate}%`}
        description="Current occupancy"
        icon={Users}
        trend={{ value: 5, isPositive: true }}
      />
      <StatsCard title="Your Bookings" value={stats.activeBookings} description="Active bookings" icon={Calendar} />
    </div>
  )
}
