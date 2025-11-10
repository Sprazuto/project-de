/**
 * Composable for dashboard data management
 * Manages dashboard state, data fetching, and business logic
 */

import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuth } from './useAuth'
import { useApi } from './useApi'
import { processRealisasiBulanData, processRealisasiTahunData } from '@/utils/realisasiDataProcessor'

export function useDashboard() {
  const { isAuthenticated } = useAuth()
  const { $api } = useApi()

  // Reactive state
  const loading = ref({
    bulan: true, // Start with loading true to show skeleton immediately
    tahun: true,
    articles: true
  })

  const error = ref({
    bulan: null,
    tahun: null,
    articles: null
  })

  const realisasiBulan = ref([])
  const realisasiTahun = ref([])
  const articles = ref([])

  // Filter state
  const selectedYear = ref(new Date().getFullYear())
  const selectedMonth = ref(new Date().getMonth() + 1)
  const selectedSatker = ref(0)
  const refreshInterval = ref(0) // 30 seconds

  // API functions
  const fetchRealisasiBulan = async (params = {}) => {
    loading.value.bulan = true
    error.value.bulan = null

    try {
      // Use GET with query parameters as confirmed by successful network log
      const queryParams = new URLSearchParams({
        // tahun: Number(selectedYear.value),
        // bulan: Number(selectedMonth.value),
        idsatker: Number(selectedSatker.value),
        ...params
      })

      const response = await $api(`/realisasi-bulan?${queryParams}`)

      const processedData = processRealisasiBulanData(response || {})
      realisasiBulan.value = processedData
    } catch (err) {
      console.error('Error fetching realization data:', err)
      error.value.bulan = err.message || 'Failed to fetch realization data'

      // Set fallback data for demo
      setFallbackRealisasiData()
    } finally {
      loading.value.bulan = false
    }
  }

  const fetchRealisasiTahun = async (params = {}) => {
    loading.value.tahun = true
    error.value.tahun = null

    try {
      // Use GET with query parameters as confirmed by successful network log
      const queryParams = new URLSearchParams({
        // tahun: Number(selectedYear.value),
        idsatker: Number(selectedSatker.value),
        ...params
      })

      const response = await $api(`/realisasi-tahun?${queryParams}`)

      realisasiTahun.value = processRealisasiTahunData(response || {})
    } catch (err) {
      console.error('Error fetching yearly realization data:', err)
    } finally {
      loading.value.tahun = false
    }
  }

  const fetchArticles = async () => {
    loading.value.articles = true
    error.value.articles = null

    try {
      const response = await $api('/articles')

      articles.value = response.data?.results?.[0]?.data || response.data?.data || response.data?.articles || []
    } catch (err) {
      console.error('Error fetching articles:', err)
    } finally {
      loading.value.articles = false
    }
  }
  // Initialization and auto-refresh
  let refreshTimer

  const refreshData = async () => {
    if (!isAuthenticated.value) return

    // Start loading immediately to show skeleton
    loading.value.bulan = true
    loading.value.tahun = true
    loading.value.articles = true

    try {
      await new Promise((resolve) => setTimeout(resolve, 3000))
      await Promise.allSettled([fetchRealisasiBulan(), fetchRealisasiTahun(), fetchArticles()])
    } catch (err) {
      console.error('Error refreshing dashboard data:', err)
    } finally {
      // Reset loading states after fetch completes
      loading.value.bulan = false
      loading.value.tahun = false
      loading.value.articles = false
    }
  }

  const startAutoRefresh = () => {
    stopAutoRefresh()
    if (refreshInterval.value > 0) {
      refreshTimer = setInterval(() => {
        refreshData()
      }, refreshInterval.value)
    }
  }

  const stopAutoRefresh = () => {
    if (refreshTimer) {
      clearInterval(refreshTimer)
      refreshTimer = null
    }
  }

  // Lifecycle
  onMounted(() => {
    refreshData()
    startAutoRefresh()
  })

  onUnmounted(() => {
    stopAutoRefresh()
  })

  // Retry functionality
  const retryFetch = async (type) => {
    switch (type) {
      case 'bulan':
        await fetchRealisasiBulan()
        break
      case 'tahun':
        await fetchRealisasiTahun()
        break
      case 'articles':
        await fetchArticles()
        break
    }
  }

  return {
    // State
    loading,
    error,
    realisasiBulan,
    realisasiTahun,
    articles,
    selectedYear,
    selectedMonth,
    selectedSatker,
    refreshInterval,

    // Actions
    refreshData,
    retryFetch,

    // Lifecycle control
    startAutoRefresh,
    stopAutoRefresh
  }
}
