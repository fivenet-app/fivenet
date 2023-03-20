<script setup lang="ts">
import { useRoute } from 'vue-router/auto';
import { ref, onMounted } from 'vue';
import { DocumentRelation } from '@arpanet/gen/resources/documents/documents_pb';
import { GetUserDocumentsRequest } from '@arpanet/gen/services/citizenstore/citizenstore_pb';
import { getCitizenStoreClient, handleGRPCError } from '../../grpc';
import { RpcError } from 'grpc-web';

const route = useRoute();

const relations = ref<Array<DocumentRelation>>([]);

const $props = defineProps({
    userId: {
        required: true,
        type: Number,
    },
});

function getUserDocuments(offset: number) {
    if (!$props.userId) return;
    const req = new GetUserDocumentsRequest();
    req.setUserId($props.userId);

    getCitizenStoreClient().
        getUserDocuments(req, null).
        then((resp) => {
            resp.getRelationsList
            relations.value = resp.getRelationsList();
        }).
        catch((err: RpcError) => {
            handleGRPCError(err, route);
        });
}

onMounted(() => {
    getUserDocuments(0);
});
</script>

<template>
    <span v-if="relations.length === 0">
        <p class="text-sm font-medium text-white">
            No Citizen Activities found.
        </p>
    </span>
    <ul v-else class="bg-white">
        <li v-for="relation in relations" :key="relation.getId()">
            {{ relation.getDocument()?.getTitle() }}
        </li>
    </ul>
</template>
