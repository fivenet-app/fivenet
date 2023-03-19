<script setup lang="ts">
import { useRoute } from 'vue-router/auto';
import { ref, onBeforeUnmount } from 'vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import DocumentView from '../../components/documents/DocumentView.vue';
import Footer from '../../components/partials/Footer.vue';

const route = useRoute();
const documentID = ref(0);

onBeforeUnmount(() => {
    //@ts-ignore it is defined but the router gen isn't correctly typping it
    documentID.value = +route.params.id;
});
</script>

<route lang="json">
{
    "name": "Documents: Info",
    "meta": {
        "requiresAuth": true,
        "permission": "DocStoreService.GetDocument",
        "breadCrumbs": [
            { "name": "Documents", "href": "/documents" },
            { "name": "Document Info: ...", "href": "/documents" }
        ]
    }
}
</route>

<template>
    <NavPageHeader title="Document" />
    <ContentWrapper>
        <DocumentView v-if="documentID > 0" :documentID="documentID" />
    </ContentWrapper>
    <Footer />
</template>
