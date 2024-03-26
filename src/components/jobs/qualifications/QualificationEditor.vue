<script lang="ts" setup>
import { Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { CheckIcon, ChevronDownIcon, LoadingIcon, PlusIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useNotificatorStore } from '~/store/notificator';
import {
    AccessLevel,
    QualificationAccess,
    type Qualification,
    QualificationRequirement,
    QualificationShort,
} from '~~/gen/ts/resources/qualifications/qualifications';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import type {
    CreateQualificationResponse,
    UpdateQualificationResponse,
} from '~~/gen/ts/services/qualifications/qualifications';
import QualificationAccessEntry from '~/components/jobs/qualifications/QualificationAccessEntry.vue';
import QualificationRequirementEntry from '~/components/jobs/qualifications/QualificationRequirementEntry.vue';
import DocEditor from '~/components/partials/DocEditor.vue';
import { useAuthStore } from '~/store/auth';

const props = defineProps<{
    id?: string;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const { t } = useI18n();

const maxAccessEntries = 10;

const canEdit = ref(true);

const canDo = computed(() => ({
    edit: true,
    access: true,
}));

interface FormData {
    weight: number;
    abbreviation: string;
    title: string;
    description: string;
    content: string;
}

const openclose = [
    { id: 0, label: t('common.open', 2), closed: false },
    { id: 1, label: t('common.close', 2), closed: true },
];

const quali = ref<{
    closed: { id: number; label: string; closed: boolean };
    requirements: QualificationRequirement[];
}>({
    closed: openclose[0],
    requirements: [],
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
const qualiAccess = ref<QualificationAccess>();

onMounted(async () => {
    if (props.id) {
        try {
            const call = $grpc.getQualificationsClient().getQualification({
                qualificationId: props.id,
            });
            const { response } = await call;

            const qualification = response.qualification;
            qualiAccess.value = response.qualification?.access;

            if (qualification) {
                setFieldValue('abbreviation', qualification.abbreviation);
                setFieldValue('title', qualification.title);
                if (qualification.description) {
                    setFieldValue('description', qualification.description);
                }
                setFieldValue('content', qualification.content);
                quali.value.closed = openclose.find((e) => e.closed === qualification.closed) as {
                    id: number;
                    label: string;
                    closed: boolean;
                };
                quali.value.requirements = qualification.requirements;
            }

            if (response.qualification?.access) {
                let accessId = 0;

                response.qualification?.access.jobs.forEach((job) => {
                    const id = accessId.toString();
                    access.value.set(id, {
                        id,
                        type: 0,
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

            await navigateTo({ name: 'jobs-qualifications' });

            return;
        }
    } else {
        const accessId = 0;
        access.value.set(accessId.toString(), {
            id: accessId.toString(),
            type: 0,
            values: {
                job: activeChar.value?.job,
                minimumGrade: 1,
                accessRole: AccessLevel.EDIT,
            },
        });
    }

    canEdit.value = true;
});

async function createQualification(values: FormData): Promise<CreateQualificationResponse> {
    const req = {
        qualification: {
            id: '0',
            job: '',
            weight: 0,
            closed: quali.value.closed.closed,
            abbreviation: values.abbreviation,
            title: values.title,
            description: values.description,
            content: values.content,
            creatorId: 0,
            creatorJob: '',
            requirements: quali.value.requirements,
            access: {
                jobs: [],
            } as QualificationAccess,
        },
    };
    access.value.forEach((entry) => {
        if (entry.values.accessRole === undefined) {
            return;
        }

        if (entry.type === 0) {
            if (!entry.values.job) {
                return;
            }

            req.qualification.access.jobs.push({
                id: '0',
                qualificationId: '0',
                job: entry.values.job,
                minimumGrade: entry.values.minimumGrade ? entry.values.minimumGrade : 0,
                access: entry.values.accessRole,
            });
        }
    });

    try {
        const call = $grpc.getQualificationsClient().createQualification(req);
        const { response } = await call;

        await navigateTo({
            name: 'jobs-qualifications-id',
            params: { id: response.qualificationId },
        });

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function updateQualification(values: FormData): Promise<UpdateQualificationResponse> {
    const req = {
        qualification: {
            id: props.id!,
            job: '',
            weight: 0,
            closed: quali.value.closed.closed,
            abbreviation: values.abbreviation,
            title: values.title,
            description: values.description,
            content: values.content,
            creatorId: 0,
            creatorJob: '',
            requirements: quali.value.requirements,
            access: {
                jobs: [],
            } as QualificationAccess,
        },
    };
    access.value.forEach((entry) => {
        if (entry.values.accessRole === undefined) {
            return;
        }

        if (entry.type === 0) {
            if (!entry.values.job) {
                return;
            }

            req.qualification.access.jobs.push({
                id: '0',
                qualificationId: '0',
                job: entry.values.job,
                minimumGrade: entry.values.minimumGrade ? entry.values.minimumGrade : 0,
                access: entry.values.accessRole,
            });
        }
    });

    try {
        const call = $grpc.getQualificationsClient().updateQualification(req);

        const { response } = await call;

        await navigateTo({
            name: 'jobs-qualifications-id',
            params: { id: response.qualificationId },
        });

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta, setFieldValue } = useForm<FormData>({
    validationSchema: {
        weight: {},
        abbreviation: { required: true, min: 3, max: 20 },
        title: { required: true, min: 3, max: 1024 },
        description: { required: true, max: 512 },
        content: { required: true, min: 20, max: 750000 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(async (values): Promise<any> => {
    if (props.id === undefined) {
        await createQualification(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
    } else {
        await updateQualification(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
    }
});
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const accessTypes = [{ id: 0, name: t('common.job', 2) }];

function addQualificationAccessEntry(): void {
    if (access.value.size > maxAccessEntries - 1) {
        notifications.dispatchNotification({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            content: {
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
        type: 0,
        values: {},
    });
}

function removeQualificationAccessEntry(event: { id: string }): void {
    access.value.delete(event.id);
}

function updateQualificationAccessEntryType(event: { id: string; type: number }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.type = event.type;
    access.value.set(event.id, accessEntry);
}

function updateQualificationAccessEntryName(event: { id: string; job?: Job; req?: Qualification }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    if (event.job) {
        accessEntry.values.job = event.job.name;
    }

    access.value.set(event.id, accessEntry);
}

function updateQualificationAccessEntryRank(event: { id: string; rank: JobGrade }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.minimumGrade = event.rank.grade;
    access.value.set(event.id, accessEntry);
}

function updateQualificationAccessEntryAccess(event: { id: string; access: AccessLevel }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.accessRole = event.access;
    access.value.set(event.id, accessEntry);
}

function updateQualificationRequirement(idx: number, qualification?: QualificationShort): void {
    if (!qualification) {
        return;
    }

    quali.value.requirements[idx].targetQualificationId = qualification.id;
    // TODO
}
</script>

<template>
    <div class="m-2">
        <form @submit.prevent="onSubmitThrottle">
            <div
                class="flex flex-col gap-2 rounded-t-lg bg-base-800 px-3 py-4 text-neutral"
                :class="!canDo.edit ? 'rounded-b-md' : ''"
            >
                <div class="flex flex-row gap-2">
                    <div class="flex-0 shrink max-w-48">
                        <label for="abbreviation" class="block text-base font-medium">
                            {{ $t('common.abbreviation') }}
                        </label>
                        <VeeField
                            name="abbreviation"
                            type="text"
                            :placeholder="$t('common.abbreviation')"
                            :label="$t('common.abbreviation')"
                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-3xl sm:leading-6"
                            :disabled="!canEdit || !canDo.edit"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="abbreviation" as="p" class="mt-2 text-sm text-error-400" />
                    </div>

                    <div class="flex-1">
                        <label for="title" class="block text-base font-medium">
                            {{ $t('common.title') }}
                        </label>
                        <VeeField
                            name="title"
                            type="text"
                            :placeholder="$t('common.title')"
                            :label="$t('common.title')"
                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-3xl sm:leading-6"
                            :disabled="!canEdit || !canDo.edit"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="title" as="p" class="mt-2 text-sm text-error-400" />
                    </div>
                </div>
                <div class="flex flex-row gap-2">
                    <div class="flex-1">
                        <label for="description" class="block text-sm font-medium">
                            {{ $t('common.description') }}
                        </label>
                        <VeeField
                            name="description"
                            as="textarea"
                            :placeholder="$t('common.description')"
                            :label="$t('common.description')"
                            class="block w-full h-20 rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:leading-6"
                            :disabled="!canEdit || !canDo.edit"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
                    </div>
                    <div class="flex-0 min-w-32">
                        <label for="closed" class="block text-sm font-medium"> {{ $t('common.close', 2) }}? </label>
                        <Listbox v-model="quali.closed" as="div" :disabled="!canEdit || !canDo.edit">
                            <div class="relative">
                                <ListboxButton
                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                >
                                    <span class="block truncate">
                                        {{ openclose.find((e) => e.closed === quali.closed.closed)?.label }}</span
                                    >
                                    <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                                        <ChevronDownIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
                                    </span>
                                </ListboxButton>

                                <transition
                                    leave-active-class="transition duration-100 ease-in"
                                    leave-from-class="opacity-100"
                                    leave-to-class="opacity-0"
                                >
                                    <ListboxOptions
                                        class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                    >
                                        <ListboxOption
                                            v-for="st in openclose"
                                            :key="st.closed.toString()"
                                            v-slot="{ active, selected }"
                                            as="template"
                                            :value="st"
                                        >
                                            <li
                                                :class="[
                                                    active ? 'bg-primary-500' : '',
                                                    'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                ]"
                                            >
                                                <span :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']">{{
                                                    st.label
                                                }}</span>

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
                                        </ListboxOption>
                                    </ListboxOptions>
                                </transition>
                            </div>
                        </Listbox>
                    </div>
                </div>
            </div>

            <div v-if="canDo.edit" class="bg-neutral">
                <VeeField
                    v-slot="{ field }"
                    name="content"
                    :placeholder="$t('common.content')"
                    :label="$t('common.content')"
                    :disabled="!canEdit || !canDo.edit"
                >
                    <DocEditor v-bind="field" :model-value="field.value ?? ''" />
                </VeeField>
                <VeeErrorMessage name="content" as="p" class="mt-2 text-sm text-error-400" />
            </div>

            <div class="my-3">
                <h2 class="text-neutral">
                    {{ $t('common.access') }}
                </h2>
                <QualificationAccessEntry
                    v-for="entry in access.values()"
                    :key="entry.id"
                    :init="entry"
                    :access-types="accessTypes"
                    :read-only="!canDo.access"
                    @type-change="updateQualificationAccessEntryType($event)"
                    @name-change="updateQualificationAccessEntryName($event)"
                    @rank-change="updateQualificationAccessEntryRank($event)"
                    @access-change="updateQualificationAccessEntryAccess($event)"
                    @delete-request="removeQualificationAccessEntry($event)"
                />
                <button
                    type="button"
                    :disabled="!canEdit || !canDo.access"
                    class="rounded-full bg-primary-500 p-2 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    data-te-toggle="tooltip"
                    :title="$t('components.documents.document_editor.add_permission')"
                    @click="addQualificationAccessEntry()"
                >
                    <PlusIcon class="h-5 w-5" aria-hidden="true" />
                </button>
            </div>

            <div class="my-3">
                <h2 class="text-neutral">
                    {{ $t('common.requirements') }}
                </h2>

                <QualificationRequirementEntry
                    v-for="(requirement, idx) in quali.requirements"
                    :key="requirement.id"
                    :requirement="requirement"
                    @update-qualification="updateQualificationRequirement(idx, $event)"
                    @remove="quali.requirements.splice(idx, 1)"
                />

                <button
                    type="button"
                    class="mt-2 rounded-full p-1.5 text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                    :disabled="!canSubmit"
                    :class="
                        !canSubmit
                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                            : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500'
                    "
                    @click="quali.requirements.push({ id: '0', qualificationId: '0', targetQualificationId: '0' })"
                >
                    <PlusIcon class="h-5 w-5" aria-hidden="true" />
                </button>
            </div>

            <div class="flex pb-14">
                <button
                    type="submit"
                    :disabled="!meta.valid || !canEdit || !canSubmit"
                    class="flex w-full justify-center rounded-md px-3.5 py-2.5 text-sm font-semibold text-neutral"
                    :class="[
                        !canEdit || !meta.valid || !canSubmit
                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                            : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                    ]"
                >
                    <template v-if="!canSubmit">
                        <LoadingIcon class="mr-2 h-5 w-5 animate-spin" aria-hidden="true" />
                    </template>
                    <template v-if="!id">
                        {{ $t('common.create') }}
                    </template>
                    <template v-else>
                        {{ $t('common.save') }}
                    </template>
                </button>
            </div>
        </form>
    </div>
</template>
