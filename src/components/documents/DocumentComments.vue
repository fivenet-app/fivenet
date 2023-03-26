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
        default: new Array<DocumentComment>(),
    }
});

// Document Comments
function getDocumentComments(): void {
    const req = new GetDocumentCommentsRequest();
    req.setPagination((new PaginationRequest()).setOffset(0));
    req.setDocumentId(props.documentId!);

    getDocStoreClient().
        getDocumentComments(req, null).
        then((resp) => {
            resp.getCommentsList().forEach((v) => {
                props.comments.push(v);
            });
        });
}

function addComment(comment: DocumentComment) {
    props.comments.unshift(comment);
}

function removedComment(comment: DocumentComment) {
    const idx = props.comments.findIndex((c) => {
        return c.getId() === comment.getId();
    });
    if (idx > -1) {
        props.comments.splice(idx, 1);
    }
}

onMounted(() => {
    if (props.documentId !== undefined && props.comments === undefined) {
        getDocumentComments();
    }
});
</script>

<template>
    <div class="pb-2">
        <DocumentCommentNew :document-id="documentId" @added="(c: DocumentComment) => addComment(c)" />
    </div>
    <div class="flow-root px-4 rounded-lg bg-base-800 text-neutral">
        <ul role="list" class="divide-y divide-gray-200">
            <DocumentCommentEntry v-for="com in comments" :key="com.getId()" :comment="com"
                @removed="(c: DocumentComment) => removedComment(c)" />
        </ul>
    </div>
</template>
