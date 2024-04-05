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
import { CheckIcon, CloseBoxIcon, CloseIcon, FileDocumentMinusIcon, LinkIcon, LockClockIcon, OpenInNewIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { ClipboardDocument, getDocument, getUser, useClipboardStore } from '~/store/clipboard';
import { DocReference, DocumentReference, DocumentShort } from '~~/gen/ts/resources/documents/documents';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';

const { $grpc } = useNuxtApp();
const clipboardStore = useClipboardStore();

const { t } = useI18n();

const props = defineProps<{
    open: boolean;
    documentId?: string;
    modelValue: Map<string, DocumentReference>;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'update:modelValue', payload: Map<string, DocumentReference>): void;
}>();

const tabs = ref<{ name: string; icon: string }[]>([
    {
        name: t('components.documents.document_managers.view_current'),
        icon: 'i-mdi-file-search',
    },
    { name: t('common.clipboard'), icon: 'i-mdi-clipboard-list' },
    {
        name: t('components.documents.document_managers.add_new'),
        icon: 'i-mdi-file-document-plus',
    },
]);

const queryDoc = ref('');

const {
    data: documents,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-references-docs-${queryDoc}`, () => listDocuments());

watchDebounced(queryDoc, async () => await refresh(), {
    debounce: 600,
    maxWait: 1750,
});

async function listDocuments(): Promise<DocumentShort[]> {
    try {
        const call = $grpc.getDocStoreClient().listDocuments({
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
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function addReference(doc: DocumentShort, reference: DocReference): void {
    const keys = Array.from(props.modelValue.keys());
    const key = !keys.length ? '1' : (parseInt(keys[keys.length - 1]) + 1).toString();

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
                <div class="fixed inset-0 bg-base-900/75 transition-opacity" />
            </TransitionChild>

            <div class="fixed inset-0 z-30 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
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
                            class="relative my-auto w-full overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left transition-all sm:my-8 sm:max-w-6xl sm:p-6"
                        >
                            <div class="absolute right-0 top-0 hidden pr-4 pt-4 sm:block">
                                <UButton
                                    class="focus:ring-primary-500 rounded-md transition-colors hover:text-base-300 focus:ring-2 focus:ring-offset-2"
                                    @click="emit('close')"
                                >
                                    <span class="sr-only">
                                        {{ $t('common.close', 1) }}
                                    </span>
                                    <CloseIcon class="size-5" />
                                </UButton>
                            </div>
                            <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                {{ $t('common.document', 1) }}
                                {{ $t('common.reference', 2) }}
                            </DialogTitle>
                            <TabGroup>
                                <TabList class="mb-4 flex flex-row">
                                    <Tab v-for="tab in tabs" :key="tab.name" v-slot="{ selected }" class="w-full flex-initial">
                                        <UButton
                                            :class="[
                                                selected
                                                    ? 'border-primary-500 text-primary-500'
                                                    : 'hover:text-accent-200 border-transparent text-base-300 hover:border-base-300',
                                                'group inline-flex w-full items-center justify-center border-b-2 px-1 py-4 text-sm font-medium transition-colors',
                                            ]"
                                            :aria-current="selected ? 'page' : undefined"
                                        >
                                            <component
                                                :is="tab.icon"
                                                :class="[
                                                    selected ? 'text-primary-500' : 'group-hover:text-accent-200 text-base-300',
                                                    '-ml-0.5 mr-2 size-5 transition-colors',
                                                ]"
                                            />
                                            <span>{{ tab.name }}</span>
                                        </UButton>
                                    </Tab>
                                </TabList>
                                <TabPanels>
                                    <div class="px-4 sm:flex sm:items-start sm:px-6 lg:px-8">
                                        <TabPanel class="w-full">
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
                                                                    v-for="[key, reference] in modelValue"
                                                                    :key="key.toString()"
                                                                >
                                                                    <td
                                                                        class="max-w-xl truncate whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-6 lg:pl-8"
                                                                    >
                                                                        {{ reference.targetDocument?.title }}
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        {{ reference.targetDocument?.state }}
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        <CitizenInfoPopover
                                                                            :user="reference.targetDocument?.creator"
                                                                        />
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        {{
                                                                            $t(
                                                                                `enums.docstore.DocReference.${
                                                                                    DocReference[reference.reference]
                                                                                }`,
                                                                            )
                                                                        }}
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        <div class="flex flex-row gap-2">
                                                                            <div class="flex">
                                                                                <NuxtLink
                                                                                    :to="{
                                                                                        name: 'documents-id',
                                                                                        params: {
                                                                                            id: reference.targetDocumentId,
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
                                                                                        class="text-primary-500 hover:text-primary-300 h-auto w-5"
                                                                                    />
                                                                                </NuxtLink>
                                                                            </div>
                                                                            <div class="flex">
                                                                                <UButton
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.remove_reference',
                                                                                        )
                                                                                    "
                                                                                    @click="removeReference(reference.id!)"
                                                                                >
                                                                                    <FileDocumentMinusIcon
                                                                                        class="h-auto w-5 text-error-400 hover:text-error-200"
                                                                                    />
                                                                                </UButton>
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
                                                                    :key="document.id"
                                                                >
                                                                    <td
                                                                        class="max-w-xl truncate whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-6 lg:pl-8"
                                                                    >
                                                                        {{ document.title }}
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        {{ document.state }}
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        <CitizenInfoPopover :user="getUser(document.creator)" />
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        {{ $t('common.created') }}
                                                                        <GenericTime
                                                                            :value="
                                                                                new Date(Date.parse(document.createdAt ?? ''))
                                                                            "
                                                                        />
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        <div class="flex flex-row gap-2">
                                                                            <div class="flex">
                                                                                <UButton
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.links',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addReferenceClipboard(
                                                                                            document,
                                                                                            DocReference.LINKED,
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <LinkIcon
                                                                                        class="h-auto w-5 text-info-500 hover:text-info-300"
                                                                                    />
                                                                                </UButton>
                                                                                <UButton
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.solves',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addReferenceClipboard(
                                                                                            document,
                                                                                            DocReference.SOLVES,
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <CheckIcon
                                                                                        class="h-auto w-5 text-success-500 hover:text-success-300"
                                                                                    />
                                                                                </UButton>
                                                                                <UButton
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.closes',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addReferenceClipboard(
                                                                                            document,
                                                                                            DocReference.CLOSES,
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <CloseBoxIcon
                                                                                        class="h-auto w-5 text-error-500 hover:text-error-300"
                                                                                    />
                                                                                </UButton>
                                                                                <UButton
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.deprecates',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addReferenceClipboard(
                                                                                            document,
                                                                                            DocReference.DEPRECATES,
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <LockClockIcon
                                                                                        class="h-auto w-5 text-warn-500 hover:text-warn-300"
                                                                                    />
                                                                                </UButton>
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
                                                <UInput
                                                    v-model="queryDoc"
                                                    type="text"
                                                    name="title"
                                                    class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    :placeholder="`${$t('common.document', 1)} ${$t('common.title')}`"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                />
                                            </div>
                                            <div class="mt-2 flow-root">
                                                <div class="-my-2 mx-0 overflow-x-auto">
                                                    <div class="inline-block min-w-full py-2 align-middle">
                                                        <DataPendingBlock
                                                            v-if="pending"
                                                            :message="$t('common.loading', [$t('common.document', 2)])"
                                                        />
                                                        <DataErrorBlock
                                                            v-else-if="error"
                                                            :title="$t('common.unable_to_load', [$t('common.document', 2)])"
                                                            :retry="refresh"
                                                        />
                                                        <DataNoDataBlock
                                                            v-else-if="documents === null || documents.length === 0"
                                                            :message="$t('components.citizens.citizens_list.no_citizens')"
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
                                                                        {{ $t('common.add') }}
                                                                        {{ $t('common.reference', 1) }}
                                                                    </th>
                                                                </tr>
                                                            </thead>
                                                            <tbody class="divide-y divide-base-500">
                                                                <tr
                                                                    v-for="document in documents.slice(0, 8)"
                                                                    :key="document.id"
                                                                >
                                                                    <td
                                                                        class="max-w-xl truncate whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-6 lg:pl-8"
                                                                    >
                                                                        {{ document.title }}
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        {{ document.state }}
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        <CitizenInfoPopover :user="document.creator" />
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        {{ $t('common.created') }}
                                                                        <GenericTime :value="document.createdAt" :ago="true" />
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        <div class="flex flex-row gap-2">
                                                                            <div class="flex">
                                                                                <UButton
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.links',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addReference(
                                                                                            document,
                                                                                            DocReference.LINKED,
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <LinkIcon
                                                                                        class="h-auto w-5 text-info-500 hover:text-info-300"
                                                                                    />
                                                                                </UButton>
                                                                                <UButton
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.solves',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addReference(
                                                                                            document,
                                                                                            DocReference.SOLVES,
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <CheckIcon
                                                                                        class="h-auto w-5 text-success-500 hover:text-success-300"
                                                                                    />
                                                                                </UButton>
                                                                                <UButton
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.closes',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addReference(
                                                                                            document,
                                                                                            DocReference.CLOSES,
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <CloseBoxIcon
                                                                                        class="h-auto w-5 text-error-500 hover:text-error-300"
                                                                                    />
                                                                                </UButton>
                                                                                <UButton
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.deprecates',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addReference(
                                                                                            document,
                                                                                            DocReference.DEPRECATES,
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <LockClockIcon
                                                                                        class="h-auto w-5 text-warn-500 hover:text-warn-300"
                                                                                    />
                                                                                </UButton>
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
                                    </div>
                                </TabPanels>
                            </TabGroup>
                            <div class="mt-5 gap-2 sm:mt-4 sm:flex sm:flex-row-reverse">
                                <UButton color="black" block class="flex-1" @click="emit('close')">
                                    {{ $t('common.close', 1) }}
                                </UButton>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
