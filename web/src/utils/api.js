import { ofetch } from 'ofetch'
import { useRouter } from 'vue-router'

export const $api = ofetch.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  async onRequest({ options }) {
    const accessToken = localStorage.getItem('access_token')

    if (accessToken && accessToken !== 'null' && accessToken !== 'undefined') {
      options.headers = {
        ...options.headers,
        Authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      }
    }
  },
  async onResponseError({ response }) {
    console.error('API Error:', response.status, response._data)
    if (response.status === 401) {
      // Clear invalid tokens
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user')

      // Only redirect if not already on login page
      if (window.location.pathname !== '/login') {
        try {
          const router = useRouter()
          await router.push('/login')
        } catch (error) {
          // Fallback to window.location if router fails
          window.location.href = '/login'
        }
      }
    }
  }
})
