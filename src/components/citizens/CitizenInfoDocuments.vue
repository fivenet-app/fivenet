<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { DocumentRelation } from '@arpanet/gen/resources/documents/documents_pb';
import { GetUserDocumentsRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { getDocStoreClient } from '../../grpc/grpc';
import { PaginationRequest } from '@arpanet/gen/resources/common/database/database_pb';
import DocumentRelations from '../documents/DocumentRelations.vue';

const relations = ref<Array<DocumentRelation>>([]);

const props = defineProps({
    userId: {
        required: true,
        type: Number,
    },
});

function getUserDocuments(pos: number) {
    if (!props.userId) return;
    if (pos < 0) pos = 0;

    const req = new GetUserDocumentsRequest();
    req.setPagination((new PaginationRequest()).setOffset(pos))
    req.setUserId(props.userId);

    getDocStoreClient().
        getUserDocuments(req, null).
        then((resp) => {
            relations.value = resp.getRelationsList();
        });
}

onMounted(() => {
    getUserDocuments(0);
});
</script>

<template>
    <span v-if="relations.length === 0">
        <p class="text-sm font-medium text-white">
            No Documents found.
        </p>
    </span>
    <DocumentRelations v-else :showDocument="true" :relations="relations" />
</template>
