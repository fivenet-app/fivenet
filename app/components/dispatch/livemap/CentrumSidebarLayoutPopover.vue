<script lang="ts" setup>
import {
    centrumSidebarPlacementOptions,
    defaultCentrumSidebarPlacement,
    useSettingsStore,
    type CentrumSidebarPlacement,
} from '~/stores/settings';

const settingsStore = useSettingsStore();
const { centrum } = storeToRefs(settingsStore);

const open = ref(false);

const sidebarPlacement = computed<CentrumSidebarPlacement>(
    () => centrum.value.centrumSidebarPlacement ?? defaultCentrumSidebarPlacement,
);

const centrumSidebarPlacementGridOptions = [
    { placement: 'top', class: 'col-start-2 row-start-1', ...centrumSidebarPlacementOptions.top },
    { placement: 'left', class: 'col-start-1 row-start-2', ...centrumSidebarPlacementOptions.left },
    { placement: 'right', class: 'col-start-3 row-start-2', ...centrumSidebarPlacementOptions.right },
    { placement: 'bottom', class: 'col-start-2 row-start-3', ...centrumSidebarPlacementOptions.bottom },
] as const;

function setLayout(placement: CentrumSidebarPlacement): void {
    centrum.value.centrumSidebarPlacement = placement;
    open.value = false;
}

function resetLayout(): void {
    centrum.value.centrumSidebarPlacement = defaultCentrumSidebarPlacement;
    open.value = false;
}
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
                <div class="space-y-6">
                    <div class="space-y-3">
                        <div>
                            <h3 class="text-sm font-semibold text-highlighted">
                                {{ $t('components.dispatch.sidebar_layout.title') }}
                            </h3>
                            <p class="text-sm text-muted">
                                {{ $t('components.dispatch.sidebar_layout.description') }}
                            </p>
                        </div>

                        <div class="rounded-lg border border-default bg-elevated/50 p-3">
                            <div class="grid grid-cols-3 grid-rows-3 place-items-center gap-2">
                                <UButton
                                    v-for="option in centrumSidebarPlacementGridOptions"
                                    :key="option.placement"
                                    :color="sidebarPlacement === option.placement ? 'primary' : 'neutral'"
                                    :variant="sidebarPlacement === option.placement ? 'solid' : 'soft'"
                                    :icon="option.icon"
                                    :class="option.class"
                                    class="w-full"
                                    :label="$t(option.labelKey)"
                                    @click="setLayout(option.placement)"
                                />
                            </div>
                        </div>
                    </div>
                </div>

                <div class="mt-4 flex items-center gap-2">
                    <UTooltip class="self-start" :text="$t('common.reset')">
                        <UButton color="neutral" icon="i-mdi-clear-box-outline" variant="soft" @click="resetLayout" />
                    </UTooltip>
                </div>
            </div>
        </template>
    </UPopover>
</template>
