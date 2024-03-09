<script lang="ts" setup>
import { LockQuestionIcon } from 'mdi-vue3';
import GenericAlert from '~/components/partials/elements/GenericAlert.vue';
import DocumentRequestAccessModal from '~/components/documents/requests/DocumentRequestAccessModal.vue';

defineProps<{
    documentId: string;
}>();

const open = ref(false);
</script>

<template v-if="can('DocStoreService.CreateDocumentReq') && attr('DocStoreService.CreateDocumentReq', 'Types', 'Access')">
    <div class="mx-auto max-w-md rounded-md">
        <DocumentRequestAccessModal :document-id="documentId" :open="open" @close="open = false" />

        <GenericAlert
            type="info"
            :title="$t('components.documents.document_request_access.title')"
            :message="$t('components.documents.document_request_access.message')"
            :icon="markRaw(LockQuestionIcon)"
            :callback-message="$t('components.documents.document_request_access.callback_message')"
            @clicked="open = true"
        />
    </div>
</template>
