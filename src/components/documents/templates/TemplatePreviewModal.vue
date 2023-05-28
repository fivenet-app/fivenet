<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from 'grpc-web';
import { useAuthStore } from '~/store/auth';
import { useClipboardStore } from '~/store/clipboard';
import { Template } from '~~/gen/ts/resources/documents/templates';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();
const clipboardStore = useClipboardStore();

const { activeChar } = storeToRefs(authStore);

const props = defineProps<{
    id: number;
    open: boolean;
}>();

defineEmits<{
    (e: 'close'): void;
}>();

const { data: template, pending, refresh, error } = useLazyAsyncData(`documents-templates-${props.id}`, () => getTemplate());

async function getTemplate(): Promise<Template> {
    return new Promise(async (res, rej) => {
        try {
            const data = clipboardStore.getTemplateData();
            data.activeChar = activeChar.value!;

            const call = $grpc.getDocStoreClient().getTemplate({
                templateId: props.id,
                data: JSON.stringify(data),
                render: true,
            });
            const { response } = await call;

            return res(response.template!);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
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

            <div class="fixed inset-0 z-10 overflow-y-auto">
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
                                class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-850 text-neutral sm:my-8 sm:w-full sm:w-screen sm:min-w-min sm:p-6"
                            >
                                <div>
                                    <div class="mt-3 text-center sm:mt-5">
                                        <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                            {{ $t('common.document', 1) }}
                                            {{ $t('common.preview') }}
                                        </DialogTitle>
                                        <div class="mt-2">
                                            <div class="pt-4">
                                                <label class="block mb-2 text-sm font-medium leading-6 text-neutral">
                                                    {{ $t('common.title') }}
                                                </label>
                                                <h1
                                                    class="p-2 mt-4 rounded-lg text-2xl font-bold text-neutral bg-base-800 break-words"
                                                >
                                                    {{ template?.title }}
                                                </h1>

                                                <label class="block mb-2 text-sm font-medium leading-6 text-neutral">
                                                    {{ $t('common.content') }}
                                                </label>
                                                <div class="p-2 mt-4 rounded-lg text-neutral bg-base-800 break-words">
                                                    <p v-html="template?.content"></p>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div class="gap-2 mt-5 sm:mt-4 sm:flex">
                                    <button
                                        type="button"
                                        class="flex-1 rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                        @click="$emit('close')"
                                        ref="cancelButtonRef"
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
