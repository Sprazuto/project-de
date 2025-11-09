/**
 * Centralized error handling utilities
 * Provides consistent error handling, logging, and user feedback across the application
 */

import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from './useAuth'

// Global error state
const globalError = ref(null)
const errorHistory = ref([])

// Error types
export const ERROR_TYPES = {
  NETWORK: 'network',
  AUTHENTICATION: 'authentication',
  AUTHORIZATION: 'authorization',
  VALIDATION: 'validation',
  SERVER: 'server',
  CLIENT: 'client',
  UNKNOWN: 'unknown'
}

// Error severity levels
export const ERROR_SEVERITY = {
  LOW: 'low',
  MEDIUM: 'medium',
  HIGH: 'high',
  CRITICAL: 'critical'
}

// Error handler
export function useErrorHandler() {
  const router = useRouter()
  const { logout } = useAuth()

  // Computed properties
  const hasError = computed(() => !!globalError.value)
  const currentError = computed(() => globalError.value)
  const recentErrors = computed(() => errorHistory.value.slice(-5)) // Last 5 errors

  // Parse error to determine type and severity
  const parseError = (error) => {
    const parsed = {
      message: 'An unexpected error occurred',
      type: ERROR_TYPES.UNKNOWN,
      severity: ERROR_SEVERITY.MEDIUM,
      statusCode: null,
      details: null,
      timestamp: new Date().toISOString(),
      id: Date.now() + Math.random()
    }

    if (!error) return parsed

    // Handle different error formats
    if (error.response) {
      // Axios/HTTP error
      parsed.statusCode = error.response.status
      parsed.message = error.response.data?.message || error.message || 'Request failed'
      parsed.details = error.response.data

      // Determine error type based on status code
      switch (error.response.status) {
        case 401:
          parsed.type = ERROR_TYPES.AUTHENTICATION
          parsed.severity = ERROR_SEVERITY.HIGH
          break
        case 403:
          parsed.type = ERROR_TYPES.AUTHORIZATION
          parsed.severity = ERROR_SEVERITY.HIGH
          break
        case 422:
          parsed.type = ERROR_TYPES.VALIDATION
          parsed.severity = ERROR_SEVERITY.MEDIUM
          break
        case 500:
        case 502:
        case 503:
        case 504:
          parsed.type = ERROR_TYPES.SERVER
          parsed.severity = ERROR_SEVERITY.HIGH
          break
        default:
          parsed.type = ERROR_TYPES.CLIENT
          parsed.severity = ERROR_SEVERITY.MEDIUM
      }
    } else if (error.name === 'NetworkError' || error.code === 'NETWORK_ERROR') {
      parsed.type = ERROR_TYPES.NETWORK
      parsed.severity = ERROR_SEVERITY.HIGH
      parsed.message = 'Network connection failed. Please check your internet connection.'
    } else if (error.name === 'TypeError' && error.message.includes('fetch')) {
      parsed.type = ERROR_TYPES.NETWORK
      parsed.severity = ERROR_SEVERITY.HIGH
      parsed.message = 'Unable to connect to the server. Please try again later.'
    } else {
      // Generic error
      parsed.message = error.message || 'An unexpected error occurred'
      parsed.type = ERROR_TYPES.UNKNOWN
    }

    return parsed
  }

  // Set global error
  const setError = (error) => {
    const parsedError = parseError(error)

    globalError.value = parsedError
    errorHistory.value.push(parsedError)

    // Log error based on severity
    if (parsedError.severity === ERROR_SEVERITY.CRITICAL) {
      console.error('CRITICAL ERROR:', parsedError)
    } else if (parsedError.severity === ERROR_SEVERITY.HIGH) {
      console.error('HIGH SEVERITY ERROR:', parsedError)
    } else {
      console.warn('ERROR:', parsedError)
    }

    return parsedError
  }

  // Clear current error
  const clearError = () => {
    globalError.value = null
  }

  // Handle authentication errors
  const handleAuthError = async (error) => {
    const parsedError = setError(error)

    if (parsedError.type === ERROR_TYPES.AUTHENTICATION) {
      // Logout and redirect to login
      try {
        await logout()
        await router.push('/login')
      } catch (e) {
        console.error('Logout after auth error failed:', e)

        // Fallback redirect
        window.location.href = '/login'
      }
    }

    return parsedError
  }

  // Handle network errors with retry logic
  const handleNetworkError = (error, retryFn = null, maxRetries = 3) => {
    const parsedError = setError(error)

    if (parsedError.type === ERROR_TYPES.NETWORK && retryFn) {
      // Implement retry logic here
      // This could be enhanced with exponential backoff
      console.log(`Retrying network request... (${maxRetries} attempts max)`)
    }

    return parsedError
  }

  // Create user-friendly error message
  const getUserFriendlyMessage = (error) => {
    const parsed = parseError(error)

    switch (parsed.type) {
      case ERROR_TYPES.NETWORK:
        return 'No internet connection. Please check your network and try again.'
      case ERROR_TYPES.AUTHENTICATION:
        return 'Your session has expired. Please log in again.'
      case ERROR_TYPES.AUTHORIZATION:
        return "You don't have permission to perform this action."
      case ERROR_TYPES.VALIDATION:
        return 'Please check your input and try again.'
      case ERROR_TYPES.SERVER:
        return 'Server error. Please try again later.'
      default:
        return parsed.message || 'Something went wrong. Please try again.'
    }
  }

  // Error reporting function (integrate with external services if needed)
  const reportError = (error, context = {}) => {
    const parsed = parseError(error)

    // Add context information
    const report = {
      ...parsed,
      context,
      userAgent: navigator.userAgent,
      url: window.location.href,
      timestamp: new Date().toISOString()
    }

    // Here you could integrate with external error reporting services
    // like Sentry, LogRocket, etc.
    console.log('Error Report:', report)

    return report
  }

  return {
    // State
    globalError,
    errorHistory,

    // Computed
    hasError,
    currentError,
    recentErrors,

    // Error handling
    setError,
    clearError,
    handleAuthError,
    handleNetworkError,
    getUserFriendlyMessage,
    reportError,

    // Constants
    ERROR_TYPES,
    ERROR_SEVERITY
  }
}

// Composition function for components
export function useErrorState() {
  const { hasError, currentError, clearError } = useErrorHandler()

  return {
    hasError: hasError.value,
    error: currentError.value,
    clearError
  }
}
