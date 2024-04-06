<script lang="ts" setup>
import { useRouteHash } from '@vueuse/router';
import {
    AccountIcon,
    CalendarEditIcon,
    CalendarIcon,
    CalendarRemoveIcon,
    CommentTextMultipleIcon,
    LockIcon,
    LockOpenVariantIcon,
    NoteCheckIcon,
    ShapeIcon,
} from 'mdi-vue3';
import AddToButton from '~/components/clipboard/AddToButton.vue';
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
import ConfirmModal from '../partials/ConfirmModal.vue';

const props = defineProps<{
    documentId: string;
}>();

const { t } = useI18n();

const { $grpc } = useNuxtApp();

const clipboardStore = useClipboardStore();

const notifications = useNotificatorStore();

const authStore = useAuthStore();
const { activeChar, isSuperuser } = storeToRefs(authStore);

const modal = useModal();

const access = ref<undefined | DocumentAccess>(undefined);
const commentCount = ref<undefined | number>();

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

const hash = useRouteHash();
if (hash.value !== undefined && hash.value !== null) {
    if (hash.value.replace(/^#/, '') === 'requests') {
        openRequestsModal();
    }
}

function openRequestsModal(): void {
    if (access.value === undefined || doc.value === undefined) {
        return;
    }

    modal.open(DocumentRequestsModal, {
        access: access.value,
        doc: doc.value!,
        onRefresh: refresh,
    });
}

const accordionItems = [
    { slot: 'relations', label: t('common.relation', 2), icon: 'i-mdi-account-multiple' },
    { slot: 'references', label: t('common.reference', 2), icon: 'i-mdi-file-document' },
    { slot: 'access', label: t('common.access'), icon: 'i-mdi-lock', defaultOpen: true },
    { slot: 'comments', label: t('common.comment', 2), icon: 'i-mdi-comment', defaultOpen: true },
    { slot: 'activity', label: t('common.activity'), icon: 'i-mdi-comment-quote' },
];
</script>

<template>
    <div>
        <UDashboardNavbar :title="$t('pages.documents.id.title')">
            <template #right>
                <UButtonGroup>
                    <IDCopyBadge
                        :id="doc?.id ?? documentId"
                        prefix="DOC"
                        :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: {} }"
                        :content="{ key: 'notifications.document_view.copy_document_id.content', parameters: {} }"
                    />

                    <AddToButton :title="$t('components.clipboard.clipboard_button.add')" :callback="addToClipboard" />
                </UButtonGroup>
            </template>
        </UDashboardNavbar>

        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.document', 2)])" />
        <template v-else-if="error">
            <DataErrorBlock :title="$t('common.unable_to_load', [$t('common.document', 2)])" :retry="refresh" />
            <DocumentRequestAccess v-if="error.message.endsWith('ErrDocViewDenied')" :document-id="documentId" class="mt-2" />
        </template>
        <DataNoDataBlock
            v-else-if="doc === null"
            icon="i-mdi-file-search"
            :message="$t('common.not_found', [$t('common.document', 2)])"
        />

        <template v-else>
            <UDashboardToolbar>
                <template #default>
                    <div class="flex flex-1 snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                        <UButton
                            v-if="
                                can('DocStoreService.ToggleDocument') &&
                                checkDocAccess(access, doc.creator, AccessLevel.STATUS, 'DocStoreService.ToggleDocument')
                            "
                            class="flex-1"
                            block
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
                                checkDocAccess(access, doc.creator, AccessLevel.ACCESS, 'DocStoreService.UpdateDocument')
                            "
                            class="flex-1"
                            block
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
                            class="flex-1"
                            block
                            icon="i-mdi-frequently-asked-questions"
                            @click="openRequestsModal"
                        >
                            {{ $t('common.request', 2) }}
                        </UButton>
                        <UButton
                            v-if="
                                (doc?.creatorJob === activeChar?.job || isSuperuser) &&
                                can('DocStoreService.ChangeDocumentOwner') &&
                                checkDocAccess(access, doc?.creator, AccessLevel.EDIT, 'DocStoreService.ChangeDocumentOwner')
                            "
                            class="flex-1"
                            block
                            :disabled="doc?.creatorId === activeChar?.userId"
                            icon="i-mdi-creation"
                            @click="
                                modal.open(ConfirmModal, {
                                    confirm: async () => changeDocumentOwner(documentId),
                                })
                            "
                        >
                            {{ $t('components.documents.document_view.take_ownership') }}
                        </UButton>
                        <UButton
                            v-if="
                                can('DocStoreService.DeleteDocument') &&
                                checkDocAccess(access, doc.creator, AccessLevel.EDIT, 'DocStoreService.DeleteDocument')
                            "
                            class="flex-1"
                            block
                            :icon="!doc.deletedAt ? 'i-mdi-trash-can' : 'i-mdi-restore'"
                            @click="
                                modal.open(ConfirmModal, {
                                    confirm: async () => deleteDocument(documentId),
                                })
                            "
                        >
                            <template v-if="!doc.deletedAt">
                                {{ $t('common.delete') }}
                            </template>
                            <template v-else>
                                {{ $t('common.restore') }}
                            </template>
                        </UButton>
                    </div>
                </template>
            </UDashboardToolbar>

            <UCard>
                <template #header>
                    <div class="mb-4">
                        <h1 class="break-words px-0.5 py-1 text-4xl font-bold sm:pl-1">
                            {{ doc?.title }}
                        </h1>
                    </div>

                    <div class="mb-2 flex gap-2">
                        <div
                            v-if="doc.category"
                            class="bg-primary-100 text-primary-500 flex flex-initial flex-row gap-1 rounded-full px-2 py-1"
                        >
                            <ShapeIcon class="h-auto w-5" />
                            <span
                                class="text-primary-800 inline-flex items-center text-sm font-medium"
                                :title="doc.category.description ?? $t('common.na')"
                            >
                                {{ doc.category.name }}
                            </span>
                        </div>

                        <div v-if="doc?.closed" class="flex flex-initial flex-row gap-1 rounded-full bg-error-100 px-2 py-1">
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
                                    class="text-primary-600 hover:text-primary-400 ml-1 font-medium"
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
                </template>

                <div>
                    <h2 class="sr-only">
                        {{ $t('common.content') }}
                    </h2>
                    <div class="break-words rounded-lg bg-base-900">
                        <!-- eslint-disable vue/no-v-html -->
                        <div ref="contentRef" class="prose prose-invert min-w-full px-4 py-2" v-html="doc.content"></div>
                    </div>
                </div>

                <template #footer>
                    <UAccordion multiple :items="accordionItems" :unmount="true">
                        <template #relations>
                            <DocumentRelations :document-id="documentId" :show-source="false" />
                        </template>
                        <template #references>
                            <DocumentReferences :document-id="documentId" :show-source="false" />
                        </template>
                        <template #access>
                            <div class="mx-4 flex flex-row flex-wrap gap-1">
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
                            <div class="mx-4 flex flex-row flex-wrap gap-1">
                                <div
                                    v-for="entry in access?.users"
                                    :key="entry.id"
                                    class="bg-secondary-100 flex flex-initial snap-start flex-row items-center gap-1 whitespace-nowrap rounded-full px-2 py-1"
                                >
                                    <span class="bg-secondary-400 size-2 rounded-full" />
                                    <span
                                        class="text-secondary-700 text-sm font-medium"
                                        :title="`${$t('common.id')} ${entry.userId}`"
                                    >
                                        {{ entry.user?.firstname }}
                                        {{ entry.user?.lastname }} -
                                        {{ $t(`enums.docstore.AccessLevel.${AccessLevel[entry.access]}`) }}
                                    </span>
                                </div>
                            </div>
                        </template>
                        <template #comments>
                            <div id="comments">
                                <DocumentComments
                                    :document-id="documentId"
                                    :closed="doc?.closed"
                                    :can-comment="checkDocAccess(access, doc.creator, AccessLevel.COMMENT)"
                                    @counted="commentCount = $event"
                                    @new-comment="commentCount && commentCount++"
                                    @deleted-comment="commentCount && commentCount > 0 && commentCount--"
                                />
                            </div>
                        </template>
                        <template #activity>
                            <DocumentActivityList :document-id="documentId" />
                        </template>
                    </UAccordion>
                </template>
            </UCard>
        </template>
    </div>
</template>

<style scoped>
.prose {
    * {
        margin-top: 4px;
        margin-bottom: 4px;
    }
}
</style>
