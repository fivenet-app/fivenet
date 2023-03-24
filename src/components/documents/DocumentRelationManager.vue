<script setup lang="ts">
import { DocumentRelation, DOC_RELATION } from '@arpanet/gen/resources/documents/documents_pb';
import { GetDocumentRequest, RemoveDcoumentReferenceRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import {
    Dialog,
    DialogPanel,
    DialogTitle,
    TransitionChild,
    TransitionRoot,
} from '@headlessui/vue';
import { XMarkIcon, ArrowTopRightOnSquareIcon, DocumentMinusIcon } from '@heroicons/vue/24/outline';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router/auto';
import { getDocStoreClient } from '../../grpc/grpc';

const router = useRouter();

const props = defineProps<{
    open: boolean,
    document: number,
}>();

const emit = defineEmits<{
    (e: 'close'): void,
}>();

const relations = ref<DocumentRelation[]>([])

onMounted(() => {
    findRelations();
});

function findRelations(): void {
    if (!props.document) return;

    const req = new GetDocumentRequest();
    req.setDocumentId(props.document);

    getDocStoreClient().
        getDocumentRelations(req, null).
        then((resp) => {
            relations.value = resp.getRelationsList();
        });
}

function removeReference(id: number): void {
    const req = new RemoveDcoumentReferenceRequest();
    req.setId(id);

    getDocStoreClient().removeDcoumentReference(req, null).then(() => {
        findRelations();
    });
}
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="emit('close')">
            <TransitionChild as="template" enter="ease-out duration-300" enter-from="opacity-0" enter-to="opacity-100"
                leave="ease-in duration-200" leave-from="opacity-100" leave-to="opacity-0">
                <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
            </TransitionChild>

            <div class="fixed inset-0 z-10 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild as="template" enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100" leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95">
                        <DialogPanel
                            class="relative transform overflow-hidden rounded-lg bg-white px-4 pt-5 pb-4 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-6xl sm:p-6">
                            <div class="absolute top-0 right-0 hidden pt-4 pr-4 sm:block">
                                <button type="button"
                                    class="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                                    @click="emit('close')">
                                    <span class="sr-only">Close</span>
                                    <XMarkIcon class="h-6 w-6" aria-hidden="true" />
                                </button>
                            </div>
                            <DialogTitle as="h3" class="text-base font-semibold leading-6 text-gray-900">Document Relations
                            </DialogTitle>
                            <div class="sm:flex sm:items-start mt-2">
                                <div class="px-4 sm:px-6 lg:px-8 w-full">
                                    <div class="mt-8 flow-root">
                                        <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
                                            <div class="inline-block min-w-full py-2 align-middle">
                                                <table class="min-w-full divide-y divide-gray-300">
                                                    <thead>
                                                        <tr>
                                                            <th scope="col"
                                                                class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6 lg:pl-8">
                                                                Name</th>
                                                            <th scope="col"
                                                                class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                                                                Creator</th>
                                                            <th scope="col"
                                                                class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                                                                Relation</th>
                                                            <th scope="col"
                                                                class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                                                                Actions</th>
                                                        </tr>
                                                    </thead>
                                                    <tbody class="divide-y divide-gray-200 bg-white">
                                                        <tr v-for="ref in relations" :key="ref.getId()">
                                                            <td
                                                                class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6 lg:pl-8 truncate">
                                                                {{
                                                                    ref.getTargetUser()?.getFirstname() }} {{
                                                                    ref.getTargetUser()?.getLastname() }}

                                                            </td>
                                                            <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{
                                                                ref.getSourceUser()?.getFirstname() }}
                                                                {{ ref.getSourceUser()?.getLastname() }}
                                                            </td>
                                                            <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{
                                                                DOC_RELATION[ref.getRelation()]?.toString() ?? ref.getRelation() }}</td>
                                                            <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                                                                <div class="flex flex-row gap-2">
                                                                    <div class="flex">
                                                                        <a :href="router.resolve({ name: 'Citizens: Info', params: { id: ref.getTargetUserId() } }).href" target="_blank">
                                                                            <ArrowTopRightOnSquareIcon
                                                                                class="w-6 h-auto text-indigo-700 hover:text-indigo-500">
                                                                            </ArrowTopRightOnSquareIcon>
                                                                        </a>
                                                                    </div>
                                                                    <div class="flex">
                                                                        <button role="button"
                                                                            @click="removeReference(ref.getId())">
                                                                            <DocumentMinusIcon
                                                                                class="w-6 h-auto text-red-700 hover:text-red-500">
                                                                            </DocumentMinusIcon>
                                                                        </button>
                                                                    </div>
                                                                </div>
                                                            </td>
                                                        </tr>
                                                    </tbody>
                                                </table>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse gap-2">
                                <button type="button"
                                    class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto"
                                    @click="emit('close')">Cancel</button>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
