<script setup>
import { computed } from 'vue'

// Props
const props = defineProps({
  title: {
    type: String,
    default: 'Realisasi Perbulan'
  },
  subtitle: {
    type: String,
    default: 'Monthly Realization Overview'
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
    validator: (value) => Array.isArray(value) && value.every((item) => item && typeof item.month === 'string' && typeof item.value === 'number')
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
  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
  const fullMonths = {
    January: 'Jan',
    February: 'Feb',
    March: 'Mar',
    April: 'Apr',
    May: 'May',
    June: 'Jun',
    July: 'Jul',
    August: 'Aug',
    September: 'Sep',
    October: 'Oct',
    November: 'Nov',
    December: 'Dec'
  }
  const currentMonthShort = fullMonths[props.currentMonth.month] || props.currentMonth.month
  const currentIndex = months.indexOf(currentMonthShort)
  const monthIndex = months.indexOf(monthName)
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
                <VIcon v-bind="tooltipProps" icon="tabler-help" size="18" :class="['cursor-pointer', 'help-icon', currentMonthColors.textColor, 'pt-7']" />
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
          <h2 :class="`font-weight-bolder mb-0 ${currentMonthColors.textColor} text-h4`" :aria-label="`${currentMonth.value} percent achievement`">{{ currentMonth.value }}%</h2>
          <p :class="`mb-0 ${currentMonthColors.textColor} text-caption`">Current Month</p>
        </VCardItem>
      </VCard>
    </VCol>

    <!-- Grid of 12 smaller cards -->
    <VCol cols="12" sm="6" lg="8" class="px-0">
      <VRow class="mx-0" dense role="list" aria-label="Monthly performance data">
        <VCol v-for="(monthData, index) in monthData" :key="index" cols="4" sm="3" md="2" lg="3" xl="2" class="px-0 my-3 py-0 pl-6" role="listitem">
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
              <h6 v-if="monthData.showPercentage" :class="`font-weight-bolder mb-0 ${monthData.colors.textColor} text-body-2`" :aria-label="`${monthData.value} percent`">{{ monthData.value }}%</h6>
              <p v-else :class="`mb-0 ${monthData.colors.textColor} text-caption opacity-6`">-</p>
            </VCardItem>
          </VCard>
        </VCol>
      </VRow>
    </VCol>
  </VRow>
</template>
