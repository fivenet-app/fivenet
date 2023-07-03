<script lang="ts" setup>
import AuditLog from '~/components/rector/audit/AuditLog.vue';

useHead({
    title: 'common.audit_log',
});
definePageMeta({
    title: 'common.audit_log',
    requiresAuth: true,
    permission: 'RectorService.ViewAuditLog',
    showQuickButtons: false,
});
</script>

<template>
    <div>
        <AuditLog />
    </div>
</template>
