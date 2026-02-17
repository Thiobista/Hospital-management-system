"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { api, ApiError } from "@/lib/api"
import AppointmentForm, { Appointment, Patient } from "@/components/AppointmentForm"
import ProtectedRoute from "@/components/ProtectedRoute"

export default function AppointmentsPage() {
  const router = useRouter()
  const [appointments, setAppointments] = useState<Appointment[]>([])
  const [patients, setPatients] = useState<Patient[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string>("")
  const [showForm, setShowForm] = useState(false)
  const [editingAppointment, setEditingAppointment] = useState<Appointment | undefined>()
  const [deletingId, setDeletingId] = useState<number | null>(null)

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      setLoading(true)
      setError("")
      const [appointmentsData, patientsData] = await Promise.all([
        api<Appointment[]>("/api/appointments"),
        api<Patient[]>("/api/patients"),
      ])
      setAppointments(appointmentsData)
      setPatients(patientsData)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch data")
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async (appointment: Appointment) => {
    try {
      await api<Appointment>("/api/appointments", "POST", appointment)
      setShowForm(false)
      fetchData()
    } catch (err) {
      throw err
    }
  }

  const handleUpdate = async (appointment: Appointment) => {
    if (!appointment.id) return
    try {
      await api<Appointment>(`/api/appointments/${appointment.id}`, "PUT", appointment)
      setEditingAppointment(undefined)
      fetchData()
    } catch (err) {
      throw err
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this appointment?")) return

    try {
      setDeletingId(id)
      await api(`/api/appointments/${id}`, "DELETE")
      fetchData()
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to delete appointment")
    } finally {
      setDeletingId(null)
    }
  }

  const getPatientName = (patientId: number) => {
    const patient = patients.find((p) => p.id === patientId)
    return patient ? patient.name : `Patient #${patientId}`
  }

  const formatDate = (dateString: string) => {
    try {
      const date = new Date(dateString)
      return date.toLocaleString("en-US", {
        year: "numeric",
        month: "short",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      })
    } catch {
      return dateString
    }
  }

  if (showForm) {
    return (
      <div className="max-w-4xl mx-auto">
        <div className="bg-white rounded-lg shadow-lg p-8">
          <h2 className="text-2xl font-bold text-gray-800 mb-6">Create New Appointment</h2>
          <AppointmentForm
            patients={patients}
            onSubmit={handleCreate}
            onCancel={() => setShowForm(false)}
          />
        </div>
      </div>
    )
  }

  if (editingAppointment) {
    return (
      <div className="max-w-4xl mx-auto">
        <div className="bg-white rounded-lg shadow-lg p-8">
          <h2 className="text-2xl font-bold text-gray-800 mb-6">Edit Appointment</h2>
          <AppointmentForm
            appointment={editingAppointment}
            patients={patients}
            onSubmit={handleUpdate}
            onCancel={() => setEditingAppointment(undefined)}
          />
        </div>
      </div>
    )
  }

  return (
    <ProtectedRoute>
      <div className="max-w-7xl mx-auto">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-800">Appointments</h1>
        <button
          onClick={() => setShowForm(true)}
          disabled={patients.length === 0}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg shadow-md transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
        >
          + Create Appointment
        </button>
      </div>

      {patients.length === 0 && (
        <div className="bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded-lg mb-6">
          Please add at least one patient before creating appointments.
        </div>
      )}

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
          {error}
        </div>
      )}

      {loading ? (
        <div className="flex justify-center items-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      ) : appointments.length === 0 ? (
        <div className="bg-white rounded-lg shadow-lg p-12 text-center">
          <p className="text-gray-500 text-lg mb-4">No appointments found</p>
          <button
            onClick={() => setShowForm(true)}
            disabled={patients.length === 0}
            className="text-blue-600 hover:text-blue-700 font-medium disabled:opacity-50"
          >
            Create your first appointment
          </button>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow-lg overflow-hidden">
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Patient
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Doctor
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Date & Time
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Notes
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {appointments.map((appointment) => (
                  <tr key={appointment.id} className="hover:bg-gray-50 transition-colors">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm font-medium text-gray-900">
                        {getPatientName(appointment.patientId)}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">{appointment.doctor}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-500">{formatDate(appointment.date)}</div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="text-sm text-gray-500 max-w-xs truncate">
                        {appointment.notes || "-"}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <div className="flex justify-end space-x-2">
                        <button
                          onClick={() => setEditingAppointment(appointment)}
                          className="text-blue-600 hover:text-blue-900 transition-colors"
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => handleDelete(appointment.id!)}
                          disabled={deletingId === appointment.id}
                          className="text-red-600 hover:text-red-900 transition-colors disabled:opacity-50"
                        >
                          {deletingId === appointment.id ? "Deleting..." : "Delete"}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
    </ProtectedRoute>
  )
}
