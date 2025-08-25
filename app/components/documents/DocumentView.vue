<script lang="ts" setup>
import AddToButton from '~/components/clipboard/AddToButton.vue';
import DocumentActivityList from '~/components/documents/activity/DocumentActivityList.vue';
import DocumentComments from '~/components/documents/comments/DocumentComments.vue';
import DocumentReferences from '~/components/documents/DocumentReferences.vue';
import DocumentRelations from '~/components/documents/DocumentRelations.vue';
import { checkDocAccess } from '~/components/documents/helpers';
import DocumentRequestAccess from '~/components/documents/requests/DocumentRequestAccess.vue';
import DocumentRequestModal from '~/components/documents/requests/DocumentRequestModal.vue';
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
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import type { ToggleDocumentPinResponse } from '~~/gen/ts/services/documents/documents';
import ConfirmModalWithReason from '../partials/ConfirmModalWithReason.vue';
import ScrollToTop from '../partials/ScrollToTop.vue';
import DocumentReminderModal from './DocumentReminderModal.vue';

const props = defineProps<{
    documentId: number;
}>();

const { t } = useI18n();

const { attr, can, activeChar, isSuperuser } = useAuth();

const clipboardStore = useClipboardStore();

const notifications = useNotificationsStore();

const overlay = useOverlay();

const documentsDocuments = await useDocumentsDocuments();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const commentCount = ref<undefined | number>();

const {
    data: doc,
    status,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}`, () => documentsDocuments.getDocument(props.documentId));

function addToClipboard(): void {
    if (doc.value?.document) {
        clipboardStore.addDocument(doc.value.document);
    }

    notifications.add({
        title: { key: 'notifications.clipboard.document_added.title', parameters: {} },
        description: { key: 'notifications.clipboard.document_added.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

const hash = useRouteHash();
if (hash.value !== undefined && hash.value !== null) {
    if (hash.value.replace(/^#/, '') === 'requests') {
        openRequestsModal();
    }
}

const documentRequestModal = overlay.create(DocumentRequestModal);
function openRequestsModal(): void {
    if (doc.value?.access === undefined || doc.value?.document === undefined) {
        return;
    }

    documentRequestModal.open({
        access: doc.value.access,
        doc: doc.value.document,
        onRefresh: () => refresh(),
    });
}

async function togglePin(documentId: number, state: boolean, personal: boolean): Promise<ToggleDocumentPinResponse> {
    try {
        const call = documentsDocumentsClient.toggleDocumentPin({
            documentId: documentId,
            state: state,
            personal: personal,
        });
        const { response } = await call;

        if (doc.value?.document) {
            doc.value.document.pin = response.pin;
        }

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function updateReminderTime(reminderTime?: Timestamp): void {
    if (!doc.value?.document) {
        return;
    }

    if (!doc.value.document.workflowUser) {
        doc.value.document.workflowUser = {
            documentId: props.documentId,
            userId: activeChar.value!.userId,
            reminderCount: 0,
            maxReminderCount: 10,
        };
    }

    doc.value.document.workflowUser.manualReminderTime = reminderTime;
}

async function toggleDocument(): Promise<void> {
    if (!doc.value?.document) {
        return;
    }

    doc.value.document!.closed = await documentsDocuments.toggleDocument(props.documentId, !doc.value.document?.closed);
}

const accordionItems = computed(() =>
    [
        { slot: 'relations' as const, label: t('common.relation', 2), icon: 'i-mdi-account-multiple' },
        { slot: 'references' as const, label: t('common.reference', 2), icon: 'i-mdi-file-document' },
        { slot: 'access' as const, label: t('common.access'), icon: 'i-mdi-lock', defaultOpen: true },
        { slot: 'comments' as const, label: t('common.comment', 2), icon: 'i-mdi-comment', defaultOpen: true },
        can('documents.DocumentsService/ListDocumentActivity').value
            ? { slot: 'activity' as const, label: t('common.activity'), icon: 'i-mdi-comment-quote' }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

defineShortcuts({
    'd-t': () => {
        if (
            !doc.value ||
            !(
                can('documents.DocumentsService/ToggleDocument').value &&
                checkDocAccess(
                    doc.value.access,
                    doc.value.document?.creator,
                    AccessLevel.STATUS,
                    'documents.DocumentsService/ToggleDocument',
                    doc.value?.document?.creatorJob,
                )
            )
        ) {
            return;
        }

        documentsDocuments.toggleDocument(props.documentId, !!doc.value?.document?.closed);
    },
    'd-e': () => {
        if (
            !doc.value ||
            !(
                can('documents.DocumentsService/UpdateDocument').value &&
                checkDocAccess(
                    doc.value.access,
                    doc.value.document?.creator,
                    AccessLevel.EDIT,
                    'documents.DocumentsService/ToggleDocument',
                    doc.value?.document?.creatorJob,
                )
            )
        ) {
            return;
        }

        navigateTo({
            name: 'documents-id-edit',
            params: { id: props.documentId },
        });
    },
    'd-r': () => {
        if (!doc.value || !can('documents.DocumentsService/ListDocumentReqs').value) {
            return;
        }

        openRequestsModal();
    },
});

const scrollRef = useTemplateRef('scrollRef');

const confirmModal = overlay.create(ConfirmModal);
const confirmModalWithReason = overlay.create(ConfirmModalWithReason);
const documentReminderModal = overlay.create(DocumentReminderModal, { props: { documentId: props.documentId } });
</script>

<template>
    <UDashboardNavbar class="print:hidden" :title="$t('pages.documents.id.title')">
        <template #right>
            <PartialsBackButton to="/documents" />

            <UButton
                icon="i-mdi-refresh"
                :label="$t('common.refresh')"
                :loading="isRequestPending(status)"
                @click="() => refresh()"
            />

            <UButtonGroup class="inline-flex">
                <IDCopyBadge
                    :id="doc?.document?.id ?? documentId"
                    prefix="DOC"
                    :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: {} }"
                    :content="{ key: 'notifications.document_view.copy_document_id.content', parameters: {} }"
                />

                <AddToButton :title="$t('components.clipboard.clipboard_button.add')" :callback="addToClipboard" />
            </UButtonGroup>
        </template>
    </UDashboardNavbar>

    <UDashboardPanelContent class="p-0 sm:pb-0">
        <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.document', 1)])" />
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
                                can('documents.DocumentsService/ToggleDocument').value &&
                                checkDocAccess(
                                    doc.access,
                                    doc.document?.creator,
                                    AccessLevel.STATUS,
                                    'documents.DocumentsService/ToggleDocument',
                                    doc?.document?.creatorJob,
                                )
                            "
                            class="flex-1"
                            :text="`${$t('common.open', 1)}/ ${$t('common.close')}`"
                            :shortcuts="['D', 'T']"
                        >
                            <UButton
                                class="flex-1 flex-col"
                                block
                                :icon="doc.document?.closed ? 'i-mdi-lock-open-variant' : 'i-mdi-lock'"
                                :ui="{ leadingIcon: doc.document?.closed ? 'text-success-500' : 'text-success-500' }"
                                @click="toggleDocument()"
                            >
                                <template v-if="doc.document?.closed">
                                    {{ $t('common.open', 1) }}
                                </template>
                                <template v-else>
                                    {{ $t('common.close', 1) }}
                                </template>
                            </UButton>
                        </UTooltip>

                        <UTooltip
                            v-if="
                                can('documents.DocumentsService/UpdateDocument').value &&
                                checkDocAccess(
                                    doc.access,
                                    doc.document?.creator,
                                    AccessLevel.ACCESS,
                                    'documents.DocumentsService/UpdateDocument',
                                    doc?.document?.creatorJob,
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
                                    params: { id: doc.document?.id },
                                }"
                                icon="i-mdi-pencil"
                            >
                                {{ $t('common.edit') }}
                            </UButton>
                        </UTooltip>

                        <UTooltip
                            v-if="can('documents.DocumentsService/ToggleDocumentPin').value"
                            class="flex flex-1"
                            :text="`${$t('common.pin', 1)}/ ${$t('common.unpin')}`"
                        >
                            <UButtonGroup class="flex flex-1">
                                <UButton
                                    class="flex-1 flex-col"
                                    block
                                    :color="doc.document?.pin?.state && doc.document?.pin?.userId ? 'error' : 'primary'"
                                    @click="togglePin(documentId, !doc.document?.pin?.userId, true)"
                                >
                                    <UIcon
                                        class="size-5"
                                        :name="
                                            doc.document?.pin?.state && doc.document?.pin?.userId
                                                ? 'i-mdi-playlist-remove'
                                                : 'i-mdi-playlist-plus'
                                        "
                                    />
                                    {{ $t('common.personal') }}
                                </UButton>

                                <UButton
                                    v-if="attr('documents.DocumentsService/ToggleDocumentPin', 'Types', 'JobWide').value"
                                    class="flex-1 flex-col"
                                    block
                                    :color="doc.document?.pin?.state && doc.document?.pin?.job ? 'error' : 'primary'"
                                    @click="togglePin(documentId, !doc.document?.pin?.job, false)"
                                >
                                    <UIcon
                                        class="size-5"
                                        :name="
                                            doc.document?.pin?.state && doc.document?.pin?.job ? 'i-mdi-pin-off' : 'i-mdi-pin'
                                        "
                                    />
                                    {{ $t('common.job') }}
                                </UButton>
                            </UButtonGroup>
                        </UTooltip>

                        <UTooltip
                            v-if="can('documents.DocumentsService/ListDocumentReqs').value"
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
                            v-if="can('documents.DocumentsService/SetDocumentReminder').value"
                            class="flex-1"
                            :text="$t('common.reminder')"
                        >
                            <UButton
                                class="flex-1 flex-col"
                                block
                                icon="i-mdi-reminder"
                                @click="
                                    documentReminderModal.open({
                                        documentId: documentId,
                                        reminderTime: doc.document?.workflowUser?.manualReminderTime ?? undefined,
                                        'onUpdate:reminderTime': ($event) => updateReminderTime($event),
                                    })
                                "
                            >
                                {{ $t('common.reminder') }}
                            </UButton>
                        </UTooltip>

                        <UTooltip
                            v-if="
                                (doc?.document?.creatorJob === activeChar?.job || isSuperuser) &&
                                can('documents.DocumentsService/ChangeDocumentOwner').value &&
                                checkDocAccess(
                                    doc.access,
                                    doc?.document?.creator,
                                    AccessLevel.EDIT,
                                    'documents.DocumentsService/ChangeDocumentOwner',
                                    doc?.document?.creatorJob,
                                )
                            "
                            class="flex-1"
                            :text="$t('components.documents.document_view.take_ownership')"
                        >
                            <UButton
                                class="flex-1 flex-col"
                                block
                                :disabled="doc?.document?.creatorId === activeChar?.userId"
                                icon="i-mdi-creation"
                                @click="
                                    confirmModal.open({
                                        confirm: async () => documentsDocuments.changeDocumentOwner(documentId),
                                    })
                                "
                            >
                                {{ $t('components.documents.document_view.take_ownership') }}
                            </UButton>
                        </UTooltip>

                        <UTooltip
                            v-if="
                                can('documents.DocumentsService/DeleteDocument').value &&
                                checkDocAccess(
                                    doc.access,
                                    doc.document?.creator,
                                    AccessLevel.EDIT,
                                    'documents.DocumentsService/DeleteDocument',
                                    doc?.document?.creatorJob,
                                )
                            "
                            class="flex-1"
                            :text="$t('common.delete')"
                        >
                            <UButton
                                class="flex-1 flex-col"
                                block
                                :color="!doc.document?.deletedAt ? 'error' : 'success'"
                                :icon="!doc.document?.deletedAt ? 'i-mdi-delete' : 'i-mdi-restore'"
                                :label="!doc.document?.deletedAt ? $t('common.delete') : $t('common.restore')"
                                @click="
                                    (doc.document?.deletedAt !== undefined ? confirmModalWithReason : confirmModal).open({
                                        confirm: async (reason?: string) =>
                                            documentsDocuments.deleteDocument(
                                                documentId,
                                                isSuperuser && doc?.document?.deletedAt !== undefined,
                                                reason,
                                            ),
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
                        <h1 class="px-0.5 py-1 text-4xl font-bold break-words sm:pl-1">
                            <span v-if="!doc.document?.title" class="italic">
                                {{ $t('common.untitled') }}
                            </span>
                            <span v-else>
                                {{ doc.document?.title }}
                            </span>
                        </h1>
                    </div>

                    <div class="mb-2 flex gap-2">
                        <DocumentCategoryBadge :category="doc.document?.category" />

                        <OpenClosedBadge :closed="doc.document?.closed" size="md" />

                        <UBadge v-if="doc.document?.state" class="inline-flex gap-1" size="md">
                            <UIcon class="size-5" name="i-mdi-note-check" />
                            <span>
                                {{ doc.document?.state }}
                            </span>
                        </UBadge>

                        <UBadge class="inline-flex gap-1" color="neutral" size="md">
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
                        <UBadge class="inline-flex gap-1" color="neutral" size="md">
                            <UIcon class="size-5" name="i-mdi-account" />
                            <span class="inline-flex items-center gap-1">
                                <span class="text-sm font-medium">{{ $t('common.created_by') }}</span>
                                <CitizenInfoPopover :user="doc.document?.creator" />
                            </span>
                        </UBadge>

                        <UBadge class="inline-flex gap-1" color="neutral" size="md">
                            <UIcon class="size-5" name="i-mdi-calendar" />
                            <span>
                                {{ $t('common.created') }}
                                <GenericTime :value="doc.document?.createdAt" type="long" />
                            </span>
                        </UBadge>

                        <UBadge v-if="doc.document?.updatedAt" class="inline-flex gap-1" color="neutral" size="md">
                            <UIcon class="size-5" name="i-mdi-calendar-edit" />
                            <span>
                                {{ $t('common.updated') }}
                                <GenericTime :value="doc.document?.updatedAt" type="long" />
                            </span>
                        </UBadge>

                        <UBadge
                            v-if="doc.document?.workflowState?.autoCloseTime"
                            class="inline-flex gap-1"
                            color="neutral"
                            size="md"
                        >
                            <UIcon class="size-5" name="i-mdi-lock-clock" />
                            <span>
                                {{ $t('common.auto_close', 2) }}
                                <GenericTime :value="doc.document?.workflowState?.autoCloseTime" ago />
                            </span>
                        </UBadge>
                        <UBadge
                            v-else-if="doc.document?.workflowState?.nextReminderTime"
                            class="inline-flex gap-1"
                            color="neutral"
                            size="md"
                        >
                            <UIcon class="size-5" name="i-mdi-reminder" />
                            <span>
                                {{ $t('common.reminder') }}
                                <GenericTime :value="doc.document?.workflowState?.nextReminderTime" ago />
                            </span>
                        </UBadge>

                        <UBadge
                            v-if="doc.document?.workflowUser?.manualReminderTime"
                            class="inline-flex gap-1"
                            color="neutral"
                            size="md"
                        >
                            <UIcon class="size-5" name="i-mdi-reminder" />
                            <span>
                                {{ $t('common.reminder') }}
                                <GenericTime :value="doc.document?.workflowUser?.manualReminderTime" type="short" />
                            </span>
                        </UBadge>

                        <UBadge v-if="doc.document?.draft" class="inline-flex gap-1" color="info" size="md">
                            <UIcon class="size-5" name="i-mdi-pencil" />
                            <span>
                                {{ $t('common.draft') }}
                            </span>
                        </UBadge>

                        <UBadge v-if="doc.document?.deletedAt" class="inline-flex gap-1" color="warning" size="md">
                            <UIcon class="size-5" name="i-mdi-calendar-remove" />
                            <span>
                                {{ $t('common.deleted') }}
                                <GenericTime :value="doc.document?.deletedAt" type="long" />
                            </span>
                        </UBadge>
                    </div>
                </template>

                <div>
                    <h2 class="sr-only">
                        {{ $t('common.content') }}
                    </h2>

                    <div class="dark:bg-base-900 mx-auto w-full max-w-(--breakpoint-xl) rounded-lg bg-neutral-100 break-words">
                        <HTMLContent
                            v-if="doc.document?.content?.content"
                            class="px-4 py-2"
                            :value="doc.document.content.content"
                        />
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
                                    v-if="!doc.access || (doc.access?.jobs.length === 0 && doc.access?.users.length === 0)"
                                    icon="i-mdi-file-search"
                                    :message="$t('common.not_found', [$t('common.access', 2)])"
                                />

                                <AccessBadges
                                    v-else
                                    :access-level="AccessLevel"
                                    :jobs="doc.access.jobs"
                                    :users="doc.access.users"
                                    i18n-key="enums.documents"
                                />
                            </UContainer>
                        </template>

                        <template #comments>
                            <UContainer>
                                <div id="comments">
                                    <DocumentComments
                                        :document-id="documentId"
                                        :closed="doc.document?.closed"
                                        :can-comment="checkDocAccess(doc.access, doc.document?.creator, AccessLevel.COMMENT)"
                                        @counted="commentCount = $event"
                                        @new-comment="commentCount && commentCount++"
                                        @deleted-comment="commentCount && commentCount > 0 && commentCount--"
                                    />
                                </div>
                            </UContainer>
                        </template>

                        <template v-if="can('documents.DocumentsService/ListDocumentActivity').value" #activity>
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
