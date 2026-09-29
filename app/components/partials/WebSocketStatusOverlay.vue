<script lang="ts" setup>
import { useDocumentVisibility, useOnline, type WebSocketStatus } from '@vueuse/core';
import { v4 as uuidv4 } from 'uuid';
import { useGRPCWebsocketTransport } from '~/composables/grpcws';

const { t } = useI18n();

const { timeouts } = useAppConfig();

const { webSocket } = useGRPCWebsocketTransport();

const toast = useToast();

const status = useDebounce(webSocket.status, 150);
const visibility = useDocumentVisibility();
const online = useOnline();

const notificationId = ref<string | undefined>();
const unavailableTimer = ref<ReturnType<typeof setTimeout> | undefined>();

const unavailableDelay = 10_000;

const isConnectionCheckAllowed = (): boolean => visibility.value === 'visible' && online.value;

function clearUnavailableTimer(): void {
    if (unavailableTimer.value === undefined) return;

    clearTimeout(unavailableTimer.value);
    unavailableTimer.value = undefined;
}

function showUnavailableNotification(): void {
    if (!isConnectionCheckAllowed() || status.value === 'OPEN' || notificationId.value !== undefined) return;

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
}

function scheduleUnavailableNotification(): void {
    if (!isConnectionCheckAllowed() || unavailableTimer.value !== undefined) return;

    unavailableTimer.value = setTimeout(() => {
        unavailableTimer.value = undefined;
        showUnavailableNotification();
    }, unavailableDelay);
}

async function checkWebSocketStatus(previousStatus: WebSocketStatus, status: WebSocketStatus): Promise<void> {
    if (status === 'OPEN') {
        clearUnavailableTimer();

        if (notificationId.value === undefined) return;

        const wasAllowed = isConnectionCheckAllowed();
        toast.remove(notificationId.value);
        notificationId.value = undefined;

        if (!wasAllowed) return;

        toast.add({
            id: uuidv4(),
            color: 'success',
            icon: 'i-mdi-check-network',
            title: t('notifications.grpc_errors.available.title'),
            description: t('notifications.grpc_errors.available.content'),
            duration: timeouts.notification,
        });
    } else if (previousStatus === 'CONNECTING' && status === 'CLOSED') {
        scheduleUnavailableNotification();
    }
}

const previousStatus = ref<WebSocketStatus>('OPEN');
const { resume } = watch(
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

watch([visibility, online], ([currentVisibility, isOnline], [previousVisibility, wasOnline]) => {
    if (currentVisibility === previousVisibility && isOnline === wasOnline) return;

    clearUnavailableTimer();
    previousStatus.value = status.value;

    // A backgrounded tab can lose a healthy socket while its heartbeat is suspended.
    // Start a fresh failure observation when the tab is visible and online again.
    if (currentVisibility !== 'visible' || !isOnline) return;

    if (status.value !== 'OPEN') {
        previousStatus.value = status.value;
    }
});

useTimeoutFn(() => {
    resume();
}, 2750);
</script>

<template>
    <div></div>
</template>
