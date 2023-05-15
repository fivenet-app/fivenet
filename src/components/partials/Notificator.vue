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

const accessToken = computed(() => authStore.getAccessToken);
const activeChar = computed(() => authStore.getActiveChar);

const stream = ref<ClientReadableStream<StreamResponse> | undefined>(undefined);

async function streamNotifications(): Promise<void> {
    if (stream.value !== undefined) return;

    const request = new StreamRequest();
    request.setLastId(store.getLastId);

    stream.value = $grpc.getNotificatorClient().
        stream(request).
        on('error', async (err: RpcError) => {
            console.debug('Notificator: Stream errored', err);
            stream.value?.cancel();
            restartStream();
        }).
        on('data', async (resp) => {
            if (resp.getLastId() > store.getLastId)
                store.setLastId(resp.getLastId());

            resp.getNotificationsList().forEach(v => {
                let nType: NotificationType = v.getType() as NotificationType ?? 'info';
                notifications.dispatchNotification({
                    title: v.getTitle(),
                    content: v.getContent(),
                    type: nType
                });
            });

            // If the response contains an (updated) token
            if (resp.hasToken()) {
                const tokenUpdate = resp.getToken()!;
                if (tokenUpdate.hasExpires()) {
                    authStore.setAccessToken(tokenUpdate.getNewToken(), toDate(tokenUpdate.getExpires()) as null | Date);

                    notifications.dispatchNotification({
                        title: 'notifications.renewed_token.title',
                        titleI18n: true,
                        content: 'notifications.renewed_token.content',
                        contentI18n: true,
                        type: 'info'
                    });
                }
            }

            if (resp.getRestartStream()) {
                console.debug('Notificator: Server requested stream to be restarted')
                cancelStream();
                restartStream();
            }
        }).
        on('end', async () => {
            console.debug('Notificator: Stream Ended');
            restartStream();
        });

    console.debug('Notificator: Stream Started');
}

async function cancelStream(): Promise<void> {
    stream.value?.cancel();
    stream.value = undefined;
}

async function restartStream(): Promise<void> {
    setTimeout(async () => {
        toggleStream();
    }, 2250);
}

async function toggleStream(): Promise<void> {
    // Only stream notifications when a character is active
    if (accessToken.value && activeChar.value) {
        streamNotifications();
    } else {
        cancelStream();
    }
}

watch(activeChar, async () => toggleStream());
onMounted(() => {
    streamNotifications();
});

onBeforeUnmount(() => {
    cancelStream();
});
</script>

<template></template>
