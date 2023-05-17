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

const { accessToken, activeChar } = storeToRefs(authStore);
const { setAccessToken, setActiveChar, setPermissions, setJobProps } = authStore;

const stream = ref<ClientReadableStream<StreamResponse> | undefined>(undefined);

// In seconds
const initialBackoffTime = 2;
let restartBackoffTime = 0;

async function streamNotifications(): Promise<void> {
    if (stream.value !== undefined) return;

    const request = new StreamRequest();
    request.setLastId(store.getLastId);

    stream.value = $grpc.getNotificatorClient().
        stream(request).
        on('error', async (err: RpcError) => {
            console.debug('Notificator: Stream errored', err);
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

                // Update active char when updated user info is received
                if (tokenUpdate.hasUserInfo()) {
                    setActiveChar(tokenUpdate.getUserInfo()!);
                    setPermissions(tokenUpdate.getPermissionsList());
                }
                if (tokenUpdate.hasJobProps()) {
                    setJobProps(tokenUpdate.getJobProps()!);
                } else {
                    setJobProps(null);
                }

                if (tokenUpdate.hasNewToken() && tokenUpdate.hasExpires()) {
                    setAccessToken(tokenUpdate.getNewToken(), toDate(tokenUpdate.getExpires()) as null | Date);

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
                console.debug('Notificator: Server requested stream to be restarted');
                restartBackoffTime = 0;
                restartStream();
            }
        }).
        on('end', async () => {
            console.debug('Notificator: Stream ended');
            restartStream();
        });

    console.debug('Notificator: Stream started');
}

async function cancelStream(): Promise<void> {
    if (stream.value === undefined) {
        return;
    }

    stream.value?.cancel();
    stream.value = undefined;
}

async function restartStream(): Promise<void> {
    // Reset back off time when over 3 minutes
    if (restartBackoffTime > 180) {
        restartBackoffTime = initialBackoffTime;
    } else {
        restartBackoffTime += initialBackoffTime;
    }

    await cancelStream();
    console.debug('Notificator: Restart back off time in', restartBackoffTime, "seconds");
    setTimeout(async () => {
        toggleStream();
    }, restartBackoffTime * 1000);
}

async function toggleStream(): Promise<void> {
    // Only stream notifications when a character is active
    if (accessToken.value !== null && activeChar.value !== null) {
        streamNotifications();
    } else {
        cancelStream();
    }
}

watch(accessToken, async () => toggleStream());
watch(activeChar, async () => toggleStream());

onMounted(async () => {
    toggleStream();
});

onBeforeUnmount(async () => {
    console.debug('Notificator: Unmount');
    cancelStream();
});
</script>

<template></template>
