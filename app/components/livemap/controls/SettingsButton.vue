<script lang="ts" setup>
import SettingsDrawer from '~/components/livemap/controls/SettingsDrawer.vue';

const overlay = useOverlay();
const settingsModal = overlay.create(SettingsDrawer);
</script>

<template>
    <LControl position="bottomright">
        <UTooltip :text="$t('common.setting', 2)">
            <UButton
                class="inset-0 inline-flex items-center justify-center rounded-md border border-black/20 bg-clip-padding text-black"
                icon="i-mdi-cog"
                size="xs"
                block
                @click="settingsModal.open()"
            />
        </UTooltip>
    </LControl>
</template>
