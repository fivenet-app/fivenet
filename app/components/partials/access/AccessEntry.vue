<script lang="ts" setup>
import { useAuth } from '~/composables/useAuth';
import { useCompletorStore } from '~/stores/completor';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import { QualificationExamMode } from '~~/gen/ts/resources/qualifications/exam/exam';
import type { QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import SelectMenu from '../SelectMenu.vue';
import type { AccessLevelEnum, AccessType, MixedAccessEntry } from './helpers';

const props = withDefaults(
    defineProps<{
        disabled?: boolean;
        requiredMode?: 'badge' | 'checkbox' | 'none';
        lockRequiredCheckbox?: boolean;
        accessTypes: AccessType[];
        accessRoles?: AccessLevelEnum[];
        jobs?: Job[] | undefined;
        hideGrade?: boolean;
        hideJobs?: string[];
        hideOtherJobs?: boolean;
        existingEntries?: MixedAccessEntry[];
        name?: string;
    }>(),
    {
        disabled: false,
        requiredMode: 'badge',
        lockRequiredCheckbox: false,
        accessRoles: () => [],
        jobs: () => [],
        hideGrade: false,
        hideJobs: () => [],
        hideOtherJobs: false,
        existingEntries: () => [],
        name: undefined,
    },
);

defineEmits<{
    (e: 'delete'): void;
}>();

const entry = defineModel<MixedAccessEntry>({ required: true });

const completorStore = useCompletorStore();
const { activeChar } = useAuth();

const { game } = useAppConfig();
const currentJob = computed(() => activeChar.value?.job);

const requiredAccessFloor = computed(() => {
    if (!entry.value.required) return undefined;

    return entry.value.requiredAccess ?? entry.value.access;
});

const accessRoleItems = computed(() => {
    if (requiredAccessFloor.value === undefined) {
        return props.accessRoles;
    }

    return props.accessRoles.filter((role) => role.value >= requiredAccessFloor.value!);
});

const requiredSubjectLocked = computed(() => props.disabled || !!entry.value.required);
const otherEntries = computed(() =>
    props.existingEntries.filter((existing) => existing.id !== entry.value.id || existing.type !== entry.value.type),
);

function isDuplicateUser(userId: number | undefined): boolean {
    return (
        userId !== undefined && otherEntries.value.some((existing) => existing.type === 'user' && existing.userId === userId)
    );
}

function isDuplicateQualification(qualificationId: number | undefined): boolean {
    return (
        qualificationId !== undefined &&
        otherEntries.value.some((existing) => existing.type === 'qualification' && existing.qualificationId === qualificationId)
    );
}

function isDuplicateJob(job: string | undefined): boolean {
    return (
        job !== undefined &&
        otherEntries.value.some(
            (existing) => existing.type === 'job' && existing.job === job && existing.minimumGrade === entry.value.minimumGrade,
        )
    );
}

const jobItems = computed(() => {
    const filteredJobs = props.jobs?.filter((j) => props.hideJobs.length === 0 || !props.hideJobs.includes(j.name)) ?? [];

    const visibleJobs =
        !props.hideOtherJobs || entry.value.type !== 'job'
            ? filteredJobs
            : currentJob.value
              ? filteredJobs.filter((job) => job.name === currentJob.value)
              : [];

    return visibleJobs.map((job) => ({ ...job, disabled: isDuplicateJob(job.name) }));
});

const gradeItems = computed(() => {
    const grades = props.jobs.find((job) => job.name === entry.value.job)?.grades ?? [];

    return grades.map((grade) => ({
        ...grade,
        disabled: otherEntries.value.some(
            (existing) => existing.type === 'job' && existing.job === entry.value.job && existing.minimumGrade === grade.grade,
        ),
    }));
});

function markDuplicateUsers(users: UserShort[]): (UserShort & { disabled?: boolean })[] {
    return users.map((user) => ({ ...user, disabled: isDuplicateUser(user.userId) }));
}

function markDuplicateQualifications(qualifications: QualificationShort[]): (QualificationShort & { disabled?: boolean })[] {
    return qualifications.map((qualification) => ({
        ...qualification,
        disabled: isDuplicateQualification(qualification.id),
    }));
}

watch(
    () => [entry.value.required, entry.value.requiredAccess, entry.value.access] as const,
    ([required, requiredAccess, access]) => {
        if (!required) return;

        const nextRequiredAccess = requiredAccess ?? access;
        const nextAccess = Math.max(access, nextRequiredAccess);

        if (nextRequiredAccess === requiredAccess && nextAccess === access) return;

        entry.value = {
            ...entry.value,
            requiredAccess: nextRequiredAccess,
            access: nextAccess,
        };
    },
    { immediate: true },
);

const selectedUser = ref<UserShort | undefined>();
let suppressUserUpdate = false;
function setSelectedUserSilently(user: UserShort | undefined): void {
    if (selectedUser.value === user) return;

    suppressUserUpdate = true;
    selectedUser.value = user;
}

watch(selectedUser, () => {
    if (suppressUserUpdate) {
        suppressUserUpdate = false;
        return;
    }

    if (entry.value.type !== 'user') return;

    entry.value.user = selectedUser.value;
    entry.value.userId = selectedUser.value?.userId;
});

const selectedQualification = ref<QualificationShort | undefined>();
let suppressQualificationUpdate = false;
function setSelectedQualificationSilently(qualification: QualificationShort | undefined): void {
    if (selectedQualification.value === qualification) return;

    suppressQualificationUpdate = true;
    selectedQualification.value = qualification;
}

watch(selectedQualification, () => {
    if (suppressQualificationUpdate) {
        suppressQualificationUpdate = false;
        return;
    }

    if (entry.value.type !== 'qualification') return;

    entry.value.qualification = selectedQualification.value;
    entry.value.qualificationId = selectedQualification.value?.id;
});

async function findUser(userId?: number): Promise<UserShort[]> {
    if (userId === undefined) return [];

    return completorStore.completeCitizens({
        search: '',
        userIds: [userId],
    });
}

function userFallback(userId: number): UserShort {
    return {
        userId,
        job: '',
        jobGrade: 0,
        firstname: '#UserID',
        lastname: String(userId),
        dateofbirth: '',
    };
}

let qualificationsClientPromise: ReturnType<typeof getQualificationsQualificationsClient> | undefined;

function getQualificationsClient(): ReturnType<typeof getQualificationsQualificationsClient> {
    return (qualificationsClientPromise ??= getQualificationsQualificationsClient());
}

let setFromPropsRun = 0;

function resetSelections(): void {
    if (entry.value.type !== 'user') setSelectedUserSilently(undefined);
    if (entry.value.type !== 'qualification') setSelectedQualificationSilently(undefined);
}

async function setFromProps(): Promise<void> {
    const run = ++setFromPropsRun;
    resetSelections();

    if (entry.value.type === 'user' && entry.value.userId !== undefined) {
        if (selectedUser.value?.userId === entry.value.userId) return;

        const userId = entry.value.userId;
        if (entry.value.user?.userId === userId) {
            setSelectedUserSilently(entry.value.user);
            return;
        }

        setSelectedUserSilently(undefined);
        let users: UserShort[] = [];
        try {
            users = await findUser(userId);
        } catch (_) {
            // Keep the ID and show a local fallback when the user lookup fails.
        }
        if (run !== setFromPropsRun || entry.value.type !== 'user' || entry.value.userId !== userId) return;

        setSelectedUserSilently(users.find((char) => char.userId === userId) ?? userFallback(userId));
    } else if (entry.value.type === 'qualification' && entry.value.qualificationId !== undefined) {
        if (selectedQualification.value?.id === entry.value.qualificationId) return;

        const qualificationId = entry.value.qualificationId;
        if (entry.value.qualification?.id === qualificationId) {
            setSelectedQualificationSilently(entry.value.qualification);
            return;
        }

        setSelectedQualificationSilently(undefined);
        try {
            const client = await getQualificationsClient();
            const { response } = await client.getQualification({
                qualificationId,
            });
            if (
                run !== setFromPropsRun ||
                entry.value.type !== 'qualification' ||
                entry.value.qualificationId !== qualificationId
            )
                return;

            setSelectedQualificationSilently(response.qualification);
        } catch (_) {
            if (
                run !== setFromPropsRun ||
                entry.value.type !== 'qualification' ||
                entry.value.qualificationId !== qualificationId
            )
                return;

            // Fallback to show qualification id
            setSelectedQualificationSilently({
                id: qualificationId,
                job: '',
                weight: 0,
                abbreviation: 'N/A',
                title: 'N/A (ID: ' + qualificationId + ')',
                closed: false,
                draft: false,
                public: false,
                creatorJob: '',
                examMode: QualificationExamMode.UNSPECIFIED,
                requirements: [],
            });
        }
    } else if (entry.value.type === 'job' && props.hideOtherJobs && currentJob.value) {
        entry.value.job = currentJob.value;
    }
}

watch(
    () => entry.value.type,
    (type, previousType) => {
        if (previousType === undefined || type === previousType) return;

        entry.value = {
            ...entry.value,
            userId: undefined,
            user: undefined,
            job: undefined,
            minimumGrade: undefined,
            qualificationId: undefined,
            qualification: undefined,
        };
    },
);

setFromProps();
watch([() => props.jobs, () => props.hideOtherJobs], () => setFromProps(), { deep: true });
watch(currentJob, () => setFromProps());
watch(
    () => [entry.value.type, entry.value.userId, entry.value.qualificationId] as const,
    () => setFromProps(),
);

watch(
    [() => entry.value.job, () => props.jobs],
    ([job, jobs]) => {
        if (!job) return;

        const grades = jobs.find((j) => j.name === job)?.grades;
        if (!grades?.length) return;

        const hasValidGrade = grades.some((grade) => grade.grade === entry.value.minimumGrade);
        if (hasValidGrade && entry.value.minimumGrade !== -1) return;

        entry.value.minimumGrade = grades[grades.length - 1]?.grade ?? game.startJobGrade;
    },
    { immediate: true, deep: true },
);
</script>

<template>
    <div class="flex flex-1 flex-col gap-2 pb-2 md:flex-row md:pb-0">
        <div class="grid grid-cols-2 gap-2 md:flex md:flex-1">
            <div class="flex flex-initial flex-row items-center gap-2">
                <UFormField v-if="requiredMode === 'badge' && entry.required" :name="`${$props.name}.required`">
                    <UTooltip :text="$t('components.access.required_notice')">
                        <UBadge color="warning" variant="soft" :label="$t('common.required')" />
                    </UTooltip>
                </UFormField>

                <UFormField v-if="requiredMode === 'checkbox'" :name="`${$props.name}.required`">
                    <UTooltip class="flex-initial" :text="$t('common.require')">
                        <UCheckbox
                            v-model="entry.required"
                            :disabled="disabled || (lockRequiredCheckbox && !!entry.required)"
                            name="required"
                        />
                    </UTooltip>
                </UFormField>

                <UFormField
                    class="h-full min-w-40 flex-initial"
                    :name="`${$props.name}.type`"
                    :label="$t('common.type')"
                    :ui="{ label: 'md:hidden' }"
                >
                    <UInput v-if="accessTypes.length === 1" type="text" disabled :model-value="accessTypes[0]?.label" />
                    <ClientOnly v-else>
                        <USelectMenu
                            v-model="entry.type"
                            class="w-full"
                            :disabled="requiredSubjectLocked"
                            :placeholder="$t('common.type')"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            value-key="value"
                            :items="accessTypes"
                        >
                            <template #empty>
                                {{ $t('common.not_found', [$t('common.type')]) }}
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>
            </div>

            <UFormField
                v-if="entry.type === 'user'"
                class="flex-1"
                :name="`${$props.name}.userId`"
                :label="$t('common.user')"
                :ui="{ label: 'md:hidden' }"
            >
                <SelectMenu
                    v-model="selectedUser"
                    class="w-full"
                    :disabled="requiredSubjectLocked"
                    :searchable="
                        async (q: string) =>
                            markDuplicateUsers(
                                await completorStore.completeCitizens({
                                    search: q,
                                    userIds: entry.userId ? [entry.userId] : [],
                                }),
                            )
                    "
                    searchable-key="completor-citizens"
                    :filter-fields="['firstname', 'lastname']"
                    :search-input="{ placeholder: $t('common.search_field') }"
                    :placeholder="$t('common.citizen', 1)"
                >
                    <template v-if="selectedUser" #default>
                        {{ userToLabel(selectedUser) }}
                    </template>

                    <template #item-label="{ item }">
                        {{ `${item?.firstname} ${item?.lastname} (${item?.dateofbirth})` }}
                    </template>

                    <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                </SelectMenu>
            </UFormField>

            <UFormField
                v-else-if="entry.type === 'qualification'"
                class="flex-1"
                :name="`${$props.name}.qualificationId`"
                :label="$t('common.qualification')"
                :ui="{ label: 'md:hidden' }"
            >
                <SelectMenu
                    v-model="selectedQualification"
                    class="w-full"
                    :disabled="requiredSubjectLocked"
                    :searchable="
                        async (q: string) => {
                            const client = await getQualificationsClient();
                            const { response } = await client.listQualifications({
                                pagination: {
                                    offset: 0,
                                },
                                search: q,
                            });
                            return markDuplicateQualifications((response?.qualifications ?? []) as QualificationShort[]);
                        }
                    "
                    searchable-key="complete-qualifications"
                    :filter-fields="['abbreviation', 'title']"
                    :search-input="{ placeholder: $t('common.search_field') }"
                    :placeholder="$t('common.qualification', 1)"
                >
                    <template v-if="selectedQualification" #default>
                        {{ selectedQualification.abbreviation }}: {{ selectedQualification.title }}
                    </template>

                    <template #item-label="{ item }">
                        {{ `${item?.abbreviation}: ${item?.title}` }}
                    </template>

                    <template #empty> {{ $t('common.not_found', [$t('common.qualification', 2)]) }} </template>
                </SelectMenu>
            </UFormField>

            <template v-else>
                <UFormField
                    v-if="!hideOtherJobs || !currentJob"
                    class="flex-1"
                    :name="`${$props.name}.job`"
                    :label="$t('common.job')"
                    :ui="{ label: 'md:hidden' }"
                >
                    <ClientOnly>
                        <USelectMenu
                            v-model="entry.job"
                            class="w-full"
                            :disabled="requiredSubjectLocked || (hideOtherJobs && !currentJob)"
                            :filter-fields="['label', 'name']"
                            value-key="name"
                            :items="jobItems"
                            :placeholder="$t('common.job')"
                            :search-input="{ placeholder: $t('common.search_field') }"
                        >
                            <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>

                <UFormField
                    v-if="!hideGrade"
                    class="flex-1"
                    :name="`${$props.name}.minimumGrade`"
                    :label="$t('common.rank')"
                    :ui="{ label: 'md:hidden' }"
                >
                    <ClientOnly>
                        <USelectMenu
                            class="w-full"
                            :model-value="gradeItems.find((grade) => grade.grade === entry.minimumGrade)"
                            :disabled="requiredSubjectLocked || !entry.job"
                            :filter-fields="['name', 'label']"
                            :items="gradeItems"
                            :placeholder="$t('common.rank')"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            @update:model-value="entry.minimumGrade = $event?.grade ?? undefined"
                        >
                            <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>
            </template>

            <UFormField
                class="min-w-60 flex-initial"
                :name="`${$props.name}.access`"
                :label="$t('common.access')"
                :ui="{ label: 'md:hidden' }"
            >
                <ClientOnly>
                    <USelectMenu
                        v-model="entry.access"
                        class="w-full"
                        :disabled="disabled || accessRoleItems.length <= 1"
                        value-key="value"
                        :items="accessRoleItems"
                        :filter-fields="['label']"
                        :placeholder="$t('common.na')"
                        :search-input="{ placeholder: $t('common.search_field') }"
                    >
                        <template #default>
                            {{ accessRoleItems?.find((a) => a.value === entry.access)?.label ?? $t('common.na') }}
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.access', 2)]) }} </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormField>
        </div>

        <UFormField class="md:mt-1" :ui="{ container: 'flex justify-end-safe md:inline' }">
            <UTooltip v-if="!disabled" :text="entry.required ? $t('common.required') : $t('components.access.remove_entry')">
                <UButton
                    class="flex-initial"
                    :color="entry.required ? 'gray' : 'red'"
                    icon="i-mdi-remove"
                    :label="$t('components.access.remove_entry')"
                    :disabled="entry.required"
                    :ui="{ label: 'md:hidden' }"
                    @click="$emit('delete')"
                />
            </UTooltip>
        </UFormField>
    </div>
</template>
