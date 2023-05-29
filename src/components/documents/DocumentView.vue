<script lang="ts" setup>
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from '@headlessui/vue';
import {
    CalendarIcon,
    ChatBubbleLeftEllipsisIcon,
    DocumentMagnifyingGlassIcon,
    FingerPrintIcon,
    LockClosedIcon,
    LockOpenIcon,
    MagnifyingGlassIcon,
    PencilIcon,
    TrashIcon,
    UserIcon,
} from '@heroicons/vue/20/solid';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { QuillEditor } from '@vueup/vue-quill';
import { useClipboard } from '@vueuse/core';
import { ref } from 'vue';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificationsStore } from '~/store/notifications';
import { ACCESS_LEVEL } from '~~/gen/ts/resources/documents/access';
import { Document, DocumentAccess } from '~~/gen/ts/resources/documents/documents';
import AddToClipboardButton from '../clipboard/AddToClipboardButton.vue';
import DataErrorBlock from '../partials/DataErrorBlock.vue';
import DataPendingBlock from '../partials/DataPendingBlock.vue';
import DocumentComments from './DocumentComments.vue';
import DocumentReferences from './DocumentReferences.vue';
import DocumentRelations from './DocumentRelations.vue';

const { $grpc } = useNuxtApp();
const clipboardStore = useClipboardStore();
const notifications = useNotificationsStore();
const clipboard = useClipboard();

const { t } = useI18n();

const access = ref<undefined | DocumentAccess>(undefined);
const commentCount = ref(-1);
const tabs = ref<{ name: string; icon: typeof LockOpenIcon }[]>([
    { name: t('common.relation', 2), icon: UserIcon },
    { name: t('common.reference', 2), icon: DocumentMagnifyingGlassIcon },
]);

const props = defineProps<{
    documentId: number;
}>();

const { data: document, pending, refresh, error } = useLazyAsyncData(`document-${props.documentId}`, () => getDocument());

async function getDocument(): Promise<Document> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getDocStoreClient().getDocument({
                documentId: props.documentId,
            });
            const { response } = await call;

            access.value = response.access;

            return res(response.document!);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function deleteDocument(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            await $grpc.getDocStoreClient().deleteDocument({
                documentId: props.documentId,
            });

            notifications.dispatchNotification({
                title: t('notifications.document_deleted.title'),
                content: t('notifications.document_deleted.content'),
                type: 'success',
            });

            await navigateTo({ name: 'documents' });

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

function addToClipboard(): void {
    if (document.value) {
        clipboardStore.addDocument(document.value);
    }

    notifications.dispatchNotification({
        title: t('notifications.clipboard.document_added.title'),
        content: t('notifications.clipboard.document_added.content'),
        duration: 3500,
        type: 'info',
    });
}

function copyDocumentIDToClipboard(): void {
    clipboard.copy('DOC-' + props.documentId);
    notifications.dispatchNotification({
        title: t('notifications.document_view.copy_document_id.title'),
        content: t('notifications.document_view.copy_document_id.content'),
        duration: 3500,
        type: 'info',
    });
}
</script>

<style>
#editor .ql-toolbar {
    display: none;
}
</style>

<template>
    <div class="mt-2">
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.document', 2)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.document', 2)])" :retry="refresh" />
        <button
            v-else-if="!document"
            type="button"
            class="relative block w-full p-12 text-center rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
        >
            <MagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-base-200">
                {{ $t('common.not_found', [$t('common.document', 2)]) }}
            </span>
        </button>
        <div v-else class="rounded-lg bg-base-850">
            <div class="h-full px-4 py-6 sm:px-6 lg:px-8">
                <div>
                    <div>
                        <div class="pb-2 md:flex md:items-center md:justify-between md:space-x-4">
                            <div>
                                <h1 class="text-2xl font-bold text-neutral break-all">
                                    {{ document?.title }}
                                </h1>
                                <p class="text-sm text-base-300">
                                    {{ $t('common.created_by') }}
                                    {{ ' ' }}
                                    <NuxtLink
                                        :to="{
                                            name: 'citizens-id',
                                            params: {
                                                id: document?.creator?.userId ?? 0,
                                            },
                                        }"
                                        class="font-medium text-primary-400 hover:text-primary-300"
                                    >
                                        {{ document?.creator?.firstname }}
                                        {{ document?.creator?.lastname }}
                                    </NuxtLink>
                                </p>
                            </div>
                            <div class="flex mt-4 space-x-3 md:mt-0">
                                <NuxtLink
                                    :to="{
                                        name: 'documents-edit-id',
                                        params: { id: document?.id ?? 0 },
                                    }"
                                    type="button"
                                    class="inline-flex justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                    v-can="'DocStoreService.UpdateDocument'"
                                >
                                    <PencilIcon class="-ml-0.5 w-5 h-auto" aria-hidden="true" />
                                    {{ $t('common.edit') }}
                                </NuxtLink>
                                <button
                                    v-can="['DocStoreService.DeleteDocument']"
                                    type="button"
                                    @click="deleteDocument"
                                    class="inline-flex justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                >
                                    <TrashIcon class="-ml-0.5 w-5 h-auto" aria-hidden="true" />
                                    {{ $t('common.delete') }}
                                </button>
                            </div>
                        </div>
                        <div class="flex flex-row gap-2">
                            <div
                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full text-base-100 bg-base-500"
                                @click="copyDocumentIDToClipboard"
                            >
                                <FingerPrintIcon class="w-5 h-auto" aria-hidden="true" />
                                <span class="text-sm font-medium text-base-100">DOC-{{ documentId }}</span>
                            </div>
                            <div class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-base-100 text-base-500">
                                <CalendarIcon class="w-5 h-auto" aria-hidden="true" />
                                <span class="text-sm font-medium text-base-700"
                                    ><time :datetime="$d(toDate(document?.createdAt)!, 'short')">
                                        {{ $d(toDate(document?.createdAt)!, 'short') }}
                                    </time>
                                </span>
                            </div>
                            <div
                                v-if="document?.closed"
                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-error-100"
                            >
                                <LockClosedIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                                <span class="text-sm font-medium text-error-700">
                                    {{ $t('common.close', 2) }}
                                </span>
                            </div>
                            <div v-else class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-success-100">
                                <LockOpenIcon class="w-5 h-5 text-green-500" aria-hidden="true" />
                                <span class="text-sm font-medium text-green-700">
                                    {{ $t('common.open') }}
                                </span>
                            </div>
                            <div
                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-primary-100 text-primary-500"
                            >
                                <ChatBubbleLeftEllipsisIcon class="w-5 h-auto" aria-hidden="true" />
                                <span class="text-sm font-medium text-primary-700">
                                    {{ commentCount >= 0 ? commentCount : '?' }}
                                    {{ $t('common.comment', 2) }}
                                </span>
                            </div>
                        </div>
                        <div class="flex flex-row gap-2 pb-3 mt-2 overflow-x-auto snap-x sm:pb-0">
                            <div
                                v-for="entry in access?.jobs"
                                :key="entry.id"
                                class="flex flex-row items-center flex-initial gap-1 px-2 py-1 rounded-full bg-info-100 whitespace-nowrap snap-start"
                            >
                                <span class="w-2 h-2 rounded-full bg-info-500" aria-hidden="true" />
                                <span class="text-sm font-medium text-info-800"
                                    >{{ entry.jobLabel
                                    }}<span :title="entry.jobGradeLabel" v-if="entry.minimumGrade > 0">
                                        ({{ $t('common.rank') }}: {{ entry.minimumGrade }})</span
                                    >
                                    -
                                    {{ $t(`enums.docstore.ACCESS_LEVEL.${ACCESS_LEVEL[entry.access]}`) }}
                                </span>
                            </div>
                            <div
                                v-for="entry in access?.users"
                                :key="entry.id"
                                class="flex flex-row items-center flex-initial gap-1 px-2 py-1 rounded-full bg-secondary-100 whitespace-nowrap snap-start"
                            >
                                <span class="w-2 h-2 rounded-full bg-secondary-400" aria-hidden="true" />
                                <span class="text-sm font-medium text-secondary-700">
                                    {{ entry.user?.firstname }}
                                    {{ entry.user?.lastname }} -
                                    {{ $t(`enums.docstore.ACCESS_LEVEL.${ACCESS_LEVEL[entry.access]}`) }}
                                </span>
                            </div>
                        </div>
                        <div>
                            <h2 class="sr-only">
                                {{ $t('common.content') }}
                            </h2>
                            <div class="p-2 mt-4 rounded-lg text-neutral bg-base-800 break-words">
                                <div id="editor">
                                    <QuillEditor
                                        content-type="html"
                                        :content="document?.content"
                                        :toolbar="[]"
                                        theme="snow"
                                        :read-only="true"
                                    />
                                </div>
                            </div>
                        </div>
                        <div>
                            <TabGroup>
                                <TabList class="flex flex-row">
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
                                    <TabPanel>
                                        <DocumentRelations :document-id="documentId" :show-document="false" />
                                    </TabPanel>
                                    <TabPanel>
                                        <DocumentReferences :document-id="documentId" :show-source="false" />
                                    </TabPanel>
                                </TabPanels>
                            </TabGroup>
                        </div>
                        <div class="mt-4" v-can="'DocStoreService.GetDocumentComments'">
                            <h2 class="text-lg font-semibold text-neutral">
                                {{ $t('common.comment', 2) }}
                            </h2>
                            <DocumentComments
                                :document-id="documentId"
                                :closed="document?.closed"
                                @counted="commentCount = $event"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <AddToClipboardButton :callback="addToClipboard" :title="$t('components.clipboard.clipboard_button.add')" />
</template>
