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
    <UDashboardPanel>
        <UDashboardNavbar :title="$t('common.audit_log')">
            <template #leading>
                <UDashboardSidebarCollapse />
            </template>

            <template #right>
                <PartialsBackButton fallback-to="/settings" />
            </template>
        </UDashboardNavbar>

        <AuditLog />
    </UDashboardPanel>
</template>
