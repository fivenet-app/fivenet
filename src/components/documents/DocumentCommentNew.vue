<script lang="ts" setup>
import { DocumentComment } from '@arpanet/gen/resources/documents/documents_pb';
import { PostDocumentCommentRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { computed, ref } from 'vue';
import { getDocStoreClient } from '../../grpc/grpc';
import { useStore } from '../../store/store';

const store = useStore();

const activeChar = computed(() => store.state.auth?.activeChar);

const emit = defineEmits<{
    (e: 'added', comment: DocumentComment): void,
}>();

const props = defineProps({
    documentId: {
        required: true,
        type: Number,
    }
});

const message = ref('');

function addComment() {
    const req = new PostDocumentCommentRequest();
    const com = new DocumentComment();
    com.setDocumentId(props.documentId);
    com.setComment(message.value);
    req.setComment(com);

    getDocStoreClient().
        postDocumentComment(req, null).
        then((resp) => {
            com.setCreatorId(activeChar.value!.getUserId());
            com.setCreator(activeChar.value!);
            com.setId(resp.getId());

            emit('added', com);
        });
}
</script>

<template>
    <div v-can="'DocStoreService.PostDocumentComment'" class="flex items-start space-x-4">
        <div class="min-w-0 flex-1">
            <form @submit.prevent="addComment()" class="relative">
                <div
                    class="overflow-hidden rounded-lg shadow-sm ring-1 ring-inset ring-gray-500 focus-within:ring-2 focus-within:ring-indigo-600">
                    <label for="comment" class="sr-only">Add your comment</label>
                    <textarea rows="3" name="comment" id="comment"
                        class="block w-full resize-none border-0 bg-transparent text-gray-50 placeholder:text-gray-400 focus:ring-0 sm:py-1.5 sm:text-sm sm:leading-6"
                        v-model="message" placeholder="Add your comment..." />

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
                        <button type="submit"
                            class="inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Post</button>
                    </div>
                </div>
            </form>
        </div>
    </div>
</template>
