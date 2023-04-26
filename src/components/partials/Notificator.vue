<script lang="ts" setup>
import { ClientReadableStream, RpcError } from 'grpc-web';
import { StreamRequest, StreamResponse } from '@fivenet/gen/services/notificator/notificator_pb';
import { useNotificatorStore } from '~/store/notificator';
import { useAuthStore } from '~/store/auth';
import { NotificationType } from '~/composables/notification/interfaces/Notification.interface';
import { useNotificationsStore } from '~/store/notifications';

const { $grpc } = useNuxtApp();
const store = useNotificatorStore();
const authStore = useAuthStore();
const notifications = useNotificationsStore();

const accessToken = computed(() => authStore.$state.accessToken);
const activeChar = computed(() => authStore.$state.activeChar);

const stream = ref<ClientReadableStream<StreamResponse> | undefined>(undefined);

async function streamNotifications(): Promise<void> {
    if (stream.value !== undefined) return;

    const request = new StreamRequest();
    request.setLastId(store.$state.lastId);

    stream.value = $grpc.getNotificatorClient().
        stream(request).
        on('error', async (err: RpcError) => {
            stream.value?.cancel();
        }).
        on('data', async (resp) => {
            if (resp.getLastId() > store.$state.lastId)
                store.setLastId(resp.getLastId());

            resp.getNotificationsList().forEach(v => {
                let nType: NotificationType = v.getType() as NotificationType ?? 'info';
                notifications.dispatchNotification({
                    title: v.getTitle(),
                    content: v.getContent(),
                    type: nType
                });
            });
        }).
        on('end', async () => {
            console.debug('Notificator Stream Ended');
            toggleStream();
        });

    console.debug('Notificator Stream Started');
}

async function cancelStream(): Promise<void> {
    stream.value?.cancel();
    stream.value = undefined;
}

async function toggleStream(): Promise<void> {
    // Only stream notifications when a character is active
    if (accessToken.value && activeChar.value) {
        streamNotifications();
    } else {
        cancelStream();
    }
}

watch(accessToken, async () => toggleStream());
watch(activeChar, async () => toggleStream());

onBeforeUnmount(() => {
    cancelStream();
});
</script>

<template></template>
