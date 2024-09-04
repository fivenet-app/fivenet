<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { Comment } from '~~/gen/ts/resources/documents/comment';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    modelValue: Comment;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', comment: Comment): void;
    (e: 'deleted', comment: Comment): void;
}>();

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
        const { response } = await getGRPCDocStoreClient().editComment({ comment });

        editing.value = false;
        resetForm();

        if (!response.comment) {
            return;
        }

        emit('update:modelValue', response.comment);
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

        emit('deleted', props.modelValue);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function resetForm(): void {
    state.comment = props.modelValue.comment;
}

onMounted(() => resetForm());

watch(props, () => resetForm());

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await editComment(props.modelValue.documentId, props.modelValue.id, event.data).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);
</script>

<template>
    <li class="py-2">
        <div v-if="!editing" class="flex space-x-3">
            <div :class="[modelValue.deletedAt ? 'bg-warn-800' : '', 'flex-1 space-y-1']">
                <div class="flex items-center justify-between">
                    <div class="flex items-center">
                        <CitizenInfoPopover :user="modelValue.creator" show-avatar-in-name />
                    </div>

                    <div class="flex flex-1 items-center">
                        <GenericTime class="ml-2 text-sm" :value="modelValue.createdAt" />
                    </div>

                    <div v-if="modelValue.deletedAt" class="flex flex-1 flex-row items-center justify-center gap-1.5">
                        <UIcon name="i-mdi-trash-can" class="size-5 shrink-0" />
                        <span>{{ $t('common.deleted') }}</span>
                    </div>

                    <div v-if="modelValue.creatorId === activeChar?.userId || permissions.includes('superuser')">
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
                            @click="
                                modal.open(ConfirmModal, {
                                    confirm: async () => deleteComment(modelValue.id),
                                })
                            "
                        />
                    </div>
                </div>

                <p class="whitespace-pre-line break-words text-sm">
                    {{ modelValue.comment }}
                </p>
            </div>
        </div>

        <template v-else>
            <div v-if="can('DocStoreService.PostComment').value" class="flex items-start space-x-4">
                <div class="min-w-0 flex-1">
                    <UForm :schema="schema" :state="state" class="relative" @submit="onSubmitThrottle">
                        <UFormGroup name="comment">
                            <UTextarea
                                v-model="state.comment"
                                ref="commentInput"
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
