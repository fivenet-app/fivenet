<script setup lang="ts">
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { max, min, required } from '@vee-validate/rules';
import { useElementVisibility, useThrottleFn, watchOnce } from '@vueuse/core';
import { CommentTextMultipleIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import { useAuthStore } from '~/store/auth';
import { RequestComment } from '~~/gen/ts/resources/jobs/requests';
import { RequestsListCommentsResponse } from '~~/gen/ts/services/jobs/jobs';
import CommentEntry from './CommentEntry.vue';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();

const { activeChar } = storeToRefs(authStore);

const props = defineProps<{
    requestId: bigint;
}>();

const emit = defineEmits<{
    (e: 'counted', count: bigint): void;
}>();

const offset = ref(0n);

const {
    data: comments,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`jobs-requests-${props.requestId}-comments-${offset}`, () => getComments(), {
    immediate: false,
});

async function getComments(): Promise<RequestsListCommentsResponse> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getJobsClient().requestsListComments({
                pagination: {
                    offset: offset.value,
                    pageSize: 5n,
                },
                requestId: props.requestId,
            });
            const { response } = await call;

            return res(response);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function addComment(documentId: bigint, values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        if (!comments.value) {
            return res();
        }

        const comment: RequestComment = {
            id: 0n,
            requestId: documentId,
            comment: values.comment,
        };

        try {
            const call = $grpc.getJobsClient().requestsPostComment({
                comment: comment,
            });
            const { response } = await call;
            if (!response.comment) return res();

            comment.id = response.comment.id;
            comment.creator = activeChar.value!;

            comments.value.comments.unshift(comment);

            resetForm();

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function removeComment(comment: RequestComment): Promise<void> {
    return new Promise(async (res, rej) => {
        if (!comments.value) {
            return res();
        }

        if (!can('DocStoreService.DeleteComment')) {
            return res();
        }

        const idx = comments.value.comments.findIndex((c) => {
            return c.id === comment.id;
        });

        if (idx > -1) {
            comments.value.comments.splice(idx, 1);
        }

        return res();
    });
}

const commentsEl = ref<HTMLDivElement | null>(null);
const isVisible = useElementVisibility(commentsEl);
watchOnce(isVisible, () => refresh());

const commentInput = ref<HTMLInputElement | null>(null);

function focusComment(): void {
    if (commentInput.value) {
        commentInput.value.focus();
    }
}

watch(offset, async () => refresh());

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

interface FormData {
    comment: string;
}

const { handleSubmit, meta, resetForm } = useForm<FormData>({
    validationSchema: {
        comment: { required: true, min: 3, max: 1536 },
    },
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await addComment(props.requestId, values).finally(() => setTimeout(() => (canSubmit.value = true), 350)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <div>
        <div ref="commentsEl" class="pb-2">
            <div v-if="can('DocStoreService.PostComment')">
                <div class="flex items-start space-x-4">
                    <div class="min-w-0 flex-1">
                        <form @submit.prevent="onSubmitThrottle" class="relative">
                            <div
                                class="overflow-hidden rounded-lg shadow-sm ring-1 ring-inset ring-gray-500 focus-within:ring-2 focus-within:ring-primary-600"
                            >
                                <label for="comment" class="sr-only">
                                    {{ $t('components.documents.document_comments.add_comment') }}
                                </label>
                                <VeeField
                                    as="textarea"
                                    rows="3"
                                    name="comment"
                                    :label="$t('common.comment')"
                                    :placeholder="$t('components.documents.document_comments.add_comment')"
                                    class="block w-full resize-none border-0 bg-transparent text-gray-50 placeholder:text-gray-400 focus:ring-0 sm:py-1.5 sm:text-sm sm:leading-6"
                                    ref="commentInput"
                                />

                                <!-- Spacer element to match the height of the toolbar -->
                                <div class="py-2" aria-hidden="true">
                                    <!-- Matches height of button in toolbar (1px border + 36px content height) -->
                                    <div class="py-px">
                                        <div class="h-9" />
                                        <div class="ml-2">
                                            <VeeErrorMessage name="comment" as="p" class="mt-2 text-sm text-error-400" />
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div class="absolute inset-x-0 bottom-0 flex justify-between py-2 pl-3 pr-2">
                                <div class="flex items-center space-x-5"></div>
                                <div class="flex-shrink-0">
                                    <button
                                        type="submit"
                                        class="flex justify-center rounded-md px-3 py-2 text-sm font-semibold text-neutral shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                        :disabled="!meta.valid || !canSubmit"
                                        :class="[
                                            !meta.valid || !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                        ]"
                                    >
                                        <template v-if="!canSubmit">
                                            <LoadingIcon class="animate-spin h-5 w-5 mr-2" />
                                        </template>
                                        {{ $t('common.post') }}
                                    </button>
                                </div>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
        <div>
            <DataNoDataBlock
                v-if="!comments || comments.comments.length === 0"
                :message="$t('components.documents.document_comments.no_comments')"
                :icon="CommentTextMultipleIcon"
                :focus="focusComment"
            />
            <div
                v-else
                class="flow-root rounded-lg bg-base-800 text-neutral rounded-lg ring-1 ring-inset ring-gray-500 shadow-sm shadow-sm focus-within:ring-2 focus-within:ring-primary-600"
            >
                <ul role="list" class="divide-y divide-gray-200 px-4">
                    <CommentEntry
                        v-for="com in comments.comments"
                        :key="com.id?.toString()"
                        :comment="com"
                        @removed="(c: RequestComment) => removeComment(c)"
                    />
                </ul>
            </div>
        </div>
        <TablePagination
            v-if="comments && comments.comments.length > 0"
            :pagination="comments.pagination"
            @offset-change="offset = $event"
            class="mt-2"
        />
    </div>
</template>
