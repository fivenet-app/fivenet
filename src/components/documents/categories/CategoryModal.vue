<script lang="ts" setup>
import { DocumentCategory } from '@fivenet/gen/resources/documents/category_pb';
import { DeleteDocumentCategoryRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { TagIcon } from '@heroicons/vue/24/solid';
import { dispatchNotification } from '../../notification';
import { RpcError } from 'grpc-web';

const { $grpc } = useNuxtApp();

const emit = defineEmits<{
    (e: 'close'): void,
}>();

const props = defineProps({
    open: {
        required: true,
        type: Boolean,
    },
    category: {
        type: DocumentCategory,
        required: false,
    },
});

async function deleteCategory(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new DeleteDocumentCategoryRequest();
        req.setIdsList([]);

        try {
            const resp = await $grpc.getDocStoreClient()
                .deleteDocumentCategory(req, null);

            dispatchNotification({ title: 'Password has been changed', content: 'Please login with your new password.', type: 'success' });
            emit('close');
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
            <TransitionChild as="template" enter="ease-out duration-300" enter-from="opacity-0" enter-to="opacity-100"
                leave="ease-in duration-200" leave-from="opacity-100" leave-to="opacity-0">
                <div class="fixed inset-0 transition-opacity bg-opacity-75 bg-base-900" />
            </TransitionChild>

            <div class="fixed inset-0 z-10 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild as="template" enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100" leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95">
                        <div>
                            <DialogPanel
                                class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-850 text-neutral sm:my-8 sm:w-full sm:max-w-6xl sm:p-6">
                                <div>
                                    <div>
                                        <div
                                            class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-base-800">
                                            <TagIcon class="h-6 w-6 text-primary-500" aria-hidden="true" />
                                        </div>
                                        <div class="mt-3 text-center sm:mt-5">
                                            <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                                {{ category?.getName() }}
                                            </DialogTitle>
                                        </div>
                                    </div>
                                    <div class="gap-2 mt-5 sm:mt-4 sm:flex">
                                        <button type="button" v-can="'DocStoreService.DeleteDocumentCategory'"
                                            class="flex-1 rounded-md bg-red-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-red-400"
                                            @click="deleteCategory()" ref="cancelButtonRef">
                                            Delete
                                        </button>
                                        <button type="button"
                                            class="flex-1 rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                            @click="$emit('close')" ref="cancelButtonRef">
                                            Close
                                        </button>
                                    </div>
                                </div>
                            </DialogPanel>
                        </div>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
