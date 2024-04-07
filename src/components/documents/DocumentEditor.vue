<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { type TranslateItem } from '~/composables/i18n';
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
import { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { CreateDocumentRequest, UpdateDocumentRequest } from '~~/gen/ts/services/docstore/docstore';
import DocumentAccessEntry from '~/components/documents/DocumentAccessEntry.vue';
import DocumentReferenceManager from '~/components/documents/DocumentReferenceManager.vue';
import DocumentRelationManager from '~/components/documents/DocumentRelationManager.vue';
import { checkDocAccess } from '~/components/documents/helpers';
import DocEditor from '~/components/partials/DocEditor.vue';

const props = defineProps<{
    documentId?: string;
}>();

const { $grpc } = useNuxtApp();

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

interface FormData {
    title: string;
    state: string;
    content: string;
    public: boolean;
}

const openclose = [
    { id: 0, label: t('common.open', 2), closed: false },
    { id: 1, label: t('common.close', 2), closed: true },
];

const doc = ref<{
    closed: { id: number; label: string; closed: boolean };
}>({
    closed: openclose[0],
});
const access = ref<
    Map<
        string,
        {
            id: string;
            type: number;
            values: {
                job?: string;
                char?: number;
                accessRole?: AccessLevel;
                minimumGrade?: number;
            };
            required?: boolean;
        }
    >
>(new Map());
const docAccess = ref<DocumentAccess>();
const docCreator = ref<UserShort | undefined>();

const openRelationManager = ref<boolean>(false);
const relationManagerData = ref<Map<string, DocumentRelation>>(new Map());
const currentRelations = ref<Readonly<DocumentRelation>[]>([]);
watch(currentRelations, () => currentRelations.value.forEach((e) => relationManagerData.value.set(e.id!, e)));

const openReferenceManager = ref<boolean>(false);
const referenceManagerData = ref<Map<string, DocumentReference>>(new Map());
const currentReferences = ref<Readonly<DocumentReference>[]>([]);
watch(currentReferences, () => currentReferences.value.forEach((e) => referenceManagerData.value.set(e.id!, e)));

const selectedCategory = ref<Category | undefined>(undefined);

const templateId = ref<undefined | string>();

onMounted(async () => {
    if (route.query.templateId) {
        const data = clipboardStore.getTemplateData();
        data.activeChar = activeChar.value!;
        console.debug('Documents: Editor - Clipboard Template Data', data);

        templateId.value = route.query.templateId as string;

        try {
            const call = $grpc.getDocStoreClient().getTemplate({
                templateId: templateId.value as string,
                data,
                render: true,
            });
            const { response } = await call;

            if (response.template === undefined) {
                throw new Error('failed to get template from server response');
            }

            const template = response.template;
            setFieldValue('title', template.contentTitle!);
            setFieldValue('state', template.state!);
            setFieldValue('content', template.content.replace(/\s+/g, ' '));
            selectedCategory.value = template?.category;

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
                        values: { char: user.userId, accessRole: user.access },
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
            $grpc.handleError(e as RpcError);
            console.error('Documents: Editor - Template Error', e);

            await navigateTo({ name: 'documents' });

            return;
        }
    } else if (props.documentId) {
        try {
            const req = { documentId: props.documentId };
            const call = $grpc.getDocStoreClient().getDocument(req);
            const { response } = await call;
            const document = response.document;
            docAccess.value = response.access;
            docCreator.value = document?.creator;

            if (document) {
                setFieldValue('title', document.title);
                setFieldValue('state', document.state);
                setFieldValue('content', document.content);
                doc.value.closed = openclose.find((e) => e.closed === document.closed) as {
                    id: number;
                    label: string;
                    closed: boolean;
                };
                selectedCategory.value = document.category;
                setFieldValue('public', document.public);

                const refs = await $grpc.getDocStoreClient().getDocumentReferences(req);
                currentReferences.value = refs.response.references;
                const rels = await $grpc.getDocStoreClient().getDocumentRelations(req);
                currentRelations.value = rels.response.relations;
            }

            if (response.access) {
                let accessId = 0;

                response.access.users.forEach((user) => {
                    const id = accessId.toString();
                    access.value.set(id, {
                        id,
                        type: 0,
                        values: { char: user.userId, accessRole: user.access },
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
            $grpc.handleError(e as RpcError);

            await navigateTo({ name: 'documents' });

            return;
        }
    } else {
        if (documentStore.$state) {
            setFieldValue('title', documentStore.$state.title);
            setFieldValue('state', documentStore.$state.state);
            setFieldValue('content', documentStore.$state.content);
            if (documentStore.$state.closed) {
                doc.value.closed = documentStore.$state.closed;
            }
        }

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

async function saveToStore(values: FormData): Promise<void> {
    if (saving.value) {
        return;
    }
    saving.value = true;

    documentStore.save({
        title: values.title,
        content: values.content,
        state: values.state,
        closed: doc.value.closed,
        category: selectedCategory.value,
    });

    useTimeoutFn(() => {
        saving.value = false;
    }, 1250);
}

const changed = ref(false);

defineRule('required', required);
defineRule('max', max);
defineRule('min', min);

const { handleSubmit, values, setFieldValue, meta } = useForm<FormData>({
    validationSchema: {
        title: { required: true, min: 3, max: 255 },
        state: { required: false, min: 2, max: 32 },
    },
    initialValues: {
        public: false,
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(async (values): Promise<void> => {
    let prom: Promise<void>;
    if (props.documentId === undefined) {
        prom = createDocument(values, doc.value.closed.closed);
    } else {
        prom = updateDocument(props.documentId, values, doc.value.closed.closed);
    }

    await prom.finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
});
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

watchOnce(meta, () => (changed.value = true));
watchDebounced(
    doc.value,
    async () => {
        if (changed.value) {
            saveToStore(values);
        }
    },
    {
        debounce: 1000,
        maxWait: 2500,
    },
);
watchDebounced(
    meta,
    async () => {
        if (changed.value) {
            saveToStore(values);
        }
    },
    {
        debounce: 1000,
        maxWait: 3500,
    },
);

const accessTypes = [
    { id: 0, name: t('common.citizen', 2) },
    { id: 1, name: t('common.job', 2) },
];

function addDocumentAccessEntry(): void {
    if (access.value.size > maxAccessEntries - 1) {
        notifications.add({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            description: {
                key: 'notifications.max_access_entry.content',
                parameters: { max: maxAccessEntries.toString() },
            } as TranslateItem,
            type: 'error',
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

function removeDocumentAccessEntry(event: { id: string }): void {
    access.value.delete(event.id);
}

function updateDocumentAccessEntryType(event: { id: string; type: number }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.type = event.type;
    access.value.set(event.id, accessEntry);
}

function updateDocumentAccessEntryName(event: { id: string; job?: Job; char?: UserShort }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    if (event.job) {
        accessEntry.values.job = event.job.name;
        accessEntry.values.char = undefined;
    } else if (event.char) {
        accessEntry.values.job = undefined;
        accessEntry.values.char = event.char.userId;
    }

    access.value.set(event.id, accessEntry);
}

function updateDocumentAccessEntryRank(event: { id: string; rank: JobGrade }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.minimumGrade = event.rank.grade;
    access.value.set(event.id, accessEntry);
}

function updateDocumentAccessEntryAccess(event: { id: string; access: AccessLevel }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.accessRole = event.access;
    access.value.set(event.id, accessEntry);
}

async function createDocument(values: FormData, closed: boolean): Promise<void> {
    // Prepare request
    const req: CreateDocumentRequest = {
        title: values.title,
        content: values.content,
        contentType: DocContentType.HTML,
        closed,
        state: values.state,
        public: values.public,
        templateId: templateId.value,
    };
    if (selectedCategory.value !== undefined) {
        req.categoryId = selectedCategory.value.id;
    }

    const reqAccess: DocumentAccess = {
        jobs: [],
        users: [],
    };
    access.value.forEach((entry) => {
        if (entry.values.accessRole === undefined) {
            return;
        }

        if (entry.type === 0) {
            if (!entry.values.char) {
                return;
            }

            reqAccess.users.push({
                id: '0',
                documentId: '0',
                userId: entry.values.char,
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
        const call = $grpc.getDocStoreClient().createDocument(req);
        const { response } = await call;

        const promises: Promise<any>[] = [];
        if (canDo.value.references) {
            referenceManagerData.value.forEach((ref) => {
                ref.sourceDocumentId = response.documentId;

                const prom = $grpc.getDocStoreClient().addDocumentReference({
                    reference: ref,
                });
                promises.push(prom.response);
            });
        }

        if (canDo.value.relations) {
            relationManagerData.value.forEach((rel) => {
                rel.documentId = response.documentId;

                const prom = $grpc.getDocStoreClient().addDocumentRelation({
                    relation: rel,
                });
                promises.push(prom.response);
            });
        }
        await Promise.all(promises);

        notifications.add({
            title: { key: 'notifications.document_created.title', parameters: {} },
            description: { key: 'notifications.document_created.content', parameters: {} },
            type: 'success',
        });
        clipboardStore.clear();
        documentStore.clear();

        await navigateTo({
            name: 'documents-id',
            params: { id: response.documentId },
        });
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function updateDocument(id: string, values: FormData, closed: boolean): Promise<void> {
    const req: UpdateDocumentRequest = {
        documentId: id,
        title: values.title,
        content: values.content,
        contentType: DocContentType.HTML,
        closed,
        state: values.state,
        public: values.public,
    };
    if (selectedCategory.value !== undefined) {
        req.categoryId = selectedCategory.value.id;
    }

    const reqAccess: DocumentAccess = {
        jobs: [],
        users: [],
    };
    access.value.forEach((entry) => {
        if (entry.values.accessRole === undefined) {
            return;
        }

        if (entry.type === 0) {
            if (!entry.values.char) {
                return;
            }

            reqAccess.users.push({
                id: '0',
                documentId: id,
                userId: entry.values.char,
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
        const call = $grpc.getDocStoreClient().updateDocument(req);
        const { response } = await call;

        if (canDo.value.references) {
            const referencesToRemove: string[] = [];
            currentReferences.value.forEach((ref) => {
                if (!referenceManagerData.value.has(ref.id!)) referencesToRemove.push(ref.id!);
            });
            referencesToRemove.forEach((id) => {
                $grpc.getDocStoreClient().removeDocumentReference({ id });
            });
            referenceManagerData.value.forEach((ref) => {
                if (currentReferences.value.find((r) => r.id === ref.id!)) {
                    return;
                }
                ref.sourceDocumentId = response.documentId;

                $grpc.getDocStoreClient().addDocumentReference({
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
                $grpc.getDocStoreClient().removeDocumentRelation({ id });
            });
            relationManagerData.value.forEach((rel) => {
                if (currentRelations.value.find((r) => r.id === rel.id!)) {
                    return;
                }
                rel.documentId = response.documentId;

                $grpc.getDocStoreClient().addDocumentRelation({
                    relation: rel,
                });
            });
        }

        notifications.add({
            title: { key: 'notifications.document_updated.title', parameters: {} },
            description: { key: 'notifications.document_updated.content', parameters: {} },
            type: 'success',
        });
        clipboardStore.clear();
        documentStore.clear();

        await navigateTo({
            name: 'documents-id',
            params: { id: response.documentId },
        });
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const { data: jobs } = useAsyncData('completor-jobs', () => completorStore.listJobs());

const canDo = computed(() => ({
    edit:
        props.documentId === undefined
            ? true
            : checkDocAccess(docAccess.value, docCreator.value, AccessLevel.EDIT, 'DocStoreService.UpdateDocument'),
    access:
        props.documentId === undefined
            ? true
            : checkDocAccess(docAccess.value, docCreator.value, AccessLevel.ACCESS, 'DocStoreService.UpdateDocument'),
    references: can('DocStoreService.AddDocumentReference'),
    relations: can('DocStoreService.AddDocumentRelation'),
}));

console.info(
    'Documents Editor - Can Do: Edit',
    canDo.value.edit,
    'Access',
    canDo.value.access,
    'References',
    canDo.value.references,
    'Relations',
    canDo.value.relations,
);

const router = useRouter();
</script>

<template>
    <div>
        <UForm class="p-2" :state="{}">
            <UDashboardToolbar>
                <template #left>
                    <UButton @click="router.back()">
                        {{ $t('common.back') }}
                    </UButton>
                </template>
                <template #right>
                    <UButton
                        block
                        :disabled="!meta.valid || !canEdit || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle"
                    >
                        <template v-if="!documentId">
                            {{ $t('common.create') }}
                        </template>
                        <template v-else>
                            {{ $t('common.save') }}
                        </template>
                    </UButton>
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

            <div class="flex flex-col gap-2">
                <UFormGroup name="title" :label="$t('common.title')">
                    <VeeField
                        name="title"
                        :placeholder="$t('common.title')"
                        :label="$t('common.title')"
                        :disabled="!canEdit || !canDo.edit"
                    >
                        <UInput
                            type="text"
                            size="xl"
                            :placeholder="$t('common.title')"
                            :disabled="!canEdit || !canDo.edit"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </VeeField>
                    <VeeErrorMessage name="title" as="p" class="mt-2 text-sm text-error-400" />
                </UFormGroup>

                <div class="flex flex-row gap-2">
                    <UFormGroup class="flex-1" :label="$t('common.category', 1)">
                        <UInputMenu
                            v-model="selectedCategory"
                            option-attribute="name"
                            :search-attributes="['name']"
                            block
                            nullable
                            :search="completorStore.completeDocumentCategories"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        >
                            <template #option-empty="{ query: search }">
                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                            </template>
                            <template #empty> {{ $t('common.not_found', [$t('common.category', 2)]) }} </template>
                        </UInputMenu>
                    </UFormGroup>

                    <UFormGroup name="state" :label="$t('common.state')" class="flex-1">
                        <VeeField
                            name="state"
                            type="text"
                            :placeholder="`${$t('common.document', 1)} ${$t('common.state')}`"
                            :label="`${$t('common.document', 1)} ${$t('common.state')}`"
                            :disabled="!canEdit || !canDo.edit"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="state" as="p" class="mt-2 text-sm text-error-400" />
                    </UFormGroup>

                    <UFormGroup name="closed" :label="`${$t('common.close', 2)}?`" class="flex-1">
                        <USelectMenu
                            v-model="doc.closed"
                            :options="openclose"
                            :placeholder="doc.closed ? doc.closed.label : $t('common.na')"
                            :disabled="!canEdit || !canDo.edit"
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

            <div v-if="canDo.edit">
                <VeeField
                    v-slot="{ field }"
                    name="content"
                    :placeholder="$t('common.document', 1)"
                    :label="$t('common.document', 1)"
                    :disabled="!canEdit || !canDo.edit"
                >
                    <DocEditor v-bind="field" :model-value="field.value ?? ''" />
                </VeeField>
                <VeeErrorMessage name="content" as="p" class="mt-2 text-sm text-error-400" />

                <template v-if="saving">
                    <div class="flex animate-pulse justify-center">
                        <UIcon name="i-mdi-content-save" class="mr-2 h-auto w-4 animate-spin" />
                        <span>{{ $t('common.save', 2) }}...</span>
                    </div>
                </template>
            </div>

            <UButtonGroup v-if="canDo.edit" class="my-2 inline-flex w-full">
                <UButton
                    v-if="canDo.relations"
                    class="flex-1"
                    :disabled="!canEdit || !canDo.edit"
                    icon="i-mdi-account-multiple"
                    @click="openRelationManager = true"
                >
                    {{ $t('common.citizen', 1) }} {{ $t('common.relation', 2) }}
                </UButton>
                <UButton
                    v-if="canDo.references"
                    class="flex-1"
                    :disabled="!canEdit || !canDo.edit"
                    icon="i-mdi-file-document"
                    @click="openReferenceManager = true"
                >
                    {{ $t('common.document', 1) }} {{ $t('common.reference', 2) }}
                </UButton>
            </UButtonGroup>

            <div class="my-2">
                <h2>
                    {{ $t('common.access') }}
                </h2>
                <DocumentAccessEntry
                    v-for="entry in access.values()"
                    :key="entry.id"
                    :init="entry"
                    :access-types="accessTypes"
                    :read-only="!canDo.access || entry.required === true"
                    :jobs="jobs"
                    @type-change="updateDocumentAccessEntryType($event)"
                    @name-change="updateDocumentAccessEntryName($event)"
                    @rank-change="updateDocumentAccessEntryRank($event)"
                    @access-change="updateDocumentAccessEntryAccess($event)"
                    @delete-request="removeDocumentAccessEntry($event)"
                />
                <UButton
                    :disabled="!canEdit || !canDo.access"
                    :ui="{ rounded: 'rounded-full' }"
                    icon="i-mdi-plus"
                    :title="$t('components.documents.document_editor.add_permission')"
                    @click="addDocumentAccessEntry()"
                />
            </div>
        </UForm>
    </div>
</template>
