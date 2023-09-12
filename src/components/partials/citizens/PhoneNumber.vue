<script lang="ts" setup>
import { useClipboard } from '@vueuse/core';
import { PhoneIcon } from 'mdi-vue3';
import { isNUIAvailable, phoneCallNumber } from '~/components/nui';
import { useNotificationsStore } from '~/store/notifications';

const props = withDefaults(
    defineProps<{
        number: string | undefined;
        showIcon?: boolean;
    }>(),
    {
        showIcon: undefined,
    },
);

const clipboard = useClipboard();

const notifications = useNotificationsStore();

function doCall(): void {
    if (props.number === undefined) return;

    if (isNUIAvailable()) {
        phoneCallNumber(props.number);
    } else {
        notifications.dispatchNotification({
            type: 'info',
            title: {
                key: 'notifications.components.partials.users.PhoneNumber.copied.title',
                parameters: [],
            },
            content: {
                key: 'notifications.components.partials.users.PhoneNumber.copied.content',
                parameters: [],
            },
        });

        clipboard.copy(props.number);
    }
}
</script>

<template>
    <div class="flex inline-flex items-center">
        <span v-if="number === undefined">N/A</span>
        <template v-else>
            <span v-for="part in (number ?? '').match(/.{1,3}/g)" class="mr-1">{{ part }}</span>

            <button
                v-if="showIcon === undefined || showIcon"
                type="button"
                class="ml-1 flex-initial text-primary-500 hover:text-primary-400"
                @click="doCall"
            >
                <PhoneIcon class="w-6 h-auto" aria-hidden="true" />
            </button>
        </template>
    </div>
</template>
