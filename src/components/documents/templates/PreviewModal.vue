<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { CloseIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/store/auth';
import { useClipboardStore } from '~/store/clipboard';
import { Template } from '~~/gen/ts/resources/documents/templates';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();
const clipboardStore = useClipboardStore();

const { activeChar } = storeToRefs(authStore);

const props = defineProps<{
    id: bigint;
    open: boolean;
}>();

defineEmits<{
    (e: 'close'): void;
}>();

const { data: template, pending, refresh, error } = useLazyAsyncData(`documents-templates-${props.id}`, () => getTemplate());

async function getTemplate(): Promise<Template> {
    try {
        const data = clipboardStore.getTemplateData();
        data.activeChar = activeChar.value!;
        console.debug('Documents: Editor - Clipboard Template Data', data);

        const call = $grpc.getDocStoreClient().getTemplate({
            templateId: props.id,
            data,
            render: true,
        });
        const { response } = await call;

        return response.template!;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="$emit('close')">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 transition-opacity bg-opacity-75 bg-base-900" />
            </TransitionChild>

            <div class="fixed inset-0 z-30 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild
                        as="template"
                        enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100"
                        leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    >
                        <div>
                            <DialogPanel
                                class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-800 text-neutral sm:my-8 w-full sm:w-screen sm:min-w-min sm:p-6"
                            >
                                <div class="absolute right-0 top-0 pr-4 pt-4 block">
                                    <button
                                        type="button"
                                        class="rounded-md bg-neutral text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                        @click="$emit('close')"
                                    >
                                        <span class="sr-only">{{ $t('common.close') }}</span>
                                        <CloseIcon class="h-6 w-6" aria-hidden="true" />
                                    </button>
                                </div>
                                <div>
                                    <div class="mt-3 text-center sm:mt-5">
                                        <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                            {{ $t('common.document', 1) }}
                                            {{ $t('common.preview') }}
                                        </DialogTitle>
                                        <div class="mt-2">
                                            <div class="pt-4">
                                                <DataPendingBlock
                                                    v-if="pending"
                                                    :message="$t('common.loading', [$t('common.template', 2)])"
                                                />
                                                <DataErrorBlock
                                                    v-else-if="error"
                                                    :title="$t('common.unable_to_load', [$t('common.template', 2)])"
                                                    :retry="refresh"
                                                />
                                                <DataNoDataBlock v-else-if="!template" :type="$t('common.template', 2)" />
                                                <template v-else>
                                                    <div>
                                                        <label class="block mb-2 text-sm font-medium leading-6 text-neutral">
                                                            {{ $t('common.title') }}
                                                        </label>
                                                        <h1
                                                            class="p-2 mt-4 rounded-lg text-2xl font-bold text-neutral bg-base-800 break-words"
                                                        >
                                                            {{ template?.title }}
                                                        </h1>
                                                    </div>
                                                    <div>
                                                        <label class="block mb-2 text-sm font-medium leading-6 text-neutral">
                                                            {{ $t('common.state') }}
                                                        </label>
                                                        <p
                                                            class="p-2 mt-4 rounded-lg text-base font-bold text-neutral bg-base-800 break-words"
                                                        >
                                                            {{ template?.state }}
                                                        </p>
                                                    </div>

                                                    <label class="block mb-2 text-sm font-medium leading-6 text-neutral">
                                                        {{ $t('common.content') }}
                                                    </label>
                                                    <div class="p-2 mt-4 rounded-lg text-neutral bg-base-800 break-words">
                                                        <!-- eslint-disable-next-line vue/no-v-html -->
                                                        <p v-html="template?.content"></p>
                                                    </div>
                                                </template>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div class="gap-2 mt-5 sm:mt-4 sm:flex">
                                    <button
                                        type="button"
                                        class="flex-1 rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.close', 1) }}
                                    </button>
                                </div>
                            </DialogPanel>
                        </div>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
