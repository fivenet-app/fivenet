<script lang="ts" setup>
import SignaturePad from '~/components/partials/SignaturePad.vue';

useHead({
    title: 'pages.documents.title',
});

definePageMeta({
    title: 'pages.documents.title',
    requiresAuth: true,
    permission: 'TODOService/TODOMethod',
});

const signature = ref('');
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar>
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <SignaturePad v-model="signature" />
        </template>
    </UDashboardPanel>
</template>
