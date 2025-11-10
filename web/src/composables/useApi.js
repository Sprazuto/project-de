/**
 * Composable for API management
 * Handles API client configuration, error handling, and request/response processing
 * Seamlessly integrates with useAuth for authenticated requests
 */

import { ref } from 'vue'
import { ofetch } from 'ofetch'
import { useRouter } from 'vue-router'
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

  // Create ofetch instance with interceptors
  const $api = ofetch.create({
    baseURL,
    timeout,
    async onRequest({ options }) {
      requestCount.value++
      const auth = useAuth()
      let token = auth.getAccessToken?.() || localStorage.getItem('access_token') || localStorage.getItem('token')

      if (token && token !== 'null' && token !== 'undefined') {
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      }
    },
    async onResponse({ response }) {
      requestCount.value--
    },
    async onResponseError({ response }) {
      requestCount.value--
      console.error('API Error:', response.status, response._data)

      if (response.status === 401) {
        // Clear invalid tokens
        try {
          const auth = useAuth()
          auth.logout?.()
        } catch (e) {
          console.warn('Auth cleanup failed:', e)
        }
        localStorage.removeItem('access_token')
        localStorage.removeItem('token')
        localStorage.removeItem('refresh_token')
        localStorage.removeItem('user')

        // Redirect to login
        if (typeof window !== 'undefined' && window.location.pathname !== '/login') {
          try {
            const router = useRouter()
            await router.push('/login')
          } catch (error) {
            window.location.href = '/login'
          }
        }
      }
    }
  })

  // Helper functions
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
      status = error.response.status
      message = error.response._data?.message || error.message
    } else if (error.status) {
      status = error.status
      message = error.data?.message || error.message
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

  // HTTP method helpers using ofetch
  const get = (url, params = {}) => {
    const queryString = new URLSearchParams(params).toString()
    const fullUrl = queryString ? `${url}?${queryString}` : url
    return $api(fullUrl, { method: 'GET' })
      .then((data) => ({ success: true, data }))
      .catch((error) => ({ success: false, error: handleError(error), data: null }))
  }

  const post = (url, data = {}) => {
    return $api(url, { method: 'POST', body: data })
      .then((data) => ({ success: true, data }))
      .catch((error) => ({ success: false, error: handleError(error), data: null }))
  }

  const put = (url, data = {}) => {
    return $api(url, { method: 'PUT', body: data })
      .then((data) => ({ success: true, data }))
      .catch((error) => ({ success: false, error: handleError(error), data: null }))
  }

  const patch = (url, data = {}) => {
    return $api(url, { method: 'PATCH', body: data })
      .then((data) => ({ success: true, data }))
      .catch((error) => ({ success: false, error: handleError(error), data: null }))
  }

  const del = (url) => {
    return $api(url, { method: 'DELETE' })
      .then((data) => ({ success: true, data }))
      .catch((error) => ({ success: false, error: handleError(error), data: null }))
  }

  // File upload
  const upload = async (url, formData) => {
    return $api(url, {
      method: 'POST',
      body: formData,
      headers: {
        // Let browser set Content-Type for FormData
      }
    })
      .then((data) => ({ success: true, data }))
      .catch((error) => ({ success: false, error: handleError(error), data: null }))
  }

  return {
    // State
    isOnline,
    requestCount,

    // ofetch instance for direct use
    $api,

    // Methods
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
