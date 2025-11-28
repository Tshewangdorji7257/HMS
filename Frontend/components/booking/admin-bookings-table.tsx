"use client"

import { useState, useEffect, useMemo } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { hostelDataService } from "@/lib/data"
import { cn } from "@/lib/utils"
import { ArrowDownToLine, Download, Clock, Hash, Loader2, Search, Filter } from "lucide-react"
import type { Booking } from "@/lib/types"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

function formatBookingId(raw: string) {
  // Expect pattern like booking-<timestamp>-<random>
  // We reduce to #<last4OfTimestamp>-<HHmm>
  const parts = raw.split('-')
  if (parts.length >= 3) {
    const ts = parts[1]
    const date = new Date(parseInt(ts, 10))
    if (!isNaN(date.getTime())) {
      const hh = date.getHours().toString().padStart(2, '0')
      const mm = date.getMinutes().toString().padStart(2, '0')
      return `#${ts.slice(-4)}-${hh}${mm}`
    }
  }
  return raw
}

function formatDate(iso: string) {
  const d = new Date(iso)
  if (isNaN(d.getTime())) return iso
  return d.toLocaleString(undefined, {
    year: 'numeric', month: 'short', day: '2-digit', hour: '2-digit', minute: '2-digit'
  })
}

function bookingsToCsv(bookings: any[]) {
  if (!bookings || bookings.length === 0) return ""
  const headers = Object.keys(bookings[0])
  const rows = bookings.map((b) => headers.map((h) => JSON.stringify(b[h] ?? "")).join(","))
  return [headers.join(","), ...rows].join("\n")
}

function cleanStudent(name: string) {
  // Remove everything after first '.' (e.g., 02230312.cst -> 02230312)
  if (!name) return name
  const firstDot = name.indexOf('.')
  return firstDot > 0 ? name.slice(0, firstDot) : name
}

export function AdminBookingsTable() {
  const [bookings, setBookings] = useState<Booking[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState("")
  const [buildingFilter, setBuildingFilter] = useState<string>("all")
  const [roomFilter, setRoomFilter] = useState<string>("all")
  const [statusFilter, setStatusFilter] = useState<string>("all")

  useEffect(() => {
    async function loadBookings() {
      try {
        const data = await hostelDataService.getAllBookings()
        setBookings(data)
      } catch (error) {
        console.error('Error loading bookings:', error)
      } finally {
        setLoading(false)
      }
    }
    loadBookings()
  }, [])

  // Get unique buildings and rooms for filters
  const buildings = useMemo(() => {
    const uniqueBuildings = Array.from(new Set(bookings.map(b => b.buildingName)))
    return uniqueBuildings.sort()
  }, [bookings])

  const rooms = useMemo(() => {
    let filteredBookings = bookings
    if (buildingFilter !== "all") {
      filteredBookings = filteredBookings.filter(b => b.buildingName === buildingFilter)
    }
    const uniqueRooms = Array.from(new Set(filteredBookings.map(b => b.roomNumber)))
    return uniqueRooms.sort()
  }, [bookings, buildingFilter])

  // Filter and search bookings
  const filteredBookings = useMemo(() => {
    let filtered = bookings

    // Apply building filter
    if (buildingFilter !== "all") {
      filtered = filtered.filter(b => b.buildingName === buildingFilter)
    }

    // Apply room filter
    if (roomFilter !== "all") {
      filtered = filtered.filter(b => b.roomNumber === roomFilter)
    }

    // Apply status filter
    if (statusFilter !== "all") {
      filtered = filtered.filter(b => b.status === statusFilter)
    }

    // Apply search query
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase()
      filtered = filtered.filter(b => 
        b.userName.toLowerCase().includes(query) ||
        b.buildingName.toLowerCase().includes(query) ||
        b.roomNumber.toLowerCase().includes(query) ||
        b.bedNumber.toString().includes(query) ||
        b.id.toLowerCase().includes(query)
      )
    }

    return filtered
  }, [bookings, buildingFilter, roomFilter, statusFilter, searchQuery])

  const enhanced = useMemo(() => filteredBookings.map(b => ({
    ...b,
    shortId: formatBookingId(b.id),
    displayDate: formatDate(b.bookingDate),
    studentDisplay: cleanStudent(b.userName)
  })), [filteredBookings])

  const handleExport = () => {
    const csv = bookingsToCsv(bookings)
    const blob = new Blob([csv], { type: "text/csv;charset=utf-8;" })
    const url = URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = url
    a.download = `hostel-bookings-${new Date().toISOString().slice(0, 10)}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="space-y-1">
          <h2 className="text-xl font-serif font-light tracking-tight flex items-center gap-2">
            <Hash className="h-5 w-5 text-rose-500" /> Bookings Overview
          </h2>
          <p className="text-xs text-muted-foreground">
            {enhanced.length} {searchQuery || buildingFilter !== "all" || roomFilter !== "all" || statusFilter !== "all" ? "filtered" : "total"} • 
            {bookings.length} total bookings • 
            Updated <Clock className="inline h-3 w-3" /> {new Date().toLocaleTimeString()}
          </p>
        </div>
        <div className="flex gap-2">
          <Button onClick={handleExport} variant="outline" className="gap-2">
            <Download className="h-4 w-4" /> Export CSV
          </Button>
        </div>
      </div>

      {/* Search and Filters */}
      <div className="flex flex-col gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            type="text"
            placeholder="Search by student, building, room, bed, or booking ID..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-9"
          />
        </div>
        
        <div className="flex flex-wrap gap-2">
          <Select value={buildingFilter} onValueChange={(value) => {
            setBuildingFilter(value)
            setRoomFilter("all") // Reset room filter when building changes
          }}>
            <SelectTrigger className="w-[180px]">
              <Filter className="h-4 w-4 mr-2" />
              <SelectValue placeholder="All Buildings" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Buildings</SelectItem>
              {buildings.map((building) => (
                <SelectItem key={building} value={building}>
                  {building}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select value={roomFilter} onValueChange={setRoomFilter}>
            <SelectTrigger className="w-[150px]">
              <SelectValue placeholder="All Rooms" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Rooms</SelectItem>
              {rooms.map((room) => (
                <SelectItem key={room} value={room}>
                  Room {room}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-[150px]">
              <SelectValue placeholder="All Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Status</SelectItem>
              <SelectItem value="active">Active</SelectItem>
              <SelectItem value="cancelled">Cancelled</SelectItem>
            </SelectContent>
          </Select>

          {(searchQuery || buildingFilter !== "all" || roomFilter !== "all" || statusFilter !== "all") && (
            <Button 
              variant="ghost" 
              size="sm"
              onClick={() => {
                setSearchQuery("")
                setBuildingFilter("all")
                setRoomFilter("all")
                setStatusFilter("all")
              }}
              className="text-muted-foreground hover:text-foreground"
            >
              Clear Filters
            </Button>
          )}
        </div>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12 text-muted-foreground">
          <Loader2 className="h-6 w-6 animate-spin mr-2" />
          Loading bookings...
        </div>
      ) : enhanced.length === 0 ? (
        <div className="text-muted-foreground">No bookings yet.</div>
      ) : (
        <div className="overflow-x-auto rounded-xl border border-slate-200/60 shadow-sm bg-white/70 backdrop-blur">
          <table className="w-full table-fixed border-collapse text-sm">
            <thead className="bg-gradient-to-r from-rose-50 to-pink-50/60 text-slate-600">
              <tr className="[&_th]:py-3 [&_th]:px-3 [&_th]:text-left [&_th]:font-medium">
                <th className="w-[110px]">Booking</th>
                <th className="min-w-[180px]">Student</th>
                <th className="w-[130px]">Building</th>
                <th className="w-[80px]">Room</th>
                <th className="w-[60px]">Bed</th>
                <th className="min-w-[170px]">Date</th>
                <th className="w-[90px]">Status</th>
              </tr>
            </thead>
            <tbody>
              {enhanced.map((b: any, idx: number) => (
                <tr
                  key={b.id}
                  className={cn("border-t border-slate-100/70 hover:bg-rose-50/40 transition-colors", idx % 2 === 1 && "bg-slate-50/40")}
                >
                  <td className="px-3 py-2 font-mono text-xs text-slate-700">{b.shortId}</td>
                  <td className="px-3 py-2 font-medium text-slate-800 truncate max-w-[160px]" title={b.studentDisplay}>{b.studentDisplay}</td>
                  <td className="px-3 py-2 text-slate-700">{b.buildingName}</td>
                  <td className="px-3 py-2"><span className="inline-flex items-center justify-center rounded-md bg-slate-100 px-2 py-1 text-[11px] font-medium text-slate-700">{b.roomNumber}</span></td>
                  <td className="px-3 py-2"><span className="inline-flex items-center justify-center rounded-md bg-slate-100 px-2 py-1 text-[11px] font-medium text-slate-700">{b.bedNumber}</span></td>
                  <td className="px-3 py-2 text-slate-600 whitespace-nowrap">{b.displayDate}</td>
                  <td className="px-3 py-2">
                    <span
                      className={cn(
                        "inline-flex items-center rounded-full px-2.5 py-1 text-[10px] font-semibold tracking-wide uppercase",
                        b.status === 'active' && "bg-emerald-100 text-emerald-700 border border-emerald-200",
                        b.status === 'cancelled' && "bg-rose-100 text-rose-700 border border-rose-200"
                      )}
                    >
                      {b.status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
