<script lang="ts" setup>
import { useRouteHash } from '@vueuse/router';
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
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import OpenClosedBadge from '../partials/OpenClosedBadge.vue';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    documentId: string;
}>();

const { t } = useI18n();

const clipboardStore = useClipboardStore();

const notifications = useNotificatorStore();

const authStore = useAuthStore();
const { activeChar, isSuperuser } = storeToRefs(authStore);

const modal = useModal();

const access = ref<undefined | DocumentAccess>(undefined);
const commentCount = ref<undefined | number>();

const {
    data: doc,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}`, () => getDocument(props.documentId));

async function getDocument(id: string): Promise<Document> {
    try {
        const call = getGRPCDocStoreClient().getDocument({
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

async function deleteDocument(id: string): Promise<void> {
    try {
        await getGRPCDocStoreClient().deleteDocument({
            documentId: id,
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

async function toggleDocument(id: string, closed: boolean): Promise<void> {
    try {
        await getGRPCDocStoreClient().toggleDocument({
            documentId: id,
            closed,
        });

        doc.value!.closed = closed;

        if (!closed) {
            notifications.add({
                title: { key: `notifications.document_toggled.open.title`, parameters: {} },
                description: { key: `notifications.document_toggled.open.content`, parameters: {} },
                type: NotificationType.SUCCESS,
            });
        } else {
            notifications.add({
                title: { key: `notifications.document_toggled.closed.title`, parameters: {} },
                description: { key: `notifications.document_toggled.closed.content`, parameters: {} },
                type: NotificationType.SUCCESS,
            });
        }
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function changeDocumentOwner(id: string): Promise<void> {
    try {
        await getGRPCDocStoreClient().changeDocumentOwner({
            documentId: id,
        });

        notifications.add({
            title: { key: 'notifications.document_take_ownership.title', parameters: {} },
            description: { key: 'notifications.document_take_ownership.content', parameters: {} },
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

const contentRef = ref<HTMLDivElement | null>(null);
function disableCheckboxes(): void {
    if (contentRef.value === null) {
        return;
    }

    const checkboxes: NodeListOf<HTMLInputElement> = contentRef.value.querySelectorAll('input[type=checkbox]');
    checkboxes.forEach((el) => {
        el.setAttribute('disabled', 'disabled');
        el.classList.add('form-checkbox');
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
    });
}

const accordionItems = computed(() =>
    [
        { slot: 'relations', label: t('common.relation', 2), icon: 'i-mdi-account-multiple' },
        { slot: 'references', label: t('common.reference', 2), icon: 'i-mdi-file-document' },
        { slot: 'access', label: t('common.access'), icon: 'i-mdi-lock', defaultOpen: true },
        { slot: 'comments', label: t('common.comment', 2), icon: 'i-mdi-comment', defaultOpen: true },
        can('DocStoreService.ListDocumentActivity')
            ? { slot: 'activity', label: t('common.activity'), icon: 'i-mdi-comment-quote' }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

defineShortcuts({
    'd-t': () => {
        if (
            !doc.value ||
            !(
                can('DocStoreService.ToggleDocument') &&
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
                can('DocStoreService.UpdateDocument') &&
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
        if (!doc.value || !can('DocStoreService.ListDocumentReqs')) {
            return;
        }

        openRequestsModal();
    },
});
</script>

<template>
    <UDashboardNavbar :title="$t('pages.documents.id.title')">
        <template #right>
            <UButtonGroup class="inline-flex">
                <UButton color="black" icon="i-mdi-arrow-back" to="/documents">
                    {{ $t('common.back') }}
                </UButton>

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
                    <UTooltip
                        v-if="
                            can('DocStoreService.ToggleDocument') &&
                            checkDocAccess(access, doc.creator, AccessLevel.STATUS, 'DocStoreService.ToggleDocument')
                        "
                        class="flex-1"
                        :text="`${$t('common.open', 1)}/ ${$t('common.close')}`"
                        :shortcuts="['D', 'T']"
                    >
                        <UButton block @click="toggleDocument(documentId, !doc.closed)">
                            <template v-if="doc.closed">
                                <UIcon name="i-mdi-lock-open-variant" class="size-5 text-success-500" />
                                {{ $t('common.open', 1) }}
                            </template>
                            <template v-else>
                                <UIcon name="i-mdi-lock" class="size-5 text-error-400" />
                                {{ $t('common.close', 1) }}
                            </template>
                        </UButton>
                    </UTooltip>

                    <UTooltip
                        v-if="
                            can('DocStoreService.UpdateDocument') &&
                            checkDocAccess(access, doc.creator, AccessLevel.ACCESS, 'DocStoreService.UpdateDocument')
                        "
                        class="flex-1"
                        :text="$t('common.edit')"
                        :shortcuts="['D', 'E']"
                    >
                        <UButton
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
                        v-if="can('DocStoreService.ListDocumentReqs')"
                        class="flex-1"
                        :text="$t('common.request', 2)"
                        :shortcuts="['D', 'R']"
                    >
                        <UButton class="flex-1" block icon="i-mdi-frequently-asked-questions" @click="openRequestsModal">
                            {{ $t('common.request', 2) }}
                        </UButton>
                    </UTooltip>

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
                        {{ doc.title }}
                    </h1>
                </div>

                <div class="mb-2 flex gap-2">
                    <UBadge v-if="doc.category" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-shape" class="size-5" />
                        <span :title="doc.category.description ?? $t('common.na')">
                            {{ doc.category.name }}
                        </span>
                    </UBadge>

                    <OpenClosedBadge :closed="doc.closed" />

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
                <div class="contentView mx-auto max-w-screen-xl break-words rounded-lg bg-base-900">
                    <!-- eslint-disable vue/no-v-html -->
                    <div ref="contentRef" class="prose prose-invert min-w-full px-4 py-2" v-html="doc.content"></div>
                </div>
            </div>

            <template #footer>
                <UAccordion multiple :items="accordionItems" :unmount="true">
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

                            <div v-else class="flex flex-col gap-2">
                                <div class="flex flex-row flex-wrap gap-1">
                                    <UBadge
                                        v-for="entry in access?.jobs"
                                        :key="entry.id"
                                        color="black"
                                        class="inline-flex gap-1"
                                        size="md"
                                    >
                                        <span class="size-2 rounded-full bg-info-500" />
                                        <span>
                                            {{ entry.jobLabel
                                            }}<span
                                                v-if="entry.minimumGrade > 0"
                                                :title="`${entry.jobLabel} - ${$t('common.rank')} ${entry.minimumGrade}`"
                                            >
                                                ({{ entry.jobGradeLabel }})</span
                                            >
                                            -
                                            {{ $t(`enums.docstore.AccessLevel.${AccessLevel[entry.access]}`) }}
                                        </span>
                                    </UBadge>
                                </div>

                                <div class="flex flex-row flex-wrap gap-1">
                                    <UBadge
                                        v-for="entry in access?.users"
                                        :key="entry.id"
                                        color="black"
                                        class="inline-flex gap-1"
                                        size="md"
                                    >
                                        <span class="size-2 rounded-full bg-amber-500" />
                                        <span :title="`${$t('common.id')} ${entry.userId}`">
                                            {{ entry.user?.firstname }}
                                            {{ entry.user?.lastname }} -
                                            {{ $t(`enums.docstore.AccessLevel.${AccessLevel[entry.access]}`) }}
                                        </span>
                                    </UBadge>
                                </div>
                            </div>
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

                    <template v-if="can('DocStoreService.ListDocumentActivity')" #activity>
                        <UContainer>
                            <DocumentActivityList :document-id="documentId" />
                        </UContainer>
                    </template>
                </UAccordion>
            </template>
        </UCard>
    </template>
</template>

<style scoped>
.contentView:deep(.prose) {
    * {
        margin-top: 4px;
        margin-bottom: 4px;
    }

    input[type='checkbox']:checked {
        opacity: 1;
    }
}
</style>
