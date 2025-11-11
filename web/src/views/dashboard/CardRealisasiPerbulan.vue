<script setup>
import { computed } from 'vue'
import { rupiahAbbreviate } from '@/@core/utils/formatters'

// Props
const props = defineProps({
  categoryKey: {
    type: String,
    default: ''
  },
  hintTitle: {
    type: String,
    default: 'Information'
  },
  hintDescription: {
    type: String,
    default: ''
  },
  currentMonth: {
    type: Object,
    required: true,
    validator: (value) => value && typeof value.month === 'string' && typeof value.value === 'number'
  },
  monthlyData: {
    type: Array,
    required: true,
    validator: (value) =>
      Array.isArray(value) &&
      value.every((item) => item && typeof item.month === 'string' && typeof item.value === 'number')
  }
})

// Computed properties for colors
const currentMonthColors = computed(() => {
  const progress = props.currentMonth.value
  if (progress >= 75) {
    return {
      bgColor: 'primary',
      textColor: 'text-white'
    }
  } else if (progress >= 50) {
    return {
      bgColor: 'secondary',
      textColor: 'text-white'
    }
  } else if (progress >= 25) {
    return {
      bgColor: 'error',
      textColor: 'text-white'
    }
  } else {
    return {
      bgColor: 'dark',
      textColor: 'text-white'
    }
  }
})

// Methods
const isFutureMonth = (monthName) => {
  const normalize = (name) => {
    if (!name) return -1
    const lower = String(name).toLowerCase()

    if (lower.startsWith('januari')) return 0
    if (lower.startsWith('februari')) return 1
    if (lower.startsWith('maret')) return 2
    if (lower.startsWith('april')) return 3
    if (lower.startsWith('mei')) return 4
    if (lower.startsWith('juni')) return 5
    if (lower.startsWith('juli')) return 6
    if (lower.startsWith('agustus')) return 7
    if (lower.startsWith('september')) return 8
    if (lower.startsWith('oktober')) return 9
    if (lower.startsWith('november')) return 10
    if (lower.startsWith('desember')) return 11

    return -1
  }

  const currentIndex = normalize(props.currentMonth.month)
  const monthIndex = normalize(monthName)

  if (currentIndex === -1 || monthIndex === -1) return false

  return monthIndex > currentIndex
}

// Computed for month data with colors
const monthData = computed(() => {
  return props.monthlyData.map((month) => {
    let colors = { bgColor: 'transparent', textColor: 'text-muted' }
    let showPercentage = true

    if (!isFutureMonth(month.month)) {
      const progress = month.value
      if (progress >= 75) {
        colors = { bgColor: 'primary', textColor: 'text-white' }
      } else if (progress >= 50) {
        colors = { bgColor: 'secondary', textColor: 'text-white' }
      } else if (progress >= 25) {
        colors = { bgColor: 'error', textColor: 'text-white' }
      } else if (progress === 0) {
        colors = { bgColor: 'dark', textColor: 'text-white' }
      } else {
        colors = { bgColor: 'dark', textColor: 'text-white' }
      }
    } else {
      colors.bgColor = 'transparent'
      colors.textColor = 'text-muted'
      showPercentage = false
    }

    return { ...month, colors, showPercentage }
  })
})

// Computed for hint
const hasHint = computed(() => Boolean(props.hintDescription?.trim()))
const hintContent = computed(() => {
  if (!hasHint.value) return null
  return {
    title: props.hintTitle?.trim() || 'Hint',
    description: props.hintDescription.trim()
  }
})

// Check if this is Anggaran category
const isAnggaranCategory = computed(() => {
  return props.categoryKey === 'anggaran'
})
</script>

<template>
  <VRow class="mx-0">
    <!-- Large card for current month -->
    <VCol cols="12" sm="6" lg="4" class="px-0 my-3 py-2 position-relative">
      <VCard
        :class="`bg-gradient-${currentMonthColors.bgColor} ${currentMonthColors.textColor}`"
        elevation="0"
        rounded="xl"
        role="region"
        :aria-label="`Current month ${currentMonth.month} with ${currentMonth.value}% achievement`"
      >
        <VCardItem class="pb-sm-0 w-100" style="position: absolute">
          <VCardTitle :class="`${currentMonthColors.textColor} d-flex align-start justify-space-between pb-5`">
            <span></span>
            <VTooltip v-if="hasHint" location="top">
              <template #activator="{ props: tooltipProps }">
                <VIcon
                  v-bind="tooltipProps"
                  icon="tabler-help"
                  size="18"
                  :class="['cursor-pointer', 'help-icon', currentMonthColors.textColor, 'pt-7']"
                />
              </template>
              <div v-if="hintContent">
                <strong v-html="hintContent.title" />
                <div v-html="hintContent.description" />
              </div>
            </VTooltip>
          </VCardTitle>
        </VCardItem>
        <VCardItem class="text-center py-3 d-flex flex-column justify-center h-100">
          <h5 :class="`mb-1 ${currentMonthColors.textColor} text-h6`">
            {{ currentMonth.month }}
          </h5>
          <h2
            :class="`font-weight-bolder mb-0 ${currentMonthColors.textColor} text-h4`"
            :aria-label="`${currentMonth.value_formatted || currentMonth.value} percent achievement`"
          >
            {{ currentMonth.value_formatted || currentMonth.value }}%
          </h2>
          <p :class="`mb-0 ${currentMonthColors.textColor} text-caption text-right`">
            {{
              isAnggaranCategory
                ? rupiahAbbreviate(currentMonth.realisasi_formatted || currentMonth.realisasi || currentMonth.value)
                : currentMonth.realisasi_formatted || currentMonth.realisasi || currentMonth.value
            }}<small
              >/{{
                isAnggaranCategory
                  ? rupiahAbbreviate(currentMonth.target_formatted || currentMonth.target || 100)
                  : currentMonth.target_formatted || currentMonth.target || 100
              }}</small
            >
          </p>
        </VCardItem>
      </VCard>
    </VCol>

    <!-- Grid of 12 smaller cards -->
    <VCol cols="12" sm="6" lg="8" class="px-0">
      <VRow class="mx-0" dense role="list" aria-label="Monthly performance data">
        <VCol
          v-for="(monthData, index) in monthData"
          :key="index"
          cols="4"
          sm="3"
          md="2"
          lg="3"
          xl="2"
          class="px-0 my-3 py-0 pl-6"
          role="listitem"
        >
          <VCard
            :class="`${monthData.colors.bgColor !== 'transparent' ? 'bg-gradient-' + monthData.colors.bgColor : ''} ${monthData.colors.textColor}`"
            elevation="0"
            rounded="xl"
            height="100"
            min-height="80"
            :aria-label="`${monthData.month}: ${monthData.showPercentage ? monthData.value + '%' : 'No data'}`"
          >
            <VCardItem class="text-center py-2 d-flex flex-column justify-center h-100">
              <p :class="`mb-1 ${monthData.colors.textColor} text-caption font-weight-medium`" style="line-height: 1.2">
                {{ monthData.month }}
              </p>
              <h6
                v-if="monthData.showPercentage"
                :class="`font-weight-bolder mb-0 ${monthData.colors.textColor} text-body-2`"
                :aria-label="`${monthData.value_formatted || monthData.value} percent`"
              >
                {{ monthData.value_formatted || monthData.value }}%

                <p :class="`mb-0 ${monthData.colors.textColor} text-caption text-right`">
                  <small
                    >{{
                      isAnggaranCategory
                        ? rupiahAbbreviate(monthData.realisasi_formatted || monthData.realisasi || monthData.value)
                        : monthData.realisasi_formatted || monthData.realisasi || monthData.value
                    }}<small
                      >/{{
                        isAnggaranCategory
                          ? rupiahAbbreviate(monthData.target_formatted || monthData.target || 100)
                          : monthData.target_formatted || monthData.target || 100
                      }}</small
                    ></small
                  >
                </p>
              </h6>
            </VCardItem>
          </VCard>
        </VCol>
      </VRow>
    </VCol>
  </VRow>
</template>
