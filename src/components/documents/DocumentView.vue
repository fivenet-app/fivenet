<script lang="ts">
import { defineComponent } from 'vue';
import { Document, GetDocumentRequest } from '@arpanet/gen/documents/documents_pb';
import { getDocumentsClient, handleGRPCError } from '../../grpc';
import { RpcError } from 'grpc-web';

export default defineComponent({
    data() {
        return {
            document: undefined as undefined | Document,
        };
    },
    props: {
        documentID: {
            required: true,
            type: Number,
        },
    },
    mounted() {
        this.getDocument();
    },
    methods: {
        getDocument(): void {
            const req = new GetDocumentRequest();
            req.setId(this.documentID);

            getDocumentsClient().
                getDocument(req, null).
                then((resp) => {
                    this.document = resp.getDocument();
                }).
                catch((err: RpcError) => {
                    handleGRPCError(err, this.$route);
                });
        }
    },
});
</script>

<template>
    {{  document }}
</template>
