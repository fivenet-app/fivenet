<script lang="ts">
import { defineComponent } from 'vue';

import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import Footer from '../../components/partials/Footer.vue';
import DocumentsList from '../../components/documents/DocumentsList.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';

export default defineComponent({
    components: {
        Footer,
        DocumentsList,
        ContentWrapper,
        NavPageHeader,
    },
});
</script>

<route lang="json">
{
    "name": "Documents",
    "meta": {
        "requiresAuth": true,
        "permission": "DocStoreService.FindDocuments",
        "breadCrumbs": [
            { "name": "Documents", "href": "/documents" }
        ]
    }
}
</route>

<template>
    <NavPageHeader title="Documents" />
    <ContentWrapper>
        <DocumentsList />
    </ContentWrapper>
    <Footer />
</template>
