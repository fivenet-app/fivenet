<script lang="ts" setup>
import { useRoute } from 'vue-router/auto';
import { ref, onMounted } from 'vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import DocumentView from '../../components/documents/DocumentView.vue';
import Footer from '../../components/partials/Footer.vue';

const route = useRoute<'Documents: Info'>();
const documentID = ref(0);

onMounted(() => {
    documentID.value = parseInt(route.params.id);
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
    <NavPageHeader title="Document" />
    <ContentWrapper>
        <DocumentView v-if="documentID > 0" :documentID="documentID" />
    </ContentWrapper>
    <Footer />
</template>
