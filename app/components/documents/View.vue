<script lang="ts" setup>
import AddToButton from '~/components/clipboard/AddToButton.vue';
import List from '~/components/documents/activity/List.vue';
import Comments from '~/components/documents/comments/Comments.vue';
import { checkDocAccess } from '~/components/documents/helpers';
import References from '~/components/documents/References.vue';
import Relations from '~/components/documents/Relations.vue';
import RequestAccess from '~/components/documents/requests/RequestAccess.vue';
import AccessBadges from '~/components/partials/access/AccessBadges.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import CategoryBadge from '~/components/partials/documents/CategoryBadge.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import { useClipboardStore } from '~/stores/clipboard';
import { AccessLevel } from '~~/gen/ts/resources/documents/access/access';
import type { DocumentData } from '~~/gen/ts/resources/documents/data/data';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import type { ToggleDocumentPinResponse } from '~~/gen/ts/services/documents/documents';
import { normalizeDocumentData } from '~/components/quickbuttons/penaltycalculator/helpers';
import ConfirmModalWithReason from '../partials/ConfirmModalWithReason.vue';
import CustomContentRenderer from '../partials/content/CustomContentRenderer.vue';
import DraftBadge from '../partials/DraftBadge.vue';
import RefreshButton from '../partials/RefreshButton.vue';
import ResponsiveActions from '../partials/ResponsiveActions.vue';
import {
    separator as actionSeparator,
    type ResponsiveActionEntry,
    type ResponsiveActionItem,
} from '../partials/ResponsiveActions.types';
import ScrollToTop from '../partials/ScrollToTop.vue';
import ApprovalDrawer from './approval/ApprovalDrawer.vue';
import ReminderDrawer from './ReminderDrawer.vue';
import RequestDrawer from './requests/RequestDrawer.vue';
import ApprovalBadge from './approval/ApprovalBadge.vue';
import { addDays, addMonths, addWeeks } from 'date-fns';

const props = defineProps<{
    documentId: number;
}>();

const { t } = useI18n();

const { attr, can, activeChar, isSuperuser } = useAuth();

const clipboardStore = useClipboardStore();

const notifications = useNotificationsStore();

const settingsStore = useSettingsStore();
const { design } = storeToRefs(settingsStore);

const overlay = useOverlay();

const documentsDocuments = await useDocumentsDocuments();

const {
    data: doc,
    status,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}`, () => documentsDocuments.getDocument(props.documentId));

const documentData = computed<DocumentData | undefined>(() => normalizeDocumentData(doc.value?.document?.data));
provide('documents:content:data', documentData);

useHead({
    title: () =>
        doc.value?.document?.title
            ? `${doc.value.document.title} - ${t('pages.documents.id.title')}`
            : t('pages.documents.id.title'),
});

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

const canDo = computed(() => ({
    status:
        can('documents.DocumentsService/ToggleDocument').value &&
        checkDocAccess(
            doc.value?.access,
            doc.value?.document?.creator,
            AccessLevel.STATUS,
            'documents.DocumentsService/ToggleDocument',
            doc.value?.document?.creatorJob,
        ),
    accessUpdate: checkDocAccess(
        doc.value?.access,
        doc.value?.document?.creator,
        AccessLevel.ACCESS,
        undefined,
        doc.value?.document?.creatorJob,
    ),
    contentUpdate: checkDocAccess(
        doc.value?.access,
        doc.value?.document?.creator,
        AccessLevel.EDIT,
        undefined,
        doc.value?.document?.creatorJob,
    ),
    pin: can('documents.DocumentsService/ToggleDocumentPin').value,
    requests: can('documents.DocumentsService/ListDocumentReqs').value,
    approve: can('documents.DocumentsService/ListDocuments').value,
    reminder: can('documents.DocumentsService/SetDocumentReminder').value,
    takeOwnership:
        (doc.value?.document?.creatorJob === activeChar.value?.job || isSuperuser.value) &&
        can('documents.DocumentsService/ChangeDocumentOwner').value &&
        checkDocAccess(
            doc.value?.access,
            doc.value?.document?.creator,
            AccessLevel.EDIT,
            'documents.DocumentsService/ChangeDocumentOwner',
            doc.value?.document?.creatorJob,
        ),
    delete:
        can('documents.DocumentsService/DeleteDocument').value &&
        checkDocAccess(
            doc.value?.access,
            doc.value?.document?.creator,
            AccessLevel.EDIT,
            'documents.DocumentsService/DeleteDocument',
            doc.value?.document?.creatorJob,
        ),
}));

const requestDrawer = overlay.create(RequestDrawer);
const approvalDrawer = overlay.create(ApprovalDrawer);

const hash = useRouteHash('', { mode: 'push' });

async function openRequestsDrawer(): Promise<void> {
    if (doc.value?.access === undefined || doc.value?.document === undefined) return;

    requestDrawer
        .open({
            access: doc.value.access,
            doc: doc.value.document,
            onRefresh: () => refresh(),
        })
        .then(() => (hash.value = ''));

    hash.value = `#requests`;
}

async function openApprovalDrawer(): Promise<void> {
    approvalDrawer
        .open({
            documentId: props.documentId,
            docCreatorId: doc.value?.document?.creatorId,
            docMeta: doc.value!.document!.meta,
            canEdit: canDo.value.contentUpdate,
            'onUpdate:docMeta': ($event) => {
                if (doc.value?.document) doc.value.document.meta = $event;
            },
        })
        .then(() => (hash.value = ''))
        .finally(() => {
            hash.value = '';
        });

    hash.value = `#approvals`;
}

async function handleHash(): Promise<void> {
    if (hash.value === undefined || hash.value === null) return;

    const val = hash.value.replace(/^#/, '');
    if (val === 'requests') {
        openRequestsDrawer();
    } else if (val === 'approvals') {
        openApprovalDrawer();
    } else if (val === 'comments') {
        nextTick(() => {
            const el = document.getElementById('comments');
            if (el) {
                el.scrollIntoView({ behavior: 'smooth' });
            }
        });
    }
}

watch(doc, () => handleHash());

async function togglePin(documentId: number, state: boolean, personal: boolean): Promise<ToggleDocumentPinResponse> {
    try {
        const response = await documentsDocuments.togglePin(documentId, state, personal);

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
    if (!doc.value?.document) return;

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
    if (!doc.value?.document) return;
    if (!doc.value?.document?.meta) {
        doc.value.document!.meta = {
            documentId: props.documentId,
            closed: false,
            draft: false,
            approved: false,
            public: false,
            state: '',
        };
    }

    doc.value.document!.meta.closed = await documentsDocuments.toggleDocument(
        props.documentId,
        !doc.value.document?.meta?.closed,
    );
}

const accordionItems = computed(() =>
    [
        { value: 'relations', slot: 'relations' as const, label: t('common.relation', 2), icon: 'i-mdi-account-multiple' },
        { value: 'references', slot: 'references' as const, label: t('common.reference', 2), icon: 'i-mdi-file-document' },
        { value: 'access', slot: 'access' as const, label: t('common.access'), icon: 'i-mdi-lock' },
        { value: 'comments', slot: 'comments' as const, label: t('common.comment', 2), icon: 'i-mdi-comment' },
        can('documents.DocumentsService/ListDocumentActivity').value
            ? { value: 'activity', slot: 'activity' as const, label: t('common.activity'), icon: 'i-mdi-comment-quote' }
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
        )
            return;

        documentsDocuments.toggleDocument(props.documentId, !!doc.value?.document?.meta?.closed);
    },
    'd-e': () => {
        if (
            !doc.value ||
            !checkDocAccess(
                doc.value.access,
                doc.value.document?.creator,
                AccessLevel.ACCESS,
                undefined,
                doc.value?.document?.creatorJob,
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
        if (!doc.value || !can('documents.DocumentsService/ListDocumentReqs').value) return;

        openRequestsDrawer();
    },
});

function setCommentCount(count: number): void {
    if (!doc.value?.document) return;

    if (!doc.value?.document?.meta) {
        doc.value.document!.meta = {
            documentId: props.documentId,
            closed: false,
            draft: false,
            approved: false,
            public: false,
            state: '',
        };
    }
    doc.value.document.meta.commentCount = count;
}

const reminderActionItems = computed<ResponsiveActionItem[]>(() => [
    {
        label: t('components.documents.document_view.reminder_duration.in_days', 1),
        icon: 'i-mdi-calendar-day',
        onClick: () =>
            reminderDrawer.open({
                documentId: props.documentId,
                reminderTime: toTimestamp(addDays(new Date(), 1)),
                'onUpdate:reminderTime': ($event) => updateReminderTime($event),
            }),
    },
    {
        label: t('components.documents.document_view.reminder_duration.in_weeks', 1),
        icon: 'i-mdi-calendar-week',
        onClick: () =>
            reminderDrawer.open({
                documentId: props.documentId,
                reminderTime: toTimestamp(addWeeks(new Date(), 1)),
                'onUpdate:reminderTime': ($event) => updateReminderTime($event),
            }),
    },
    {
        label: t('components.documents.document_view.reminder_duration.in_months', 1),
        icon: 'i-mdi-calendar-month',
        onClick: () =>
            reminderDrawer.open({
                documentId: props.documentId,
                reminderTime: toTimestamp(addMonths(new Date(), 1)),
                'onUpdate:reminderTime': ($event) => updateReminderTime($event),
            }),
    },
    {
        label: t('components.documents.document_view.reminder_duration.custom'),
        icon: 'i-mdi-calendar-time',
        onClick: () => {
            reminderDrawer.open({
                documentId: props.documentId,
                reminderTime: doc.value?.document?.workflowUser?.manualReminderTime ?? undefined,
                'onUpdate:reminderTime': ($event) => updateReminderTime($event),
            });
        },
    },
]);

const actionItems = computed<ResponsiveActionEntry[]>(() => {
    if (!doc.value) return [];

    const items: ResponsiveActionEntry[] = [];

    function pushSeparator(): void {
        const lastItem = items.at(-1);
        if (lastItem?.kind !== 'separator') {
            items.push(actionSeparator());
        }
    }

    if (canDo.value.status) {
        items.push({
            label: doc.value.document?.meta?.closed ? t('common.open', 1) : t('common.close', 1),
            tooltip: `${t('common.open', 1)}/ ${t('common.close')}`,
            icon: doc.value.document?.meta?.closed ? 'i-mdi-lock-open-variant' : 'i-mdi-lock',
            color: doc.value.document?.meta?.closed ? 'success' : 'error',
            kbds: ['D', 'T'],
            onClick: () => {
                void toggleDocument();
            },
        });
    }

    if (canDo.value.accessUpdate || canDo.value.contentUpdate) {
        items.push({
            label: t('common.edit'),
            tooltip: t('common.edit'),
            icon: 'i-mdi-pencil',
            kbds: ['D', 'E'],
            color: 'neutral',
            to: {
                name: 'documents-id-edit',
                params: { id: doc.value.document?.id },
            },
        });
    }

    if (canDo.value.pin) {
        pushSeparator();

        const pinChildren: ResponsiveActionItem[] = [
            {
                label: t('common.personal'),
                icon:
                    doc.value.document?.pin?.state && doc.value.document?.pin?.userId
                        ? 'i-mdi-playlist-remove'
                        : 'i-mdi-playlist-plus',
                color: doc.value.document?.pin?.state && doc.value.document?.pin?.userId ? 'primary' : undefined,
                onClick: () => {
                    void togglePin(props.documentId, !doc.value?.document?.pin?.userId, true);
                },
            },
        ];

        if (attr('documents.DocumentsService/ToggleDocumentPin', 'Types', 'JobWide').value) {
            pinChildren.push({
                label: t('common.job'),
                icon: doc.value.document?.pin?.state && doc.value.document?.pin?.job ? 'i-mdi-pin-off' : 'i-mdi-pin',
                color: doc.value.document?.pin?.state && doc.value.document?.pin?.job ? 'primary' : undefined,
                onClick: () => {
                    void togglePin(props.documentId, !doc.value?.document?.pin?.job, false);
                },
            });
        }

        items.push({
            label: t('common.pin'),
            tooltip: `${t('common.pin', 1)}/ ${t('common.unpin')}`,
            icon: 'i-mdi-pin',
            color: 'neutral',
            children: pinChildren,
        });
    }

    if (canDo.value.requests || canDo.value.approve || canDo.value.reminder) {
        pushSeparator();
    }

    if (canDo.value.requests) {
        items.push({
            label: t('common.request', 2),
            tooltip: t('common.request', 2),
            icon: 'i-mdi-frequently-asked-questions',
            kbds: ['D', 'R'],
            color: 'neutral',
            onClick: () => {
                void openRequestsDrawer();
            },
        });
    }

    if (canDo.value.approve) {
        items.push({
            label: t('common.approve'),
            tooltip: doc.value.document?.meta?.draft
                ? t('components.documents.approval.document_not_published')
                : t('common.approve'),
            icon: 'i-mdi-approval',
            disabled: !!doc.value.document?.meta?.draft,
            color: 'neutral',
            onClick: () => {
                void openApprovalDrawer();
            },
        });
    }

    if (canDo.value.reminder) {
        items.push({
            label: t('common.reminder'),
            tooltip: t('common.reminder'),
            icon: 'i-mdi-reminder',
            color: 'neutral',
            children: reminderActionItems.value,
        });
    }

    if (canDo.value.takeOwnership) {
        pushSeparator();

        items.push({
            label: t('components.documents.document_view.take_ownership'),
            tooltip: t('components.documents.document_view.take_ownership'),
            icon: 'i-mdi-creation',
            color: 'neutral',
            disabled: doc.value.document?.creatorId === activeChar?.value?.userId,
            onClick: () => {
                confirmModal.open({
                    confirm: async () => documentsDocuments.changeDocumentOwner(props.documentId).then(() => refresh()),
                });
            },
        });
    }

    if (canDo.value.delete) {
        items.push({
            label: !doc.value.document?.deletedAt ? t('common.delete') : t('common.restore'),
            tooltip: !doc.value.document?.deletedAt ? t('common.delete') : t('common.restore'),
            color: !doc.value.document?.deletedAt ? 'error' : 'success',
            icon: !doc.value.document?.deletedAt ? 'i-mdi-delete' : 'i-mdi-restore',
            onClick: () => {
                (doc.value?.document?.deletedAt === undefined ? confirmModalWithReason : confirmModal).open({
                    confirm: async (reason?: string) =>
                        (await documentsDocuments.deleteDocument(
                            props.documentId,
                            isSuperuser && doc.value?.document?.deletedAt !== undefined,
                            reason,
                        ))
                            ? refresh()
                            : undefined,
                });
            },
        });
    }

    return items;
});

const scrollRef = useTemplateRef('scrollRef');

const confirmModal = overlay.create(ConfirmModal);
const confirmModalWithReason = overlay.create(ConfirmModalWithReason);
const reminderDrawer = overlay.create(ReminderDrawer, { props: { documentId: props.documentId } });
</script>

<template>
    <UDashboardPanel :ui="{ root: 'pb-(--page-content-bottom-offset)' }">
        <template #header>
            <UDashboardNavbar class="print:hidden" :title="$t('pages.documents.id.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton to="/documents" />

                    <RefreshButton :loading="isRequestPending(status)" @click="() => refresh()" />

                    <UFieldGroup class="inline-flex">
                        <IDCopyBadge
                            :id="doc?.document?.id ?? documentId"
                            prefix="DOC"
                            :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: {} }"
                            :content="{ key: 'notifications.document_view.copy_document_id.content', parameters: {} }"
                        />

                        <AddToButton :title="$t('components.clipboard.clipboard_button.add')" :callback="addToClipboard" />
                    </UFieldGroup>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar v-if="doc && actionItems.length" class="p-1 print:hidden">
                <ResponsiveActions :items="actionItems" :label="$t('common.action', 2)" />
            </UDashboardToolbar>

            <UDashboardToolbar v-if="doc" class="print:hidden">
                <div class="mx-auto my-2 w-full max-w-(--breakpoint-xl)">
                    <div class="mb-2">
                        <h1
                            class="px-0.5 py-1 text-4xl font-bold break-words sm:pl-1"
                            :class="design.documents.viewCollapsedTitle && 'line-clamp-1'"
                        >
                            <span v-if="!doc.document?.title" class="italic">
                                {{ $t('common.untitled') }}
                            </span>
                            <span v-else>
                                {{ doc.document?.title }}
                            </span>
                        </h1>
                    </div>

                    <div class="mb-2 flex gap-2">
                        <CategoryBadge :category="doc.document?.category" />

                        <OpenClosedBadge :closed="doc.document?.meta?.closed" size="md" />

                        <ApprovalBadge :meta="doc.document?.meta" />

                        <UBadge
                            v-if="doc.document?.meta?.state"
                            class="inline-flex gap-1"
                            size="md"
                            icon="i-mdi-note-check"
                            :label="doc.document?.meta?.state"
                        />

                        <UBadge
                            class="inline-flex gap-1"
                            color="neutral"
                            size="md"
                            icon="i-mdi-comment-text-multiple"
                            :label="
                                doc.document?.meta?.commentCount !== undefined
                                    ? $t('common.comments', doc.document?.meta?.commentCount)
                                    : '? ' + $t('common.comment', 2)
                            "
                        />
                    </div>

                    <div class="flex flex-row pb-3 sm:pb-0">
                        <div class="flex-1">
                            <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto">
                                <UBadge class="inline-flex gap-1" color="neutral" size="md" icon="i-mdi-account">
                                    <span class="inline-flex items-center gap-1">
                                        {{ $t('common.created_by') }}
                                        <CitizenInfoPopover :user="doc.document?.creator" size="xs" />
                                    </span>
                                </UBadge>

                                <UBadge class="inline-flex gap-1" color="neutral" size="md" icon="i-mdi-calendar">
                                    {{ $t('common.created') }}
                                    <GenericTime :value="doc.document?.createdAt" type="long" />
                                </UBadge>

                                <UBadge
                                    v-if="doc.document?.updatedAt"
                                    class="inline-flex gap-1"
                                    color="neutral"
                                    size="md"
                                    icon="i-mdi-calendar-edit"
                                >
                                    {{ $t('common.updated') }}
                                    <GenericTime :value="doc.document?.updatedAt" type="long" />
                                </UBadge>

                                <UBadge
                                    v-if="doc.document?.workflowState?.autoCloseTime"
                                    class="inline-flex gap-1"
                                    color="neutral"
                                    size="md"
                                    icon="i-mdi-lock-clock"
                                >
                                    {{ $t('common.auto_close', 2) }}
                                    <GenericTime :value="doc.document?.workflowState?.autoCloseTime" ago />
                                </UBadge>
                                <UBadge
                                    v-else-if="doc.document?.workflowState?.nextReminderTime"
                                    class="inline-flex gap-1"
                                    color="neutral"
                                    size="md"
                                    icon="i-mdi-reminder"
                                >
                                    {{ $t('common.reminder') }}
                                    <GenericTime :value="doc.document?.workflowState?.nextReminderTime" ago />
                                </UBadge>

                                <UBadge
                                    v-if="doc.document?.workflowUser?.manualReminderTime"
                                    class="inline-flex gap-1"
                                    color="neutral"
                                    size="md"
                                    icon="i-mdi-reminder"
                                >
                                    {{ $t('common.reminder') }}
                                    <GenericTime :value="doc.document?.workflowUser?.manualReminderTime" type="short" />
                                </UBadge>

                                <DraftBadge v-if="doc.document?.meta?.draft" />

                                <UBadge
                                    v-if="doc.document?.deletedAt"
                                    class="inline-flex gap-1"
                                    color="warning"
                                    size="md"
                                    icon="i-mdi-calendar-remove"
                                >
                                    {{ $t('common.deleted') }}
                                    <GenericTime :value="doc.document?.deletedAt" type="long" />
                                </UBadge>
                            </div>
                        </div>

                        <div>
                            <UTooltip :text="$t('common.expand_collapse')">
                                <UButton
                                    class="group place-self-end"
                                    icon="i-mdi-chevron-double-down"
                                    variant="link"
                                    size="sm"
                                    :ui="{
                                        leadingIcon: 'transition-transform duration-200 group-data-[state=open]:rotate-180',
                                    }"
                                    :data-state="design.documents.viewCollapsedTitle ? 'open' : 'closed'"
                                    @click="
                                        () => {
                                            design.documents.viewCollapsedTitle = !design.documents.viewCollapsedTitle;
                                        }
                                    "
                                />
                            </UTooltip>
                        </div>
                    </div>
                </div>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.document', 1)])" />
            <template v-else-if="error">
                <DataErrorBlock
                    :title="$t('common.unable_to_load', [$t('common.document', 1)])"
                    :error="error"
                    :retry="refresh"
                />
                <RequestAccess
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
                <div ref="scrollRef">
                    <div
                        class="mx-auto w-full max-w-(--breakpoint-xl) rounded-lg bg-neutral-100 p-4 break-words dark:bg-neutral-800"
                    >
                        <CustomContentRenderer
                            v-if="doc.document?.content"
                            :value="doc.document.content"
                            :placeholder="$t('common.no_content')"
                        />
                    </div>
                </div>

                <div class="mx-auto w-full max-w-(--breakpoint-xl)">
                    <UAccordion
                        class="print:hidden"
                        :default-value="['access', 'comments']"
                        type="multiple"
                        :items="accordionItems"
                    >
                        <template #relations>
                            <UContainer class="mb-2">
                                <Relations :document-id="documentId" :show-document="false" />
                            </UContainer>
                        </template>

                        <template #references>
                            <UContainer class="mb-2">
                                <References :document-id="documentId" :show-source="false" />
                            </UContainer>
                        </template>

                        <template #access>
                            <UContainer class="mb-2">
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
                            <UContainer class="mb-2">
                                <Comments
                                    id="comments"
                                    :document-id="documentId"
                                    :closed="doc.document?.meta?.closed"
                                    :can-comment="checkDocAccess(doc.access, doc.document?.creator, AccessLevel.COMMENT)"
                                    @counted="($event) => setCommentCount($event)"
                                    @new-comment="doc.document?.meta?.commentCount && doc.document.meta.commentCount++"
                                    @deleted-comment="
                                        doc?.document?.meta?.commentCount &&
                                        doc.document?.meta?.commentCount > 0 &&
                                        doc.document.meta.commentCount--
                                    "
                                />
                            </UContainer>
                        </template>

                        <template v-if="can('documents.DocumentsService/ListDocumentActivity').value" #activity>
                            <UContainer class="mb-2">
                                <List :document-id="documentId" />
                            </UContainer>
                        </template>
                    </UAccordion>
                </div>

                <ScrollToTop :element="scrollRef" />
            </template>
        </template>
    </UDashboardPanel>
</template>
