<script setup lang="ts">
import {
    Dialog,
    DialogPanel,
    DialogTitle,
    Tab,
    TabGroup,
    TabList,
    TabPanel,
    TabPanels,
    TransitionChild,
    TransitionRoot,
} from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { watchDebounced } from '@vueuse/core';
import {
    CheckIcon,
    ClipboardListIcon,
    CloseBoxIcon,
    CloseIcon,
    FileDocumentMinusIcon,
    FileDocumentMultipleIcon,
    FileDocumentPlusIcon,
    FileSearchIcon,
    LinkIcon,
    LockClockIcon,
    OpenInNewIcon,
} from 'mdi-vue3';
import { type DefineComponent } from 'vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import Time from '~/components/partials/elements/Time.vue';
import { ClipboardDocument, getDocument, getUser, useClipboardStore } from '~/store/clipboard';
import { DocReference, DocumentReference, DocumentShort } from '~~/gen/ts/resources/documents/documents';
import CitizenInfoPopover from '../partials/citizens/CitizenInfoPopover.vue';

const { $grpc } = useNuxtApp();
const clipboardStore = useClipboardStore();

const { t } = useI18n();

const props = defineProps<{
    open: boolean;
    document?: bigint;
    modelValue: Map<bigint, DocumentReference>;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'update:modelValue', payload: Map<bigint, DocumentReference>): void;
}>();

const tabs = ref<{ name: string; icon: DefineComponent }[]>([
    {
        name: t('components.documents.document_managers.view_current'),
        icon: markRaw(FileSearchIcon),
    },
    { name: t('common.clipboard'), icon: markRaw(ClipboardListIcon) },
    {
        name: t('components.documents.document_managers.add_new'),
        icon: markRaw(FileDocumentPlusIcon),
    },
]);

const queryDoc = ref('');

const {
    data: documents,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.document?.toString()}-references-docs-${queryDoc}`, () => listDocuments());

watchDebounced(queryDoc, async () => await refresh(), {
    debounce: 600,
    maxWait: 1750,
});

async function listDocuments(): Promise<DocumentShort[]> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getDocStoreClient().listDocuments({
                pagination: {
                    offset: 0n,
                    pageSize: 8n,
                },
                orderBy: [],
                search: queryDoc.value,
                categoryIds: [],
                creatorIds: [],
            });
            const { response } = await call;

            return res(
                response.documents.filter(
                    (doc) =>
                        !Array.from(props.modelValue.values()).find(
                            (r) => r.targetDocumentId === doc.id || doc.id === props.document,
                        ),
                ),
            );
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

function addReference(doc: DocumentShort, reference: DocReference): void {
    const keys = Array.from(props.modelValue.keys());
    const key = !keys.length ? 1n : keys[keys.length - 1] + 1n;

    props.modelValue.set(key, {
        id: key,
        sourceDocumentId: props.document ?? 0n,
        reference: reference,
        targetDocumentId: doc.id,
        targetDocument: doc,
    });
    refresh();
}

function addReferenceClipboard(doc: ClipboardDocument, reference: DocReference): void {
    addReference(getDocument(doc), reference);
}

function removeReference(id: bigint): void {
    props.modelValue.delete(id);
    listDocuments();
}
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="emit('close')">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 transition-opacity bg-opacity-75 bg-base-900" />
            </TransitionChild>

            <div class="fixed inset-0 z-30 overflow-y-auto">
                <div class="flex items-end justify-center min-h-full p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild
                        as="template"
                        enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100"
                        leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    >
                        <DialogPanel
                            class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-800 text-neutral sm:my-8 w-full sm:max-w-6xl sm:p-6 my-auto"
                        >
                            <div class="absolute top-0 right-0 hidden pt-4 pr-4 sm:block">
                                <button
                                    type="button"
                                    class="transition-colors rounded-md hover:text-base-300 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                    @click="emit('close')"
                                >
                                    <span class="sr-only">
                                        {{ $t('common.close', 1) }}
                                    </span>
                                    <CloseIcon class="w-6 h-6" aria-hidden="true" />
                                </button>
                            </div>
                            <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                {{ $t('common.document', 1) }}
                                {{ $t('common.reference', 2) }}
                            </DialogTitle>
                            <TabGroup>
                                <TabList class="flex flex-row mb-4">
                                    <Tab v-for="tab in tabs" :key="tab.name" v-slot="{ selected }" class="flex-initial w-full">
                                        <button
                                            :class="[
                                                selected
                                                    ? 'border-primary-500 text-primary-500'
                                                    : 'border-transparent text-base-300 hover:border-base-300 hover:text-base-200',
                                                'group inline-flex items-center border-b-2 py-4 px-1 text-m font-medium w-full justify-center transition-colors',
                                            ]"
                                            :aria-current="selected ? 'page' : undefined"
                                        >
                                            <component
                                                :is="tab.icon"
                                                :class="[
                                                    selected ? 'text-primary-500' : 'text-base-300 group-hover:text-base-200',
                                                    '-ml-0.5 mr-2 h-5 w-5 transition-colors',
                                                ]"
                                                aria-hidden="true"
                                            />
                                            <span>{{ tab.name }}</span>
                                        </button>
                                    </Tab>
                                </TabList>
                                <TabPanels>
                                    <div class="px-4 sm:flex sm:items-start sm:px-6 lg:px-8">
                                        <TabPanel class="w-full">
                                            <div class="flow-root">
                                                <div class="mx-0 -my-2 overflow-x-auto">
                                                    <div class="inline-block min-w-full py-2 align-middle">
                                                        <table class="min-w-full divide-y divide-base-200 text-neutral">
                                                            <thead>
                                                                <tr>
                                                                    <th
                                                                        scope="col"
                                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-6 lg:pl-8"
                                                                    >
                                                                        {{ $t('common.title') }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.state') }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.creator') }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.reference') }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.reference') }}
                                                                    </th>
                                                                </tr>
                                                            </thead>
                                                            <tbody class="divide-y divide-base-500">
                                                                <tr
                                                                    v-for="[key, reference] in $props.modelValue"
                                                                    :key="key.toString()"
                                                                >
                                                                    <td
                                                                        class="py-4 pl-4 pr-3 text-sm font-medium truncate whitespace-nowrap sm:pl-6 lg:pl-8 max-w-xl"
                                                                    >
                                                                        {{ reference.targetDocument?.title }}
                                                                    </td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        {{ reference.targetDocument?.state }}
                                                                    </td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        <CitizenInfoPopover
                                                                            :user="reference.targetDocument?.creator"
                                                                        />
                                                                    </td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        {{
                                                                            $t(
                                                                                `enums.docstore.DocReference.${
                                                                                    DocReference[reference.reference]
                                                                                }`,
                                                                            )
                                                                        }}
                                                                    </td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        <div class="flex flex-row gap-2">
                                                                            <div class="flex">
                                                                                <NuxtLink
                                                                                    :to="{
                                                                                        name: 'documents-id',
                                                                                        params: {
                                                                                            id: reference.targetDocumentId.toString(),
                                                                                        },
                                                                                    }"
                                                                                    target="_blank"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.open_document',
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <OpenInNewIcon
                                                                                        class="w-6 h-auto text-primary-500 hover:text-primary-300"
                                                                                    />
                                                                                </NuxtLink>
                                                                            </div>
                                                                            <div class="flex">
                                                                                <button
                                                                                    role="button"
                                                                                    @click="removeReference(reference.id!)"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.remove_reference',
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <FileDocumentMinusIcon
                                                                                        class="w-6 h-auto text-error-400 hover:text-error-200"
                                                                                    />
                                                                                </button>
                                                                            </div>
                                                                        </div>
                                                                    </td>
                                                                </tr>
                                                            </tbody>
                                                        </table>
                                                    </div>
                                                </div>
                                            </div>
                                        </TabPanel>
                                        <TabPanel class="w-full">
                                            <div class="flow-root mt-2">
                                                <div class="mx-0 -my-2 overflow-x-auto">
                                                    <div class="inline-block min-w-full py-2 align-middle">
                                                        <DataNoDataBlock
                                                            v-if="clipboardStore.$state.documents.length === 0"
                                                            :type="$t('common.reference', 2)"
                                                            :icon="FileDocumentMultipleIcon"
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
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.state') }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.creator') }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.created_at') }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{
                                                                            $t(
                                                                                'components.documents.document_managers.add_reference',
                                                                            )
                                                                        }}
                                                                    </th>
                                                                </tr>
                                                            </thead>
                                                            <tbody class="divide-y divide-base-500">
                                                                <tr
                                                                    v-for="document in clipboardStore.$state.documents"
                                                                    :key="document.id?.toString()"
                                                                >
                                                                    <td
                                                                        class="py-4 pl-4 pr-3 text-sm font-medium truncate whitespace-nowrap sm:pl-6 lg:pl-8 max-w-xl"
                                                                    >
                                                                        {{ document.title }}
                                                                    </td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        {{ document.state }}
                                                                    </td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        <CitizenInfoPopover :user="getUser(document.creator)" />
                                                                    </td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        {{ $t('common.created') }}
                                                                        <Time
                                                                            :value="
                                                                                new Date(Date.parse(document.createdAt ?? ''))
                                                                            "
                                                                        />
                                                                    </td>
                                                                    <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                        <div class="flex flex-row gap-2">
                                                                            <div class="flex">
                                                                                <button
                                                                                    role="button"
                                                                                    @click="
                                                                                        addReferenceClipboard(
                                                                                            document,
                                                                                            DocReference.LINKED,
                                                                                        )
                                                                                    "
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.links',
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <LinkIcon
                                                                                        class="w-6 h-auto text-info-500 hover:text-info-300"
                                                                                    />
                                                                                </button>
                                                                                <button
                                                                                    role="button"
                                                                                    @click="
                                                                                        addReferenceClipboard(
                                                                                            document,
                                                                                            DocReference.SOLVES,
                                                                                        )
                                                                                    "
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.solves',
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <CheckIcon
                                                                                        class="w-6 h-auto text-success-500 hover:text-success-300"
                                                                                    />
                                                                                </button>
                                                                                <button
                                                                                    role="button"
                                                                                    @click="
                                                                                        addReferenceClipboard(
                                                                                            document,
                                                                                            DocReference.CLOSES,
                                                                                        )
                                                                                    "
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.closes',
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <CloseBoxIcon
                                                                                        class="w-6 h-auto text-error-500 hover:text-error-300"
                                                                                    />
                                                                                </button>
                                                                                <button
                                                                                    role="button"
                                                                                    @click="
                                                                                        addReferenceClipboard(
                                                                                            document,
                                                                                            DocReference.DEPRECATES,
                                                                                        )
                                                                                    "
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.deprecates',
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <LockClockIcon
                                                                                        class="w-6 h-auto text-warn-500 hover:text-warn-300"
                                                                                    />
                                                                                </button>
                                                                            </div>
                                                                        </div>
                                                                    </td>
                                                                </tr>
                                                            </tbody>
                                                        </table>
                                                    </div>
                                                </div>
                                            </div>
                                        </TabPanel>
                                        <TabPanel class="w-full">
                                            <div>
                                                <label for="title" class="sr-only"
                                                    >{{ $t('common.document', 1) }} {{ $t('common.title') }}</label
                                                >
                                                <input
                                                    type="text"
                                                    name="title"
                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    :placeholder="`${$t('common.document', 1)} ${$t('common.title')}`"
                                                    v-model="queryDoc"
                                                />
                                            </div>
                                            <div class="flow-root mt-2">
                                                <div class="mx-0 -my-2 overflow-x-auto">
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
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.state') }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.creator') }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.created_at') }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.add') }}
                                                                        {{ $t('common.reference', 1) }}
                                                                    </th>
                                                                </tr>
                                                            </thead>
                                                            <tbody class="divide-y divide-base-500">
                                                                <template v-if="documents">
                                                                    <tr
                                                                        v-for="document in documents.slice(0, 8)"
                                                                        :key="document.id?.toString()"
                                                                    >
                                                                        <td
                                                                            class="py-4 pl-4 pr-3 text-sm font-medium truncate whitespace-nowrap sm:pl-6 lg:pl-8 max-w-xl"
                                                                        >
                                                                            {{ document.title }}
                                                                        </td>
                                                                        <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                            {{ document.state }}
                                                                        </td>
                                                                        <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                            <CitizenInfoPopover :user="document.creator" />
                                                                        </td>
                                                                        <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                            {{ $t('common.created') }}
                                                                            <Time :value="document.createdAt" :ago="true" />
                                                                        </td>
                                                                        <td class="px-3 py-4 text-sm whitespace-nowrap">
                                                                            <div class="flex flex-row gap-2">
                                                                                <div class="flex">
                                                                                    <button
                                                                                        role="button"
                                                                                        @click="
                                                                                            addReference(
                                                                                                document,
                                                                                                DocReference.LINKED,
                                                                                            )
                                                                                        "
                                                                                        data-te-toggle="tooltip"
                                                                                        :title="
                                                                                            $t(
                                                                                                'components.documents.document_managers.links',
                                                                                            )
                                                                                        "
                                                                                    >
                                                                                        <LinkIcon
                                                                                            class="w-6 h-auto text-info-500 hover:text-info-300"
                                                                                        />
                                                                                    </button>
                                                                                    <button
                                                                                        role="button"
                                                                                        @click="
                                                                                            addReference(
                                                                                                document,
                                                                                                DocReference.SOLVES,
                                                                                            )
                                                                                        "
                                                                                        data-te-toggle="tooltip"
                                                                                        :title="
                                                                                            $t(
                                                                                                'components.documents.document_managers.solves',
                                                                                            )
                                                                                        "
                                                                                    >
                                                                                        <CheckIcon
                                                                                            class="w-6 h-auto text-success-500 hover:text-success-300"
                                                                                        />
                                                                                    </button>
                                                                                    <button
                                                                                        role="button"
                                                                                        @click="
                                                                                            addReference(
                                                                                                document,
                                                                                                DocReference.CLOSES,
                                                                                            )
                                                                                        "
                                                                                        data-te-toggle="tooltip"
                                                                                        :title="
                                                                                            $t(
                                                                                                'components.documents.document_managers.closes',
                                                                                            )
                                                                                        "
                                                                                    >
                                                                                        <CloseBoxIcon
                                                                                            class="w-6 h-auto text-error-500 hover:text-error-300"
                                                                                        />
                                                                                    </button>
                                                                                    <button
                                                                                        role="button"
                                                                                        @click="
                                                                                            addReference(
                                                                                                document,
                                                                                                DocReference.DEPRECATES,
                                                                                            )
                                                                                        "
                                                                                        data-te-toggle="tooltip"
                                                                                        :title="
                                                                                            $t(
                                                                                                'components.documents.document_managers.deprecates',
                                                                                            )
                                                                                        "
                                                                                    >
                                                                                        <LockClockIcon
                                                                                            class="w-6 h-auto text-warn-500 hover:text-warn-300"
                                                                                        />
                                                                                    </button>
                                                                                </div>
                                                                            </div>
                                                                        </td>
                                                                    </tr>
                                                                </template>
                                                            </tbody>
                                                        </table>
                                                    </div>
                                                </div>
                                            </div>
                                        </TabPanel>
                                    </div>
                                </TabPanels>
                            </TabGroup>
                            <div class="gap-2 mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
                                <button
                                    type="button"
                                    class="rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="emit('close')"
                                >
                                    {{ $t('common.close', 1) }}
                                </button>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
