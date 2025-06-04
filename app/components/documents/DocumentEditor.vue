<script lang="ts" setup>
import type { UForm } from '#components';
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import DocumentReferenceManagerModal from '~/components/documents/DocumentReferenceManagerModal.vue';
import DocumentRelationManagerModal from '~/components/documents/DocumentRelationManagerModal.vue';
import { checkDocAccess, logger } from '~/components/documents/helpers';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { availableIcons, fallbackIcon } from '~/components/partials/icons';
import { useClipboardStore } from '~/stores/clipboard';
import { useCompletorStore } from '~/stores/completor';
import { useNotificatorStore } from '~/stores/notificator';
import type { DocumentContent, DocumentMeta } from '~/types/history';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { DocumentJobAccess, DocumentUserAccess } from '~~/gen/ts/resources/documents/access';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { DocumentReference, DocumentRelation } from '~~/gen/ts/resources/documents/documents';
import type { File } from '~~/gen/ts/resources/file/file';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { GetDocumentResponse, UpdateDocumentRequest } from '~~/gen/ts/services/documents/documents';
import ConfirmModal from '../partials/ConfirmModal.vue';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';

const props = defineProps<{
    documentId: number;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { can } = useAuth();

const modal = useModal();

const clipboardStore = useClipboardStore();

const completorStore = useCompletorStore();

const notifications = useNotificatorStore();

const historyStore = useHistoryStore();

const {
    data: document,
    pending: loading,
    error,
    refresh,
} = useLazyAsyncData(`documents-${props.documentId}-editor`, () => getDocument(props.documentId));

async function getDocument(id: number): Promise<GetDocumentResponse> {
    try {
        const call = $grpc.documents.documents.getDocument({
            documentId: id,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);

        await navigateTo({ name: 'documents' });
        throw e;
    }
}

const { maxAccessEntries } = useAppConfig();

const { ydoc, provider } = useCollabDoc('documents', props.documentId);

watchOnce(document, () => provider.connect());

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
    state.files = document.value.document.files ?? [];
}
provider.once('loadContent', () => setFromProps());

watch(document, async () => {
    const [refs, rels] = await Promise.all([
        $grpc.documents.documents.getDocumentReferences({
            documentId: props.documentId,
        }),
        $grpc.documents.documents.getDocumentRelations({
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
    closed: z.boolean(),
    draft: z.boolean(),
    public: z.boolean(),
    category: z.custom<Category>(),
    access: z.object({
        jobs: z.custom<DocumentJobAccess>().array().max(maxAccessEntries),
        users: z.custom<DocumentUserAccess>().array().max(maxAccessEntries),
    }),
    files: z.custom<File>().array().max(5),
    references: z.custom<DocumentReference>().array().max(15),
    relations: z.custom<DocumentRelation>().array().max(15),
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

async function saveHistory(values: Schema, type = 'document'): Promise<void> {
    if (saving.value) {
        return;
    }
    saving.value = true;

    historyStore.addVersion<DocumentContent, DocumentMeta>(
        type,
        props.documentId,
        {
            content: values.content,
            files: values.files,
        },
        {
            title: values.title,
            category: values.category,
            access: values.access,
            closed: values.closed,
            state: values.state,
            references: values.references,
            relations: values.relations,
        },
    );

    useTimeoutFn(() => {
        saving.value = false;
    }, 1750);
}

watchDebounced(
    state,
    async () => {
        if (changed.value) {
            saveHistory(state);
        } else {
            changed.value = true;
        }
    },
    {
        debounce: 750,
        maxWait: 2500,
    },
);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    await updateDocument(props.documentId, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

async function updateDocument(id: number, values: Schema): Promise<void> {
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
        const call = $grpc.documents.documents.updateDocument(req);
        const { response } = await call;

        if (canDo.value.references) {
            // Remove references that are no longer present
            const referencesToRemove: number[] = [];
            state.references.forEach((ref) => {
                if (!state.references.some((r) => r.id === ref.id)) {
                    referencesToRemove.push(ref.id!);
                }
            });
            referencesToRemove.forEach((id) => {
                $grpc.documents.documents.removeDocumentReference({
                    id: id,
                });
            });
            // Add new references
            state.references.forEach((ref) => {
                if (state.references.some((r) => r.id === ref.id)) {
                    return;
                }
                ref.sourceDocumentId = response.document!.id!;
                $grpc.documents.documents.addDocumentReference({
                    reference: ref,
                });
            });
        }

        if (canDo.value.relations) {
            // Remove relations that are no longer present
            const relationsToRemove: number[] = [];
            state.relations.forEach((rel) => {
                if (!state.relations.some((r) => r.id === rel.id)) {
                    relationsToRemove.push(rel.id!);
                }
            });
            relationsToRemove.forEach((id) => {
                $grpc.documents.documents.removeDocumentRelation({ id });
            });
            // Add new relations
            state.relations.forEach((rel) => {
                if (state.relations.some((r) => r.id === rel.id)) {
                    return;
                }
                rel.documentId = response.document!.id;
                $grpc.documents.documents.addDocumentRelation({
                    relation: rel,
                });
            });
        }

        notifications.add({
            title: { key: 'notifications.document_updated.title', parameters: {} },
            description: { key: 'notifications.document_updated.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
        clipboardStore.clear();

        await navigateTo({
            name: 'documents-id',
            params: { id: response.document!.id },
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const items = [
    {
        slot: 'content',
        label: t('common.content'),
        icon: 'i-mdi-pencil',
    },
    {
        slot: 'access',
        label: t('common.access', 1),
        icon: 'i-mdi-key',
    },
];

const router = useRouter();

const selectedTab = computed({
    get() {
        const index = items.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items[value]?.slot }, hash: '#' });
    },
});

const categoriesLoading = ref(false);

const canDo = computed(() => ({
    edit:
        props.documentId === undefined
            ? true
            : checkDocAccess(
                  state.access,
                  document.value?.document?.creator,
                  AccessLevel.EDIT,
                  'documents.DocumentsService.UpdateDocument',
              ),
    access:
        props.documentId === undefined
            ? true
            : checkDocAccess(
                  state.access,
                  document.value?.document?.creator,
                  AccessLevel.ACCESS,
                  'documents.DocumentsService.UpdateDocument',
              ),
    references: can('documents.DocumentsService.AddDocumentReference').value,
    relations: can('documents.DocumentsService.AddDocumentRelation').value,
}));

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

useYText(ydoc.getText('title'), toRef(state, 'title'), { provider: provider });
useYText(ydoc.getText('state'), toRef(state, 'state'), { provider: provider });
const detailsYdoc = ydoc.getMap('details');
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

provide('yjsDoc', ydoc);
provide('yjsProvider', provider);

const formRef = useTemplateRef<typeof UForm>('formRef');
</script>

<template>
    <UForm
        ref="formRef"
        class="min-h-dscreen flex w-full max-w-full flex-1 flex-col overflow-y-auto"
        :schema="schema"
        :state="state"
        @submit="onSubmitThrottle"
    >
        <UDashboardNavbar :title="$t('pages.documents.edit.title')">
            <template #right>
                <PartialsBackButton :fallback-to="{ name: 'documents-id', params: { id: documentId } }" />

                <UButton type="submit" trailing-icon="i-mdi-content-save" :disabled="!canSubmit" :loading="!canSubmit">
                    <span class="hidden truncate sm:block">
                        {{ $t('common.save') }}
                    </span>
                </UButton>

                <UButton
                    v-if="document?.document?.draft"
                    type="submit"
                    color="info"
                    trailing-icon="i-mdi-publish"
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    @click.prevent="
                        modal.open(ConfirmModal, {
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

        <UDashboardPanelContent class="p-0 sm:pb-0">
            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.page', 1)])" />
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
                class="flex flex-1 flex-col"
                :items="items"
                :ui="{
                    wrapper: 'space-y-0 overflow-y-hidden',
                    container: 'flex flex-1 flex-col overflow-y-hidden',
                    base: 'flex flex-1 flex-col overflow-y-hidden',
                    list: { rounded: '' },
                }"
            >
                <template #content>
                    <UDashboardToolbar>
                        <template #default>
                            <div class="flex w-full flex-col gap-2">
                                <UFormGroup name="title" :label="$t('common.title')" required>
                                    <UInput
                                        v-model="state.title"
                                        type="text"
                                        size="xl"
                                        :placeholder="$t('common.title')"
                                        :disabled="!canDo.edit"
                                    />
                                </UFormGroup>

                                <div class="flex flex-row gap-2">
                                    <UFormGroup class="flex-1" name="category" :label="$t('common.category', 1)">
                                        <ClientOnly>
                                            <USelectMenu
                                                v-model="state.category"
                                                option-attribute="name"
                                                :search-attributes="['name']"
                                                block
                                                nullable
                                                :searchable="
                                                    async (search: string) => {
                                                        try {
                                                            categoriesLoading = true;
                                                            const categories =
                                                                await completorStore.completeDocumentCategories(search);
                                                            categoriesLoading = false;
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
                                                        } finally {
                                                            categoriesLoading = false;
                                                        }
                                                    }
                                                "
                                                searchable-lazy
                                                :searchable-placeholder="$t('common.search_field')"
                                            >
                                                <template #label>
                                                    <span
                                                        v-if="state.category"
                                                        class="inline-flex gap-1"
                                                        :class="`bg-${state.category.color}-500`"
                                                    >
                                                        <component
                                                            :is="
                                                                availableIcons.find(
                                                                    (item) => item.name === state.category?.icon,
                                                                ) ?? fallbackIcon
                                                            "
                                                            v-if="state.category.icon"
                                                            class="size-5"
                                                        />
                                                        <span class="truncate">{{ state.category.name }}</span>
                                                    </span>
                                                    <span v-else> &nbsp; </span>
                                                </template>

                                                <template #option="{ option }">
                                                    <span class="inline-flex gap-1" :class="`bg-${option.color}-500`">
                                                        <component
                                                            :is="
                                                                availableIcons.find((item) => item.name === option.icon) ??
                                                                fallbackIcon
                                                            "
                                                            v-if="option.icon"
                                                            class="size-5"
                                                        />
                                                        <span class="truncate">{{ option.name }}</span>
                                                    </span>
                                                </template>

                                                <template #option-empty="{ query: search }">
                                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                                </template>

                                                <template #empty>
                                                    {{ $t('common.not_found', [$t('common.category', 2)]) }}
                                                </template>
                                            </USelectMenu>
                                        </ClientOnly>
                                    </UFormGroup>

                                    <UFormGroup class="flex-1" name="state" :label="$t('common.state')">
                                        <UInput
                                            v-model="state.state"
                                            type="text"
                                            :placeholder="`${$t('common.document', 1)} ${$t('common.state')}`"
                                            :disabled="!canDo.edit"
                                        />
                                    </UFormGroup>

                                    <UFormGroup class="flex-initial" name="closed" :label="`${$t('common.close', 2)}?`">
                                        <UToggle v-model="state.closed" :disabled="!canDo.edit" />
                                    </UFormGroup>
                                </div>
                            </div>
                        </template>
                    </UDashboardToolbar>

                    <UFormGroup
                        v-if="canDo.edit"
                        class="flex flex-1 overflow-y-hidden"
                        name="content"
                        :ui="{ container: 'flex flex-1 flex-col mt-0 overflow-y-hidden', label: { wrapper: 'hidden' } }"
                        label="&nbsp;"
                    >
                        <ClientOnly>
                            <TiptapEditor
                                v-model="state.content"
                                v-model:files="state.files"
                                class="mx-auto w-full max-w-screen-xl flex-1 overflow-y-hidden"
                                :disabled="!canDo.edit"
                                rounded="rounded-none"
                                :target-id="document.document?.id"
                                filestore-namespace="documents"
                                :filestore-service="(opts) => $grpc.documents.documents.uploadFile(opts)"
                                @file-uploaded="(file) => state.files.push(file)"
                            >
                                <template #footer>
                                    <div v-if="saving" class="place-self-start">
                                        <UIcon class="h-4 w-4 animate-spin" name="i-mdi-content-save" />
                                        <span>{{ $t('common.save', 2) }}...</span>
                                    </div>
                                </template>
                            </TiptapEditor>
                        </ClientOnly>
                    </UFormGroup>

                    <UDashboardToolbar
                        class="flex shrink-0 justify-between border-b-0 border-t border-gray-200 px-3 py-3.5 dark:border-gray-700"
                    >
                        <UButtonGroup v-if="canDo.edit" class="inline-flex w-full">
                            <UButton
                                v-if="canDo.relations"
                                class="flex-1"
                                block
                                :disabled="!canDo.edit"
                                icon="i-mdi-account-multiple"
                                @click="
                                    modal.open(DocumentRelationManagerModal, {
                                        relations: state.relations,
                                        documentId: documentId,
                                        'onUpdate:relations': (value) => (state.relations = value),
                                    })
                                "
                            >
                                {{ $t('common.citizen', 1) }} {{ $t('common.relation', 2) }}
                            </UButton>

                            <UButton
                                v-if="canDo.references"
                                class="flex-1"
                                block
                                :disabled="!canDo.edit"
                                icon="i-mdi-file-document"
                                @click="
                                    modal.open(DocumentReferenceManagerModal, {
                                        references: state.references,
                                        documentId: documentId,
                                        'onUpdate:references': (value) => (state.references = value),
                                    })
                                "
                            >
                                {{ $t('common.document', 1) }} {{ $t('common.reference', 2) }}
                            </UButton>
                        </UButtonGroup>
                    </UDashboardToolbar>
                </template>

                <template #access>
                    <div class="flex flex-1 flex-col gap-2 overflow-y-scroll px-2">
                        <h2 class="text- text-gray-900 dark:text-white">
                            {{ $t('common.access') }}
                        </h2>

                        <AccessManager
                            v-model:jobs="state.access.jobs"
                            v-model:users="state.access.users"
                            :target-id="documentId ?? 0"
                            :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.documents.AccessLevel')"
                            :disabled="!canDo.access"
                        />
                    </div>
                </template>
            </UTabs>
        </UDashboardPanelContent>
    </UForm>
</template>
