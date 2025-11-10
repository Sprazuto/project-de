import { useApi } from '@/composables/useApi'

const AuthService = {
  login: async (username, password) => {
    try {
      // Match backend API structure: { email, username, password }
      const loginData = {
        username,
        password
      }

      // Get API instance from composable
      const { $api } = useApi()

      // Fixed: Use relative path so $api can properly combine with baseURL
      const response = await $api(`/user/login`, {
        method: 'POST',
        body: loginData
      })

      // Handle different response formats
      if (response.user && response.token) {
        // Response format: { user, token, message }
        // Backend returns: { token: { access_token: "jwt", refresh_token: "jwt" } }
        // We need to extract the access_token from the object
        const actualToken = typeof response.token === 'object' ? response.token.access_token : response.token

        localStorage.setItem('user', JSON.stringify(response.user))
        localStorage.setItem('access_token', actualToken)

        return {
          user: response.user,
          token: actualToken,
          message: response.message || 'Login successful'
        }
      } else if (response.token) {
        // Response format: { token, message }
        // Backend returns: { token: { access_token: "jwt", refresh_token: "jwt" } }
        // We need to extract the access_token from the object
        const actualToken = typeof response.token === 'object' ? response.token.access_token : response.token

        const user = JSON.parse(localStorage.getItem('user')) || null
        localStorage.setItem('access_token', actualToken)

        return {
          user: user,
          token: actualToken,
          message: response.message || 'Login successful'
        }
      } else {
        throw new Error('Invalid response format from server')
      }
    } catch (error) {
      console.error('Login error:', error.status, error.data || error.message)
      throw error
    }
  },

  register: async (user) => {
    try {
      const { $api } = useApi()
      const response = await $api('/user/register', {
        method: 'POST',
        body: user
      })

      if (response.user) {
        localStorage.setItem('user', JSON.stringify(response.user))
      }

      return response
    } catch (error) {
      console.error('Register error:', error.status, error.data || error.message)
      throw error
    }
  },

  logout: () => {
    localStorage.removeItem('user')
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  },

  getCurrentUser: () => {
    return JSON.parse(localStorage.getItem('user')) || null
  },

  getAccessToken: () => {
    return localStorage.getItem('access_token') || null
  },

  // Additional token management methods
  getRefreshToken: () => {
    return localStorage.getItem('refresh_token') || null
  },

  isTokenExpired: () => {
    const token = localStorage.getItem('access_token')
    if (!token) return true

    try {
      // JWT token expiry check
      const payload = JSON.parse(atob(token.split('.')[1]))
      const currentTime = Date.now() / 1000

      return payload.exp < currentTime
    } catch (e) {
      // If JWT parsing fails, assume token is invalid
      return true
    }
  },

  hasValidToken: () => {
    return !!(localStorage.getItem('access_token') && !AuthService.isTokenExpired())
  },

  // Auth state persistence
  saveAuthState: (user, token, refreshToken = null) => {
    if (user) localStorage.setItem('user', JSON.stringify(user))
    if (token) localStorage.setItem('access_token', token)
    if (refreshToken) localStorage.setItem('refresh_token', refreshToken)
  },

  clearAuthState: () => {
    localStorage.removeItem('user')
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }
}

export default AuthService
