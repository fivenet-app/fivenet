<script lang="ts" setup>
import {
    centrumSidebarPlacementOptions,
    defaultCentrumSidebarPlacement,
    type CentrumSidebarPlacement,
} from '~/stores/settings';

const settingsStore = useSettingsStore();
const { centrum } = storeToRefs(settingsStore);

const sidebarPlacement = computed<CentrumSidebarPlacement>({
    get: () => centrum.value.centrumSidebarPlacement ?? defaultCentrumSidebarPlacement,
    set: (value) => {
        centrum.value.centrumSidebarPlacement = value;
    },
});

const options = [
    { placement: 'top', class: 'col-start-2 row-start-1', ...centrumSidebarPlacementOptions.top },
    { placement: 'left', class: 'col-start-1 row-start-2', ...centrumSidebarPlacementOptions.left },
    { placement: 'right', class: 'col-start-3 row-start-2', ...centrumSidebarPlacementOptions.right },
    { placement: 'bottom', class: 'col-start-2 row-start-3', ...centrumSidebarPlacementOptions.bottom },
] as const;

function resetLayout(): void {
    sidebarPlacement.value = defaultCentrumSidebarPlacement;
}
</script>

<template>
    <div class="space-y-4">
        <div>
            <h3 class="text-sm font-semibold text-highlighted">
                {{ $t('components.dispatch.sidebar_layout.title') }}
            </h3>
            <p class="text-sm text-muted">
                {{ $t('components.dispatch.sidebar_layout.description') }}
            </p>
        </div>

        <div class="mx-auto grid max-w-sm grid-cols-3 grid-rows-3 place-items-center gap-2">
            <UButton
                v-for="option in options"
                :key="option.placement"
                :color="sidebarPlacement === option.placement ? 'primary' : 'neutral'"
                :variant="sidebarPlacement === option.placement ? 'solid' : 'soft'"
                :icon="option.icon"
                :class="option.class"
                class="w-full"
                :label="$t(option.labelKey)"
                @click="sidebarPlacement = option.placement"
            />
        </div>

        <UButton color="neutral" variant="soft" :label="$t('common.reset')" @click="resetLayout" />
    </div>
</template>
