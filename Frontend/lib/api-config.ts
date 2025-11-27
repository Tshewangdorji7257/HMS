// API Configuration for Backend Integration
export const API_CONFIG = {
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000',
  apiBaseURL: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8000/api',
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
}

// Helper to get auth headers with token
export const getAuthHeader = (token?: string | null) => {
  const headers: Record<string, string> = { ...API_CONFIG.headers }
  
  // Try to get token from parameter, localStorage, or return without auth
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  } else if (typeof window !== 'undefined') {
    const storedAuth = localStorage.getItem('hostel-auth-token')
    if (storedAuth) {
      headers['Authorization'] = `Bearer ${storedAuth}`
    }
  }
  
  return headers
}

// Helper for admin auth headers
export const getAdminAuthHeader = (token?: string | null) => {
  const headers: Record<string, string> = { ...API_CONFIG.headers }
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  } else if (typeof window !== 'undefined') {
    const storedAuth = localStorage.getItem('hostel-admin-auth-token')
    if (storedAuth) {
      headers['Authorization'] = `Bearer ${storedAuth}`
    }
  }
  
  return headers
}

// API Response types
export interface ApiResponse<T = any> {
  success: boolean
  message?: string
  data?: T
  error?: string
}

// Error handling helper
export const handleApiError = (error: any): string => {
  if (error.response?.data?.message) {
    return error.response.data.message
  }
  if (error.message) {
    return error.message
  }
  return 'An unexpected error occurred'
}

// Fetch wrapper with error handling
export async function apiFetch<T = any>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = endpoint.startsWith('http') 
    ? endpoint 
    : `${API_CONFIG.baseURL}${endpoint}`
  
  console.log('🌐 API Request:', {
    method: options.method || 'GET',
    url,
    hasAuth: !!(options.headers as any)?.['Authorization']
  })
  
  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        ...API_CONFIG.headers,
        ...options.headers,
      },
    })

    const data = await response.json()
    
    console.log('🌐 API Response:', {
      status: response.status,
      ok: response.ok,
      dataKeys: Object.keys(data)
    })

    if (!response.ok) {
      throw new Error(data.message || data.error || `HTTP ${response.status}`)
    }

    return data
  } catch (error: any) {
    console.error('❌ API Fetch Error:', error)
    throw error
  }
}
