<script lang="ts">
import { defineComponent, ref } from 'vue';
import { mapState } from 'vuex';
import { Quill, QuillEditor } from '@vueup/vue-quill';
import '@vueup/vue-quill/dist/vue-quill.snow.css';
import { getDocStoreClient, handleGRPCError } from '../../grpc';
import { CreateOrUpdateDocumentRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { DocumentAccess, DocumentJobAccess, DOC_ACCESS, DOC_CONTENT_TYPE } from '@arpanet/gen/resources/documents/documents_pb';
import { RpcError } from 'grpc-web';
import { dispatchNotification } from '../notification';
import { User } from '@arpanet/gen/resources/users/users_pb';
import AccessEntry from '../partials/AccessEntry.vue';
import PlusIcon from '@heroicons/vue/20/solid/PlusIcon';


export default defineComponent({
    components: {
        QuillEditor,
        AccessEntry,
        PlusIcon,
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
            access: [] as { id: number, type: string, values: { name: string, accessrole: string, minimumrank: string } }[]
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

    },
    methods: {
        addAccessEntry(): void {
            if (this.access.length > 4) {
                dispatchNotification({ title: 'Maximum amount of Access entries exceeded', content: 'There can only be a maximum of 5 access entries on a Document', type: 'error' })
                return;
            }

            this.access.push({
                id: this.access.length > 0 ? this.access[this.access.length - 1].id + 1 : 0,
                type: 'jobs',
                values: {
                    name: '',
                    accessrole: '',
                    minimumrank: ''
                }
            })
        },
        updateAccesEntry(data: any): void {
            const accessIndex = this.access.findIndex(e => e.id === data.id);
            if (!accessIndex) return;

            this.access[accessIndex].type = data.selectedAccessType.id

            console.log(this.access);
        },
        submitForm(): void {
            if (this.saving) {
                return;
            }

            this.saving = true;
            const req = new CreateOrUpdateDocumentRequest();
            req.setTitle(this.title);
            req.setContent(this.content);
            req.setContentType(DOC_CONTENT_TYPE.HTML);
            req.setClosed(this.closed);
            req.setState(this.state);
            req.setPublic(this.public);
            // req.setAccess(this.access);

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
                createOrUpdateDocument(req, null).then((resp) => {
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
