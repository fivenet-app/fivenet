<script lang="ts" setup>
import { useClipboard } from '@vueuse/core';
import { PhoneIcon } from 'mdi-vue3';
import { isNUIAvailable, phoneCallNumber } from '~/composables/nui';
import { useNotificatorStore } from '~/store/notificator';

const props = withDefaults(
    defineProps<{
        number: string | undefined;
        showIcon?: boolean;
        hideNumber?: boolean;
        showLabel?: boolean;
        width?: string;
    }>(),
    {
        showIcon: undefined,
        width: 'w-6',
    },
);

const clipboard = useClipboard();

const notifications = useNotificatorStore();

function doCall(): void {
    if (props.number === undefined) return;

    if (isNUIAvailable()) {
        phoneCallNumber(props.number);
    } else {
        notifications.dispatchNotification({
            type: 'info',
            title: {
                key: 'notifications.components.partials.users.PhoneNumber.copied.title',
                parameters: {},
            },
            content: {
                key: 'notifications.components.partials.users.PhoneNumber.copied.content',
                parameters: {},
            },
        });

        clipboard.copy(props.number);
    }
}
</script>

<template>
    <div class="inline-flex items-center">
        <span v-if="number === undefined">N/A</span>
        <template v-else>
            <template v-if="hideNumber === undefined || !hideNumber">
                <span v-for="part in (number ?? '').match(/.{1,3}/g)" class="mr-1">{{ part }}</span>
            </template>

            <button
                v-if="showIcon === undefined || showIcon"
                type="button"
                class="ml-1 flex-initial inline-flex items-center text-primary-500 hover:text-primary-400"
                @click="doCall"
            >
                <PhoneIcon class="h-auto" :class="width" aria-hidden="true" />
                <span v-if="showLabel" class="ml-1">{{ $t('common.call') }}</span>
            </button>
        </template>
    </div>
</template>
