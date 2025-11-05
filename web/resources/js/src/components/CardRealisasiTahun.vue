<template>
    <b-card
        no-body
        :class="`bg-gradient-${color} ${textColorClass} card-rounded-2rem`"
    >
        <b-card-header class="w-100 pb-0">
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
                    ></p>
                </div>
                <feather-icon
                    v-if="hintDescription"
                    v-b-popover.click.top.html="{
                        content: hintDescription,
                        title: hintTitle,
                        variant: color,
                    }"
                    icon="HelpCircleIcon"
                    size="18"
                    :class="`cursor-pointer ${textColorClass}`"
                    style="opacity: 0.8"
                />
            </div>
        </b-card-header>
        <div
            class="d-flex justify-content-between align-items-center mb-50 w-100 px-2"
        >
            <span>{{ subtitle }}</span>
            <span>{{ progress }}%</span>
        </div>
        <b-progress
            :value="progress"
            max="100"
            height="6px"
            :variant="bgColor"
            class="mx-2"
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
                        <h4
                            :class="`font-weight-bolder mb-50 ${textColorClass}`"
                            v-html="item.value"
                            v-b-popover.hover.top.html="getPopoverContent(item)"
                        ></h4>
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
                        <h4
                            :class="`font-weight-bolder mb-50 ${textColorClass}`"
                            v-html="item.value"
                            v-b-popover.hover.top.html="getPopoverContent(item)"
                        ></h4>
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
    BProgress,
    VBPopover,
} from "bootstrap-vue";
import { getCardColorsByProgress } from "@/utils/cardUtils";

export default {
    components: {
        BCard,
        BCardHeader,
        BRow,
        BCardText,
        BCol,
        BProgress,
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
        color: {
            type: String,
            default: "dark",
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
