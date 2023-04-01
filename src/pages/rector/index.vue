<script lang="ts" setup>
import ContentCenterWrapper from '../components/partials/ContentCenterWrapper.vue';

useHead({
    title: 'Control Panel',
});
definePageMeta({
    title: 'Control Panel',
    requiresAuth: true,
    permission: 'RectorService.GetRoles',
});
</script>

<template>
    <ContentCenterWrapper>
        <!-- TODO -->
        <h2>RectorService</h2>
    </ContentCenterWrapper>
</template>
