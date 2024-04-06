<script setup lang="ts">
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useAuthStore } from '~/store/auth';
import { type Comment } from '~~/gen/ts/resources/documents/comment';
import { GetCommentsResponse } from '~~/gen/ts/services/docstore/docstore';
import DocumentCommentEntry from '~/components/documents/DocumentCommentEntry.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();

const { activeChar } = storeToRefs(authStore);

const props = withDefaults(
    defineProps<{
        documentId: string;
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

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending, refresh, error } = useLazyAsyncData(
    `document-${props.documentId}-comments-${page.value}`,
    () => getComments(),
    {
        immediate: false,
    },
);

async function getComments(): Promise<GetCommentsResponse> {
    try {
        const call = $grpc.getDocStoreClient().getComments({
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
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

interface FormData {
    comment: string;
}

async function addComment(documentId: string, values: FormData): Promise<void> {
    if (data.value === null) {
        return;
    }

    const comment: Comment = {
        id: '0',
        documentId,
        comment: values.comment,
        creatorJob: '',
    };

    try {
        const call = $grpc.getDocStoreClient().postComment({ comment });
        const { response } = await call;

        comment.id = response.id;
        comment.creator = activeChar.value!;

        data.value.comments.unshift(comment);

        resetForm();

        emit('newComment');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function removeComment(comment: Comment): Promise<void> {
    if (!data.value) {
        return;
    }

    if (!can('DocStoreService.DeleteComment')) {
        return;
    }

    const idx = data.value.comments.findIndex((c) => {
        return c.id === comment.id;
    });

    if (idx > -1) {
        data.value.comments.splice(idx, 1);
    }

    emit('deletedComment');
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

const { handleSubmit, meta, resetForm } = useForm<FormData>({
    validationSchema: {
        comment: { required: true, min: 3, max: 1536 },
    },
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await addComment(props.documentId, values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <UContainer>
        <div ref="commentsEl">
            <template v-if="can('DocStoreService.PostComment')">
                <div v-if="!closed && canComment" class="flex items-start space-x-4">
                    <div class="min-w-0 flex-1">
                        <UForm :state="{}" class="relative">
                            <div>
                                <VeeField
                                    v-slot="{ handleChange, value }"
                                    ref="commentInput"
                                    rows="3"
                                    name="comment"
                                    :label="$t('common.comment')"
                                    class="block w-full resize-none border-0 bg-transparent text-gray-50 placeholder:text-gray-400 focus:ring-0 sm:py-1.5 sm:text-sm sm:leading-6"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                >
                                    <UTextarea
                                        :placeholder="$t('components.documents.document_comments.add_comment')"
                                        :model-value="value"
                                        @change="handleChange"
                                    />
                                </VeeField>
                                <VeeErrorMessage name="comment" as="p" class="mt-2 text-sm text-error-400" />
                            </div>

                            <div class="mt-2 shrink-0">
                                <UButton :disabled="!meta.valid || !canSubmit" :loading="!canSubmit" @click="onSubmitThrottle">
                                    {{ $t('common.post') }}
                                </UButton>
                            </div>
                        </UForm>
                    </div>
                </div>
            </template>
        </div>

        <div>
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.comment', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.comment', 2)])"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-if="!data || !data.comments || data?.comments.length === 0"
                :message="$t('components.documents.document_comments.no_comments')"
                icon="i-mdi-comment-text-multiple"
                :focus="focusComment"
            />
            <ul
                v-if="data && data.comments && data.comments.length > 0"
                role="list"
                class="divide-y divide-gray-200 dark:divide-gray-700"
            >
                <DocumentCommentEntry
                    v-for="(comment, idx) in data?.comments"
                    :key="comment.id"
                    v-model:comment="data!.comments[idx]"
                    @deleted="removeComment($event)"
                />
            </ul>

            <div class="flex justify-end px-3 py-3.5">
                <UPagination
                    v-model="page"
                    :page-count="data?.pagination?.pageSize ?? 0"
                    :total="data?.pagination?.totalCount ?? 0"
                />
            </div>
        </div>
    </UContainer>
</template>
