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
</script>

<template>
    <UDashboardPanel>
        <template #body>
            <SignaturePad />
        </template>
    </UDashboardPanel>
</template>
