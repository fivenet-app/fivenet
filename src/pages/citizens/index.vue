<script setup lang="ts">
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import CitizensList from '../../components/citizens/CitizensList.vue';
import Footer from '../../components/partials/Footer.vue';
</script>

<route lang="json">
{
    "name": "Citizens",
    "meta": {
        "requiresAuth": true,
        "permission": "CitizenStoreService.FindUsers",
        "breadCrumbs": [
            { "name": "Citizens", "href": "/citizens" }
        ]
    }
}
</route>

<template>
    <NavPageHeader title="Citizens" />
    <ContentWrapper>
        <CitizensList />
    </ContentWrapper>
    <Footer />
</template>
