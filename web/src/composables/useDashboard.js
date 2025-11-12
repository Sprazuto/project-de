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
    articles: true,
    rankingsOpd: true,
    rankingsKecamatan: true
  })

  const error = ref({
    bulan: null,
    tahun: null,
    perbulan: null,
    articles: null,
    rankingsOpd: null,
    rankingsKecamatan: null
  })

  const realisasiBulan = ref([])
  const realisasiTahun = ref([])
  const realisasiPerbulan = ref(null)
  const articles = ref([])

  // Peringkat Kinerja state (alias-based from /sijagur/peringkat-kinerja)
  const rankingsOpd = ref([])
  const rankingsKecamatan = ref([])

  // Ranking filters (shared for both OPD & Kecamatan APIs)
  const rankingDimension = ref('kumulatif')
  const rankingCategory = ref('all')

  // Filter state
  const selectedYear = ref(new Date().getFullYear())
  const selectedMonth = ref(new Date().getMonth() + 1)
  const selectedSatker = ref(0)

  // Auto-refresh interval in milliseconds (10 seconds)
  // Used by startAutoRefresh to periodically refresh dashboard data.
  const refreshInterval = ref(0) // 1000 : 1sec

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

      // Normalize response structure for bulan
      const normalizedResponse = response && response.data && response.data.results ? response.data : response
      const processedData = processRealisasiBulanData(normalizedResponse || {})
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

      // Normalize response structure for tahun
      const normalizedResponse = response && response.data && response.data.results ? response.data : response
      realisasiTahun.value = processRealisasiTahunData(normalizedResponse || {})
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
      // - Axios: response.data = { results: [{ data: [...], meta: {...} }] }
      // - Direct: response = { results: [{ data: [...], meta: {...} }] }
      const d = response && response.data && response.data.results ? response.data : response

      if (!d || !Array.isArray(d.results) || !d.results[0] || !Array.isArray(d.results[0].data)) {
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
      d.results[0].data.forEach((categoryData) => {
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

  // Map raw ranking row -> CardRankings item (includes formatted values when provided by backend)
  const mapRankingRowToCard = (row) => ({
    name: row.nama_opd || `OPD ${row.id}`,
    total_score: Number(row.score_total) || 0,
    // formatted overall score (e.g., "91.39%")
    total_score_formatted: row.score_total_formatted + '%' || null,
    score_status: row.score_status || null,
    categories: [
      {
        key: 'barjas',
        title: 'Realisasi Barjas',
        subtitle: 'Capaian Tahunan',
        percentage: Number(row.score_barjas) || 0,
        // use backend formatted value if present
        formatted: row.score_barjas_formatted + '%' || null,
        icon: 'LayersIcon'
      },
      {
        key: 'fisik',
        title: 'Realisasi Fisik',
        subtitle: 'Capaian Tahunan',
        percentage: Number(row.score_fisik) || 0,
        formatted: row.score_fisik_formatted + '%' || null,
        icon: 'MapPinIcon'
      },
      {
        key: 'anggaran',
        title: 'Realisasi Anggaran',
        subtitle: 'Capaian Tahunan',
        percentage: Number(row.score_anggaran) || 0,
        formatted: row.score_anggaran_formatted + '%' || null,
        icon: 'ShoppingBagIcon'
      },
      {
        key: 'kinerja',
        title: 'Realisasi Kinerja',
        subtitle: 'Capaian Tahunan',
        percentage: Number(row.score_kinerja) || 0,
        formatted: row.score_kinerja_formatted + '%' || null,
        icon: 'SlidersIcon'
      }
    ]
  })

  // Fetch Peringkat Kinerja OPD (scope=skpd) using single backend endpoint with scope filter.
  const fetchRankingsOpd = async () => {
    loading.value.rankingsOpd = true
    error.value.rankingsOpd = null

    try {
      const params = new URLSearchParams({
        year: String(selectedYear.value || new Date().getFullYear()),
        month: String(selectedMonth.value || new Date().getMonth() + 1),
        idsatker: String(selectedSatker.value || 0),
        category: rankingCategory.value || 'all',
        dimension: rankingDimension.value || 'kumulatif',
        scope: 'skpd'
        // No client-side page/pageSize: backend should return complete result set for this scope
      })

      const response = await $api(`/sijagur/peringkat-kinerja?${params.toString()}`)
      const payload = response && response.status && Array.isArray(response.data) ? response : response

      if (!payload || payload.status !== 'success' || !Array.isArray(payload.data)) {
        console.error('Unexpected peringkat-kinerja OPD (scope=skpd) payload:', payload)
        throw new Error('Invalid peringkat-kinerja OPD (scope=skpd) response')
      }

      rankingsOpd.value = payload.data.map(mapRankingRowToCard)
    } catch (err) {
      console.error('Error fetching peringkat kinerja OPD (scope=skpd) rankings:', err)
      error.value.rankingsOpd = err.message || 'Failed to fetch peringkat kinerja OPD'
      rankingsOpd.value = []
    } finally {
      loading.value.rankingsOpd = false
    }
  }

  // Fetch Peringkat Kinerja Kecamatan (scope=kecamatan) using same endpoint with scope filter.
  const fetchRankingsKecamatan = async () => {
    loading.value.rankingsKecamatan = true
    error.value.rankingsKecamatan = null

    try {
      const params = new URLSearchParams({
        year: String(selectedYear.value || new Date().getFullYear()),
        month: String(selectedMonth.value || new Date().getMonth() + 1),
        category: rankingCategory.value || 'all',
        dimension: rankingDimension.value || 'kumulatif',
        scope: 'kecamatan'
      })

      const response = await $api(`/sijagur/peringkat-kinerja?${params.toString()}`)
      const payload = response && response.status && Array.isArray(response.data) ? response : response

      if (!payload || payload.status !== 'success' || !Array.isArray(payload.data)) {
        console.error('Unexpected peringkat-kinerja Kecamatan (scope=kecamatan) payload:', payload)
        throw new Error('Invalid peringkat-kinerja Kecamatan (scope=kecamatan) response')
      }

      rankingsKecamatan.value = payload.data.map(mapRankingRowToCard)
    } catch (err) {
      console.error('Error fetching peringkat kinerja Kecamatan (scope=kecamatan) rankings:', err)
      error.value.rankingsKecamatan = err.message || 'Failed to fetch peringkat kinerja Kecamatan'
      rankingsKecamatan.value = []
    } finally {
      loading.value.rankingsKecamatan = false
    }
  }

  // Initialization and auto-refresh
  let refreshTimer

  const refreshData = async () => {
    if (!isAuthenticated.value) return

    // Start loading immediately so all sections (including perbulan) show skeletons
    loading.value.bulan = true
    loading.value.tahun = true
    loading.value.perbulan = true
    loading.value.rankings = true
    loading.value.articles = true

    try {
      // Brief defer so initial skeleton DOM can paint before heavy work
      // await new Promise((resolve) => setTimeout(resolve, 3000))

      await Promise.allSettled([
        fetchRealisasiBulan(),
        fetchRealisasiTahun(),
        fetchRealisasiPerbulan(),
        // fetchArticles(),
        fetchRankingsOpd(),
        fetchRankingsKecamatan()
      ])
    } catch (err) {
      console.error('Error refreshing dashboard data:', err)
    } finally {
      // Let each fetch* manage its own loading flag.
      // Do NOT forcibly override loading.value.perbulan here,
      // otherwise CardRealisasiPerbulanSection skeleton never becomes visible.
      loading.value.bulan = false
      loading.value.tahun = false
      loading.value.perbulan = false
      loading.value.articles = false
      // Rankings loading flags are finalized inside fetchRankingsOpd/Kecamatan
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
      case 'perbulan':
        await fetchRealisasiPerbulan()
        break
      case 'rankingsOpd':
        await fetchRankingsOpd()
        break
      case 'rankingsKecamatan':
        await fetchRankingsKecamatan()
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
    realisasiPerbulan,
    articles,
    rankingsOpd,
    rankingsKecamatan,
    selectedYear,
    selectedMonth,
    selectedSatker,
    refreshInterval,
    rankingDimension,
    rankingCategory,

    // Actions
    refreshData,
    retryFetch,
    fetchRankingsOpd,
    fetchRankingsKecamatan,

    // Lifecycle control
    startAutoRefresh,
    stopAutoRefresh
  }
}
