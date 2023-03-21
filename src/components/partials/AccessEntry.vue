<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import {
    Listbox,
    ListboxButton,
    ListboxOption,
    ListboxOptions,
    Combobox,
    ComboboxButton,
    ComboboxInput,
    ComboboxOption,
    ComboboxOptions
} from '@headlessui/vue';
import {
    CheckIcon,
    ChevronDownIcon,
    XMarkIcon,
} from '@heroicons/vue/20/solid';
import { watchDebounced } from '@vueuse/core';
import { CompleteCharNamesRequest, CompleteJobNamesRequest } from '@arpanet/gen/services/completor/completor_pb';
import { getCompletorClient, handleGRPCError } from '../../grpc';
import { RpcError } from 'grpc-web';
import { Job, JobGrade } from '@arpanet/gen/resources/jobs/jobs_pb';
import { UserShort } from '@arpanet/gen/resources/users/users_pb';
import { DOC_ACCESS } from '@arpanet/gen/resources/documents/documents_pb';
import { toTitleCase } from '../../utils/strings';

const props = defineProps<{
    init: { id: number, type: number, values: { job?: string, char?: number, accessrole?: DOC_ACCESS, minimumrank?: number } }
}>();

const emit = defineEmits<{
    (e: 'typeChange', payload: { id: number, type: number }): void,
    (e: 'nameChange', payload: { id: number, job: Job | undefined, char: UserShort | undefined }): void,
    (e: 'rankChange', payload: { id: number, rank: JobGrade }): void,
    (e: 'accessChange', payload: { id: number, access: DOC_ACCESS }): void,
    (e: 'deleteRequest', payload: { id: number }): void,
}>()

const accessTypes = [
    { id: 0, name: 'Citizen' },
    { id: 1, name: 'Jobs' },
];
const selectedAccessType = ref<{ id: number, name: string }>({ id: -1, name: '' });

let entriesChars = [] as UserShort[];
const queryChar = ref('');
const selectedChar = ref<undefined | UserShort>(undefined);

let entriesJobs = [] as Job[];
const queryJob = ref('');
const selectedJob = ref<Job>();

let entriesMinimumRank = [] as JobGrade[];
const queryMinimumRank = ref('');
const selectedMinimumRank = ref<JobGrade | undefined>(undefined);

let entriesAccessRole = Object.keys(DOC_ACCESS);
const queryAccessRole = ref('');
const selectedAccessRole = ref();

if (props.init.type === 0 && props.init.values.char && props.init.values.accessrole) {
    selectedChar.value = entriesChars.find(char => char.getUserId() === props.init.values.char);
    selectedAccessRole.value = accessTypes.find(type => type.id === props.init.values.accessrole)
} else if (props.init.type === 1 && props.init.values.job && props.init.values.minimumrank && props.init.values.accessrole) {
    selectedJob.value = entriesJobs.find(job => job.getName() === props.init.values.job);
    selectedMinimumRank.value = entriesMinimumRank.find(rank => rank.getGrade() === props.init.values.minimumrank);
    selectedAccessRole.value = accessTypes.find(type => type.id === props.init.values.accessrole)
}

function findJobs(): void {
    const req = new CompleteJobNamesRequest();
    req.setSearch(queryJob.value);

    getCompletorClient().completeJobNames(req, null).then((resp) => {
        entriesJobs = resp.getJobsList();
    }).catch((err: RpcError) => {
        handleGRPCError(err);
    })
}

function findChars(): void {
    const req = new CompleteCharNamesRequest();
    req.setSearch(queryJob.value);

    getCompletorClient().completeCharNames(req, null).then((resp) => {
        entriesChars = resp.getUsersList();
    }).catch((err: RpcError) => {
        handleGRPCError(err);
    })
}

onMounted(() => {
    const passedType = accessTypes.find(e => e.id === props.init.type);
    if (passedType) selectedAccessType.value = passedType;
});

watchDebounced(queryJob, () => findJobs(), { debounce: 750, maxWait: 2000 });
watchDebounced(queryChar, () => findChars(), { debounce: 750, maxWait: 2000 });

watch(selectedAccessType, () => {
    emit('typeChange', { id: props.init.id, type: selectedAccessType.value.id });

    selectedChar.value = undefined;
    selectedJob.value = undefined;
    selectedMinimumRank.value = undefined;

    if (selectedAccessType.value.id === 0) {
        findChars();
    } else {
        findJobs();
    }
});

watch(selectedJob, () => {
    if (!selectedJob.value) return;
    emit('nameChange', { id: props.init.id, job: selectedJob.value, char: undefined });

    entriesMinimumRank = selectedJob.value.getGradesList()
});

watch(selectedChar, () => {
    if (!selectedChar.value) return;
    emit('nameChange', { id: props.init.id, job: undefined, char: selectedChar.value });
});

watch(selectedMinimumRank, () => {
    if (!selectedMinimumRank.value) return;
    emit('rankChange', { id: props.init.id, rank: selectedMinimumRank.value });
});

watch(selectedAccessRole, () => {
    if (!selectedAccessRole.value) return;
    emit('accessChange', { id: props.init.id, access: selectedAccessRole.value });
});
</script>

<template>
    <div class="flex flex-row items-center my-2">
        <div class="flex-initial w-60 mr-2">
            <Listbox as="div" v-model="selectedAccessType">
                <div class="relative">
                    <ListboxButton
                        class="relative w-full cursor-default rounded-md bg-white py-1.5 pl-3 pr-10 text-left text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-600 sm:text-sm sm:leading-6">
                        <span class="block truncate">{{ selectedAccessType?.name }}</span>
                        <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                            <ChevronDownIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
                        </span>
                    </ListboxButton>

                    <transition leave-active-class="transition ease-in duration-100" leave-from-class="opacity-100"
                        leave-to-class="opacity-0">
                        <ListboxOptions
                            class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                            <ListboxOption as="template" v-for="accessType in accessTypes" :key="accessType.id"
                                :value="accessType" v-slot="{ active, selected }">
                                <li
                                    :class="[active ? 'bg-indigo-600 text-white' : 'text-gray-900', 'relative cursor-default select-none py-2 pl-8 pr-4']">
                                    <span :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']">{{
                                        accessType.name
                                    }}</span>

                                    <span v-if="selected"
                                        :class="[active ? 'text-white' : 'text-indigo-600', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
                                        <CheckIcon class="h-5 w-5" aria-hidden="true" />
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
                                class="w-full rounded-md border-0 bg-white py-1.5 pl-3 pr-10 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                                @change="queryChar = $event.target.value"
                                :display-value="(char: any) => `${char?.getFirstname()} ${char?.getLastname()}`" />
                        </ComboboxButton>

                        <ComboboxOptions v-if="entriesChars.length > 0"
                            class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                            <ComboboxOption v-for="char in entriesChars" :key="char.getIdentifier()" :value="char" as="char"
                                v-slot="{ active, selected }">
                                <li
                                    :class="['relative cursor-default select-none py-2 pl-8 pr-4', active ? 'bg-indigo-600 text-white' : 'text-gray-900']">
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ char.getFirstname() }} {{ char.getLastname() }}
                                    </span>

                                    <span v-if="selected"
                                        :class="['absolute inset-y-0 left-0 flex items-center pl-1.5', active ? 'text-white' : 'text-indigo-600']">
                                        <CheckIcon class="h-5 w-5" aria-hidden="true" />
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
                                class="w-full rounded-md border-0 bg-white py-1.5 pl-3 pr-10 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                                @change="queryJob = $event.target.value" :display-value="(job: any) => job?.getLabel()" />
                        </ComboboxButton>

                        <ComboboxOptions v-if="entriesJobs.length > 0"
                            class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                            <ComboboxOption v-for="job in entriesJobs" :key="job.getName()" :value="job" as="job"
                                v-slot="{ active, selected }">
                                <li
                                    :class="['relative cursor-default select-none py-2 pl-8 pr-4', active ? 'bg-indigo-600 text-white' : 'text-gray-900']">
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ job.getLabel() }}
                                    </span>

                                    <span v-if="selected"
                                        :class="['absolute inset-y-0 left-0 flex items-center pl-1.5', active ? 'text-white' : 'text-indigo-600']">
                                        <CheckIcon class="h-5 w-5" aria-hidden="true" />
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
                                class="w-full rounded-md border-0 bg-white py-1.5 pl-3 pr-10 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                                @change="queryMinimumRank = $event.target.value"
                                :display-value="(rank: any) => rank?.getLabel()" />
                        </ComboboxButton>

                        <ComboboxOptions v-if="entriesMinimumRank.length > 0"
                            class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                            <ComboboxOption v-for="rank in entriesMinimumRank" :key="rank.getGrade()" :value="rank"
                                as="minimumrank" v-slot="{ active, selected }">
                                <li
                                    :class="['relative cursor-default select-none py-2 pl-8 pr-4', active ? 'bg-indigo-600 text-white' : 'text-gray-900']">
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ rank.getLabel() }}
                                    </span>

                                    <span v-if="selected"
                                        :class="['absolute inset-y-0 left-0 flex items-center pl-1.5', active ? 'text-white' : 'text-indigo-600']">
                                        <CheckIcon class="h-5 w-5" aria-hidden="true" />
                                    </span>
                                </li>
                            </ComboboxOption>
                        </ComboboxOptions>
                    </div>
                </Combobox>
            </div>
        </div>
        <div class="flex-inital w-60 mr-2">
            <Combobox as="div" v-model="selectedAccessRole">
                <div class="relative">
                    <ComboboxButton as="div">
                        <ComboboxInput
                            class="w-full rounded-md border-0 bg-white py-1.5 pl-3 pr-10 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                            @change="queryAccessRole = $event.target.value"
                            :display-value="(role: any) => toTitleCase(role.toLowerCase())" />
                    </ComboboxButton>

                    <ComboboxOptions v-if="entriesAccessRole.length > 0"
                        class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                        <ComboboxOption v-for="role in entriesAccessRole" :key="role" :value="role" as="accessrole"
                            v-slot="{ active, selected }">
                            <li
                                :class="['relative cursor-default select-none py-2 pl-8 pr-4', active ? 'bg-indigo-600 text-white' : 'text-gray-900']">
                                <span :class="['block truncate', selected && 'font-semibold']">
                                    {{ toTitleCase(role.toLowerCase()) }}
                                </span>

                                <span v-if="selected"
                                    :class="['absolute inset-y-0 left-0 flex items-center pl-1.5', active ? 'text-white' : 'text-indigo-600']">
                                    <CheckIcon class="h-5 w-5" aria-hidden="true" />
                                </span>
                            </li>
                        </ComboboxOption>
                    </ComboboxOptions>
                </div>
            </Combobox>
        </div>
        <div class="flex-initial">
            <button type="button"
                class="rounded-full bg-indigo-600 p-1.5 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
                <XMarkIcon class="h-6 w-6" @click="$emit('deleteRequest', { id: props.init.id })" aria-hidden="true" />
            </button>
        </div>
    </div>
</template>
