import axios from 'axios'

const getAPI_URL = () => process.env.MIX_GIN_API_URL || axios.defaults.baseURL

const ADMIN_CREDENTIALS = {
    email: process.env.MIX_GIN_USER_EMAIL,
    password: process.env.MIX_GIN_USER_PASSWORD,
    name: process.env.MIX_GIN_USER_NAME
}

const AuthService = {
    login: async (email, password) => {
        try {
            const response = await axios.post(`${getAPI_URL()}/user/login`, { email, password })
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

    autoAuth: async () => {
        if (AuthService.getCurrentUser()) {
            AuthService.setAuthHeader()
            return true
        }

        try {
            // Try login first
            await AuthService.login(ADMIN_CREDENTIALS.email, ADMIN_CREDENTIALS.password)
            return true
        } catch (loginError) {
            if (loginError.response?.status === 406 || loginError.response?.status === 401) {
                // User not exists, register then login
                try {
                    await AuthService.register(ADMIN_CREDENTIALS)
                    await AuthService.login(ADMIN_CREDENTIALS.email, ADMIN_CREDENTIALS.password)
                    return true
                } catch (registerError) {
                    console.error('Auto register failed:', registerError)
                    return false
                }
            } else {
                console.error('Auto login failed:', loginError)
                return false
            }
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
