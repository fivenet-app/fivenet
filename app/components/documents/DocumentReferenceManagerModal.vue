<script setup lang="ts">
import { UButton, UTooltip } from '#components';
import type { TabsItem } from '@nuxt/ui';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DocumentInfoPopover from '~/components/partials/documents/DocumentInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { type ClipboardDocument, getDocument, useClipboardStore } from '~/stores/clipboard';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import { type DocumentReference, type DocumentShort, DocReference } from '~~/gen/ts/resources/documents/documents';

const props = defineProps<{
    documentId?: number;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const modelValue = defineModel<DocumentReference[]>('references', {
    type: Array,
    required: true,
});

const { t } = useI18n();

const clipboardStore = useClipboardStore();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const items = ref<TabsItem[]>([
    {
        label: t('components.documents.document_managers.view_current'),
        icon: 'i-mdi-file-search',
        slot: 'current' as const,
        value: 'current',
    },
    {
        label: t('common.clipboard'),
        icon: 'i-mdi-clipboard-list',
        slot: 'clipboard' as const,
        value: 'clipboard',
    },
    {
        label: t('components.documents.document_managers.add_new'),
        icon: 'i-mdi-file-document-plus',
        slot: 'new' as const,
        value: 'new',
    },
]);

const queryDoc = ref('');

const {
    data: documents,
    status,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-references-docs-${queryDoc.value}`, () => listDocuments());

watchDebounced(queryDoc, async () => await refresh(), {
    debounce: 200,
    maxWait: 1750,
});

async function listDocuments(): Promise<DocumentShort[]> {
    try {
        const call = documentsDocumentsClient.listDocuments({
            pagination: {
                offset: 0,
                pageSize: 8,
            },
            search: queryDoc.value,
            categoryIds: [],
            creatorIds: [],
            documentIds: [],
        });
        const { response } = await call;

        return response.documents.filter(
            (doc) => !modelValue.value.find((r) => r.targetDocumentId === doc.id || doc.id === props.documentId),
        );
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

let lastId = 0;

async function addReference(doc: DocumentShort, reference: DocReference): Promise<void> {
    modelValue.value.push({
        id: lastId--,
        sourceDocumentId: props.documentId ?? 0,
        reference: reference,
        targetDocumentId: doc.id,
        targetDocument: doc,
    });

    await refresh();
}

function addReferenceClipboard(doc: ClipboardDocument, reference: DocReference): void {
    addReference(getDocument(doc), reference);
}

function removeReference(id: number): void {
    const idx = modelValue.value.findIndex((r) => r.id === id);
    if (idx > -1) {
        modelValue.value.splice(idx, 1);
    }

    refresh();
}

const columnsCurrent = [
    {
        accessorKey: 'title',
        label: t('common.title'),
    },
    {
        accessorKey: 'creator',
        label: t('common.creator'),
    },
    {
        accessorKey: 'reference',
        label: t('common.reference', 1),
    },
    {
        accessorKey: 'actions',
        label: t('common.action', 2),
    },
];

const columnsClipboard = [
    {
        accessorKey: 'title',
        label: t('common.title'),
    },
    {
        accessorKey: 'creator',
        label: t('common.creator'),
    },
    {
        accessorKey: 'createdAt',
        label: t('common.created_at'),
    },
    {
        accessorKey: 'references',
        label: t('components.documents.document_managers.add_reference'),
    },
];

const columnsNew = [
    {
        accessorKey: 'title',
        label: t('common.title'),
    },
    {
        accessorKey: 'creator',
        label: t('common.creator'),
    },
    {
        accessorKey: 'createdAt',
        label: t('common.created_at'),
    },
    {
        accessorKey: 'references',
        label: t('components.documents.document_managers.add_reference'),
    },
];
</script>

<template>
    <UModal>
        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl leading-6 font-semibold">
                        {{ $t('common.document', 1) }}
                        {{ $t('common.reference', 2) }}
                    </h3>

                    <UButton
                        class="-my-1"
                        color="neutral"
                        variant="ghost"
                        icon="i-mdi-window-close"
                        @click="$emit('close', false)"
                    />
                </div>
            </template>

            <div>
                <UTabs :items="items">
                    <template #current>
                        <UTable
                            :columns="columnsCurrent"
                            :data="modelValue"
                            :empty-state="{ icon: 'i-mdi-file', label: $t('common.not_found', [$t('common.reference', 2)]) }"
                        >
                            <template #title-cell="{ row }">
                                <DocumentInfoPopover
                                    :document="!row.original.targetDocument?.id ? undefined : row.original.targetDocument"
                                    :document-id="row.original.targetDocumentId"
                                />
                            </template>

                            <template #creator-cell="{ row }">
                                <CitizenInfoPopover
                                    :user="
                                        !row.original.targetDocument?.creator ? undefined : row.original.targetDocument?.creator
                                    "
                                    :user-id="row.original.targetDocument?.creatorId"
                                    :trailing="false"
                                />
                            </template>

                            <template #reference-cell="{ row }">
                                {{ $t(`enums.documents.DocReference.${DocReference[row.original.reference]}`) }}
                            </template>

                            <template #actions-cell="{ row }">
                                <div class="flex flex-row gap-2">
                                    <UTooltip :text="$t('components.documents.document_managers.open_document')">
                                        <UButton
                                            :to="{
                                                name: 'documents-id',
                                                params: {
                                                    id: row.original.targetDocumentId,
                                                },
                                            }"
                                            target="_blank"
                                            variant="link"
                                            icon="i-mdi-open-in-new"
                                        />
                                    </UTooltip>

                                    <UTooltip :text="$t('components.documents.document_managers.remove_reference')">
                                        <UButton
                                            icon="i-mdi-file-document-minus"
                                            variant="link"
                                            color="error"
                                            @click="removeReference(row.original.id!)"
                                        />
                                    </UTooltip>
                                </div>
                            </template>
                        </UTable>
                    </template>

                    <template #clipboard>
                        <div>
                            <UTable
                                :columns="columnsClipboard"
                                :data="clipboardStore.$state.documents"
                                :empty-state="{
                                    icon: 'i-mdi-file',
                                    label: $t('common.not_found', [$t('common.document', 2)]),
                                }"
                            >
                                <template #title-cell="{ row }">
                                    <DocumentInfoPopover :document="getDocument(row.original)" />
                                </template>

                                <template #creator-cell="{ row }">
                                    <CitizenInfoPopover :user="row.original.creator" :trailing="false" />
                                </template>

                                <template #createdAt-cell="{ row }">
                                    <GenericTime :value="new Date(row.original.createdAt ?? 'now')" ago />
                                </template>

                                <template #references-cell="{ row }">
                                    <UButtonGroup>
                                        <UTooltip :text="$t('components.documents.document_managers.links')">
                                            <UButton
                                                color="blue"
                                                icon="i-mdi-link"
                                                @click="addReferenceClipboard(row.original, DocReference.LINKED)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.solves')">
                                            <UButton
                                                color="green"
                                                icon="i-mdi-check"
                                                @click="addReferenceClipboard(row.original, DocReference.SOLVES)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.closes')">
                                            <UButton
                                                color="error"
                                                icon="i-mdi-close-box"
                                                @click="addReferenceClipboard(row.original, DocReference.CLOSES)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.deprecates')">
                                            <UButton
                                                color="warning"
                                                icon="i-mdi-lock-clock"
                                                @click="addReferenceClipboard(row.original, DocReference.DEPRECATES)"
                                            />
                                        </UTooltip>
                                    </UButtonGroup>
                                </template>
                            </UTable>
                        </div>
                    </template>

                    <template #new>
                        <UFormField class="mb-2" name="title" :label="$t('common.search')">
                            <UInput
                                v-model="queryDoc"
                                type="text"
                                name="title"
                                :placeholder="`${$t('common.document', 1)} ${$t('common.title')}`"
                                leading-icon="i-mdi-search"
                            />
                        </UFormField>

                        <div>
                            <DataErrorBlock
                                v-if="error"
                                :title="$t('common.unable_to_load', [$t('common.document', 2)])"
                                :error="error"
                                :retry="refresh"
                            />
                            <UTable
                                v-else
                                :columns="columnsNew"
                                :loading="isRequestPending(status)"
                                :data="documents"
                                :empty-state="{
                                    icon: 'i-mdi-file',
                                    label: $t('common.not_found', [$t('common.reference', 2)]),
                                }"
                            >
                                <template #title-cell="{ row }">
                                    <DocumentInfoPopover :document="row.original" />
                                </template>

                                <template #creator-cell="{ row }">
                                    <CitizenInfoPopover :user="row.original.creator" :trailing="false" />
                                </template>

                                <template #createdAt-cell="{ row }">
                                    <GenericTime :value="row.original.createdAt" ago />
                                </template>

                                <template #references-cell="{ row }">
                                    <UButtonGroup>
                                        <UTooltip :text="$t('components.documents.document_managers.links')">
                                            <UButton
                                                color="blue"
                                                icon="i-mdi-link"
                                                @click="addReference(row.original, DocReference.LINKED)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.solves')">
                                            <UButton
                                                color="green"
                                                icon="i-mdi-check"
                                                @click="addReference(row.original, DocReference.SOLVES)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.closes')">
                                            <UButton
                                                color="error"
                                                icon="i-mdi-close-box"
                                                @click="addReference(row.original, DocReference.CLOSES)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.deprecates')">
                                            <UButton
                                                color="warning"
                                                icon="i-mdi-lock-clock"
                                                @click="addReference(row.original, DocReference.DEPRECATES)"
                                            />
                                        </UTooltip>
                                    </UButtonGroup>
                                </template>
                            </UTable>
                        </div>
                    </template>
                </UTabs>
            </div>

            <template #footer>
                <UButton class="flex-1" block color="neutral" @click="$emit('close', false)">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
