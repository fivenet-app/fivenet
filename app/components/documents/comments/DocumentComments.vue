<script setup lang="ts">
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import Pagination from '~/components/partials/Pagination.vue';
import type { Content } from '~/types/history';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import type { Comment } from '~~/gen/ts/resources/documents/comment';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { GetCommentsResponse } from '~~/gen/ts/services/documents/documents';
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

const notifications = useNotificationsStore();

const historyStore = useHistoryStore();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const page = useRouteQuery('page', '1', { transform: Number });

const { data, status, refresh, error } = useLazyAsyncData(
    `document-${props.documentId}-comments-${page.value}`,
    () => getComments(),
    {
        immediate: false,
    },
);

async function getComments(): Promise<GetCommentsResponse> {
    try {
        const call = documentsDocumentsClient.getComments({
            pagination: {
                offset: calculateOffset(page.value, data.value?.pagination),
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
    content: z.string().min(3).max(1536),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    content: '',
});

const changed = ref(false);
const saving = ref(false);

// Track last saved string and timestamp
let lastSavedString = '';
let lastSaveTimestamp = 0;

async function saveHistory(values: Schema, name: string | undefined = undefined, type = 'document_comments'): Promise<void> {
    if (saving.value) {
        return;
    }

    const now = Date.now();
    // Skip if identical to last saved or if within MIN_GAP
    if (state.content === lastSavedString || now - lastSaveTimestamp < 5000) {
        return;
    }

    saving.value = true;

    historyStore.addVersion<Content>(
        type,
        props.documentId,
        {
            content: values.content,
            files: [],
        },
        name,
    );

    useTimeoutFn(() => {
        saving.value = false;
    }, 1750);

    lastSavedString = state.content;
    lastSaveTimestamp = now;
}

historyStore.handleRefresh(() => saveHistory(state, 'document'));

watchDebounced(
    state,
    () => {
        if (changed.value) {
            saveHistory(state);
        } else {
            changed.value = true;
        }
    },
    {
        debounce: 1_000,
        maxWait: 2_500,
    },
);

async function addComment(documentId: number, values: Schema): Promise<void> {
    if (data.value === null) {
        return;
    }

    const comment: Comment = {
        id: 0,
        documentId,
        content: {
            rawContent: values.content,
        },
        creatorJob: '',
    };

    try {
        const call = documentsDocumentsClient.postComment({ comment });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        if (response.comment) {
            data.value?.comments.unshift(response.comment);
        }

        state.content = '';

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

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await addComment(props.documentId, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div>
        <div ref="commentsEl">
            <template v-if="canComment">
                <div v-if="!closed && canComment" class="flex items-start space-x-4">
                    <div class="min-w-0 flex-1">
                        <UForm class="relative" :schema="schema" :state="state" @submit="onSubmitThrottle">
                            <UFormField name="comment">
                                <ClientOnly>
                                    <TiptapEditor
                                        v-model="state.content"
                                        :placeholder="$t('components.documents.document_comments.add_comment')"
                                        wrapper-class="min-h-44"
                                        disable-images
                                        :limit="1250"
                                        history-type="document_comments"
                                        :saving="saving"
                                    />
                                </ClientOnly>
                            </UFormField>

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
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.comment', 2)])" />
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

            <ul v-else class="divide-y divide-gray-100 dark:divide-gray-800" role="list">
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
                :status="status"
                :refresh="refresh"
                disable-border
            />
        </div>
    </div>
</template>
