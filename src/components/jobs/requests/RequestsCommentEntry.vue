<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { max, min, required } from '@vee-validate/rules';
import { useConfirmDialog, useThrottleFn } from '@vueuse/core';
import { LoadingIcon, PencilIcon, TrashCanIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useAuthStore } from '~/store/auth';
import { RequestComment } from '~~/gen/ts/resources/jobs/requests';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();

const { activeChar, permissions } = storeToRefs(authStore);

const emit = defineEmits<{
    (e: 'update:comment', comment: RequestComment): void;
    (e: 'deleted', comment: RequestComment): void;
}>();

const props = defineProps<{
    comment: RequestComment;
}>();

const editing = ref(false);

interface FormData {
    comment: string;
}

async function editComment(requestId: string, commentId: string, values: FormData): Promise<void> {
    const comment: RequestComment = {
        id: commentId,
        requestId,
        comment: values.comment,
    };

    try {
        const { response } = await $grpc.getJobsRequestsClient().postRequestComment({ comment });

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
        await $grpc.getJobsRequestsClient().deleteRequestComment({
            id,
        });

        emit('deleted', props.comment);
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();

onConfirm(async (id) => deleteComment(id));

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta, resetForm } = useForm<FormData>({
    validationSchema: {
        comment: { required: true, min: 3, max: 1536 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await editComment(props.comment.requestId, props.comment.id, values).finally(() =>
            setTimeout(() => (canSubmit.value = true), 400),
        ),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(props.comment.id)" />

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
                    <div v-if="comment.deletedAt" class="flex flex-1 flex-row items-center justify-center text-base-100">
                        <TrashCanIcon type="button" class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-300" aria-hidden="true" />
                        {{ $t('common.deleted') }}
                    </div>
                    <div v-if="comment.creatorId === activeChar?.userId || permissions.includes('superuser')">
                        <button v-if="can('JobsRequestsService.CreateRequest')" @click="editing = true">
                            <PencilIcon class="ml-auto mr-2.5 w-5 h-auto" aria-hidden="true" />
                        </button>
                        <button v-if="can('JobsRequestsService.DeleteRequestComment')" type="button" @click="reveal()">
                            <TrashCanIcon class="ml-auto mr-2.5 w-5 h-auto" aria-hidden="true" />
                        </button>
                    </div>
                </div>
                <p class="break-words text-sm">
                    {{ comment.comment }}
                </p>
            </div>
        </div>
        <template v-else>
            <div v-if="can('JobsRequestsService.CreateRequest')" class="flex items-start space-x-4">
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
                                        <LoadingIcon class="mr-2 h-5 w-5 animate-spin" aria-hidden="true" />
                                    </template>
                                    {{ $t('common.edit') }}
                                </button>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </template>
    </li>
</template>
