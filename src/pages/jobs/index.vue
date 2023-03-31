<script lang="ts" setup>
import ContentWrapper from '../../components/partials/ContentWrapper.vue';

useHead({
    title: 'Jobs',
});
definePageMeta({
    requiresAuth: true,
    permission: 'Jobs.View',
});
</script>

<template>
    <ContentWrapper>
        <h1 class="text-lg font-semibold leading-6 text-neutral">EMPLOYEE INFO AND MANAGEMENT SYSTEM</h1>
    </ContentWrapper>
</template>
