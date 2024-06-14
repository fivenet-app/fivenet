<script lang="ts" setup>
import type { WebSocketStatus } from '@vueuse/core';
import { v4 as uuidv4 } from 'uuid';
import { useGRPCWebsocketTransport } from '~/composables/grpcws';

const { t } = useI18n();

const { timeouts } = useAppConfig();

const { webSocket, wsInitiated } = useGRPCWebsocketTransport();

const toast = useToast();

const status = useDebounce(webSocket.status, 150);

const notificationId = ref<string | undefined>();

async function checkWebSocketStatus(previousStatus: WebSocketStatus, status: WebSocketStatus): Promise<void> {
    if (status === 'OPEN') {
        if (notificationId.value !== undefined) {
            toast.remove(notificationId.value);
            notificationId.value = undefined;
        }

        toast.add({
            id: uuidv4(),
            color: 'green',
            icon: 'i-mdi-check-network',
            title: t('notifications.grpc_errors.available.title'),
            description: t('notifications.grpc_errors.available.content'),
            timeout: timeouts.notification,
        });
    } else if (previousStatus === 'CONNECTING' && status === 'CLOSED') {
        if (notificationId.value !== undefined) {
            return;
        }

        notificationId.value = uuidv4();
        toast.add({
            id: notificationId.value,
            color: 'red',
            icon: 'i-mdi-close-network',
            title: t('notifications.grpc_errors.unavailable.title'),
            description: t('notifications.grpc_errors.unavailable.content'),
            timeout: 0,
            ui: {
                closeButton: {
                    ui: { rounded: 'rounded-full', base: 'hidden' },
                },
            },
            actions: [
                {
                    label: t('common.refresh'),
                    icon: 'i-mdi-refresh',
                    click: () => reloadNuxtApp({}),
                },
            ],
        });
    }
}

const previousStatus = ref<WebSocketStatus>('OPEN');
const { resume } = watchPausable(
    status,
    async () => {
        if (wsInitiated.value && previousStatus.value !== status.value) {
            checkWebSocketStatus(previousStatus.value, status.value);
            previousStatus.value = status.value;
        }
    },
    {
        immediate: false,
    },
);

useTimeoutFn(() => {
    resume();
}, 3000);
</script>

<template>
    <div v-if="notificationId" class="relative z-50">
        <div class="fixed inset-0 bg-gray-200/75 transition-opacity dark:bg-gray-800/75" />

        <div class="fixed inset-0 overflow-y-auto">
            <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0"></div>
        </div>
    </div>
</template>
