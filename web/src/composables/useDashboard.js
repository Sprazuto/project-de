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
    perbulan: true,
    articles: true
  })

  const error = ref({
    bulan: null,
    tahun: null,
    perbulan: null,
    articles: null
  })

  const realisasiBulan = ref([])
  const realisasiTahun = ref([])
  const realisasiPerbulan = ref(null)
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

  const fetchRealisasiPerbulan = async (params = {}) => {
    loading.value.perbulan = true
    error.value.perbulan = null

    try {
      const queryParams = new URLSearchParams({
        tahun: Number(selectedYear.value),
        idsatker: Number(selectedSatker.value),
        ...params
      })

      const response = await $api(`/realisasi-perbulan?${queryParams}`)

      // Normalize Axios or direct-fetch style response:
      // - Axios: response.data = { data: [...], meta: {...} }
      // - Direct: response = { data: [...], meta: {...} }
      const d = response && response.data && response.data.data ? response.data : response

      if (!d || !Array.isArray(d.data)) {
        throw new Error('Invalid response format for realisasi perbulan')
      }

      const mapCategory = (categoryData, hintTitle, hintDescription) => {
        if (!categoryData || !categoryData.current_month || !Array.isArray(categoryData.monthly)) return null

        return {
          key: categoryData.category,
          hintTitle,
          hintDescription,
          currentMonth: {
            month: categoryData.current_month.month,
            value: Number(categoryData.current_month.value) || 0,
            value_formatted: categoryData.current_month.value_formatted || '0',
            realisasi: Number(categoryData.current_month.realisasi) || 0,
            target: Number(categoryData.current_month.target) || 100,
            realisasi_formatted: categoryData.current_month.realisasi_formatted || null,
            target_formatted: categoryData.current_month.target_formatted || null
          },
          monthlyData: categoryData.monthly.map((m) => ({
            month: m.month,
            value: Number(m.value) || 0,
            value_formatted: m.value_formatted || '0',
            realisasi: Number(m.realisasi) || 0,
            target: Number(m.target) || 100,
            realisasi_formatted: m.realisasi_formatted || null,
            target_formatted: m.target_formatted || null
          }))
        }
      }

      const categories = {}
      d.data.forEach((categoryData) => {
        const hintTitles = {
          barjas: 'Informasi',
          fisik: 'Informasi',
          anggaran: 'Informasi',
          kinerja: 'Informasi'
        }
        const hintDescriptions = {
          barjas: 'Data realisasi bulanan untuk kategori Barang dan Jasa.',
          fisik: 'Data realisasi bulanan untuk kategori Fisik.',
          anggaran: 'Data realisasi bulanan untuk kategori Anggaran.',
          kinerja: 'Data realisasi bulanan untuk kategori Kinerja.'
        }

        categories[categoryData.category] = mapCategory(
          categoryData,
          hintTitles[categoryData.category] || 'Informasi',
          hintDescriptions[categoryData.category] || ''
        )
      })

      realisasiPerbulan.value = categories
    } catch (err) {
      console.error('Error fetching realisasi perbulan data:', err)
      error.value.perbulan = err.message || 'Failed to fetch realisasi perbulan data'
      realisasiPerbulan.value = null
    } finally {
      loading.value.perbulan = false
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
      await Promise.allSettled([
        fetchRealisasiBulan(),
        fetchRealisasiTahun(),
        fetchRealisasiPerbulan()
        // fetchArticles()
      ])
    } catch (err) {
      console.error('Error refreshing dashboard data:', err)
    } finally {
      // Reset loading states after fetch completes
      loading.value.bulan = false
      loading.value.tahun = false
      loading.value.perbulan = false
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
      case 'perbulan':
        await fetchRealisasiPerbulan()
        break
    }
  }

  return {
    // State
    loading,
    error,
    realisasiBulan,
    realisasiTahun,
    realisasiPerbulan,
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
