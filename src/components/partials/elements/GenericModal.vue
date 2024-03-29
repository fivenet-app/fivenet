<script lang="ts" setup>
import { Dialog, DialogPanel, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { CloseIcon } from 'mdi-vue3';

withDefaults(
    defineProps<{
        open: boolean;
        dialogClass?: unknown;
        unmount?: boolean;
    }>(),
    {
        dialogClass: '' as any,
        unmount: true,
    },
);

defineEmits<{
    (e: 'close'): void;
}>();
</script>

<template>
    <TransitionRoot as="template" :show="open" :unmount="unmount">
        <Dialog as="div" class="relative z-30" :unmount="unmount" @close="$emit('close')">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
                :unmount="unmount"
            >
                <div class="fixed inset-0 bg-base-900/75 transition-opacity" />
            </TransitionChild>

            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                enter-to="opacity-100 translate-y-0 sm:scale-100"
                leave="ease-in duration-200"
                leave-from="opacity-100 translate-y-0 sm:scale-100"
                leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                :unmount="unmount"
            >
                <div class="fixed inset-0 z-30 overflow-y-auto">
                    <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                        <DialogPanel
                            class="relative w-full overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left text-neutral transition-all sm:my-8 sm:max-w-6xl sm:p-6"
                            :class="dialogClass"
                        >
                            <div class="absolute right-0 top-0 block pr-4 pt-4">
                                <button
                                    type="button"
                                    class="rounded-md bg-neutral text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                    @click="$emit('close')"
                                >
                                    <span class="sr-only">{{ $t('common.close') }}</span>
                                    <CloseIcon class="size-5" aria-hidden="true" />
                                </button>
                            </div>
                            <div>
                                <div class="text-center">
                                    <slot />
                                </div>
                            </div>
                            <div class="mt-5 gap-2 sm:mt-4 sm:flex">
                                <button
                                    type="button"
                                    class="w-full flex-1 rounded-md bg-base-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="$emit('close')"
                                >
                                    {{ $t('common.close', 1) }}
                                </button>
                            </div>
                        </DialogPanel>
                    </div>
                </div>
            </TransitionChild>
        </Dialog>
    </TransitionRoot>
</template>
