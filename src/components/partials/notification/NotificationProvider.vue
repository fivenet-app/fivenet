<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { notificationTypeToIcon, notificationTypeToColor } from '~/components/partials/notification/helpers';

const { t } = useI18n();

const authStore = useAuthStore();
const { username, activeChar } = storeToRefs(authStore);

const notificatorStore = useNotificatorStore();
const { notifications } = storeToRefs(notificatorStore);
const { startStream, stopStream } = notificatorStore;

async function toggleStream(): Promise<void> {
    console.log('toggleStream', username.value, activeChar.value);
    // Only stream notifications when a user is logged in and has a character selected
    if (username.value !== null && activeChar.value !== null) {
        return startStream();
    } else {
        await stopStream();
        notificatorStore.$reset();
    }
}

watch(username, async () => toggleStream());
watch(activeChar, async () => toggleStream());

onMounted(async () => toggleStream());

onBeforeUnmount(async () => stopStream());

const toast = useToast();

watchArray(
    notifications,
    (_, _0, added) => {
        added.forEach((notification) => {
            toast.add({
                id: notification.id,
                title: t(notification.title.key, notification.title.parameters ?? {}),
                description: t(notification.description.key, notification.description.parameters ?? {}),
                icon: notificationTypeToIcon(notification.type),
                color: notificationTypeToColor(notification.type),
                timeout: notification.timeout ?? 3500,
                actions: notification.onClick
                    ? [
                          {
                              label: notification.onClickText
                                  ? t(notification.onClickText.key, notification.onClickText.parameters ?? {})
                                  : t('common.click_here'),
                              click: notification.onClick,
                          },
                      ]
                    : [],
                callback: () => {
                    if (notification.callback) notification.callback();
                    if (notification.id) notificatorStore.remove(notification.id);
                },
            });
        });
    },
    { deep: true },
);
</script>

<template>
    <div></div>
</template>
