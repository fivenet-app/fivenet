<script lang="ts" setup>
import { Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { CheckIcon, ChevronDownIcon, LoadingIcon, PlusIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useNotificatorStore } from '~/store/notificator';
import type { AccessLevel, Qualification } from '~~/gen/ts/resources/jobs/qualifications';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import type { CreateQualificationResponse, UpdateQualificationResponse } from '~~/gen/ts/services/jobs/qualifications';
import QualificationAccessEntry from '~/components/jobs/qualifications/QualificationAccessEntry.vue';

const props = defineProps<{
    id?: string;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

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
    summary: string;
    description: string;
}

const openclose = [
    { id: 0, label: t('common.open', 2), closed: false },
    { id: 1, label: t('common.close', 2), closed: true },
];

const qualification = ref<{
    closed: { id: number; label: string; closed: boolean };
    public: boolean;
}>({
    closed: openclose[0],
    public: false,
});

const access = ref<
    Map<
        string,
        {
            id: string;
            type: number;
            values: {
                job?: string;
                quali?: string;
                accessRole?: AccessLevel;
                minimumGrade?: number;
            };
        }
    >
>(new Map());

async function createQualification(values: FormData): Promise<CreateQualificationResponse> {
    try {
        const call = $grpc.getJobsQualificationsClient().createQualification({
            qualification: {
                id: '0',
                job: '',
                weight: 0,
                closed: false,
                abbreviation: values.abbreviation,
                title: values.title,
                summary: values.summary,
                description: values.description,
                creatorId: 0,
                creatorJob: '',
            },
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function updateQualification(values: FormData): Promise<UpdateQualificationResponse> {
    try {
        const call = $grpc.getJobsQualificationsClient().updateQualification({
            qualification: {
                id: '0',
                job: '',
                weight: 0,
                closed: false,
                abbreviation: values.abbreviation,
                title: values.title,
                summary: values.summary,
                description: values.description,
                creatorId: 0,
                creatorJob: '',
            },
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        registrationToken: { required: true, digits: 6 },
        password: { required: true, min: 6, max: 70 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(async (values): Promise<any> => {
    if (props.id === undefined) {
        await createQualification(values).finally(() => setTimeout(() => (canSubmit.value = true), 400));
    } else {
        await updateQualification(values).finally(() => setTimeout(() => (canSubmit.value = true), 400));
    }
});
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const accessTypes = [
    { id: 0, name: t('common.qualifications', 2) },
    { id: 1, name: t('common.job', 2) },
];

function addDocumentAccessEntry(): void {
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

function updateDocumentAccessEntryName(event: { id: string; job?: Job; req?: Qualification }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    if (event.job) {
        accessEntry.values.job = event.job.name;
        accessEntry.values.quali = undefined;
    } else if (event.req) {
        accessEntry.values.job = undefined;
        accessEntry.values.quali = event.req.id;
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
</script>

<template>
    <div class="m-2">
        <form @submit.prevent="onSubmitThrottle">
            <div
                class="flex flex-col gap-2 rounded-t-lg bg-base-800 px-3 py-4 text-neutral"
                :class="!canDo.edit ? 'rounded-b-md' : ''"
            >
                <div>
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
                <div class="flex flex-row gap-2">
                    <div class="flex-1">
                        <label for="category" class="block text-sm font-medium">
                            {{ $t('common.category') }}
                        </label>
                        TODO
                    </div>
                    <div class="flex-1">
                        <label for="state" class="block text-sm font-medium">
                            {{ $t('common.state') }}
                        </label>
                        <VeeField
                            name="state"
                            type="text"
                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            :placeholder="`${$t('common.qualifications', 1)} ${$t('common.state')}`"
                            :label="`${$t('common.qualifications', 1)} ${$t('common.state')}`"
                            :disabled="!canEdit || !canDo.edit"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="state" as="p" class="mt-2 text-sm text-error-400" />
                    </div>
                    <div class="flex-1">
                        <label for="closed" class="block text-sm font-medium"> {{ $t('common.close', 2) }}? </label>
                        <Listbox v-model="qualification.closed" as="div" :disabled="!canEdit || !canDo.edit">
                            <div class="relative">
                                <ListboxButton
                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                >
                                    <span class="block truncate">
                                        {{ openclose.find((e) => e.closed === qualification.closed.closed)?.label }}</span
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
                    @type-change="updateDocumentAccessEntryType($event)"
                    @name-change="updateDocumentAccessEntryName($event)"
                    @rank-change="updateDocumentAccessEntryRank($event)"
                    @access-change="updateDocumentAccessEntryAccess($event)"
                    @delete-request="removeDocumentAccessEntry($event)"
                />
                <button
                    type="button"
                    :disabled="!canEdit || !canDo.access"
                    class="rounded-full bg-primary-500 p-2 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    data-te-toggle="tooltip"
                    :title="$t('components.documents.document_editor.add_permission')"
                    @click="addDocumentAccessEntry()"
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
