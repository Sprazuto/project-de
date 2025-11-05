<template>
    <!-- Rankings List -->
    <div class="rankings-list" role="list" aria-label="Performance rankings">
        <div
            v-for="(item, index) in rankings"
            :key="index"
            class="ranking-item d-flex align-items-center justify-content-between py-1 px-1 mb-1 rounded"
            :class="getRankingClass(index)"
            role="listitem"
            :aria-label="`${item.name} ranked ${index + 1} with ${
                item.total_score
            }% overall score`"
        >
            <div class="d-flex align-items-center">
                <div
                    class="rank-badge mr-2"
                    :class="getRankBadgeClass(index)"
                    aria-hidden="true"
                >
                    <small style="font-size: 0.6rem">#</small>{{ index + 1 }}
                </div>
                <div class="status-score-section mr-3">
                    <b-badge
                        :variant="getScoreBadgeVariant(item.total_score)"
                        class="mb-1"
                        :aria-label="`Status: ${getScoreStatusLabel(
                            item.total_score
                        )}`"
                    >
                        {{ getScoreStatusLabel(item.total_score) }}
                    </b-badge>
                    <div class="score-value" aria-label="Total score">
                        {{ item.total_score }}%
                    </div>
                </div>
                <div class="instance-name-section mr-3">
                    <h6
                        class="mb-0 text-primary font-weight-bold"
                        aria-label="Organization name"
                    >
                        {{ item.name }}
                    </h6>
                </div>
            </div>
            <div
                class="categories d-flex flex-wrap align-items-stretch"
                role="list"
                aria-label="Performance categories"
            >
                <div
                    v-for="(category, catIndex) in item.categories"
                    :key="catIndex"
                    class="category-card position-relative ml-2 d-flex flex-column"
                    role="listitem"
                    :aria-label="`${category.title}: ${category.percentage}% achievement`"
                >
                    <b-card
                        no-body
                        class="mini-card bg-transparent mb-0 overflow-hidden"
                        :class="`border-${getCategoryColor(
                            category.percentage
                        )}`"
                    >
                        <div
                            class="ribbon-flag"
                            v-if="isTopInCategory(category, catIndex)"
                        >
                            <div
                                class="flag-pole w-10 pt-1"
                                :class="`bg-${getCategoryColor(
                                    category.percentage
                                )}`"
                            ></div>
                            <div class="flag-body">
                                <div class="flag-text">TOP</div>
                            </div>
                        </div>
                        <b-card-body class="p-1 position-relative">
                            <div
                                class="watermark-icon position-absolute"
                                style="
                                    top: 50%;
                                    right: -25px;
                                    transform: translateY(-50%) rotate(-12deg);
                                    opacity: 0.05;
                                    z-index: 0;
                                    filter: drop-shadow(
                                        0 4px 18px rgba(0, 0, 0, 0.28)
                                    );
                                    pointer-events: none;
                                "
                            >
                                <feather-icon :icon="category.icon" size="70" />
                            </div>
                            <div
                                class="text-center mb-1 position-relative"
                                style="z-index: 2"
                            >
                                <div
                                    class="position-absolute"
                                    style="top: 0; left: 0; opacity: 0.8"
                                >
                                    <feather-icon
                                        :icon="category.icon"
                                        size="12"
                                    />
                                </div>
                                <small
                                    class="category-subtitle d-block"
                                    style="
                                        font-size: 0.6rem;
                                        margin-bottom: 1px;
                                    "
                                >
                                    {{ category.subtitle }}
                                </small>
                                <small
                                    class="category-title d-block"
                                    style="font-size: 0.7rem; font-weight: bold"
                                >
                                    {{ category.title }}
                                </small>
                            </div>
                            <div
                                class="progress-container mt-1"
                                style="position: relative; z-index: 2"
                            >
                                <b-progress
                                    :value="category.percentage"
                                    max="100"
                                    height="4px"
                                    :variant="
                                        getCategoryColor(category.percentage)
                                    "
                                    class="mb-1"
                                    :aria-label="`Progress: ${category.percentage}%`"
                                />
                                <small
                                    class="percentage text-center d-block"
                                    aria-hidden="true"
                                    >{{ category.percentage }}%</small
                                >
                            </div>
                        </b-card-body>
                    </b-card>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import {
    BCard,
    BCardHeader,
    BCardBody,
    BProgress,
    BBadge,
    VBPopover,
} from "bootstrap-vue";

export default {
    components: {
        BCard,
        BCardHeader,
        BCardBody,
        BProgress,
        BBadge,
    },
    directives: {
        "b-popover": VBPopover,
    },
    props: {
        title: {
            type: String,
            default: "Rankings",
        },
        subtitle: {
            type: String,
            default: "Top performers across categories",
        },
        hintTitle: {
            type: String,
            default: "Information",
        },
        hintDescription: {
            type: String,
            default: "",
        },
        rankings: {
            type: Array,
            default: () => [],
        },
    },
    methods: {
        getRankingClass(index) {
            if (index === 0) return "bg-light-primary";
            if (index === 1) return "bg-light-secondary";
            // if (index === 2) return "bg-light-info";
            return "bg-transparent";
        },
        getRankBadgeClass(index) {
            if (index === 0) return "badge-gold";
            if (index === 1) return "badge-silver";
            if (index === 2) return "badge-bronze";
            return "badge-default";
        },
        getCategoryColor(percentage) {
            if (percentage >= 75) return "primary";
            if (percentage >= 50) return "secondary";
            if (percentage >= 25) return "danger";
            return "dark";
        },
        getScoreBadgeClass(score) {
            if (score >= 75) return "score-melesat";
            if (score >= 50) return "score-berlari";
            if (score >= 25) return "score-berjalan";
            return "score-diam";
        },
        getScoreStatusLabel(score) {
            if (score >= 75) return "Melesat";
            if (score >= 50) return "Berlari";
            if (score >= 25) return "Berjalan";
            return "Diam";
        },
        getScoreBadgeVariant(score) {
            if (score >= 75) return "light-primary";
            if (score >= 50) return "light-secondary";
            if (score >= 25) return "light-danger";
            return "light-dark";
        },
        isTopInCategory(category, catIndex) {
            // Assuming rankings are sorted by total_score descending
            // Find the max percentage for this category across all rankings
            const maxPercentage = Math.max(
                ...this.rankings.map((r) => r.categories[catIndex].percentage)
            );
            return category.percentage === maxPercentage;
        },
    },
};
</script>

<style scoped>
.ranking-item {
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
    font-weight: bold;
    color: white;
}

.status-score-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    min-width: 80px;
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

.instance-name {
    font-size: 0.9rem;
    margin: 2px 0;
}

.score-value {
    font-size: 0.8rem;
    font-weight: bold;
    color: #495057;
}

.categories {
    flex: 1;
    justify-content: flex-end;
}

.category-card {
    min-width: 150px;
    flex: 1;
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
</style>
