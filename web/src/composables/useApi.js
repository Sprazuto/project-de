/**
 * Composable for API management
 * Handles API client configuration, error handling, and request/response processing
 * Seamlessly integrates with useAuth for authenticated requests
 */

import { ref } from 'vue'
import { useAuth } from './useAuth'

export function useApi() {
  // Configuration
  const baseURL = import.meta.env.VITE_API_BASE_URL
  const timeout = 30000 // 30 seconds

  // State
  const isOnline = ref(navigator.onLine)
  const requestCount = ref(0)

  // Set up online/offline detection
  if (typeof window !== 'undefined') {
    window.addEventListener('online', () => {
      isOnline.value = true
    })
    window.addEventListener('offline', () => {
      isOnline.value = false
    })
  }

  // Helper functions
  const getAuthHeaders = () => {
    // Get token from useAuth composable if available
    const auth = useAuth()
    let token = auth.getAccessToken?.() || localStorage.getItem('access_token') || localStorage.getItem('token')

    return token ? { Authorization: `Bearer ${token}` } : {}
  }

  const handleError = (error, defaultMessage = 'An error occurred') => {
    console.error('API Error:', error)

    if (!isOnline.value) {
      return 'No internet connection. Please check your network.'
    }

    if (error.name === 'AbortError') {
      return 'Request was cancelled'
    }

    // Handle different error response formats
    let status = null
    let message = null

    if (error.response) {
      // Ofetch error format
      status = error.response.status
      message = error.response._data?.message || error.message
    } else if (error.status) {
      // Custom error format
      status = error.status
      message = error.data?.message || error.message
    }

    if (status === 401) {
      // Clear invalid tokens
      try {
        const auth = useAuth()

        auth.logout?.()
      } catch (e) {
        console.warn('Auth cleanup failed:', e)
      }

      // Clear localStorage as fallback
      localStorage.removeItem('access_token')
      localStorage.removeItem('token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user')

      // Redirect to login
      if (typeof window !== 'undefined' && window.location.pathname !== '/login') {
        // TEMPORARILY DISABLED: Redirect to login for debugging
        // window.location.href = '/login'
        console.log('401 Unauthorized - redirect disabled for debugging')
      }

      return 'Authentication failed. Please log in again.'
    }

    switch (status) {
      case 403:
        return 'You do not have permission to perform this action.'
      case 404:
        return 'The requested resource was not found.'
      case 422:
        return 'Validation error. Please check your input.'
      case 500:
        return 'Server error. Please try again later.'
      default:
        return message || error.message || defaultMessage
    }
  }

  // Request function
  const request = async (options) => {
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), timeout)

    try {
      requestCount.value++

      const config = {
        baseURL,
        timeout,
        signal: controller.signal,
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders(),
          ...options.headers
        },
        ...options
      }

      // Properly construct the full URL
      const fullUrl = config.url.startsWith('http') ? config.url : config.baseURL ? `${config.baseURL}${config.url}` : config.url

      const response = await fetch(fullUrl, config)

      clearTimeout(timeoutId)

      // Handle HTTP errors
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        const error = new Error(errorData.message || `HTTP ${response.status}`)

        error.response = { status: response.status, data: errorData }
        throw error
      }

      // Parse response
      const contentType = response.headers.get('content-type')

      const data = contentType?.includes('application/json') ? await response.json() : await response.text()

      return {
        success: true,
        data,
        status: response.status
      }
    } catch (error) {
      clearTimeout(timeoutId)

      // For authentication errors, try to refresh user state
      if (error.response?.status === 401) {
        try {
          const auth = useAuth()

          auth.refreshUser?.()
        } catch (e) {
          console.warn('Auth refresh failed:', e)
        }
      }

      return {
        success: false,
        error: handleError(error),
        data: null
      }
    } finally {
      requestCount.value--
    }
  }

  // HTTP method helpers
  const get = (url, params = {}) => {
    const queryString = new URLSearchParams(params).toString()
    const fullUrl = queryString ? `${url}?${queryString}` : url

    return request({ url: fullUrl, method: 'GET' })
  }

  const post = (url, data = {}) => {
    return request({
      url,
      method: 'POST',
      body: JSON.stringify(data)
    })
  }

  const put = (url, data = {}) => {
    return request({
      url,
      method: 'PUT',
      body: JSON.stringify(data)
    })
  }

  const patch = (url, data = {}) => {
    return request({
      url,
      method: 'PATCH',
      body: JSON.stringify(data)
    })
  }

  const del = (url) => {
    return request({ url, method: 'DELETE' })
  }

  // File upload
  const upload = async (url, formData) => {
    return request({
      url,
      method: 'POST',
      headers: {
        ...getAuthHeaders()

        // Don't set Content-Type for FormData, let browser set it with boundary
      },
      body: formData
    })
  }

  return {
    // State
    isOnline,
    requestCount,

    // Methods
    request,
    get,
    post,
    put,
    patch,
    delete: del,
    upload,
    handleError
  }
}

// Convenience factory for authenticated API calls
export function createAuthenticatedApi(baseUrl) {
  const api = useApi()

  return {
    ...api,
    baseUrl,

    // Helper for authenticated requests with automatic token handling
    authenticatedRequest: async (options) => {
      const auth = useAuth()

      // Ensure user is authenticated
      if (!auth.isAuthenticated.value) {
        throw new Error('User is not authenticated')
      }

      // Merge base URL if provided
      const url = baseUrl ? `${baseUrl}${options.url || ''}` : options.url
      const method = options.method || 'GET'
      const data = options.data || options.body

      switch (method.toUpperCase()) {
        case 'GET':
          return api.get(url, data)
        case 'POST':
          return api.post(url, data)
        case 'PUT':
          return api.put(url, data)
        case 'PATCH':
          return api.patch(url, data)
        case 'DELETE':
          return api.delete(url)
        default:
          return api.request({ url, method, data })
      }
    }
  }
}
