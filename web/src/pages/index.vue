<script setup>
// Vue 3 Composition API with seamless useAuth and useApi integration
import { computed, onMounted } from 'vue'
import { useDashboard } from '@/composables/useDashboard'

// Fixed component imports - using the correct paths
import CardHeader from '@/views/dashboard/CardHeader.vue'
import CardRealisasiBulanSection from '@/views/dashboard/CardRealisasiBulanSection.vue'
import CardRealisasiTahunSection from '@/views/dashboard/CardRealisasiTahunSection.vue'
import CardRealisasiPerbulanSection from '@/views/dashboard/CardRealisasiPerbulanSection.vue'
import CardRankingSection from '@/views/dashboard/CardRankingSection.vue'

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

const rankingsOPDKumulatif = computed(() => {
  const src = dashboard.rankingsOpd
  return Array.isArray(src) ? src : Array.isArray(src?.value) ? src.value : []
})

const rankingsKecamatanKumulatif = computed(() => {
  const src = dashboard.rankingsKecamatan
  return Array.isArray(src) ? src : Array.isArray(src?.value) ? src.value : []
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

    <!-- Peringkat Kinerja OPD - Kumulatif (live via useDashboard + CardRankingSection) -->
    <CardHeader
      v-if="!dashboard.loading.value.rankingsOpd && rankingsOPDKumulatif"
      title="Peringkat Kinerja"
      subtitle="Organisasi Perangkat Daerah"
      icon="tabler-award"
    />
    <CardRankingSection
      title="Peringkat Kinerja"
      subtitle="Organisasi Perangkat Daerah"
      :rankings="rankingsOPDKumulatif"
      :loading="dashboard.loading.value.rankingsOpd"
      :error="dashboard.error.rankingsOpd"
    />

    <!-- Peringkat Kinerja Kecamatan - Kumulatif (placeholder, no data yet) -->
    <CardHeader
      v-if="!dashboard.loading.value.rankingsKecamatan && rankingsKecamatanKumulatif"
      title="Peringkat Kinerja"
      subtitle="Kecamatan"
      icon="tabler-award"
    />
    <CardRankingSection
      title="Peringkat Kinerja"
      subtitle="Kecamatan"
      :rankings="rankingsKecamatanKumulatif"
      :loading="dashboard.loading.value.rankingsKecamatan"
      :error="dashboard.error.rankingsKecamatan"
    />
  </div>
</template>
