<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import TemplatesList from '~/components/documents/templates/TemplatesList.vue';

useHead({
    title: 'Templates',
});
definePageMeta({
    title: 'Templates',
    requiresAuth: true,
    permission: 'DocStoreService.FindDocuments',
});
</script>

<template>
    <ContentWrapper>
        <div class="pt-2">
            <TemplatesList />
        </div>
    </ContentWrapper>
</template>
