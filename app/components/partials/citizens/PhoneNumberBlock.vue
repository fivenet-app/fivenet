<script lang="ts" setup>
import { phoneCallNumber } from '~/composables/nui';
import { useSettingsStore } from '~/stores/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = withDefaults(
    defineProps<{
        number: string | undefined;
        showIcon?: boolean;
        hideNumber?: boolean;
        showLabel?: boolean;
        hideNaText?: boolean;
        disableTruncate?: boolean;
    }>(),
    {
        showIcon: true,
        hideNumber: false,
        showLabel: false,
        hideNaText: false,
        disableTruncate: false,
    },
);

defineOptions({
    inheritAttrs: false,
});

const settingsStore = useSettingsStore();
const { nuiEnabled, streamerMode } = storeToRefs(settingsStore);

const notifications = useNotificationsStore();

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
    <div class="inline-flex items-center gap-1">
        <template v-if="number">
            <UTooltip v-if="showIcon" class="w-full" :text="$t('common.call')">
                <UButton class="shrink-0" variant="link" icon="i-mdi-phone" @click="doCall">
                    <span v-if="showLabel" :class="!disableTruncate && 'truncate'">{{ $t('common.call') }}</span>
                </UButton>
            </UTooltip>

            <span v-if="!hideNumber" class="inline-flex gap-1" :class="[streamerMode ? 'blur' : '']">
                <span v-for="(part, idx) in (number ?? '').match(/.{1,3}/g)" :key="idx">{{ part }}</span>
            </span>
        </template>

        <template v-else-if="!hideNaText">
            {{ $t('common.na') }}
        </template>
    </div>
</template>
