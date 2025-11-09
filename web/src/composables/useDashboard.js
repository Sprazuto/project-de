/**
 * Composable for dashboard data management
 * Manages dashboard state, data fetching, and business logic
 */

import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuth } from './useAuth'
import { $api } from '@/utils/api'
import { processRealisasiBulanData, processRealisasiTahunData } from '@/utils/realisasiDataProcessor'

// TypeScript interfaces
/**
 * @typedef {Object} DashboardStats
 * @property {number} totalBudget - Total budget amount
 * @property {number} totalRealization - Total realization amount
 * @property {number} realizationRate - Realization percentage
 * @property {number} variance - Budget variance
 * @property {string} trend - Performance trend
 */

/**
 * @typedef {Object} RealisasiData
 * @property {string} nama_satker - Organization name
 * @property {number} total_budget - Total budget
 * @property {number} total_realization - Total realization
 * @property {number} realization_percentage - Realization percentage
 * @property {number} variance - Budget variance
 */

/**
 * @typedef {Object} MonthlyDataPoint
 * @property {string} month - Month abbreviation
 * @property {number} value - Percentage value
 */

export function useDashboard() {
  const { isAuthenticated } = useAuth()

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

  // Computed properties
  const monthlyStats = computed(() => {
    if (!realisasiBulan.value.length) {
      return {
        total_budget: 0,
        total_realisasi: 0,
        realization_percentage: 0,
        variance: 0
      }
    }

    const stats = {
      total_budget: realisasiBulan.value.reduce((sum, item) => sum + (item.total_budget || 0), 0),
      total_realisasi: realisasiBulan.value.reduce((sum, item) => sum + (item.total_realisasi || 0), 0),
      realization_percentage: 0,
      variance: 0
    }

    if (stats.total_budget > 0) {
      stats.realization_percentage = (stats.total_realisasi / stats.total_budget) * 100
      stats.variance = stats.total_budget - stats.total_realisasi
    }

    return stats
  })

  const yearlyStats = computed(() => {
    if (!realisasiTahun.value.length) {
      return {
        total_budget: 0,
        total_realisasi: 0,
        realization_percentage: 0,
        variance: 0
      }
    }

    const stats = {
      total_budget: realisasiTahun.value.reduce((sum, item) => sum + (item.total_budget || 0), 0),
      total_realisasi: realisasiTahun.value.reduce((sum, item) => sum + (item.total_realisasi || 0), 0),
      realization_percentage: 0,
      variance: 0
    }

    if (stats.total_budget > 0) {
      stats.realization_percentage = (stats.total_realisasi / stats.total_budget) * 100
      stats.variance = stats.total_budget - stats.total_realisasi
    }

    return stats
  })

  const filteredData = computed(() => {
    let filtered = realisasiBulan.value
    if (selectedSatker.value > 0) {
      filtered = filtered.filter((item) => item.idsatker === selectedSatker.value)
    }

    return filtered
  })

  // Utility functions
  const formatCurrency = (value) => {
    if (!value) return 'Rp 0'

    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(value)
  }

  const formatPercentage = (value) => {
    return `${Math.round(value || 0)}%`
  }

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

  // Data processing functions
  const processRealisasiData = (data) => {
    if (!Array.isArray(data)) return []

    return data.map((item) => ({
      id: item.id || Math.random(),
      nama_satker: item.nama_satker || 'Unknown',
      total_budget: parseFloat(item.total_budget) || 0,
      total_realisasi: parseFloat(item.total_realisasi) || 0,
      realization_percentage: parseFloat(item.realization_percentage) || 0,
      variance: parseFloat(item.total_budget) - parseFloat(item.total_realisasi) || 0,
      idsatker: item.idsatker || 0
    }))
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

  // Filter management
  const updateFilters = (newFilters) => {
    if (newFilters.year !== undefined) selectedYear.value = newFilters.year
    if (newFilters.month !== undefined) selectedMonth.value = newFilters.month
    if (newFilters.satker !== undefined) selectedSatker.value = newFilters.satker
    if (newFilters.interval !== undefined) {
      refreshInterval.value = newFilters.interval
      startAutoRefresh()
    }

    // Refresh data with new filters
    refreshData()
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

    // Computed
    monthlyStats,
    yearlyStats,
    filteredData,

    // Actions
    refreshData,
    retryFetch,
    updateFilters,

    // Utilities
    formatCurrency,
    formatPercentage,

    // Lifecycle control
    startAutoRefresh,
    stopAutoRefresh
  }
}
