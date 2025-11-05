<template>
    <div>
        <b-row v-if="loading" class="match-height">
            <b-col lg="3" v-for="n in 4" :key="n">
                <b-card class="text-center">
                    <b-skeleton height="200px" />
                    <b-skeleton width="60%" class="mt-1" />
                    <b-skeleton width="40%" class="mt-1" />
                </b-card>
            </b-col>
        </b-row>

        <b-alert v-else-if="error" variant="danger" show>
            <div class="alert-body">
                <feather-icon icon="AlertTriangleIcon" class="mr-50" />
                {{ error }}
            </div>
        </b-alert>

        <b-row v-else class="match-height">
            <b-col lg="3" v-for="(card, index) in realisasiBulan" :key="index">
                <card-realisasi-bulan
                    :title="card.title"
                    :subtitle="card.subtitle"
                    :hint-title="card.hintTitle"
                    :hint-description="card.hintDescription"
                    :items="card.items"
                    :progress="card.progress"
                    :color="card.color"
                    :layout="card.layout || 'columns'"
                />
            </b-col>
        </b-row>
    </div>
</template>

<script>
import { BRow, BCol, BCard, BSkeleton, BAlert } from "bootstrap-vue";
import CardRealisasiBulan from "@/components/CardRealisasiBulan.vue";

export default {
    name: "CardRealisasiBulanSection",
    components: {
        BRow,
        BCol,
        BCard,
        BSkeleton,
        BAlert,
        CardRealisasiBulan,
    },
    props: {
        realisasiBulan: {
            type: Array,
            default: () => [],
        },
        loading: {
            type: Boolean,
            default: false,
        },
        error: {
            type: String,
            default: null,
        },
    },
};
</script>

<style scoped></style>
