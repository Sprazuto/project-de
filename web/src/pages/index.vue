<script setup>
// Vue 3 Composition API with seamless useAuth and useApi integration
import { ref, computed, onMounted } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { useApi } from '@/composables/useApi'
import { useDashboard } from '@/composables/useDashboard'

// Import required Vuetify components
import { VCard, VCardText, VRow, VCol } from 'vuetify/components'

// Fixed component imports - using the correct paths
import CardHeader from '@/views/dashboard/CardHeader.vue'
import CardRealisasiBulanSection from '@/views/dashboard/CardRealisasiBulanSection.vue'
import CardRealisasiTahunSection from '@/views/dashboard/CardRealisasiTahunSection.vue'
import CardRealisasiPerbulanSection from '@/views/dashboard/CardRealisasiPerbulanSection.vue'
import CardRealisasiBulan from '@/views/dashboard/CardRealisasiBulan.vue'
import CardRankings from '@/views/dashboard/CardRankings.vue'

// Page meta
definePage({
  meta: {
    requiresAuth: true,
    title: 'Dashboard - SIJAGUR',
    description: 'Sistem Informasi Realisasi Anggaran'
  }
})

// Composables for seamless API integration
const { currentUser, isAuthenticated } = useAuth()
const dashboard = useDashboard()

// Processed realization data computed property
const processedRealisasiBulan = computed(() => {
  return dashboard.realisasiBulan.value || []
})

const processedRealisasiTahun = computed(() => {
  return dashboard.realisasiTahun.value || []
})

// Monthly data for perbulan sections - matching Home.vue structure
const monthlyData = ref({
  barjas: {
    currentMonth: { month: 'October', value: 85 },
    data: [
      { month: 'Jan', value: 75 },
      { month: 'Feb', value: 80 },
      { month: 'Mar', value: 70 },
      { month: 'Apr', value: 85 },
      { month: 'May', value: 90 },
      { month: 'Jun', value: 78 },
      { month: 'Jul', value: 82 },
      { month: 'Aug', value: 88 },
      { month: 'Sep', value: 92 },
      { month: 'Oct', value: 85 },
      { month: 'Nov', value: 0 },
      { month: 'Dec', value: 95 }
    ]
  },
  fisik: {
    currentMonth: { month: 'October', value: 20 },
    data: [
      { month: 'Jan', value: 15 },
      { month: 'Feb', value: 18 },
      { month: 'Mar', value: 12 },
      { month: 'Apr', value: 22 },
      { month: 'May', value: 25 },
      { month: 'Jun', value: 19 },
      { month: 'Jul', value: 21 },
      { month: 'Aug', value: 23 },
      { month: 'Sep', value: 26 },
      { month: 'Oct', value: 20 },
      { month: 'Nov', value: 0 },
      { month: 'Dec', value: 0 }
    ]
  },
  anggaran: {
    currentMonth: { month: 'October', value: 60 },
    data: [
      { month: 'Jan', value: 55 },
      { month: 'Feb', value: 58 },
      { month: 'Mar', value: 52 },
      { month: 'Apr', value: 62 },
      { month: 'May', value: 65 },
      { month: 'Jun', value: 0 },
      { month: 'Jul', value: 61 },
      { month: 'Aug', value: 63 },
      { month: 'Sep', value: 66 },
      { month: 'Oct', value: 60 },
      { month: 'Nov', value: 0 },
      { month: 'Dec', value: 0 }
    ]
  },
  kinerja: {
    currentMonth: { month: 'October', value: 35 },
    data: [
      { month: 'Jan', value: 30 },
      { month: 'Feb', value: 33 },
      { month: 'Mar', value: 27 },
      { month: 'Apr', value: 37 },
      { month: 'May', value: 40 },
      { month: 'Jun', value: 34 },
      { month: 'Jul', value: 36 },
      { month: 'Aug', value: 38 },
      { month: 'Sep', value: 41 },
      { month: 'Oct', value: 35 },
      { month: 'Nov', value: 0 },
      { month: 'Dec', value: 0 }
    ]
  }
})

// Rankings computed property - using legacy mock data structure
const rankings = computed(() => {
  // Mock data matching legacy structure exactly
  return [
    {
      name: 'DINAS PERIKANAN DAN PETERNAKAN',
      total_score: 85.75,
      categories: [
        {
          title: 'Realisasi Barjas',
          subtitle: 'Capaian Tahunan',
          percentage: 85,
          icon: 'LayersIcon'
        },
        {
          title: 'Realisasi Fisik',
          subtitle: 'Capaian Tahunan',
          percentage: 78,
          icon: 'MapPinIcon'
        },
        {
          title: 'Realisasi Anggaran',
          subtitle: 'Capaian Tahunan',
          percentage: 74,
          icon: 'ShoppingBagIcon'
        },
        {
          title: 'Realisasi Kinerja',
          subtitle: 'Capaian Tahunan',
          percentage: 45,
          icon: 'SlidersIcon'
        }
      ]
    },
    {
      name: 'DINAS PENANAMAN MODAL DAN PELAYANAN TERPADU SATU PINTU',
      total_score: 82.75,
      categories: [
        {
          title: 'Realisasi Barjas',
          subtitle: 'Capaian Tahunan',
          percentage: 29,
          icon: 'LayersIcon'
        },
        {
          title: 'Realisasi Fisik',
          subtitle: 'Capaian Tahunan',
          percentage: 75,
          icon: 'MapPinIcon'
        },
        {
          title: 'Realisasi Anggaran',
          subtitle: 'Capaian Tahunan',
          percentage: 62,
          icon: 'ShoppingBagIcon'
        },
        {
          title: 'Realisasi Kinerja',
          subtitle: 'Capaian Tahunan',
          percentage: 85,
          icon: 'SlidersIcon'
        }
      ]
    },
    {
      name: 'DINAS KEPENDUDUKAN DAN PENCATATAN SIPIL',
      total_score: 79.75,
      categories: [
        {
          title: 'Realisasi Barjas',
          subtitle: 'Capaian Tahunan',
          percentage: 79,
          icon: 'LayersIcon'
        },
        {
          title: 'Realisasi Fisik',
          subtitle: 'Capaian Tahunan',
          percentage: 40,
          icon: 'MapPinIcon'
        },
        {
          title: 'Realisasi Anggaran',
          subtitle: 'Capaian Tahunan',
          percentage: 86,
          icon: 'ShoppingBagIcon'
        },
        {
          title: 'Realisasi Kinerja',
          subtitle: 'Capaian Tahunan',
          percentage: 55,
          icon: 'SlidersIcon'
        }
      ]
    },
    {
      name: 'DINAS PERHUBUNGAN',
      total_score: 76.75,
      categories: [
        {
          title: 'Realisasi Barjas',
          subtitle: 'Capaian Tahunan',
          percentage: 76,
          icon: 'LayersIcon'
        },
        {
          title: 'Realisasi Fisik',
          subtitle: 'Capaian Tahunan',
          percentage: 69,
          icon: 'MapPinIcon'
        },
        {
          title: 'Realisasi Anggaran',
          subtitle: 'Capaian Tahunan',
          percentage: 99,
          icon: 'ShoppingBagIcon'
        },
        {
          title: 'Realisasi Kinerja',
          subtitle: 'Capaian Tahunan',
          percentage: 79,
          icon: 'SlidersIcon'
        }
      ]
    },
    {
      name: 'DINAS PENDIDIKAN',
      total_score: 73.75,
      categories: [
        {
          title: 'Realisasi Barjas',
          subtitle: 'Capaian Tahunan',
          percentage: 99,
          icon: 'LayersIcon'
        },
        {
          title: 'Realisasi Fisik',
          subtitle: 'Capaian Tahunan',
          percentage: 66,
          icon: 'MapPinIcon'
        },
        {
          title: 'Realisasi Anggaran',
          subtitle: 'Capaian Tahunan',
          percentage: 80,
          icon: 'ShoppingBagIcon'
        },
        {
          title: 'Realisasi Kinerja',
          subtitle: 'Capaian Tahunan',
          percentage: 76,
          icon: 'SlidersIcon'
        }
      ]
    }
  ]
})

// Initialize data - the useDashboard composable handles this automatically
onMounted(async () => {
  // Let Vue render first (skeleton visible)
  await new Promise((resolve) => setTimeout(resolve, 0))
})
</script>

<template>
  <div>
    <!-- Realisasi Bulan Section using processed data from useDashboard -->
    <CardRealisasiBulanSection :realisasi-bulan="processedRealisasiBulan" :loading="dashboard.loading.value.bulan" :error="dashboard.error.bulan" />

    <!-- Realisasi Tahun Section using processed data from useDashboard -->
    <CardRealisasiTahunSection :realisasi-tahun="processedRealisasiTahun" :loading="dashboard.loading.value.tahun" :error="dashboard.error.tahun" />

    <!-- Realisasi Perbulan Barjas Section -->
    <CardHeader title="Realisasi Perbulan Barjas" icon="tabler-stack-pop" />
    <CardRealisasiPerbulanSection
      :realisasi-perbulan="[
        {
          title: 'Realisasi Perbulan Barjas',
          subtitle: 'Monthly Barjas Realization',
          hintTitle: 'Information',
          hintDescription: 'Monthly realization data for Barjas category',
          currentMonth: monthlyData.barjas.currentMonth,
          monthlyData: monthlyData.barjas.data
        }
      ]"
      :loading="false"
      :error="null"
    />

    <!-- Realisasi Perbulan Fisik Section -->
    <CardHeader title="Realisasi Perbulan Fisik" icon="tabler-map-pin" />
    <CardRealisasiPerbulanSection
      :realisasi-perbulan="[
        {
          title: 'Realisasi Perbulan Fisik',
          subtitle: 'Monthly Physical Realization',
          hintTitle: 'Information',
          hintDescription: 'Monthly realization data for Physical category',
          currentMonth: monthlyData.fisik.currentMonth,
          monthlyData: monthlyData.fisik.data
        }
      ]"
      :loading="false"
      :error="null"
    />

    <!-- Realisasi Perbulan Anggaran Section -->
    <CardHeader title="Realisasi Perbulan Anggaran" icon="tabler-shopping-bag" />
    <CardRealisasiPerbulanSection
      :realisasi-perbulan="[
        {
          title: 'Realisasi Perbulan Anggaran',
          subtitle: 'Monthly Budget Realization',
          hintTitle: 'Information',
          hintDescription: 'Monthly realization data for Budget category',
          currentMonth: monthlyData.anggaran.currentMonth,
          monthlyData: monthlyData.anggaran.data
        }
      ]"
      :loading="false"
      :error="null"
    />

    <!-- Realisasi Perbulan Kinerja Section -->
    <CardHeader title="Realisasi Perbulan Kinerja" icon="tabler-adjustments-alt" />
    <CardRealisasiPerbulanSection
      :realisasi-perbulan="[
        {
          title: 'Realisasi Perbulan Kinerja',
          subtitle: 'Monthly Performance Realization',
          hintTitle: 'Information',
          hintDescription: 'Monthly realization data for Performance category',
          currentMonth: monthlyData.kinerja.currentMonth,
          monthlyData: monthlyData.kinerja.data
        }
      ]"
      :loading="false"
      :error="null"
    />

    <!-- Card Header for Rankings -->
    <CardHeader title="Peringkat Kinerja" subtitle="Organisasi Perangkat Daerah" icon="tabler-award" />

    <VRow class="match-height">
      <VCol cols="12" lg="12">
        <CardRankings :rankings="rankings" />
      </VCol>
    </VRow>
  </div>
</template>
