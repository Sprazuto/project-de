<script setup>
import { computed } from 'vue'

const props = defineProps({
  realisasiTahun: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: null
  }
})

// Import the card component
import CardRealisasiTahun from './CardRealisasiTahun.vue'
</script>

<template>
  <div>
    <!-- Loading State -->
    <VRow v-if="loading" class="match-height">
      <VCol v-for="n in 4" :key="`skeleton-${n}`" cols="12" sm="6" md="3" lg="3" xl="3">
        <VCard elevation="0" rounded="xl">
          <VCardItem class="pb-sm-0">
            <VCardTitle class="d-flex align-start justify-space-between pb-5">
              <VSkeletonLoader type="text" width="70%" />
              <VSkeletonLoader type="avatar" size="18" />
            </VCardTitle>
            <VCardSubtitle class="d-flex justify-space-between align-items-center mb-2 w-100">
              <VSkeletonLoader type="text" width="40%" />
              <VSkeletonLoader type="text" width="20%" />
            </VCardSubtitle>
            <VSkeletonLoader type="image" height="6" width="100%" class="mb-2" />
          </VCardItem>

          <VCardItem class="text-center pt-0 pb-1">
            <VRow class="mx-0 mt-0">
              <VCol cols="6" class="d-flex align-items-center flex-column">
                <VSkeletonLoader type="text" width="80%" class="mb-1" />
                <VSkeletonLoader type="text" width="90%" height="16" />
              </VCol>
              <VCol cols="6" class="d-flex align-items-center flex-column">
                <VSkeletonLoader type="text" width="60%" class="mb-1" />
                <VSkeletonLoader type="text" width="90%" height="16" />
              </VCol>
            </VRow>
            <VRow class="mx-0 mt-0">
              <VCol cols="6" class="d-flex align-items-center flex-column">
                <VSkeletonLoader type="text" width="70%" class="mb-1" />
                <VSkeletonLoader type="text" width="90%" height="16" />
              </VCol>
              <VCol cols="6" class="d-flex align-items-center flex-column">
                <VSkeletonLoader type="text" width="60%" class="mb-1" />
                <VSkeletonLoader type="text" width="90%" height="16" />
              </VCol>
            </VRow>
          </VCardItem>
        </VCard>
      </VCol>
    </VRow>

    <!-- Error State -->
    <VAlert v-else-if="error" type="error" variant="tonal" prominent role="alert" aria-live="assertive">
      <template #prepend>
        <VIcon icon="mdi-alert-circle" />
      </template>
      {{ error }}
    </VAlert>

    <!-- Content State -->
    <VRow v-else class="match-height" role="list" aria-label="Yearly realization data">
      <VCol
        v-for="(card, index) in realisasiTahun"
        :key="`tahun-${index}`"
        cols="12"
        sm="6"
        md="3"
        lg="3"
        xl="3"
        role="listitem"
      >
        <CardRealisasiTahun
          :title="card.title"
          :subtitle="card.subtitle"
          :hint-title="card.hintTitle"
          :hint-description="card.hintDescription"
          :items="card.items"
          :progress="card.progress"
          :capaian="card.capaian"
          :color="card.color"
          :layout="card.layout || 'columns'"
        />
      </VCol>
    </VRow>
  </div>
</template>
