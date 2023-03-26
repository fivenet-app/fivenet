<script setup lang="ts">
import { PaginationRequest } from '@arpanet/gen/resources/common/database/database_pb';
import { DocumentComment } from '@arpanet/gen/resources/documents/documents_pb';
import { GetDocumentCommentsRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { ref, onMounted } from 'vue';
import { getDocStoreClient } from '../../grpc/grpc';
import DocumentCommentNew from './DocumentCommentNew.vue';
import DocumentCommentEntry from './DocumentCommentEntry.vue';

const props = defineProps({
    documentId: {
        required: true,
        type: Number,
    },
    comments: {
        required: false,
        type: Array<DocumentComment>,
    }
});

const comments = ref<DocumentComment[]>([]);

// Document Comments
function getDocumentComments(): void {
    const req = new GetDocumentCommentsRequest();
    req.setPagination((new PaginationRequest()).setOffset(0));
    req.setDocumentId(props.documentId!);

    getDocStoreClient().
        getDocumentComments(req, null).
        then((resp) => {
            comments.value = resp.getCommentsList();
        });
}

onMounted(() => {
    if (props.documentId !== undefined && props.comments === undefined) {
        getDocumentComments();
    }
});
</script>

<template>
    <DocumentCommentNew :document-id="documentId" @added="(c: DocumentComment) => comments.push(c)" />
    <div class="flow-root px-4 rounded-lg bg-base-800 text-neutral">
        <ul role="list" class="divide-y divide-gray-200">
            <DocumentCommentEntry v-for="com in (props.comments ?? comments)" :key="com.getId()" :comment="com" />
        </ul>
    </div>
</template>
