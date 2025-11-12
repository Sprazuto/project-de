<script setup>
// Vue 3 Composition API with seamless useAuth and useApi integration
import { computed, onMounted } from 'vue'
import { useDashboard } from '@/composables/useDashboard'

// Fixed component imports - using the correct paths
import CardHeader from '@/views/dashboard/CardHeader.vue'
import CardRealisasiBulanSection from '@/views/dashboard/CardRealisasiBulanSection.vue'
import CardRealisasiTahunSection from '@/views/dashboard/CardRealisasiTahunSection.vue'
import CardRealisasiPerbulanSection from '@/views/dashboard/CardRealisasiPerbulanSection.vue'
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
const dashboard = useDashboard()

// Processed realization data computed property
const processedRealisasiBulan = computed(() => {
  return dashboard.realisasiBulan.value || []
})

const processedRealisasiTahun = computed(() => {
  return dashboard.realisasiTahun.value || []
})

// Computed perbulan cards from useDashboard (driven by /realisasi-perbulan)
const perbulanBarjas = computed(() => {
  const src = dashboard.realisasiPerbulan.value?.barjas
  if (!src || !src.currentMonth || !Array.isArray(src.monthlyData)) return null

  return {
    key: 'barjas',
    title: 'Realisasi Perbulan Barjas',
    icon: 'tabler-stack-pop',
    hintTitle: 'Informasi',
    hintDescription: 'Data realisasi bulanan untuk kategori Barang dan Jasa.',
    currentMonth: src.currentMonth,
    monthlyData: src.monthlyData
  }
})

const perbulanFisik = computed(() => {
  const src = dashboard.realisasiPerbulan.value?.fisik
  if (!src || !src.currentMonth || !Array.isArray(src.monthlyData)) return null

  return {
    key: 'fisik',
    title: 'Realisasi Perbulan Fisik',
    icon: 'tabler-map-pin',
    hintTitle: 'Informasi',
    hintDescription: 'Data realisasi bulanan untuk kategori Fisik.',
    currentMonth: src.currentMonth,
    monthlyData: src.monthlyData
  }
})

const perbulanAnggaran = computed(() => {
  const src = dashboard.realisasiPerbulan.value?.anggaran
  if (!src || !src.currentMonth || !Array.isArray(src.monthlyData)) return null

  return {
    key: 'anggaran',
    title: 'Realisasi Perbulan Anggaran',
    icon: 'tabler-shopping-bag',
    hintTitle: 'Informasi',
    hintDescription: 'Data realisasi bulanan untuk kategori Anggaran.',
    currentMonth: src.currentMonth,
    monthlyData: src.monthlyData
  }
})

const perbulanKinerja = computed(() => {
  const src = dashboard.realisasiPerbulan.value?.kinerja
  if (!src || !src.currentMonth || !Array.isArray(src.monthlyData)) return null

  return {
    key: 'kinerja',
    title: 'Realisasi Perbulan Kinerja',
    icon: 'tabler-adjustments-alt',
    hintTitle: 'Informasi',
    hintDescription: 'Data realisasi bulanan untuk kategori Kinerja.',
    currentMonth: src.currentMonth,
    monthlyData: src.monthlyData
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
    <CardRealisasiBulanSection
      :realisasi-bulan="processedRealisasiBulan"
      :loading="dashboard.loading.value.bulan"
      :error="dashboard.error.bulan"
    />

    <!-- Realisasi Tahun Section using processed data from useDashboard -->
    <CardRealisasiTahunSection
      :realisasi-tahun="processedRealisasiTahun"
      :loading="dashboard.loading.value.tahun"
      :error="dashboard.error.tahun"
    />

    <!-- Realisasi Perbulan Barjas -->
    <CardHeader
      v-if="!dashboard.loading.value.perbulan && perbulanBarjas"
      :title="perbulanBarjas.title"
      :icon="perbulanBarjas.icon"
    />
    <CardRealisasiPerbulanSection
      :realisasi-perbulan="perbulanBarjas ? [perbulanBarjas] : []"
      :loading="dashboard.loading.value.perbulan"
      :error="dashboard.error.perbulan"
    />

    <!-- Realisasi Perbulan Fisik -->
    <CardHeader
      v-if="!dashboard.loading.value.perbulan && perbulanFisik"
      :title="perbulanFisik.title"
      :icon="perbulanFisik.icon"
    />
    <CardRealisasiPerbulanSection
      :realisasi-perbulan="perbulanFisik ? [perbulanFisik] : []"
      :loading="dashboard.loading.value.perbulan"
      :error="dashboard.error.perbulan"
    />

    <!-- Realisasi Perbulan Anggaran -->
    <CardHeader
      v-if="!dashboard.loading.value.perbulan && perbulanAnggaran"
      :title="perbulanAnggaran.title"
      :icon="perbulanAnggaran.icon"
    />
    <CardRealisasiPerbulanSection
      :realisasi-perbulan="perbulanAnggaran ? [perbulanAnggaran] : []"
      :loading="dashboard.loading.value.perbulan"
      :error="dashboard.error.perbulan"
    />

    <!-- Realisasi Perbulan Kinerja -->
    <CardHeader
      v-if="!dashboard.loading.value.perbulan && perbulanKinerja"
      :title="perbulanKinerja.title"
      :icon="perbulanKinerja.icon"
    />
    <CardRealisasiPerbulanSection
      :realisasi-perbulan="perbulanKinerja ? [perbulanKinerja] : []"
      :loading="dashboard.loading.value.perbulan"
      :error="dashboard.error.perbulan"
    />

    <!-- Card Header for Rankings -->
    <CardHeader title="Peringkat Kinerja" subtitle="Organisasi Perangkat Daerah" icon="tabler-award" />

    <CardRankings :rankings="rankings" />

    <!-- Card Header for Rankings -->
    <CardHeader title="Peringkat Kinerja" subtitle="Kecamatan" icon="tabler-award" />

    <CardRankings :rankings="rankings" />
  </div>
</template>
