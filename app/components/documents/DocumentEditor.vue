<script lang="ts" setup>
import type { UForm } from '#components';
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import DocumentReferenceManager from '~/components/documents/DocumentReferenceManager.vue';
import DocumentRelationManager from '~/components/documents/DocumentRelationManager.vue';
import { checkDocAccess, logger } from '~/components/documents/helpers';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { availableIcons, fallbackIcon } from '~/components/partials/icons';
import { getDocument as clipboardGetDocument, getUser, useClipboardStore } from '~/stores/clipboard';
import { useCompletorStore } from '~/stores/completor';
import { useDocumentEditorStore } from '~/stores/documenteditor';
import { useNotificatorStore } from '~/stores/notificator';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { DocumentJobAccess, DocumentUserAccess } from '~~/gen/ts/resources/documents/access';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { DocumentReference, DocumentRelation } from '~~/gen/ts/resources/documents/documents';
import { DocReference, DocRelation } from '~~/gen/ts/resources/documents/documents';
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

const { can, activeChar } = useAuth();

const modal = useModal();

const clipboardStore = useClipboardStore();

const completorStore = useCompletorStore();

const documentStore = useDocumentEditorStore();

const notifications = useNotificatorStore();

const {
    data: document,
    pending: loading,
    error,
    refresh,
} = useLazyAsyncData(`documents-${props.documentId}`, () => getDocument(props.documentId));

async function getDocument(id: number): Promise<GetDocumentResponse> {
    try {
        const req = { documentId: id };
        const call = $grpc.documents.documents.getDocument(req);
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

async function setFromProps(): Promise<void> {
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

    const refs = await $grpc.documents.documents.getDocumentReferences({
        documentId: props.documentId,
    });
    currentReferences.value = refs.response.references;
    const rels = await $grpc.documents.documents.getDocumentRelations({
        documentId: props.documentId,
    });
    currentRelations.value = rels.response.relations;
}
provider.once('loadContent', async () => await setFromProps());

const route = useRoute();

const canEdit = ref(false);

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
});

const openRelationManager = ref<boolean>(false);
const relationManagerData = ref(new Map<number, DocumentRelation>());
const currentRelations = ref<Readonly<DocumentRelation>[]>([]);
watch(currentRelations, () => currentRelations.value.forEach((e) => relationManagerData.value.set(e.id!, e)));

const openReferenceManager = ref<boolean>(false);
const referenceManagerData = ref(new Map<number, DocumentReference>());
const currentReferences = ref<Readonly<DocumentReference>[]>([]);
watch(currentReferences, () => currentReferences.value.forEach((e) => referenceManagerData.value.set(e.id!, e)));

onMounted(async () => {
    clipboardStore.activeStack.documents.forEach((doc, idx) => {
        referenceManagerData.value.set(idx, {
            id: idx,
            sourceDocumentId: props.documentId ?? 0,
            targetDocumentId: doc.id!,
            targetDocument: clipboardGetDocument(doc),
            creatorId: activeChar.value!.userId,
            creator: activeChar.value!,
            reference: DocReference.SOLVES,
        });
    });

    clipboardStore.activeStack.users.forEach((user, idx) => {
        relationManagerData.value.set(idx, {
            id: idx,
            documentId: props.documentId ?? 0,
            targetUserId: user.userId!,
            targetUser: getUser(user),
            sourceUserId: activeChar.value!.userId,
            sourceUser: activeChar.value!,
            relation: DocRelation.CAUSED,
        });
    });

    canEdit.value = true;
});

const saving = ref(false);

async function saveToStore(values: Schema): Promise<void> {
    if (saving.value) {
        return;
    }
    saving.value = true;

    documentStore.save({
        title: values.title,
        content: values.content,
        state: values.state,
        closed: values.closed,
        category: values.category,
    });

    useTimeoutFn(() => {
        saving.value = false;
    }, 1250);
}

const changed = ref(false);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    await updateDocument(props.documentId, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

watchDebounced(
    state,
    async () => {
        if (changed.value) {
            saveToStore(state);
        } else {
            changed.value = true;
        }
    },
    {
        debounce: 750,
        maxWait: 2500,
    },
);

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
            const referencesToRemove: number[] = [];
            currentReferences.value.forEach((ref) => {
                if (!referenceManagerData.value.has(ref.id!)) {
                    referencesToRemove.push(ref.id!);
                }
            });
            referencesToRemove.forEach((id) => {
                $grpc.documents.documents.removeDocumentReference({
                    id: id,
                });
            });
            referenceManagerData.value.forEach((ref) => {
                if (currentReferences.value.find((r) => r.id === ref.id!)) {
                    return;
                }
                ref.sourceDocumentId = response.document!.id!;

                $grpc.documents.documents.addDocumentReference({
                    reference: ref,
                });
            });
        }

        if (canDo.value.relations) {
            const relationsToRemove: number[] = [];
            currentRelations.value.forEach((rel) => {
                if (!relationManagerData.value.has(rel.id!)) relationsToRemove.push(rel.id!);
            });
            relationsToRemove.forEach((id) => {
                $grpc.documents.documents.removeDocumentRelation({ id });
            });
            relationManagerData.value.forEach((rel) => {
                if (currentRelations.value.find((r) => r.id === rel.id!)) {
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
        documentStore.clear();

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

useYArrayFiltered<File>(
    ydoc.getArray('files'),
    toRef(state, 'files'),
    {
        omit: ['createdAt', 'meta'],
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

                <UButtonGroup class="inline-flex">
                    <UButton
                        type="submit"
                        trailing-icon="i-mdi-content-save"
                        :disabled="!canEdit || !canSubmit"
                        :loading="!canSubmit"
                    >
                        <span class="hidden truncate sm:block">
                            {{ $t('common.save') }}
                        </span>
                    </UButton>

                    <UButton
                        v-if="document?.document?.draft"
                        type="submit"
                        color="info"
                        trailing-icon="i-mdi-publish"
                        :disabled="!canEdit || !canSubmit"
                        :loading="!canSubmit"
                        @click.prevent="
                            modal.open(ConfirmModal, {
                                title: $t('common.publish_confirm.title', { type: $t('common.document', 1) }),
                                description: $t('common.publish_confirm.description'),
                                color: 'info',
                                iconClass: 'text-info-500 dark:text-info-400',
                                icon: 'i-mdi-publish',
                                confirm: () => {
                                    state.draft = !state.draft;
                                    formRef?.submit();
                                },
                            })
                        "
                    >
                        <span class="hidden truncate sm:block">
                            {{ $t('common.publish') }}
                        </span>
                    </UButton>
                </UButtonGroup>
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
                                        :disabled="!canEdit || !canDo.edit"
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
                                                            categories.unshift(emptyCategory);
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
                                            :disabled="!canEdit || !canDo.edit"
                                        />
                                    </UFormGroup>

                                    <UFormGroup class="flex-initial" name="closed" :label="`${$t('common.close', 2)}?`">
                                        <UToggle v-model="state.closed" :disabled="!canDo.edit" />
                                    </UFormGroup>
                                </div>
                            </div>
                        </template>
                    </UDashboardToolbar>

                    <DocumentRelationManager
                        v-model="relationManagerData"
                        :open="openRelationManager"
                        :document-id="documentId"
                        @close="openRelationManager = false"
                    />
                    <DocumentReferenceManager
                        v-model="referenceManagerData"
                        :open="openReferenceManager"
                        :document-id="documentId"
                        @close="openReferenceManager = false"
                    />

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
                                :disabled="!canEdit || !canDo.edit"
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
                                :disabled="!canEdit || !canDo.edit"
                                icon="i-mdi-account-multiple"
                                @click="openRelationManager = true"
                            >
                                {{ $t('common.citizen', 1) }} {{ $t('common.relation', 2) }}
                            </UButton>
                            <UButton
                                v-if="canDo.references"
                                class="flex-1"
                                block
                                :disabled="!canEdit || !canDo.edit"
                                icon="i-mdi-file-document"
                                @click="openReferenceManager = true"
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
                            :disabled="!canEdit || !canDo.access"
                        />
                    </div>
                </template>
            </UTabs>
        </UDashboardPanelContent>
    </UForm>
</template>
