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
    <!-- Loading State: clean, minimal skeleton that mimics CardRealisasiPerbulan layout -->
    <VRow v-if="loading" class="match-height">
      <VCol cols="12">
        <!-- Outer wrapper matches real card section spacing and radius -->
        <VCard elevation="0" rounded="xl" class="pa-4">
          <!-- Header skeleton: title + optional icon space, neutral grayscale -->
          <VCardItem class="pb-4">
            <VCardTitle class="d-flex align-center justify-center">
              <VSkeletonLoader type="avatar" size="20" class="bg-grey-lighten-3" />
              <VSkeletonLoader type="text" width="35%" class="bg-grey-lighten-3" />
            </VCardTitle>
            <VCardSubtitle>
              <VSkeletonLoader type="text" width="22%" height="10" class="mt-1 bg-grey-lighten-3" />
            </VCardSubtitle>
          </VCardItem>

          <!-- Body skeleton: structure = 1 large current month + 12 small cards -->
          <VCardItem class="pt-0">
            <VRow class="mx-0">
              <!-- Large current month skeleton -->
              <VCol cols="12" sm="6" lg="4" class="px-0 my-3 py-2">
                <VCard elevation="0" rounded="xl" class="d-flex flex-column justify-center align-center px-4">
                  <!-- Title placeholder -->
                  <VSkeletonLoader type="text" width="40%" height="12" class="mb-3 bg-grey-lighten-3" />
                  <!-- Main value placeholder -->
                  <VSkeletonLoader type="text" width="55%" height="26" class="mb-2 bg-grey-lighten-3" />
                  <!-- Supporting text placeholder -->
                  <VSkeletonLoader type="text" width="60%" height="10" class="bg-grey-lighten-3" />
                </VCard>
              </VCol>

              <!-- Grid of 12 smaller month cards -->
              <VCol cols="12" sm="6" lg="8" class="px-0">
                <VRow class="mx-0" dense>
                  <VCol
                    v-for="n in 12"
                    :key="`skeleton-month-${n}`"
                    cols="4"
                    sm="3"
                    md="2"
                    lg="3"
                    xl="2"
                    class="px-0 my-3 py-0 pl-6"
                  >
                    <VCard
                      elevation="0"
                      rounded="xl"
                      height="100"
                      min-height="80"
                      class="d-flex flex-column justify-center align-center px-3"
                    >
                      <!-- Month label placeholder -->
                      <VSkeletonLoader type="text" width="60%" height="9" class="mb-2 bg-grey-lighten-3" />
                      <!-- Percentage placeholder -->
                      <VSkeletonLoader type="text" width="40%" height="14" class="mb-1 bg-grey-lighten-3" />
                      <!-- Realisasi/target line placeholder -->
                      <VSkeletonLoader type="text" width="65%" height="8" class="bg-grey-lighten-3" />
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
          :category-key="card.key"
          :hint-title="card.hintTitle"
          :hint-description="card.hintDescription"
          :current-month="card.currentMonth"
          :monthly-data="card.monthlyData"
        />
      </VCol>
    </VRow>
  </div>
</template>
