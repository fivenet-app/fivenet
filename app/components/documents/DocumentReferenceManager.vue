<script setup lang="ts">
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { ClipboardDocument } from '~/store/clipboard';
import { getDocument, getUser, useClipboardStore } from '~/store/clipboard';
import type { DocumentReference, DocumentShort } from '~~/gen/ts/resources/documents/documents';
import { DocReference } from '~~/gen/ts/resources/documents/documents';

const props = defineProps<{
    open: boolean;
    documentId?: string;
    modelValue: Map<string, DocumentReference>;
}>();

defineEmits<{
    (e: 'close'): void;
    (e: 'update:modelValue', payload: Map<string, DocumentReference>): void;
}>();

const { t } = useI18n();

const clipboardStore = useClipboardStore();

const tabs = ref<{ key: string; label: string; icon: string }[]>([
    {
        key: 'current',
        label: t('components.documents.document_managers.view_current'),
        icon: 'i-mdi-file-search',
    },
    {
        key: 'clipboard',
        label: t('common.clipboard'),
        icon: 'i-mdi-clipboard-list',
    },
    {
        key: 'new',
        label: t('components.documents.document_managers.add_new'),
        icon: 'i-mdi-file-document-plus',
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
        const call = getGRPCDocStoreClient().listDocuments({
            pagination: {
                offset: 0,
                pageSize: 8,
            },
            orderBy: [],
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

function addReference(doc: DocumentShort, reference: DocReference): void {
    const keys = Array.from(props.modelValue.keys());
    const key = !keys.length ? '1' : (parseInt(keys[keys.length - 1]!) + 1).toString();

    props.modelValue.set(key, {
        id: key,
        sourceDocumentId: props.documentId ?? '0',
        reference,
        targetDocumentId: doc.id,
        targetDocument: doc,
    });
    refresh();
}

function addReferenceClipboard(doc: ClipboardDocument, reference: DocReference): void {
    addReference(getDocument(doc), reference);
}

function removeReference(id: string): void {
    props.modelValue.delete(id);
    listDocuments();
}
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

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="$emit('close')" />
                </div>
            </template>

            <div>
                <UTabs :items="tabs">
                    <template #item="{ item }">
                        <template v-if="item.key === 'current'">
                            <div class="flow-root">
                                <div class="-my-2 mx-0 overflow-x-auto">
                                    <div class="inline-block min-w-full py-2 align-middle">
                                        <table class="min-w-full divide-y divide-base-200">
                                            <thead>
                                                <tr>
                                                    <th
                                                        scope="col"
                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-6 lg:pl-8"
                                                    >
                                                        {{ $t('common.title') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.state') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.creator') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.reference') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.action', 2) }}
                                                    </th>
                                                </tr>
                                            </thead>
                                            <tbody class="divide-y divide-base-500">
                                                <tr v-for="[key, reference] in modelValue" :key="key.toString()">
                                                    <td
                                                        class="max-w-xl truncate whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-6 lg:pl-8"
                                                    >
                                                        {{ reference.targetDocument?.title }}
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        {{ reference.targetDocument?.state }}
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        <CitizenInfoPopover
                                                            :user="reference.targetDocument?.creator"
                                                            :trailing="false"
                                                        />
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        {{
                                                            $t(
                                                                `enums.docstore.DocReference.${
                                                                    DocReference[reference.reference]
                                                                }`,
                                                            )
                                                        }}
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        <div class="flex flex-row gap-2">
                                                            <UButton
                                                                :to="{
                                                                    name: 'documents-id',
                                                                    params: {
                                                                        id: reference.targetDocumentId,
                                                                    },
                                                                }"
                                                                target="_blank"
                                                                :title="
                                                                    $t('components.documents.document_managers.open_document')
                                                                "
                                                                variant="link"
                                                                icon="i-mdi-open-in-new"
                                                            />

                                                            <UButton
                                                                :title="
                                                                    $t(
                                                                        'components.documents.document_managers.remove_reference',
                                                                    )
                                                                "
                                                                icon="i-mdi-file-document-minus"
                                                                color="red"
                                                                @click="removeReference(reference.id!)"
                                                            />
                                                        </div>
                                                    </td>
                                                </tr>
                                            </tbody>
                                        </table>
                                    </div>
                                </div>
                            </div>
                        </template>
                        <template v-else-if="item.key === 'clipboard'">
                            <div class="mt-2 flow-root">
                                <div class="-my-2 mx-0 overflow-x-auto">
                                    <div class="inline-block min-w-full py-2 align-middle">
                                        <DataNoDataBlock
                                            v-if="clipboardStore.$state.documents.length === 0"
                                            :type="$t('common.reference', 2)"
                                            icon="i-mdi-file-document-multiple"
                                        />
                                        <table v-else class="min-w-full divide-y divide-base-200">
                                            <thead>
                                                <tr>
                                                    <th
                                                        scope="col"
                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-6 lg:pl-8"
                                                    >
                                                        {{ $t('common.title') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.state') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.creator') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.created_at') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('components.documents.document_managers.add_reference') }}
                                                    </th>
                                                </tr>
                                            </thead>
                                            <tbody class="divide-y divide-base-500">
                                                <tr v-for="document in clipboardStore.$state.documents" :key="document.id">
                                                    <td
                                                        class="max-w-xl truncate whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-6 lg:pl-8"
                                                    >
                                                        {{ document.title }}
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        {{ document.state }}
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        <CitizenInfoPopover
                                                            :user="getUser(document.creator)"
                                                            :trailing="false"
                                                        />
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        {{ $t('common.created') }}
                                                        <GenericTime :value="new Date(Date.parse(document.createdAt ?? ''))" />
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        <div class="flex flex-row gap-2">
                                                            <UButton
                                                                :title="$t('components.documents.document_managers.links')"
                                                                color="blue"
                                                                icon="i-mdi-link"
                                                                @click="addReferenceClipboard(document, DocReference.LINKED)"
                                                            />

                                                            <UButton
                                                                :title="$t('components.documents.document_managers.solves')"
                                                                color="green"
                                                                icon="i-mdi-check"
                                                                @click="addReferenceClipboard(document, DocReference.SOLVES)"
                                                            />

                                                            <UButton
                                                                :title="$t('components.documents.document_managers.closes')"
                                                                color="red"
                                                                icon="i-mdi-close-box"
                                                                @click="addReferenceClipboard(document, DocReference.CLOSES)"
                                                            />

                                                            <UButton
                                                                :title="$t('components.documents.document_managers.deprecates')"
                                                                color="amber"
                                                                icon="i-mdi-lock-clock"
                                                                @click="
                                                                    addReferenceClipboard(document, DocReference.DEPRECATES)
                                                                "
                                                            />
                                                        </div>
                                                    </td>
                                                </tr>
                                            </tbody>
                                        </table>
                                    </div>
                                </div>
                            </div>
                        </template>
                        <template v-else-if="item.key === 'new'">
                            <UFormGroup name="title" :label="$t('common.search')">
                                <UInput
                                    v-model="queryDoc"
                                    type="text"
                                    name="title"
                                    :placeholder="`${$t('common.document', 1)} ${$t('common.title')}`"
                                />
                            </UFormGroup>

                            <div class="mt-2 flow-root">
                                <div class="-my-2 mx-0 overflow-x-auto">
                                    <div class="inline-block min-w-full py-2 align-middle">
                                        <DataPendingBlock
                                            v-if="loading"
                                            :message="$t('common.loading', [$t('common.document', 2)])"
                                        />
                                        <DataErrorBlock
                                            v-else-if="error"
                                            :title="$t('common.unable_to_load', [$t('common.document', 2)])"
                                            :retry="refresh"
                                        />
                                        <DataNoDataBlock
                                            v-else-if="!documents || documents.length === 0"
                                            :message="$t('components.citizens.CitizensList.no_citizens')"
                                        />
                                        <table v-else class="min-w-full divide-y divide-base-200">
                                            <thead>
                                                <tr>
                                                    <th
                                                        scope="col"
                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-6 lg:pl-8"
                                                    >
                                                        {{ $t('common.title') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.state') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.creator') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.created_at') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('components.documents.document_managers.add_reference') }}
                                                    </th>
                                                </tr>
                                            </thead>
                                            <tbody class="divide-y divide-base-500">
                                                <tr v-for="document in documents.slice(0, 8)" :key="document.id">
                                                    <td
                                                        class="max-w-xl truncate whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-6 lg:pl-8"
                                                    >
                                                        {{ document.title }}
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        {{ document.state }}
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        <CitizenInfoPopover :user="document.creator" :trailing="false" />
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        <GenericTime :value="document.createdAt" :ago="true" />
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        <div class="flex flex-row gap-2">
                                                            <UButton
                                                                :title="$t('components.documents.document_managers.links')"
                                                                color="blue"
                                                                icon="i-mdi-link"
                                                                @click="addReference(document, DocReference.LINKED)"
                                                            />

                                                            <UButton
                                                                :title="$t('components.documents.document_managers.solves')"
                                                                color="green"
                                                                icon="i-mdi-check"
                                                                @click="addReference(document, DocReference.SOLVES)"
                                                            />

                                                            <UButton
                                                                :title="$t('components.documents.document_managers.closes')"
                                                                color="red"
                                                                icon="i-mdi-close-box"
                                                                @click="addReference(document, DocReference.CLOSES)"
                                                            />

                                                            <UButton
                                                                :title="$t('components.documents.document_managers.deprecates')"
                                                                color="amber"
                                                                icon="i-mdi-lock-clock"
                                                                @click="addReference(document, DocReference.DEPRECATES)"
                                                            />
                                                        </div>
                                                    </td>
                                                </tr>
                                            </tbody>
                                        </table>
                                    </div>
                                </div>
                            </div>
                        </template>
                    </template>
                </UTabs>
            </div>

            <template #footer>
                <UButton block class="flex-1" color="black" @click="$emit('close')">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
