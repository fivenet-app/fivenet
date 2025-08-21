<script lang="ts" setup>
import CategoryList from '~/components/documents/categories/CategoryList.vue';

useHead({
    title: 'pages.documents.categories.title',
});

definePageMeta({
    title: 'pages.documents.categories.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService/ListCategories',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <CategoryList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
