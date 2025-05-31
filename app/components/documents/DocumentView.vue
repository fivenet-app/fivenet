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
import type { ToggleDocumentPinResponse } from '~~/gen/ts/services/documents/documents';
import ConfirmModalWithReason from '../partials/ConfirmModalWithReason.vue';
import ScrollToTop from '../partials/ScrollToTop.vue';
import DocumentReminderModal from './DocumentReminderModal.vue';

const props = defineProps<{
    documentId: number;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { attr, can, activeChar, isSuperuser } = useAuth();

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
        const call = $grpc.documents.documents.getDocument({
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
        await $grpc.documents.documents.deleteDocument({
            documentId: id,
            reason: reason,
        });

        // Navigate to document list when deletedAt timestamp is undefined
        if (doc.value?.deletedAt === undefined) {
            notifications.add({
                title: { key: 'notifications.document_deleted.title', parameters: {} },
                description: { key: 'notifications.document_deleted.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });

            await navigateTo({ name: 'documents' });
        } else {
            notifications.add({
                title: { key: 'notifications.document_restored.title', parameters: {} },
                description: { key: 'notifications.document_restored.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });

            await refresh();
        }
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function toggleDocument(id: number, closed: boolean): Promise<void> {
    try {
        await $grpc.documents.documents.toggleDocument({
            documentId: id,
            closed,
        });

        doc.value!.closed = closed;

        if (!closed) {
            notifications.add({
                title: { key: `notifications.documents.document_toggled.open.title`, parameters: {} },
                description: { key: `notifications.documents.document_toggled.open.content`, parameters: {} },
                type: NotificationType.SUCCESS,
            });
        } else {
            notifications.add({
                title: { key: `notifications.documents.document_toggled.closed.title`, parameters: {} },
                description: { key: `notifications.documents.document_toggled.closed.content`, parameters: {} },
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
        await $grpc.documents.documents.changeDocumentOwner({
            documentId: id,
        });

        notifications.add({
            title: { key: 'notifications.documents.document_take_ownership.title', parameters: {} },
            description: { key: 'notifications.documents.document_take_ownership.content', parameters: {} },
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

async function togglePin(documentId: number, state: boolean, personal: boolean): Promise<ToggleDocumentPinResponse> {
    try {
        const call = $grpc.documents.documents.toggleDocumentPin({
            documentId: documentId,
            state: state,
            personal: personal,
        });
        const { response } = await call;

        if (doc.value) {
            doc.value.pin = response.pin;
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
        can('documents.DocumentsService.ListDocumentActivity').value
            ? { slot: 'activity', label: t('common.activity'), icon: 'i-mdi-comment-quote' }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

defineShortcuts({
    'd-t': () => {
        if (
            !doc.value ||
            !(
                can('documents.DocumentsService.ToggleDocument').value &&
                checkDocAccess(access.value, doc.value.creator, AccessLevel.STATUS, 'documents.DocumentsService.ToggleDocument')
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
                can('documents.DocumentsService.UpdateDocument').value &&
                checkDocAccess(access.value, doc.value.creator, AccessLevel.EDIT, 'documents.DocumentsService.ToggleDocument')
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
        if (!doc.value || !can('documents.DocumentsService.ListDocumentReqs').value) {
            return;
        }

        openRequestsModal();
    },
});

const scrollRef = useTemplateRef('scrollRef');
</script>

<template>
    <UDashboardNavbar class="print:hidden" :title="$t('pages.documents.id.title')">
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

    <UDashboardPanelContent class="p-0 sm:pb-0">
        <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.document', 1)])" />
        <template v-else-if="error">
            <DataErrorBlock :title="$t('common.unable_to_load', [$t('common.document', 1)])" :error="error" :retry="refresh" />
            <DocumentRequestAccess
                v-if="error.message.includes('ErrDocViewDenied')"
                class="mt-2 w-full"
                :document-id="documentId"
            />
        </template>
        <DataNoDataBlock
            v-else-if="!doc"
            icon="i-mdi-file-search"
            :message="$t('common.not_found', [$t('common.document', 1)])"
        />

        <template v-else>
            <UDashboardToolbar class="print:hidden">
                <template #default>
                    <div class="flex flex-1 snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                        <UTooltip
                            v-if="
                                can('documents.DocumentsService.ToggleDocument').value &&
                                checkDocAccess(
                                    access,
                                    doc.creator,
                                    AccessLevel.STATUS,
                                    'documents.DocumentsService.ToggleDocument',
                                )
                            "
                            class="flex-1"
                            :text="`${$t('common.open', 1)}/ ${$t('common.close')}`"
                            :shortcuts="['D', 'T']"
                        >
                            <UButton
                                class="flex-1 flex-col"
                                block
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
                                can('documents.DocumentsService.UpdateDocument').value &&
                                checkDocAccess(
                                    access,
                                    doc.creator,
                                    AccessLevel.ACCESS,
                                    'documents.DocumentsService.UpdateDocument',
                                )
                            "
                            class="flex-1"
                            :text="$t('common.edit')"
                            :shortcuts="['D', 'E']"
                        >
                            <UButton
                                class="flex-1 flex-col"
                                block
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
                            v-if="can('documents.DocumentsService.ToggleDocumentPin').value"
                            class="flex flex-1"
                            :text="`${$t('common.pin', 1)}/ ${$t('common.unpin')}`"
                        >
                            <UButtonGroup class="flex flex-1">
                                <UButton
                                    class="flex-1 flex-col"
                                    block
                                    :color="doc.pin?.state && doc.pin?.userId ? 'error' : 'primary'"
                                    @click="togglePin(documentId, !doc.pin?.userId, true)"
                                >
                                    <UIcon
                                        class="size-5"
                                        :name="
                                            doc.pin?.state && doc.pin?.userId ? 'i-mdi-playlist-remove' : 'i-mdi-playlist-plus'
                                        "
                                    />
                                    {{ $t('common.personal') }}
                                </UButton>

                                <UButton
                                    v-if="attr('documents.DocumentsService.ToggleDocumentPin', 'Types', 'JobWide').value"
                                    class="flex-1 flex-col"
                                    block
                                    :color="doc.pin?.state && doc.pin?.job ? 'error' : 'primary'"
                                    @click="togglePin(documentId, !doc.pin?.job, false)"
                                >
                                    <UIcon
                                        class="size-5"
                                        :name="doc.pin?.state && doc.pin?.job ? 'i-mdi-pin-off' : 'i-mdi-pin'"
                                    />
                                    {{ $t('common.job') }}
                                </UButton>
                            </UButtonGroup>
                        </UTooltip>

                        <UTooltip
                            v-if="can('documents.DocumentsService.ListDocumentReqs').value"
                            class="flex-1"
                            :text="$t('common.request', 2)"
                            :shortcuts="['D', 'R']"
                        >
                            <UButton
                                class="flex-1 flex-col"
                                block
                                icon="i-mdi-frequently-asked-questions"
                                @click="openRequestsModal"
                            >
                                {{ $t('common.request', 2) }}
                            </UButton>
                        </UTooltip>

                        <UTooltip
                            v-if="can('documents.DocumentsService.SetDocumentReminder').value"
                            class="flex-1"
                            :text="$t('common.reminder')"
                        >
                            <UButton
                                class="flex-1 flex-col"
                                block
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
                        </UTooltip>

                        <UTooltip
                            v-if="
                                (doc?.creatorJob === activeChar?.job || isSuperuser) &&
                                can('documents.DocumentsService.ChangeDocumentOwner').value &&
                                checkDocAccess(
                                    access,
                                    doc?.creator,
                                    AccessLevel.EDIT,
                                    'documents.DocumentsService.ChangeDocumentOwner',
                                )
                            "
                            class="flex-1"
                            :text="$t('components.documents.document_view.take_ownership')"
                        >
                            <UButton
                                class="flex-1 flex-col"
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
                        </UTooltip>

                        <UTooltip
                            v-if="
                                can('documents.DocumentsService.DeleteDocument').value &&
                                checkDocAccess(
                                    access,
                                    doc.creator,
                                    AccessLevel.EDIT,
                                    'documents.DocumentsService.DeleteDocument',
                                )
                            "
                            class="flex-1"
                            :text="$t('common.delete')"
                        >
                            <UButton
                                class="flex-1 flex-col"
                                block
                                :color="!doc.deletedAt ? 'error' : 'success'"
                                :icon="!doc.deletedAt ? 'i-mdi-delete' : 'i-mdi-restore'"
                                :label="!doc.deletedAt ? $t('common.delete') : $t('common.restore')"
                                @click="
                                    modal.open(doc.deletedAt !== undefined ? ConfirmModal : ConfirmModalWithReason, {
                                        confirm: async (reason?: string) => deleteDocument(documentId, reason),
                                    })
                                "
                            />
                        </UTooltip>
                    </div>
                </template>
            </UDashboardToolbar>

            <UCard ref="scrollRef" class="relative overflow-x-auto">
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
                            <UIcon class="size-5" name="i-mdi-note-check" />
                            <span>
                                {{ doc.state }}
                            </span>
                        </UBadge>

                        <UBadge class="inline-flex gap-1" color="black" size="md">
                            <UIcon class="size-5" name="i-mdi-comment-text-multiple" />
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
                        <UBadge class="inline-flex gap-1" color="black" size="md">
                            <UIcon class="size-5" name="i-mdi-account" />
                            <span class="inline-flex items-center gap-1">
                                <span class="text-sm font-medium">{{ $t('common.created_by') }}</span>
                                <CitizenInfoPopover :user="doc.creator" />
                            </span>
                        </UBadge>

                        <UBadge class="inline-flex gap-1" color="black" size="md">
                            <UIcon class="size-5" name="i-mdi-calendar" />
                            <span>
                                {{ $t('common.created') }}
                                <GenericTime :value="doc.createdAt" type="long" />
                            </span>
                        </UBadge>

                        <UBadge v-if="doc.updatedAt" class="inline-flex gap-1" color="black" size="md">
                            <UIcon class="size-5" name="i-mdi-calendar-edit" />
                            <span>
                                {{ $t('common.updated') }}
                                <GenericTime :value="doc.updatedAt" type="long" />
                            </span>
                        </UBadge>

                        <UBadge v-if="doc.workflowState?.autoCloseTime" class="inline-flex gap-1" color="black" size="md">
                            <UIcon class="size-5" name="i-mdi-lock-clock" />
                            <span>
                                {{ $t('common.auto_close', 2) }}
                                <GenericTime :value="doc.workflowState.autoCloseTime" ago />
                            </span>
                        </UBadge>
                        <UBadge
                            v-else-if="doc.workflowState?.nextReminderTime"
                            class="inline-flex gap-1"
                            color="black"
                            size="md"
                        >
                            <UIcon class="size-5" name="i-mdi-reminder" />
                            <span>
                                {{ $t('common.reminder') }}
                                <GenericTime :value="doc.workflowState.nextReminderTime" ago />
                            </span>
                        </UBadge>

                        <UBadge v-if="doc.workflowUser?.manualReminderTime" class="inline-flex gap-1" color="black" size="md">
                            <UIcon class="size-5" name="i-mdi-reminder" />
                            <span>
                                {{ $t('common.reminder') }}
                                <GenericTime :value="doc.workflowUser.manualReminderTime" type="short" />
                            </span>
                        </UBadge>

                        <UBadge v-if="doc.draft" class="inline-flex gap-1" color="info" size="md">
                            <UIcon class="size-5" name="i-mdi-pencil" />
                            <span>
                                {{ $t('common.draft') }}
                            </span>
                        </UBadge>

                        <UBadge v-if="doc.deletedAt" class="inline-flex gap-1" color="amber" size="md">
                            <UIcon class="size-5" name="i-mdi-calendar-remove" />
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
                    <UAccordion class="print:hidden" multiple :items="accordionItems" :unmount="true">
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
                                    i18n-key="enums.documents"
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

                        <template v-if="can('documents.DocumentsService.ListDocumentActivity').value" #activity>
                            <UContainer>
                                <DocumentActivityList :document-id="documentId" />
                            </UContainer>
                        </template>
                    </UAccordion>
                </template>
            </UCard>

            <ScrollToTop :element="scrollRef?.$el" />
        </template>
    </UDashboardPanelContent>
</template>
