<script lang="ts" setup>
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { Quill, QuillEditor } from '@vueup/vue-quill';
import { useConfirmDialog } from '@vueuse/core';
import {
    AccountMultipleIcon,
    CalendarEditIcon,
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
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Time from '~/components/partials/elements/Time.vue';
import Divider from '~/composables/quill/divider/divider';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import { Document, DocumentAccess } from '~~/gen/ts/resources/documents/documents';
import Comments from './Comments.vue';
import References from './References.vue';
import Relations from './Relations.vue';
import { checkDocAccess } from './helpers';

const { $grpc } = useNuxtApp();
const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

const { t } = useI18n();

const access = ref<undefined | DocumentAccess>(undefined);
const commentCount = ref<bigint | undefined>();
const tabs = ref<{ name: string; icon: DefineComponent }[]>([
    { name: t('common.relation', 2), icon: markRaw(AccountMultipleIcon) },
    { name: t('common.reference', 2), icon: markRaw(FileDocumentIcon) },
]);

const props = defineProps<{
    documentId: bigint;
}>();

const {
    data: doc,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}`, () => getDocument(props.documentId));

async function getDocument(id: bigint): Promise<Document> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getDocStoreClient().getDocument({
                documentId: id,
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

async function deleteDocument(id: bigint): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            await $grpc.getDocStoreClient().deleteDocument({
                documentId: id,
            });

            notifications.dispatchNotification({
                title: { key: 'notifications.document_deleted.title', parameters: {} },
                content: { key: 'notifications.document_deleted.content', parameters: {} },
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

async function toggleDocument(id: bigint, closed: boolean): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            await $grpc.getDocStoreClient().toggleDocument({
                documentId: id,
                closed: closed,
            });

            doc.value!.closed = closed;

            if (!closed) {
                notifications.dispatchNotification({
                    title: { key: `notifications.document_toggled.open.title`, parameters: {} },
                    content: { key: `notifications.document_toggled.open.content`, parameters: {} },
                    type: 'success',
                });
            } else {
                notifications.dispatchNotification({
                    title: { key: `notifications.document_toggled.closed.title`, parameters: {} },
                    content: { key: `notifications.document_toggled.closed.content`, parameters: {} },
                    type: 'success',
                });
            }

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

function addToClipboard(): void {
    if (doc.value) {
        clipboardStore.addDocument(doc.value);
    }

    notifications.dispatchNotification({
        title: { key: 'notifications.clipboard.document_added.title', parameters: {} },
        content: { key: 'notifications.clipboard.document_added.content', parameters: {} },
        duration: 3250,
        type: 'info',
    });
}

Quill.register(Divider);

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();
onConfirm(async (id: bigint) => deleteDocument(id));
</script>

<style>
#editor .ql-toolbar {
    display: none;
}
#editor .ql-editor {
    min-height: 100% !important;
}
#editor .ql-editor hr {
    color: white;
}
</style>

<template>
    <div class="m-2">
        <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(documentId)" />

        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.document', 2)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.document', 2)])" :retry="refresh" />
        <DataNoDataBlock
            v-else-if="!doc"
            :icon="FileSearchIcon"
            :message="$t('common.not_found', [$t('common.document', 2)])"
        />
        <div v-else class="rounded-lg bg-base-700">
            <div class="h-full px-4 py-6 sm:px-6 lg:px-8">
                <div>
                    <div>
                        <div class="pb-2 md:flex md:items-center md:justify-between md:space-x-4">
                            <div>
                                <h1
                                    class="py-1 pl-0.5 pr-0.5 text-2xl font-bold text-neutral sm:pl-0 max-w-5xl flex items-center"
                                >
                                    <span
                                        v-if="doc.category"
                                        class="mr-2 inline-flex items-center rounded-md bg-primary-400/10 px-2 py-1 text-xs font-medium text-primary-400 ring-1 ring-inset ring-primary-400/30 break-keep"
                                    >
                                        {{ doc.category.name }}
                                    </span>
                                    <span class="break-all">
                                        {{ doc?.title }}
                                    </span>
                                </h1>
                                <p class="text-sm text-base-300 inline-flex">
                                    {{ $t('common.created_by') }}
                                    <CitizenInfoPopover
                                        :user="doc.creator"
                                        class="ml-1 font-medium text-primary-400 hover:text-primary-300"
                                    />
                                </p>
                            </div>
                            <div class="flex mt-1 space-x-3 md:mt-0">
                                <div
                                    v-if="
                                        can('DocStoreService.ToggleDocument') &&
                                        checkDocAccess(
                                            access,
                                            doc.creator,
                                            AccessLevel.STATUS,
                                            'DocStoreService.ToggleDocument',
                                        )
                                    "
                                >
                                    <button
                                        type="button"
                                        @click="toggleDocument(documentId, !doc?.closed)"
                                        class="inline-flex items-center justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                    >
                                        <template v-if="doc?.closed">
                                            <LockOpenVariantIcon class="w-5 h-5 text-success-500" aria-hidden="true" />
                                            {{ $t('common.open', 1) }}
                                        </template>
                                        <template v-else>
                                            <LockIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                                            {{ $t('common.close', 1) }}
                                        </template>
                                    </button>
                                </div>
                                <NuxtLink
                                    v-if="
                                        can('DocStoreService.UpdateDocument') &&
                                        checkDocAccess(
                                            access,
                                            doc.creator,
                                            AccessLevel.ACCESS,
                                            'DocStoreService.UpdateDocument',
                                        )
                                    "
                                    :to="{
                                        name: 'documents-edit-id',
                                        params: { id: doc?.id.toString() ?? 0 },
                                    }"
                                    type="button"
                                    class="inline-flex justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                >
                                    <PencilIcon class="-ml-0.5 w-5 h-auto" aria-hidden="true" />
                                    {{ $t('common.edit') }}
                                </NuxtLink>
                                <button
                                    v-if="
                                        can('DocStoreService.DeleteDocument') &&
                                        checkDocAccess(access, doc.creator, AccessLevel.EDIT, 'DocStoreService.DeleteDocument')
                                    "
                                    type="button"
                                    @click="reveal(documentId)"
                                    class="inline-flex justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                >
                                    <TrashCanIcon class="-ml-0.5 w-5 h-auto" aria-hidden="true" />
                                    {{ $t('common.delete') }}
                                </button>
                            </div>
                        </div>
                        <div class="flex flex-row flex-wrap gap-2 pb-3 overflow-x-auto snap-x sm:pb-0">
                            <IDCopyBadge
                                :id="doc.id"
                                prefix="DOC"
                                :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: {} }"
                                :content="{ key: 'notifications.document_view.copy_document_id.content', parameters: {} }"
                            />
                            <div class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-base-100 text-base-500">
                                <CalendarIcon class="w-5 h-auto" aria-hidden="true" />
                                <span class="text-sm font-medium text-base-700">
                                    {{ $t('common.created_at') }}
                                    <Time :value="doc.createdAt" type="long" />
                                </span>
                            </div>
                            <div
                                v-if="doc.updatedAt"
                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-base-100 text-base-500"
                            >
                                <CalendarEditIcon class="w-5 h-auto" aria-hidden="true" />
                                <span class="text-sm font-medium text-base-700">
                                    {{ $t('common.updated_at') }}
                                    <Time :value="doc.updatedAt" type="long" />
                                </span>
                            </div>
                            <div
                                v-if="doc?.closed"
                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-error-100"
                            >
                                <LockIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                                <span class="text-sm font-medium text-error-700">
                                    {{ $t('common.close', 2) }}
                                </span>
                            </div>
                            <div v-else class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-success-100">
                                <LockOpenVariantIcon class="w-5 h-5 text-success-500" aria-hidden="true" />
                                <span class="text-sm font-medium text-success-700">
                                    {{ $t('common.open', 2) }}
                                </span>
                            </div>
                            <div
                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-primary-100 text-primary-500"
                            >
                                <CommentTextMultipleIcon class="w-5 h-auto" aria-hidden="true" />
                                <span class="text-sm font-medium text-primary-700">
                                    {{
                                        commentCount !== undefined
                                            ? $t('common.comments', parseInt(commentCount?.toString()))
                                            : '? ' + $t('common.comment', 2)
                                    }}
                                </span>
                            </div>
                        </div>
                        <div
                            v-if="access && (access?.jobs.length > 0 || access?.users.length > 0)"
                            class="flex flex-row flex-wrap gap-2 pb-3 mt-2 overflow-x-auto snap-x sm:pb-0"
                        >
                            <div
                                v-for="entry in access?.jobs"
                                :key="entry.id?.toString()"
                                class="flex flex-row items-center flex-initial gap-1 px-2 py-1 rounded-full bg-info-100 whitespace-nowrap snap-start"
                            >
                                <span class="w-2 h-2 rounded-full bg-info-500" aria-hidden="true" />
                                <span class="text-sm font-medium text-info-800"
                                    >{{ entry.jobLabel
                                    }}<span
                                        :title="`${entry.jobLabel} - ${$t('common.rank')} ${entry.minimumGrade}`"
                                        v-if="entry.minimumGrade > 0"
                                    >
                                        ({{ entry.jobGradeLabel }})</span
                                    >
                                    -
                                    {{ $t(`enums.docstore.AccessLevel.${AccessLevel[entry.access]}`) }}
                                </span>
                            </div>
                            <div
                                v-for="entry in access?.users"
                                :key="entry.id?.toString()"
                                class="flex flex-row items-center flex-initial gap-1 px-2 py-1 rounded-full bg-secondary-100 whitespace-nowrap snap-start"
                            >
                                <span class="w-2 h-2 rounded-full bg-secondary-400" aria-hidden="true" />
                                <span
                                    class="text-sm font-medium text-secondary-700"
                                    :title="`${$t('common.id')} ${entry.userId}`"
                                >
                                    {{ entry.user?.firstname }}
                                    {{ entry.user?.lastname }} -
                                    {{ $t(`enums.docstore.AccessLevel.${AccessLevel[entry.access]}`) }}
                                </span>
                            </div>
                        </div>
                        <div>
                            <h2 class="sr-only">
                                {{ $t('common.content') }}
                            </h2>
                            <div class="mt-4 mb-2 rounded-lg text-neutral bg-base-800 break-words">
                                <div id="editor">
                                    <QuillEditor
                                        class="bg-base-900 rounded-md"
                                        content-type="html"
                                        :content="doc?.content"
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
                                                    : 'border-base-600 text-base-300 hover:border-base-300 hover:text-base-200',
                                                'border rounded-t-md group inline-flex items-center border-b-2 py-4 px-1 text-m font-medium w-full justify-center transition-colors',
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
                            <Comments
                                :document-id="documentId"
                                :closed="doc?.closed"
                                :can-comment="checkDocAccess(access, doc.creator, AccessLevel.COMMENT)"
                                @counted="commentCount = $event"
                                @new-comment="commentCount && commentCount++"
                                @deleted-comment="commentCount && commentCount > 0 && commentCount--"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <AddToButton :callback="addToClipboard" :title="$t('components.clipboard.clipboard_button.add')" />
</template>
