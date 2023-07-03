<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { NotificationType } from '~/composables/notification/interfaces/Notification.interface';
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';
import { NOTIFICATION_CATEGORY } from '~~/gen/ts/resources/notifications/notifications';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();
const notifications = useNotificationsStore();

const { getLastId } = storeToRefs(notifications);
const { setLastId } = notifications;
const { accessToken, activeChar } = storeToRefs(authStore);
const { setAccessToken, setActiveChar, setPermissions, setJobProps } = authStore;

const abort = ref<AbortController | undefined>();

// In seconds
const initialBackoffTime = 2;
let restartBackoffTime = 0;

async function streamNotifications(): Promise<void> {
    if (abort.value !== undefined) return;

    console.debug('Notificator: Stream starting, starting at ID:', getLastId.value);
    try {
        abort.value = new AbortController();

        const call = $grpc.getNotificatorClient().stream(
            {
                lastId: getLastId.value,
            },
            {
                abort: abort.value.signal,
            }
        );

        for await (let resp of call.responses) {
            if (resp.lastId > getLastId.value) setLastId(resp.lastId);

            console.debug('Notifications:', resp.notifications);
            resp.notifications.forEach((n) => {
                let nType: NotificationType = (n.type as NotificationType) ?? 'info';

                switch (n.category) {
                    case NOTIFICATION_CATEGORY.GENERAL:
                        notifications.dispatchNotification({
                            title: n.title!,
                            content: n.content!,
                            type: nType,
                            category: n.category,
                            data: n.data,
                        });
                        break;

                    default:
                        notifications.dispatchNotification({
                            title: n.title!,
                            content: n.content!,
                            type: nType,
                            category: n.category,
                            data: n.data,
                        });
                        break;
                }
            });

            // If the response contains an (updated) token
            if (resp.token) {
                const tokenUpdate = resp.token;

                // Update active char when updated user info is received
                if (tokenUpdate.userInfo) {
                    console.debug('Notificator: Updated UserInfo received');

                    setActiveChar(tokenUpdate.userInfo);
                    setPermissions(tokenUpdate.permissions);
                    if (tokenUpdate.jobProps) {
                        setJobProps(tokenUpdate.jobProps!);
                    } else {
                        setJobProps(null);
                    }
                }

                if (tokenUpdate.newToken && tokenUpdate.expires) {
                    console.debug('Notificator: New Token received');

                    setAccessToken(tokenUpdate.newToken, toDate(tokenUpdate.expires) as null | Date);

                    notifications.dispatchNotification({
                        title: { key: 'notifications.renewed_token.title', parameters: [] },
                        content: { key: 'notifications.renewed_token.content', parameters: [] },
                        type: 'info',
                    });
                }
            }

            if (resp.restartStream) {
                console.debug('Notificator: Server requested stream to be restarted');
                restartBackoffTime = 0;
                restartStream();
            }
        }

        console.debug('Notificator: Stream ended');
    } catch (e) {
        const err = e as RpcError;
        if (err.code == 'CANCELLED') {
            return;
        }

        console.debug('Notificator: Stream errored', err);
        restartStream();
    }
}

async function cancelStream(): Promise<void> {
    if (abort.value === undefined) {
        return;
    }
    console.debug('Notificator: Stream cancelled');

    abort.value?.abort();
    abort.value = undefined;
}

async function restartStream(): Promise<void> {
    // Reset back off time when over 3 minutes
    if (restartBackoffTime > 180) {
        restartBackoffTime = initialBackoffTime;
    } else {
        restartBackoffTime += initialBackoffTime;
    }

    await cancelStream();
    console.debug('Notificator: Restart back off time in', restartBackoffTime, 'seconds');
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
        notifications.$reset();
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
