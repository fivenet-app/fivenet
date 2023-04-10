<script lang="ts" setup>
import AuditLog from '~/components/rector/AuditLog.vue';

useHead({
    title: 'Audit Log',
});
definePageMeta({
    title: 'Audit Log',
    requiresAuth: true,
    permission: 'RectorService.ViewAuditLog',
});
</script>

<template>
    <div>
        <AuditLog />
    </div>
</template>
