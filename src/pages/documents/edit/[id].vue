<script lang="ts" setup>
import ContentWrapper from '../../../components/partials/ContentWrapper.vue';
import DocumentEditor from '../../../components/documents/DocumentEditor.vue';
import Footer from '../../../components/partials/Footer.vue';
import { useRoute } from 'vue-router/auto';

const route = useRoute<'Documents: Edit'>();
</script>

<route lang="json">
{
    "name": "Documents: Edit",
    "meta": {
        "requiresAuth": true,
        "permission": "DocStoreService.CreateDocument"
    }
}
</route>

<template>
    <ContentWrapper>
        <DocumentEditor :id="parseInt(route.params.id)" />
    </ContentWrapper>
    <Footer />
</template>
