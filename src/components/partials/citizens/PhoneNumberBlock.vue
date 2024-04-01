<script lang="ts" setup>
import { PhoneIcon } from 'mdi-vue3';
import { isNUIAvailable, phoneCallNumber } from '~/composables/nui';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';

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

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const notifications = useNotificatorStore();

async function doCall(): Promise<void> {
    if (props.number === undefined) {
        return;
    }

    if (isNUIAvailable()) {
        return phoneCallNumber(props.number);
    } else {
        notifications.add({
            type: 'info',
            title: {
                key: 'notifications.components.partials.users.PhoneNumber.copied.title',
                parameters: {},
            },
            description: {
                key: 'notifications.components.partials.users.PhoneNumber.copied.content',
                parameters: {},
            },
        });

        return copyToClipboardWrapper(props.number);
    }
}
</script>

<template>
    <div class="inline-flex items-center">
        <span v-if="number === undefined">N/A</span>
        <template v-else>
            <span
                v-if="showIcon === undefined || showIcon"
                class="mr-1 inline-flex flex-initial items-center text-primary-500 hover:text-primary-400"
                @click="doCall"
            >
                <PhoneIcon class="hidden h-auto sm:block" :class="width" aria-hidden="true" />
                <span v-if="showLabel" class="ml-1">{{ $t('common.call') }}</span>
            </span>

            <span v-if="hideNumber === undefined || !hideNumber" :class="streamerMode ? 'blur' : ''">
                <span v-for="(part, idx) in (number ?? '').match(/.{1,3}/g)" :key="idx" class="mr-1">{{ part }}</span>
            </span>
        </template>
    </div>
</template>
