<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import type { JSONContent } from '@tiptap/core';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import CustomContentRenderer from '~/components/partials/content/CustomContentRenderer.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import type { HistoryContent } from '~/types/history';
import { contentToTiptapValue, tiptapToContent } from '~/utils/content';
import { getDocumentsCommentsClient } from '~~/gen/ts/clients';
import type { Comment } from '~~/gen/ts/resources/documents/comment/comment';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { isEmptyDoc } from '~/utils/tiptap';

const props = withDefaults(
    defineProps<{
        documentId: number;
        canComment?: boolean;
    }>(),
    {
        canComment: false,
    },
);

const emit = defineEmits<{
    (e: 'deleted', id: number | undefined): void;
    (e: 'restored', id: number | undefined): void;
}>();

const comment = defineModel<Comment | undefined>();

const { t } = useI18n();

const overlay = useOverlay();

const { can, activeChar, isSuperuser } = useAuth();

const { custom } = useAppConfig();

const notifications = useNotificationsStore();

const historyStore = useHistoryStore();

const documentsCommentsClient = await getDocumentsCommentsClient();

const editing = ref<boolean>(false);

const schema = z.object({
    content: z.custom<JSONContent | string>().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    content: undefined,
});

function serializeCommentContent(value: Schema): string {
    if (isEmptyDoc(value.content as JSONContent | undefined)) return '__empty__';
    return JSON.stringify(value.content);
}

const { snapshotDirty, syncSnapshot } = useSnapshotChanges(state, {
    dirty: editing,
    serializer: serializeCommentContent,
});

const saving = ref<boolean>(false);

// Track last saved string and timestamp
let lastSavedString: JSONContent | string | undefined = undefined;
let lastSaveTimestamp = 0;

async function saveHistory(values: Schema, type = 'document_comments'): Promise<void> {
    if (saving.value) return;

    const now = Date.now();
    // Skip if identical to last saved or if within MIN_GAP
    if (state.content === lastSavedString || now - lastSaveTimestamp < 5000) return;

    saving.value = true;

    historyStore.addVersion<HistoryContent>(
        type,
        props.documentId,
        {
            content: values.content,
            files: [],
        },
        `${t('common.comment')}: DOC-${props.documentId}`,
    );

    useTimeoutFn(() => {
        saving.value = false;
    }, 1750);

    lastSavedString = state.content;
    lastSaveTimestamp = now;
}

historyStore.handleRefresh(() => saveHistory(state));

watchDebounced(
    state,
    () => {
        if (snapshotDirty.value) {
            saveHistory(state);
        }
    },
    {
        debounce: 1_000,
        maxWait: 2_500,
    },
);

async function editComment(documentId: number, commentId: number, values: Schema): Promise<void> {
    try {
        const { response } = await documentsCommentsClient.editComment({
            comment: {
                id: commentId,
                documentId,
                content: tiptapToContent(values.content),
                creatorJob: '',
            },
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        editing.value = false;
        setFromProps();

        if (!response.comment) return;

        comment.value = response.comment;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteComment(id: number): Promise<void> {
    try {
        await documentsCommentsClient.deleteComment({
            commentId: id,
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        if (comment.value) {
            if (comment.value.deletedAt) {
                comment.value.deletedAt = undefined;
                emit('restored', comment.value?.id);
            } else {
                emit('deleted', comment.value?.id);
                comment.value.deletedAt = toTimestamp();
            }
        }
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function setFromProps(): void {
    if (!comment.value) return;

    state.content = contentToTiptapValue(comment.value.content);
    syncSnapshot();
}

onBeforeMount(() => setFromProps());
watch(props, () => setFromProps());

function cancelEdit(): void {
    editing.value = false;
    nextTick(() => setFromProps());
}

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (!comment.value) return;

    canSubmit.value = false;
    await editComment(comment.value.documentId, comment.value.id, event.data).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);

const confirmModal = overlay.create(ConfirmModal);
</script>

<template>
    <div
        v-if="comment"
        class="group relative rounded-md bg-neutral-100 px-2 py-2 text-default ring ring-default dark:bg-neutral-900 dark:ring-neutral-700"
        :class="comment.deletedAt ? custom.classes.deletedAt : ''"
    >
        <div v-if="!editing" class="relative">
            <div
                v-if="comment.creatorId === activeChar?.userId || isSuperuser"
                class="absolute top-2 right-2 z-10 opacity-0 transition-opacity group-focus-within:opacity-100 group-hover:opacity-100"
            >
                <UFieldGroup>
                    <UTooltip v-if="canComment" :text="$t('common.edit')">
                        <UButton variant="link" icon="i-mdi-pencil" @click="editing = true" />
                    </UTooltip>

                    <UTooltip
                        v-if="can('documents.CommentsService/DeleteComment').value"
                        :text="!comment.deletedAt ? $t('common.delete') : $t('common.restore')"
                    >
                        <UButton
                            :color="!comment.deletedAt ? 'error' : 'success'"
                            :icon="!comment.deletedAt ? 'i-mdi-delete' : 'i-mdi-restore'"
                            variant="link"
                            @click="
                                () =>
                                    confirmModal.open({
                                        confirm: async () => deleteComment(comment!.id),
                                    })
                            "
                        />
                    </UTooltip>
                </UFieldGroup>
            </div>

            <div class="rounded-lg p-2">
                <CustomContentRenderer v-if="comment.content" :value="comment.content" />
            </div>
        </div>

        <div v-else-if="canComment" class="flex items-start space-x-4">
            <div class="min-w-0 flex-1">
                <UForm class="relative" :schema="schema" :state="state" @submit="onSubmitThrottle">
                    <UFormField name="content" :ui="{ error: 'hidden' }">
                        <ClientOnly>
                            <TiptapEditor
                                v-model="state.content"
                                name="content"
                                wrapper-class="min-h-44"
                                disable-images
                                :limit="1250"
                                :saving="saving"
                                history-type="document_comments"
                            />
                        </ClientOnly>
                    </UFormField>

                    <div class="mt-2 flex shrink-0 justify-between">
                        <UButton
                            type="submit"
                            :disabled="!canSubmit"
                            :label="$t('common.edit')"
                            :loading="!canSubmit"
                            trailing-icon="i-mdi-comment-edit"
                        />

                        <UButton
                            type="button"
                            color="error"
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                            :label="$t('common.cancel')"
                            @click="cancelEdit"
                        />
                    </div>
                </UForm>
            </div>
        </div>
    </div>
</template>
