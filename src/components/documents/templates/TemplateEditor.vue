<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { max, min, numeric, required } from '@vee-validate/rules';
import { useThrottleFn, watchDebounced } from '@vueuse/core';
import { CheckIcon, LoadingIcon, PlusIcon } from 'mdi-vue3';
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
import TemplateSyntaxHint from '~/components/documents/templates/partials/TemplateSyntaxHint.vue';

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
        await createOrUpdateTemplate(values, props.templateId).finally(() => setTimeout(() => (canSubmit.value = true), 400)),
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
        notifications.dispatchNotification({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            content: { key: 'notifications.max_access_entry.content', parameters: { max: maxAccessEntries.toString() } },
            type: 'error',
        });
        return;
    }

    const id = access.value.size > 0 ? parseInt([...access.value.keys()].pop() ?? '0', 10) + 1 : 0;
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
        notifications.dispatchNotification({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            content: { key: 'notifications.max_access_entry.content', parameters: { max: maxAccessEntries.toString() } },
            type: 'error',
        });
        return;
    }

    const id = contentAccess.value.size > 0 ? parseInt([...contentAccess.value.keys()].pop() ?? '0', 10) + 1 : 0;
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
            weight: values.weight as number,
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
        },
    };

    try {
        if (templateId === undefined) {
            const call = $grpc.getDocStoreClient().createTemplate(req);
            const { response } = await call;

            notifications.dispatchNotification({
                title: { key: 'notifications.templates.created.title', parameters: {} },
                content: { key: 'notifications.templates.created.title', parameters: {} },
                type: 'success',
            });

            await navigateTo({
                name: 'documents-templates-id',
                params: { id: response.id },
            });
        } else {
            const call = $grpc.getDocStoreClient().updateTemplate(req);
            await call;

            notifications.dispatchNotification({
                title: { key: 'notifications.templates.updated.title', parameters: {} },
                content: { key: 'notifications.templates.updated.content', parameters: {} },
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
</script>

<template>
    <div class="m-2 text-neutral">
        <form @submit.prevent="onSubmitThrottle">
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
                        :label="$t('common.weight')"
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                        :label="$t('common.title')"
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                    <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
                </div>
            </div>
            <div>
                <div class="my-3">
                    <h2 class="text-neutral">{{ $t('common.template') }} {{ $t('common.access') }}</h2>
                    <DocumentAccessEntry
                        v-for="entry in access.values()"
                        :key="entry.id"
                        :init="entry"
                        :access-types="accessTypes"
                        :access-roles="[AccessLevel.VIEW, AccessLevel.EDIT]"
                        @type-change="updateDocumentAccessEntryType($event)"
                        @name-change="updateDocumentAccessEntryName($event)"
                        @rank-change="updateDocumentAccessEntryRank($event)"
                        @access-change="updateDocumentAccessEntryAccess($event)"
                        @delete-request="removeDocumentAccessEntry($event)"
                    />
                    <button
                        type="button"
                        class="rounded-full bg-primary-500 p-2 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                        data-te-toggle="tooltip"
                        :title="$t('components.documents.document_editor.add_permission')"
                        @click="addDocumentAccessEntry()"
                    >
                        <PlusIcon class="h-5 w-5" aria-hidden="true" />
                    </button>
                </div>
            </div>
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
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                    <VeeErrorMessage name="contentTitle" as="p" class="mt-2 text-sm text-error-400" />
                    <TemplateSyntaxHint />
                </div>
            </div>
            <div>
                <label for="contentCategory" class="mt-2 block text-sm font-medium">
                    {{ $t('common.category') }}
                </label>
                <div>
                    <Combobox v-model="selectedCategory" as="div" nullable>
                        <div class="relative">
                            <ComboboxButton as="div">
                                <ComboboxInput
                                    autocomplete="off"
                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    :display-value="(category: any) => category?.name"
                                    @change="queryCategories = $event.target.value"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </ComboboxButton>

                            <ComboboxOptions
                                v-if="entriesCategories.length > 0"
                                class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                            >
                                <ComboboxOption
                                    v-for="category in entriesCategories"
                                    v-slot="{ active, selected }"
                                    :key="category.id"
                                    :value="category"
                                    as="category"
                                >
                                    <li
                                        :class="[
                                            'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                            active ? 'bg-primary-500' : '',
                                        ]"
                                    >
                                        <span :class="['block truncate', selected && 'font-semibold']">
                                            {{ category.name }}
                                        </span>

                                        <span
                                            v-if="selected"
                                            :class="[
                                                active ? 'text-neutral' : 'text-primary-500',
                                                'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                            ]"
                                        >
                                            <CheckIcon class="h-5 w-5" aria-hidden="true" />
                                        </span>
                                    </li>
                                </ComboboxOption>
                            </ComboboxOptions>
                        </div>
                    </Combobox>
                </div>
            </div>
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
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    />
                    <VeeErrorMessage name="contentState" as="p" class="mt-2 text-sm text-error-400" />
                    <TemplateSyntaxHint />
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
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    />
                    <VeeErrorMessage name="content" as="p" class="mt-2 text-sm text-error-400" />
                    <TemplateSyntaxHint />
                </div>
            </div>
            <TemplateSchemaEditor v-model="schema" class="mt-2" />
            <div class="my-3">
                <h2 class="text-neutral">{{ $t('common.content') }} {{ $t('common.access') }}</h2>
                <DocumentAccessEntry
                    v-for="entry in contentAccess.values()"
                    :key="entry.id"
                    :init="entry"
                    :access-types="contentAccessTypes"
                    :show-required="true"
                    @type-change="updateContentDocumentAccessEntryType($event)"
                    @name-change="updateContentDocumentAccessEntryName($event)"
                    @rank-change="updateContentDocumentAccessEntryRank($event)"
                    @access-change="updateContentDocumentAccessEntryAccess($event)"
                    @delete-request="removeContentDocumentAccessEntry($event)"
                    @required-change="updateContentDocumentAccessEntryRequired($event)"
                />
                <button
                    type="button"
                    class="rounded-full bg-primary-500 p-2 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    data-te-toggle="tooltip"
                    :title="$t('components.documents.document_editor.add_permission')"
                    @click="addContentDocumentAccessEntry()"
                >
                    <PlusIcon class="h-5 w-5" aria-hidden="true" />
                </button>
            </div>
            <div>
                <button
                    type="submit"
                    class="mt-4 flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold text-neutral transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                    :disabled="!meta.valid || !canSubmit"
                    :class="[
                        !meta.valid || !canSubmit
                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                            : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                    ]"
                >
                    <template v-if="!canSubmit">
                        <LoadingIcon class="mr-2 h-5 w-5 animate-spin" aria-hidden="true" />
                    </template>
                    {{ templateId ? $t('common.save') : $t('common.create') }}
                </button>
            </div>
        </form>
    </div>
</template>
