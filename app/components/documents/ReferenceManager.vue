<script setup lang="ts">
import { UButton, UFieldGroup, UTooltip } from '#components';
import type { TableColumn, TabsItem } from '@nuxt/ui';
import { h } from 'vue';
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

const modelValue = defineModel<DocumentReference[]>({
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

const columnsCurrent = computed(
    () =>
        [
            {
                accessorKey: 'title',
                header: t('common.title'),
                cell: ({ row }) =>
                    h(DocumentInfoPopover, {
                        document: !row.original.targetDocument?.id ? undefined : row.original.targetDocument,
                        documentId: row.original.targetDocumentId,
                    }),
            },
            {
                accessorKey: 'creator',
                header: t('common.creator'),
                cell: ({ row }) =>
                    h(CitizenInfoPopover, {
                        user: !row.original.targetDocument?.creator ? undefined : row.original.targetDocument?.creator,
                        userId: row.original.targetDocument?.creatorId,
                        trailing: false,
                    }),
            },
            {
                accessorKey: 'reference',
                header: t('common.reference', 1),
                cell: ({ row }) => t(`enums.documents.DocReference.${DocReference[row.original.reference]}`),
            },
            {
                id: 'actions',
                cell: ({ row }) =>
                    h('div', { class: 'flex flex-row gap-2' }, [
                        h(
                            UTooltip,
                            { text: t('components.documents.document_managers.open_document') },
                            {
                                default: () =>
                                    h(UButton, {
                                        to: {
                                            name: 'documents-id',
                                            params: { id: row.original.targetDocumentId },
                                        },
                                        target: '_blank',
                                        variant: 'link',
                                        icon: 'i-mdi-open-in-new',
                                    }),
                            },
                        ),
                        h(
                            UTooltip,
                            { text: t('components.documents.document_managers.remove_reference') },
                            {
                                default: () =>
                                    h(UButton, {
                                        icon: 'i-mdi-file-document-minus',
                                        variant: 'link',
                                        color: 'error',
                                        onClick: () => removeReference(row.original.id!),
                                    }),
                            },
                        ),
                    ]),
            },
        ] as TableColumn<DocumentReference>[],
);

const columnsClipboard = computed(
    () =>
        [
            {
                accessorKey: 'title',
                header: t('common.title'),
                cell: ({ row }) => h(DocumentInfoPopover, { document: getDocument(row.original) }),
            },
            {
                accessorKey: 'creator',
                header: t('common.creator'),
                cell: ({ row }) => h(CitizenInfoPopover, { user: row.original.creator, trailing: false }),
            },
            {
                accessorKey: 'createdAt',
                header: t('common.created_at'),
                cell: ({ row }) => h(GenericTime, { value: new Date(row.original.createdAt ?? 'now'), ago: true }),
            },
            {
                accessorKey: 'references',
                header: t('components.documents.document_managers.add_reference'),
                cell: ({ row }) =>
                    h(
                        UFieldGroup,
                        {},
                        {
                            default: () => [
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.links') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'blue',
                                                icon: 'i-mdi-link',
                                                onClick: () => addReferenceClipboard(row.original, DocReference.LINKED),
                                            }),
                                    },
                                ),
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.solves') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'green',
                                                icon: 'i-mdi-check',
                                                onClick: () => addReferenceClipboard(row.original, DocReference.SOLVES),
                                            }),
                                    },
                                ),
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.closes') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'error',
                                                icon: 'i-mdi-close-box',
                                                onClick: () => addReferenceClipboard(row.original, DocReference.CLOSES),
                                            }),
                                    },
                                ),
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.deprecates') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'warning',
                                                icon: 'i-mdi-lock-clock',
                                                onClick: () => addReferenceClipboard(row.original, DocReference.DEPRECATES),
                                            }),
                                    },
                                ),
                            ],
                        },
                    ),
            },
        ] as TableColumn<ClipboardDocument>[],
);

const columnsNew = computed(
    () =>
        [
            {
                accessorKey: 'title',
                header: t('common.title'),
                cell: ({ row }) => h(DocumentInfoPopover, { document: row.original }),
            },
            {
                accessorKey: 'creator',
                header: t('common.creator'),
                cell: ({ row }) => h(CitizenInfoPopover, { user: row.original.creator, trailing: false }),
            },
            {
                accessorKey: 'createdAt',
                header: t('common.created_at'),
                cell: ({ row }) => h(GenericTime, { value: row.original.createdAt, ago: true }),
            },
            {
                accessorKey: 'references',
                header: t('components.documents.document_managers.add_reference'),
                cell: ({ row }) =>
                    h(
                        UFieldGroup,
                        {},
                        {
                            default: () => [
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.links') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'blue',
                                                icon: 'i-mdi-link',
                                                onClick: () => addReference(row.original, DocReference.LINKED),
                                            }),
                                    },
                                ),
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.solves') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'green',
                                                icon: 'i-mdi-check',
                                                onClick: () => addReference(row.original, DocReference.SOLVES),
                                            }),
                                    },
                                ),
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.closes') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'error',
                                                icon: 'i-mdi-close-box',
                                                onClick: () => addReference(row.original, DocReference.CLOSES),
                                            }),
                                    },
                                ),
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.deprecates') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'warning',
                                                icon: 'i-mdi-lock-clock',
                                                onClick: () => addReference(row.original, DocReference.DEPRECATES),
                                            }),
                                    },
                                ),
                            ],
                        },
                    ),
            },
        ] as TableColumn<DocumentShort>[],
);
</script>

<template>
    <UTabs default-value="current" :items="items" variant="link">
        <template #current>
            <UTable :columns="columnsCurrent" :data="modelValue" :empty="$t('common.not_found', [$t('common.reference', 2)])" />
        </template>

        <template #clipboard>
            <UTable
                :columns="columnsClipboard"
                :data="clipboardStore.$state.documents"
                :empty="$t('common.not_found', [$t('common.document', 2)])"
            />
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
                    :empty="$t('common.not_found', [$t('common.reference', 2)])"
                />
            </div>
        </template>
    </UTabs>
</template>
