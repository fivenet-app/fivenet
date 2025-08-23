<script lang="ts" setup>
import type { WebSocketStatus } from '@vueuse/core';
import { v4 as uuidv4 } from 'uuid';
import { useGRPCWebsocketTransport } from '~/composables/grpc/grpcws';

withDefaults(
    defineProps<{
        hideOverlay?: boolean;
    }>(),
    {
        hideOverlay: false,
    },
);

const { t } = useI18n();

const { timeouts } = useAppConfig();

const { webSocket } = useGRPCWebsocketTransport();

const toast = useToast();

const status = useDebounce(webSocket.status, 150);

const notificationId = ref<string | undefined>();

const overlayRef = useTemplateRef('overlayRef');

async function checkWebSocketStatus(previousStatus: WebSocketStatus, status: WebSocketStatus): Promise<void> {
    if (notificationId.value !== undefined && status === 'OPEN') {
        toast.remove(notificationId.value);
        notificationId.value = undefined;

        toast.add({
            id: uuidv4(),
            color: 'success',
            icon: 'i-mdi-check-network',
            title: t('notifications.grpc_errors.available.title'),
            description: t('notifications.grpc_errors.available.content'),
            duration: timeouts.notification,
        });

        overlayRef.value?.blur();
    } else if (previousStatus === 'CONNECTING' && status === 'CLOSED') {
        if (notificationId.value !== undefined) {
            return;
        }

        notificationId.value = uuidv4();
        toast.add({
            id: notificationId.value,
            color: 'error',
            icon: 'i-mdi-close-network',
            title: t('notifications.grpc_errors.unavailable.title'),
            description: t('notifications.grpc_errors.unavailable.content'),
            duration: 0,
            close: false,
            actions: [
                {
                    label: t('common.retrying'),
                    icon: 'i-mdi-circle-arrows',
                    loading: true,
                    active: true,
                    disabled: true,
                },
                {
                    label: t('common.refresh'),
                    icon: 'i-mdi-reload',
                    onClick: () => reloadNuxtApp({}),
                },
            ],
        });

        overlayRef.value?.focus();
    }
}

const previousStatus = ref<WebSocketStatus>('OPEN');
const { resume } = watchPausable(
    status,
    async () => {
        if (previousStatus.value !== status.value) {
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
}, 2750);
</script>

<template>
    <div
        v-if="notificationId && !hideOverlay"
        ref="overlayRef"
        class="z-999999 relative"
        :class="hideOverlay && 'pointer-events-none'"
    >
        <div class="fixed inset-0 bg-gray-200/75 transition-opacity dark:bg-gray-800/75" />

        <div class="fixed inset-0 overflow-y-auto">
            <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                <UIcon class="size-32 animate-spin" name="i-mdi-loading" />
            </div>
        </div>
    </div>
</template>
