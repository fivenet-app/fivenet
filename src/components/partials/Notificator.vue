<script lang="ts" setup>
import { ClientReadableStream, RpcError } from 'grpc-web';
import { StreamRequest, StreamResponse } from '@fivenet/gen/services/notificator/notificator_pb';
import { useNotificatorStore } from '~/store/notificator';
import { dispatchNotification } from '~/components/partials/notification';
import { NotificationType } from '~/components/partials/notification/interfaces';
import { useAuthStore } from '~/store/auth';

const { $grpc } = useNuxtApp();
const store = useNotificatorStore();
const authStore = useAuthStore();

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
            if (resp.getLastId() > store.$state.lastId) {
                store.setLastId(resp.getLastId());
            }
            resp.getNotificationsList().forEach(v => {
                let nType: NotificationType = 'info';
                dispatchNotification({ title: v.getTitle(), content: v.getContent(), type: nType });
            });
        }).
        on('end', async () => {
            console.debug('Notificator Stream Ended');
        });
}

watch(activeChar, () => {
    // Only stream notifications when a character is active
    if (activeChar.value) {
        streamNotifications();
    } else {
        stream.value?.cancel();
        stream.value = undefined;
    }
});
</script>

<template></template>
