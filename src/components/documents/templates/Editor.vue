<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { max, min, numeric, required } from '@vee-validate/rules';
import { useThrottleFn, watchDebounced } from '@vueuse/core';
import { CheckIcon, LoadingIcon, PlusIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import AccessEntry from '~/components/documents/AccessEntry.vue';
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import { Category } from '~~/gen/ts/resources/documents/category';
import { DocumentAccess } from '~~/gen/ts/resources/documents/documents';
import { ObjectSpecs, TemplateJobAccess, TemplateRequirements } from '~~/gen/ts/resources/documents/templates';
import { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import { CreateTemplateRequest, UpdateTemplateRequest } from '~~/gen/ts/services/docstore/docstore';
import SchemaEditor, { ObjectSpecsValue, SchemaEditorValue } from './SchemaEditor.vue';
import TemplateHint from './partials/TemplateHint.vue';

const props = defineProps<{
    templateId?: bigint;
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
        await createOrUpdateTemplate(values, props.templateId).finally(() => setTimeout(() => (canSubmit.value = true), 350)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const schema = ref<SchemaEditorValue>({
    users: {
        req: false,
        min: 0n,
        max: 0n,
    },

    documents: {
        req: false,
        min: 0n,
        max: 0n,
    },

    vehicles: {
        req: false,
        min: 0n,
        max: 0n,
    },
});
const access = ref<
    Map<
        bigint,
        {
            id: bigint;
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

function addAccessEntry(): void {
    if (access.value.size > maxAccessEntries - 1) {
        notifications.dispatchNotification({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            content: { key: 'notifications.max_access_entry.content', parameters: { max: maxAccessEntries.toString() } },
            type: 'error',
        });
        return;
    }

    const id = access.value.size > 0 ? ([...access.value.keys()].pop() as bigint) + 1n : 0n;
    access.value.set(id, {
        id,
        type: 1,
        values: {},
    });
}

function removeAccessEntry(event: { id: bigint }): void {
    access.value.delete(event.id);
}

function updateAccessEntryType(event: { id: bigint; type: number }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.type = event.type;
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryName(event: { id: bigint; job?: Job }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    if (event.job) {
        accessEntry.values.job = event.job.name;

        access.value.set(event.id, accessEntry);
    }
}

function updateAccessEntryRank(event: { id: bigint; rank: JobGrade }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.minimumGrade = event.rank.grade;
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryAccess(event: { id: bigint; access: AccessLevel }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.accessRole = event.access;
    access.value.set(event.id, accessEntry);
}

const contentAccess = ref<
    Map<
        bigint,
        {
            id: bigint;
            type: number;
            values: {
                job?: string;
                char?: number;
                accessRole?: AccessLevel;
                minimumGrade?: number;
            };
        }
    >
>(new Map());

const contentAccessTypes = [
    { id: 0, name: t('common.citizen', 2) },
    { id: 1, name: t('common.job', 2) },
];

function addContentAccessEntry(): void {
    if (contentAccess.value.size > maxAccessEntries - 1) {
        notifications.dispatchNotification({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            content: { key: 'notifications.max_access_entry.content', parameters: { max: maxAccessEntries.toString() } },
            type: 'error',
        });
        return;
    }

    const id = contentAccess.value.size > 0 ? ([...contentAccess.value.keys()].pop() as bigint) + 1n : 0n;
    contentAccess.value.set(id, {
        id,
        type: 1,
        values: {},
    });
}

function removeContentAccessEntry(event: { id: bigint }): void {
    contentAccess.value.delete(event.id);
}

function updateContentAccessEntryType(event: { id: bigint; type: number }): void {
    const accessEntry = contentAccess.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.type = event.type;
    contentAccess.value.set(event.id, accessEntry);
}

function updateContentAccessEntryName(event: { id: bigint; job?: Job }): void {
    const accessEntry = contentAccess.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    if (event.job) {
        accessEntry.values.job = event.job.name;

        contentAccess.value.set(event.id, accessEntry);
    }
}

function updateContentAccessEntryRank(event: { id: bigint; rank: JobGrade }): void {
    const accessEntry = contentAccess.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.minimumGrade = event.rank.grade;
    contentAccess.value.set(event.id, accessEntry);
}

function updateContentAccessEntryAccess(event: { id: bigint; access: AccessLevel }): void {
    const accessEntry = contentAccess.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.accessRole = event.access;
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

async function createOrUpdateTemplate(values: FormData, templateId?: bigint): Promise<void> {
    return new Promise(async (res, rej) => {
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
                    id: 0n,
                    templateId: templateId ?? 0n,
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
                    id: 0n,
                    documentId: 0n,
                    userId: entry.values.char,
                    access: entry.values.accessRole,
                });
            } else if (entry.type === 1) {
                if (!entry.values.job) {
                    return;
                }

                reqAccess.jobs.push({
                    id: 0n,
                    documentId: 0n,
                    job: entry.values.job!,
                    minimumGrade: entry.values.minimumGrade ? entry.values.minimumGrade : 0,
                    access: entry.values.accessRole,
                });
            }
        });

        if (typeof values.weight === 'string') values.weight = parseInt(values.weight as string);

        const req: CreateTemplateRequest | UpdateTemplateRequest = {
            template: {
                id: templateId ?? 0n,
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
                    params: { id: response.id?.toString() },
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

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
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
            if (tpl.category) {
                selectedCategory.value = tpl.category;
            }

            const tplAccess = tpl.jobAccess;
            if (tplAccess) {
                let accessId = 0n;

                tplAccess.forEach((job) => {
                    access.value.set(accessId, {
                        id: accessId,
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

            const ctAccess = tpl.contentAccess;
            if (ctAccess) {
                let accessId = 0n;

                ctAccess.users.forEach((user) => {
                    contentAccess.value.set(accessId, {
                        id: accessId,
                        type: 0,
                        values: { char: user.userId, accessRole: user.access },
                    });
                    accessId++;
                });

                ctAccess.jobs.forEach((job) => {
                    contentAccess.value.set(accessId, {
                        id: accessId,
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

            schema.value.users.req = tpl.schema?.requirements?.users?.required ?? false;
            schema.value.users.min = tpl.schema?.requirements?.users?.min ?? 0n;
            schema.value.users.max = tpl.schema?.requirements?.users?.max ?? 0n;

            schema.value.documents.req = tpl.schema?.requirements?.documents?.required ?? false;
            schema.value.documents.min = tpl.schema?.requirements?.documents?.min ?? 0n;
            schema.value.documents.max = tpl.schema?.requirements?.documents?.max ?? 0n;

            schema.value.vehicles.req = tpl.schema?.requirements?.vehicles?.required ?? false;
            schema.value.vehicles.min = tpl.schema?.requirements?.vehicles?.min ?? 0n;
            schema.value.vehicles.max = tpl.schema?.requirements?.vehicles?.max ?? 0n;
        } catch (e) {
            $grpc.handleError(e as RpcError);
        }
    } else {
        setValues({
            weight: 0,
        });

        access.value.set(0n, {
            id: 0n,
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
    <div class="text-neutral">
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
                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    />
                </div>
            </div>
            <div>
                <label for="title" class="block font-medium text-sm mt-2">
                    {{ $t('common.template') }} {{ $t('common.title') }}
                </label>
                <div>
                    <VeeField
                        as="textarea"
                        rows="1"
                        name="title"
                        :label="$t('common.title')"
                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    />
                    <VeeErrorMessage name="title" as="p" class="mt-2 text-sm text-error-400" />
                </div>
            </div>
            <div>
                <label for="description" class="block font-medium text-sm mt-2">
                    {{ $t('common.template') }} {{ $t('common.description') }}
                </label>
                <div>
                    <VeeField
                        as="textarea"
                        rows="4"
                        name="description"
                        :label="$t('common.description')"
                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    />
                    <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
                </div>
            </div>
            <div>
                <div class="my-3">
                    <h2 class="text-neutral">{{ $t('common.template') }} {{ $t('common.access') }}</h2>
                    <AccessEntry
                        v-for="entry in access.values()"
                        :key="entry.id?.toString()"
                        :init="entry"
                        :access-types="accessTypes"
                        :access-roles="[AccessLevel.VIEW, AccessLevel.EDIT]"
                        @typeChange="updateAccessEntryType($event)"
                        @nameChange="updateAccessEntryName($event)"
                        @rankChange="updateAccessEntryRank($event)"
                        @accessChange="updateAccessEntryAccess($event)"
                        @deleteRequest="removeAccessEntry($event)"
                    />
                    <button
                        type="button"
                        class="p-2 rounded-full bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                        data-te-toggle="tooltip"
                        :title="$t('components.documents.document_editor.add_permission')"
                        @click="addAccessEntry()"
                    >
                        <PlusIcon class="w-5 h-5" aria-hidden="true" />
                    </button>
                </div>
            </div>
            <div>
                <label for="contentTitle" class="block font-medium text-sm mt-2">
                    {{ $t('common.content') }} {{ $t('common.title') }}
                </label>
                <div>
                    <VeeField
                        as="textarea"
                        rows="2"
                        name="contentTitle"
                        :label="$t('common.title')"
                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    />
                    <VeeErrorMessage name="contentTitle" as="p" class="mt-2 text-sm text-error-400" />
                    <TemplateHint />
                </div>
            </div>
            <div>
                <label for="contentCategory" class="block font-medium text-sm mt-2">
                    {{ $t('common.category') }}
                </label>
                <div>
                    <Combobox as="div" v-model="selectedCategory" nullable>
                        <div class="relative">
                            <ComboboxButton as="div">
                                <ComboboxInput
                                    autocomplete="off"
                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    @change="queryCategories = $event.target.value"
                                    :display-value="(category: any) => category?.name"
                                />
                            </ComboboxButton>

                            <ComboboxOptions
                                v-if="entriesCategories.length > 0"
                                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                            >
                                <ComboboxOption
                                    v-for="category in entriesCategories"
                                    :key="category.id?.toString()"
                                    :value="category"
                                    as="category"
                                    v-slot="{ active, selected }"
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
                                            <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                        </span>
                                    </li>
                                </ComboboxOption>
                            </ComboboxOptions>
                        </div>
                    </Combobox>
                </div>
            </div>
            <div>
                <label for="contentState" class="block font-medium text-sm mt-2">
                    {{ $t('common.content') }} {{ $t('common.state') }}
                </label>
                <div>
                    <VeeField
                        as="textarea"
                        rows="2"
                        name="contentState"
                        :label="$t('common.state')"
                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    />
                    <VeeErrorMessage name="contentState" as="p" class="mt-2 text-sm text-error-400" />
                    <TemplateHint />
                </div>
            </div>
            <div>
                <label for="content" class="block font-medium text-sm mt-2">
                    {{ $t('common.content') }} {{ $t('common.template') }}
                </label>
                <div>
                    <VeeField
                        as="textarea"
                        rows="6"
                        name="content"
                        :label="$t('common.template')"
                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    />
                    <VeeErrorMessage name="content" as="p" class="mt-2 text-sm text-error-400" />
                    <TemplateHint />
                </div>
            </div>
            <SchemaEditor v-model="schema" class="mt-2" />
            <div class="my-3">
                <h2 class="text-neutral">{{ $t('common.content') }} {{ $t('common.access') }}</h2>
                <AccessEntry
                    v-for="entry in contentAccess.values()"
                    :key="entry.id?.toString()"
                    :init="entry"
                    :access-types="contentAccessTypes"
                    @typeChange="updateContentAccessEntryType($event)"
                    @nameChange="updateContentAccessEntryName($event)"
                    @rankChange="updateContentAccessEntryRank($event)"
                    @accessChange="updateContentAccessEntryAccess($event)"
                    @deleteRequest="removeContentAccessEntry($event)"
                />
                <button
                    type="button"
                    class="p-2 rounded-full bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    data-te-toggle="tooltip"
                    :title="$t('components.documents.document_editor.add_permission')"
                    @click="addContentAccessEntry()"
                >
                    <PlusIcon class="w-5 h-5" aria-hidden="true" />
                </button>
            </div>
            <div>
                <button
                    type="submit"
                    class="flex justify-center w-full mt-4 px-3 py-2 text-sm font-semibold transition-colors rounded-md text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                    :disabled="!meta.valid || !canSubmit"
                    :class="[
                        !meta.valid || !canSubmit
                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                            : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                    ]"
                >
                    <template v-if="!canSubmit">
                        <LoadingIcon class="animate-spin h-5 w-5 mr-2" />
                    </template>
                    {{ templateId ? $t('common.save') : $t('common.create') }}
                </button>
            </div>
        </form>
    </div>
</template>
