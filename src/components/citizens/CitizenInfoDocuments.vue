<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { DocumentRelation } from '@arpanet/gen/resources/documents/documents_pb';
import { GetUserDocumentsRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { getDocStoreClient } from '../../grpc/grpc';
import { PaginationRequest } from '@arpanet/gen/resources/common/database/database_pb';
import { DOC_RELATION_Util } from '@arpanet/gen/resources/documents/documents.pb_enums';

const relations = ref<Array<DocumentRelation>>([]);

const props = defineProps({
    userId: {
        required: true,
        type: Number,
    },
});

function getUserDocuments(pos: number) {
    if (!props.userId) return;
    if (pos < 0) return;

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
            No Citizen Activities found.
        </p>
    </span>
    <ul v-else class="bg-white">
        <li v-for="relation in relations" :key="relation.getId()">
            <router-link :to="{ name: 'Documents: Info', params: { id: relation.getDocumentId() } }">
                {{ relation.getDocument()?.getTitle() }} (Category: {{ relation.getDocument()?.getCategory() }})
                <br />
                Relation: {{ DOC_RELATION_Util.toEnumKey(relation.getRelation()) }}
                <br />
                {{ relation.getDocument()?.getCreator() }}
            </router-link>
        </li>
    </ul>
</template>
