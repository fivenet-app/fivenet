<script lang="ts" setup>
import type { UForm } from '#components';
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { checkDocAccess, logger } from '~/components/documents/helpers';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { useClipboardStore } from '~/stores/clipboard';
import { useCompletorStore } from '~/stores/completor';
import type { Content } from '~/types/history';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import { type DocumentJobAccess, type DocumentUserAccess, AccessLevel } from '~~/gen/ts/resources/documents/access';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { DocumentReference, DocumentRelation } from '~~/gen/ts/resources/documents/documents';
import type { File } from '~~/gen/ts/resources/file/file';
import { ObjectType } from '~~/gen/ts/resources/notifications/client_view';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UpdateDocumentRequest } from '~~/gen/ts/services/documents/documents';
import { jobAccessEntry, userAccessEntry } from '~~/shared/types/validation';
import ConfirmModal from '../partials/ConfirmModal.vue';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';
import CategoryBadge from '../partials/documents/CategoryBadge.vue';
import SelectMenu from '../partials/SelectMenu.vue';
import DocumentReferenceManager from './DocumentReferenceManager.vue';
import DocumentRelationManager from './DocumentRelationManager.vue';

const props = defineProps<{
    documentId: number;
}>();

const { t } = useI18n();

const { can } = useAuth();

const overlay = useOverlay();

const clipboardStore = useClipboardStore();

const completorStore = useCompletorStore();

const notifications = useNotificationsStore();

const historyStore = useHistoryStore();

const documentsDocuments = await useDocumentsDocuments();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const {
    data: document,
    status,
    error,
    refresh,
} = useLazyAsyncData(`documents-${props.documentId}-editor`, () => documentsDocuments.getDocument(props.documentId));

useHead({
    title: () =>
        document.value?.document?.title
            ? `${document.value.document.title} - ${t('pages.documents.edit.title')}`
            : t('pages.documents.edit.title'),
});

const { maxAccessEntries } = useAppConfig();

const { ydoc, provider } = await useCollabDoc('documents', props.documentId);

function setFromProps(): void {
    if (!document.value?.document) return;

    state.title = document.value.document.title;
    state.state = document.value.document.state;
    state.content = document.value.document.content?.rawContent ?? '';
    state.category = document.value.document.category ?? emptyCategory;
    state.closed = document.value.document.closed;
    state.draft = document.value.document.draft;
    state.public = document.value.document.public;
    if (document.value.access) {
        state.access.jobs = document.value.access.jobs;
        state.access.users = document.value.access.users;
    }
    state.files = document.value.document.files;
}

const onSync = (s: boolean) => {
    if (!s) return;
    logger.debug('DocumentEditor - Sync received, setting state from props', s);

    if (ydoc.getXmlFragment('content').length === 0) {
        logger.info('DocumentEditor - Content is empty, setting from props');
        // If the content is empty, we need to set it from the props
        setFromProps();
    }
    provider.off('sync', onSync);
};
provider.on('sync', onSync);

watch(document, async () => {
    const [refs, rels] = await Promise.all([
        documentsDocumentsClient.getDocumentReferences({
            documentId: props.documentId,
        }),
        documentsDocumentsClient.getDocumentRelations({
            documentId: props.documentId,
        }),
    ]);
    state.references = refs.response.references;
    state.relations = rels.response.relations;
});

const route = useRoute();

const emptyCategory: Category = {
    id: 0,
    name: t('common.categories', 0),
};

const schema = z.object({
    title: z.string().min(3).max(255),
    state: z.union([z.string().length(0), z.string().min(3).max(32)]),
    content: z.string().min(3).max(1750000),
    closed: z.coerce.boolean(),
    draft: z.coerce.boolean(),
    public: z.coerce.boolean(),
    category: z.custom<Category>().default({ ...emptyCategory }),
    access: z
        .object({
            jobs: jobAccessEntry.array().max(maxAccessEntries).default([]),
            users: userAccessEntry.array().max(maxAccessEntries).default([]),
        })
        .default({ jobs: [], users: [] }),
    files: z.custom<File>().array().max(5).default([]),
    references: z.custom<DocumentReference>().array().max(15).default([]),
    relations: z.custom<DocumentRelation>().array().max(15).default([]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    title: '',
    state: '',
    content: '',
    closed: false,
    draft: false,
    public: false,
    category: emptyCategory,
    access: {
        jobs: [],
        users: [],
    },
    files: [],
    references: [],
    relations: [],
});

const changed = ref(false);
const saving = ref(false);

// Track last saved string and timestamp
let lastSavedString = '';
let lastSaveTimestamp = 0;

async function saveHistory(values: Schema, name: string | undefined = undefined, type = 'document'): Promise<void> {
    if (saving.value) {
        return;
    }

    const now = Date.now();
    // Skip if identical to last saved or if within MIN_GAP
    if (state.content === lastSavedString || now - lastSaveTimestamp < 5000) {
        return;
    }

    saving.value = true;

    historyStore.addVersion<Content>(
        type,
        props.documentId,
        {
            content: values.content,
            files: values.files,
        },
        name,
    );

    useTimeoutFn(() => {
        saving.value = false;
    }, 1750);

    lastSavedString = state.content;
    lastSaveTimestamp = now;
}

historyStore.handleRefresh(() => saveHistory(state, 'document'));

watchDebounced(
    state,
    () => {
        if (changed.value) {
            saveHistory(state);
        } else {
            changed.value = true;
        }
    },
    {
        debounce: 1_000,
        maxWait: 2_500,
    },
);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    await updateDocument(props.documentId, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

async function updateDocument(id: number, values: Schema): Promise<void> {
    values.access.users.forEach((user) => {
        if (user.id < 0) user.id = 0;
        user.user = undefined; // Clear user object to avoid sending unnecessary data
    });
    values.access.jobs.forEach((job) => job.id < 0 && (job.id = 0));

    const req: UpdateDocumentRequest = {
        documentId: id,
        title: values.title,
        content: {
            rawContent: values.content,
        },
        contentType: ContentType.HTML,
        state: values.state,
        closed: values.closed,
        draft: values.draft,
        public: values.public,
        categoryId: values.category?.id !== 0 ? values.category?.id : undefined,
        access: values.access,
        files: values.files,
    };

    try {
        const call = documentsDocumentsClient.updateDocument(req);
        const { response } = await call;

        if (canDo.value.references) {
            // Remove references that are no longer present
            const referencesToRemove: number[] = [];
            state.references
                .filter((r) => r.id !== undefined && r.id > 0)
                .forEach((ref) => {
                    if (!state.references.some((r) => r.id === ref.id)) {
                        referencesToRemove.push(ref.id!);
                    }
                });
            referencesToRemove.forEach((id) => {
                documentsDocumentsClient.removeDocumentReference({
                    id: id,
                });
            });
            // Add new references
            state.references
                .filter((r) => r.id === undefined || r.id <= 0)
                .forEach((ref) => {
                    if (state.references.some((r) => r.id === ref.id)) {
                        return;
                    }
                    ref.sourceDocumentId = response.document!.id!;
                    documentsDocumentsClient.addDocumentReference({
                        reference: ref,
                    });
                });
        }

        if (canDo.value.relations) {
            // Remove relations that are no longer present
            const relationsToRemove: number[] = [];
            state.relations
                .filter((r) => r.id !== undefined && r.id > 0)
                .forEach((rel) => {
                    if (!state.relations.some((r) => r.id === rel.id)) {
                        relationsToRemove.push(rel.id!);
                    }
                });
            relationsToRemove.forEach((id) => {
                documentsDocumentsClient.removeDocumentRelation({ id });
            });
            // Add new relations
            state.relations
                .filter((r) => r.id === undefined || r.id <= 0)
                .forEach((rel) => {
                    rel.documentId = response.document!.id!;
                    documentsDocumentsClient.addDocumentRelation({
                        relation: rel,
                    });
                });
        }

        notifications.add({
            title: { key: 'notifications.document_updated.title', parameters: {} },
            description: { key: 'notifications.document_updated.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        if (response.document) {
            state.draft = response.document.draft;
        }

        clipboardStore.clear();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const items = [
    {
        slot: 'content' as const,
        label: t('common.content'),
        icon: 'i-mdi-pencil',
        value: 'content',
    },
    {
        slot: 'access' as const,
        label: t('common.access', 1),
        icon: 'i-mdi-key',
        value: 'access',
    },
    {
        slot: 'references' as const,
        label: t('common.reference', 2),
        icon: 'i-mdi-file-document',
        value: 'references',
    },
    {
        slot: 'relations' as const,
        label: t('common.relation', 2),
        icon: 'i-mdi-account-multiple',
        value: 'relations',
    },
];

const router = useRouter();

const selectedTab = computed({
    get() {
        return (route.query.tab as string) || 'content';
    },
    set(tab) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.push({ query: { tab: tab }, hash: '#control-active-item' });
    },
});

const canDo = computed(() => ({
    edit: checkDocAccess(
        state.access,
        document.value?.document?.creator,
        AccessLevel.EDIT,
        'documents.DocumentsService/UpdateDocument',
        document.value?.document?.creatorJob,
    ),
    access: checkDocAccess(
        state.access,
        document.value?.document?.creator,
        AccessLevel.ACCESS,
        'documents.DocumentsService/UpdateDocument',
        document.value?.document?.creatorJob,
    ),
    references: can('documents.DocumentsService/AddDocumentReference').value,
    relations: can('documents.DocumentsService/AddDocumentRelation').value,
}));

// Handle the client update event
const { sendClientView } = useClientUpdate(ObjectType.DOCUMENT, () =>
    notifications.add({
        title: { key: 'notifications.documents.client_view_update.title', parameters: {} },
        description: { key: 'notifications.documents.client_view_update.content', parameters: {} },
        duration: 7500,
        type: NotificationType.INFO,
        actions: [
            {
                label: { key: 'common.refresh', parameters: {} },
                icon: 'i-mdi-refresh',
                onClick: () => refresh(),
            },
        ],
    }),
);
sendClientView(props.documentId);

useYText(ydoc.getText('title'), toRef(state, 'title'), { provider: provider });
useYText(ydoc.getText('state'), toRef(state, 'state'), { provider: provider });
const detailsYdoc = ydoc.getMap('details');
useYBoolean(detailsYdoc, 'draft', toRef(state, 'draft'), { provider: provider });
useYBoolean(detailsYdoc, 'closed', toRef(state, 'closed'), { provider: provider });
useYBoolean(detailsYdoc, 'draft', toRef(state, 'draft'), { provider: provider });
useYBoolean(detailsYdoc, 'public', toRef(state, 'public'), { provider: provider });
const categoryYdoc = ydoc.getMap<Primitive>('category');
useYObject<Category>(
    categoryYdoc,
    toRef(state, 'category'),
    {
        omit: ['createdAt', 'deletedAt'],
    },
    {
        provider: provider,
    },
);

// Access
useYArrayFiltered<DocumentJobAccess>(
    ydoc.getArray('access_jobs'),
    toRef(state.access, 'jobs'),
    { omit: ['createdAt', 'user'] },
    { provider: provider },
);
useYArrayFiltered<DocumentUserAccess>(
    ydoc.getArray('access_users'),
    toRef(state.access, 'users'),
    {
        omit: ['createdAt', 'user'],
    },
    { provider: provider },
);

// Files
useYArrayFiltered<File>(
    ydoc.getArray('files'),
    toRef(state, 'files'),
    {
        omit: ['createdAt', 'meta'],
    },
    { provider: provider },
);

// References and Relations
useYArrayFiltered<DocumentReference>(
    ydoc.getArray('doc_references'),
    toRef(state, 'references'),
    {
        omit: ['createdAt', 'sourceDocument', 'targetUser', 'deletedAt', 'updatedAt'],
    },
    { provider: provider },
);
useYArrayFiltered<DocumentRelation>(
    ydoc.getArray('doc_relations'),
    toRef(state, 'relations'),
    {
        omit: ['createdAt', 'document', 'sourceUser'],
    },
    { provider: provider },
);

logger.info(
    'Editor - Can Do: Edit',
    canDo.value.edit,
    'Access',
    canDo.value.access,
    'References',
    canDo.value.references,
    'Relations',
    canDo.value.relations,
);

const confirmModal = overlay.create(ConfirmModal);

const formRef = useTemplateRef('formRef');

provide('yjsDoc', ydoc);
provide('yjsProvider', provider);
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0 overflow-y-hidden' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.documents.edit.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton :fallback-to="{ name: 'documents-id', params: { id: documentId } }" />

                    <UButton
                        trailing-icon="i-mdi-content-save"
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                        @click="() => formRef?.submit()"
                    >
                        <span class="hidden truncate sm:block">
                            {{ $t('common.save') }}
                        </span>
                    </UButton>

                    <UButton
                        v-if="state.draft"
                        color="info"
                        trailing-icon="i-mdi-publish"
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                        @click="
                            confirmModal.open({
                                title: $t('common.publish_confirm.title', { type: $t('common.document', 1) }),
                                description: $t('common.publish_confirm.description'),
                                color: 'info',
                                iconClass: 'text-info-500 dark:text-info-400',
                                icon: 'i-mdi-publish',
                                confirm: () => {
                                    state.draft = false;
                                    formRef?.submit();
                                },
                            })
                        "
                    >
                        <span class="hidden truncate sm:block">
                            {{ $t('common.publish') }}
                        </span>
                    </UButton>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <UForm
                ref="formRef"
                class="flex min-h-0 w-full flex-1 flex-col overflow-y-hidden"
                :schema="schema"
                :state="state"
                @submit="onSubmitThrottle"
            >
                <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.document', 1)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.page', 1)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock
                    v-else-if="!document"
                    icon="i-mdi-file-search"
                    :message="$t('common.not_found', [$t('common.page', 1)])"
                />

                <UTabs
                    v-else
                    v-model="selectedTab"
                    class="flex-1 flex-col overflow-y-hidden"
                    :items="items"
                    variant="link"
                    :unmount-on-hide="false"
                    :ui="{ list: 'mx-auto', content: 'flex flex-col flex-1 min-h-0 max-h-full overflow-y-hidden' }"
                >
                    <template #content>
                        <UDashboardToolbar>
                            <template #default>
                                <div class="mx-auto mb-2 flex w-full max-w-(--breakpoint-xl) flex-col gap-2">
                                    <UFormField name="title" :label="$t('common.title')" required>
                                        <UInput
                                            v-model="state.title"
                                            type="text"
                                            size="xl"
                                            class="w-full"
                                            :placeholder="$t('common.title')"
                                            :disabled="!canDo.edit"
                                        />
                                    </UFormField>

                                    <div class="flex flex-row gap-2">
                                        <UFormField class="flex-1" name="category" :label="$t('common.category', 1)">
                                            <SelectMenu
                                                v-model="state.category"
                                                :filter-fields="['name']"
                                                block
                                                nullable
                                                :disabled="!canDo.edit"
                                                class="w-full"
                                                :searchable="
                                                    async (q: string) => {
                                                        try {
                                                            const categories =
                                                                await completorStore.completeDocumentCategories(q);
                                                            if (!categories.find((c) => c.id === state.category.id)) {
                                                                categories.unshift(state.category);
                                                            }
                                                            if (!categories.find((c) => c.id === emptyCategory.id)) {
                                                                categories.unshift(emptyCategory);
                                                            }
                                                            return categories;
                                                        } catch (e) {
                                                            handleGRPCError(e as RpcError);
                                                            throw e;
                                                        }
                                                    }
                                                "
                                                searchable-key="completor-document-categories"
                                                :search-input="{ placeholder: $t('common.search_field') }"
                                            >
                                                <template
                                                    v-if="state.category && state.category.id !== emptyCategory.id"
                                                    #default
                                                >
                                                    <CategoryBadge :category="state.category" />
                                                </template>

                                                <template #item="{ item }">
                                                    <CategoryBadge :category="item" />
                                                </template>

                                                <template #empty>
                                                    {{ $t('common.not_found', [$t('common.category', 2)]) }}
                                                </template>
                                            </SelectMenu>
                                        </UFormField>

                                        <UFormField class="flex-1" name="state" :label="$t('common.state')">
                                            <UInput
                                                v-model="state.state"
                                                type="text"
                                                class="w-full"
                                                :placeholder="`${$t('common.document', 1)} ${$t('common.state')}`"
                                                :disabled="!canDo.edit"
                                            />
                                        </UFormField>

                                        <UFormField class="flex-initial" name="closed" :label="`${$t('common.close', 2)}?`">
                                            <USwitch v-model="state.closed" :disabled="!canDo.edit" />
                                        </UFormField>
                                    </div>
                                </div>
                            </template>
                        </UDashboardToolbar>

                        <ClientOnly>
                            <TiptapEditor
                                v-model="state.content"
                                v-model:files="state.files"
                                class="m-2 mx-auto w-full max-w-(--breakpoint-xl) flex-1"
                                :disabled="!canDo.edit"
                                history-type="document"
                                :saving="saving"
                                enable-collab
                                :target-id="document.document?.id"
                                filestore-namespace="documents"
                                :filestore-service="(opts) => documentsDocumentsClient.uploadFile(opts)"
                            />
                        </ClientOnly>
                    </template>

                    <template #access>
                        <UContainer class="p-4 sm:p-4">
                            <UPageCard :title="$t('common.access')">
                                <UFormField name="access">
                                    <AccessManager
                                        v-model:jobs="state.access.jobs"
                                        v-model:users="state.access.users"
                                        :disabled="!canDo.access"
                                        :target-id="documentId ?? 0"
                                        :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.documents.AccessLevel')"
                                        name="access"
                                    />
                                </UFormField>
                            </UPageCard>
                        </UContainer>
                    </template>

                    <template #references>
                        <UContainer class="p-4 sm:p-4">
                            <UPageCard :title="`${$t('common.document', 1)} ${$t('common.reference', 2)}`">
                                <DocumentReferenceManager v-model="state.references" :document-id="documentId" />
                            </UPageCard>
                        </UContainer>
                    </template>

                    <template #relations>
                        <UContainer class="p-4 sm:p-4">
                            <UPageCard :title="`${$t('common.document', 1)} ${$t('common.relation', 2)}`">
                                <DocumentRelationManager v-model="state.relations" :document-id="documentId" />
                            </UPageCard>
                        </UContainer>
                    </template>
                </UTabs>
            </UForm>
        </template>
    </UDashboardPanel>
</template>
