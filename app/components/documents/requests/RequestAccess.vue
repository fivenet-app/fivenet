<script lang="ts" setup>
import RequestAccessModal from '~/components/documents/requests/RequestAccessModal.vue';

defineProps<{
    documentId: number;
}>();

const { attr, can } = useAuth();

const overlay = useOverlay();
const documentRequestAccessModal = overlay.create(RequestAccessModal);
</script>

<template>
    <div
        v-if="
            can('documents.DocumentsService/CreateDocumentReq').value &&
            attr('documents.DocumentsService/CreateDocumentReq', 'Types', 'Access').value
        "
        class="mx-auto max-w-md rounded-md"
    >
        <UAlert
            color="primary"
            variant="subtle"
            icon="i-mdi-lock-question"
            :title="$t('components.documents.document_request_access.title')"
            :description="$t('components.documents.document_request_access.message')"
            :actions="[
                {
                    variant: 'solid',
                    color: 'primary',
                    label: $t('components.documents.document_request_access.callback_message'),
                    onClick: () =>
                        documentRequestAccessModal.open({
                            documentId: documentId,
                        }),
                },
            ]"
        />
    </div>
</template>
