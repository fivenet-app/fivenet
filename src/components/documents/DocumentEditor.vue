<script setup lang="ts">
import { computed, ref } from 'vue';
import { useStore } from 'vuex';
import { Quill, QuillEditor } from '@vueup/vue-quill';
import '@vueup/vue-quill/dist/vue-quill.snow.css';
import { getDocStoreClient, handleGRPCError } from '../../grpc';
import { CreateDocumentRequest, GetDocumentRequest, UpdateDocumentRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { DocumentAccess, DocumentJobAccess, DocumentUserAccess, DOC_ACCESS, DOC_CONTENT_TYPE } from '@arpanet/gen/resources/documents/documents_pb';
import { RpcError } from 'grpc-web';
import { dispatchNotification } from '../notification';
import AccessEntry from '../partials/AccessEntry.vue';
import PlusIcon from '@heroicons/vue/20/solid/PlusIcon';
import { useRouter, useRoute } from 'vue-router/auto';
import { Job, JobGrade } from '@arpanet/gen/resources/jobs/jobs_pb';
import { UserShort } from '@arpanet/gen/resources/users/users_pb';

const store = useStore();
const router = useRouter();
const route = useRoute();

const props = defineProps({
    id: {
        required: false,
        type: Number
    },
});

const activeChar = computed(() => store.state.activeChar);

const title = ref('');
const content = ref('');
const categoryID = ref(0);
const closed = ref(false);
const state = ref('');
const isPublic = ref(false);
const access = ref<Map<number, { id: number, type: number, values: { job?: string, char?: number, accessrole?: DOC_ACCESS, minimumrank?: number } }>>(new Map());

const modules = [] as Quill.Module[];

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
    req.setClosed(closed.value);
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
        }).catch((err: RpcError) => {
            console.log(err);
            handleGRPCError(err, route);
        });
}

function editForm(): void {
    const req = new UpdateDocumentRequest();
    req.setDocumentId(9);
    req.setTitle(title.value);
    req.setContent(content.value);
    req.setContentType(DOC_CONTENT_TYPE.HTML);
    req.setClosed(closed.value);
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
        }).catch((err: RpcError) => {
            console.log(err);
            handleGRPCError(err, route);
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
            closed.value = document.getClosed();
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
    }).catch((err: RpcError) => {
        console.log(err);
        handleGRPCError(err, route);
    })
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
    <div class="bg-white">
        <QuillEditor v-model:content="content" contentType="html" toolbar="full" theme="snow" :modules="modules" />
    </div>
    <div class="my-3">
        <h2 class="text-neutral">Access</h2>
        <AccessEntry v-for="entry in access.values()" :key="entry.id"
            :initializationData="entry"
            @typeChange="$event => updateAccessEntryType($event)"
            @nameChange="$event => updateAccessEntryName($event)"
            @rankChange="$event => updateAccessEntryRank($event)"
            @accessChange="$event => updateAccessEntryAccess($event)"
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
