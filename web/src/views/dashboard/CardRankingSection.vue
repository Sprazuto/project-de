<script setup>
import CardRankings from './CardRankings.vue'

const props = defineProps({
  title: {
    type: String,
    default: 'Peringkat Kinerja'
  },
  subtitle: {
    type: String,
    default: 'Top performers across categories'
  },
  rankings: {
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
</script>

<template>
  <div>
    <!-- Loading State:
         Match other dashboard cards:
         - Single rounded card
         - Neutral tonal background
         - Skeletons laid out to mirror final ranking structure
    -->
    <!-- Loading skeleton: simulate actual CardRankings layout using Vuetify-supported skeleton types -->
    <VRow v-if="loading" class="match-height">
      <VCol cols="12">
        <VCard elevation="0" rounded="xl" class="pa-4 mb-3">
          <VCardText class="py-2">
            <!-- Header skeleton: title + secondary control area -->
            <VCardItem class="pb-4">
              <VCardTitle class="d-flex align-center justify-center">
                <VSkeletonLoader type="avatar" size="20" class="bg-grey-lighten-3" />
                <VSkeletonLoader type="sentences" width="35%" class="bg-grey-lighten-3" />
              </VCardTitle>
            </VCardItem>

            <!-- List skeleton: each row mirrors .ranking-item structure -->
            <div v-for="i in 5" :key="i" class="d-flex align-center justify-space-between py-2 px-1 mb-1">
              <!-- Left cluster: rank badge + status/score + name -->
              <div class="d-flex align-center flex-grow-1">
                <!-- Rank badge -->
                <VSkeletonLoader type="avatar" class="mr-3" />

                <!-- Status + total score -->
                <VSkeletonLoader type="sentences" class="mr-3" style="width: 70px" />

                <!-- Instance name -->
                <VSkeletonLoader type="text" class="mr-3 flex-grow-1" style="max-width: 320px" />
              </div>

              <!-- Right cluster: 4 compact category cards -->
              <div class="d-flex align-center" style="gap: 8px">
                <div v-for="j in 4" :key="j" class="d-flex flex-column justify-center" style="width: 150px">
                  <VCard
                    elevation="0"
                    rounded="xl"
                    height="100"
                    min-height="80"
                    class="d-flex flex-column justify-center align-center px-1"
                    style="top: -20px"
                  >
                    <!-- Month label placeholder -->
                    <VSkeletonLoader type="text" width="80%" height="9" class="mb-2 bg-grey-lighten-3" />
                    <!-- Percentage placeholder -->
                    <VSkeletonLoader type="text" width="85%" height="14" class="mb-1 bg-grey-lighten-3" />
                    <!-- Realisasi/target line placeholder -->
                    <VSkeletonLoader type="text" width="65%" height="8" class="bg-grey-lighten-3" />
                  </VCard>
                </div>
              </div>
            </div>
          </VCardText>
        </VCard>
      </VCol>
    </VRow>

    <!-- Error State: local to this section -->
    <VAlert v-else-if="error" type="error" variant="tonal" density="compact" class="mb-2">
      {{ error }}
    </VAlert>

    <!-- Empty State -->
    <VCard v-else-if="!rankings || !rankings.length" elevation="0" rounded="xl" class="pa-3">
      <div class="text-caption text-disabled">Tidak ada data peringkat kinerja.</div>
    </VCard>

    <!-- Data State: delegate rendering to CardRankings to preserve behavior -->
    <div v-else>
      <CardRankings :title="title" :subtitle="subtitle" :rankings="rankings" />
    </div>
  </div>
</template>
