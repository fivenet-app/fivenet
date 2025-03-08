<script setup lang="ts">
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useNotificatorStore } from '~/store/notificator';
import type { Comment } from '~~/gen/ts/resources/documents/comment';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { GetCommentsResponse } from '~~/gen/ts/services/docstore/docstore';
import DocumentCommentEntry from './DocumentCommentEntry.vue';

const props = withDefaults(
    defineProps<{
        documentId: number;
        closed?: boolean;
        canComment?: boolean;
    }>(),
    {
        closed: false,
        canComment: false,
    },
);

const emit = defineEmits<{
    (e: 'counted', count: number): void;
    (e: 'newComment'): void;
    (e: 'deletedComment'): void;
}>();

const { $grpc } = useNuxtApp();

const { can } = useAuth();

const notifications = useNotificatorStore();

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-comments-${page.value}`, () => getComments(), {
    immediate: false,
});

async function getComments(): Promise<GetCommentsResponse> {
    try {
        const call = $grpc.docstore.docStore.getComments({
            pagination: {
                offset: offset.value,
                pageSize: 5,
            },
            documentId: props.documentId,
        });
        const { response } = await call;

        if (response.pagination) {
            emit('counted', response.pagination.totalCount);
        }

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const schema = z.object({
    comment: z.string().min(3).max(1536),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    comment: '',
});

async function addComment(documentId: number, values: Schema): Promise<void> {
    if (data.value === null) {
        return;
    }

    const comment: Comment = {
        id: 0,
        documentId,
        content: {
            rawContent: values.comment,
        },
        creatorJob: '',
    };

    try {
        const call = $grpc.docstore.docStore.postComment({ comment });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        if (response.comment) {
            data.value?.comments.unshift(response.comment);
        }

        state.comment = '';

        emit('newComment');
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function removeComment(id: number): Promise<void> {
    if (!data.value) {
        return;
    }

    const idx = data.value.comments.findIndex((c) => {
        return c.id === id;
    });

    if (idx > -1) {
        data.value.comments.splice(idx, 1);
    }

    emit('deletedComment');
}

const commentsEl = useTemplateRef('commentsEl');
const isVisible = useElementVisibility(commentsEl);

watchOnce(isVisible, async () => refresh());

watch(offset, async () => refresh());

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await addComment(props.documentId, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div>
        <div ref="commentsEl">
            <template v-if="can('DocStoreService.PostComment').value">
                <div v-if="!closed && canComment" class="flex items-start space-x-4">
                    <div class="min-w-0 flex-1">
                        <UForm :schema="schema" :state="state" class="relative" @submit="onSubmitThrottle">
                            <UFormGroup name="comment">
                                <ClientOnly>
                                    <TiptapEditor
                                        v-model="state.comment"
                                        :placeholder="$t('components.documents.document_comments.add_comment')"
                                        wrapper-class="min-h-44"
                                        comment-mode
                                        :limit="1250"
                                    />
                                </ClientOnly>
                            </UFormGroup>

                            <div class="mt-2 shrink-0">
                                <UButton type="submit" :disabled="!canSubmit" :loading="!canSubmit">
                                    {{ $t('common.post') }}
                                </UButton>
                            </div>
                        </UForm>
                    </div>
                </div>
            </template>
        </div>

        <div class="mt-2">
            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.comment', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.comment', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="!data?.comments || data?.comments.length === 0"
                :message="$t('components.documents.document_comments.no_comments')"
                icon="i-mdi-comment-text-multiple"
            />

            <ul v-else role="list" class="divide-y divide-gray-100 dark:divide-gray-800">
                <DocumentCommentEntry
                    v-for="(comment, idx) in data.comments"
                    :key="comment.id"
                    v-model="data.comments[idx]"
                    @deleted="removeComment(comment.id)"
                />
            </ul>

            <Pagination
                v-if="data?.pagination?.totalCount && data?.pagination?.totalCount > 0"
                v-model="page"
                :pagination="data?.pagination"
                :loading="loading"
                :refresh="refresh"
                disable-border
            />
        </div>
    </div>
</template>
