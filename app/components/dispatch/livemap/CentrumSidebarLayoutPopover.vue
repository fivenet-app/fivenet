<script lang="ts" setup>
import CentrumSidebarPlacementSettings from '~/components/user-settings/preferences/CentrumSidebarPlacementSettings.vue';

const open = ref(false);
</script>

<template>
    <UPopover v-model:open="open" :content="{ align: 'end', sideOffset: 8 }">
        <UTooltip :text="$t('components.dispatch.sidebar_layout.trigger')">
            <UButton
                color="neutral"
                variant="ghost"
                icon="i-mdi-view-split-vertical"
                :aria-label="$t('components.dispatch.sidebar_layout.trigger')"
            />
        </UTooltip>

        <template #content>
            <div class="w-86 max-w-[calc(100vw-2rem)] p-3">
                <CentrumSidebarPlacementSettings />
            </div>
        </template>
    </UPopover>
</template>
