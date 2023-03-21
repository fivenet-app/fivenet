<script setup lang="ts">
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import Footer from '../../components/partials/Footer.vue';
</script>

<route lang="json">
{
    "name": "Jobs",
    "meta": {
        "requiresAuth": true,
        "permission": "Jobs.View"
    }
}
</route>

<template>
    <NavPageHeader title="Job" />
    <ContentWrapper>
        <h1 class="text-lg font-semibold leading-6 text-white">EMPLOYEE INFO AND MANAGEMENT SYSTEM</h1>
    </ContentWrapper>
    <Footer />
</template>
