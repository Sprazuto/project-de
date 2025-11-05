<template>
    <b-card
        no-body
        :class="`bg-gradient-${bgColor} ${textColorClass} card-rounded-2rem`"
    >
        <b-card-header class="w-100">
            <div class="d-flex justify-content-between align-items-start w-100">
                <div class="flex-grow-1">
                    <h4 :class="`mb-1 ${textColorClass}`">{{ title }}</h4>
                    <p
                        :class="`mb-0 ${textColorClass}`"
                        style="
                            font-size: 0.875rem;
                            opacity: 0.8;
                            line-height: 1.3;
                        "
                    >
                        {{ subtitle }}
                    </p>
                </div>
                <feather-icon
                    v-if="hintDescription"
                    v-b-popover.click.top.html="{
                        content: hintDescription,
                        title: hintTitle,
                        variant: bgColor,
                    }"
                    icon="HelpCircleIcon"
                    size="18"
                    :class="`cursor-pointer ${textColorClass}`"
                    style="opacity: 0.8"
                />
            </div>
        </b-card-header>

        <!-- apex chart -->
        <vue-apex-charts
            type="radialBar"
            height="245"
            :options="chartOptions"
            :series="goalOverviewRadialBar.series"
        />
        <div class="text-center py-1">
            <template v-if="layout === 'rows'">
                <b-row
                    v-for="(item, index) in items"
                    :key="index"
                    class="mx-0 my-50"
                >
                    <b-col
                        cols="12"
                        class="d-flex align-items-between flex-column"
                    >
                        <b-card-text :class="`mb-0 ${textColorClass}`">
                            {{ item.label }}
                        </b-card-text>
                        <h3
                            :class="`font-weight-bolder mb-50 ${textColorClass}`"
                            v-html="item.value"
                            v-b-popover.hover.top.html="getPopoverContent(item)"
                        ></h3>
                    </b-col>
                </b-row>
            </template>
            <template v-else>
                <b-row
                    v-for="(row, rowIndex) in itemRows"
                    :key="rowIndex"
                    class="mx-0 my-50"
                >
                    <b-col
                        v-for="(item, colIndex) in row"
                        :key="`${rowIndex}-${colIndex}`"
                        :cols="12 / row.length"
                        class="d-flex align-items-between flex-column"
                    >
                        <b-card-text :class="`mb-0 ${textColorClass}`">
                            {{ item.label }}
                        </b-card-text>
                        <h3
                            :class="`font-weight-bolder mb-50 ${textColorClass}`"
                            v-html="item.value"
                            v-b-popover.hover.top.html="getPopoverContent(item)"
                        ></h3>
                    </b-col>
                </b-row>
            </template>
        </div>
    </b-card>
</template>

<script>
import {
    BCard,
    BCardHeader,
    BRow,
    BCol,
    BCardText,
    VBPopover,
} from "bootstrap-vue";
import VueApexCharts from "vue-apexcharts";
import { $themeColors } from "@themeConfig";
import { getCardColorsByProgress } from "@/utils/cardUtils";

export default {
    components: {
        VueApexCharts,
        BCard,
        BCardHeader,
        BRow,
        BCardText,
        BCol,
    },
    directives: {
        "b-popover": VBPopover,
    },
    props: {
        title: {
            type: String,
            default: "Realisasi",
        },
        subtitle: {
            type: String,
            default: "",
        },
        hintTitle: {
            type: String,
            default: "",
        },
        hintDescription: {
            type: String,
            default: "",
        },
        items: {
            type: Array,
            default: () => [],
        },
        layout: {
            type: String,
            default: "columns", // 'columns' or 'rows'
        },
        progress: {
            type: [Number, String],
            default: 0,
        },
    },
    computed: {
        cardColors() {
            return getCardColorsByProgress(this.progress);
        },
        bgColor() {
            return this.cardColors.bgColor;
        },
        textColorClass() {
            return this.cardColors.textColor;
        },
        itemRows() {
            if (this.items.length <= 3) {
                return [this.items];
            } else {
                // For more than 3 items, split into rows of 2
                const rows = [];
                for (let i = 0; i < this.items.length; i += 2) {
                    rows.push(this.items.slice(i, i + 2));
                }
                return rows;
            }
        },
        chartOptions() {
            return {
                chart: {
                    sparkline: { enabled: true },
                    dropShadow: {
                        enabled: true,
                        blur: 3,
                        left: 1,
                        top: 1,
                        opacity: 0.1,
                    },
                },
                colors: ["#ebe9f1"],
                plotOptions: {
                    radialBar: {
                        offsetY: -10,
                        startAngle: -150,
                        endAngle: 150,
                        hollow: { size: "77%" },
                        track: {
                            background: "#ebe9f111",
                            strokeWidth: "50%",
                        },
                        dataLabels: {
                            name: { show: false },
                            value: {
                                color:
                                    this.textColorClass === "text-light"
                                        ? "#ebe9f1"
                                        : "#5e5873",
                                fontSize: "2.86rem",
                                fontWeight: "600",
                            },
                        },
                    },
                },
                fill: {
                    type: "gradient",
                    gradient: {
                        shade: "dark",
                        type: "horizontal",
                        shadeIntensity: 0.5,
                        gradientToColors: [$themeColors[this.bgColor]],
                        inverseColors: true,
                        opacityFrom: 1,
                        opacityTo: 1,
                        stops: [0, 100],
                    },
                },
                stroke: { lineCap: "round" },
                grid: { padding: { bottom: 30 } },
            };
        },
    },
    data() {
        return {
            goalOverviewRadialBar: {
                series: [this.progress],
            },
        };
    },
    methods: {
        getPopoverContent(item) {
            if (item.popoverTitle) {
                return {
                    title: item.popoverTitle,
                    content: item.popoverContent || "",
                    variant: "primary",
                };
            }
            return null;
        },
    },
};
</script>
