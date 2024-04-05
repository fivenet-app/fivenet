<script lang="ts" setup>
import { max, min, numeric, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import DocumentAccessEntry from '~/components/documents/DocumentAccessEntry.vue';
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { AccessLevel, DocumentAccess } from '~~/gen/ts/resources/documents/access';
import { Category } from '~~/gen/ts/resources/documents/category';
import { ObjectSpecs, TemplateJobAccess, TemplateRequirements } from '~~/gen/ts/resources/documents/templates';
import { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import { CreateTemplateRequest, UpdateTemplateRequest } from '~~/gen/ts/services/docstore/docstore';
import TemplateSchemaEditor, { type SchemaEditorValue } from '~/components/documents/templates/TemplateSchemaEditor.vue';
import type { ObjectSpecsValue } from '~/components/documents/templates/types';
import type { Template } from '~~/gen/ts/resources/documents/templates';
import SingleHint from '~/components/SingleHint.vue';

const props = defineProps<{
    templateId?: string;
}>();

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();
const notifications = useNotificatorStore();
const completorStore = useCompletorStore();

const { activeChar } = storeToRefs(authStore);

const { t } = useI18n();

const maxAccessEntries = 10;

defineRule('required', required);
defineRule('numeric', numeric);
defineRule('min', min);
defineRule('max', max);

interface FormData {
    weight: number;
    title: string;
    description: string;
    contentTitle: string;
    content: string;
    contentState: string;
}

const { handleSubmit, setValues, meta } = useForm<FormData>({
    validationSchema: {
        weight: { required: true, numeric: { min: 0, max: 4294967295 } },
        title: { required: true, min: 3, max: 255 },
        description: { required: true, min: 3, max: 255 },
        contentTitle: { required: true, min: 3, max: 2048 },
        content: { required: true, min: 3, max: 1500000 },
        contentState: { required: false, min: 0, max: 2048 },
    },
    validateOnMount: true,
    initialValues: {
        contentState: '',
    },
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await createOrUpdateTemplate(values, props.templateId).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const schema = ref<SchemaEditorValue>({
    users: {
        req: false,
        min: 0,
        max: 0,
    },

    documents: {
        req: false,
        min: 0,
        max: 0,
    },

    vehicles: {
        req: false,
        min: 0,
        max: 0,
    },
});
const access = ref<
    Map<
        string,
        {
            id: string;
            type: number;
            values: {
                job?: string;
                accessRole?: AccessLevel;
                minimumGrade?: number;
            };
        }
    >
>(new Map());

const accessTypes = [{ id: 1, name: t('common.job', 2) }];

function addDocumentAccessEntry(): void {
    if (access.value.size > maxAccessEntries - 1) {
        notifications.add({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            description: { key: 'notifications.max_access_entry.content', parameters: { max: maxAccessEntries.toString() } },
            type: 'error',
        });
        return;
    }

    const id = access.value.size > 0 ? parseInt([...access.value.keys()].pop() ?? '0') + 1 : 0;
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

function updateDocumentAccessEntryName(event: { id: string; job?: Job }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    if (event.job) {
        accessEntry.values.job = event.job.name;

        access.value.set(event.id, accessEntry);
    }
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

const contentAccess = ref<
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

const contentAccessTypes = [
    { id: 0, name: t('common.citizen', 2) },
    { id: 1, name: t('common.job', 2) },
];

function addContentDocumentAccessEntry(): void {
    if (contentAccess.value.size > maxAccessEntries - 1) {
        notifications.add({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            description: { key: 'notifications.max_access_entry.content', parameters: { max: maxAccessEntries.toString() } },
            type: 'error',
        });
        return;
    }

    const id = contentAccess.value.size > 0 ? parseInt([...contentAccess.value.keys()].pop() ?? '0') + 1 : 0;
    contentAccess.value.set(id.toString(), {
        id: id.toString(),
        type: 1,
        values: {},
        required: false,
    });
}

function removeContentDocumentAccessEntry(event: { id: string }): void {
    contentAccess.value.delete(event.id);
}

function updateContentDocumentAccessEntryType(event: { id: string; type: number }): void {
    const accessEntry = contentAccess.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.type = event.type;
    contentAccess.value.set(event.id, accessEntry);
}

function updateContentDocumentAccessEntryName(event: { id: string; job?: Job }): void {
    const accessEntry = contentAccess.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    if (event.job) {
        accessEntry.values.job = event.job.name;
        contentAccess.value.set(event.id, accessEntry);
    }
}

function updateContentDocumentAccessEntryRank(event: { id: string; rank: JobGrade }): void {
    const accessEntry = contentAccess.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.minimumGrade = event.rank.grade;
    contentAccess.value.set(event.id, accessEntry);
}

function updateContentDocumentAccessEntryAccess(event: { id: string; access: AccessLevel }): void {
    const accessEntry = contentAccess.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.accessRole = event.access;
    contentAccess.value.set(event.id, accessEntry);
}

function updateContentDocumentAccessEntryRequired(event: { id: string; required?: boolean }): void {
    const accessEntry = contentAccess.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.required = event.required;
    contentAccess.value.set(event.id, accessEntry);
}

function createObjectSpec(v: ObjectSpecsValue): ObjectSpecs {
    const o: ObjectSpecs = {
        required: v.req,
        min: v.min,
        max: v.max,
    };
    return o;
}

async function createOrUpdateTemplate(values: FormData, templateId?: string): Promise<void> {
    const tRequirements: TemplateRequirements = {
        users: createObjectSpec(schema.value.users),
        documents: createObjectSpec(schema.value.documents),
        vehicles: createObjectSpec(schema.value.vehicles),
    };

    const jobAccesses: TemplateJobAccess[] = [];
    access.value.forEach((entry) => {
        if (entry.values.accessRole === undefined) {
            return;
        }

        if (entry.type === 1) {
            if (!entry.values.job) {
                return;
            }

            jobAccesses.push({
                id: '0',
                templateId: templateId ?? '0',
                job: entry.values.job,
                minimumGrade: entry.values.minimumGrade ? entry.values.minimumGrade : 0,
                access: entry.values.accessRole,
            });
        }
    });

    const reqAccess: DocumentAccess = {
        jobs: [],
        users: [],
    };
    contentAccess.value.forEach((entry) => {
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
                required: entry.required,
            });
        } else if (entry.type === 1) {
            if (!entry.values.job) {
                return;
            }

            reqAccess.jobs.push({
                id: '0',
                documentId: '0',
                job: entry.values.job!,
                minimumGrade: entry.values.minimumGrade ? entry.values.minimumGrade : 0,
                access: entry.values.accessRole,
                required: entry.required,
            });
        }
    });

    if (typeof values.weight === 'string') {
        values.weight = parseInt(values.weight as string);
    }

    const req: CreateTemplateRequest | UpdateTemplateRequest = {
        template: {
            id: templateId ?? '0',
            weight: values.weight,
            title: values.title,
            description: values.description,
            contentTitle: values.contentTitle,
            content: values.content,
            state: values.contentState,
            schema: {
                requirements: tRequirements,
            },
            contentAccess: reqAccess,
            jobAccess: jobAccesses,
            category: selectedCategory.value,
            creatorJob: '',
        },
    };

    try {
        if (templateId === undefined) {
            const call = $grpc.getDocStoreClient().createTemplate(req);
            const { response } = await call;

            notifications.add({
                title: { key: 'notifications.templates.created.title', parameters: {} },
                description: { key: 'notifications.templates.created.title', parameters: {} },
                type: 'success',
            });

            await navigateTo({
                name: 'documents-templates-id',
                params: { id: response.id },
            });
        } else {
            const call = $grpc.getDocStoreClient().updateTemplate(req);
            const { response } = await call;
            if (response.template) {
                setValuesFromTemplate(response.template);
            }

            notifications.add({
                title: { key: 'notifications.templates.updated.title', parameters: {} },
                description: { key: 'notifications.templates.updated.content', parameters: {} },
                type: 'success',
            });
        }
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const entriesCategories = ref<Category[]>([]);
const queryCategories = ref('');
const selectedCategory = ref<Category | undefined>(undefined);

watchDebounced(queryCategories, () => findCategories(), {
    debounce: 600,
    maxWait: 1400,
});

async function findCategories(): Promise<void> {
    entriesCategories.value = await completorStore.completeDocumentCategories(queryCategories.value);
}

function setValuesFromTemplate(tpl: Template): void {
    setValues({
        weight: tpl.weight,
        title: tpl.title,
        description: tpl.description,
        contentTitle: tpl.contentTitle,
        content: tpl.content,
        contentState: tpl.state,
    });
    if (tpl.category !== undefined) {
        selectedCategory.value = tpl.category;
    }

    const tplAccess = tpl.jobAccess;
    if (tplAccess !== undefined) {
        let accessId = 0;

        tplAccess.forEach((job) => {
            const id = accessId.toString();
            access.value.set(id, {
                id,
                type: 1,
                values: {
                    job: job.job,
                    minimumGrade: job.minimumGrade,
                    accessRole: job.access,
                },
            });
            accessId++;
        });
    }

    const ctAccess = tpl.contentAccess;
    if (ctAccess !== undefined) {
        let accessId = 0n;

        ctAccess.users.forEach((access) => {
            const id = accessId.toString();
            contentAccess.value.set(id, {
                id,
                type: 0,
                values: { char: access.userId, accessRole: access.access },
                required: access.required,
            });
            accessId++;
        });

        ctAccess.jobs.forEach((access) => {
            const id = accessId.toString();
            contentAccess.value.set(id, {
                id,
                type: 1,
                values: {
                    job: access.job,
                    accessRole: access.access,
                    minimumGrade: access.minimumGrade,
                },
                required: access.required,
            });
            accessId++;
        });
    }

    schema.value.users.req = tpl.schema?.requirements?.users?.required ?? false;
    schema.value.users.min = tpl.schema?.requirements?.users?.min ?? 0;
    schema.value.users.max = tpl.schema?.requirements?.users?.max ?? 0;

    schema.value.documents.req = tpl.schema?.requirements?.documents?.required ?? false;
    schema.value.documents.min = tpl.schema?.requirements?.documents?.min ?? 0;
    schema.value.documents.max = tpl.schema?.requirements?.documents?.max ?? 0;

    schema.value.vehicles.req = tpl.schema?.requirements?.vehicles?.required ?? false;
    schema.value.vehicles.min = tpl.schema?.requirements?.vehicles?.min ?? 0;
    schema.value.vehicles.max = tpl.schema?.requirements?.vehicles?.max ?? 0;
}

onMounted(async () => {
    if (props.templateId) {
        try {
            const call = $grpc.getDocStoreClient().getTemplate({
                templateId: props.templateId,
                render: false,
            });
            const { response } = await call;

            const tpl = response.template;
            if (!tpl) {
                return;
            }

            setValuesFromTemplate(tpl);
        } catch (e) {
            $grpc.handleError(e as RpcError);
        }
    } else {
        setValues({
            weight: 0,
        });

        access.value.set('0', {
            id: '0',
            type: 1,
            values: {
                job: activeChar.value?.job,
                minimumGrade: 1,
                accessRole: AccessLevel.VIEW,
            },
        });
    }

    findCategories();
});

const { data: jobs } = useLazyAsyncData('completor-jobs', () => completorStore.listJobs());
</script>

<template>
    <div class="m-2">
        <UForm :state="{}">
            <UCard class="bg-base-800">
                <div>
                    <label for="content" class="block text-sm font-medium leading-6 text-gray-100">
                        {{ $t('common.template', 2) }} {{ $t('common.weight') }}
                    </label>
                    <div class="mt-2">
                        <VeeField
                            type="number"
                            name="weight"
                            min="0"
                            max="4294967295"
                            class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            :label="$t('common.weight')"
                            :placeholder="$t('common.weight')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </div>
                </div>
                <div>
                    <label for="title" class="mt-2 block text-sm font-medium">
                        {{ $t('common.template') }} {{ $t('common.title') }}
                    </label>
                    <div>
                        <VeeField
                            as="textarea"
                            rows="1"
                            name="title"
                            class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            :label="$t('common.title')"
                            :placeholder="$t('common.title')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="title" as="p" class="mt-2 text-sm text-error-400" />
                    </div>
                </div>
                <div>
                    <label for="description" class="mt-2 block text-sm font-medium">
                        {{ $t('common.template') }} {{ $t('common.description') }}
                    </label>
                    <div>
                        <VeeField
                            as="textarea"
                            rows="4"
                            name="description"
                            :label="$t('common.description')"
                            class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
                    </div>
                </div>
                <div>
                    <div class="my-3">
                        <h2 class="text-sm">{{ $t('common.template') }} {{ $t('common.access') }}</h2>
                        <DocumentAccessEntry
                            v-for="entry in access.values()"
                            :key="entry.id"
                            :init="entry"
                            :access-types="accessTypes"
                            :access-roles="[AccessLevel.VIEW, AccessLevel.EDIT]"
                            :jobs="jobs"
                            @type-change="updateDocumentAccessEntryType($event)"
                            @name-change="updateDocumentAccessEntryName($event)"
                            @rank-change="updateDocumentAccessEntryRank($event)"
                            @access-change="updateDocumentAccessEntryAccess($event)"
                            @delete-request="removeDocumentAccessEntry($event)"
                        />
                        <UButton
                            icon="i-mdi-plus"
                            data-te-toggle="tooltip"
                            :title="$t('components.documents.document_editor.add_permission')"
                            @click="addDocumentAccessEntry()"
                        />
                    </div>
                </div>
            </UCard>

            <div class="mt-2 flex flex-col sm:flex-row">
                <SingleHint
                    class="min-w-full"
                    hint-id="template_editor_templating"
                    to="https://github.com/galexrt/fivenet/blob/main/docs/features/documents_templates.md"
                    :external="true"
                    link-target="_blank"
                />
            </div>

            <UCard class="bg-base-800">
                <div>
                    <label for="contentTitle" class="mt-2 block text-sm font-medium">
                        {{ $t('common.content') }} {{ $t('common.title') }}
                    </label>
                    <div>
                        <VeeField
                            as="textarea"
                            rows="2"
                            name="contentTitle"
                            :label="$t('common.title')"
                            class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="contentTitle" as="p" class="mt-2 text-sm text-error-400" />
                    </div>
                </div>

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

                <div>
                    <label for="contentState" class="mt-2 block text-sm font-medium">
                        {{ $t('common.content') }} {{ $t('common.state') }}
                    </label>
                    <div>
                        <VeeField
                            as="textarea"
                            rows="2"
                            name="contentState"
                            :label="$t('common.state')"
                            class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="contentState" as="p" class="mt-2 text-sm text-error-400" />
                    </div>
                </div>
                <div>
                    <label for="content" class="mt-2 block text-sm font-medium">
                        {{ $t('common.content') }} {{ $t('common.template') }}
                    </label>
                    <div>
                        <VeeField
                            as="textarea"
                            rows="6"
                            name="content"
                            :label="$t('common.template')"
                            class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="content" as="p" class="mt-2 text-sm text-error-400" />
                    </div>
                </div>
                <TemplateSchemaEditor v-model="schema" class="mt-2" />
                <div class="my-3">
                    <h2 class="text-sm">{{ $t('common.content') }} {{ $t('common.access') }}</h2>
                    <DocumentAccessEntry
                        v-for="entry in contentAccess.values()"
                        :key="entry.id"
                        :init="entry"
                        :access-types="contentAccessTypes"
                        :show-required="true"
                        :jobs="jobs"
                        @type-change="updateContentDocumentAccessEntryType($event)"
                        @name-change="updateContentDocumentAccessEntryName($event)"
                        @rank-change="updateContentDocumentAccessEntryRank($event)"
                        @access-change="updateContentDocumentAccessEntryAccess($event)"
                        @delete-request="removeContentDocumentAccessEntry($event)"
                        @required-change="updateContentDocumentAccessEntryRequired($event)"
                    />
                    <UButton
                        data-te-toggle="tooltip"
                        icon="i-mdi-plus"
                        :title="$t('components.documents.document_editor.add_permission')"
                        @click="addContentDocumentAccessEntry()"
                    />
                </div>
            </UCard>

            <div>
                <UButton
                    class="mt-4 flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                    :disabled="!meta.valid || !canSubmit"
                    :loading="!canSubmit"
                >
                    {{ templateId ? $t('common.save') : $t('common.create') }}
                </UButton>
            </div>
        </UForm>
    </div>
</template>
