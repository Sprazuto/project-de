<script setup>
import { computed } from 'vue'

// Props
const props = defineProps({
  realisasiPerbulan: {
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
import CardRealisasiPerbulan from './CardRealisasiPerbulan.vue'
</script>

<template>
  <div>
    <!-- Loading State -->
    <VRow v-if="loading" class="match-height">
      <VCol cols="12">
        <VCard elevation="2" rounded="xl">
          <VCardItem class="pb-2">
            <VCardTitle class="d-flex align-start justify-space-between pb-3">
              <VSkeletonLoader type="text" width="70%" />
              <VSkeletonLoader type="avatar" size="18" />
            </VCardTitle>
            <VCardSubtitle>
              <VSkeletonLoader type="text" width="50%" />
            </VCardSubtitle>
          </VCardItem>

          <VCardItem class="px-4 pb-2">
            <VRow class="mx-0 mb-2">
              <!-- Large card skeleton for current month -->
              <VCol cols="12" md="4" class="px-1 mb-1">
                <VCard elevation="0" rounded="xl" height="200">
                  <VCardItem class="text-center py-2 d-flex flex-column justify-center h-100">
                    <VSkeletonLoader type="text" width="60%" class="mb-2" />
                    <VSkeletonLoader type="text" width="80%" height="32" class="mb-1" />
                    <VSkeletonLoader type="text" width="50%" height="16" />
                  </VCardItem>
                </VCard>
              </VCol>

              <!-- Grid of 12 smaller cards skeleton -->
              <VCol cols="12" md="8" class="px-0">
                <VRow class="mx-0">
                  <VCol v-for="n in 12" :key="`skeleton-month-${n}`" cols="3" class="pl-1 pr-0 mb-1">
                    <VCard elevation="0" rounded="xl" height="100">
                      <VCardItem class="text-center py-1 d-flex flex-column justify-center h-100">
                        <VSkeletonLoader type="text" width="70%" height="14" class="mb-1" />
                        <VSkeletonLoader type="text" width="60%" height="16" />
                      </VCardItem>
                    </VCard>
                  </VCol>
                </VRow>
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
    <VRow v-else class="match-height" role="list" aria-label="Monthly realization data">
      <VCol cols="12" role="listitem">
        <CardRealisasiPerbulan
          v-for="(card, index) in realisasiPerbulan"
          :key="`perbulan-${index}`"
          :title="card.title"
          :subtitle="card.subtitle"
          :hint-title="card.hintTitle"
          :hint-description="card.hintDescription"
          :current-month="card.currentMonth"
          :monthly-data="card.monthlyData"
        />
      </VCol>
    </VRow>
  </div>
</template>
