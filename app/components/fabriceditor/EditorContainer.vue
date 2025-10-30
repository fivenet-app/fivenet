<script setup lang="ts">
import EditorToolbar from './EditorToolbar.vue';
import EditorWrapper from './EditorWrapper.vue';

withDefaults(
    defineProps<{
        maxHeight?: number;
        maxWidth?: number;
    }>(),
    {
        maxHeight: 350,
        maxWidth: 900,
    },
);
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar class="print:hidden" :title="$t('pages.documents.id.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton to="/documents" />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar class="p-1 print:hidden">
                <template #default>
                    <EditorToolbar />
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <EditorWrapper :max-height="maxHeight" :max-width="maxWidth" />
        </template>
    </UDashboardPanel>
</template>
