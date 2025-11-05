<template>
    <div>
        <card-realisasi-bulan-section
            :realisasi-bulan="realisasiBulan"
            :loading="loading.bulan"
            :error="error.bulan"
        />

        <card-realisasi-tahun-section
            :realisasi-tahun="realisasiTahun"
            :loading="loading.tahun"
            :error="error.tahun"
        />

        <card-header title="Realisasi Perbulan Barjas" :icon="'LayersIcon'" />

        <b-row class="match-height">
            <b-col lg="12">
                <card-realisasi-perbulan
                    :current-month="monthlyData.barjas.currentMonth"
                    :monthly-data="monthlyData.barjas.data"
                />
            </b-col>
        </b-row>

        <card-header title="Realisasi Perbulan Fisik" :icon="'MapPinIcon'" />

        <b-row class="match-height">
            <b-col lg="12">
                <card-realisasi-perbulan
                    :current-month="monthlyData.fisik.currentMonth"
                    :monthly-data="monthlyData.fisik.data"
                />
            </b-col>
        </b-row>

        <card-header
            title="Realisasi Perbulan Anggaran"
            :icon="'ShoppingBagIcon'"
        />

        <b-row class="match-height">
            <b-col lg="12">
                <card-realisasi-perbulan
                    :current-month="monthlyData.anggaran.currentMonth"
                    :monthly-data="monthlyData.anggaran.data"
                />
            </b-col>
        </b-row>

        <card-header title="Realisasi Perbulan Kinerja" :icon="'SlidersIcon'" />

        <b-row class="match-height">
            <b-col lg="12">
                <card-realisasi-perbulan
                    :current-month="monthlyData.kinerja.currentMonth"
                    :monthly-data="monthlyData.kinerja.data"
                />
            </b-col>
        </b-row>

        <card-header
            title="Peringkat Kinerja"
            subtitle="Organisasi Perangkat Daerah"
            :icon="'AwardIcon'"
        />

        <b-row class="match-height">
            <b-col lg="12">
                <card-rankings :rankings="rankings" />
            </b-col>
        </b-row>

        <b-card
            v-if="articles.length > 0"
            title="Articles from API"
            role="region"
            aria-label="Articles section"
        >
            <b-list-group role="list" aria-label="List of articles">
                <b-list-group-item
                    v-for="article in articles"
                    :key="article.id"
                    role="listitem"
                    :aria-label="`Article: ${article.title}`"
                >
                    <h6>{{ article.title }}</h6>
                    <p>{{ article.content }}</p>
                </b-list-group-item>
            </b-list-group>
        </b-card>
    </div>
</template>

<script>
import {
    BCard,
    BCardText,
    BLink,
    BListGroup,
    BListGroupItem,
    BRow,
    BCol,
} from "bootstrap-vue";
import {
    processRealisasiBulanData,
    processRealisasiTahunData,
} from "@/utils/realisasiDataProcessor";

export default {
    components: {
        BCard,
        BCardText,
        BLink,
        BListGroup,
        BListGroupItem,
        BRow,
        BCol,
        CardRealisasiBulanSection: () =>
            import("@/components/CardRealisasiBulanSection.vue"),
        CardRealisasiTahunSection: () =>
            import("@/components/CardRealisasiTahunSection.vue"),
        CardRealisasiPerbulan: () =>
            import("@/components/CardRealisasiPerbulan.vue"),
        CardHeader: () => import("@/components/CardHeader.vue"),
        CardRankings: () => import("@/components/CardRankings.vue"),
    },
    data() {
        return {
            articles: [],
            realisasiBulan: [],
            realisasiTahun: [],
            loading: {
                bulan: true,
                tahun: true,
            },
            error: {
                bulan: null,
                tahun: null,
            },
            monthlyData: {
                barjas: {
                    currentMonth: { month: "October", value: 85 },
                    data: [
                        { month: "Jan", value: 75 },
                        { month: "Feb", value: 80 },
                        { month: "Mar", value: 70 },
                        { month: "Apr", value: 85 },
                        { month: "May", value: 90 },
                        { month: "Jun", value: 78 },
                        { month: "Jul", value: 82 },
                        { month: "Aug", value: 88 },
                        { month: "Sep", value: 92 },
                        { month: "Oct", value: 85 },
                        { month: "Nov", value: 0 },
                        { month: "Dec", value: 95 },
                    ],
                },
                fisik: {
                    currentMonth: { month: "October", value: 20 },
                    data: [
                        { month: "Jan", value: 15 },
                        { month: "Feb", value: 18 },
                        { month: "Mar", value: 12 },
                        { month: "Apr", value: 22 },
                        { month: "May", value: 25 },
                        { month: "Jun", value: 19 },
                        { month: "Jul", value: 21 },
                        { month: "Aug", value: 23 },
                        { month: "Sep", value: 26 },
                        { month: "Oct", value: 20 },
                        { month: "Nov", value: 0 },
                        { month: "Dec", value: 0 },
                    ],
                },
                anggaran: {
                    currentMonth: { month: "October", value: 60 },
                    data: [
                        { month: "Jan", value: 55 },
                        { month: "Feb", value: 58 },
                        { month: "Mar", value: 52 },
                        { month: "Apr", value: 62 },
                        { month: "May", value: 65 },
                        { month: "Jun", value: 0 },
                        { month: "Jul", value: 61 },
                        { month: "Aug", value: 63 },
                        { month: "Sep", value: 66 },
                        { month: "Oct", value: 60 },
                        { month: "Nov", value: 0 },
                        { month: "Dec", value: 0 },
                    ],
                },
                kinerja: {
                    currentMonth: { month: "October", value: 35 },
                    data: [
                        { month: "Jan", value: 30 },
                        { month: "Feb", value: 33 },
                        { month: "Mar", value: 27 },
                        { month: "Apr", value: 37 },
                        { month: "May", value: 40 },
                        { month: "Jun", value: 34 },
                        { month: "Jul", value: 36 },
                        { month: "Aug", value: 38 },
                        { month: "Sep", value: 41 },
                        { month: "Oct", value: 35 },
                        { month: "Nov", value: 0 },
                        { month: "Dec", value: 0 },
                    ],
                },
            },
        };
    },
    computed: {
        rankings() {
            return this.$store.getters["rankings/getRankings"] || [];
        },
    },
    async mounted() {
        // let Vue render first (skeleton visible)
        await this.$nextTick();

        // now start loading data
        this.initializeData();
    },
    methods: {
        async initializeData() {
            try {
                await Promise.allSettled([
                    this.fetchArticles(),
                    this.fetchRealisasiBulan(),
                    this.fetchRealisasiTahun(),
                    this.$store.dispatch("rankings/fetchRankings"),
                ]);
            } catch (error) {
                console.error("Error initializing data:", error);
            }
        },
        async fetchArticles() {
            try {
                const response = await this.$http.get("/articles");
                this.articles =
                    response.data.results?.[0]?.data ||
                    response.data.data ||
                    response.data.articles ||
                    [];
            } catch (error) {
                console.error("Error fetching articles:", error);
                this.articles = [];
            }
        },
        async fetchRealisasiBulan(params = {}) {
            this.loading.bulan = true;
            this.error.bulan = null;
            try {
                const response = await this.$http.get("/realisasi-bulan", {
                    params: {
                        tahun: params.tahun,
                        bulan: params.bulan,
                        idsatker: params.idsatker || 0,
                    },
                });
                const processed = processRealisasiBulanData(response.data);
                this.realisasiBulan.splice(
                    0,
                    this.realisasiBulan.length,
                    ...processed
                );
            } catch (error) {
                this.error.bulan =
                    error.message || "Failed to fetch realisasi bulan data";
                console.error("Error fetching realisasi bulan:", error);
            } finally {
                this.loading.bulan = false;
            }
        },
        async fetchRealisasiTahun(params = {}) {
            this.loading.tahun = true;
            this.error.tahun = null;
            try {
                const response = await this.$http.get("/realisasi-tahun", {
                    params: {
                        tahun: params.tahun,
                        idsatker: params.idsatker || 0,
                    },
                });
                const processed = processRealisasiTahunData(response.data);
                this.realisasiTahun.splice(
                    0,
                    this.realisasiTahun.length,
                    ...processed
                );
            } catch (error) {
                this.error.tahun =
                    error.message || "Failed to fetch realisasi tahun data";
                console.error("Error fetching realisasi tahun:", error);
            } finally {
                this.loading.tahun = false;
            }
        },
    },
};
</script>

<style></style>
