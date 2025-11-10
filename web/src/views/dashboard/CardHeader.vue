<script setup>
const props = defineProps({
  title: {
    type: String,
    default: 'Title'
  },
  subtitle: {
    type: String,
    default: ''
  },
  icon: {
    type: String,
    default: 'mdi-chart-bar'
  },
  iconSize: {
    type: [String, Number],
    default: 28
  },
  titleClass: {
    type: String,
    default: 'text-gradient-primary mb-1'
  },
  iconClass: {
    type: String,
    default: 'mr-1 text-primary'
  },
  hintTitle: {
    type: String,
    default: ''
  },
  hintDescription: {
    type: String,
    default: ''
  }
})
</script>

<template>
  <VCard
    rounded="xl"
    class="match-height"
    :style="{ 'border-block-start': `4px solid rgba(var(--v-theme-primary), var(--v-disabled-opacity))` }"
  >
    <VCardText>
      <VRow>
        <VCol cols="12">
          <div class="d-flex flex-column align-center h-100 justify-center">
            <div class="d-flex flex-row align-center gap-2">
              <VAvatar variant="tonal" size="48" rounded="xl" color="primary">
                <VIcon :icon="icon" size="36" />
              </VAvatar>
              <h3 class="text-primary font-weight-medium">{{ title }}</h3>
            </div>
            <p class="text-body-1 text-center mb-2">
              {{ subtitle }}
            </p>
            <!-- <VBtn>Get Premium Courses</VBtn> -->
          </div>
        </VCol>
      </VRow>
    </VCardText>
  </VCard>
</template>
