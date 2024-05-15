<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useAuthStore } from '~/store/auth';
import { Comment } from '~~/gen/ts/resources/documents/comment';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    comment: Comment;
}>();

const emit = defineEmits<{
    (e: 'update:comment', comment: Comment): void;
    (e: 'deleted', comment: Comment): void;
}>();

const { $grpc } = useNuxtApp();

const modal = useModal();

const authStore = useAuthStore();
const { activeChar, permissions } = storeToRefs(authStore);

const notifications = useNotificatorStore();

const editing = ref(false);

const schema = z.object({
    comment: z.string().min(3).max(1536),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    comment: '',
});

async function editComment(documentId: string, commentId: string, values: Schema): Promise<void> {
    const comment: Comment = {
        id: commentId,
        documentId,
        comment: values.comment,
        creatorJob: '',
    };

    try {
        const { response } = await $grpc.getDocStoreClient().editComment({ comment });

        editing.value = false;
        resetForm();

        emit('update:comment', response.comment!);
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function deleteComment(id: string): Promise<void> {
    try {
        await $grpc.getDocStoreClient().deleteComment({
            commentId: id,
        });

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('deleted', props.comment);
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function resetForm(): void {
    state.comment = props.comment.comment;
}

onMounted(() => resetForm());

watch(props, () => resetForm());

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await editComment(props.comment.documentId, props.comment.id, event.data).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);
</script>

<template>
    <li class="py-2">
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
                        <UIcon name="i-mdi-trash-can" class="size-5 shrink-0" />
                        <span>{{ $t('common.deleted') }}</span>
                    </div>

                    <div v-if="comment.creatorId === activeChar?.userId || permissions.includes('superuser')">
                        <UButton
                            v-if="can('DocStoreService.PostComment')"
                            variant="link"
                            icon="i-mdi-pencil"
                            @click="editing = true"
                        />

                        <UButton
                            v-if="can('DocStoreService.DeleteComment')"
                            variant="link"
                            icon="i-mdi-trash-can"
                            @click="
                                modal.open(ConfirmModal, {
                                    confirm: async () => deleteComment(comment.id),
                                })
                            "
                        />
                    </div>
                </div>

                <p class="whitespace-pre-line break-words text-sm">
                    {{ comment.comment }}
                </p>
            </div>
        </div>

        <template v-else>
            <div v-if="can('DocStoreService.PostComment')" class="flex items-start space-x-4">
                <div class="min-w-0 flex-1">
                    <UForm :schema="schema" :state="state" class="relative" @submit="onSubmitThrottle">
                        <UFormGroup name="comment">
                            <UTextarea
                                v-model="state.comment"
                                ref="commentInput"
                                :rows="5"
                                :placeholder="$t('components.documents.document_comments.add_comment')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <div class="mt-2 shrink-0">
                            <UButton type="submit" :disabled="!canSubmit" :loading="!canSubmit">
                                {{ $t('common.edit') }}
                            </UButton>
                        </div>
                    </UForm>
                </div>
            </div>
        </template>
    </li>
</template>
