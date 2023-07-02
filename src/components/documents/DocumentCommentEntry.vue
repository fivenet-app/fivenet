<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiPencil, mdiTrashCan } from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { useAuthStore } from '~/store/auth';
import { DocumentComment } from '~~/gen/ts/resources/documents/documents';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();

const { activeChar, permissions } = storeToRefs(authStore);

const emit = defineEmits<{
    (e: 'removed', comment: DocumentComment): void;
}>();

const props = defineProps<{
    comment: DocumentComment;
}>();

const editing = ref(false);
const message = ref(props.comment.comment);

async function editComment(): Promise<void> {
    return new Promise(async (res, rej) => {
        const comment: DocumentComment = {
            id: props.comment.id,
            documentId: props.comment.documentId,
            comment: message.value,
        };

        try {
            await $grpc.getDocStoreClient().editDocumentComment({
                comment: comment,
            });

            editing.value = false;
            props.comment.comment = message.value;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function deleteComment(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            await $grpc.getDocStoreClient().deleteDocumentComment({
                commentId: props.comment.id,
            });

            emit('removed', props.comment);

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <li class="py-2">
        <div v-if="!editing" class="flex space-x-3">
            <div :class="[comment.deletedAt ? 'bg-warn-800' : 'bg-base-850', 'flex-1 space-y-1']">
                <div class="flex items-center justify-between">
                    <NuxtLink
                        :to="{ name: 'citizens-id', params: { id: comment.creatorId! } }"
                        class="text-sm font-medium text-primary-400 hover:text-primary-300"
                    >
                        {{ comment.creator?.firstname }}
                        {{ comment.creator?.lastname }}
                    </NuxtLink>
                    <div v-if="comment.deletedAt" class="flex flex-row items-center justify-center flex-1 text-base-100">
                        <SvgIcon
                            class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                            aria-hidden="true"
                            type="mdi"
                            :path="mdiTrashCan"
                        />
                        {{ $t('common.deleted') }}
                    </div>
                    <div v-if="comment.creatorId === activeChar?.userId || permissions.includes('superuser')">
                        <button v-can="'DocStoreService.PostDocumentComment'" @click="editing = true">
                            <SvgIcon class="w-5 h-auto ml-auto mr-2.5" type="mdi" :path="mdiPencil" />
                        </button>
                        <button v-can="'DocStoreService.DeleteDocumentComment'" @click="deleteComment()">
                            <SvgIcon class="w-5 h-auto ml-auto mr-2.5" type="mdi" :path="mdiTrashCan" />
                        </button>
                    </div>
                </div>
                <p class="text-sm break-words">
                    {{ comment.comment }}
                </p>
            </div>
        </div>
        <div v-else v-can="'DocStoreService.PostDocumentComment'" class="flex items-start space-x-4">
            <div class="min-w-0 flex-1">
                <form @submit.prevent="editComment" class="relative">
                    <div
                        class="overflow-hidden rounded-lg shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-primary-600"
                    >
                        <label for="comment" class="sr-only">
                            {{ $t('components.documents.document_comment_entry.edit_comment') }}
                        </label>
                        <textarea
                            rows="3"
                            name="comment"
                            class="block w-full resize-none border-0 bg-transparent text-gray-50 placeholder:text-gray-400 focus:ring-0 sm:py-1.5 sm:text-sm sm:leading-6"
                            v-model="message"
                            :placeholder="$t('components.documents.document_comment_entry.edit_comment')"
                        />

                        <!-- Spacer element to match the height of the toolbar -->
                        <div class="py-2" aria-hidden="true">
                            <!-- Matches height of button in toolbar (1px border + 36px content height) -->
                            <div class="py-px">
                                <div class="h-9" />
                            </div>
                        </div>
                    </div>

                    <div class="absolute inset-x-0 bottom-0 flex justify-between py-2 pl-3 pr-2">
                        <div class="flex items-center space-x-5"></div>
                        <div class="flex-shrink-0">
                            <button
                                type="submit"
                                class="inline-flex items-center rounded-md bg-primary-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600"
                            >
                                {{ $t('common.edit') }}
                            </button>
                        </div>
                    </div>
                </form>
            </div>
        </div>
    </li>
</template>
