<script lang="ts" setup>
import AuditLog from '~/components/settings/AuditLog.vue';

useHead({
    title: 'common.audit_log',
});

definePageMeta({
    title: 'common.audit_log',
    requiresAuth: true,
    permission: 'settings.SettingsService/ViewAuditLog',
});
</script>

<template>
    <AuditLog />
</template>
