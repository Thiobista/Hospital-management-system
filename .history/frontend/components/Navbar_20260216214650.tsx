"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { useEffect, useState } from "react"

export default function Navbar() {
  const pathname = usePathname()
  const [mounted, setMounted] = useState(false)
  const [token, setToken] = useState<string | null>(null)

  useEffect(() => {
    setMounted(true)
    setToken(localStorage.getItem("token"))
  }, [])

  const logout = () => {
    localStorage.removeItem("token")
    window.location.href = "/login"
  }

  if (!mounted || !token) return null

  const isActive = (path: string) => pathname === path

  return (
    <nav className="bg-gradient-to-r from-blue-600 to-blue-700 shadow-lg">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center space-x-1">
            <div className="flex-shrink-0 flex items-center">
              <span className="text-white text-xl font-bold">🏥 Clinic System</span>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            <Link
              href="/dashboard"
              className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                isActive("/dashboard")
                  ? "bg-blue-800 text-white"
                  : "text-blue-100 hover:bg-blue-800 hover:text-white"
              }`}
            >
              Dashboard
            </Link>
            <Link
              href="/patients"
              className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                isActive("/patients") || pathname?.startsWith("/patients")
                  ? "bg-blue-800 text-white"
                  : "text-blue-100 hover:bg-blue-800 hover:text-white"
              }`}
            >
              Patients
            </Link>
            <Link
              href="/appointments"
              className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                isActive("/appointments") || pathname?.startsWith("/appointments")
                  ? "bg-blue-800 text-white"
                  : "text-blue-100 hover:bg-blue-800 hover:text-white"
              }`}
            >
              Appointments
            </Link>
            <button
              onClick={logout}
              className="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors shadow-sm"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>
  )
}
