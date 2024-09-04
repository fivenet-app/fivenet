<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import DocumentAccessEntry from '~/components/documents/DocumentAccessEntry.vue';
import DocumentReferenceManager from '~/components/documents/DocumentReferenceManager.vue';
import DocumentRelationManager from '~/components/documents/DocumentRelationManager.vue';
import { checkDocAccess, logger } from '~/components/documents/helpers';
import DocEditor from '~/components/partials/DocEditor.vue';
import { useAuthStore } from '~/store/auth';
import { getDocument, getUser, useClipboardStore } from '~/store/clipboard';
import { useCompletorStore } from '~/store/completor';
import { useDocumentEditorStore } from '~/store/documenteditor';
import { useNotificatorStore } from '~/store/notificator';
import { AccessLevel, DocumentAccess } from '~~/gen/ts/resources/documents/access';
import { Category } from '~~/gen/ts/resources/documents/category';
import {
    DocContentType,
    DocReference,
    DocRelation,
    DocumentReference,
    DocumentRelation,
} from '~~/gen/ts/resources/documents/documents';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { CreateDocumentRequest, UpdateDocumentRequest } from '~~/gen/ts/services/docstore/docstore';

const props = defineProps<{
    documentId?: string;
}>();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const clipboardStore = useClipboardStore();

const completorStore = useCompletorStore();

const documentStore = useDocumentEditorStore();

const notifications = useNotificatorStore();

const { t } = useI18n();

const route = useRoute();

const maxAccessEntries = 10;

const canEdit = ref(false);

const schema = z.object({
    title: z.string().min(3).max(255),
    state: z.union([z.string().length(0), z.string().min(3).max(32)]),
    content: z.string().min(20).max(1750000),
    public: z.boolean(),
    closed: z.boolean(),
    category: z.custom<Category>().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    title: '',
    state: '',
    content: '',
    public: false,
    closed: false,
    category: undefined,
});

const access = ref(
    new Map<
        string,
        {
            id: string;
            type: number;
            values: {
                job?: string;
                userId?: number;
                accessRole?: AccessLevel;
                minimumGrade?: number;
            };
            required?: boolean;
        }
    >(),
);
const docAccess = ref<DocumentAccess>();
const docCreator = ref<UserShort | undefined>();

const openRelationManager = ref<boolean>(false);
const relationManagerData = ref(new Map<string, DocumentRelation>());
const currentRelations = ref<Readonly<DocumentRelation>[]>([]);
watch(currentRelations, () => currentRelations.value.forEach((e) => relationManagerData.value.set(e.id!, e)));

const openReferenceManager = ref<boolean>(false);
const referenceManagerData = ref(new Map<string, DocumentReference>());
const currentReferences = ref<Readonly<DocumentReference>[]>([]);
watch(currentReferences, () => currentReferences.value.forEach((e) => referenceManagerData.value.set(e.id!, e)));

const templateId = ref<undefined | string>();

onMounted(async () => {
    if (route.query.templateId) {
        const data = clipboardStore.getTemplateData();
        data.activeChar = activeChar.value!;
        logger.debug('Editor - Clipboard Template Data', data);

        templateId.value = route.query.templateId as string;

        try {
            const call = getGRPCDocStoreClient().getTemplate({
                templateId: templateId.value as string,
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
            state.category = template.category;

            if (template?.contentAccess) {
                if (authStore.activeChar !== null) {
                    docCreator.value = authStore.activeChar;
                }
                const docAccess = template.contentAccess!;
                let accessId = 0;
                docAccess.users.forEach((user) => {
                    const id = accessId.toString();
                    access.value.set(id, {
                        id,
                        type: 0,
                        values: { userId: user.userId, accessRole: user.access },
                        required: user.required,
                    });
                    accessId++;
                });

                docAccess.jobs.forEach((job) => {
                    const id = accessId.toString();
                    access.value.set(id, {
                        id,
                        type: 1,
                        values: {
                            job: job.job,
                            accessRole: job.access,
                            minimumGrade: job.minimumGrade,
                        },
                        required: job.required,
                    });
                    accessId++;
                });
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
            const call = getGRPCDocStoreClient().getDocument(req);
            const { response } = await call;

            const document = response.document;
            docAccess.value = response.access;
            docCreator.value = document?.creator;
            if (document) {
                state.title = document.title;
                state.state = document.state;
                state.content = document.content;
                state.category = document.category;
                state.closed = document.closed;
                state.public = document.public;

                const refs = await getGRPCDocStoreClient().getDocumentReferences(req);
                currentReferences.value = refs.response.references;
                const rels = await getGRPCDocStoreClient().getDocumentRelations(req);
                currentRelations.value = rels.response.relations;
            }

            if (response.access) {
                let accessId = 0;

                response.access.users.forEach((user) => {
                    const id = accessId.toString();
                    access.value.set(id, {
                        id,
                        type: 0,
                        values: { userId: user.userId, accessRole: user.access },
                    });
                    accessId++;
                });

                response.access.jobs.forEach((job) => {
                    const id = accessId.toString();
                    access.value.set(id, {
                        id,
                        type: 1,
                        values: {
                            job: job.job,
                            accessRole: job.access,
                            minimumGrade: job.minimumGrade,
                        },
                    });
                    accessId++;
                });
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

        const accessId = 0;
        access.value.set(accessId.toString(), {
            id: accessId.toString(),
            type: 1,
            values: {
                job: activeChar.value?.job,
                minimumGrade: 1,
                accessRole: AccessLevel.EDIT,
            },
        });
    }

    clipboardStore.activeStack.documents.forEach((doc, i) => {
        const id = i.toString();
        referenceManagerData.value.set(id, {
            id,
            sourceDocumentId: props.documentId ?? '0',
            targetDocumentId: doc.id!,
            targetDocument: getDocument(doc),
            creatorId: activeChar.value!.userId,
            creator: activeChar.value!,
            reference: DocReference.SOLVES,
        });
    });
    clipboardStore.activeStack.users.forEach((user, i) => {
        const id = i.toString();
        relationManagerData.value.set(id, {
            id,
            documentId: props.documentId ?? '0',
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

const accessTypes = [
    { id: 0, name: t('common.citizen', 2) },
    { id: 1, name: t('common.job', 2) },
];

function addAccessEntry(): void {
    if (access.value.size > maxAccessEntries - 1) {
        notifications.add({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            description: {
                key: 'notifications.max_access_entry.content',
                parameters: { max: maxAccessEntries.toString() },
            },
            type: NotificationType.ERROR,
        });
        return;
    }

    const id = access.value.size > 0 ? parseInt([...access.value.keys()]?.pop() ?? '1', 10) + 1 : 0;
    access.value.set(id.toString(), {
        id: id.toString(),
        type: 1,
        values: {},
    });
}

function removeAccessEntry(event: { id: string }): void {
    access.value.delete(event.id);
}

function updateAccessEntryType(event: { id: string; type: number }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.type = event.type;
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryName(event: { id: string; job?: Job; char?: UserShort }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    if (event.job) {
        accessEntry.values.job = event.job.name;
        accessEntry.values.userId = undefined;
    } else if (event.char) {
        accessEntry.values.job = undefined;
        accessEntry.values.userId = event.char.userId;
    }

    access.value.set(event.id, accessEntry);
}

function updateAccessEntryRank(event: { id: string; rank: JobGrade }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.minimumGrade = event.rank.grade;
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryAccess(event: { id: string; access: AccessLevel }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.accessRole = event.access;
    access.value.set(event.id, accessEntry);
}

async function createDocument(values: Schema): Promise<void> {
    // Prepare request
    const req: CreateDocumentRequest = {
        title: values.title,
        content: values.content,
        contentType: DocContentType.HTML,
        closed: values.closed,
        state: values.state,
        public: values.public,
        templateId: templateId.value,
        categoryId: values.category?.id,
    };

    const reqAccess: DocumentAccess = {
        jobs: [],
        users: [],
    };
    access.value.forEach((entry) => {
        if (entry.values.accessRole === undefined) {
            return;
        }

        if (entry.type === 0) {
            if (!entry.values.userId) {
                return;
            }

            reqAccess.users.push({
                id: '0',
                documentId: '0',
                userId: entry.values.userId,
                access: entry.values.accessRole,
            });
        } else if (entry.type === 1) {
            if (!entry.values.job) {
                return;
            }

            reqAccess.jobs.push({
                id: '0',
                documentId: '0',
                job: entry.values.job,
                minimumGrade: entry.values.minimumGrade ? entry.values.minimumGrade : 0,
                access: entry.values.accessRole,
            });
        }
    });
    req.access = reqAccess;

    // Try to submit to server
    try {
        const call = getGRPCDocStoreClient().createDocument(req);
        const { response } = await call;

        const promises: Promise<any>[] = [];
        if (canDo.value.references) {
            referenceManagerData.value.forEach((ref) => {
                ref.sourceDocumentId = response.documentId;

                const prom = getGRPCDocStoreClient().addDocumentReference({
                    reference: ref,
                });
                promises.push(prom.response);
            });
        }

        if (canDo.value.relations) {
            relationManagerData.value.forEach((rel) => {
                rel.documentId = response.documentId;

                const prom = getGRPCDocStoreClient().addDocumentRelation({
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

async function updateDocument(id: string, values: Schema): Promise<void> {
    const req: UpdateDocumentRequest = {
        documentId: id,
        title: values.title,
        content: values.content,
        contentType: DocContentType.HTML,
        closed: values.closed,
        state: values.state,
        public: values.public,
        categoryId: values.category?.id,
    };

    const reqAccess: DocumentAccess = {
        jobs: [],
        users: [],
    };
    access.value.forEach((entry) => {
        if (entry.values.accessRole === undefined) {
            return;
        }

        if (entry.type === 0) {
            if (!entry.values.userId) {
                return;
            }

            reqAccess.users.push({
                id: '0',
                documentId: id,
                userId: entry.values.userId,
                access: entry.values.accessRole,
            });
        } else if (entry.type === 1) {
            if (!entry.values.job) {
                return;
            }

            reqAccess.jobs.push({
                id: '0',
                documentId: id,
                job: entry.values.job,
                minimumGrade: entry.values.minimumGrade ? entry.values.minimumGrade : 0,
                access: entry.values.accessRole,
            });
        }
    });
    req.access = reqAccess;

    try {
        const call = getGRPCDocStoreClient().updateDocument(req);
        const { response } = await call;

        if (canDo.value.references) {
            const referencesToRemove: string[] = [];
            currentReferences.value.forEach((ref) => {
                if (!referenceManagerData.value.has(ref.id!)) referencesToRemove.push(ref.id!);
            });
            referencesToRemove.forEach((id) => {
                getGRPCDocStoreClient().removeDocumentReference({ id });
            });
            referenceManagerData.value.forEach((ref) => {
                if (currentReferences.value.find((r) => r.id === ref.id!)) {
                    return;
                }
                ref.sourceDocumentId = response.documentId;

                getGRPCDocStoreClient().addDocumentReference({
                    reference: ref,
                });
            });
        }

        if (canDo.value.relations) {
            const relationsToRemove: string[] = [];
            currentRelations.value.forEach((rel) => {
                if (!relationManagerData.value.has(rel.id!)) relationsToRemove.push(rel.id!);
            });
            relationsToRemove.forEach((id) => {
                getGRPCDocStoreClient().removeDocumentRelation({ id });
            });
            relationManagerData.value.forEach((rel) => {
                if (currentRelations.value.find((r) => r.id === rel.id!)) {
                    return;
                }
                rel.documentId = response.documentId;

                getGRPCDocStoreClient().addDocumentRelation({
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

const categoriesLoading = ref(false);

const canDo = computed(() => ({
    edit:
        props.documentId === undefined
            ? true
            : checkDocAccess(docAccess.value, docCreator.value, AccessLevel.EDIT, 'DocStoreService.UpdateDocument'),
    access:
        props.documentId === undefined
            ? true
            : checkDocAccess(docAccess.value, docCreator.value, AccessLevel.ACCESS, 'DocStoreService.UpdateDocument'),
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

const { data: jobs } = useAsyncData('completor-jobs', () => completorStore.listJobs());
</script>

<template>
    <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
        <UDashboardNavbar :title="$t('pages.documents.edit.title')">
            <template #right>
                <UButtonGroup class="inline-flex">
                    <UButton
                        color="black"
                        icon="i-mdi-arrow-left"
                        :to="documentId ? { name: 'documents-id', params: { id: documentId } } : `/documents`"
                    >
                        {{ $t('common.back') }}
                    </UButton>

                    <UButton
                        type="submit"
                        trailing-icon="i-mdi-content-save"
                        :disabled="!canEdit || !canSubmit"
                        :loading="!canSubmit"
                    >
                        <template v-if="!documentId">
                            {{ $t('common.create') }}
                        </template>
                        <template v-else>
                            {{ $t('common.save') }}
                        </template>
                    </UButton>
                </UButtonGroup>
            </template>
        </UDashboardNavbar>

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
                            <UInputMenu
                                v-model="state.category"
                                option-attribute="name"
                                :search-attributes="['name']"
                                block
                                nullable
                                :search="
                                    async (search: string) => {
                                        try {
                                            categoriesLoading = true;
                                            const categories = await completorStore.completeDocumentCategories(search);
                                            categoriesLoading = false;
                                            return categories;
                                        } catch (e) {
                                            handleGRPCError(e as RpcError);
                                            throw e;
                                        } finally {
                                            categoriesLoading = false;
                                        }
                                    }
                                "
                                search-lazy
                                :search-placeholder="$t('common.search_field')"
                            >
                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>
                                <template #empty> {{ $t('common.not_found', [$t('common.category', 2)]) }} </template>
                            </UInputMenu>
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
                        </UFormGroup>
                    </div>
                </div>
            </template>
        </UDashboardToolbar>

        <DocumentRelationManager
            v-model="relationManagerData"
            :open="openRelationManager"
            :document="documentId"
            @close="openRelationManager = false"
        />
        <DocumentReferenceManager
            v-model="referenceManagerData"
            :open="openReferenceManager"
            :document-id="documentId"
            @close="openReferenceManager = false"
        />

        <template v-if="canDo.edit">
            <UFormGroup name="content">
                <ClientOnly>
                    <DocEditor v-model="state.content" :disabled="!canEdit || !canDo.edit" />
                </ClientOnly>
            </UFormGroup>

            <template v-if="saving">
                <div class="flex animate-pulse justify-center">
                    <UIcon name="i-mdi-content-save" class="mr-2 h-auto w-4 animate-spin" />
                    <span>{{ $t('common.save', 2) }}...</span>
                </div>
            </template>
        </template>

        <div class="mt-2 flex flex-col gap-2 px-2">
            <UButtonGroup v-if="canDo.edit" class="mt-2 inline-flex w-full">
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

            <div>
                <h2 class="text- text-gray-900 dark:text-white">
                    {{ $t('common.access') }}
                </h2>

                <DocumentAccessEntry
                    v-for="entry in access.values()"
                    :key="entry.id"
                    :init="entry"
                    :access-types="accessTypes"
                    :read-only="!canDo.access || entry.required === true"
                    :jobs="jobs"
                    @type-change="updateAccessEntryType($event)"
                    @name-change="updateAccessEntryName($event)"
                    @rank-change="updateAccessEntryRank($event)"
                    @access-change="updateAccessEntryAccess($event)"
                    @delete-request="removeAccessEntry($event)"
                />

                <UButton
                    :disabled="!canEdit || !canDo.access"
                    :ui="{ rounded: 'rounded-full' }"
                    icon="i-mdi-plus"
                    :title="$t('components.documents.document_editor.add_permission')"
                    @click="addAccessEntry()"
                />
            </div>
        </div>
    </UForm>
</template>
