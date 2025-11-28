// Backend data service for hostel management
import type { HostelData, Building, Room, Bed, Booking } from "./types"
import { apiFetch, handleApiError } from './api-config'
import { authService } from './auth'

// Data management functions with backend integration
class HostelDataService {
  private readonly bookingsKey = "hostel-bookings-cache"
  private readonly buildingsKey = "hostel-buildings-cache"

  // Get all buildings from backend
  async getBuildings(): Promise<Building[]> {
    try {
      console.log('🏢 getBuildings: Calling /api/buildings...')
      const response = await apiFetch<{ success: boolean; buildings: any[] }>('/api/buildings')
      console.log('🏢 getBuildings: Response received:', { 
        success: response?.success, 
        buildingsCount: response?.buildings?.length 
      })
      
      if (response.success && response.buildings) {
        // Transform backend data to match frontend types
        console.log('🏢 getBuildings: Transforming buildings...')
        const buildings = response.buildings.map(this.transformBuilding)
        console.log('🏢 getBuildings: Transformed', buildings.length, 'buildings')
        
        // Cache buildings
        if (typeof globalThis.window !== 'undefined') {
          localStorage.setItem(this.buildingsKey, JSON.stringify(buildings))
          console.log('🏢 getBuildings: Buildings cached')
        }
        
        return buildings
      }
      
      console.log('🏢 getBuildings: No buildings in response, using cache')
      return this.getCachedBuildings()
    } catch (error) {
      console.error("❌ Error fetching buildings:", error)
      return this.getCachedBuildings()
    }
  }

  // Get specific building by ID
  async getBuildingById(buildingId: string): Promise<Building | null> {
    try {
      const response = await apiFetch<{ success: boolean; building: any }>(
        `/api/buildings/${buildingId}`
      )
      
      if (response.success && response.building) {
        return this.transformBuilding(response.building)
      }
      
      return null
    } catch (error) {
      console.error("Error fetching building:", error)
      return null
    }
  }

  // Get specific room by ID
  async getRoomById(buildingId: string, roomId: string): Promise<Room | null> {
    try {
      const response = await apiFetch<{ success: boolean; room: any }>(
        `/api/buildings/${buildingId}/rooms/${roomId}`
      )
      
      if (response.success && response.room) {
        return this.transformRoom(response.room)
      }
      
      return null
    } catch (error) {
      console.error("Error fetching room:", error)
      return null
    }
  }

  // Search buildings
  async searchBuildings(query: string): Promise<Building[]> {
    try {
      const response = await apiFetch<{ success: boolean; buildings: any[] }>(
        `/api/buildings/search?q=${encodeURIComponent(query)}`
      )
      
      if (response.success && response.buildings) {
        return response.buildings.map(this.transformBuilding)
      }
      
      return []
    } catch (error) {
      console.error("Error searching buildings:", error)
      return []
    }
  }

  // Get all data (buildings + bookings)
  async getData(): Promise<HostelData> {
    try {
      console.log('📦 getData: Fetching buildings from backend...')
      const buildings = await this.getBuildings()
      console.log('📦 getData: Buildings fetched:', buildings.length)
      
      console.log('📦 getData: Fetching bookings from backend...')
      const bookings = await this.getBookings()
      console.log('📦 getData: Bookings fetched:', bookings.length)
      
      return {
        buildings,
        bookings,
      }
    } catch (error) {
      console.error("❌ Error loading hostel data:", error)
      return {
        buildings: this.getCachedBuildings(),
        bookings: [],
      }
    }
  }

  // Create a booking
  async createBooking(bookingData: {
    buildingId: string
    buildingName: string
    roomId: string
    roomNumber: string
    bedId: string
    bedNumber: number
  }): Promise<{ success: boolean; booking?: Booking; error?: string }> {
    try {
      const token = authService.getToken()
      if (!token) {
        return { success: false, error: 'Authentication required' }
      }

      const authState = authService.getAuthState()
      if (!authState.user) {
        return { success: false, error: 'User not found' }
      }

      const response = await apiFetch<{ success: boolean; booking: any; message: string }>(
        '/api/bookings',
        {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token}`,
          },
          body: JSON.stringify({
            user_id: authState.user.id,
            user_name: authState.user.name,
            user_email: authState.user.email,
            building_id: bookingData.buildingId,
            building_name: bookingData.buildingName,
            room_id: bookingData.roomId,
            room_number: bookingData.roomNumber,
            bed_id: bookingData.bedId,
            bed_number: bookingData.bedNumber,
          }),
        }
      )

      if (response.success && response.booking) {
        return { 
          success: true, 
          booking: this.transformBooking(response.booking)
        }
      }

      return { success: false, error: response.message || 'Booking failed' }
    } catch (error: any) {
      console.error("Error creating booking:", error)
      return { success: false, error: handleApiError(error) }
    }
  }

  // Get user bookings
  async getBookings(userId?: string): Promise<Booking[]> {
    try {
      const token = authService.getToken()
      if (!token) {
        return this.getCachedBookings()
      }

      const authState = authService.getAuthState()
      const targetUserId = userId || authState.user?.id

      if (!targetUserId) {
        return this.getCachedBookings()
      }

      const response = await apiFetch<{ success: boolean; bookings: any[] }>(
        `/api/bookings/users/${targetUserId}`,
        {
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        }
      )

      if (response.success && response.bookings) {
        const bookings = response.bookings.map(this.transformBooking)
        
        // Cache bookings
        if (typeof globalThis.window !== 'undefined') {
          localStorage.setItem(this.bookingsKey, JSON.stringify(bookings))
        }
        
        return bookings
      }

      return this.getCachedBookings()
    } catch (error) {
      console.error("Error fetching bookings:", error)
      return this.getCachedBookings()
    }
  }

  // Get all bookings (admin only)
  async getAllBookings(adminToken?: string): Promise<Booking[]> {
    try {
      // Use admin token if provided, otherwise try regular token
      let token = adminToken
      if (!token && typeof globalThis.window !== 'undefined') {
        token = localStorage.getItem('hostel-admin-auth-token') || authService.getToken()
      }
      
      if (!token) {
        console.warn('No auth token available for getAllBookings')
        return []
      }

      console.log('📋 Fetching all bookings with token:', token ? 'Token present' : 'No token')

      const response = await apiFetch<{ success: boolean; bookings: any[] }>(
        '/api/bookings',
        {
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        }
      )

      console.log('📋 getAllBookings response:', { 
        success: response.success, 
        bookingsCount: response.bookings?.length 
      })

      if (response.success && response.bookings) {
        return response.bookings.map(this.transformBooking)
      }

      return []
    } catch (error) {
      console.error("Error fetching all bookings:", error)
      return []
    }
  }

  // Cancel a booking
  async cancelBooking(bookingId: string): Promise<{ success: boolean; error?: string }> {
    try {
      const token = authService.getToken()
      if (!token) {
        return { success: false, error: 'Authentication required' }
      }

      const authState = authService.getAuthState()
      const userEmail = authState.user?.email || ''

      const response = await apiFetch<{ success: boolean; message: string }>(
        `/api/bookings/${bookingId}/cancel?user_email=${encodeURIComponent(userEmail)}`,
        {
          method: 'PUT',
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        }
      )

      return { 
        success: response.success,
        error: response.success ? undefined : response.message
      }
    } catch (error: any) {
      console.error("Error cancelling booking:", error)
      return { success: false, error: handleApiError(error) }
    }
  }

  // Transform backend building data to frontend format
  private transformBuilding = (backendBuilding: any): Building => {
    return {
      id: backendBuilding.id,
      name: backendBuilding.name,
      description: backendBuilding.description,
      totalRooms: backendBuilding.total_rooms,
      totalBeds: backendBuilding.total_beds,
      availableBeds: backendBuilding.available_beds,
      rooms: backendBuilding.rooms?.map(this.transformRoom) || [],
      amenities: backendBuilding.amenities || [],
      image: backendBuilding.image || `/placeholder.svg?height=200&width=300&query=${encodeURIComponent(`${backendBuilding.name} hostel`)}`,
    }
  }

  // Transform backend room data to frontend format
  private transformRoom = (backendRoom: any): Room => {
    return {
      id: backendRoom.id,
      number: backendRoom.number,
      buildingId: backendRoom.building_id,
      totalBeds: backendRoom.total_beds,
      availableBeds: backendRoom.available_beds,
      beds: backendRoom.beds?.map(this.transformBed) || [],
      type: backendRoom.type,
      amenities: backendRoom.amenities || [],
      price: backendRoom.price,
    }
  }

  // Transform backend bed data to frontend format
  private transformBed = (backendBed: any): Bed => {
    return {
      id: backendBed.id,
      number: backendBed.number,
      isOccupied: backendBed.is_occupied,
      occupiedBy: backendBed.occupied_by,
      occupiedByName: backendBed.occupied_by_name,
    }
  }

  // Transform backend booking data to frontend format
  private transformBooking = (backendBooking: any): Booking => {
    return {
      id: backendBooking.id,
      userId: backendBooking.user_id,
      userName: backendBooking.user_name,
      buildingId: backendBooking.building_id,
      buildingName: backendBooking.building_name,
      roomId: backendBooking.room_id,
      roomNumber: backendBooking.room_number,
      bedId: backendBooking.bed_id,
      bedNumber: backendBooking.bed_number,
      bookingDate: backendBooking.booking_date || backendBooking.created_at,
      status: backendBooking.status,
    }
  }

  // Get cached buildings from localStorage
  private getCachedBuildings(): Building[] {
    if (typeof globalThis.window === 'undefined') {
      return []
    }

    try {
      const cached = localStorage.getItem(this.buildingsKey)
      return cached ? JSON.parse(cached) : []
    } catch (error) {
      console.error("Error loading cached buildings:", error)
      return []
    }
  }

  // Get cached bookings from localStorage
  private getCachedBookings(): Booking[] {
    if (typeof globalThis.window === 'undefined') {
      return []
    }

    try {
      const cached = localStorage.getItem(this.bookingsKey)
      return cached ? JSON.parse(cached) : []
    } catch (error) {
      console.error("Error loading cached bookings:", error)
      return []
    }
  }

  // Save bookings (deprecated - now handled by backend)
  saveBookings(bookings: Booking[]): void {
    if (typeof globalThis.window !== 'undefined') {
      localStorage.setItem(this.bookingsKey, JSON.stringify(bookings))
    }
  }
}

export const hostelDataService = new HostelDataService()
