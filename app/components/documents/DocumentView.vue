<script lang="ts" setup>
import AddToButton from '~/components/clipboard/AddToButton.vue';
import DocumentActivityList from '~/components/documents/activity/DocumentActivityList.vue';
import DocumentComments from '~/components/documents/comments/DocumentComments.vue';
import DocumentReferences from '~/components/documents/DocumentReferences.vue';
import DocumentRelations from '~/components/documents/DocumentRelations.vue';
import { checkDocAccess } from '~/components/documents/helpers';
import DocumentRequestAccess from '~/components/documents/requests/DocumentRequestAccess.vue';
import DocumentRequestsModal from '~/components/documents/requests/DocumentRequestsModal.vue';
import AccessBadges from '~/components/partials/access/AccessBadges.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import HTMLContent from '~/components/partials/content/HTMLContent.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DocumentCategoryBadge from '~/components/partials/documents/DocumentCategoryBadge.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import { useClipboardStore } from '~/stores/clipboard';
import { useNotificatorStore } from '~/stores/notificator';
import type { DocumentAccess } from '~~/gen/ts/resources/documents/access';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import type { Document } from '~~/gen/ts/resources/documents/documents';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import type { ToggleDocumentPinResponse } from '~~/gen/ts/services/docstore/docstore';
import DocumentReminderModal from './DocumentReminderModal.vue';

const props = defineProps<{
    documentId: number;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { can, activeChar, isSuperuser } = useAuth();

const clipboardStore = useClipboardStore();

const notifications = useNotificatorStore();

const modal = useModal();

const access = ref<undefined | DocumentAccess>(undefined);
const commentCount = ref<undefined | number>();

const {
    data: doc,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}`, () => getDocument(props.documentId));

async function getDocument(id: number): Promise<Document> {
    try {
        const call = $grpc.docstore.docStore.getDocument({
            documentId: id,
        });
        const { response } = await call;

        access.value = response.access;

        return response.document!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteDocument(id: number, reason?: string): Promise<void> {
    try {
        await $grpc.docstore.docStore.deleteDocument({
            documentId: id,
            reason: reason,
        });

        notifications.add({
            title: { key: 'notifications.document_deleted.title', parameters: {} },
            description: { key: 'notifications.document_deleted.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await navigateTo({ name: 'documents' });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function toggleDocument(id: number, closed: boolean): Promise<void> {
    try {
        await $grpc.docstore.docStore.toggleDocument({
            documentId: id,
            closed,
        });

        doc.value!.closed = closed;

        if (!closed) {
            notifications.add({
                title: { key: `notifications.docstore.document_toggled.open.title`, parameters: {} },
                description: { key: `notifications.docstore.document_toggled.open.content`, parameters: {} },
                type: NotificationType.SUCCESS,
            });
        } else {
            notifications.add({
                title: { key: `notifications.docstore.document_toggled.closed.title`, parameters: {} },
                description: { key: `notifications.docstore.document_toggled.closed.content`, parameters: {} },
                type: NotificationType.SUCCESS,
            });
        }
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function changeDocumentOwner(id: number): Promise<void> {
    try {
        await $grpc.docstore.docStore.changeDocumentOwner({
            documentId: id,
        });

        notifications.add({
            title: { key: 'notifications.docstore.document_take_ownership.title', parameters: {} },
            description: { key: 'notifications.docstore.document_take_ownership.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await refresh();
    } catch (e) {
        handleGRPCError(e as RpcError);
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
        type: NotificationType.INFO,
    });
}

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
        onRefresh: () => refresh(),
    });
}

async function togglePin(documentId: number, state: boolean): Promise<ToggleDocumentPinResponse> {
    try {
        const call = $grpc.docstore.docStore.toggleDocumentPin({
            documentId: documentId,
            state: state,
        });
        const { response } = await call;

        if (doc.value) {
            doc.value.pinned = response.state;
        }

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function updateReminderTime(reminderTime?: Timestamp): void {
    if (!doc.value) {
        return;
    }

    if (!doc.value.workflowUser) {
        doc.value.workflowUser = {
            documentId: props.documentId,
            userId: activeChar.value!.userId,
        };
    }

    doc.value.workflowUser.manualReminderTime = reminderTime;
}

const accordionItems = computed(() =>
    [
        { slot: 'relations', label: t('common.relation', 2), icon: 'i-mdi-account-multiple' },
        { slot: 'references', label: t('common.reference', 2), icon: 'i-mdi-file-document' },
        { slot: 'access', label: t('common.access'), icon: 'i-mdi-lock', defaultOpen: true },
        { slot: 'comments', label: t('common.comment', 2), icon: 'i-mdi-comment', defaultOpen: true },
        can('DocStoreService.ListDocumentActivity').value
            ? { slot: 'activity', label: t('common.activity'), icon: 'i-mdi-comment-quote' }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

defineShortcuts({
    'd-t': () => {
        if (
            !doc.value ||
            !(
                can('DocStoreService.ToggleDocument').value &&
                checkDocAccess(access.value, doc.value.creator, AccessLevel.STATUS, 'DocStoreService.ToggleDocument')
            )
        ) {
            return;
        }

        toggleDocument(props.documentId, doc.value?.closed);
    },
    'd-e': () => {
        if (
            !doc.value ||
            !(
                can('DocStoreService.UpdateDocument').value &&
                checkDocAccess(access.value, doc.value.creator, AccessLevel.EDIT, 'DocStoreService.ToggleDocument')
            )
        ) {
            return;
        }

        navigateTo({
            name: 'documents-id-edit',
            params: { id: doc.value.id },
        });
    },
    'd-r': () => {
        if (!doc.value || !can('DocStoreService.ListDocumentReqs').value) {
            return;
        }

        openRequestsModal();
    },
});
</script>

<template>
    <UDashboardNavbar :title="$t('pages.documents.id.title')" class="print:hidden">
        <template #right>
            <PartialsBackButton to="/documents" />

            <UButtonGroup class="inline-flex">
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

    <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.document', 2)])" />
    <template v-else-if="error">
        <DataErrorBlock :title="$t('common.unable_to_load', [$t('common.document', 2)])" :error="error" :retry="refresh" />
        <DocumentRequestAccess
            v-if="error.message.includes('ErrDocViewDenied')"
            :document-id="documentId"
            class="mt-2 w-full"
        />
    </template>
    <DataNoDataBlock v-else-if="!doc" icon="i-mdi-file-search" :message="$t('common.not_found', [$t('common.document', 2)])" />

    <template v-else>
        <UDashboardToolbar class="print:hidden">
            <template #default>
                <div class="flex flex-1 snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                    <UTooltip
                        v-if="
                            can('DocStoreService.ToggleDocument').value &&
                            checkDocAccess(access, doc.creator, AccessLevel.STATUS, 'DocStoreService.ToggleDocument')
                        "
                        class="flex-1"
                        :text="`${$t('common.open', 1)}/ ${$t('common.close')}`"
                        :shortcuts="['D', 'T']"
                    >
                        <UButton
                            block
                            class="flex-1 flex-col"
                            :icon="doc.closed ? 'i-mdi-lock-open-variant' : 'i-mdi-lock'"
                            :ui="{ icon: { base: doc.closed ? 'text-success-500' : 'text-success-500' } }"
                            @click="toggleDocument(documentId, !doc.closed)"
                        >
                            <template v-if="doc.closed">
                                {{ $t('common.open', 1) }}
                            </template>
                            <template v-else>
                                {{ $t('common.close', 1) }}
                            </template>
                        </UButton>
                    </UTooltip>

                    <UTooltip
                        v-if="
                            can('DocStoreService.UpdateDocument').value &&
                            checkDocAccess(access, doc.creator, AccessLevel.ACCESS, 'DocStoreService.UpdateDocument')
                        "
                        class="flex-1"
                        :text="$t('common.edit')"
                        :shortcuts="['D', 'E']"
                    >
                        <UButton
                            block
                            class="flex-1 flex-col"
                            :to="{
                                name: 'documents-id-edit',
                                params: { id: doc.id },
                            }"
                            icon="i-mdi-pencil"
                        >
                            {{ $t('common.edit') }}
                        </UButton>
                    </UTooltip>

                    <UTooltip
                        v-if="can('DocStoreService.ToggleDocumentPin').value"
                        class="flex-1"
                        :text="`${$t('common.pin', 1)}/ ${$t('common.unpin')}`"
                    >
                        <UButton block class="flex-1 flex-col" @click="togglePin(documentId, !doc.pinned)">
                            <template v-if="!doc.pinned">
                                <UIcon name="i-mdi-pin" class="size-5" />
                                {{ $t('common.pin') }}
                            </template>
                            <template v-else>
                                <UIcon name="i-mdi-pin-off" class="size-5" />
                                {{ $t('common.unpin') }}
                            </template>
                        </UButton>
                    </UTooltip>

                    <UTooltip
                        v-if="can('DocStoreService.ListDocumentReqs').value"
                        class="flex-1"
                        :text="$t('common.request', 2)"
                        :shortcuts="['D', 'R']"
                    >
                        <UButton
                            block
                            class="flex-1 flex-col"
                            icon="i-mdi-frequently-asked-questions"
                            @click="openRequestsModal"
                        >
                            {{ $t('common.request', 2) }}
                        </UButton>
                    </UTooltip>

                    <UButton
                        v-if="can('DocStoreService.SetDocumentReminder').value"
                        block
                        class="flex-1 flex-col"
                        icon="i-mdi-reminder"
                        @click="
                            modal.open(DocumentReminderModal, {
                                documentId: documentId,
                                reminderTime: doc.workflowUser?.manualReminderTime ?? undefined,
                                'onUpdate:reminderTime': () => updateReminderTime($event),
                            })
                        "
                    >
                        {{ $t('common.reminder') }}
                    </UButton>

                    <UButton
                        v-if="
                            (doc?.creatorJob === activeChar?.job || isSuperuser) &&
                            can('DocStoreService.ChangeDocumentOwner').value &&
                            checkDocAccess(access, doc?.creator, AccessLevel.EDIT, 'DocStoreService.ChangeDocumentOwner')
                        "
                        block
                        class="flex-1 flex-col"
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
                            can('DocStoreService.DeleteDocument').value &&
                            checkDocAccess(access, doc.creator, AccessLevel.EDIT, 'DocStoreService.DeleteDocument')
                        "
                        block
                        class="flex-1 flex-col"
                        :color="!doc.deletedAt ? 'red' : 'green'"
                        :icon="!doc.deletedAt ? 'i-mdi-trash-can' : 'i-mdi-restore'"
                        :label="!doc.deletedAt ? $t('common.delete') : $t('common.restore')"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteDocument(documentId),
                            })
                        "
                    />
                </div>
            </template>
        </UDashboardToolbar>

        <UCard class="relative overflow-x-auto">
            <template #header>
                <div class="mb-4">
                    <h1 class="break-words px-0.5 py-1 text-4xl font-bold sm:pl-1">
                        {{ doc.title }}
                    </h1>
                </div>

                <div class="mb-2 flex gap-2">
                    <DocumentCategoryBadge :category="doc.category" />

                    <OpenClosedBadge :closed="doc.closed" size="md" />

                    <UBadge v-if="doc.state" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-note-check" class="size-5" />
                        <span>
                            {{ doc.state }}
                        </span>
                    </UBadge>

                    <UBadge color="black" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-comment-text-multiple" class="size-5" />
                        <span>
                            {{
                                commentCount !== undefined
                                    ? $t('common.comments', commentCount)
                                    : '? ' + $t('common.comment', 2)
                            }}
                        </span>
                    </UBadge>
                </div>

                <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-0">
                    <UBadge color="black" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-account" class="size-5" />
                        <span class="inline-flex items-center gap-1">
                            <span class="text-sm font-medium">{{ $t('common.created_by') }}</span>
                            <CitizenInfoPopover :user="doc.creator" />
                        </span>
                    </UBadge>

                    <UBadge color="black" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-calendar" class="size-5" />
                        <span>
                            {{ $t('common.created_at') }}
                            <GenericTime :value="doc.createdAt" type="long" />
                        </span>
                    </UBadge>

                    <UBadge v-if="doc.updatedAt" color="black" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-calendar-edit" class="size-5" />
                        <span>
                            {{ $t('common.updated_at') }}
                            <GenericTime :value="doc.updatedAt" type="long" />
                        </span>
                    </UBadge>

                    <UBadge v-if="doc.workflowState?.autoCloseTime" color="black" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-lock-clock" class="size-5" />
                        <span>
                            {{ $t('common.auto_close', 2) }}
                            <GenericTime :value="doc.workflowState.autoCloseTime" ago />
                        </span>
                    </UBadge>
                    <UBadge v-else-if="doc.workflowState?.nextReminderTime" color="black" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-reminder" class="size-5" />
                        <span>
                            {{ $t('common.reminder') }}
                            <GenericTime :value="doc.workflowState.nextReminderTime" ago />
                        </span>
                    </UBadge>

                    <UBadge v-if="doc.workflowUser?.manualReminderTime" color="black" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-reminder" class="size-5" />
                        <span>
                            {{ $t('common.reminder') }}
                            <GenericTime :value="doc.workflowUser.manualReminderTime" type="short" />
                        </span>
                    </UBadge>

                    <UBadge v-if="doc.deletedAt" color="amber" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-calendar-remove" class="size-5" />
                        <span>
                            {{ $t('common.deleted') }}
                            <GenericTime :value="doc.deletedAt" type="long" />
                        </span>
                    </UBadge>
                </div>
            </template>

            <div>
                <h2 class="sr-only">
                    {{ $t('common.content') }}
                </h2>
                <div class="mx-auto w-full max-w-screen-xl break-words rounded-lg bg-neutral-100 dark:bg-base-900">
                    <HTMLContent v-if="doc.content?.content" class="px-4 py-2" :value="doc.content.content" />
                </div>
            </div>

            <template #footer>
                <UAccordion multiple :items="accordionItems" :unmount="true" class="print:hidden">
                    <template #relations>
                        <UContainer>
                            <DocumentRelations :document-id="documentId" :show-document="false" />
                        </UContainer>
                    </template>

                    <template #references>
                        <UContainer>
                            <DocumentReferences :document-id="documentId" :show-source="false" />
                        </UContainer>
                    </template>

                    <template #access>
                        <UContainer>
                            <DataNoDataBlock
                                v-if="!access || (access?.jobs.length === 0 && access?.users.length === 0)"
                                icon="i-mdi-file-search"
                                :message="$t('common.not_found', [$t('common.access', 2)])"
                            />

                            <AccessBadges
                                v-else
                                :access-level="AccessLevel"
                                :jobs="access.jobs"
                                :users="access.users"
                                i18n-key="enums.docstore"
                            />
                        </UContainer>
                    </template>

                    <template #comments>
                        <UContainer>
                            <div id="comments">
                                <DocumentComments
                                    :document-id="documentId"
                                    :closed="doc.closed"
                                    :can-comment="checkDocAccess(access, doc.creator, AccessLevel.COMMENT)"
                                    @counted="commentCount = $event"
                                    @new-comment="commentCount && commentCount++"
                                    @deleted-comment="commentCount && commentCount > 0 && commentCount--"
                                />
                            </div>
                        </UContainer>
                    </template>

                    <template v-if="can('DocStoreService.ListDocumentActivity').value" #activity>
                        <UContainer>
                            <DocumentActivityList :document-id="documentId" />
                        </UContainer>
                    </template>
                </UAccordion>
            </template>
        </UCard>
    </template>
</template>
