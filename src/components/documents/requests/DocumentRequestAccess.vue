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
            title="You don't have access to this document"
            message="But you can request access to this document using the button below."
            type="info"
            :icon="markRaw(LockQuestionIcon)"
            callback-message="Request Access here"
            @clicked="open = true"
        />
    </div>
</template>
