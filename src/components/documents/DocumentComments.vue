<script setup lang="ts">
import { PaginationRequest, PaginationResponse } from '@fivenet/gen/resources/common/database/database_pb';
import { DocumentComment } from '@fivenet/gen/resources/documents/documents_pb';
import { DeleteDocumentCommentRequest, GetDocumentCommentsRequest, PostDocumentCommentRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { computed, ref, onMounted } from 'vue';
import DocumentCommentEntry from './DocumentCommentEntry.vue';
import { useAuthStore } from '~/store/auth';
import { ChatBubbleLeftEllipsisIcon } from '@heroicons/vue/20/solid';
import TablePagination from '~/components/partials/TablePagination.vue';
import { RpcError } from 'grpc-web';

const { $grpc } = useNuxtApp();
const store = useAuthStore();

const activeChar = computed(() => store.$state.activeChar);

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

const pagination = ref<PaginationResponse>();
const offset = ref(0);

// Document Comments
async function getDocumentComments(): Promise<void> {
    const req = new GetDocumentCommentsRequest();
    req.setPagination((new PaginationRequest()).setOffset(offset.value));
    req.setDocumentId(props.documentId!);

    try {
        const resp = await $grpc.getDocStoreClient().
            getDocumentComments(req, null);

        resp.getCommentsList().forEach((v) => {
            pagination.value = resp.getPagination();
            props.comments.push(v);
        });
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);
        return;
    }
}
const message = ref('');

async function addComment(): Promise<void> {
    const req = new PostDocumentCommentRequest();
    const com = new DocumentComment();
    com.setDocumentId(props.documentId);
    com.setComment(message.value);
    req.setComment(com);

    try {
        const resp = await $grpc.getDocStoreClient().
            postDocumentComment(req, null);

        com.setId(resp.getId());
        com.setCreatorId(activeChar.value!.getUserId());
        com.setCreator(activeChar.value!);

        props.comments.unshift(com);
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);
        return;
    }
}

async function removeComment(comment: DocumentComment): Promise<void> {
    const idx = props.comments.findIndex((c) => {
        return c.getId() === comment.getId();
    });

    const req = new DeleteDocumentCommentRequest();
    req.setCommentId(comment.getId());

    try {
        await $grpc.getDocStoreClient().
            deleteDocumentComment(req, null);

        if (idx > -1) {
            props.comments.splice(idx, 1);
        }
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);
        return;
    }
}

const commentInput = ref<HTMLInputElement | null>(null);
function focusComment(): void {
    if (commentInput.value) {
        commentInput.value.focus();
    }
}

watch(offset, async () => getDocumentComments());
onMounted(async () => {
    if (props.documentId !== undefined && props.comments === undefined) {
        getDocumentComments();
    }
});
</script>

<template>
    <div class="pb-2">
        <div v-can="'DocStoreService.PostDocumentComment'" class="flex items-start space-x-4">
            <div class="min-w-0 flex-1">
                <form @submit.prevent="addComment()" class="relative">
                    <div
                        class="overflow-hidden rounded-lg shadow-sm ring-1 ring-inset ring-gray-500 focus-within:ring-2 focus-within:ring-indigo-600">
                        <label for="comment" class="sr-only">Add your comment</label>
                        <textarea rows="3" name="comment" id="comment"
                            class="block w-full resize-none border-0 bg-transparent text-gray-50 placeholder:text-gray-400 focus:ring-0 sm:py-1.5 sm:text-sm sm:leading-6"
                            ref="commentInput" v-model="message" placeholder="Add your comment..." />

                        <!-- Spacer element to match the height of the toolbar -->
                        <div class="py-2" aria-hidden="true">
                            <!-- Matches height of button in toolbar (1px border + 36px content height) -->
                            <div class="py-px">
                                <div class="h-9" />
                            </div>
                        </div>
                    </div>

                    <div class="absolute inset-x-0 bottom-0 flex justify-between py-2 pl-3 pr-2">
                        <div class="flex items-center space-x-5"></div>
                        <div class="flex-shrink-0">
                            <button type="submit"
                                class="inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Post</button>
                        </div>
                    </div>
                </form>
            </div>
        </div>
    </div>
    <div class="bg-base-800">
        <button v-if="comments.length == 0" type="button" @click="focusComment()"
            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2">
            <ChatBubbleLeftEllipsisIcon class="w-12 h-12 mx-auto text-neutral" />
            <span v-can="'DocStoreService.PostDocumentComment'" class="block mt-2 text-sm font-semibold text-gray-300">No
                comments have been posted yet</span>
        </button>
        <div v-else class="flow-root px-4 rounded-lg text-neutral">
            <ul role="list" class="divide-y divide-gray-200">
                <DocumentCommentEntry v-for="com in comments" :key="com.getId()" :comment="com"
                    @removed="(c: DocumentComment) => removeComment(c)" />
            </ul>
        </div>
    </div>
    <TablePagination v-if="comments.length > 0" :pagination="pagination" @offset-change="offset = $event" class="mt-2" />
</template>
