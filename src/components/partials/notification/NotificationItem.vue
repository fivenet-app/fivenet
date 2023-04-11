<script lang="ts" setup>
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
    <transition enter-active-class="transition duration-300 ease-out transform"
        enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
        enter-to-class="translate-y-0 opacity-100 sm:translate-x-0" leave-active-class="transition duration-100 ease-in"
        leave-from-class="opacity-100" leave-to-class="opacity-0">
        <div class="z-50 w-full max-w-sm overflow-hidden bg-base-800 rounded-lg pointer-events-auto shadow-float text-neutral">
            <div class="p-4">
                <div class="flex items-start">
                    <div class="flex-shrink-0 w-8 my-auto" v-if="notification.type">
                        <CheckCircleIcon v-if="notification.type === 'success'" class="text-success-400" aria-hidden="true" />
                        <InformationCircleIcon v-else-if="notification.type === 'info'" class="text-info-400" aria-hidden="true" />
                        <ExclamationTriangleIcon v-else-if="notification.type === 'warning'" class="text-warn-400" aria-hidden="true" />
                        <ExclamationCircleIcon v-else-if="notification.type === 'error'" class="text-error-400" aria-hidden="true" />
                    </div>
                    <div class="ml-3 w-0 flex-1 pt-0.5">
                        <p class="text-sm font-semibold" v-if="notification.title">
                            {{ notification.title }}
                        </p>
                        <p :class="`${notification.title ? 'mt-1' : ''} text-sm leading-5`">
                            {{ notification.content }}
                        </p>
                    </div>
                    <div class="flex flex-shrink-0 ml-4">
                        <button @click="() => closeNotification(notification.id)" type="button"
                            class="inline-flex text-neutral hover:text-base-300 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2">
                            <span class="sr-only">Close</span>
                            <XMarkIcon class="w-5 h-5" aria-hidden="true" />
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </transition>
</template>
