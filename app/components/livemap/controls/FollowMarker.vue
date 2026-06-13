<script lang="ts" setup>
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { t } = useI18n();

const livemapStore = useLivemapStore();
const { gotoCoords } = livemapStore;
const { followMarker, selectedMarker, showLocationMarker } = storeToRefs(livemapStore);

const notifications = useNotificationsStore();

watch(followMarker, (val) => {
    if (val && selectedMarker.value) {
        showLocationMarker.value = false;

        gotoCoords(selectedMarker.value, true);

        notifications.add({
            title: { key: 'notifications.livemap.follow_marker.title' },
            description: {
                key: 'notifications.livemap.follow_marker.content',
                parameters: {
                    name: selectedMarker.value?.user ? userToLabel(selectedMarker.value?.user) : t('common.unknown'),
                },
            },
            type: NotificationType.SUCCESS,
            duration: 3000,
        });
    } else {
        notifications.add({
            title: { key: 'notifications.livemap.unfollow_marker.title' },
            description: {
                key: 'notifications.livemap.unfollow_marker.content',
                parameters: {
                    name: selectedMarker.value?.user ? userToLabel(selectedMarker.value?.user) : t('common.unknown'),
                },
            },
            type: NotificationType.INFO,
            duration: 3000,
        });
    }
});

function stopFollowing(): void {
    followMarker.value = false;
    selectedMarker.value = undefined;
}

watch(
    selectedMarker,
    async () => {
        if (!selectedMarker.value || !followMarker.value) return;

        gotoCoords(selectedMarker.value, true);
    },
    { deep: true },
);
</script>

<template>
    <LControl position="topleft">
        <UTooltip v-if="followMarker && selectedMarker" :text="$t('components.livemap.unfollow_marker')">
            <UButton
                class="inset-0 inline-flex items-center justify-center rounded-md border border-black/20 bg-clip-padding text-black"
                icon="i-mdi-track-changes"
                size="xs"
                block
                @click="stopFollowing"
            />
        </UTooltip>
    </LControl>
</template>
