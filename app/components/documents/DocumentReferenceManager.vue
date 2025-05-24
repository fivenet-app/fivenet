<script setup lang="ts">
import type { TabItem } from '#ui/types';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DocumentInfoPopover from '~/components/partials/documents/DocumentInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { ClipboardDocument } from '~/stores/clipboard';
import { getDocument, useClipboardStore } from '~/stores/clipboard';
import type { DocumentReference, DocumentShort } from '~~/gen/ts/resources/documents/documents';
import { DocReference } from '~~/gen/ts/resources/documents/documents';

const props = defineProps<{
    open: boolean;
    documentId?: number;
    modelValue: Map<number, DocumentReference>;
}>();

defineEmits<{
    (e: 'close'): void;
    (e: 'update:modelValue', payload: Map<number, DocumentReference>): void;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const clipboardStore = useClipboardStore();

const items = ref<TabItem[]>([
    {
        label: t('components.documents.document_managers.view_current'),
        icon: 'i-mdi-file-search',
        slot: 'current',
    },
    {
        label: t('common.clipboard'),
        icon: 'i-mdi-clipboard-list',
        slot: 'clipboard',
    },
    {
        label: t('components.documents.document_managers.add_new'),
        icon: 'i-mdi-file-document-plus',
        slot: 'new',
    },
]);

const queryDoc = ref('');

const {
    data: documents,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-references-docs-${queryDoc.value}`, () => listDocuments());

watchDebounced(queryDoc, async () => await refresh(), {
    debounce: 200,
    maxWait: 1750,
});

async function listDocuments(): Promise<DocumentShort[]> {
    try {
        const call = $grpc.documents.documents.listDocuments({
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
            (doc) =>
                !Array.from(props.modelValue.values()).find(
                    (r) => r.targetDocumentId === doc.id || doc.id === props.documentId,
                ),
        );
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function addReference(doc: DocumentShort, reference: DocReference): Promise<void> {
    const keys = Array.from(props.modelValue.keys());
    const key = !keys.length ? 1 : keys[keys.length - 1]! + 1;

    props.modelValue.set(key, {
        id: key,
        sourceDocumentId: props.documentId ?? 0,
        reference,
        targetDocumentId: doc.id,
        targetDocument: doc,
    });

    await refresh();
}

function addReferenceClipboard(doc: ClipboardDocument, reference: DocReference): void {
    addReference(getDocument(doc), reference);
}

function removeReference(id: number): void {
    props.modelValue.delete(id);
    listDocuments();
}

const columnsCurrent = [
    {
        key: 'title',
        label: t('common.title'),
    },
    {
        key: 'creator',
        label: t('common.creator'),
    },
    {
        key: 'reference',
        label: t('common.reference', 1),
    },
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
];

const columnsClipboard = [
    {
        key: 'title',
        label: t('common.title'),
    },
    {
        key: 'creator',
        label: t('common.creator'),
    },
    {
        key: 'createdAt',
        label: t('common.created_at'),
    },
    {
        key: 'references',
        label: t('components.documents.document_managers.add_reference'),
        sortable: false,
    },
];

const columnsNew = [
    {
        key: 'title',
        label: t('common.title'),
    },
    {
        key: 'creator',
        label: t('common.creator'),
    },
    {
        key: 'createdAt',
        label: t('common.created_at'),
    },
    {
        key: 'references',
        label: t('components.documents.document_managers.add_reference'),
        sortable: false,
    },
];
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }" :model-value="open" @update:model-value="$emit('close')">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.document', 1) }}
                        {{ $t('common.reference', 2) }}
                    </h3>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="$emit('close')" />
                </div>
            </template>

            <div>
                <UTabs :items="items">
                    <template #current>
                        <UTable
                            :columns="columnsCurrent"
                            :rows="[...modelValue.values()]"
                            :empty-state="{ icon: 'i-mdi-file', label: $t('common.not_found', [$t('common.reference', 2)]) }"
                        >
                            <template #title-data="{ row }">
                                <DocumentInfoPopover :document="row.targetDocument" />
                            </template>

                            <template #creator-data="{ row }">
                                <CitizenInfoPopover :user="row.targetDocument?.creator" :trailing="false" />
                            </template>

                            <template #reference-data="{ row }">
                                {{ $t(`enums.documents.DocReference.${DocReference[row.reference]}`) }}
                            </template>

                            <template #actions-data="{ row }">
                                <div class="flex flex-row gap-2">
                                    <UTooltip :text="$t('components.documents.document_managers.open_document')">
                                        <UButton
                                            :to="{
                                                name: 'documents-id',
                                                params: {
                                                    id: row.targetDocumentId,
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
                                            @click="removeReference(row.id!)"
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
                                :rows="clipboardStore.$state.documents"
                                :empty-state="{
                                    icon: 'i-mdi-file',
                                    label: $t('common.not_found', [$t('common.document', 2)]),
                                }"
                            >
                                <template #title-data="{ row }">
                                    <DocumentInfoPopover :document="getDocument(row)" />
                                </template>

                                <template #creator-data="{ row }">
                                    <CitizenInfoPopover :user="row.creator" :trailing="false" />
                                </template>

                                <template #createdAt-data="{ row }">
                                    <GenericTime :value="row.createdAt" ago />
                                </template>

                                <template #references-data="{ row }">
                                    <UButtonGroup>
                                        <UTooltip :text="$t('components.documents.document_managers.links')">
                                            <UButton
                                                color="blue"
                                                icon="i-mdi-link"
                                                @click="addReferenceClipboard(row, DocReference.LINKED)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.solves')">
                                            <UButton
                                                color="green"
                                                icon="i-mdi-check"
                                                @click="addReferenceClipboard(row, DocReference.SOLVES)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.closes')">
                                            <UButton
                                                color="error"
                                                icon="i-mdi-close-box"
                                                @click="addReferenceClipboard(row, DocReference.CLOSES)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.deprecates')">
                                            <UButton
                                                color="amber"
                                                icon="i-mdi-lock-clock"
                                                @click="addReferenceClipboard(row, DocReference.DEPRECATES)"
                                            />
                                        </UTooltip>
                                    </UButtonGroup>
                                </template>
                            </UTable>
                        </div>
                    </template>

                    <template #new>
                        <UFormGroup class="mb-2" name="title" :label="$t('common.search')">
                            <UInput
                                v-model="queryDoc"
                                type="text"
                                name="title"
                                :placeholder="`${$t('common.document', 1)} ${$t('common.title')}`"
                                leading-icon="i-mdi-search"
                            />
                        </UFormGroup>

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
                                :loading="loading"
                                :rows="documents"
                                :empty-state="{
                                    icon: 'i-mdi-file',
                                    label: $t('common.not_found', [$t('common.reference', 2)]),
                                }"
                            >
                                <template #title-data="{ row }">
                                    <DocumentInfoPopover :document="row" />
                                </template>

                                <template #creator-data="{ row }">
                                    <CitizenInfoPopover :user="row.creator" :trailing="false" />
                                </template>

                                <template #createdAt-data="{ row }">
                                    <GenericTime :value="row.createdAt" ago />
                                </template>

                                <template #references-data="{ row }">
                                    <UButtonGroup>
                                        <UTooltip :text="$t('components.documents.document_managers.links')">
                                            <UButton
                                                color="blue"
                                                icon="i-mdi-link"
                                                @click="addReference(row, DocReference.LINKED)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.solves')">
                                            <UButton
                                                color="green"
                                                icon="i-mdi-check"
                                                @click="addReference(row, DocReference.SOLVES)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.closes')">
                                            <UButton
                                                color="error"
                                                icon="i-mdi-close-box"
                                                @click="addReference(row, DocReference.CLOSES)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.deprecates')">
                                            <UButton
                                                color="amber"
                                                icon="i-mdi-lock-clock"
                                                @click="addReference(row, DocReference.DEPRECATES)"
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
                <UButton class="flex-1" block color="black" @click="$emit('close')">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
