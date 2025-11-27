export interface Bed {
  id: string
  number: number
  isOccupied: boolean
  occupiedBy?: string // User ID who booked this bed
  occupiedByName?: string // User name for display
}

export interface Room {
  id: string
  number: string
  buildingId: string
  totalBeds: number
  availableBeds: number
  beds: Bed[]
  amenities: string[]
  type: "single" | "double" | "triple" | "quad"
  price?: number
}

export interface Building {
  id: string
  name: string
  description: string
  totalRooms: number
  totalBeds: number
  availableBeds: number
  rooms: Room[]
  amenities: string[]
  image: string
}

export interface Booking {
  id: string
  userId: string
  userName: string
  buildingId: string
  buildingName: string
  roomId: string
  roomNumber: string
  bedId: string
  bedNumber: number
  bookingDate: string
  status: "active" | "cancelled"
}

export interface HostelData {
  buildings: Building[]
  bookings: Booking[]
}
