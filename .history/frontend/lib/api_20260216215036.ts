const API_URL = "http://localhost:8080"

export interface ApiError {
  error: string
}

export const api = async <T = unknown>(
  endpoint: string,
  method = "GET",
  body?: Record<string, unknown>
): Promise<T> => {
  const token = typeof window !== "undefined" ? localStorage.getItem("token") : null

  const headers: HeadersInit = {
    "Content-Type": "application/json",
  }

  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  try {
    const res = await fetch(`${API_URL}${endpoint}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    })

    const data = await res.json()

    if (!res.ok) {
      throw new Error(data.error || `HTTP error! status: ${res.status}`)
    }

    return data as T
  } catch (error) {
    if (error instanceof Error) {
      throw error
    }
    throw new Error("An unexpected error occurred")
  }
}
