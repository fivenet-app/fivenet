<script lang="ts" setup>
import { HELP_GUIDE_ACTION_EVENT, type HelpGuideActionDetail } from '~/composables/useHelpGuides';
import DispatchLayoutSettings from '~/components/user-settings/preferences/DispatchLayoutSettings.vue';

const props = withDefaults(
    defineProps<{
        hideInnerPanes?: boolean;
    }>(),
    {
        hideInnerPanes: false,
    },
);

const open = ref(false);

function onHelpGuideAction(event: Event): void {
    const action = event instanceof CustomEvent ? (event.detail as HelpGuideActionDetail | undefined) : undefined;
    if (!action) return;

    if (action.actionId === 'dispatch-center-open-layout' && action.phase === 'enter') {
        open.value = true;
    }

    if (action.actionId === 'dispatch-center-close-layout' && action.phase === 'leave') {
        open.value = false;
    }
}

onMounted(() => window.addEventListener(HELP_GUIDE_ACTION_EVENT, onHelpGuideAction));
onBeforeUnmount(() => window.removeEventListener(HELP_GUIDE_ACTION_EVENT, onHelpGuideAction));
</script>

<template>
    <UPopover v-model:open="open" :content="{ align: 'end', sideOffset: 8 }">
        <UTooltip :text="$t('components.dispatch.layout_popover.trigger')">
            <UButton
                data-tour="dispatch-center-layout"
                color="neutral"
                variant="ghost"
                icon="i-mdi-view-split-vertical"
                :aria-label="$t('components.dispatch.layout_popover.trigger')"
            />
        </UTooltip>

        <template #content>
            <div
                data-tour="dispatch-center-layout-popover"
                class="max-w-[calc(100vw-2rem)] p-4"
                :class="!props.hideInnerPanes ? 'w-[28rem]' : 'w-[24rem]'"
            >
                <DispatchLayoutSettings :key="String(open)" :hide-inner-panes="props.hideInnerPanes" />
            </div>
        </template>
    </UPopover>
</template>
