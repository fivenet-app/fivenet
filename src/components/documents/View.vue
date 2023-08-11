<script lang="ts" setup>
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from '@headlessui/vue';

import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { QuillEditor } from '@vueup/vue-quill';
import {
    AccountMultipleIcon,
    CalendarIcon,
    CommentTextMultipleIcon,
    FileDocumentIcon,
    FileSearchIcon,
    LockIcon,
    LockOpenVariantIcon,
    PencilIcon,
    TrashCanIcon,
} from 'mdi-vue3';
import { DefineComponent } from 'vue';
import AddToButton from '~/components/clipboard/AddToButton.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Time from '~/components/partials/elements/Time.vue';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificationsStore } from '~/store/notifications';
import { ACCESS_LEVEL } from '~~/gen/ts/resources/documents/access';
import { Document, DocumentAccess } from '~~/gen/ts/resources/documents/documents';
import Comments from './Comments.vue';
import References from './References.vue';
import Relations from './Relations.vue';

const { $grpc } = useNuxtApp();
const clipboardStore = useClipboardStore();
const notifications = useNotificationsStore();

const { t } = useI18n();

const access = ref<undefined | DocumentAccess>(undefined);
const commentCount = ref(-1n);
const tabs = ref<{ name: string; icon: DefineComponent }[]>([
    { name: t('common.relation', 2), icon: markRaw(AccountMultipleIcon) },
    { name: t('common.reference', 2), icon: markRaw(FileDocumentIcon) },
]);

const props = defineProps<{
    documentId: bigint;
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
                title: { key: 'notifications.document_deleted.title', parameters: [] },
                content: { key: 'notifications.document_deleted.content', parameters: [] },
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

async function toggleDocument(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const closed = !document.value?.closed ?? true;
            await $grpc.getDocStoreClient().toggleDocument({
                documentId: props.documentId,
                closed: closed,
            });

            document.value!.closed = !document.value?.closed;

            notifications.dispatchNotification({
                title: { key: `notifications.document_toggled.${!closed ? 'open' : 'closed'}.title`, parameters: [] },
                content: { key: `notifications.document_toggled.${!closed ? 'open' : 'closed'}.content`, parameters: [] },
                type: 'success',
            });

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
        title: { key: 'notifications.clipboard.document_added.title', parameters: [] },
        content: { key: 'notifications.clipboard.document_added.content', parameters: [] },
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
        <DataNoDataBlock
            v-else-if="!document"
            :icon="FileSearchIcon"
            :message="$t('common.not_found', [$t('common.document', 2)])"
        />
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
                                <div v-if="can('DocStoreService.ToggleDocument')">
                                    <button
                                        v-if="document?.closed"
                                        type="button"
                                        @click="toggleDocument"
                                        class="inline-flex justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                    >
                                        <LockOpenVariantIcon class="w-5 h-5 text-green-500" aria-hidden="true" />
                                        {{ $t('common.open', 2) }}
                                    </button>
                                    <button
                                        v-else
                                        type="button"
                                        @click="toggleDocument"
                                        class="inline-flex justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                    >
                                        <LockIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                                        {{ $t('common.close', 1) }}
                                    </button>
                                </div>
                                <NuxtLink
                                    v-if="can('DocStoreService.UpdateDocument')"
                                    :to="{
                                        name: 'documents-edit-id',
                                        params: { id: document?.id.toString() ?? 0 },
                                    }"
                                    type="button"
                                    class="inline-flex justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                >
                                    <PencilIcon class="-ml-0.5 w-5 h-auto" aria-hidden="true" />
                                    {{ $t('common.edit') }}
                                </NuxtLink>
                                <button
                                    v-if="can('DocStoreService.DeleteDocument')"
                                    type="button"
                                    @click="deleteDocument"
                                    class="inline-flex justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                >
                                    <TrashCanIcon class="-ml-0.5 w-5 h-auto" aria-hidden="true" />
                                    {{ $t('common.delete') }}
                                </button>
                            </div>
                        </div>
                        <div class="flex flex-row gap-2">
                            <IDCopyBadge
                                :id="document.id"
                                prefix="DOC"
                                :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: [] }"
                                :content="{ key: 'notifications.document_view.copy_document_id.content', parameters: [] }"
                            />
                            <div class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-base-100 text-base-500">
                                <CalendarIcon class="w-5 h-auto" aria-hidden="true" />
                                <span class="text-sm font-medium text-base-700">
                                    <Time :value="document.createdAt" type="long" />
                                </span>
                            </div>
                            <div
                                v-if="document?.closed"
                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-error-100"
                            >
                                <LockIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                                <span class="text-sm font-medium text-error-700">
                                    {{ $t('common.close', 2) }}
                                </span>
                            </div>
                            <div v-else class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-success-100">
                                <LockOpenVariantIcon class="w-5 h-5 text-green-500" aria-hidden="true" />
                                <span class="text-sm font-medium text-green-700">
                                    {{ $t('common.open') }}
                                </span>
                            </div>
                            <div
                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-primary-100 text-primary-500"
                            >
                                <CommentTextMultipleIcon class="w-5 h-auto" aria-hidden="true" />
                                <span class="text-sm font-medium text-primary-700">
                                    {{ commentCount >= 0 ? commentCount : '?' }}
                                    {{ $t('common.comment', 2) }}
                                </span>
                            </div>
                        </div>
                        <div class="flex flex-row flex-wrap gap-2 pb-3 mt-2 overflow-x-auto snap-x sm:pb-0">
                            <div
                                v-for="entry in access?.jobs"
                                :key="entry.id?.toString()"
                                class="flex flex-row items-center flex-initial gap-1 px-2 py-1 rounded-full bg-info-100 whitespace-nowrap snap-start"
                            >
                                <span class="w-2 h-2 rounded-full bg-info-500" aria-hidden="true" />
                                <span class="text-sm font-medium text-info-800"
                                    >{{ entry.jobLabel
                                    }}<span :title="`$t('common.rank') {{ entry.minimumGrade }}`" v-if="entry.minimumGrade > 0">
                                        ({{ entry.jobGradeLabel }})</span
                                    >
                                    -
                                    {{ $t(`enums.docstore.ACCESS_LEVEL.${ACCESS_LEVEL[entry.access]}`) }}
                                </span>
                            </div>
                            <div
                                v-for="entry in access?.users"
                                :key="entry.id?.toString()"
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
                                                :class="[
                                                    selected ? 'text-primary-500' : 'text-base-300 group-hover:text-base-200',
                                                    '-ml-0.5 mr-2 h-5 w-5 transition-colors',
                                                ]"
                                                aria-hidden="true"
                                                :is="tab.icon"
                                            />
                                            <span>{{ tab.name }}</span>
                                        </button>
                                    </Tab>
                                </TabList>
                                <TabPanels>
                                    <TabPanel>
                                        <Relations :document-id="documentId" :show-document="false" />
                                    </TabPanel>
                                    <TabPanel>
                                        <References :document-id="documentId" :show-source="false" />
                                    </TabPanel>
                                </TabPanels>
                            </TabGroup>
                        </div>
                        <div class="mt-4">
                            <h2 class="text-lg font-semibold text-neutral">
                                {{ $t('common.comment', 2) }}
                            </h2>
                            <Comments :document-id="documentId" :closed="document?.closed" @counted="commentCount = $event" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <AddToButton :callback="addToClipboard" :title="$t('components.clipboard.clipboard_button.add')" />
</template>
