<script lang="ts" setup>
import { useStore } from '../../store/store';
import { DocumentComment } from '@arpanet/gen/resources/documents/documents_pb';
import { EditDocumentCommentRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { PencilIcon } from '@heroicons/vue/20/solid';
import { computed, ref } from 'vue';
import { getDocStoreClient } from '../../grpc/grpc';

const store = useStore();

const activeCharId = computed(() => store.state.auth?.activeChar?.getUserId());

const props = defineProps({
    comment: {
        required: true,
        type: DocumentComment,
    },
});

const editing = ref(false);
const message = ref(props.comment.getComment());

function editComment() {
    const req = new EditDocumentCommentRequest();
    const c = new DocumentComment();
    c.setId(props.comment.getId());
    c.setDocumentId(props.comment.getDocumentId());
    c.setComment(message.value);
    req.setComment(c);

    getDocStoreClient().
        editDocumentComment(req, null).
        then((resp) => {
            props.comment.setComment(message.value);
            editing.value = false;
        });
}
</script>

<template>
    <li class="py-4">
        <div v-if="!editing" class="flex space-x-3">
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <router-link :to="{ name: 'Citizens: Info', params: { id: comment.getCreatorId() } }"
                        class="text-sm font-medium text-primary-400 hover:text-primary-300">
                        {{ comment.getCreator()?.getFirstname() }} {{ comment.getCreator()?.getLastname() }}
                    </router-link>
                    <button v-can="'DocStoreService.PostDocumentComment'" v-if="comment.getCreatorId() === activeCharId"
                        @click="editing = true">
                        <PencilIcon class="w-5 h-auto ml-auto mr-2.5" />
                    </button>
                </div>
                <p class="text-sm">{{ comment.getComment() }}</p>
            </div>
        </div>
        <div v-else v-can="'DocStoreService.PostDocumentComment'" class="flex items-start space-x-4">
            <div class="min-w-0 flex-1">
                <form @submit.prevent="editComment" class="relative">
                    <div
                        class="overflow-hidden rounded-lg shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-indigo-600">
                        <label for="comment" class="sr-only">Edit your comment</label>
                        <textarea rows="3" name="comment" id="comment"
                            class="block w-full resize-none border-0 bg-transparent text-gray-50 placeholder:text-gray-400 focus:ring-0 sm:py-1.5 sm:text-sm sm:leading-6"
                            v-model="message" placeholder="Edit your comment..." />

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
                                class="inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Edit</button>
                        </div>
                    </div>
                </form>
            </div>
        </div>
    </li>
</template>
