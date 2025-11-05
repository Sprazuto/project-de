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
                    :current-month="{ month: 'October', value: 85 }"
                    :monthly-data="[
                        { month: 'Jan', value: 75 },
                        { month: 'Feb', value: 80 },
                        { month: 'Mar', value: 70 },
                        { month: 'Apr', value: 85 },
                        { month: 'May', value: 90 },
                        { month: 'Jun', value: 78 },
                        { month: 'Jul', value: 82 },
                        { month: 'Aug', value: 88 },
                        { month: 'Sep', value: 92 },
                        { month: 'Oct', value: 85 },
                        { month: 'Nov', value: 0 },
                        { month: 'Dec', value: 95 },
                    ]"
                />
            </b-col>
        </b-row>

        <card-header title="Realisasi Perbulan Fisik" :icon="'MapPinIcon'" />

        <b-row class="match-height">
            <b-col lg="12">
                <card-realisasi-perbulan
                    :current-month="{ month: 'October', value: 20 }"
                    :monthly-data="[
                        { month: 'Jan', value: 15 },
                        { month: 'Feb', value: 18 },
                        { month: 'Mar', value: 12 },
                        { month: 'Apr', value: 22 },
                        { month: 'May', value: 25 },
                        { month: 'Jun', value: 19 },
                        { month: 'Jul', value: 21 },
                        { month: 'Aug', value: 23 },
                        { month: 'Sep', value: 26 },
                        { month: 'Oct', value: 20 },
                        { month: 'Nov', value: 0 },
                        { month: 'Dec', value: 0 },
                    ]"
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
                    :current-month="{ month: 'October', value: 60 }"
                    :monthly-data="[
                        { month: 'Jan', value: 55 },
                        { month: 'Feb', value: 58 },
                        { month: 'Mar', value: 52 },
                        { month: 'Apr', value: 62 },
                        { month: 'May', value: 65 },
                        { month: 'Jun', value: 0 },
                        { month: 'Jul', value: 61 },
                        { month: 'Aug', value: 63 },
                        { month: 'Sep', value: 66 },
                        { month: 'Oct', value: 60 },
                        { month: 'Nov', value: 0 },
                        { month: 'Dec', value: 0 },
                    ]"
                />
            </b-col>
        </b-row>

        <card-header title="Realisasi Perbulan Kinerja" :icon="'SlidersIcon'" />

        <b-row class="match-height">
            <b-col lg="12">
                <card-realisasi-perbulan
                    :current-month="{ month: 'October', value: 35 }"
                    :monthly-data="[
                        { month: 'Jan', value: 30 },
                        { month: 'Feb', value: 33 },
                        { month: 'Mar', value: 27 },
                        { month: 'Apr', value: 37 },
                        { month: 'May', value: 40 },
                        { month: 'Jun', value: 34 },
                        { month: 'Jul', value: 36 },
                        { month: 'Aug', value: 38 },
                        { month: 'Sep', value: 41 },
                        { month: 'Oct', value: 35 },
                        { month: 'Nov', value: 0 },
                        { month: 'Dec', value: 0 },
                    ]"
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

        <b-card v-if="articles.length > 0" title="Articles from API">
            <b-list-group>
                <b-list-group-item
                    v-for="article in articles"
                    :key="article.id"
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
import CardRealisasiBulanSection from "@/components/CardRealisasiBulanSection.vue";
import CardRealisasiTahunSection from "@/components/CardRealisasiTahunSection.vue";
import CardRealisasiPerbulan from "@/components/CardRealisasiPerbulan.vue";
import CardHeader from "@/components/CardHeader.vue";
import CardRankings from "@/components/CardRankings.vue";
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
        CardRealisasiBulanSection,
        CardRealisasiTahunSection,
        CardRealisasiPerbulan,
        CardHeader,
        CardRankings,
    },
    data() {
        return {
            articles: [],
            realisasiBulan: [],
            realisasiTahun: [],
            loading: {
                bulan: false,
                tahun: false,
            },
            error: {
                bulan: null,
                tahun: null,
            },
        };
    },
    computed: {
        rankings() {
            return [
                {
                    name: "Kabupaten A",
                    total_score: 85.75,
                    categories: [
                        {
                            title: "Realisasi Barjas",
                            subtitle: "Capaian Tahunan",
                            percentage: 85,
                            icon: "LayersIcon",
                        },
                        {
                            title: "Realisasi Fisik",
                            subtitle: "Capaian Tahunan",
                            percentage: 78,
                            icon: "MapPinIcon",
                        },
                        {
                            title: "Realisasi Anggaran",
                            subtitle: "Capaian Tahunan",
                            percentage: 92,
                            icon: "ShoppingBagIcon",
                        },
                        {
                            title: "Realisasi Kinerja",
                            subtitle: "Capaian Tahunan",
                            percentage: 88,
                            icon: "SlidersIcon",
                        },
                    ],
                },
                {
                    name: "Kabupaten B",
                    total_score: 82.75,
                    categories: [
                        {
                            title: "Realisasi Barjas",
                            subtitle: "Capaian Tahunan",
                            percentage: 82,
                            icon: "LayersIcon",
                        },
                        {
                            title: "Realisasi Fisik",
                            subtitle: "Capaian Tahunan",
                            percentage: 75,
                            icon: "MapPinIcon",
                        },
                        {
                            title: "Realisasi Anggaran",
                            subtitle: "Capaian Tahunan",
                            percentage: 89,
                            icon: "ShoppingBagIcon",
                        },
                        {
                            title: "Realisasi Kinerja",
                            subtitle: "Capaian Tahunan",
                            percentage: 85,
                            icon: "SlidersIcon",
                        },
                    ],
                },
                {
                    name: "Kabupaten C",
                    total_score: 79.75,
                    categories: [
                        {
                            title: "Realisasi Barjas",
                            subtitle: "Capaian Tahunan",
                            percentage: 79,
                            icon: "LayersIcon",
                        },
                        {
                            title: "Realisasi Fisik",
                            subtitle: "Capaian Tahunan",
                            percentage: 72,
                            icon: "MapPinIcon",
                        },
                        {
                            title: "Realisasi Anggaran",
                            subtitle: "Capaian Tahunan",
                            percentage: 86,
                            icon: "ShoppingBagIcon",
                        },
                        {
                            title: "Realisasi Kinerja",
                            subtitle: "Capaian Tahunan",
                            percentage: 82,
                            icon: "SlidersIcon",
                        },
                    ],
                },
                {
                    name: "Kabupaten D",
                    total_score: 76.75,
                    categories: [
                        {
                            title: "Realisasi Barjas",
                            subtitle: "Capaian Tahunan",
                            percentage: 76,
                            icon: "LayersIcon",
                        },
                        {
                            title: "Realisasi Fisik",
                            subtitle: "Capaian Tahunan",
                            percentage: 69,
                            icon: "MapPinIcon",
                        },
                        {
                            title: "Realisasi Anggaran",
                            subtitle: "Capaian Tahunan",
                            percentage: 83,
                            icon: "ShoppingBagIcon",
                        },
                        {
                            title: "Realisasi Kinerja",
                            subtitle: "Capaian Tahunan",
                            percentage: 79,
                            icon: "SlidersIcon",
                        },
                    ],
                },
                {
                    name: "Kabupaten E",
                    total_score: 73.75,
                    categories: [
                        {
                            title: "Realisasi Barjas",
                            subtitle: "Capaian Tahunan",
                            percentage: 73,
                            icon: "LayersIcon",
                        },
                        {
                            title: "Realisasi Fisik",
                            subtitle: "Capaian Tahunan",
                            percentage: 66,
                            icon: "MapPinIcon",
                        },
                        {
                            title: "Realisasi Anggaran",
                            subtitle: "Capaian Tahunan",
                            percentage: 80,
                            icon: "ShoppingBagIcon",
                        },
                        {
                            title: "Realisasi Kinerja",
                            subtitle: "Capaian Tahunan",
                            percentage: 76,
                            icon: "SlidersIcon",
                        },
                    ],
                },
            ];
        },
    },
    mounted() {
        this.fetchArticles();
        this.fetchRealisasiBulan();
        this.fetchRealisasiTahun();
    },
    methods: {
        fetchArticles() {
            this.$http
                .get("/articles")
                .then((response) => {
                    this.articles =
                        (response.data.results &&
                            response.data.results[0] &&
                            response.data.results[0].data) ||
                        response.data.data ||
                        response.data.articles ||
                        [];
                })
                .catch((error) => {
                    console.error("Error fetching articles:", error);
                });
        },
        fetchRealisasiBulan(params = {}) {
            this.loading.bulan = true;
            this.error.bulan = null;
            this.$http
                .get("/realisasi-bulan", {
                    params: {
                        tahun: params.tahun,
                        bulan: params.bulan,
                        idsatker: params.idsatker || 0,
                    },
                })
                .then((response) => {
                    const processed = processRealisasiBulanData(response.data);
                    this.realisasiBulan.splice(
                        0,
                        this.realisasiBulan.length,
                        ...processed
                    );
                })
                .catch((error) => {
                    this.error.bulan =
                        error.message || "Failed to fetch realisasi bulan data";
                    console.error("Error fetching realisasi bulan:", error);
                })
                .finally(() => {
                    this.loading.bulan = false;
                });
        },
        fetchRealisasiTahun(params = {}) {
            this.loading.tahun = true;
            this.error.tahun = null;
            this.$http
                .get("/realisasi-tahun", {
                    params: {
                        tahun: params.tahun,
                        idsatker: params.idsatker || 0,
                    },
                })
                .then((response) => {
                    const processed = processRealisasiTahunData(response.data);
                    this.realisasiTahun.splice(
                        0,
                        this.realisasiTahun.length,
                        ...processed
                    );
                })
                .catch((error) => {
                    this.error.tahun =
                        error.message || "Failed to fetch realisasi tahun data";
                    console.error("Error fetching realisasi tahun:", error);
                })
                .finally(() => {
                    this.loading.tahun = false;
                });
        },
    },
};
</script>

<style></style>
