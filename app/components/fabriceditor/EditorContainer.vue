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

const modelValue = defineModel<string>({ required: true });
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar class="print:hidden" :title="$t('pages.documents.stamps.title')">
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
            <EditorWrapper v-model="modelValue" :max-height="maxHeight" :max-width="maxWidth" />
        </template>
    </UDashboardPanel>
</template>
