import { ref, computed } from 'vue'
import AuthService from '@/utils/auth'

// Global authentication state
const currentUser = ref(null)
const isAuthenticated = ref(false)
const isLoading = ref(false)
const authChecked = ref(false)

// Initialize authentication state on composable load
const initializeAuth = () => {
  if (!authChecked.value) {
    const user = AuthService.getCurrentUser()
    if (user) {
      currentUser.value = user
      isAuthenticated.value = true
    }
    authChecked.value = true
  }
}

// Computed properties
const userName = computed(() => {
  return currentUser.value?.name || currentUser.value?.username || 'User'
})

const userEmail = computed(() => {
  return currentUser.value?.email || ''
})

const userRole = computed(() => {
  return currentUser.value?.role || 'user'
})

// Authentication methods
const login = async (email, username, password) => {
  isLoading.value = true

  try {
    // Use username for login since the form uses username
    const response = await AuthService.login(username, password)

    // Update state
    currentUser.value = response.user
    isAuthenticated.value = true

    return response
  } catch (error) {
    // Clear state on error
    currentUser.value = null
    isAuthenticated.value = false
    throw error
  } finally {
    isLoading.value = false
  }
}

const logout = async () => {
  try {
    // Call logout service
    AuthService.logout()

    // Clear state
    currentUser.value = null
    isAuthenticated.value = false

    // Use window.location for redirect to avoid router issues
    window.location.href = '/login'

    return true
  } catch (error) {
    console.error('Logout error:', error)

    // Even if there's an error, clear local state
    currentUser.value = null
    isAuthenticated.value = false

    return false
  }
}

const refreshUser = async () => {
  try {
    const user = AuthService.getCurrentUser()
    if (user) {
      currentUser.value = user
      isAuthenticated.value = true

      return user
    } else {
      currentUser.value = null
      isAuthenticated.value = false

      return null
    }
  } catch (error) {
    console.error('Refresh user error:', error)
    currentUser.value = null
    isAuthenticated.value = false

    return null
  }
}

// Route protection middleware
const authMiddleware = async (to, from, next) => {
  try {
    // Check if user has valid authentication token
    const token = AuthService.getAccessToken()
    const user = AuthService.getCurrentUser()

    if (token && user) {
      // Update state
      currentUser.value = user
      isAuthenticated.value = true

      // User is authenticated, allow access
      next()
    } else {
      // User is not authenticated, redirect to login
      const redirectPath = to.fullPath

      next(`/login?redirect=${encodeURIComponent(redirectPath)}`)
    }
  } catch (error) {
    console.error('Auth middleware error:', error)

    // On error, redirect to login
    const redirectPath = to.fullPath

    next(`/login?redirect=${encodeURIComponent(redirectPath)}`)
  }
}

// Legacy route protection helper (for backward compatibility)
const requireAuth = (to, from, next) => {
  authMiddleware(to, from, next)
}

// Check if user has specific role
const hasRole = (role) => {
  return currentUser.value?.role === role
}

// Check if user has any of the specified roles
const hasAnyRole = (roles) => {
  const userRole = currentUser.value?.role

  return roles.includes(userRole)
}

// Check authentication status (used by router guards)
const checkAuthStatus = async () => {
  try {
    const token = AuthService.getAccessToken()
    const user = AuthService.getCurrentUser()

    if (token && user) {
      // Update state with existing authentication
      currentUser.value = user
      isAuthenticated.value = true

      return true
    } else {
      // No valid authentication found
      currentUser.value = null
      isAuthenticated.value = false

      return false
    }
  } catch (error) {
    console.error('Check auth status error:', error)
    currentUser.value = null
    isAuthenticated.value = false

    return false
  }
}

// Initialize auth state
initializeAuth()

// Token management methods
const getAccessToken = () => {
  return AuthService.getAccessToken()
}

const getRefreshToken = () => {
  return AuthService.getRefreshToken?.() || localStorage.getItem('refresh_token')
}

const isTokenExpired = () => {
  return AuthService.isTokenExpired?.() || false
}

const hasValidToken = () => {
  return AuthService.hasValidToken?.() || false
}

const saveAuthState = (user, token, refreshToken) => {
  return AuthService.saveAuthState?.(user, token, refreshToken)
}

const clearAuthState = () => {
  return AuthService.clearAuthState?.()
}

// Export as a composable function for use with router guards
export function useAuth() {
  return {
    // State
    currentUser,
    isAuthenticated,
    isLoading,

    // Computed
    userName,
    userEmail,
    userRole,

    // Methods
    login,
    logout,
    refreshUser,
    authMiddleware,
    requireAuth,
    hasRole,
    hasAnyRole,
    checkAuthStatus,

    // Token management methods
    getAccessToken,
    getRefreshToken,
    isTokenExpired,
    hasValidToken,
    saveAuthState,
    clearAuthState,

    // Initialize
    initializeAuth
  }
}

export {
  // State
  currentUser,
  isAuthenticated,
  isLoading,

  // Computed
  userName,
  userEmail,
  userRole,

  // Methods
  login,
  logout,
  refreshUser,
  authMiddleware,
  requireAuth,
  hasRole,
  hasAnyRole,
  checkAuthStatus,

  // Token management methods
  getAccessToken,
  getRefreshToken,
  isTokenExpired,
  hasValidToken,
  saveAuthState,
  clearAuthState,

  // Initialize
  initializeAuth
}
