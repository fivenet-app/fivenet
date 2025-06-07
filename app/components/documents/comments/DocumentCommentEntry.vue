<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import HTMLContent from '~/components/partials/content/HTMLContent.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useNotificatorStore } from '~/stores/notificator';
import type { Content } from '~/types/history';
import type { Comment } from '~~/gen/ts/resources/documents/comment';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = withDefaults(
    defineProps<{
        modelValue?: Comment;
        canComment?: boolean;
    }>(),
    {
        modelValue: undefined,
        canComment: false,
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', comment: Comment | undefined): void;
    (e: 'deleted', id: number | undefined): void;
}>();

const comment = useVModel(props, 'modelValue', emit);

const { $grpc } = useNuxtApp();

const modal = useModal();

const { can, activeChar, isSuperuser } = useAuth();

const notifications = useNotificatorStore();

const historyStore = useHistoryStore();

const editing = ref(false);

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
        props.modelValue!.documentId,
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

async function editComment(documentId: number, commentId: number, values: Schema): Promise<void> {
    try {
        const { response } = await $grpc.documents.documents.editComment({
            comment: {
                id: commentId,
                documentId,
                content: {
                    rawContent: values.content,
                },
                creatorJob: '',
            },
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        editing.value = false;
        resetForm();

        if (!response.comment) {
            return;
        }

        comment.value = response.comment;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteComment(id: number): Promise<void> {
    try {
        await $grpc.documents.documents.deleteComment({
            commentId: id,
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('deleted', comment.value?.id);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function resetForm(): void {
    if (!comment.value) {
        return;
    }

    state.content = comment.value.content?.rawContent ?? '';
}

onMounted(() => resetForm());
watch(props, () => resetForm());

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (!comment.value) {
        return;
    }

    canSubmit.value = false;
    await editComment(comment.value.documentId, comment.value.id, event.data).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);
</script>

<template>
    <li v-if="comment" class="py-2">
        <div v-if="!editing" class="flex space-x-3">
            <div :class="[comment.deletedAt ? 'bg-warn-800' : '', 'flex-1 space-y-1']">
                <div class="flex items-center justify-between">
                    <div class="flex items-center">
                        <CitizenInfoPopover :user="comment.creator" show-avatar-in-name />
                    </div>

                    <div class="flex flex-1 items-center">
                        <GenericTime class="ml-2 text-sm" :value="comment.createdAt" />
                    </div>

                    <div v-if="comment.deletedAt" class="flex flex-1 flex-row items-center justify-center gap-1.5">
                        <UIcon class="size-5 shrink-0" name="i-mdi-delete" />
                        <span>{{ $t('common.deleted') }}</span>
                    </div>

                    <div v-if="comment.creatorId === activeChar?.userId || isSuperuser">
                        <UTooltip v-if="canComment" :text="$t('common.edit')">
                            <UButton v-if="canComment" variant="link" icon="i-mdi-pencil" @click="editing = true" />
                        </UTooltip>

                        <UTooltip v-if="can('documents.DocumentsService.DeleteComment').value" :text="$t('common.delete')">
                            <UButton
                                variant="link"
                                icon="i-mdi-delete"
                                color="error"
                                @click="
                                    modal.open(ConfirmModal, {
                                        confirm: async () => deleteComment(comment!.id),
                                    })
                                "
                            />
                        </UTooltip>
                    </div>
                </div>

                <div class="rounded-lg bg-neutral-100 dark:bg-base-900">
                    <HTMLContent v-if="comment.content?.content" class="px-4 py-2" :value="comment.content.content" />
                </div>
            </div>
        </div>

        <div v-else-if="canComment" class="flex items-start space-x-4">
            <div class="min-w-0 flex-1">
                <UForm class="relative" :schema="schema" :state="state" @submit="onSubmitThrottle">
                    <UFormGroup name="comment">
                        <ClientOnly>
                            <TiptapEditor
                                v-model="state.content"
                                wrapper-class="min-h-44"
                                disable-images
                                :limit="1250"
                                disable-collab
                                :saving="saving"
                                history-type="document_comments"
                            />
                        </ClientOnly>
                    </UFormGroup>

                    <div class="mt-2 shrink-0">
                        <UButton type="submit" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.edit') }}
                        </UButton>
                    </div>
                </UForm>
            </div>
        </div>
    </li>
</template>
