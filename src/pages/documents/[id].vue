<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import DocumentView from '../../components/documents/DocumentView.vue';
import ClipboardButton from '../../components/clipboard/ClipboardButton.vue';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';

useHead({
    title: 'Documents: View',
});
definePageMeta({
    title: 'Documents: View',
    requiresAuth: true,
    permission: 'DocStoreService.GetDocument',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('documents-id');
const documentId = ref(0);

onMounted(() => {
    documentId.value = parseInt(route.params.id);
});
</script>

<template>
    <ContentWrapper>
        <DocumentView v-if="documentId > 0" :documentId="documentId" />
        <ClipboardButton />
    </ContentWrapper>
</template>
