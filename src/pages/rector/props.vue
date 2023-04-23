<script lang="ts" setup>
import JobProps from '~/components/rector/JobProps.vue';

useHead({
    title: 'Job Properties',
});
definePageMeta({
    title: 'Job Properties',
    requiresAuth: true,
    permission: 'RectorService.ViewAuditLog',
});
</script>

<template>
    <div>
        <JobProps />
    </div>
</template>
