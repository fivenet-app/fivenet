<script setup lang="ts">
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { useTimeoutFn } from '@vueuse/core';
import { BellOutlineIcon, CloseIcon } from 'mdi-vue3';
import NotificationsList from '~/components/partials/notification/NotificationsList.vue';
import { useNotificatorStore } from '~/store/notificator';

const notificatorStore = useNotificatorStore();
const { getNotificationsCount } = storeToRefs(notificatorStore);

const newNotification = ref(false);

const { start } = useTimeoutFn(() => (newNotification.value = false), 1000, {
    immediate: false,
});

const currentCount = ref(0);
watch(getNotificationsCount, () => {
    if (getNotificationsCount.value > currentCount.value) {
        newNotification.value = true;
        start();
    }

    currentCount.value = getNotificationsCount.value;
});

const open = ref(false);
</script>

<template>
    <div class="relative flex-shrink-0">
        <span v-if="getNotificationsCount > 0" class="absolute left-0 top-0 -mr-1 -mt-1 flex h-4 w-4">
            <span
                class="absolute inline-flex h-full w-full rounded-full bg-error-400 opacity-75"
                :class="newNotification ? 'animate-ping' : ''"
            ></span>
            <span
                class="relative inline-flex h-4 w-4 rounded-full bg-error-500 flex justify-center items-center items text-xs text-neutral"
            >
                <span v-if="getNotificationsCount <= 9">{{ getNotificationsCount }} </span>
                <span v-else> 9+ </span>
            </span>
        </span>
        <button
            class="flex rounded-full bg-base-800 text-sm ring-2 ring-base-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
            @click="open = true"
        >
            <span class="sr-only">{{ $t('common.notification', 2) }}</span>
            <BellOutlineIcon
                class="h-10 w-10 rounded-full bg-base-800 p-1 text-base-300 hover:text-base-100 hover:transition-colors"
                aria-hidden="true"
            />
        </button>

        <TransitionRoot as="template" :show="open">
            <Dialog as="div" class="relative z-30" @close="open = false">
                <TransitionChild
                    as="template"
                    enter="ease-in-out duration-500"
                    enter-from="opacity-0"
                    enter-to="opacity-100"
                    leave="ease-in-out duration-500"
                    leave-from="opacity-100"
                    leave-to="opacity-0"
                >
                    <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
                </TransitionChild>

                <div class="fixed inset-0 overflow-hidden">
                    <div class="absolute inset-0 overflow-hidden">
                        <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-xl pl-10 sm:pl-16">
                            <TransitionChild
                                as="template"
                                enter="transform transition ease-in-out duration-100 sm:duration-200"
                                enter-from="translate-x-full"
                                enter-to="translate-x-0"
                                leave="transform transition ease-in-out duration-100 sm:duration-200"
                                leave-from="translate-x-0"
                                leave-to="translate-x-full"
                            >
                                <DialogPanel class="pointer-events-auto w-screen max-w-3xl">
                                    <div class="flex h-full flex-col divide-y divide-gray-200 bg-primary-900 shadow-xl">
                                        <div class="h-0 flex-1 overflow-y-auto">
                                            <div class="bg-primary-700 px-4 py-6 sm:px-6">
                                                <div class="flex items-center justify-between">
                                                    <DialogTitle
                                                        class="inline-flex text-base font-semibold leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.notification', 2) }}
                                                    </DialogTitle>
                                                    <div class="ml-3 flex h-7 items-center">
                                                        <button
                                                            type="button"
                                                            class="rounded-md bg-gray-100 text-gray-500 hover:text-gray-400 focus:outline-none focus:ring-2 focus:ring-neutral"
                                                            @click="open = false"
                                                        >
                                                            <span class="sr-only">{{ $t('common.close') }}</span>
                                                            <CloseIcon class="h-5 w-5" aria-hidden="true" />
                                                        </button>
                                                    </div>
                                                </div>
                                            </div>
                                            <div class="flex flex-1 flex-col justify-between">
                                                <div class="divide-y divide-gray-200 px-2">
                                                    <div class="mt-1">
                                                        <div class="my-2 space-y-24">
                                                            <div class="form-control flex-1">
                                                                <NotificationsList :compact="true" @clicked="open = false" />
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="flex flex-shrink-0 justify-end px-4 py-4">
                                            <span class="isolate inline-flex w-full rounded-md pr-4 shadow-sm">
                                                <button
                                                    type="button"
                                                    class="relative inline-flex w-full items-center rounded-l-md bg-primary-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline-primary-500"
                                                    @click="
                                                        open = false;
                                                        navigateTo({ name: 'notifications' });
                                                    "
                                                >
                                                    {{ $t('components.partials.sidebar_notifications') }}
                                                </button>
                                                <button
                                                    type="button"
                                                    class="relative -ml-px inline-flex w-full items-center rounded-r-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-200 hover:text-gray-900"
                                                    @click="open = false"
                                                >
                                                    {{ $t('common.close', 1) }}
                                                </button>
                                            </span>
                                        </div>
                                    </div>
                                </DialogPanel>
                            </TransitionChild>
                        </div>
                    </div>
                </div>
            </Dialog>
        </TransitionRoot>
    </div>
</template>
