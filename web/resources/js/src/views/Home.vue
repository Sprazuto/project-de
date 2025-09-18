<template>
    <div>
        <b-card title="Kick start your project 🚀">
            <b-card-text>All the best for your new project.</b-card-text>
            <b-card-text
                >Please make sure to read our
                <b-link
                    href="https://pixinvent.com/demo/vuexy-vuejs-admin-dashboard-template/documentation/guide/development/installation.html"
                    target="_blank"
                >
                    Template Documentation
                </b-link>
                to understand where to go from here and how to use our
                template.</b-card-text
            >
        </b-card>

        <b-card title="Want to integrate with JWT? 🔒">
            <b-card-text
                >We carefully crafted JWT flow so you can implement JWT with
                ease and with minimum efforts.</b-card-text
            >
            <b-card-text
                >Please read our JWT Documentation to get more out of JWT
                authentication.</b-card-text
            >
        </b-card>
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
} from "bootstrap-vue";

export default {
    components: {
        BCard,
        BCardText,
        BLink,
        BListGroup,
        BListGroupItem,
    },
    data() {
        return {
            articles: [],
        };
    },
    mounted() {
        this.$http
            .get("/articles")
            .then((response) => {
                this.articles = response.data.results && response.data.results[0] && response.data.results[0].data || response.data.data || response.data.articles || [];
            })
            .catch((error) => {
                console.error("Error fetching articles:", error);
            });
    },
};
</script>

<style></style>
