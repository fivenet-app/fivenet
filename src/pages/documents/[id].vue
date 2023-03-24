<script lang="ts" setup>
import { useRoute } from 'vue-router/auto';
import { ref, onMounted } from 'vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import DocumentView from '../../components/documents/DocumentView.vue';
import Footer from '../../components/partials/Footer.vue';

const route = useRoute<'Documents: Info'>();
const documentId = ref(0);

onMounted(() => {
    documentId.value = parseInt(route.params.id);
});
</script>

<route lang="json">
{
    "name": "Documents: Info",
    "meta": {
        "requiresAuth": true,
        "permission": "DocStoreService.GetDocument"
    }
}
</route>

<template>
    <ContentWrapper>
        <DocumentView v-if="documentId > 0" :documentId="documentId" />
    </ContentWrapper>
    <Footer />
</template>
