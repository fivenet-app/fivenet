<script lang="ts" setup>
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { t } = useI18n();

const livemapStore = useLivemapStore();
const { gotoCoords } = livemapStore;
const { followMarker, selectedMarker: marker, showLocationMarker } = storeToRefs(livemapStore);

const notifications = useNotificationsStore();

watch(followMarker, (val) => {
    if (val && marker.value) {
        showLocationMarker.value = false;

        gotoCoords(marker.value, true);

        notifications.add({
            title: { key: 'notifications.livemap.follow_marker.title' },
            description: {
                key: 'notifications.livemap.follow_marker.content',
                parameters: { name: marker.value?.user ? userToLabel(marker.value?.user) : t('common.unknown') },
            },
            type: NotificationType.SUCCESS,
            duration: 3000,
        });
    } else {
        notifications.add({
            title: { key: 'notifications.livemap.unfollow_marker.title' },
            description: {
                key: 'notifications.livemap.unfollow_marker.content',
                parameters: { name: marker.value?.user ? userToLabel(marker.value?.user) : t('common.unknown') },
            },
            type: NotificationType.INFO,
            duration: 3000,
        });
    }
});

watch(
    marker,
    async () => {
        if (!marker.value || !followMarker.value) return;

        gotoCoords(marker.value, true);
    },
    { deep: true },
);
</script>

<template>
    <LControl position="topleft">
        <UTooltip v-if="followMarker" :text="$t('components.livemap.unfollow_marker')">
            <UButton
                class="inset-0 inline-flex items-center justify-center rounded-md border border-black/20 bg-clip-padding text-black"
                icon="i-mdi-track-changes"
                size="xs"
                block
                @click="followMarker = false"
            />
        </UTooltip>
    </LControl>
</template>
