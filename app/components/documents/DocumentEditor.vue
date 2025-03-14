<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import DocumentReferenceManager from '~/components/documents/DocumentReferenceManager.vue';
import DocumentRelationManager from '~/components/documents/DocumentRelationManager.vue';
import { checkDocAccess, logger } from '~/components/documents/helpers';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { availableIcons, fallbackIcon } from '~/components/partials/icons';
import { getDocument, getUser, useClipboardStore } from '~/stores/clipboard';
import { useCompletorStore } from '~/stores/completor';
import { useDocumentEditorStore } from '~/stores/documenteditor';
import { useNotificatorStore } from '~/stores/notificator';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { DocumentJobAccess, DocumentUserAccess } from '~~/gen/ts/resources/documents/access';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { DocumentReference, DocumentRelation } from '~~/gen/ts/resources/documents/documents';
import { DocReference, DocRelation } from '~~/gen/ts/resources/documents/documents';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { CreateDocumentRequest, UpdateDocumentRequest } from '~~/gen/ts/services/docstore/docstore';

const props = defineProps<{
    documentId?: number;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { can, activeChar } = useAuth();

const { game } = useAppConfig();

const clipboardStore = useClipboardStore();

const completorStore = useCompletorStore();

const documentStore = useDocumentEditorStore();

const notifications = useNotificatorStore();

const { maxAccessEntries } = useAppConfig();

const route = useRoute();

const canEdit = ref(false);

const schema = z.object({
    title: z.string().min(3).max(255),
    state: z.union([z.string().length(0), z.string().min(3).max(32)]),
    content: z.string().min(3).max(1750000),
    public: z.boolean(),
    closed: z.boolean(),
    category: z.custom<Category>().optional(),
    access: z.object({
        jobs: z.custom<DocumentJobAccess>().array().max(maxAccessEntries),
        users: z.custom<DocumentUserAccess>().array().max(maxAccessEntries),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    title: '',
    state: '',
    content: '',
    public: false,
    closed: false,
    category: undefined,
    access: {
        jobs: [],
        users: [],
    },
});

const docCreator = ref<UserShort | undefined>();

const openRelationManager = ref<boolean>(false);
const relationManagerData = ref(new Map<number, DocumentRelation>());
const currentRelations = ref<Readonly<DocumentRelation>[]>([]);
watch(currentRelations, () => currentRelations.value.forEach((e) => relationManagerData.value.set(e.id!, e)));

const openReferenceManager = ref<boolean>(false);
const referenceManagerData = ref(new Map<number, DocumentReference>());
const currentReferences = ref<Readonly<DocumentReference>[]>([]);
watch(currentReferences, () => currentReferences.value.forEach((e) => referenceManagerData.value.set(e.id!, e)));

const emptyCategory: Category = {
    id: 0,
    name: t('common.categories', 0),
};
const templateId = ref<undefined | number>();

onMounted(async () => {
    if (route.query.templateId) {
        const data = clipboardStore.getTemplateData();
        data.activeChar = activeChar.value!;
        logger.debug('Editor - Clipboard Template Data', data);

        templateId.value = parseInt(route.query.templateId as string);

        try {
            const call = $grpc.docstore.docStore.getTemplate({
                templateId: templateId.value,
                data,
                render: true,
            });
            const { response } = await call;

            if (response.template === undefined) {
                throw new Error('failed to get template from server response');
            }

            const template = response.template;
            state.title = template.contentTitle;
            state.state = template.state;
            state.content = template.content;
            state.category = template.category ?? emptyCategory;
            if (template?.contentAccess) {
                state.access = template.contentAccess;
            }
            if (activeChar.value !== null) {
                docCreator.value = activeChar.value;
            }
        } catch (e) {
            handleGRPCError(e as RpcError);
            logger.error('Editor - Template Error', e);

            await navigateTo({ name: 'documents' });

            return;
        }
    } else if (props.documentId) {
        try {
            const req = { documentId: props.documentId };
            const call = $grpc.docstore.docStore.getDocument(req);
            const { response } = await call;

            const document = response.document;
            docCreator.value = document?.creator;
            if (document) {
                state.title = document.title;
                state.state = document.state;
                state.content = document.content?.rawContent ?? '';
                state.category = document.category ?? emptyCategory;
                state.closed = document.closed;
                state.public = document.public;
                state.access = response.access!;

                const refs = await $grpc.docstore.docStore.getDocumentReferences(req);
                currentReferences.value = refs.response.references;
                const rels = await $grpc.docstore.docStore.getDocumentRelations(req);
                currentRelations.value = rels.response.relations;
            }
        } catch (e) {
            handleGRPCError(e as RpcError);

            await navigateTo({ name: 'documents' });

            return;
        }
    } else {
        state.title = documentStore.$state.title;
        state.state = documentStore.$state.state;
        state.content = documentStore.$state.content;
        state.category = documentStore.$state.category;
        state.closed = documentStore.$state.closed;

        state.access.jobs.push({
            id: 0,
            targetId: props.documentId ?? 0,
            job: activeChar.value!.job,
            minimumGrade: game.startJobGrade,
            access: AccessLevel.EDIT,
        });
    }

    clipboardStore.activeStack.documents.forEach((doc, idx) => {
        referenceManagerData.value.set(idx, {
            id: idx,
            sourceDocumentId: props.documentId ?? 0,
            targetDocumentId: doc.id!,
            targetDocument: getDocument(doc),
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
    let prom: Promise<void>;
    if (props.documentId === undefined) {
        prom = createDocument(event.data);
    } else {
        prom = updateDocument(props.documentId, event.data);
    }

    await prom.finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
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

async function createDocument(values: Schema): Promise<void> {
    // Prepare request
    const req: CreateDocumentRequest = {
        title: values.title,
        content: {
            rawContent: values.content,
        },
        contentType: ContentType.HTML,
        closed: values.closed,
        state: values.state,
        public: values.public,
        templateId: templateId.value,
        categoryId: values.category?.id !== 0 ? values.category?.id : undefined,
        access: values.access,
    };

    // Try to submit to server
    try {
        const call = $grpc.docstore.docStore.createDocument(req);
        const { response } = await call;

        const promises: Promise<unknown>[] = [];
        if (canDo.value.references) {
            referenceManagerData.value.forEach((ref) => {
                ref.sourceDocumentId = response.documentId;

                const prom = $grpc.docstore.docStore.addDocumentReference({
                    reference: ref,
                });
                promises.push(prom.response);
            });
        }

        if (canDo.value.relations) {
            relationManagerData.value.forEach((rel) => {
                rel.documentId = response.documentId;

                const prom = $grpc.docstore.docStore.addDocumentRelation({
                    relation: rel,
                });
                promises.push(prom.response);
            });
        }
        await Promise.all(promises);

        notifications.add({
            title: { key: 'notifications.document_created.title', parameters: {} },
            description: { key: 'notifications.document_created.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
        clipboardStore.clear();
        documentStore.clear();

        await navigateTo({
            name: 'documents-id',
            params: { id: response.documentId },
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function updateDocument(id: number, values: Schema): Promise<void> {
    const req: UpdateDocumentRequest = {
        documentId: id,
        title: values.title,
        content: {
            rawContent: values.content,
        },
        contentType: ContentType.HTML,
        closed: values.closed,
        state: values.state,
        public: values.public,
        categoryId: values.category?.id !== 0 ? values.category?.id : undefined,
        access: values.access,
    };

    try {
        const call = $grpc.docstore.docStore.updateDocument(req);
        const { response } = await call;

        if (canDo.value.references) {
            const referencesToRemove: number[] = [];
            currentReferences.value.forEach((ref) => {
                if (!referenceManagerData.value.has(ref.id!)) {
                    referencesToRemove.push(ref.id!);
                }
            });
            referencesToRemove.forEach((id) => {
                $grpc.docstore.docStore.removeDocumentReference({
                    id: id,
                });
            });
            referenceManagerData.value.forEach((ref) => {
                if (currentReferences.value.find((r) => r.id === ref.id!)) {
                    return;
                }
                ref.sourceDocumentId = response.documentId;

                $grpc.docstore.docStore.addDocumentReference({
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
                $grpc.docstore.docStore.removeDocumentRelation({ id });
            });
            relationManagerData.value.forEach((rel) => {
                if (currentRelations.value.find((r) => r.id === rel.id!)) {
                    return;
                }
                rel.documentId = response.documentId;

                $grpc.docstore.docStore.addDocumentRelation({
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
            params: { id: response.documentId },
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const items = [
    {
        slot: 'edit',
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
            : checkDocAccess(state.access, docCreator.value, AccessLevel.EDIT, 'DocStoreService.UpdateDocument'),
    access:
        props.documentId === undefined
            ? true
            : checkDocAccess(state.access, docCreator.value, AccessLevel.ACCESS, 'DocStoreService.UpdateDocument'),
    references: can('DocStoreService.AddDocumentReference').value,
    relations: can('DocStoreService.AddDocumentRelation').value,
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
</script>

<template>
    <UForm
        :schema="schema"
        :state="state"
        class="min-h-dscreen flex w-full max-w-full flex-1 flex-col overflow-y-auto"
        @submit="onSubmitThrottle"
    >
        <UDashboardNavbar :title="documentId ? $t('pages.documents.edit.title') : $t('pages.documents.create.title')">
            <template #right>
                <PartialsBackButton
                    :fallback-to="documentId ? { name: 'documents-id', params: { id: documentId } } : `/documents`"
                />

                <UButtonGroup class="inline-flex">
                    <UButton
                        type="submit"
                        trailing-icon="i-mdi-content-save"
                        :disabled="!canEdit || !canSubmit"
                        :loading="!canSubmit"
                    >
                        <span class="hidden truncate sm:block">
                            <template v-if="!documentId">
                                {{ $t('common.create') }}
                            </template>
                            <template v-else>
                                {{ $t('common.save') }}
                            </template>
                        </span>
                    </UButton>
                </UButtonGroup>
            </template>
        </UDashboardNavbar>

        <UDashboardPanelContent class="p-0">
            <UTabs
                v-model="selectedTab"
                :items="items"
                class="flex flex-1 flex-col"
                :ui="{
                    wrapper: 'space-y-0 overflow-y-hidden',
                    container: 'flex flex-1 flex-col overflow-y-hidden',
                    base: 'flex flex-1 flex-col overflow-y-hidden',
                    list: { rounded: '' },
                }"
            >
                <template #edit>
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
                                    <UFormGroup name="category" :label="$t('common.category', 1)" class="flex-1">
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

                                    <UFormGroup name="state" :label="$t('common.state')" class="flex-1">
                                        <UInput
                                            v-model="state.state"
                                            type="text"
                                            :placeholder="`${$t('common.document', 1)} ${$t('common.state')}`"
                                            :disabled="!canEdit || !canDo.edit"
                                        />
                                    </UFormGroup>

                                    <UFormGroup name="closed" :label="`${$t('common.close', 2)}?`" class="flex-1">
                                        <ClientOnly>
                                            <USelectMenu
                                                v-model="state.closed"
                                                :disabled="!canEdit || !canDo.edit"
                                                :options="[
                                                    { label: $t('common.open', 2), closed: false },
                                                    { label: $t('common.close', 2), closed: true },
                                                ]"
                                                value-attribute="closed"
                                                :searchable-placeholder="$t('common.search_field')"
                                            >
                                                <template #option-empty="{ query: search }">
                                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                                </template>

                                                <template #empty>
                                                    {{ $t('common.not_found', [$t('common.close', 1)]) }}
                                                </template>
                                            </USelectMenu>
                                        </ClientOnly>
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
                        name="content"
                        class="flex flex-1 overflow-y-hidden"
                        :ui="{ container: 'flex flex-1 mt-0 overflow-y-hidden', label: { wrapper: 'hidden' } }"
                        label="&nbsp;"
                    >
                        <ClientOnly>
                            <TiptapEditor
                                v-model="state.content"
                                :disabled="!canEdit || !canDo.edit"
                                class="mx-auto w-full max-w-screen-xl flex-1 overflow-y-hidden"
                                rounded="rounded-none"
                            >
                                <template #footer>
                                    <div v-if="saving" class="place-self-start">
                                        <UIcon name="i-mdi-content-save" class="h-4 w-4 animate-spin" />
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
                            :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.docstore.AccessLevel')"
                            :disabled="!canEdit || !canDo.access"
                        />
                    </div>
                </template>
            </UTabs>
        </UDashboardPanelContent>
    </UForm>
</template>
