<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import SingleHint from '~/components/SingleHint.vue';
import DocumentAccessEntry from '~/components/documents/DocumentAccessEntry.vue';
import TemplateSchemaEditor, { type SchemaEditorValue } from '~/components/documents/templates/TemplateSchemaEditor.vue';
import type { ObjectSpecsValue } from '~/components/documents/templates/types';
import DocEditor from '~/components/partials/DocEditor.vue';
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import type { DocumentAccess } from '~~/gen/ts/resources/documents/access';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { ObjectSpecs, Template, TemplateJobAccess, TemplateRequirements } from '~~/gen/ts/resources/documents/templates';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import type { CreateTemplateRequest, UpdateTemplateRequest } from '~~/gen/ts/services/docstore/docstore';

const props = defineProps<{
    templateId?: string;
}>();

const authStore = useAuthStore();

const notifications = useNotificatorStore();

const completorStore = useCompletorStore();

const { activeChar } = storeToRefs(authStore);

const { t } = useI18n();

const maxAccessEntries = 10;

const schema = z.object({
    weight: z.number().min(0).max(999_999),
    title: z.string().min(3).max(255),
    description: z.string().min(3).max(255),
    contentTitle: z.string().min(3).max(2048),
    content: z.string().min(3).max(1500000),
    contentState: z.union([z.string().min(1).max(2048), z.string().length(0)]),
    category: z.custom<Category>().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    weight: 0,
    title: '',
    description: '',
    contentTitle: '',
    content: '',
    contentState: '',
    category: undefined,
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateTemplate(event.data, props.templateId).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const schemaEditor = ref<SchemaEditorValue>({
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
const access = ref(
    new Map<
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
    >(),
);

const accessTypes = [{ id: 1, name: t('common.job', 2) }];

function addDocumentAccessEntry(): void {
    if (access.value.size > maxAccessEntries - 1) {
        notifications.add({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            description: { key: 'notifications.max_access_entry.content', parameters: { max: maxAccessEntries.toString() } },
            type: NotificationType.ERROR,
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

const contentAccess = ref(
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

const contentAccessTypes = [
    { id: 0, name: t('common.citizen', 2) },
    { id: 1, name: t('common.job', 2) },
];

function addContentDocumentAccessEntry(): void {
    if (contentAccess.value.size > maxAccessEntries - 1) {
        notifications.add({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            description: { key: 'notifications.max_access_entry.content', parameters: { max: maxAccessEntries.toString() } },
            type: NotificationType.ERROR,
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

async function createOrUpdateTemplate(values: Schema, templateId?: string): Promise<void> {
    const tRequirements: TemplateRequirements = {
        users: createObjectSpec(schemaEditor.value.users),
        documents: createObjectSpec(schemaEditor.value.documents),
        vehicles: createObjectSpec(schemaEditor.value.vehicles),
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
            if (!entry.values.userId) {
                return;
            }

            reqAccess.users.push({
                id: '0',
                documentId: '0',
                userId: entry.values.userId,
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
            category: state.category,
            creatorJob: '',
        },
    };

    try {
        if (templateId === undefined) {
            const call = getGRPCDocStoreClient().createTemplate(req);
            const { response } = await call;

            notifications.add({
                title: { key: 'notifications.templates.created.title', parameters: {} },
                description: { key: 'notifications.templates.created.title', parameters: {} },
                type: NotificationType.SUCCESS,
            });

            await navigateTo({
                name: 'documents-templates-id',
                params: { id: response.id },
            });
        } else {
            const call = getGRPCDocStoreClient().updateTemplate(req);
            const { response } = await call;
            if (response.template) {
                setValuesFromTemplate(response.template);
            }

            notifications.add({
                title: { key: 'notifications.templates.updated.title', parameters: {} },
                description: { key: 'notifications.templates.updated.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });
        }
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const entriesCategories = ref<Category[]>([]);
const queryCategories = ref('');

watchDebounced(queryCategories, () => findCategories(), {
    debounce: 200,
    maxWait: 1250,
});

async function findCategories(): Promise<void> {
    entriesCategories.value = await completorStore.completeDocumentCategories(queryCategories.value);
}

function setValuesFromTemplate(tpl: Template): void {
    state.weight = tpl.weight;
    state.title = tpl.title;
    state.description = tpl.description;
    state.contentTitle = tpl.contentTitle;
    state.content = tpl.content;
    state.contentState = tpl.state;
    state.category = tpl.category;

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
        let accessId = 0;

        ctAccess.users.forEach((access) => {
            const id = accessId.toString();
            contentAccess.value.set(id, {
                id,
                type: 0,
                values: { userId: access.userId, accessRole: access.access },
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

    schemaEditor.value.users.req = tpl.schema?.requirements?.users?.required ?? false;
    schemaEditor.value.users.min = tpl.schema?.requirements?.users?.min ?? 0;
    schemaEditor.value.users.max = tpl.schema?.requirements?.users?.max ?? 0;

    schemaEditor.value.documents.req = tpl.schema?.requirements?.documents?.required ?? false;
    schemaEditor.value.documents.min = tpl.schema?.requirements?.documents?.min ?? 0;
    schemaEditor.value.documents.max = tpl.schema?.requirements?.documents?.max ?? 0;

    schemaEditor.value.vehicles.req = tpl.schema?.requirements?.vehicles?.required ?? false;
    schemaEditor.value.vehicles.min = tpl.schema?.requirements?.vehicles?.min ?? 0;
    schemaEditor.value.vehicles.max = tpl.schema?.requirements?.vehicles?.max ?? 0;
}

onMounted(async () => {
    if (props.templateId) {
        try {
            const call = getGRPCDocStoreClient().getTemplate({
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
            handleGRPCError(e as RpcError);
        }
    } else {
        state.weight = 0;

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

const categoriesLoading = ref(false);

const { data: jobs } = useAsyncData('completor-jobs', () => completorStore.listJobs());
</script>

<template>
    <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
        <UDashboardNavbar :title="$t('pages.documents.templates.edit.title')">
            <template #right>
                <UButton
                    color="black"
                    icon="i-mdi-arrow-left"
                    :to="templateId ? { name: 'documents-templates-id', params: { id: templateId } } : `/documents/templates`"
                >
                    {{ $t('common.back') }}
                </UButton>

                <UButtonGroup class="inline-flex">
                    <UButton type="submit" trailing-icon="i-mdi-content-save" :disabled="!canSubmit" :loading="!canSubmit">
                        {{ templateId ? $t('common.save') : $t('common.create') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UDashboardNavbar>

        <UContainer class="w-full">
            <div>
                <UFormGroup name="weight" :label="`${$t('common.template', 1)} ${$t('common.weight')}`">
                    <UInput
                        v-model="state.weight"
                        type="number"
                        name="weight"
                        :min="0"
                        :max="999999"
                        :placeholder="$t('common.weight')"
                    />
                </UFormGroup>

                <UFormGroup name="title" :label="`${$t('common.template')} ${$t('common.title')}`" required>
                    <UTextarea v-model="state.title" name="title" :rows="1" :placeholder="$t('common.title')" />
                </UFormGroup>

                <UFormGroup name="description" :label="`${$t('common.template')} ${$t('common.description')}`" required>
                    <UTextarea v-model="state.description" name="description" :rows="4" :label="$t('common.description')" />
                </UFormGroup>
            </div>

            <div class="my-2">
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
                    :ui="{ rounded: 'rounded-full' }"
                    :title="$t('components.documents.document_editor.add_permission')"
                    @click="addDocumentAccessEntry()"
                />
            </div>

            <SingleHint
                class="my-2"
                hint-id="template_editor_templating"
                to="https://fivenet.app/user-guides/documents/templates"
                :external="true"
                link-target="_blank"
            />

            <div>
                <UFormGroup name="contentTitle" :label="`${$t('common.content')} ${$t('common.title')}`" required>
                    <UTextarea v-model="state.contentTitle" name="contentTitle" :rows="2" />
                </UFormGroup>

                <UFormGroup name="category" :label="$t('common.category', 1)">
                    <ClientOnly>
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
                    </ClientOnly>
                </UFormGroup>

                <UFormGroup name="contentState" :label="`${$t('common.content')} ${$t('common.state')}`">
                    <UTextarea v-model="state.contentState" name="contentState" :rows="2" />
                </UFormGroup>

                <UFormGroup name="content" :label="`${$t('common.content')} ${$t('common.template')}`" required>
                    <ClientOnly>
                        <DocEditor v-model="state.content" split-screen />
                    </ClientOnly>
                </UFormGroup>
            </div>

            <div class="my-2">
                <h2 class="text-sm">{{ $t('common.requirements', 2) }}</h2>

                <TemplateSchemaEditor v-model="schemaEditor" class="mt-2" />
            </div>

            <div class="my-2">
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
                    icon="i-mdi-plus"
                    :ui="{ rounded: 'rounded-full' }"
                    :title="$t('components.documents.document_editor.add_permission')"
                    @click="addContentDocumentAccessEntry()"
                />
            </div>
        </UContainer>
    </UForm>
</template>
