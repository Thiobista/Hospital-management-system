"use client"

import { useState, useEffect } from "react"

export interface Patient {
  id: number
  name: string
  age: number
  gender: string
  phone: string
}

export interface Appointment {
  id?: number
  patientId: number
  doctor: string
  date: string
  notes: string
}

interface AppointmentFormProps {
  appointment?: Appointment
  patients: Patient[]
  onSubmit: (appointment: Appointment) => Promise<void>
  onCancel: () => void
  isLoading?: boolean
}

export default function AppointmentForm({
  appointment,
  patients,
  onSubmit,
  onCancel,
  isLoading = false,
}: AppointmentFormProps) {
  const [formData, setFormData] = useState<Appointment>({
    patientId: 0,
    doctor: "",
    date: "",
    notes: "",
  })
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (appointment) {
      // Convert ISO date string to datetime-local format
      let dateValue = ""
      if (appointment.date) {
        try {
          const date = new Date(appointment.date)
          const year = date.getFullYear()
          const month = String(date.getMonth() + 1).padStart(2, "0")
          const day = String(date.getDate()).padStart(2, "0")
          const hours = String(date.getHours()).padStart(2, "0")
          const minutes = String(date.getMinutes()).padStart(2, "0")
          dateValue = `${year}-${month}-${day}T${hours}:${minutes}`
        } catch {
          dateValue = ""
        }
      }
      setFormData({
        patientId: appointment.patientId || 0,
        doctor: appointment.doctor || "",
        date: dateValue,
        notes: appointment.notes || "",
      })
    } else {
      // Set default date to today
      const now = new Date()
      const year = now.getFullYear()
      const month = String(now.getMonth() + 1).padStart(2, "0")
      const day = String(now.getDate()).padStart(2, "0")
      const hours = String(now.getHours()).padStart(2, "0")
      const minutes = String(now.getMinutes()).padStart(2, "0")
      setFormData({
        patientId: 0,
        doctor: "",
        date: `${year}-${month}-${day}T${hours}:${minutes}`,
        notes: "",
      })
    }
  }, [appointment])

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {}

    if (!formData.patientId || formData.patientId === 0) {
      newErrors.patientId = "Please select a patient"
    }

    if (!formData.doctor.trim()) {
      newErrors.doctor = "Doctor name is required"
    }

    if (!formData.date) {
      newErrors.date = "Date and time are required"
    } else {
      const selectedDate = new Date(formData.date)
      const now = new Date()
      if (selectedDate < now) {
        newErrors.date = "Appointment date cannot be in the past"
      }
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!validate()) return

    try {
      // Convert datetime-local format to ISO string for backend
      const submitData = {
        ...formData,
        date: new Date(formData.date).toISOString(),
      }
      await onSubmit(submitData)
    } catch (error) {
      console.error("Form submission error:", error)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label htmlFor="patientId" className="block text-sm font-medium text-gray-700 mb-1">
          Patient <span className="text-red-500">*</span>
        </label>
        <select
          id="patientId"
          value={formData.patientId}
          onChange={(e) => setFormData({ ...formData, patientId: parseInt(e.target.value) })}
          className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
            errors.patientId ? "border-red-500" : "border-gray-300"
          }`}
          disabled={!!appointment}
        >
          <option value={0}>Select a patient</option>
          {patients.map((patient) => (
            <option key={patient.id} value={patient.id}>
              {patient.name} ({patient.age} years, {patient.gender}) - {patient.phone}
            </option>
          ))}
        </select>
        {errors.patientId && <p className="mt-1 text-sm text-red-600">{errors.patientId}</p>}
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label htmlFor="doctor" className="block text-sm font-medium text-gray-700 mb-1">
            Doctor Name <span className="text-red-500">*</span>
          </label>
          <input
            id="doctor"
            type="text"
            value={formData.doctor}
            onChange={(e) => setFormData({ ...formData, doctor: e.target.value })}
            className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
              errors.doctor ? "border-red-500" : "border-gray-300"
            }`}
            placeholder="Enter doctor name"
          />
          {errors.doctor && <p className="mt-1 text-sm text-red-600">{errors.doctor}</p>}
        </div>

        <div>
          <label htmlFor="date" className="block text-sm font-medium text-gray-700 mb-1">
            Date & Time <span className="text-red-500">*</span>
          </label>
          <input
            id="date"
            type="datetime-local"
            value={formData.date}
            onChange={(e) => setFormData({ ...formData, date: e.target.value })}
            className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
              errors.date ? "border-red-500" : "border-gray-300"
            }`}
          />
          {errors.date && <p className="mt-1 text-sm text-red-600">{errors.date}</p>}
        </div>
      </div>

      <div>
        <label htmlFor="notes" className="block text-sm font-medium text-gray-700 mb-1">
          Notes
        </label>
        <textarea
          id="notes"
          value={formData.notes}
          onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
          rows={4}
          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          placeholder="Enter appointment notes (optional)"
        />
      </div>

      <div className="flex justify-end space-x-3 pt-4">
        <button
          type="button"
          onClick={onCancel}
          className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
          disabled={isLoading}
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={isLoading}
          className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isLoading ? "Saving..." : appointment ? "Update Appointment" : "Create Appointment"}
        </button>
      </div>
    </form>
  )
}
