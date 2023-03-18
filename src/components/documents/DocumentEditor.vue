<script lang="ts">
import { defineComponent, ref, computed, ComputedRef } from 'vue';
import { mapState } from 'vuex';
import { Quill, QuillEditor } from '@vueup/vue-quill';
import '@vueup/vue-quill/dist/vue-quill.snow.css';
import { getDocStoreClient, handleGRPCError, getCompletorClient } from '../../grpc';
import { CreateDocumentRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { DocumentAccess, DocumentJobAccess, DOC_ACCESS, DOC_CONTENT_TYPE } from '@arpanet/gen/resources/documents/documents_pb';
import { RpcError } from 'grpc-web';
import { dispatchNotification } from '../notification';
import { User } from '@arpanet/gen/resources/users/users_pb';
import { PlusIcon, CheckIcon, ChevronUpDownIcon } from '@heroicons/vue/20/solid'
import { watchDebounced } from '@vueuse/core';
import { CompleteCharNamesRequest, CompleteJobNamesRequest } from '@arpanet/gen/services/completor/completor_pb';
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

export default defineComponent({
    components: {
        QuillEditor,
        PlusIcon,
        CheckIcon,
        ChevronUpDownIcon,
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
            saving: false,
            title: "",
            content: "",
            categoryID: 0,
            closed: false,
            state: "",
            public: false,
            accessTypes: [
                { id: 0, name: 'Citizen' },
                { id: 1, name: 'Jobs' },
            ],
            selectedAccessType: ref<null | { id: number, name: string }>(null),
            people: [] as { id: string | number, label: string }[],
            query: { value: '' },
            selectedPerson: ref(null),
        };
    },
    computed: {
        ...mapState({
            activeChar: 'activeChar',
        }),
    },
    updated() {
        console.log(this.content);
    },
    setup() {
        const modules = [] as Quill.Module[];

        return {
            modules,
        };
    },
    mounted() {
        this.selectedAccessType = this.accessTypes[1];

        watchDebounced(this.query, () => {
            if (this.selectedAccessType?.id === 0) {
                this.findChars();
            } else {
                this.findJobs();
            }
        }, { debounce: 750, maxWait: 2000 });
    },
    props: {
        'targetDocumentID': {
            required: false,
            type: Number,
            default: 0,
        },
    },
    methods: {
        findJobs(): void {
            const req = new CompleteJobNamesRequest();
            req.setSearch(this.query.value);

            getCompletorClient().completeJobNames(req, null).then((resp) => {
                this.people = [];
                resp.getJobsList().forEach((job) => {
                    this.people.push({ id: job.getName(), label: job.getLabel() })
                })
                console.log(this.people);
            }).catch((err: RpcError) => {
                handleGRPCError(err, this.$route);
            })
        },
        findChars(): void {
            const req = new CompleteCharNamesRequest();
            req.setSearch(this.query.value);

            getCompletorClient().completeCharNames(req, null).then((resp) => {
                this.people = [];
                resp.getUsersList().forEach((user) => {
                    this.people.push({ id: user.getUserId(), label: `${user.getFirstname()} ${user.getLastname()}` })
                })
                console.log(this.people);
            }).catch((err: RpcError) => {
                handleGRPCError(err, this.$route);
            })
        },
        submitForm(): void {
            if (this.saving) {
                return;
            }

            this.saving = true;
            const req = new CreateDocumentRequest();
            req.setTitle(this.title);
            req.setContent(this.content);
            req.setContentType(DOC_CONTENT_TYPE.HTML);
            req.setClosed(this.closed);
            req.setState(this.state);
            req.setPublic(this.public);
            req.setTargetDocumentId(this.targetDocumentID);

            const access = new DocumentAccess();
            const jobsAccessList = new Array<DocumentJobAccess>();
            const jobAccess = new DocumentJobAccess();
            jobAccess.setAccess(DOC_ACCESS.VIEW);
            const activeChar = this.activeChar as null | User;
            jobAccess.setJob(activeChar?.getJob());
            jobsAccessList.push(jobAccess);

            access.setJobsList(jobsAccessList);

            req.setAccess(access);

            getDocStoreClient().
                createDocument(req, null).then((resp) => {
                    dispatchNotification({ title: "Document created!", content: "Document has been created." });
                    this.saving = false;
                    this.$router.push('/documents/' + resp.getId());
                }).catch((err: RpcError) => {
                    handleGRPCError(err, this.$route);
                    this.saving = false;
                });
        },
    },
});
</script>

<route lang="json">
{
    "name": "documents-new",
    "meta": {
        "requiresAuth": true
    }
}
</route>

<template>
    <div
        class="rounded-md px-3 pt-2.5 pb-1.5 shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-indigo-600 bg-white">
        <label for="name" class="block text-xs font-medium text-gray-900">Title</label>
        <input v-model="title" type="text" name="name"
            class="block w-full border-0 p-0 text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6"
            placeholder="Document Title" />
    </div>
    <div class="bg-white">
        <QuillEditor v-model:content="content" contentType="html" toolbar="full" theme="snow" :modules="modules" />
    </div>
    <div class="my-3">
        <h2 class="text-neutral">Access</h2>
        <div class="flex flex-col">
            <div class="flex-1">
                <Listbox as="div" v-model="selectedAccessType">
                    <ListboxLabel class="block text-sm font-medium leading-6 text-neutral">Type</ListboxLabel>
                    <div class="relative mt-2">
                        <ListboxButton
                            class="relative w-full cursor-default rounded-md bg-white py-1.5 pl-3 pr-10 text-left text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-600 sm:text-sm sm:leading-6">
                            <span class="block truncate">{{ selectedAccessType?.name }}</span>
                            <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                                <ChevronUpDownIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
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
            <div class="flex-1">
                <Combobox as="div" v-model="selectedPerson">
                    <div class="relative mt-2">
                        <ComboboxInput
                            class="w-full rounded-md border-0 bg-white py-1.5 pl-3 pr-10 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                            @change="query.value = $event.target.value" @click="selectedAccessType?.id === 0 ? findChars() : findJobs()" :display-value="(person: any) => person?.label" />
                        <ComboboxButton
                            class="absolute inset-y-0 right-0 flex items-center rounded-r-md px-2 focus:outline-none">
                            <ChevronUpDownIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
                        </ComboboxButton>

                        <ComboboxOptions v-if="people.length > 0"
                            class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                            <ComboboxOption v-for="person in people" :key="person.id" :value="person" as="template"
                                v-slot="{ active, selected }">
                                <li
                                    :class="['relative cursor-default select-none py-2 pl-8 pr-4', active ? 'bg-indigo-600 text-white' : 'text-gray-900']">
                                    <span :class="['block truncate', selected && 'font-semibold']">
                                        {{ person.label }}
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
    </div>
    <button @click="submitForm()"
        class="rounded-md bg-white/10 py-2.5 px-3.5 text-sm font-semibold text-white shadow-sm hover:bg-white/20">Submit</button>
</template>
