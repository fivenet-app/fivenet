<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useAuthStore } from '~/store/auth';
import { Comment } from '~~/gen/ts/resources/documents/comment';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import ConfirmModal from '../partials/ConfirmModal.vue';

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

const editing = ref(false);

interface FormData {
    comment: string;
}

async function editComment(documentId: string, commentId: string, values: FormData): Promise<void> {
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

        emit('deleted', props.comment);
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta, setValues } = useForm<FormData>({
    validationSchema: {
        comment: { required: true, min: 3, max: 1536 },
    },
    validateOnMount: true,
});

function resetForm(): void {
    setValues({
        comment: props.comment.comment,
    });
}

onMounted(() => resetForm());

watch(props, () => resetForm());

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await editComment(props.comment.documentId, props.comment.id, values).finally(() =>
            useTimeoutFn(() => (canSubmit.value = true), 400),
        ),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <li class="py-2">
        <div v-if="!editing" class="flex space-x-3">
            <div :class="[comment.deletedAt ? 'bg-warn-800' : '', 'flex-1 space-y-1']">
                <div class="flex items-center justify-between">
                    <div class="flex items-center">
                        <CitizenInfoPopover
                            :user="comment.creator"
                            class="text-sm font-medium text-primary-400 hover:text-primary-300"
                        />
                    </div>
                    <div class="flex flex-1 items-center text-accent-200">
                        <GenericTime class="ml-2 text-sm" :value="comment.createdAt" />
                    </div>
                    <div v-if="comment.deletedAt" class="flex flex-1 flex-row items-center justify-center">
                        <UIcon name="i-mdi-trash-can" class="mr-1.5 size-5 shrink-0" />
                        {{ $t('common.deleted') }}
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
                    <form class="relative" @submit.prevent="onSubmitThrottle">
                        <div
                            class="overflow-hidden rounded-lg shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-primary-600"
                        >
                            <label for="comment" class="sr-only">
                                {{ $t('components.documents.document_comment_entry.edit_comment') }}
                            </label>
                            <VeeField
                                ref="commentInput"
                                as="textarea"
                                rows="3"
                                name="comment"
                                :label="$t('common.comment')"
                                :placeholder="$t('components.documents.document_comment_entry.edit_comment')"
                                class="block w-full resize-none border-0 bg-transparent text-gray-50 placeholder:text-gray-400 focus:ring-0 sm:py-1.5 sm:text-sm sm:leading-6"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />

                            <!-- Spacer element to match the height of the toolbar -->
                            <div class="py-2">
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
                            <div class="shrink-0">
                                <UButton
                                    type="submit"
                                    class="flex justify-center rounded-md px-3 py-2 text-sm font-semibold shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                    :disabled="!meta.valid || !canSubmit"
                                    :loading="!canSubmit"
                                >
                                    {{ $t('common.edit') }}
                                </UButton>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </template>
    </li>
</template>
