<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { DocumentRelation } from '@arpanet/gen/resources/documents/documents_pb';
import { GetUserDocumentsRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { getDocStoreClient, handleRPCError } from '../../grpc/grpc';
import { PaginationRequest } from '@arpanet/gen/resources/common/database/database_pb';
import DocumentRelations from '../documents/DocumentRelations.vue';
import { DocumentTextIcon } from '@heroicons/vue/24/outline';
import { RpcError } from 'grpc-web';

const relations = ref<Array<DocumentRelation>>([]);

const props = defineProps({
    userId: {
        required: true,
        type: Number,
    },
});

async function getUserDocuments(pos: number): Promise<void> {
    if (!props.userId) return;
    if (pos < 0) pos = 0;

    const req = new GetUserDocumentsRequest();
    req.setPagination((new PaginationRequest()).setOffset(pos))
    req.setUserId(props.userId);

    try {
        const resp = await getDocStoreClient().
            getUserDocuments(req, null);
        relations.value = resp.getRelationsList();
    } catch (e) {
        handleRPCError(e as RpcError);
        return;
    }
}

onMounted(async () => {
    getUserDocuments(0);
});
</script>

<template>
    <div class="mt-3">
        <button v-if="relations.length == 0" type="button"
            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
            disabled>
            <DocumentTextIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-gray-300">
                No User Documents found.
            </span>
        </button>
        <DocumentRelations v-else :relations="relations" />
    </div>
</template>
