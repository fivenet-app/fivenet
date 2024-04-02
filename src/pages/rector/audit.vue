<script lang="ts" setup>
import AuditLog from '~/components/rector/audit/AuditLog.vue';

useHead({
    title: 'common.audit_log',
});
definePageMeta({
    title: 'common.audit_log',
    requiresAuth: true,
    permission: 'RectorService.ViewAuditLog',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('common.audit_log')"> </UDashboardNavbar>

            <AuditLog />
        </UDashboardPanel>
    </UDashboardPage>
</template>
