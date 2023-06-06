<script setup lang="ts">
import { ChatBubbleLeftEllipsisIcon } from '@heroicons/vue/20/solid';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import TablePagination from '~/components/partials/TablePagination.vue';
import { useAuthStore } from '~/store/auth';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import { DocumentComment } from '~~/gen/ts/resources/documents/documents';
import DocumentCommentEntry from './DocumentCommentEntry.vue';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();

const { activeChar } = storeToRefs(authStore);

const props = withDefaults(
    defineProps<{
        documentId: bigint;
        closed?: boolean;
    }>(),
    {
        closed: false,
    }
);

const emit = defineEmits<{
    (e: 'counted', count: bigint): void;
}>();

const pagination = ref<PaginationResponse>();
const offset = ref(0n);

const {
    data: comments,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-comments-${offset}`, () => getDocumentComments());

async function getDocumentComments(): Promise<Array<DocumentComment>> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getDocStoreClient().getDocumentComments({
                pagination: {
                    offset: offset.value,
                    pageSize: 5n,
                },
                documentId: props.documentId,
            });
            const { response } = await call;

            pagination.value = response.pagination;
            if (pagination.value) {
                emit('counted', pagination.value?.totalCount);
            }

            return res(response.comments);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const message = ref('');

async function addComment(): Promise<void> {
    return new Promise(async (res, rej) => {
        if (!comments.value) {
            return res();
        }

        const comment: DocumentComment = {
            id: 0n,
            documentId: props.documentId,
            comment: message.value,
        };

        try {
            const call = $grpc.getDocStoreClient().postDocumentComment({
                comment: comment,
            });
            const { response } = await call;

            comment.id = response.id;
            comment.creator = activeChar.value!;

            comments.value.unshift(comment);

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function removeComment(comment: DocumentComment): Promise<void> {
    return new Promise(async (res, rej) => {
        if (!comments.value) {
            return res();
        }

        const idx = comments.value.findIndex((c) => {
            return c.id === comment.id;
        });

        if (idx > -1) {
            comments.value.splice(idx, 1);
        }

        return res();
    });
}

const commentInput = ref<HTMLInputElement | null>(null);

function focusComment(): void {
    if (commentInput.value) {
        commentInput.value.focus();
    }
}

watch(offset, async () => refresh());
</script>

<template>
    <div class="pb-2">
        <div v-if="!closed" v-can="'DocStoreService.PostDocumentComment'" class="flex items-start space-x-4">
            <div class="min-w-0 flex-1">
                <form @submit.prevent="addComment()" class="relative">
                    <div
                        class="overflow-hidden rounded-lg shadow-sm ring-1 ring-inset ring-gray-500 focus-within:ring-2 focus-within:ring-indigo-600"
                    >
                        <label for="comment" class="sr-only">
                            {{ $t('components.documents.document_comments.add_comment') }}
                        </label>
                        <textarea
                            rows="3"
                            name="comment"
                            class="block w-full resize-none border-0 bg-transparent text-gray-50 placeholder:text-gray-400 focus:ring-0 sm:py-1.5 sm:text-sm sm:leading-6"
                            ref="commentInput"
                            v-model="message"
                            :placeholder="$t('components.documents.document_comments.add_comment')"
                        />

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
                            <button
                                type="submit"
                                class="inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                            >
                                {{ $t('common.post') }}
                            </button>
                        </div>
                    </div>
                </form>
            </div>
        </div>
    </div>
    <div class="bg-base-800">
        <button
            v-if="!comments || comments.length === 0"
            type="button"
            @click="focusComment()"
            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
        >
            <ChatBubbleLeftEllipsisIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-gray-300">
                {{ $t('components.documents.document_comments.no_comments') }}
            </span>
        </button>
        <div v-else v-can="'DocStoreService.DeleteDocumentComment'" class="flow-root px-4 rounded-lg text-neutral">
            <ul role="list" class="divide-y divide-gray-200">
                <DocumentCommentEntry
                    v-for="com in comments"
                    :key="com.id?.toString()"
                    :comment="com"
                    @removed="(c: DocumentComment) => removeComment(c)"
                />
            </ul>
        </div>
    </div>
    <TablePagination
        v-if="comments && comments.length > 0"
        :pagination="pagination"
        @offset-change="offset = $event"
        class="mt-2"
    />
</template>
