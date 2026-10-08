<script lang="ts" setup>
import SettingsDrawer from '~/components/livemap/controls/SettingsDrawer.vue';
import { HELP_GUIDE_ACTION_EVENT, type HelpGuideActionDetail } from '~/composables/useHelpGuides';

const overlay = useOverlay();
const settingsModal = overlay.create(SettingsDrawer);

function onHelpGuideAction(event: Event): void {
    const action = event instanceof CustomEvent ? (event.detail as HelpGuideActionDetail | undefined) : undefined;
    if (!action) return;

    if (action.actionId === 'livemap-open-settings' && action.phase === 'enter') {
        settingsModal.open();
    }

    if (action.actionId === 'livemap-close-settings' && action.phase === 'leave') {
        settingsModal.close();
    }
}

onMounted(() => window.addEventListener(HELP_GUIDE_ACTION_EVENT, onHelpGuideAction));
onBeforeUnmount(() => window.removeEventListener(HELP_GUIDE_ACTION_EVENT, onHelpGuideAction));
</script>

<template>
    <LControl position="bottomright">
        <UTooltip :text="$t('common.setting', 2)">
            <UButton
                data-tour="livemap-settings"
                class="inset-0 inline-flex items-center justify-center rounded-md border border-black/20 bg-clip-padding text-black"
                icon="i-mdi-cog"
                size="xs"
                block
                @click="settingsModal.open()"
            />
        </UTooltip>
    </LControl>
</template>
