<script lang="ts" setup>
import CategoriesList from '~/components/documents/categories/CategoriesList.vue';

useHead({
    title: 'pages.documents.categories.title',
});
definePageMeta({
    title: 'pages.documents.categories.title',
    requiresAuth: true,
    permission: 'DocStoreService.ListCategories',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <CategoriesList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
