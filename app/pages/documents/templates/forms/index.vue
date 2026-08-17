<script setup lang="ts">
useHead({
    title: 'Document Forms',
});

definePageMeta({
    title: 'Document Forms',
    requiresAuth: true,
    permission: 'TODOService/TODOMethod',
});
</script>

<template>
    <UDashboardPanel :ui="{ root: 'pb-(--page-content-bottom-offset)' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.documents.templates.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton to="/documents" />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div>TODO Forms List here</div>
        </template>
    </UDashboardPanel>
</template>
