<template>
    <b-row class="mx-0 mb-2">
        <!-- Large card for current month -->
        <b-col lg="4" md="12" class="px-1 mb-1">
            <b-card
                no-body
                :class="`bg-gradient-${currentMonthColors.bgColor} ${currentMonthColors.textColor} card-rounded-2rem h-100 row`"
            >
                <b-card-body
                    class="text-center py-2 d-flex flex-column justify-content-center h-100"
                >
                    <h5 :class="`mb-1 ${currentMonthColors.textColor}`">
                        {{ currentMonth.month }}
                    </h5>
                    <h2
                        :class="`font-weight-bolder mb-0 ${currentMonthColors.textColor}`"
                    >
                        {{ currentMonth.value }}%
                    </h2>
                    <p
                        :class="`mb-0 ${currentMonthColors.textColor}`"
                        style="font-size: 0.75rem"
                    >
                        Current Month
                    </p>
                </b-card-body>
            </b-card>
        </b-col>

        <!-- Grid of 12 smaller cards -->
        <b-col lg="8" md="12" class="px-0">
            <b-row class="mx-0">
                <b-col
                    v-for="(monthData, index) in monthData"
                    :key="index"
                    cols="3"
                    class="pl-3 pr-0 mb-1"
                >
                    <b-card
                        no-body
                        :class="`${
                            monthData.colors.bgColor !== 'transparent'
                                ? 'bg-gradient-' + monthData.colors.bgColor
                                : ''
                        } ${
                            monthData.colors.textColor
                        } card-rounded-2rem h-100`"
                    >
                        <b-card-body
                            class="text-center py-1 d-flex flex-column justify-content-center h-100"
                        >
                            <p
                                :class="`mb-0 ${monthData.colors.textColor}`"
                                style="font-size: 0.75rem; line-height: 1.2"
                            >
                                {{ monthData.month }}
                            </p>
                            <h6
                                v-if="monthData.showPercentage"
                                :class="`font-weight-bolder mb-0 ${monthData.colors.textColor}`"
                            >
                                {{ monthData.value }}%
                            </h6>
                        </b-card-body>
                    </b-card>
                </b-col>
            </b-row>
        </b-col>
    </b-row>
</template>

<script>
import {
    BCard,
    BCardHeader,
    BCardBody,
    BRow,
    BCol,
    VBPopover,
} from "bootstrap-vue";
import { getCardColorsByProgress } from "@/utils/cardUtils";

export default {
    components: {
        BCard,
        BCardHeader,
        BCardBody,
        BRow,
        BCol,
    },
    directives: {
        "b-popover": VBPopover,
    },
    props: {
        title: {
            type: String,
            default: "Realisasi Perbulan",
        },
        subtitle: {
            type: String,
            default: "Monthly Realization Overview",
        },
        hintTitle: {
            type: String,
            default: "Information",
        },
        hintDescription: {
            type: String,
            default: "",
        },
        currentMonth: {
            type: Object,
            default: () => ({
                month: "October",
                value: 85,
            }),
        },
        monthlyData: {
            type: Array,
            default: () => [
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
                { month: "Dec", value: 0 },
            ],
        },
    },
    methods: {
        isFutureMonth(monthName) {
            const months = [
                "Jan",
                "Feb",
                "Mar",
                "Apr",
                "May",
                "Jun",
                "Jul",
                "Aug",
                "Sep",
                "Oct",
                "Nov",
                "Dec",
            ];
            const fullMonths = {
                January: "Jan",
                February: "Feb",
                March: "Mar",
                April: "Apr",
                May: "May",
                June: "Jun",
                July: "Jul",
                August: "Aug",
                September: "Sep",
                October: "Oct",
                November: "Nov",
                December: "Dec",
            };
            const currentMonthShort =
                fullMonths[this.currentMonth.month] || this.currentMonth.month;
            const currentIndex = months.indexOf(currentMonthShort);
            const monthIndex = months.indexOf(monthName);
            return monthIndex > currentIndex;
        },
    },
    computed: {
        currentMonthColors() {
            return getCardColorsByProgress(this.currentMonth.value);
        },
        monthData() {
            return this.monthlyData.map((month) => {
                let colors = getCardColorsByProgress(month.value);
                let showPercentage = true;
                if (this.isFutureMonth(month.month)) {
                    colors.bgColor = "transparent";
                    colors.textColor = "text-muted";
                    showPercentage = false;
                } else if (month.value === 0) {
                    colors.bgColor = "dark";
                    colors.textColor = "text-light";
                }
                return { ...month, colors, showPercentage };
            });
        },
    },
};
</script>

<style scoped>
.card-rounded-2rem {
    border-radius: 2rem;
}
</style>
