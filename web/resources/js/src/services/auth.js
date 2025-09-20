import axios from '@axios'

const getAPI_URL = () => process.env.MIX_GIN_API_URL || axios.defaults.baseURL

const AuthService = {
    login: async (email, username, password) => {
        try {
            const response = await axios.post(`${getAPI_URL()}/user/login`, { email, username, password })
            if (response.data.token) {
                localStorage.setItem('user', JSON.stringify(response.data.user))
                localStorage.setItem('access_token', response.data.token.access_token)
                localStorage.setItem('refresh_token', response.data.token.refresh_token || '')
                AuthService.setAuthHeader()
            }
            return response.data
        } catch (error) {
            console.error('Login error:', error.response?.status, error.response?.data || error.message)
            throw error
        }
    },

    register: async (user) => {
        try {
            const response = await axios.post(`${getAPI_URL()}/user/register`, user)
            if (response.data.user) {
                localStorage.setItem('user', JSON.stringify(response.data.user))
            }
            return response.data
        } catch (error) {
            console.error('Register error:', error.response?.status, error.response?.data || error.message)
            throw error
        }
    },

    setAuthHeader() {
        const token = this.getAccessToken()
        if (token) {
            axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        } else {
            delete axios.defaults.headers.common['Authorization']
        }
    },

    logout: () => {
        localStorage.removeItem('user')
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        axios.defaults.headers.common['Authorization'] = null
    },

    getCurrentUser: () => {
        return JSON.parse(localStorage.getItem('user')) || null
    },

    getAccessToken: () => {
        return localStorage.getItem('access_token') || null
    }
}

export default AuthService
