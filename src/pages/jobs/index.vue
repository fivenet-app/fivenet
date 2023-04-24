<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';

useHead({
    title: 'Jobs',
});
definePageMeta({
    title: 'Jobs',
    requiresAuth: true,
    permission: 'Jobs.View',
});
</script>

<template>
    <ContentWrapper>
        <h1 class="text-lg font-semibold leading-6 text-neutral">{{ $t('pages.jobs.employee_info_and_management_system') }}</h1>
    </ContentWrapper>
</template>
