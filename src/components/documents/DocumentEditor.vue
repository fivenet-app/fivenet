<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useStore } from '../../store/store';
import { Quill, QuillEditor } from '@vueup/vue-quill';
import '@vueup/vue-quill/dist/vue-quill.snow.css';
import { getCompletorClient, getDocStoreClient } from '../../grpc/grpc';
import { CreateDocumentRequest, GetDocumentRequest, UpdateDocumentRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { DocumentAccess, DocumentCategory, DocumentJobAccess, DocumentUserAccess, DOC_ACCESS, DOC_CONTENT_TYPE } from '@arpanet/gen/resources/documents/documents_pb';
import { dispatchNotification } from '../notification';
import AccessEntry from '../partials/AccessEntry.vue';
import {
    PlusIcon,
    ChevronDownIcon,
    CheckIcon,
} from '@heroicons/vue/20/solid';
import { useRouter } from 'vue-router/auto';
import { Job, JobGrade } from '@arpanet/gen/resources/jobs/jobs_pb';
import { UserShort } from '@arpanet/gen/resources/users/users_pb';
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
import { CompleteDocumentCategoryRequest } from '@arpanet/gen/services/completor/completor_pb';
import { watchDebounced } from '@vueuse/core';

const store = useStore();
const router = useRouter();

const props = defineProps({
    id: {
        required: false,
        type: Number
    },
});

const activeChar = computed(() => store.state.auth?.activeChar);

const openclose = [
    { id: 0, label: 'Open', closed: false },
    { id: 1, label: 'Closed', closed: true },
]

const title = ref('');
const content = ref('');
const closed = ref(openclose[0]);
const state = ref('');
const isPublic = ref(false);
const access = ref<Map<number, { id: number, type: number, values: { job?: string, char?: number, accessrole?: DOC_ACCESS, minimumrank?: number } }>>(new Map());

let entriesCategory = [] as DocumentCategory[];
const queryCategory = ref('');
const selectedCategory = ref<DocumentCategory | undefined>(undefined);

const modules = [] as Quill.Module[];

onMounted(() => {
    findCategories();
});

watchDebounced(queryCategory, () => findCategories(), { debounce: 750, maxWait: 2000 });

function findCategories(): void {
    const req = new CompleteDocumentCategoryRequest();
    req.setSearch(queryCategory.value);

    getCompletorClient().
        completeDocumentCategory(req, null).
        then((resp) => {
            entriesCategory = resp.getCategoriesList();
        });
}

function addAccessEntry(): void {
    if (access.value.size > 4) {
        dispatchNotification({ title: 'Maximum amount of Access entries exceeded', content: 'There can only be a maximum of 5 access entries on a Document', type: 'error' })
        return;
    }

    let id = access.value.size > 0 ? [...access.value.keys()].pop() as number + 1 : 0;
    access.value.set(id, {
        id,
        type: 1,
        values: {}
    })
}

function removeAccessEntry(event: {
    id: number
}): void {
    // access.value.delete(event.id);
}

function updateAccessEntryType(event: {
    id: number,
    type: number
}): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.type = event.type;
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryName(event: {
    id: number,
    job?: Job,
    char?: UserShort
}): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    if (event.job) {
        accessEntry.values.job = event.job.getName();
        accessEntry.values.char = undefined;
    } else if (event.char) {
        accessEntry.values.job = undefined;
        accessEntry.values.char = event.char.getUserId();
    }

    access.value.set(event.id, accessEntry);
}

function updateAccessEntryRank(event: {
    id: number,
    rank: JobGrade
}): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.values.minimumrank = event.rank.getGrade();
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryAccess(event: {
    id: number,
    access: DOC_ACCESS
}): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.values.accessrole = event.access;
    access.value.set(event.id, accessEntry);
}

function submitForm(): void {
    const req = new CreateDocumentRequest();
    req.setTitle(title.value);
    req.setContent(content.value);
    req.setContentType(DOC_CONTENT_TYPE.HTML);
    req.setClosed(closed.value.closed);
    req.setState(state.value);
    req.setPublic(isPublic.value);

    const reqAccess = new DocumentAccess();
    access.value.forEach(entry => {
        if (entry.values.accessrole === undefined) return;

        if (entry.type === 0) {
            if (!entry.values.char) return;

            const user = new DocumentUserAccess();
            user.setAccess(DOC_ACCESS[entry.values.accessrole]);
            user.setUserId(entry.values.char);

            reqAccess.addUsers(user);
        } else if (entry.type === 1) {
            if (!entry.values.job) return;

            const job = new DocumentJobAccess();
            job.setJob(entry.values.job.getName());
            job.setMinimumgrade(entry.values.minimumrank ? entry.values.minimumrank : 0);
            job.setAccess(DOC_ACCESS[entry.values.accessrole]);
            job.setCreatorId(activeChar.value);

            reqAccess.addJobs(job);
        }
    });
    req.setAccess(reqAccess);

    // const access = new DocumentAccess();
    // const jobsAccessList = new Array<DocumentJobAccess>();
    // const jobAccess = new DocumentJobAccess();
    // jobAccess.setAccess(DOC_ACCESS.VIEW);
    // jobAccess.setJob(activeChar.getJob());
    // jobsAccessList.push(jobAccess);
    // access.setJobsList(jobsAccessList);

    getDocStoreClient().
        createDocument(req, null).
        then((resp) => {
            dispatchNotification({ title: "Document created!", content: "Document has been created." });
            router.push('/documents/' + resp.getDocumentId());
        });
}

function editForm(): void {
    const req = new UpdateDocumentRequest();
    req.setDocumentId(9);
    req.setTitle(title.value);
    req.setContent(content.value);
    req.setContentType(DOC_CONTENT_TYPE.HTML);
    req.setClosed(closed.value.closed);
    req.setState(state.value);
    req.setPublic(isPublic.value);

    const reqAccess = new DocumentAccess();
    access.value.forEach(entry => {
        if (entry.values.accessrole === undefined) return;

        if (entry.type === 0) {
            if (!entry.values.char) return;

            const user = new DocumentUserAccess();
            console.log(DOC_ACCESS[entry.values.accessrole])
            console.log(entry.values.accessrole)
            user.setAccess(DOC_ACCESS[entry.values.accessrole]);
            user.setUserId(entry.values.char);

            reqAccess.addUsers(user);
        } else if (entry.type === 1) {
            if (!entry.values.job) return;

            const job = new DocumentJobAccess();
            job.setJob(entry.values.job.getName());
            job.setMinimumgrade(entry.values.minimumrank ? entry.values.minimumrank : 0);
            job.setAccess(DOC_ACCESS[entry.values.accessrole]);
            job.setCreatorId(activeChar.value);

            reqAccess.addJobs(job);
        }
    });
    req.setAccess(reqAccess);

    // const access = new DocumentAccess();
    // const jobsAccessList = new Array<DocumentJobAccess>();
    // const jobAccess = new DocumentJobAccess();
    // jobAccess.setAccess(DOC_ACCESS.VIEW);
    // jobAccess.setJob(activeChar.getJob());
    // jobsAccessList.push(jobAccess);
    // access.setJobsList(jobsAccessList);

    getDocStoreClient().
        updateDocument(req, null).
        then((resp) => {
            dispatchNotification({ title: "Document updated!", content: "Document has been updated." });
        });
}

if (props.id) {
    const req = new GetDocumentRequest();
    req.setDocumentId(props.id);

    getDocStoreClient().getDocument(req, null).then((resp) => {
        const document = resp.getDocument();
        const docAccess = resp.getAccess();

        if (document) {
            title.value = document.getTitle();
            content.value = document.getContent();
            closed.value = openclose.find(e => e.closed === document.getClosed()) as { id: number; label: string; closed: boolean; };
            state.value = document.getState();
            isPublic.value = document.getPublic();
        };

        if (docAccess) {
            let accessId = 0;
            console.log(docAccess.toObject());

            docAccess.getUsersList().forEach(user => {
                access.value.set(accessId, { id: accessId, type: 0, values: { char: user.getUserId(), accessrole: user.getAccess() } })
                accessId++;
            });

            docAccess.getJobsList().forEach(job => {
                access.value.set(accessId, { id: accessId, type: 1, values: { job: job.getJob(), accessrole: job.getAccess(), minimumrank: job.getMinimumgrade() } })
                accessId++;
            });
        }
    });
}
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
    <div class="flex flex-row">
        <div class="flex-1">
            <!-- Category -->
            <Combobox as="div" v-model="selectedCategory">
                <div class="relative">
                    <ComboboxButton as="div">
                        <ComboboxInput
                            class="w-full rounded-md border-0 bg-white py-1.5 pl-3 pr-10 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                            @change="queryCategory = $event.target.value"
                            :display-value="(category: any) => category?.getName()" />
                    </ComboboxButton>

                    <ComboboxOptions v-if="entriesCategory.length > 0"
                        class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                        <ComboboxOption v-for="category in entriesCategory" :key="category.getId()" :value="category"
                            as="category" v-slot="{ active, selected }">
                            <li
                                :class="['relative cursor-default select-none py-2 pl-8 pr-4', active ? 'bg-indigo-600 text-white' : 'text-gray-900']">
                                <span :class="['block truncate', selected && 'font-semibold']">
                                    {{ category.getName() }}
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
        <div class="flex-1 rounded-md px-3 pt-2.5 pb-1.5 shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-indigo-600 bg-white">
            <!-- State -->
            <input v-model="state" type="text" name="state"
            class="block w-full border-0 p-0 text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6"
            placeholder="Document State" />
        </div>
        <div class="flex-1">
            <!-- Open/Close -->
            <Listbox as="div" v-model="closed">
                <div class="relative">
                    <ListboxButton
                        class="relative w-full cursor-default rounded-md bg-white py-1.5 pl-3 pr-10 text-left text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-600 sm:text-sm sm:leading-6">
                        <span class="block truncate">{{ openclose.find(e => e.closed === closed.closed)?.label }}</span>
                        <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                            <ChevronDownIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
                        </span>
                    </ListboxButton>

                    <transition leave-active-class="transition ease-in duration-100" leave-from-class="opacity-100"
                        leave-to-class="opacity-0">
                        <ListboxOptions
                            class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                            <ListboxOption as="template" v-for="type in openclose" :key="type.id" :value="type"
                                v-slot="{ active, selected }">
                                <li
                                    :class="[active ? 'bg-indigo-600 text-white' : 'text-gray-900', 'relative cursor-default select-none py-2 pl-8 pr-4']">
                                    <span :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']">{{
                                        type.label
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
    </div>
    <div class="bg-white">
        <QuillEditor v-model:content="content" contentType="html" toolbar="full" theme="snow" :modules="modules" />
    </div>
    <div class="my-3">
        <h2 class="text-neutral">Access</h2>
        <AccessEntry v-for="entry in access.values()" :key="entry.id" :init="entry"
            @typeChange="$event => updateAccessEntryType($event)" @nameChange="$event => updateAccessEntryName($event)"
            @rankChange="$event => updateAccessEntryRank($event)" @accessChange="$event => updateAccessEntryAccess($event)"
            @deleteRequest="$event => removeAccessEntry($event)" />
        <button type="button"
            class="rounded-full bg-indigo-600 p-2 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            data-te-toggle="tooltip" title="Add Permission" @click="addAccessEntry()">
            <PlusIcon class="h-5 w-5" aria-hidden="true" />
        </button>
    </div>
    <button @click="submitForm()"
        class="rounded-md bg-white/10 py-2.5 px-3.5 text-sm font-semibold text-white shadow-sm hover:bg-white/20">Submit</button>
    <button @click="editForm()"
        class="rounded-md bg-white/10 py-2.5 px-3.5 text-sm font-semibold text-white shadow-sm hover:bg-white/20">Edit</button>
</template>
