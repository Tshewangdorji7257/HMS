"use client"

import { useState, useEffect, useCallback } from "react"
import type { HostelData, Booking } from "@/lib/types"
import { hostelDataService } from "@/lib/data"
import { useAuth } from "./use-auth"

export function useHostelData() {
  const [data, setData] = useState<HostelData>({ buildings: [], bookings: [] })
  const [loading, setLoading] = useState(true)
  const { user, isAuthenticated } = useAuth()

  const refreshData = useCallback(async () => {
    try {
      console.log('🔄 Refreshing hostel data...')
      const hostelData = await hostelDataService.getData()
      console.log('✅ Hostel data loaded:', {
        buildingsCount: hostelData.buildings.length,
        bookingsCount: hostelData.bookings.length
      })
      setData(hostelData)
    } catch (error) {
      console.error('❌ Error loading hostel data:', error)
    }
  }, [])

  useEffect(() => {
    if (!isAuthenticated) {
      setLoading(false)
      return
    }
    
    console.log('🚀 useHostelData initialized, fetching data...')
    refreshData().finally(() => {
      setLoading(false)
      console.log('✅ Loading complete')
    })
  }, [isAuthenticated, refreshData])

  const bookBed = useCallback(async (buildingId: string, roomId: string, bedId: string): Promise<{ success: boolean; error?: string; booking?: Booking }> => {
    if (!user) {
      return { success: false, error: "User not authenticated" }
    }

    // Check if user already has a booking
    const existingBooking = data.bookings.find((booking) => booking.userId === user.id && booking.status === "active")

    if (existingBooking) {
      return { success: false, error: "You already have an active booking. Cancel it first to book a new bed." }
    }

    // Find the building, room, and bed
    const building = data.buildings.find((b) => b.id === buildingId)
    const room = building?.rooms.find((r) => r.id === roomId)
    const bed = room?.beds.find((b) => b.id === bedId)

    if (!building || !room || !bed) {
      return { success: false, error: "Bed not found" }
    }

    if (bed.isOccupied) {
      return { success: false, error: "Bed is already occupied" }
    }

    // Call backend to create booking
    const result = await hostelDataService.createBooking({
      buildingId,
      buildingName: building.name,
      roomId,
      roomNumber: room.number,
      bedId,
      bedNumber: bed.number,
    })

    if (result.success) {
      // Refresh data to update occupancy
      await refreshData()
    }

    return result
  }, [user, data.bookings, data.buildings, refreshData])

  const cancelBooking = useCallback(async (bookingId: string): Promise<{ success: boolean; error?: string }> => {
    const booking = data.bookings.find((b) => b.id === bookingId)

    if (!booking) {
      return { success: false, error: "Booking not found" }
    }

    if (booking.userId !== user?.id) {
      return { success: false, error: "You can only cancel your own bookings" }
    }

    // Call backend to cancel booking
    const result = await hostelDataService.cancelBooking(bookingId)

    if (result.success) {
      // Refresh data to update occupancy
      await refreshData()
    }

    return result
  }, [user?.id, data.bookings, refreshData])

  const getUserBookings = useCallback(async (): Promise<Booking[]> => {
    if (!user) return []
    try {
      return await hostelDataService.getBookings(user.id)
    } catch (error) {
      console.error('Error fetching user bookings:', error)
      return data.bookings.filter((booking) => booking.userId === user.id)
    }
  }, [user, data.bookings])

  const getActiveUserBooking = useCallback((): Booking | null => {
    if (!user) return null
    return data.bookings.find((booking) => booking.userId === user.id && booking.status === "active") || null
  }, [user, data.bookings])

  return {
    data,
    loading,
    refreshData,
    bookBed,
    cancelBooking,
    getUserBookings,
    getActiveUserBooking,
  }
}
