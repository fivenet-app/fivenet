<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { listEnumValues } from '@protobuf-ts/runtime';
import { CheckIcon } from 'mdi-vue3';
import { type ArrayElement } from '~/utils/types';
import { AccessLevel } from '~~/gen/ts/resources/qualifications/qualifications';
import { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';

type AccessType = { id: number; name: string };

const props = withDefaults(
    defineProps<{
        readOnly?: boolean;
        init: {
            id: string;
            type: number;
            values: {
                job?: string;
                accessRole?: AccessLevel;
                minimumGrade?: number;
            };
        };
        accessTypes: AccessType[];
        accessRoles?: undefined | AccessLevel[];
        jobs: Job[] | null;
    }>(),
    {
        readOnly: false,
        accessRoles: undefined,
    },
);

const emit = defineEmits<{
    (e: 'typeChange', payload: { id: string; type: number }): void;
    (
        e: 'nameChange',
        payload: {
            id: string;
            job: Job | undefined;
        },
    ): void;
    (e: 'rankChange', payload: { id: string; rank: JobGrade; required?: boolean }): void;
    (e: 'accessChange', payload: { id: string; access: AccessLevel; required?: boolean }): void;
    (e: 'deleteRequest', payload: { id: string }): void;
}>();

const { t } = useI18n();

const selectedAccessType = ref<AccessType>({
    id: -1,
    name: '',
});

const queryJobRaw = ref('');
const queryJob = computed(() => queryJobRaw.value.toLowerCase());
const filteredJobs = computed(() =>
    props.jobs?.filter((j) => j.name.toLowerCase().includes(queryJob.value) || j.label.toLowerCase().includes(queryJob.value)),
);
const selectedJob = ref<Job>();

const queryMinimumRankRaw = ref('');
const queryMinimumRank = computed(() => queryMinimumRankRaw.value.toLowerCase());
const entriesMinimumRank = ref<JobGrade[]>([]);
const filteredJobRanks = computed(() =>
    entriesMinimumRank.value.filter(
        (j) =>
            j.grade.toString().toLowerCase().includes(queryMinimumRank.value) ||
            j.label.toLowerCase().includes(queryMinimumRank.value),
    ),
);
const selectedMinimumRank = ref<JobGrade | undefined>(undefined);

const entriesAccessRoles: {
    id: AccessLevel;
    label: string;
    value: string;
}[] = [];
if (props.accessRoles === undefined || props.accessRoles.length === 0) {
    entriesAccessRoles.push(
        ...listEnumValues(AccessLevel)
            .map((e, k) => {
                return {
                    id: k,
                    label: t(`enums.qualifications.AccessLevel.${e.name}`),
                    value: e.name,
                };
            })
            .filter((e) => e.id !== 0),
    );
} else {
    props.accessRoles.forEach((e) => {
        entriesAccessRoles.push({
            id: e,
            label: t(`enums.qualifications.AccessLevel.${AccessLevel[e]}`),
            value: AccessLevel[e],
        });
    });
}

const queryAccessRole = ref('');
const selectedAccessRole = ref<ArrayElement<typeof entriesAccessRoles>>();

async function setFromProps(): Promise<void> {
    if (props.init.type === 0 && props.init.values.job !== undefined && props.init.values.minimumGrade !== undefined) {
        selectedJob.value = props.jobs?.find((j) => j.name === props.init.values.job);
        if (selectedJob.value) {
            entriesMinimumRank.value = selectedJob.value.grades;
            selectedMinimumRank.value = entriesMinimumRank.value.find((rank) => rank.grade === props.init.values.minimumGrade);
        }
    }

    selectedAccessRole.value = entriesAccessRoles.find((type) => type.id === props.init.values.accessRole);

    const passedType = props.accessTypes.find((e) => e.id === props.init.type);
    if (passedType) {
        selectedAccessType.value = passedType;
    }
}

onMounted(async () => setFromProps());

watch(props, () => setFromProps());

watch(selectedAccessType, async () => {
    emit('typeChange', {
        id: props.init.id,
        type: selectedAccessType.value.id,
    });

    selectedJob.value = undefined;
    selectedMinimumRank.value = undefined;

    if (selectedAccessType.value.id === 0) {
        queryJobRaw.value = '';
        queryMinimumRankRaw.value = '';
    }
});

watch(selectedJob, () => {
    if (!selectedJob.value) {
        return;
    }

    emit('nameChange', {
        id: props.init.id,
        job: selectedJob.value,
    });

    entriesMinimumRank.value = selectedJob.value.grades;
});

watch(selectedMinimumRank, () => {
    if (!selectedMinimumRank.value) {
        return;
    }

    emit('rankChange', { id: props.init.id, rank: selectedMinimumRank.value });
});

watch(selectedAccessRole, () => {
    if (!selectedAccessRole.value) {
        return;
    }

    emit('accessChange', {
        id: props.init.id,
        access: selectedAccessRole.value.id,
    });
});
</script>

<template>
    <div class="my-2 flex flex-row items-center">
        <div class="mr-2 w-60 flex-initial">
            <UInput v-if="accessTypes.length === 1" type="text" disabled :value="accessTypes[0].name" />
            <USelectMenu
                v-else
                v-model="selectedAccessType"
                :disabled="readOnly"
                :options="accessTypes"
                :placeholder="selectedAccessType ? $t(selectedAccessType.name) : $t('common.na')"
            >
                <template #label>
                    <span v-if="selectedAccessType" class="truncate">{{ selectedAccessType.name }}</span>
                </template>
                <template #option="{ option }">
                    <span class="truncate">{{ option.name }}</span>
                </template>
                <template #option-empty="{ query: search }">
                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                </template>
                <template #empty>
                    {{ $t('common.not_found', [$t('common.access', 1)]) }}
                </template>
            </USelectMenu>
        </div>
        <div v-if="selectedAccessType?.id === 0" class="flex grow">
            <div class="mr-2 flex-1">
                <Combobox v-model="selectedJob" as="div" :disabled="readOnly">
                    <div class="relative">
                        <ComboboxButton as="div">
                            <ComboboxInput
                                autocomplete="off"
                                :display-value="(job: any) => job?.label"
                                :class="readOnly ? 'disabled' : ''"
                                @change="queryJobRaw = $event.target.value"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </ComboboxButton>

                        <ComboboxOptions
                            class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                        >
                            <ComboboxOption
                                v-for="job in filteredJobs"
                                :key="job.name"
                                v-slot="{ active, selected }"
                                :value="job"
                            >
                                <li
                                    :class="[
                                        'relative cursor-default select-none py-2 pl-8 pr-4',
                                        active ? 'bg-primary-500' : '',
                                    ]"
                                >
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ job.label }}
                                    </span>

                                    <span
                                        v-if="selected"
                                        :class="[
                                            active ? 'text-neutral' : 'text-primary-500',
                                            'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                        ]"
                                    >
                                        <CheckIcon class="size-5" />
                                    </span>
                                </li>
                            </ComboboxOption>
                        </ComboboxOptions>
                    </div>
                </Combobox>
            </div>
            <div class="mr-2 flex-1">
                <Combobox v-model="selectedMinimumRank" as="div" :disabled="readOnly || selectedJob === undefined">
                    <div class="relative">
                        <ComboboxButton as="div">
                            <ComboboxInput
                                autocomplete="off"
                                :class="readOnly ? 'disabled' : ''"
                                :display-value="(rank: any) => rank?.label"
                                @change="queryMinimumRankRaw = $event.target.value"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </ComboboxButton>

                        <ComboboxOptions
                            class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                        >
                            <ComboboxOption
                                v-for="rank in filteredJobRanks"
                                :key="rank.grade"
                                v-slot="{ active, selected }"
                                :value="rank"
                                as="minimumGrade"
                            >
                                <li
                                    :class="[
                                        'relative cursor-default select-none py-2 pl-8 pr-4',
                                        active ? 'bg-primary-500' : '',
                                    ]"
                                >
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ rank.label }}
                                    </span>

                                    <span
                                        v-if="selected"
                                        :class="[
                                            active ? 'text-neutral' : 'text-primary-500',
                                            'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                        ]"
                                    >
                                        <CheckIcon class="size-5" />
                                    </span>
                                </li>
                            </ComboboxOption>
                        </ComboboxOptions>
                    </div>
                </Combobox>
            </div>
        </div>
        <div class="mr-2 w-60 flex-initial">
            <Combobox v-model="selectedAccessRole" as="div" :disabled="readOnly">
                <div class="relative">
                    <ComboboxButton as="div">
                        <ComboboxInput
                            autocomplete="off"
                            :class="readOnly ? 'disabled' : ''"
                            :display-value="(role: any) => role.label"
                            @change="queryAccessRole = $event.target.value"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </ComboboxButton>

                    <ComboboxOptions
                        class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                    >
                        <ComboboxOption
                            v-for="role in entriesAccessRoles"
                            :key="role.id"
                            v-slot="{ active, selected }"
                            :value="role"
                            as="accessRole"
                        >
                            <li :class="['relative cursor-default select-none py-2 pl-8 pr-4', active ? 'bg-primary-500' : '']">
                                <span :class="['block truncate', selected && 'font-semibold']">
                                    {{ role.label }}
                                </span>

                                <span
                                    v-if="selected"
                                    :class="[
                                        active ? 'text-neutral' : 'text-primary-500',
                                        'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                    ]"
                                >
                                    <CheckIcon class="size-5" />
                                </span>
                            </li>
                        </ComboboxOption>
                    </ComboboxOptions>
                </div>
            </Combobox>
        </div>
        <div v-if="!readOnly" class="flex-initial">
            <UButton
                :ui="{ rounded: 'rounded-full' }"
                icon="i-mdi-close"
                @click="$emit('deleteRequest', { id: props.init.id })"
            />
        </div>
    </div>
</template>
