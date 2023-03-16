<script lang="ts">
import { defineComponent } from 'vue';
import { DocumentRelation } from '@arpanet/gen/resources/documents/documents_pb';
import { GetUserDocumentsRequest } from '@arpanet/gen/services/citizenstore/citizenstore_pb';
import { getCitizenStoreClient, handleGRPCError } from '../../grpc';
import { RpcError } from 'grpc-web';

export default defineComponent({
    data() {
        return {
            loading: false,
            relations: [] as Array<DocumentRelation>,
        };
    },
    props: {
        userId: {
            required: true,
            type: Number,
        },
    },
    mounted() {
        this.getUserDocuments(0);
    },
    methods: {
        getUserDocuments(offset: number) {
            if (!this.userId) return;
            if (this.loading) return;
            const req = new GetUserDocumentsRequest();
            req.setUserId(this.userId);

            getCitizenStoreClient().
                getUserDocuments(req, null).
                then((resp) => {
                    resp.getRelationsList
                    this.relations = resp.getRelationsList();
                }).
                catch((err: RpcError) => {
                    handleGRPCError(err, this.$route);
                });
        },
    },
});
</script>

<template>
    <ul class="bg-white">
        <li v-for="relation in relations" :key="relation.getId()">
            {{ relation.getDocument()?.getTitle() }}
        </li>
    </ul>
</template>
