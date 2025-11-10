<script setup>
import { computed } from 'vue'
import { VCardItem } from 'vuetify/lib/components/index.mjs'

const props = defineProps({
  title: {
    type: String,
    default: 'Realisasi'
  },
  subtitle: {
    type: String,
    default: ''
  },
  hintTitle: {
    type: String,
    default: ''
  },
  hintDescription: {
    type: String,
    default: ''
  },
  items: {
    type: Array,
    default: () => []
  },
  layout: {
    type: String,
    default: 'columns'
  },
  progress: {
    type: [Number, String],
    default: 0
  },
  color: {
    type: String,
    default: ''
  }
})

// Computed properties - using only primary, secondary, error, and dark
const cardColors = computed(() => {
  const progress = Number(props.progress)
  if (progress >= 75) {
    return {
      bgColor: 'primary',
      textColor: 'text-white',
      chartColors: ['#028C86']
    }
  } else if (progress >= 50) {
    return {
      bgColor: 'secondary',
      textColor: 'text-white',
      chartColors: ['#B1D663']
    }
  } else if (progress >= 25) {
    return {
      bgColor: 'error',
      textColor: 'text-white',
      chartColors: ['#EF4444']
    }
  } else {
    return {
      bgColor: 'dark',
      textColor: 'text-white',
      chartColors: ['#6B7280']
    }
  }
})

const bgColor = computed(() => cardColors.value.bgColor)
const textColorClass = computed(() => cardColors.value.textColor)

const hasHint = computed(() => Boolean(props.hintDescription?.trim()))
const hintContent = computed(() => {
  if (!hasHint.value) return null
  return {
    title: props.hintTitle?.trim() || 'Hint',
    description: props.hintDescription.trim()
  }
})

const itemRows = computed(() => {
  if (props.items.length <= 3) {
    return [props.items]
  }

  const rows = []
  for (let i = 0; i < props.items.length; i += 2) {
    rows.push(props.items.slice(i, i + 2))
  }

  return rows
})

// Methods
const getPopoverContent = (item) => {
  if (item.popoverTitle) {
    return {
      title: item.popoverTitle,
      content: item.popoverContent || '',
      variant: 'primary'
    }
  }
  return null
}
</script>

<template>
  <VCard :class="`bg-gradient-${props.color} ${textColorClass}`" elevation="2" rounded="xl">
    <VCardItem class="pb-sm-0">
      <VCardTitle
        :class="`${textColorClass} d-flex align-start justify-space-between pb-5`"
        :style="{ whiteSpace: 'normal', wordBreak: 'break-word' }"
      >
        <span v-html="props.title"></span>
        <VTooltip v-if="hasHint" location="top">
          <template #activator="{ props: tooltipProps }">
            <VIcon
              v-bind="tooltipProps"
              icon="tabler-help"
              size="18"
              :class="['cursor-pointer', 'help-icon', textColorClass, 'pt-7']"
            />
          </template>
          <div v-if="hintContent">
            <strong v-html="hintContent.title" />
            <div v-html="hintContent.description" />
          </div>
        </VTooltip>
      </VCardTitle>
      <VCardSubtitle :class="`${textColorClass} d-flex justify-space-between align-items-center mb-2 w-100`">
        <span class="text-white">{{ props.subtitle }}</span>
        <span class="text-white font-weight-bold">{{ props.progress }}%</span></VCardSubtitle
      >
      <VProgressLinear :model-value="Number(props.progress)" height="6" :color="bgColor" class="mb-2" rounded />
    </VCardItem>

    <VCardItem class="text-center pt-0">
      <template v-if="props.layout === 'rows'">
        <VRow v-for="(item, index) in props.items" :key="index" class="mx-0 mt-0">
          <VCol cols="12" class="d-flex align-items-center flex-column px-1">
            <VCardText :class="`mb-0 px-0 py-1 ${textColorClass}`">
              {{ item.label }}
            </VCardText>
            <VTooltip v-if="getPopoverContent(item)" location="top">
              <template #activator="{ props: tooltipProps }">
                <h4 v-bind="tooltipProps" :class="`font-weight-bold ${textColorClass}`" v-html="item.value" />
              </template>
              <div>
                <strong>{{ getPopoverContent(item).title }}</strong>
                <div v-html="getPopoverContent(item).content" />
              </div>
            </VTooltip>
            <h4 v-else :class="`font-weight-bold ${textColorClass}`" v-html="item.value" />
          </VCol>
        </VRow>
      </template>
      <template v-else>
        <VRow v-for="(row, rowIndex) in itemRows" :key="rowIndex" class="mx-0 mt-0">
          <VCol
            v-for="(item, colIndex) in row"
            :key="`${rowIndex}-${colIndex}`"
            :cols="12 / row.length"
            class="d-flex align-items-center flex-column px-1"
          >
            <VCardText :class="`mb-0 px-0 py-1 ${textColorClass}`">
              {{ item.label }}
            </VCardText>
            <VTooltip v-if="getPopoverContent(item)" location="top">
              <template #activator="{ props: tooltipProps }">
                <h4 v-bind="tooltipProps" :class="`font-weight-bold ${textColorClass}`" v-html="item.value" />
              </template>
              <div>
                <strong>{{ getPopoverContent(item).title }}</strong>
                <div v-html="getPopoverContent(item).content" />
              </div>
            </VTooltip>
            <h4 v-else :class="`font-weight-bold ${textColorClass}`" v-html="item.value" />
          </VCol>
        </VRow>
      </template>
    </VCardItem>
  </VCard>
</template>
