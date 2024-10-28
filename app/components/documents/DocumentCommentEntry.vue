<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useNotificatorStore } from '~/store/notificator';
import type { Comment } from '~~/gen/ts/resources/documents/comment';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    modelValue?: Comment;
}>();

const emits = defineEmits<{
    (e: 'update:modelValue', comment: Comment | undefined): void;
    (e: 'deleted', id: string | undefined): void;
}>();

const comment = useVModel(props, 'modelValue', emits);

const modal = useModal();

const { can, activeChar, isSuperuser } = useAuth();

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
    try {
        const { response } = await getGRPCDocStoreClient().editComment({
            comment: {
                id: commentId,
                documentId,
                comment: values.comment,
                creatorJob: '',
            },
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

async function deleteComment(id: string): Promise<void> {
    try {
        await getGRPCDocStoreClient().deleteComment({
            commentId: id,
        });

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emits('deleted', comment.value?.id);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function resetForm(): void {
    if (!comment.value) {
        return;
    }

    state.comment = comment.value.comment;
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
                        <UIcon name="i-mdi-trash-can" class="size-5 shrink-0" />
                        <span>{{ $t('common.deleted') }}</span>
                    </div>

                    <div v-if="comment.creatorId === activeChar?.userId || isSuperuser">
                        <UButton
                            v-if="can('DocStoreService.PostComment').value"
                            variant="link"
                            icon="i-mdi-pencil"
                            @click="editing = true"
                        />

                        <UButton
                            v-if="can('DocStoreService.DeleteComment').value"
                            variant="link"
                            icon="i-mdi-trash-can"
                            color="red"
                            @click="
                                modal.open(ConfirmModal, {
                                    confirm: async () => deleteComment(comment!.id),
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
            <div v-if="can('DocStoreService.PostComment').value" class="flex items-start space-x-4">
                <div class="min-w-0 flex-1">
                    <UForm :schema="schema" :state="state" class="relative" @submit="onSubmitThrottle">
                        <UFormGroup name="comment">
                            <UTextarea
                                ref="commentInput"
                                v-model="state.comment"
                                :rows="5"
                                :placeholder="$t('components.documents.document_comments.add_comment')"
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
