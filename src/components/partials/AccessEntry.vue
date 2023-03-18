import { defineComponent } from 'vue';
<script lang="ts">
import { defineComponent, ref } from 'vue';
import {
    Listbox,
    ListboxButton,
    ListboxLabel,
    ListboxOption,
    ListboxOptions,
    Combobox,
    ComboboxButton,
    ComboboxInput,
    ComboboxOption,
    ComboboxOptions
} from '@headlessui/vue'
import {
    PlusIcon,
    CheckIcon,
    ChevronDownIcon
} from '@heroicons/vue/20/solid';
import { watchDebounced } from '@vueuse/core';
import { CompleteCharNamesRequest, CompleteJobNamesRequest } from '@arpanet/gen/services/completor/completor_pb';
import { getCompletorClient, handleGRPCError } from '../../grpc';
import { RpcError } from 'grpc-web';
import { Job } from '@arpanet/gen/resources/jobs/jobs_pb';
import { UserShort } from '@arpanet/gen/resources/users/users_pb';

export default defineComponent({
    components: {
        PlusIcon,
        CheckIcon,
        ChevronDownIcon,
        Listbox,
        ListboxButton,
        ListboxLabel,
        ListboxOption,
        ListboxOptions,
        Combobox,
        ComboboxButton,
        ComboboxInput,
        ComboboxOption,
        ComboboxOptions,
    },
    data() {
        return {
            accessTypes: [
                { id: 0, name: 'Citizen' },
                { id: 1, name: 'Jobs' },
            ],
            selectedAccessType: ref<null | { id: number, name: string }>(null),
            entriesChars: [] as UserShort[],
            queryChar: { value: '' },
            selectedChar: ref<null | UserShort>(null),
            entriesJobs: [] as Job[],
            queryJob: { value: '' },
            selectedJob: ref<null | Job>(null),
            entriesAccessRole: [] as { id: string | number, label: string }[],
            queryAccessRole: { value: '' },
            selectedAccessRole: ref(null),
            entriesMinimumRank: [] as { id: string | number, label: string }[],
            queryMinimumRank: { value: '' },
            selectedMinimumRank: ref(null),
        }
    },
    mounted() {
        this.selectedAccessType = this.accessTypes[1];

        watchDebounced(this.queryJob, () => {
            if (this.selectedAccessType?.id === 0) {
                this.findChars();
            } else {
                this.findJobs();
            }
        }, { debounce: 750, maxWait: 2000 });
    },
    methods: {
        findJobs(): void {
            const req = new CompleteJobNamesRequest();
            req.setSearch(this.queryJob.value);

            getCompletorClient().completeJobNames(req, null).then((resp) => {
                this.entriesJobs = resp.getJobsList();
            }).catch((err: RpcError) => {
                handleGRPCError(err, this.$route);
            })
        },
        findChars(): void {
            const req = new CompleteCharNamesRequest();
            req.setSearch(this.queryJob.value);

            getCompletorClient().completeCharNames(req, null).then((resp) => {
                this.entriesChars = resp.getUsersList();
            }).catch((err: RpcError) => {
                handleGRPCError(err, this.$route);
            })
        },
    }
})

</script>

<template>
    <div class="flex flex-row">
        <div class="flex-initial w-60 mr-2">
            <Listbox as="div" v-model="selectedAccessType">
                <div class="relative mt-2">
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
                <Combobox as="div" v-model="selectedJob">
                    <div class="relative mt-2">
                        <ComboboxButton as="div">
                            <ComboboxInput
                                class="w-full rounded-md border-0 bg-white py-1.5 pl-3 pr-10 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                                @change="queryJob.value = $event.target.value"
                                @click="selectedAccessType?.id === 0 ? findChars() : findJobs()"
                                :display-value="(entry: any) => entry?.getLabel()" />
                        </ComboboxButton>

                        <ComboboxOptions v-if="entriesJobs.length > 0"
                            class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                            <ComboboxOption v-for="entry in entriesJobs" :key="entry.getName()" :value="entry" as="job"
                                v-slot="{ active, selected }">
                                <li
                                    :class="['relative cursor-default select-none py-2 pl-8 pr-4', active ? 'bg-indigo-600 text-white' : 'text-gray-900']">
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ entry.getLabel() }}
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
                    <div class="relative mt-2">
                        <ComboboxButton as="div">
                            <ComboboxInput
                                class="w-full rounded-md border-0 bg-white py-1.5 pl-3 pr-10 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                                @change="queryJob.value = $event.target.value"
                                @click="selectedAccessType?.id === 0 ? findChars() : findJobs()"
                                :display-value="(entry: any) => entry?.getLabel()" />
                        </ComboboxButton>

                        <ComboboxOptions v-if="entriesJobs.length > 0"
                            class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                            <ComboboxOption v-for="entry in entriesJobs" :key="entry.getName()" :value="entry" as="job"
                                v-slot="{ active, selected }">
                                <li
                                    :class="['relative cursor-default select-none py-2 pl-8 pr-4', active ? 'bg-indigo-600 text-white' : 'text-gray-900']">
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ entry.getLabel() }}
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
            <div class="flex-1 mr-2" :hidden="selectedAccessType?.id === 0">
                <Combobox as="div" v-model="selectedMinimumRank">
                    <div class="relative mt-2">
                        <ComboboxButton as="div">
                            <ComboboxInput
                                class="w-full rounded-md border-0 bg-white py-1.5 pl-3 pr-10 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                                @change="queryMinimumRank.value = $event.target.value"
                                @click="selectedAccessType?.id === 0 ? findChars() : findJobs()"
                                :display-value="(person: any) => person?.label" />
                        </ComboboxButton>

                        <ComboboxOptions v-if="entriesMinimumRank.length > 0"
                            class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                            <ComboboxOption v-for="entry in entriesMinimumRank" :key="entry.id" :value="entry" as="minimumrank"
                                v-slot="{ active, selected }">
                                <li
                                    :class="['relative cursor-default select-none py-2 pl-8 pr-4', active ? 'bg-indigo-600 text-white' : 'text-gray-900']">
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ entry.label }}
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
        <div class="flex-inital w-60">
                <Combobox as="div" v-model="selectedAccessRole">
                    <div class="relative mt-2">
                        <ComboboxButton as="div">
                            <ComboboxInput
                                class="w-full rounded-md border-0 bg-white py-1.5 pl-3 pr-10 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                                @change="queryJob.value = $event.target.value"
                                @click="selectedAccessType?.id === 0 ? findChars() : findJobs()"
                                :display-value="(entry: any) => entry?.label" />
                        </ComboboxButton>

                        <ComboboxOptions v-if="entriesAccessRole.length > 0"
                            class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                            <ComboboxOption v-for="entry in entriesAccessRole" :key="entry.id" :value="entry" as="accessrole"
                                v-slot="{ active, selected }">
                                <li
                                    :class="['relative cursor-default select-none py-2 pl-8 pr-4', active ? 'bg-indigo-600 text-white' : 'text-gray-900']">
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ entry.label }}
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
    <button type="button"
        class="rounded-full bg-indigo-600 p-2 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
        data-te-toggle="tooltip" title="Add Permission">
        <PlusIcon class="h-5 w-5" aria-hidden="true" />
    </button>
</template>
