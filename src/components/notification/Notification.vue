<script setup lang="ts">
import store from './store';
import type { Notification } from './interfaces';

import { CheckCircleIcon, ExclamationCircleIcon, ExclamationTriangleIcon, InformationCircleIcon } from '@heroicons/vue/24/outline';
import { XMarkIcon } from '@heroicons/vue/20/solid';

defineProps<{
    notification: Notification,
}>();

const closeNotification = (id: string) => {
    store.actions.removeNotification(id);
}
</script>

<template>
    <transition enter-active-class="transform ease-out duration-300 transition"
        enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
        enter-to-class="translate-y-0 opacity-100 sm:translate-x-0" leave-active-class="transition ease-in duration-100"
        leave-from-class="opacity-100" leave-to-class="opacity-0">
        <div class="z-20 pointer-events-auto w-full max-w-sm overflow-hidden rounded-lg bg-white shadow-lg ring-1 ring-black ring-opacity-5">
            <div class="p-4">
                <div class="flex items-start">
                    <div class="flex-shrink-0" v-if="notification.type">
                        <CheckCircleIcon v-if="notification.type === 'success'" class="h-6 w-6 text-green-400" aria-hidden="true" />
                        <InformationCircleIcon v-else-if="notification.type === 'info'" class="h-6 w-6 text-blue-400" aria-hidden="true" />
                        <ExclamationTriangleIcon v-else-if="notification.type === 'warning'" class="h-6 w-6 text-yellow-400" aria-hidden="true" />
                        <ExclamationCircleIcon v-else-if="notification.type === 'error'" class="h-6 w-6 text-red-400" aria-hidden="true" />
                    </div>
                    <div class="ml-3 w-0 flex-1 pt-0.5">
                        <p class="text-sm font-medium text-gray-900" v-if="notification.title">
                            {{ notification.title }}
                        </p>
                        <p :class="`${notification.title ? 'mt-1' : ''} text-sm leading-5 text-gray-500`">
                            {{ notification.content }}
                        </p>
                    </div>
                    <div class="ml-4 flex flex-shrink-0">
                        <button @click="() => closeNotification(notification.id)" type="button"
                            class="inline-flex rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
                            <span class="sr-only">Close</span>
                            <XMarkIcon class="h-5 w-5" aria-hidden="true" />
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </transition>
</template>
