<script lang="ts" setup>
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
import { CompleteCitizensRequest, CompleteJobsRequest } from '@fivenet/gen/services/completor/completor_pb';
import { Job, JobGrade } from '@fivenet/gen/resources/jobs/jobs_pb';
import { UserShort } from '@fivenet/gen/resources/users/users_pb';
import { DOC_ACCESS } from '@fivenet/gen/resources/documents/documents_pb';
import { toTitleCase } from '~/utils/strings';
import { ArrayElement } from '~/utils/types';
import { DOC_ACCESS_Util } from '@fivenet/gen/resources/documents/documents.pb_enums';

const { $grpc } = useNuxtApp();

const props = defineProps<{
    init: { id: number, type: number, values: { job?: string, char?: number, accessrole?: DOC_ACCESS, minimumrank?: number } }
}>();

const emit = defineEmits<{
    (e: 'typeChange', payload: { id: number, type: number }): void,
    (e: 'nameChange', payload: { id: number, job: Job | undefined, char: UserShort | undefined }): void,
    (e: 'rankChange', payload: { id: number, rank: JobGrade }): void,
    (e: 'accessChange', payload: { id: number, access: DOC_ACCESS }): void,
    (e: 'deleteRequest', payload: { id: number }): void,
}>();

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

let entriesAccessRole = Object.keys(DOC_ACCESS).map(e => { return { id: DOC_ACCESS_Util.fromString(e), value: e } });
const queryAccessRole = ref('');
const selectedAccessRole = ref<ArrayElement<typeof entriesAccessRole>>();

async function findJobs(): Promise<void> {
    const req = new CompleteJobsRequest();
    req.setSearch(queryJob.value);

    const resp = await $grpc.getCompletorClient().
        completeJobs(req, null);
    entriesJobs = resp.getJobsList();
}

async function findChars(): Promise<void> {
    const req = new CompleteCitizensRequest();
    req.setSearch(queryChar.value);

    const resp = await $grpc.getCompletorClient().
        completeCitizens(req, null);
    entriesChars = resp.getUsersList();
}

onMounted(async () => {
    const passedType = accessTypes.find(e => e.id === props.init.type);
    if (passedType) selectedAccessType.value = passedType;

    if (props.init.type === 0 && props.init.values.char !== undefined && props.init.values.accessrole !== undefined) {
        await findChars();
        selectedChar.value = entriesChars.find(char => char.getUserId() === props.init.values.char);
        selectedAccessRole.value = entriesAccessRole.find(type => type.id === props.init.values.accessrole);
    } else if (props.init.type === 1 && props.init.values.job !== undefined && props.init.values.minimumrank !== undefined && props.init.values.accessrole !== undefined) {
        await findJobs();
        selectedJob.value = entriesJobs.find(job => job.getName() === props.init.values.job);
        if (selectedJob.value) entriesMinimumRank = selectedJob.value.getGradesList();
        selectedMinimumRank.value = entriesMinimumRank.find(rank => rank.getGrade() === props.init.values.minimumrank);
        selectedAccessRole.value = entriesAccessRole.find(type => type.id === props.init.values.accessrole);
    }
});

watchDebounced(queryJob, async () => await findJobs(), { debounce: 700, maxWait: 1850 });
watchDebounced(queryChar, async () => await findChars(), { debounce: 700, maxWait: 1850 });

watch(selectedAccessType, () => {
    emit('typeChange', { id: props.init.id, type: selectedAccessType.value.id });

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
    emit('accessChange', { id: props.init.id, access: selectedAccessRole.value.id });
});
</script>

<template>
    <div class="flex flex-row items-center my-2">
        <div class="flex-initial mr-2 w-60">
            <Listbox as="div" v-model="selectedAccessType">
                <div class="relative">
                    <ListboxButton
                        class="block pl-3 text-left w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6">
                        <span class="block truncate">{{ selectedAccessType?.name }}</span>
                        <span class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                            <ChevronDownIcon class="w-5 h-5 text-gray-400" aria-hidden="true" />
                        </span>
                    </ListboxButton>

                    <transition leave-active-class="transition duration-100 ease-in" leave-from-class="opacity-100"
                        leave-to-class="opacity-0">
                        <ListboxOptions
                            class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                            <ListboxOption as="template" v-for="accessType in accessTypes" :key="accessType.id"
                                :value="accessType" v-slot="{ active, selected }">
                                <li
                                    :class="[active ? 'bg-primary-500' : '', 'text-neutral relative cursor-default select-none py-2 pl-8 pr-4']">
                                    <span :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']">{{
                                        accessType.name
                                    }}</span>

                                    <span v-if="selected"
                                        :class="[active ? 'text-neutral' : 'text-primary-500', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
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
                                :display-value="(char: any) => `${char?.getFirstname()} ${char?.getLastname()}`" />
                        </ComboboxButton>

                        <ComboboxOptions v-if="entriesChars.length > 0"
                            class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                            <ComboboxOption v-for="char in entriesChars" :key="char.getIdentifier()" :value="char" as="char"
                                v-slot="{ active, selected }">
                                <li
                                    :class="['relative cursor-default select-none py-2 pl-8 pr-4 text-neutral', active ? 'bg-primary-500' : '']">
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ char.getFirstname() }} {{ char.getLastname() }}
                                    </span>

                                    <span v-if="selected"
                                        :class="[active ? 'text-neutral' : 'text-primary-500', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
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
                                @change="queryJob = $event.target.value" :display-value="(job: any) => job?.getLabel()" />
                        </ComboboxButton>

                        <ComboboxOptions v-if="entriesJobs.length > 0"
                            class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                            <ComboboxOption v-for="job in entriesJobs" :key="job.getName()" :value="job" as="job"
                                v-slot="{ active, selected }">
                                <li
                                    :class="['relative cursor-default select-none py-2 pl-8 pr-4 text-neutral', active ? 'bg-primary-500' : '']">
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ job.getLabel() }}
                                    </span>

                                    <span v-if="selected"
                                        :class="[active ? 'text-neutral' : 'text-primary-500', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
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
                                :display-value="(rank: any) => rank?.getLabel()" />
                        </ComboboxButton>

                        <ComboboxOptions v-if="entriesMinimumRank.length > 0"
                            class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                            <ComboboxOption v-for="rank in entriesMinimumRank" :key="rank.getGrade()" :value="rank"
                                as="minimumrank" v-slot="{ active, selected }">
                                <li
                                    :class="['relative cursor-default select-none py-2 pl-8 pr-4 text-neutral', active ? 'bg-primary-500' : '']">
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ rank.getLabel() }}
                                    </span>

                                    <span v-if="selected"
                                        :class="[active ? 'text-neutral' : 'text-primary-500', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
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
                            :display-value="(role: any) => toTitleCase(role.value?.toLowerCase())" />
                    </ComboboxButton>

                    <ComboboxOptions v-if="entriesAccessRole.length > 0"
                        class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                        <ComboboxOption v-for="role in entriesAccessRole" :key="role.id" :value="role" as="accessrole"
                            v-slot="{ active, selected }">
                            <li
                                :class="['relative cursor-default select-none py-2 pl-8 pr-4 text-neutral', active ? 'bg-primary-500' : '']">
                                <span :class="['block truncate', selected && 'font-semibold']">
                                    {{ toTitleCase(role.value.toLowerCase()) }}
                                </span>

                                <span v-if="selected"
                                    :class="[active ? 'text-neutral' : 'text-primary-500', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
                                    <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                </span>
                            </li>
                        </ComboboxOption>
                    </ComboboxOptions>
                </div>
            </Combobox>
        </div>
        <div class="flex-initial">
            <button type="button"
                class="rounded-full bg-primary-500 p-1.5 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
                <XMarkIcon class="w-6 h-6" @click="$emit('deleteRequest', { id: props.init.id })" aria-hidden="true" />
            </button>
        </div>
    </div>
</template>
