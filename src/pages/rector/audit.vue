<script lang="ts" setup>
import AuditLog from '~/components/rector/AuditLog.vue';

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
            <UDashboardNavbar :title="$t('common.audit_log')">
                <template #right>
                    <UButton color="black" icon="i-mdi-arrow-back" to="/rector">
                        {{ $t('common.back') }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <AuditLog />
        </UDashboardPanel>
    </UDashboardPage>
</template>
