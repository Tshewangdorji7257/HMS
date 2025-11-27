import { Card, CardContent } from "@/components/ui/card"
import type { LucideIcon } from "lucide-react"
import { TrendingUp, TrendingDown } from "lucide-react"

interface StatsCardProps {
  title: string
  value: string | number
  description?: string
  icon: LucideIcon
  trend?: {
    value: number
    isPositive: boolean
  }
}

export function StatsCard({ title, value, description, icon: Icon, trend }: StatsCardProps) {
  return (
    <Card className="group hover:shadow-xl transition-all duration-500 hover:-translate-y-1 border-0 bg-white/80 backdrop-blur-sm overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-br from-slate-50/50 to-white/30 opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
      <CardContent className="p-6 relative">
        <div className="flex items-start justify-between mb-4">
          <div className="space-y-3 flex-1">
            <p className="text-sm font-medium text-slate-500 uppercase tracking-wide">{title}</p>
            <div className="space-y-1">
              <p className="text-3xl font-serif font-light text-slate-800 leading-none">{value}</p>
              {description && <p className="text-sm text-slate-500 leading-relaxed">{description}</p>}
            </div>
          </div>
          <div className="relative">
            <div className="h-14 w-14 rounded-2xl bg-gradient-to-br from-rose-100 to-pink-100 flex items-center justify-center group-hover:scale-110 transition-transform duration-300">
              <Icon className="h-7 w-7 text-rose-600" />
            </div>
            <div className="absolute -inset-1 bg-gradient-to-br from-rose-200/50 to-pink-200/50 rounded-2xl blur opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
          </div>
        </div>

        {trend && (
          <div className="flex items-center justify-between pt-4 border-t border-slate-100">
            <div className="flex items-center space-x-2">
              {trend.isPositive ? (
                <TrendingUp className="h-4 w-4 text-emerald-600" />
              ) : (
                <TrendingDown className="h-4 w-4 text-rose-600" />
              )}
              <span className={`text-sm font-medium ${trend.isPositive ? "text-emerald-600" : "text-rose-600"}`}>
                {trend.isPositive ? "+" : ""}
                {trend.value}%
              </span>
            </div>
            <span className="text-xs text-slate-400">vs last month</span>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
