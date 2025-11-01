<template>
    <b-card
        no-body
        class="card-rounded-1rem"
        style="border: 1px solid rgba(0, 123, 255, 0.1)"
    >
        <b-card-body class="text-center py-2">
            <h2 class="text-primary mb-1">
                <feather-icon
                    :icon="icon"
                    size="28"
                    class="mr-1 text-primary"
                />
                {{ title }}
            </h2>
            <b-card-text class="text-muted">
                {{ subtitle }}
            </b-card-text>
            <feather-icon
                v-if="hintDescription"
                v-b-popover.click.top.html="{
                    content: hintDescription,
                    title: hintTitle,
                    variant: 'primary',
                }"
                icon="HelpCircleIcon"
                size="20"
                class="cursor-pointer text-muted"
                style="opacity: 0.8"
            />
        </b-card-body>
    </b-card>
</template>

<script>
import { BCard, BCardBody, BCardText, VBPopover } from "bootstrap-vue";

export default {
    components: {
        BCard,
        BCardBody,
        BCardText,
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
        icon: {
            type: String,
            default: "BarChart2Icon",
        },
        iconSize: {
            type: [String, Number],
            default: 28,
        },
        titleClass: {
            type: String,
            default: "text-gradient-primary mb-1",
        },
        iconClass: {
            type: String,
            default: "mr-1 text-primary",
        },
        hintTitle: {
            type: String,
            default: "Information",
        },
        hintDescription: {
            type: String,
            default: "",
        },
    },
};
</script>
