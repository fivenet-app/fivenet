<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';

const authStore = useAuthStore();
const { accessToken, activeChar } = storeToRefs(authStore);

const notifications = useNotificatorStore();
const { startStream, stopStream } = notifications;

async function toggleStream(): Promise<void> {
    // Only stream notifications when a user is logged in and has a character selected
    if (accessToken.value !== null && activeChar.value !== null) {
        return startStream();
    } else {
        await stopStream();
        notifications.$reset();
    }
}

watch(accessToken, async () => toggleStream());
watch(activeChar, async () => toggleStream());

onBeforeMount(async () => toggleStream());

onBeforeUnmount(async () => stopStream());
</script>

<!-- eslint-disable vue/valid-template-root -->
<template></template>
