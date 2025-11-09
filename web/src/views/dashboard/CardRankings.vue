<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: {
    type: String,
    default: 'Rankings'
  },
  subtitle: {
    type: String,
    default: 'Top performers across categories'
  },
  rankings: {
    type: Array,
    default: () => []
  }
})

// Icon mapping from Feather to Tabler icons
const iconMapping = {
  LayersIcon: 'tabler-stack-pop',
  MapPinIcon: 'tabler-map-pin',
  ShoppingBagIcon: 'tabler-shopping-bag',
  SlidersIcon: 'tabler-adjustments-alt',
  AwardIcon: 'tabler-award'
}

// Methods
const getRankingClass = (index) => {
  if (index === 0) return 'bg-light-primary'
  if (index === 1) return 'bg-light-secondary'
  return 'bg-transparent'
}

const getRankBadgeClass = (index) => {
  if (index === 0) return 'badge-gold'
  if (index === 1) return 'badge-silver'
  if (index === 2) return 'badge-bronze'
  return 'badge-default'
}

const getCategoryColor = (percentage) => {
  if (percentage >= 75) return 'primary'
  if (percentage >= 50) return 'secondary'
  if (percentage >= 25) return 'danger'
  return 'dark'
}

const getScoreStatusLabel = (score) => {
  if (score >= 75) return 'Melesat'
  if (score >= 50) return 'Berlari'
  if (score >= 25) return 'Berjalan'
  return 'Diam'
}

const getScoreBadgeVariant = (score) => {
  if (score >= 75) return 'primary'
  if (score >= 50) return 'secondary'
  if (score >= 25) return 'danger'
  return 'dark'
}

const getTablerIcon = (featherIcon) => {
  return iconMapping[featherIcon] || 'tabler-star'
}

const getRibbonType = (itemIndex, catIndex) => {
  const percentages = props.rankings.map((r, i) => ({ percentage: r.categories[catIndex].percentage, index: i }))
  percentages.sort((a, b) => b.percentage - a.percentage)
  const position = percentages.findIndex((p) => p.index === itemIndex)
  if (position === percentages.length - 1 && percentages.length > 1) return 'dark'
  if (position === 0) return 'primary'
  if (position === 1) return 'secondary'
  return null
}

const getRibbonIcon = (type) => {
  if (type === 'primary') return 'tabler-number-1'
  if (type === 'secondary') return 'tabler-number-2'
  if (type === 'dark') return 'tabler-ribbon-health'
  return null
}
</script>

<template>
  <!-- Rankings List -->
  <div role="list" aria-label="Performance rankings">
    <div
      v-for="(item, index) in props.rankings"
      :key="index"
      class="ranking-item d-flex align-items-center justify-content-between py-1 px-1 mb-1"
      :class="getRankingClass(index)"
      role="listitem"
      :aria-label="`${item.name} ranked ${index + 1} with ${item.total_score}% overall score`"
    >
      <div class="d-flex align-items-center">
        <div class="rank-badge mr-2" :class="getRankBadgeClass(index)" aria-hidden="true"><small style="font-size: 0.6rem">#</small>{{ index + 1 }}</div>
        <div class="status-score-section mr-3">
          <VChip :color="getScoreBadgeVariant(item.total_score)" size="x-small" variant="flat" :aria-label="`Status: ${getScoreStatusLabel(item.total_score)}`">
            {{ getScoreStatusLabel(item.total_score) }}
          </VChip>
          <div aria-label="Total score">{{ item.total_score }}%</div>
        </div>
        <div class="instance-name-section mr-3">
          <h4 class="mb-0 text-primary font-weight-bold" aria-label="Organization name">
            {{ item.name }}
          </h4>
        </div>
      </div>
      <div class="categories d-flex align-items-stretch" role="list" aria-label="Performance categories">
        <div
          v-for="(category, catIndex) in item.categories"
          :key="catIndex"
          class="category-card position-relative ml-2 d-flex flex-column"
          role="listitem"
          :aria-label="`${category.title}: ${category.percentage}% achievement`"
        >
          <div
            v-if="getRibbonType(index, catIndex)"
            class="ribbon position-absolute"
            :class="`bg-${getRibbonType(index, catIndex)} text-white`"
            style="
              top: -4px;
              right: 14px;
              z-index: 10;
              width: 20px;
              height: 25px;
              clip-path: polygon(0 0, 100% 0, 100% 70%, 50% 100%, 0 70%);
              display: flex;
              align-items: center;
              justify-content: center;
              font-size: 0.6rem;
              font-weight: bold;
            "
          >
            <VIcon :icon="getRibbonIcon(getRibbonType(index, catIndex))" size="12" />
          </div>
          <VCard flat class="mini-card bg-transparent mb-0 overflow-hidden">
            <VCardText class="py-4 px-4 position-relative">
              <div
                class="watermark-icon position-absolute"
                style="top: 50%; right: -25px; transform: translateY(-50%) rotate(-12deg); opacity: 0.05; z-index: 0; filter: drop-shadow(0 4px 18px rgba(0, 0, 0, 0.28)); pointer-events: none"
              >
                <VIcon :icon="getTablerIcon(category.icon)" size="70" />
              </div>
              <div class="text-center mb-1 position-relative" style="z-index: 2">
                <div class="position-absolute" style="top: 0; left: 0; opacity: 0.8">
                  <VIcon :icon="getTablerIcon(category.icon)" size="18" />
                </div>
                <small class="category-subtitle d-block" style="font-size: 0.6rem; margin-bottom: 1px">
                  {{ category.subtitle }}
                </small>
                <small class="category-title d-block" style="font-size: 0.7rem; font-weight: bold">
                  {{ category.title }}
                </small>
              </div>
              <div class="progress-container mt-1" style="position: relative; z-index: 2">
                <VProgressLinear :model-value="category.percentage" height="4" :color="getCategoryColor(category.percentage)" class="mb-1" :aria-label="`Progress: ${category.percentage}%`" />
                <small class="percentage text-center d-block" aria-hidden="true">{{ category.percentage }}%</small>
              </div>
            </VCardText>
          </VCard>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.ranking-item {
  border-radius: 1.3rem;
  transition: all 0.3s ease;
  /* border: 1px solid #e9ecef; */
}

.ranking-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.rank-badge {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  align-self: center;
  font-weight: bold;
  flex: 0 0 30px;
}

.status-score-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  width: 80px;
  flex: 0 0 80px;
}

.instance-name-section {
  display: flex;
  align-items: center;
  min-width: 120px;
}

.status-label {
  font-size: 0.7rem;
  font-weight: bold;
  color: #6c757d;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.categories {
  flex: 1;
  justify-content: flex-end;
}

.category-card {
  min-width: 15vw;
  position: relative;
  margin: 0 2px 2px 0;
  display: flex;
  flex-direction: column;
}

.category-title {
  font-size: 0.75rem;
  font-weight: bold;
  /* color: white; */
  /* text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5); */
  line-height: 1.2;
  margin-bottom: 2px;
}

.progress-container {
  text-align: center;
}

.percentage {
  font-size: 0.8rem;
  font-weight: bold;
  /* color: white; */
  /* text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5); */
}

.mini-card {
  border-radius: 1rem;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .ranking-item {
    flex-direction: column;
    align-items: flex-start;
  }

  .categories {
    width: 100%;
    margin-top: 8px;
  }

  .category-card {
    min-width: 120px;
  }

  .status-score-section {
    min-width: 60px;
  }
}
</style>
