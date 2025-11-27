"use client"

import { Header } from "@/components/layout/header"
import { MyBookingsContent } from "@/components/bookings/my-bookings-content"

export default function BookingsPage() {
  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="container mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground mb-2">My Bookings</h1>
          <p className="text-muted-foreground">Manage your hostel room reservations</p>
        </div>
        <MyBookingsContent />
      </main>
    </div>
  )
}
