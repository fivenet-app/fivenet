<script setup lang="ts">
import { computed, ref } from 'vue';
import { useStore } from 'vuex';
import { Quill, QuillEditor } from '@vueup/vue-quill';
import '@vueup/vue-quill/dist/vue-quill.snow.css';
import { getDocStoreClient, handleGRPCError } from '../../grpc';
import { CreateOrUpdateDocumentRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { DocumentAccess, DocumentJobAccess, DOC_ACCESS, DOC_CONTENT_TYPE } from '@arpanet/gen/resources/documents/documents_pb';
import { RpcError } from 'grpc-web';
import { dispatchNotification } from '../notification';
import AccessEntry from '../partials/AccessEntry.vue';
import PlusIcon from '@heroicons/vue/20/solid/PlusIcon';
import { useRouter, useRoute } from 'vue-router/auto';

const store = useStore();
const router = useRouter();
const route = useRoute();

const activeChar = computed(() => store.state.activeChar);

const title = ref('');
const content = ref('');
const categoryID = ref(0);
const closed = ref(false);
const state = ref('');
const isPublic = ref(false);
const access = ref<{ id: number, type: string, values: { name: string, accessrole: string, minimumrank: string } }[]>([]);

const modules = [] as Quill.Module[];

function addAccessEntry(): void {
    if (access.value.length > 4) {
        dispatchNotification({ title: 'Maximum amount of Access entries exceeded', content: 'There can only be a maximum of 5 access entries on a Document', type: 'error' })
        return;
    }

    access.value.push({
        id: access.value.length > 0 ? access.value[access.value.length - 1].id + 1 : 0,
        type: 'jobs',
        values: {
            name: '',
            accessrole: '',
            minimumrank: ''
        }
    })
}

function updateAccesEntry(data: any): void {
    const accessIndex = access.value.findIndex(e => e.id === data.id);
    if (!accessIndex) return;

    access.value[accessIndex].type = data.selectedAccessType.id

    console.log(access);
}

function submitForm(): void {
    const req = new CreateOrUpdateDocumentRequest();
    req.setTitle(title.value);
    req.setContent(content.value);
    req.setContentType(DOC_CONTENT_TYPE.HTML);
    req.setClosed(closed.value);
    req.setState(state.value);
    req.setPublic(isPublic.value);
    // req.setAccess(access);

    const access = new DocumentAccess();
    const jobsAccessList = new Array<DocumentJobAccess>();
    const jobAccess = new DocumentJobAccess();
    jobAccess.setAccess(DOC_ACCESS.VIEW);
    jobAccess.setJob(activeChar.getJob());
    jobsAccessList.push(jobAccess);

    access.setJobsList(jobsAccessList);

    req.setAccess(access);

    getDocStoreClient().
        createOrUpdateDocument(req, null).then((resp) => {
            dispatchNotification({ title: "Document created!", content: "Document has been created." });
            router.push('/documents/' + resp.getId());
        }).catch((err: RpcError) => {
            handleGRPCError(err, route);
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
    <div class="bg-white">
        <QuillEditor v-model:content="content" contentType="html" toolbar="full" theme="snow" :modules="modules" />
    </div>
    <div class="my-3">
        <h2 class="text-neutral">Access</h2>
        <AccessEntry v-for="entry in access" :type="entry.type" :key="entry.id" :id="entry.id"
            @typeChange="$event => updateAccesEntry($event)" />
        <button type="button"
            class="rounded-full bg-indigo-600 p-2 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            data-te-toggle="tooltip" title="Add Permission" @click="addAccessEntry()">
            <PlusIcon class="h-5 w-5" aria-hidden="true" />
        </button>
    </div>
    <button @click="submitForm()"
        class="rounded-md bg-white/10 py-2.5 px-3.5 text-sm font-semibold text-white shadow-sm hover:bg-white/20">Submit</button>
</template>
