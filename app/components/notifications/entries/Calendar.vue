<script lang="ts" setup>
import type { Notification } from '~~/gen/ts/resources/notifications/notifications';
import Base from './Base.vue';
import { notificationToIcon } from '../helpers';

defineProps<{ notification: Notification; full?: boolean }>();
defineEmits<{
    (e: 'navigate', notification: Notification): void;
    (e: 'state-changed', notification: Notification): void;
}>();
</script>

<template>
    <Base
        :notification="notification"
        :icon="notificationToIcon(notification)"
        :full="full"
        @navigate="$emit('navigate', $event)"
        @state-changed="$emit('state-changed', $event)"
    />
</template>
