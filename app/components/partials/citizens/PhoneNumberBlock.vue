<script lang="ts" setup>
import { phoneCallNumber } from '~/composables/nui';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = withDefaults(
    defineProps<{
        number: string | undefined;
        showIcon?: boolean;
        hideNumber?: boolean;
        showLabel?: boolean;
        width?: string;
        padded?: boolean;
        hideNa?: boolean;
    }>(),
    {
        showIcon: true,
        hideNumber: false,
        showLabel: false,
        width: 'w-6',
        hideNa: false,
    },
);

const settingsStore = useSettingsStore();
const { nuiEnabled, streamerMode } = storeToRefs(settingsStore);

const notifications = useNotificatorStore();

async function doCall(): Promise<void> {
    if (props.number === undefined) {
        return;
    }

    if (nuiEnabled.value) {
        return phoneCallNumber(props.number);
    } else {
        notifications.add({
            title: {
                key: 'notifications.components.partials.users.PhoneNumber.copied.title',
                parameters: {},
            },
            description: {
                key: 'notifications.components.partials.users.PhoneNumber.copied.content',
                parameters: {},
            },
            type: NotificationType.INFO,
        });

        return copyToClipboardWrapper(props.number);
    }
}
</script>

<template>
    <div class="inline-flex items-center" :class="!padded && 'gap-1'">
        <span v-if="number === undefined">{{ $t('common.na') }}</span>
        <template v-else>
            <UButton v-if="showIcon" variant="link" icon="i-mdi-phone" :padded="padded" @click="doCall">
                <span class="sr-only">{{ $t('common.call') }}</span>
                <span v-if="showLabel" class="truncate">{{ $t('common.call') }}</span>
            </UButton>

            <span v-if="!hideNumber" class="inline-flex gap-1" :class="[streamerMode ? 'blur' : '']">
                <span v-for="(part, idx) in (number ?? '').match(/.{1,3}/g)" :key="idx">{{ part }}</span>
            </span>
        </template>
    </div>
</template>
