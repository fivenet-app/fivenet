<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
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
import { useCompletorStore } from '~/store/completor';

const props = defineProps<{
    qualificationId?: string;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const completorStore = useCompletorStore();

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
    discordSettingsSyncEnabled: boolean;
    discordSettingsRoleName?: string;
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
    if (props.qualificationId) {
        try {
            const call = $grpc.getQualificationsClient().getQualification({
                qualificationId: props.qualificationId,
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
                setFieldValue('discordSettingsSyncEnabled', qualification.discordSettings?.syncEnabled ?? false);
                setFieldValue('discordSettingsRoleName', qualification.discordSettings?.roleName);
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
            discordSettings: {
                syncEnabled: values.discordSettingsSyncEnabled,
                roleName: values.discordSettingsRoleName,
            },
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
            id: props.qualificationId!,
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
            discordSettings: {
                syncEnabled: values.discordSettingsSyncEnabled,
                roleName: values.discordSettingsRoleName,
            },
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
                qualificationId: props.qualificationId!,
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
        discordSettingsSyncEnabled: {},
        discordSettingsRoleName: { required: false, min: 3, max: 50 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(async (values): Promise<any> => {
    if (props.qualificationId === undefined) {
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

    const id = access.value.size > 0 ? parseInt([...access.value.keys()]?.pop() ?? '1') + 1 : 0;
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

    quali.value.requirements[idx].qualificationId = props.qualificationId ?? '0';
    quali.value.requirements[idx].targetQualificationId = qualification.id;
}

const { data: jobs } = useAsyncData('completor-jobs', () => completorStore.listJobs());
</script>

<template>
    <div class="m-2">
        <UForm :state="{}">
            <div class="flex flex-col gap-2 rounded-t-lg bg-base-800 px-3 py-4" :class="!canDo.edit ? 'rounded-b-md' : ''">
                <div class="flex flex-row gap-2">
                    <div class="max-w-48 shrink">
                        <label for="abbreviation" class="block text-base font-medium">
                            {{ $t('common.abbreviation') }}
                        </label>
                        <VeeField
                            name="abbreviation"
                            type="text"
                            :placeholder="$t('common.abbreviation')"
                            :label="$t('common.abbreviation')"
                            class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-3xl sm:leading-6"
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
                            class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-3xl sm:leading-6"
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
                            class="placeholder:text-accent-200 block h-20 w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:leading-6"
                            :disabled="!canEdit || !canDo.edit"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
                    </div>
                    <div class="min-w-32">
                        <label for="closed" class="block text-sm font-medium"> {{ $t('common.close', 2) }}? </label>
                        <USelectMenu
                            v-model="quali.closed"
                            :options="openclose"
                            :placeholder="quali.closed ? quali.closed.label : $t('common.na')"
                            :disabled="!canEdit || !canDo.edit"
                        >
                            <template #option-empty="{ query: search }">
                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                            </template>
                            <template #empty>
                                {{ $t('common.not_found', [$t('common.close', 1)]) }}
                            </template>
                        </USelectMenu>
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
                <h2>
                    {{ $t('common.access') }}
                </h2>
                <QualificationAccessEntry
                    v-for="entry in access.values()"
                    :key="entry.id"
                    :init="entry"
                    :access-types="accessTypes"
                    :read-only="!canDo.access"
                    :jobs="jobs"
                    @type-change="updateQualificationAccessEntryType($event)"
                    @name-change="updateQualificationAccessEntryName($event)"
                    @rank-change="updateQualificationAccessEntryRank($event)"
                    @access-change="updateQualificationAccessEntryAccess($event)"
                    @delete-request="removeQualificationAccessEntry($event)"
                />
                <UButton
                    :ui="{ rounded: 'rounded-full' }"
                    :disabled="!canEdit || !canDo.access"
                    icon="i-mdi-plus"
                    :title="$t('components.documents.document_editor.add_permission')"
                    @click="addQualificationAccessEntry()"
                />
            </div>

            <div class="my-3">
                <h2>
                    {{ $t('common.requirements') }}
                </h2>

                <QualificationRequirementEntry
                    v-for="(requirement, idx) in quali.requirements"
                    :key="requirement.id"
                    :requirement="requirement"
                    @update-qualification="updateQualificationRequirement(idx, $event)"
                    @remove="quali.requirements.splice(idx, 1)"
                />

                <UButton
                    :ui="{ rounded: 'rounded-full' }"
                    :disabled="!canSubmit"
                    icon="i-mdi-plus"
                    @click="quali.requirements.push({ id: '0', qualificationId: '0', targetQualificationId: '0' })"
                />
            </div>

            <div class="my-3">
                <h2>
                    {{ $t('common.discord') }}
                </h2>

                <UAccordion :items="[{ slot: 'discord', label: $t('common.discord') }]">
                    <template #discord>
                        <div>
                            <VeeField v-slot="{ handleInput }" name="discordSettingsSyncEnabled">
                                <div class="flex items-center">
                                    <UToggle :disabled="!canEdit || !canDo.edit" @update:model-value="handleInput($event)">
                                        <span class="sr-only">
                                            {{ $t('common.enabled') }}
                                        </span>
                                    </UToggle>
                                    <span class="ml-3 text-sm font-medium">{{ $t('common.enabled') }}</span>
                                </div>
                            </VeeField>

                            <label for="discordSettingsRoleName" class="block text-base font-medium">
                                {{ $t('common.role') }}
                            </label>
                            <VeeField
                                name="discordSettingsRoleName"
                                type="text"
                                :placeholder="$t('common.role')"
                                :label="$t('common.role')"
                                :disabled="!canEdit || !canDo.edit"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                            <VeeErrorMessage name="discordSettingsRoleName" as="p" class="mt-2 text-sm text-error-400" />
                        </div>
                    </template>
                </UAccordion>
            </div>

            <div class="flex pb-14">
                <UButton
                    block
                    :disabled="!meta.valid || !canEdit || !canSubmit"
                    :loading="!canSubmit"
                    @click="onSubmitThrottle"
                >
                    <template v-if="!qualificationId">
                        {{ $t('common.create') }}
                    </template>
                    <template v-else>
                        {{ $t('common.save') }}
                    </template>
                </UButton>
            </div>
        </UForm>
    </div>
</template>
