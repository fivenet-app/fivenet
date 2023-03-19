<script setup lang="ts">
import Footer from '../../components/partials/Footer.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
</script>

<route lang="json">
{
    "name": "dispatches",
    "meta": {
        "requiresAuth": true,
        "permission": "Dispatches.View",
        "breadCrumbs": [
            { "name": "Dispatches", "href": "/dispatches" }
        ]
    }
}
</route>

<template>
    <NavPageHeader title="Dispatches" />
    <ContentWrapper>
        <h1 class="text-lg font-semibold leading-6 text-white">DISPATCH CENTER AND MANAGEMENT SYSTEM</h1>
    </ContentWrapper>
    <Footer />
</template>
