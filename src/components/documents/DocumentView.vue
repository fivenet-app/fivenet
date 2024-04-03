<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { useConfirmDialog, useTimeoutFn, watchOnce } from '@vueuse/core';
import { useRouteHash } from '@vueuse/router';
import {
    AccountIcon,
    AccountMultipleIcon,
    CalendarEditIcon,
    CalendarIcon,
    CalendarRemoveIcon,
    ChevronDownIcon,
    CommentIcon,
    CommentQuoteIcon,
    CommentTextMultipleIcon,
    CreationIcon,
    FileDocumentIcon,
    FileSearchIcon,
    FrequentlyAskedQuestionsIcon,
    LockIcon,
    LockOpenVariantIcon,
    NoteCheckIcon,
    PencilIcon,
    RestoreIcon,
    ShapeIcon,
    TrashCanIcon,
} from 'mdi-vue3';
import AddToButton from '~/components/clipboard/AddToButton.vue';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import { AccessLevel, DocumentAccess } from '~~/gen/ts/resources/documents/access';
import { type Document } from '~~/gen/ts/resources/documents/documents';
import DocumentComments from '~/components/documents/DocumentComments.vue';
import DocumentReferences from '~/components/documents/DocumentReferences.vue';
import DocumentRelations from '~/components/documents/DocumentRelations.vue';
import { checkDocAccess } from '~/components/documents/helpers';
import DocumentActivityList from '~/components/documents/DocumentActivityList.vue';
import DocumentRequestsModal from '~/components/documents/requests/DocumentRequestsModal.vue';
import { useAuthStore } from '~/store/auth';
import DocumentRequestAccess from '~/components/documents/requests/DocumentRequestAccess.vue';

const { $grpc } = useNuxtApp();

const clipboardStore = useClipboardStore();

const notifications = useNotificatorStore();

const authStore = useAuthStore();
const { activeChar, isSuperuser } = storeToRefs(authStore);

const props = defineProps<{
    documentId: string;
}>();

const access = ref<undefined | DocumentAccess>(undefined);
const commentCount = ref<number | undefined>();

const {
    data: doc,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}`, () => getDocument(props.documentId));

async function getDocument(id: string): Promise<Document> {
    try {
        const call = $grpc.getDocStoreClient().getDocument({
            documentId: id,
        });
        const { response } = await call;

        access.value = response.access;

        return response.document!;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function deleteDocument(id: string): Promise<void> {
    try {
        await $grpc.getDocStoreClient().deleteDocument({
            documentId: id,
        });

        notifications.add({
            title: { key: 'notifications.document_deleted.title', parameters: {} },
            description: { key: 'notifications.document_deleted.content', parameters: {} },
            type: 'success',
        });

        await navigateTo({ name: 'documents' });
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function toggleDocument(id: string, closed: boolean): Promise<void> {
    try {
        await $grpc.getDocStoreClient().toggleDocument({
            documentId: id,
            closed,
        });

        doc.value!.closed = closed;

        if (!closed) {
            notifications.add({
                title: { key: `notifications.document_toggled.open.title`, parameters: {} },
                description: { key: `notifications.document_toggled.open.content`, parameters: {} },
                type: 'success',
            });
        } else {
            notifications.add({
                title: { key: `notifications.document_toggled.closed.title`, parameters: {} },
                description: { key: `notifications.document_toggled.closed.content`, parameters: {} },
                type: 'success',
            });
        }
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function changeDocumentOwner(id: string): Promise<void> {
    try {
        await $grpc.getDocStoreClient().changeDocumentOwner({
            documentId: id,
        });

        notifications.add({
            title: { key: 'notifications.document_take_ownership.title', parameters: {} },
            description: { key: 'notifications.document_take_ownership.content', parameters: {} },
            type: 'success',
        });

        await refresh();
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function addToClipboard(): void {
    if (doc.value) {
        clipboardStore.addDocument(doc.value);
    }

    notifications.add({
        title: { key: 'notifications.clipboard.document_added.title', parameters: {} },
        description: { key: 'notifications.clipboard.document_added.content', parameters: {} },
        timeout: 3250,
        type: 'info',
    });
}

const contentRef = ref<HTMLDivElement | null>(null);
function disableCheckboxes(): void {
    if (contentRef.value === null) {
        return;
    }

    const checkboxes: NodeListOf<HTMLInputElement> = contentRef.value.querySelectorAll('input[type=checkbox]');
    checkboxes.forEach((el) => {
        el.setAttribute('disabled', 'disabled');
        el.onchange = (ev) => ev.preventDefault();
    });
}

watchOnce(doc, () => useTimeoutFn(disableCheckboxes, 50));

const {
    isRevealed: isRevealedChangeOwner,
    reveal: revealChangeOwner,
    confirm: confirmChangeOwner,
    cancel: cancelChangeOwner,
    onConfirm: onConfirmChangeOwner,
} = useConfirmDialog();
onConfirmChangeOwner(async (id: string) => changeDocumentOwner(id));

const {
    isRevealed: isRevealedDelete,
    reveal: revealDelete,
    confirm: confirmDelete,
    cancel: cancelDelete,
    onConfirm: onConfirmDelete,
} = useConfirmDialog();
onConfirmDelete(async (id: string) => deleteDocument(id));

const openRequests = ref(false);

const hash = useRouteHash();
if (hash.value !== undefined && hash.value !== null) {
    if (hash.value.replace(/^#/, '') === 'requests') {
        openRequests.value = true;
    }
}
</script>

<template>
    <div class="m-2">
        <ConfirmDialog
            :open="isRevealedChangeOwner"
            :cancel="cancelChangeOwner"
            :confirm="() => confirmChangeOwner(documentId)"
        />
        <ConfirmDialog :open="isRevealedDelete" :cancel="cancelDelete" :confirm="() => confirmDelete(documentId)" />

        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.document', 2)])" />
        <template v-else-if="error">
            <DataErrorBlock :title="$t('common.unable_to_load', [$t('common.document', 2)])" :retry="refresh" />
            <DocumentRequestAccess v-if="error.message.endsWith('ErrDocViewDenied')" :document-id="documentId" />
        </template>
        <DataNoDataBlock
            v-else-if="doc === null"
            icon="i-mdi-file-search"
            :message="$t('common.not_found', [$t('common.document', 2)])"
        />

        <div v-else class="rounded-lg bg-base-700">
            <DocumentRequestsModal
                v-if="can('DocStoreService.ListDocumentReqs') && access !== undefined"
                :open="openRequests"
                :access="access"
                :doc="doc"
                @close="openRequests = false"
                @refresh="
                    openRequests = false;
                    refresh();
                "
            />

            <div class="h-full px-4 py-6 sm:px-6 lg:px-8">
                <div>
                    <div>
                        <div class="flex snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                            <IDCopyBadge
                                :id="doc.id"
                                prefix="DOC"
                                :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: {} }"
                                :content="{ key: 'notifications.document_view.copy_document_id.content', parameters: {} }"
                            />

                            <div class="flex space-x-2 self-end">
                                <UButton
                                    v-if="
                                        can('DocStoreService.ToggleDocument') &&
                                        checkDocAccess(
                                            access,
                                            doc.creator,
                                            AccessLevel.STATUS,
                                            'DocStoreService.ToggleDocument',
                                        )
                                    "
                                    @click="toggleDocument(documentId, !doc?.closed)"
                                >
                                    <template v-if="doc?.closed">
                                        <LockOpenVariantIcon class="size-5 text-success-500" />
                                        {{ $t('common.open', 1) }}
                                    </template>
                                    <template v-else>
                                        <LockIcon class="size-5 text-error-400" />
                                        {{ $t('common.close', 1) }}
                                    </template>
                                </UButton>
                                <UButton
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
                                        name: 'documents-id-edit',
                                        params: { id: doc.id },
                                    }"
                                    icon="i-mdi-pencil"
                                >
                                    {{ $t('common.edit') }}
                                </UButton>
                                <UButton
                                    v-if="can('DocStoreService.ListDocumentReqs')"
                                    icon="i-mdi-frequently-asked-questions"
                                    @click="openRequests = true"
                                >
                                    {{ $t('common.request', 2) }}
                                </UButton>
                                <UButton
                                    v-if="
                                        (doc?.creatorJob === activeChar?.job || isSuperuser) &&
                                        can('DocStoreService.ChangeDocumentOwner') &&
                                        checkDocAccess(
                                            access,
                                            doc?.creator,
                                            AccessLevel.EDIT,
                                            'DocStoreService.ChangeDocumentOwner',
                                        )
                                    "
                                    :class="doc?.creatorId === activeChar?.userId ? 'disabled' : ''"
                                    :disabled="doc?.creatorId === activeChar?.userId"
                                    icon="i-mdi-creation"
                                    @click="revealChangeOwner(documentId)"
                                >
                                    {{ $t('components.documents.document_view.take_ownership') }}
                                </UButton>
                                <UButton
                                    v-if="
                                        can('DocStoreService.DeleteDocument') &&
                                        checkDocAccess(access, doc.creator, AccessLevel.EDIT, 'DocStoreService.DeleteDocument')
                                    "
                                    :icon="!doc.deletedAt ? 'i-mdi-trash-can' : 'i-mdi-restore'"
                                    @click="revealDelete(documentId)"
                                >
                                    <template v-if="!doc.deletedAt">
                                        {{ $t('common.delete') }}
                                    </template>
                                    <template v-else>
                                        {{ $t('common.restore') }}
                                    </template>
                                </UButton>
                            </div>
                        </div>

                        <div class="my-4">
                            <h1 class="break-words px-0.5 py-1 text-4xl font-bold sm:pl-1">
                                {{ doc?.title }}
                            </h1>
                        </div>

                        <div class="mb-2 flex gap-2">
                            <div
                                v-if="doc.category"
                                class="flex flex-initial flex-row gap-1 rounded-full bg-primary-100 px-2 py-1 text-primary-500"
                            >
                                <ShapeIcon class="h-auto w-5" />
                                <span
                                    class="inline-flex items-center text-sm font-medium text-primary-800"
                                    :title="doc.category.description ?? $t('common.na')"
                                >
                                    {{ doc.category.name }}
                                </span>
                            </div>

                            <div
                                v-if="doc?.closed"
                                class="flex flex-initial flex-row gap-1 rounded-full bg-error-100 px-2 py-1"
                            >
                                <LockIcon class="size-5 text-error-400" />
                                <span class="text-sm font-medium text-error-700">
                                    {{ $t('common.close', 2) }}
                                </span>
                            </div>
                            <div v-else class="flex flex-initial flex-row gap-1 rounded-full bg-success-100 px-2 py-1">
                                <LockOpenVariantIcon class="size-5 text-success-500" />
                                <span class="text-sm font-medium text-success-700">
                                    {{ $t('common.open', 2) }}
                                </span>
                            </div>

                            <div
                                v-if="doc?.state"
                                class="flex flex-initial flex-row gap-1 rounded-full bg-info-100 px-2 py-1 text-info-500"
                            >
                                <NoteCheckIcon class="h-auto w-5" />
                                <span class="text-sm font-medium text-info-800">
                                    {{ doc?.state }}
                                </span>
                            </div>
                            <div class="flex flex-initial flex-row gap-1 rounded-full bg-base-100 px-2 py-1 text-base-500">
                                <CommentTextMultipleIcon class="h-auto w-5" />
                                <span class="text-sm font-medium text-base-700">
                                    {{
                                        commentCount !== undefined
                                            ? $t('common.comments', commentCount)
                                            : '? ' + $t('common.comment', 2)
                                    }}
                                </span>
                            </div>
                        </div>

                        <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-0">
                            <div class="flex flex-initial flex-row gap-1 rounded-full bg-base-100 px-2 py-1 text-base-500">
                                <AccountIcon class="h-auto w-5" />
                                <span class="inline-flex items-center text-sm font-medium text-base-700">
                                    {{ $t('common.created_by') }}
                                    <CitizenInfoPopover
                                        :user="doc.creator"
                                        class="ml-1 font-medium text-primary-600 hover:text-primary-400"
                                    />
                                </span>
                            </div>

                            <div class="flex flex-initial flex-row gap-1 rounded-full bg-base-100 px-2 py-1 text-base-500">
                                <CalendarIcon class="h-auto w-5" />
                                <span class="text-sm font-medium text-base-700">
                                    {{ $t('common.created_at') }}
                                    <GenericTime :value="doc.createdAt" type="long" />
                                </span>
                            </div>
                            <div
                                v-if="doc.updatedAt"
                                class="flex flex-initial flex-row gap-1 rounded-full bg-base-100 px-2 py-1 text-base-500"
                            >
                                <CalendarEditIcon class="h-auto w-5" />
                                <span class="text-sm font-medium text-base-700">
                                    {{ $t('common.updated_at') }}
                                    <GenericTime :value="doc.updatedAt" type="long" />
                                </span>
                            </div>
                            <div
                                v-if="doc.deletedAt"
                                class="flex flex-initial flex-row gap-1 rounded-full bg-base-100 px-2 py-1 text-base-500"
                            >
                                <CalendarRemoveIcon class="h-auto w-5" />
                                <span class="text-sm font-medium text-base-700">
                                    {{ $t('common.deleted') }}
                                    <GenericTime :value="doc.deletedAt" type="long" />
                                </span>
                            </div>
                        </div>

                        <div class="my-2">
                            <h2 class="sr-only">
                                {{ $t('common.content') }}
                            </h2>
                            <div class="break-words rounded-lg bg-base-800">
                                <!-- eslint-disable vue/no-v-html -->
                                <div
                                    ref="contentRef"
                                    class="prose prose-invert min-w-full rounded-md bg-base-900 p-4"
                                    v-html="doc.content"
                                ></div>
                            </div>
                        </div>

                        <div class="my-2">
                            <Disclosure v-slot="{ open }" as="div" class="border-neutral/20 hover:border-neutral/70">
                                <DisclosureButton
                                    :class="[
                                        open ? 'rounded-t-lg border-b-0' : 'rounded-lg',
                                        'flex w-full items-start justify-between border-2 border-inherit p-2 text-left transition-colors',
                                    ]"
                                >
                                    <span class="inline-flex items-center text-base font-semibold leading-7">
                                        <AccountMultipleIcon class="mr-2 h-auto w-5" />
                                        {{ $t('common.relation', 2) }}
                                    </span>
                                    <span class="ml-6 flex h-7 items-center">
                                        <ChevronDownIcon
                                            :class="[open ? 'upsidedown' : '', 'h-auto w-5 transition-transform']"
                                        />
                                    </span>
                                </DisclosureButton>
                                <DisclosurePanel class="rounded-b-lg border-2 border-t-0 border-inherit transition-colors">
                                    <div class="mx-4 pb-2">
                                        <DocumentRelations :document-id="documentId" :show-document="false" />
                                    </div>
                                </DisclosurePanel>
                            </Disclosure>
                        </div>

                        <div class="my-2">
                            <Disclosure v-slot="{ open }" as="div" class="border-neutral/20 hover:border-neutral/70">
                                <DisclosureButton
                                    :class="[
                                        open ? 'rounded-t-lg border-b-0' : 'rounded-lg',
                                        'flex w-full items-start justify-between border-2 border-inherit p-2 text-left transition-colors',
                                    ]"
                                >
                                    <span class="inline-flex items-center text-base font-semibold leading-7">
                                        <FileDocumentIcon class="mr-2 h-auto w-5" />
                                        {{ $t('common.reference', 2) }}
                                    </span>
                                    <span class="ml-6 flex h-7 items-center">
                                        <ChevronDownIcon
                                            :class="[open ? 'upsidedown' : '', 'h-auto w-5 transition-transform']"
                                        />
                                    </span>
                                </DisclosureButton>
                                <DisclosurePanel class="rounded-b-lg border-2 border-t-0 border-inherit transition-colors">
                                    <div class="mx-4 pb-2">
                                        <DocumentReferences :document-id="documentId" :show-source="false" />
                                    </div>
                                </DisclosurePanel>
                            </Disclosure>
                        </div>

                        <div class="w-full">
                            <Disclosure
                                v-slot="{ open }"
                                as="div"
                                class="w-full border-neutral/20 hover:border-neutral/70"
                                :default-open="true"
                            >
                                <DisclosureButton
                                    :class="[
                                        open ? 'rounded-t-lg border-b-0' : 'rounded-lg',
                                        'flex w-full items-start justify-between border-2 border-inherit p-2 text-left transition-colors',
                                    ]"
                                >
                                    <span class="inline-flex items-center text-base font-semibold leading-7">
                                        <LockIcon class="mr-2 h-auto w-5" />
                                        {{ $t('common.access') }}
                                    </span>
                                    <span class="ml-6 flex h-7 items-center">
                                        <ChevronDownIcon
                                            :class="[open ? 'upsidedown' : '', 'h-auto w-5 transition-transform']"
                                        />
                                    </span>
                                </DisclosureButton>
                                <DisclosurePanel class="rounded-b-lg border-2 border-t-0 border-inherit transition-colors">
                                    <div class="mx-4 flex flex-row flex-wrap gap-1 pb-2">
                                        <DataNoDataBlock
                                            v-if="!access || (access?.jobs.length === 0 && access?.users.length === 0)"
                                            icon="i-mdi-file-search"
                                            :message="$t('common.not_found', [$t('common.access', 2)])"
                                        />

                                        <div
                                            v-for="entry in access?.jobs"
                                            :key="entry.id"
                                            class="flex flex-initial snap-x snap-start items-center gap-1 overflow-x-auto whitespace-nowrap rounded-full bg-info-100 px-2 py-1"
                                        >
                                            <span class="size-2 rounded-full bg-info-500" />
                                            <span class="text-sm font-medium text-info-800"
                                                >{{ entry.jobLabel
                                                }}<span
                                                    v-if="entry.minimumGrade > 0"
                                                    :title="`${entry.jobLabel} - ${$t('common.rank')} ${entry.minimumGrade}`"
                                                >
                                                    ({{ entry.jobGradeLabel }})</span
                                                >
                                                -
                                                {{ $t(`enums.docstore.AccessLevel.${AccessLevel[entry.access]}`) }}
                                            </span>
                                        </div>
                                    </div>
                                    <div class="mx-4 flex flex-row flex-wrap gap-1 pb-2">
                                        <div
                                            v-for="entry in access?.users"
                                            :key="entry.id"
                                            class="flex flex-initial snap-start flex-row items-center gap-1 whitespace-nowrap rounded-full bg-secondary-100 px-2 py-1"
                                        >
                                            <span class="size-2 rounded-full bg-secondary-400" />
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
                                </DisclosurePanel>
                            </Disclosure>

                            <div id="comments" class="my-2">
                                <div
                                    class="w-full rounded-lg border-2 border-neutral/20 transition-colors hover:border-neutral/70"
                                >
                                    <h2 class="inline-flex items-center p-2 text-left text-lg font-semibold transition-colors">
                                        <CommentIcon class="mr-2 h-auto w-5" />
                                        {{ $t('common.comment', 2) }}
                                    </h2>

                                    <div class="px-2 pb-2">
                                        <DocumentComments
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

                        <div
                            v-if="
                                can('DocStoreService.ListDocumentActivity') &&
                                checkDocAccess(access, doc.creator, AccessLevel.STATUS)
                            "
                            class="my-2"
                        >
                            <Disclosure v-slot="{ open }" as="div" class="border-neutral/20 hover:border-neutral/70">
                                <DisclosureButton
                                    :class="[
                                        open ? 'rounded-t-lg border-b-0' : 'rounded-lg',
                                        'flex w-full items-start justify-between border-2 border-inherit p-2 text-left transition-colors',
                                    ]"
                                >
                                    <span class="inline-flex items-center text-base font-semibold leading-7">
                                        <CommentQuoteIcon class="mr-2 h-auto w-5" />
                                        {{ $t('common.activity') }}
                                    </span>
                                    <span class="ml-6 flex h-7 items-center">
                                        <ChevronDownIcon :class="[open ? 'upsidedown' : '', 'size-5 transition-transform']" />
                                    </span>
                                </DisclosureButton>
                                <DisclosurePanel class="rounded-b-lg border-2 border-t-0 border-inherit transition-colors">
                                    <div class="mx-4 pb-2">
                                        <DocumentActivityList :document-id="documentId" />
                                    </div>
                                </DisclosurePanel>
                            </Disclosure>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <AddToButton :callback="addToClipboard" :title="$t('components.clipboard.clipboard_button.add')" />
</template>

<style scoped>
.prose {
    * {
        margin-top: 4px;
        margin-bottom: 4px;
    }
}
</style>
