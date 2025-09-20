import Vue from 'vue'
import jwtDefaultConfig from './jwtDefaultConfig'

export default class JwtService {
    // Will be used by this service for making API calls
    axiosIns = null

    // jwtConfig <= Will be used by this service
    jwtConfig = { ...jwtDefaultConfig }

    // For Refreshing Token
    isAlreadyFetchingAccessToken = false

    // For Refreshing Token
    subscribers = []

    constructor(axiosIns, jwtOverrideConfig) {
        this.axiosIns = axiosIns
        this.jwtConfig = { ...this.jwtConfig, ...jwtOverrideConfig }

        // Request Interceptor
        this.axiosIns.interceptors.request.use(
            config => {
                // Skip adding Authorization header for refresh token request
                if (config.url.endsWith('/token/refresh')) {
                    delete config.headers.Authorization
                    return config
                }

                // Get token from localStorage
                const accessToken = this.getToken()

                // If token is present add it to request's Authorization Header
                if (accessToken) {
                    // eslint-disable-next-line no-param-reassign
                    config.headers.Authorization = `${this.jwtConfig.tokenType} ${accessToken}`
                }
                return config
            },
            error => Promise.reject(error),
        )

        // Add request/response interceptor
        this.axiosIns.interceptors.response.use(
            response => response,
            error => {
                // const { config, response: { status } } = error
                const { config, response } = error
                const originalRequest = config

                // Skip refresh logic for refresh endpoint itself to avoid loops
                if (config.url.endsWith('/token/refresh') && response && response.status === 401) {
                    return Promise.reject(error)
                }

                // if (status === 401) {
                if (response && response.status === 401) {
                    if (!this.isAlreadyFetchingAccessToken) {
                        this.isAlreadyFetchingAccessToken = true
                        this.refreshToken().then(r => {
                            this.isAlreadyFetchingAccessToken = false

                            // Update accessToken in localStorage
                            this.setToken(r.data.accessToken)
                            this.setRefreshToken(r.data.refreshToken)

                            this.onAccessTokenFetched(r.data.accessToken)
                        }).catch(error => {
                            this.isAlreadyFetchingAccessToken = false
                            console.error('Token refresh failed:', error)
                            this.logout()
                            // Show error message from response using toast
                            const errorMsg = error.response?.data?.message || 'Session expired. Please log in again.'
                            Vue.$toast({
                                component: 'ToastificationContent',
                                props: {
                                    title: 'Session Error',
                                    icon: 'AlertTriangleIcon',
                                    variant: 'danger',
                                    text: errorMsg,
                                },
                            })
                            // Force redirect to login
                            window.location.href = '/login'
                        })
                    }
                    const retryOriginalRequest = new Promise(resolve => {
                        this.addSubscriber(accessToken => {
                            // Make sure to assign accessToken according to your response.
                            // Check: https://pixinvent.ticksy.com/ticket/2413870
                            // Change Authorization header
                            originalRequest.headers.Authorization = `${this.jwtConfig.tokenType} ${accessToken}`
                            resolve(this.axiosIns(originalRequest))
                        })
                    })
                    return retryOriginalRequest
                }
                return Promise.reject(error)
            },
        )
    }

    onAccessTokenFetched(accessToken) {
        this.subscribers = this.subscribers.filter(callback => callback(accessToken))
    }

    addSubscriber(callback) {
        this.subscribers.push(callback)
    }

    getToken() {
        return localStorage.getItem(this.jwtConfig.storageTokenKeyName)
    }

    getRefreshToken() {
        return localStorage.getItem(this.jwtConfig.storageRefreshTokenKeyName)
    }

    setToken(value) {
        localStorage.setItem(this.jwtConfig.storageTokenKeyName, value)
    }

    setRefreshToken(value) {
        localStorage.setItem(this.jwtConfig.storageRefreshTokenKeyName, value)
    }

    login(...args) {
        return this.axiosIns.post(this.jwtConfig.loginEndpoint, ...args)
    }

    register(...args) {
        return this.axiosIns.post(this.jwtConfig.registerEndpoint, ...args)
    }

    refreshToken() {
        // Temporarily remove Authorization header from defaults for refresh request
        const originalHeader = this.axiosIns.defaults.headers.common['Authorization']
        delete this.axiosIns.defaults.headers.common['Authorization']

        return this.axiosIns.post(this.jwtConfig.refreshEndpoint, {
            refresh_token: this.getRefreshToken(),
        }).finally(() => {
            // Restore original header after request
            if (originalHeader) {
                this.axiosIns.defaults.headers.common['Authorization'] = originalHeader
            }
        })
    }

    logout() {
        localStorage.removeItem('user')
        localStorage.removeItem(this.jwtConfig.storageTokenKeyName)
        localStorage.removeItem(this.jwtConfig.storageRefreshTokenKeyName)
        delete this.axiosIns.defaults.headers.common['Authorization']
    }
}
