<script lang="ts">
import { defineComponent } from 'vue';
import { Quill, QuillEditor } from '@vueup/vue-quill';
import '@vueup/vue-quill/dist/vue-quill.snow.css';
import { getDocStoreClient, handleGRPCError } from '../../grpc';
import { CreateDocumentRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { DOCUMENT_CONTENT_TYPE } from '@arpanet/gen/resources/documents/documents_pb';
import { RpcError } from 'grpc-web';

export default defineComponent({
    components: {
        QuillEditor,
    },
    data() {
        return {
            title: "",
            content: "",
            contentType: DOCUMENT_CONTENT_TYPE.HTML,
            categoryID: 0,
            closed: false,
            state: "",
            public: false,
            targetDocumentID: 0,
        };
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
    props: {
        'targetDocumentID': {
            required: false,
            type: Number,
            default: 0,
        },
    },
    methods: {
        submitForm(): void {
            const req = new CreateDocumentRequest();
            req.setTitle(this.title);
            req.setContent(this.content);
            req.setContentType(this.contentType);
            req.setClosed(this.closed);
            req.setState(this.state);
            req.setPublic(this.public);
            req.setTargetdocumentid(this.targetDocumentID);

            getDocStoreClient().
                createDocument(req, null).then((resp) => {
                    // TODO
                }).catch((err: RpcError) => {
                    handleGRPCError(err, this.$route);
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
    <button
        class="rounded-md bg-white/10 py-2.5 px-3.5 text-sm font-semibold text-white shadow-sm hover:bg-white/20">Submit</button>
</template>
