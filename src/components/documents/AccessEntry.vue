<script lang="ts" setup>
import {
    Combobox,
    ComboboxButton,
    ComboboxInput,
    ComboboxOption,
    ComboboxOptions,
    Listbox,
    ListboxButton,
    ListboxOption,
    ListboxOptions,
} from '@headlessui/vue';

import { listEnumValues } from '@protobuf-ts/runtime';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { watchDebounced } from '@vueuse/core';
import { CheckIcon, ChevronDownIcon, CloseIcon } from 'mdi-vue3';
import { ArrayElement } from '~/utils/types';
import { ACCESS_LEVEL } from '~~/gen/ts/resources/documents/access';
import { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import { UserShort } from '~~/gen/ts/resources/users/users';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const props = defineProps<{
    init: {
        id: bigint;
        type: number;
        values: {
            job?: string;
            char?: number;
            accessrole?: ACCESS_LEVEL;
            minimumrank?: number;
        };
    };
    accessTypes: { id: number; name: string }[];
    accessRoles?: ACCESS_LEVEL[];
}>();

const emit = defineEmits<{
    (e: 'typeChange', payload: { id: bigint; type: number }): void;
    (
        e: 'nameChange',
        payload: {
            id: bigint;
            job: Job | undefined;
            char: UserShort | undefined;
        },
    ): void;
    (e: 'rankChange', payload: { id: bigint; rank: JobGrade }): void;
    (e: 'accessChange', payload: { id: bigint; access: ACCESS_LEVEL }): void;
    (e: 'deleteRequest', payload: { id: bigint }): void;
}>();

const selectedAccessType = ref<{ id: number; name: string }>({
    id: -1,
    name: '',
});

let entriesChars = [] as UserShort[];
const queryChar = ref('');
const selectedChar = ref<undefined | UserShort>(undefined);

let entriesJobs = [] as Job[];
const queryJob = ref('');
const selectedJob = ref<Job>();

let entriesMinimumRank = [] as JobGrade[];
const queryMinimumRank = ref('');
const selectedMinimumRank = ref<JobGrade | undefined>(undefined);

let entriesAccessRoles: {
    id: ACCESS_LEVEL;
    label: string;
    value: string;
}[] = [];
if (!props.accessRoles || props.accessRoles.length === 0) {
    const enumVals = listEnumValues(ACCESS_LEVEL);
    entriesAccessRoles = enumVals.map((e, k) => {
        return {
            id: k,
            label: t(`enums.docstore.ACCESS_LEVEL.${e.name}`),
            value: e.name,
        };
    });
} else {
    props.accessRoles.forEach((e) => {
        entriesAccessRoles.push({
            id: e,
            label: t(`enums.docstore.ACCESS_LEVEL.${ACCESS_LEVEL[e]}`),
            value: ACCESS_LEVEL[e],
        });
    });
}
const queryAccessRole = ref('');
const selectedAccessRole = ref<ArrayElement<typeof entriesAccessRoles>>();

async function findJobs(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCompletorClient().completeJobs({
                search: queryJob.value,
            });
            const { response } = await call;

            entriesJobs = response.jobs;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function findChars(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCompletorClient().completeCitizens({
                search: queryChar.value,
            });
            const { response } = await call;

            entriesChars = response.users;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

onMounted(async () => {
    const passedType = props.accessTypes.find((e) => e.id === props.init.type);
    if (passedType) selectedAccessType.value = passedType;

    if (props.init.type === 0 && props.init.values.char !== undefined && props.init.values.accessrole !== undefined) {
        await findChars();
        selectedChar.value = entriesChars.find((char) => char.userId === props.init.values.char);
        selectedAccessRole.value = entriesAccessRoles.find((type) => type.id === props.init.values.accessrole);
    } else if (
        props.init.type === 1 &&
        props.init.values.job !== undefined &&
        props.init.values.minimumrank !== undefined &&
        props.init.values.accessrole !== undefined
    ) {
        await findJobs();
        selectedJob.value = entriesJobs.find((job) => job.name === props.init.values.job);
        if (selectedJob.value) entriesMinimumRank = selectedJob.value.grades;
        selectedMinimumRank.value = entriesMinimumRank.find((rank) => rank.grade === props.init.values.minimumrank);
        selectedAccessRole.value = entriesAccessRoles.find((type) => type.id === props.init.values.accessrole);
    }
});

watchDebounced(queryJob, async () => await findJobs(), {
    debounce: 600,
    maxWait: 1750,
});
watchDebounced(queryChar, async () => await findChars(), {
    debounce: 600,
    maxWait: 1750,
});

watch(selectedAccessType, () => {
    emit('typeChange', {
        id: props.init.id,
        type: selectedAccessType.value.id,
    });

    selectedChar.value = undefined;
    selectedJob.value = undefined;
    selectedMinimumRank.value = undefined;

    if (selectedAccessType.value.id === 0) {
        queryChar.value = '';
        findChars();
    } else {
        queryJob.value = '';
        queryMinimumRank.value = '';
        findJobs();
    }
});

watch(selectedJob, () => {
    if (!selectedJob.value) return;
    emit('nameChange', {
        id: props.init.id,
        job: selectedJob.value,
        char: undefined,
    });

    entriesMinimumRank = selectedJob.value.grades;
});

watch(selectedChar, () => {
    if (!selectedChar.value) return;
    emit('nameChange', {
        id: props.init.id,
        job: undefined,
        char: selectedChar.value,
    });
});

watch(selectedMinimumRank, () => {
    if (!selectedMinimumRank.value) return;
    emit('rankChange', { id: props.init.id, rank: selectedMinimumRank.value });
});

watch(selectedAccessRole, () => {
    if (!selectedAccessRole.value) return;
    emit('accessChange', {
        id: props.init.id,
        access: selectedAccessRole.value.id,
    });
});
</script>

<template>
    <div class="flex flex-row items-center my-2">
        <div class="flex-initial mr-2 w-60">
            <input
                v-if="accessTypes.length === 1"
                type="text"
                disabled
                :value="accessTypes[0].name"
                class="block pl-3 text-left w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
            />
            <Listbox v-else as="div" v-model="selectedAccessType">
                <div class="relative">
                    <ListboxButton
                        class="block pl-3 text-left w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    >
                        <span class="block truncate">{{ selectedAccessType?.name }}</span>
                        <span class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                            <ChevronDownIcon class="w-5 h-5 text-gray-400" aria-hidden="true" />
                        </span>
                    </ListboxButton>

                    <transition
                        leave-active-class="transition duration-100 ease-in"
                        leave-from-class="opacity-100"
                        leave-to-class="opacity-0"
                    >
                        <ListboxOptions
                            class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                        >
                            <ListboxOption
                                as="template"
                                v-for="accessType in accessTypes"
                                :key="accessType.id?.toString()"
                                :value="accessType"
                                v-slot="{ active, selected }"
                            >
                                <li
                                    :class="[
                                        active ? 'bg-primary-500' : '',
                                        'text-neutral relative cursor-default select-none py-2 pl-8 pr-4',
                                    ]"
                                >
                                    <span :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']">{{
                                        accessType.name
                                    }}</span>

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
                            </ListboxOption>
                        </ListboxOptions>
                    </transition>
                </div>
            </Listbox>
        </div>
        <div v-if="selectedAccessType?.id === 0" class="flex flex-grow">
            <div class="flex-1 mr-2">
                <Combobox as="div" v-model="selectedChar">
                    <div class="relative">
                        <ComboboxButton as="div">
                            <ComboboxInput
                                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                @change="queryChar = $event.target.value"
                                :display-value="(char: any) => `${char?.firstname} ${char?.lastname}`"
                            />
                        </ComboboxButton>

                        <ComboboxOptions
                            v-if="entriesChars.length > 0"
                            class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                        >
                            <ComboboxOption
                                v-for="char in entriesChars"
                                :key="char.identifier"
                                :value="char"
                                as="char"
                                v-slot="{ active, selected }"
                            >
                                <li
                                    :class="[
                                        'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                        active ? 'bg-primary-500' : '',
                                    ]"
                                >
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ char.firstname }} {{ char.lastname }}
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
        <div v-else class="flex flex-grow">
            <div class="flex-1 mr-2">
                <Combobox as="div" v-model="selectedJob">
                    <div class="relative">
                        <ComboboxButton as="div">
                            <ComboboxInput
                                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                @change="queryJob = $event.target.value"
                                :display-value="(job: any) => job?.label"
                            />
                        </ComboboxButton>

                        <ComboboxOptions
                            v-if="entriesJobs.length > 0"
                            class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                        >
                            <ComboboxOption
                                v-for="job in entriesJobs"
                                :key="job.name"
                                :value="job"
                                as="job"
                                v-slot="{ active, selected }"
                            >
                                <li
                                    :class="[
                                        'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
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
                                        <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                    </span>
                                </li>
                            </ComboboxOption>
                        </ComboboxOptions>
                    </div>
                </Combobox>
            </div>
            <div class="flex-1 mr-2">
                <Combobox as="div" v-model="selectedMinimumRank">
                    <div class="relative">
                        <ComboboxButton as="div">
                            <ComboboxInput
                                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                @change="queryMinimumRank = $event.target.value"
                                :display-value="(rank: any) => rank?.label"
                            />
                        </ComboboxButton>

                        <ComboboxOptions
                            v-if="entriesMinimumRank.length > 0"
                            class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                        >
                            <ComboboxOption
                                v-for="rank in entriesMinimumRank"
                                :key="rank.grade"
                                :value="rank"
                                as="minimumrank"
                                v-slot="{ active, selected }"
                            >
                                <li
                                    :class="[
                                        'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
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
                                        <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                    </span>
                                </li>
                            </ComboboxOption>
                        </ComboboxOptions>
                    </div>
                </Combobox>
            </div>
        </div>
        <div class="mr-2 flex-inital w-60">
            <Combobox as="div" v-model="selectedAccessRole">
                <div class="relative">
                    <ComboboxButton as="div">
                        <ComboboxInput
                            class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            @change="queryAccessRole = $event.target.value"
                            :display-value="(role: any) => role.label"
                        />
                    </ComboboxButton>

                    <ComboboxOptions
                        v-if="entriesAccessRoles.length > 0"
                        class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                    >
                        <ComboboxOption
                            v-for="role in entriesAccessRoles"
                            :key="role.id?.toString()"
                            :value="role"
                            as="accessrole"
                            v-slot="{ active, selected }"
                        >
                            <li
                                :class="[
                                    'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                    active ? 'bg-primary-500' : '',
                                ]"
                            >
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
                                    <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                </span>
                            </li>
                        </ComboboxOption>
                    </ComboboxOptions>
                </div>
            </Combobox>
        </div>
        <div class="flex-initial">
            <button
                type="button"
                class="rounded-full bg-primary-500 p-1.5 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600"
            >
                <CloseIcon class="w-6 h-6" @click="$emit('deleteRequest', { id: props.init.id })" aria-hidden="true" />
            </button>
        </div>
    </div>
</template>
