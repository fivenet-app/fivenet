<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { ClipboardListIcon } from 'mdi-vue3';
import { useClipboardStore } from '~/store/clipboard';
import ClipboardCitizens from '~/components/clipboard/modal/ClipboardCitizens.vue';
import ClipboardDocuments from '~/components/clipboard/modal/ClipboardDocuments.vue';
import ClipboardVehicles from '~/components/clipboard/modal/ClipboardVehicles.vue';

const clipboardStore = useClipboardStore();

defineProps<{
    open: boolean;
}>();

defineEmits<{
    (e: 'close'): void;
}>();
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
                <div class="fixed inset-0 bg-base-900/75 transition-opacity" />
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
                        <DialogPanel
                            class="relative w-full overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left transition-all sm:my-8 sm:max-w-6xl sm:p-6"
                        >
                            <div>
                                <div class="mx-auto flex size-12 items-center justify-center rounded-full bg-base-700">
                                    <ClipboardListIcon class="size-5 text-primary-500" />
                                </div>
                                <div class="mt-3 text-center sm:mt-5">
                                    <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                        {{ $t('components.clipboard.clipboard_modal.title') }}
                                    </DialogTitle>
                                    <div class="mt-2">
                                        <ClipboardCitizens />
                                        <ClipboardVehicles />
                                        <ClipboardDocuments />
                                    </div>
                                </div>
                            </div>
                            <div class="mt-5 gap-2 sm:mt-4 sm:flex">
                                <span class="isolate inline-flex w-full rounded-md pr-4 shadow-sm">
                                    <UButton
                                        class="relative inline-flex w-full items-center rounded-l-md bg-base-500 px-3.5 py-2.5 text-sm font-semibold hover:bg-base-400"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.close', 1) }}
                                    </UButton>
                                    <UButton
                                        class="relative -ml-px inline-flex w-full items-center rounded-r-md bg-error-500 px-3.5 py-2.5 text-sm font-semibold hover:bg-error-400"
                                        @click="clipboardStore.clear()"
                                    >
                                        {{ $t('components.clipboard.clipboard_modal.clear') }}
                                    </UButton>
                                </span>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
