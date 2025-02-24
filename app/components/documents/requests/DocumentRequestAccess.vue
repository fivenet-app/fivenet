<script lang="ts" setup>
import DocumentRequestAccessModal from '~/components/documents/requests/DocumentRequestAccessModal.vue';

defineProps<{
    documentId: number;
}>();

const { attr, can } = useAuth();

const modal = useModal();
</script>

<template>
    <div
        v-if="
            can('DocStoreService.CreateDocumentReq').value && attr('DocStoreService.CreateDocumentReq', 'Types', 'Access').value
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
                    click: () =>
                        modal.open(DocumentRequestAccessModal, {
                            documentId: documentId,
                        }),
                },
            ]"
        />
    </div>
</template>
