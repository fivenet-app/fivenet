<script lang="ts" setup>
import { defaultDispatchCenterOuterPaneLayout, useSettingsStore, type DispatchCenterOuterPane } from '~/stores/settings';

const settingsStore = useSettingsStore();
const { centrum } = storeToRefs(settingsStore);

const open = ref(false);

const sidebarPaneLayout = computed<DispatchCenterOuterPane[]>(() =>
    centrum.value.centrumSidebarPaneLayout?.length === 2
        ? centrum.value.centrumSidebarPaneLayout
        : defaultDispatchCenterOuterPaneLayout,
);

const isReversed = computed(() => sidebarPaneLayout.value[0] === 'sidebar');

function swapLayout(): void {
    centrum.value.centrumSidebarPaneLayout = sidebarPaneLayout.value[0] === 'map' ? ['sidebar', 'map'] : ['map', 'sidebar'];
    open.value = false;
}

function resetLayout(): void {
    centrum.value.centrumSidebarPaneLayout = [...defaultDispatchCenterOuterPaneLayout];
    open.value = false;
}
</script>

<template>
    <UPopover v-model:open="open" :content="{ align: 'end', sideOffset: 8 }">
        <UTooltip :text="$t('components.dispatch.sidebar_layout.trigger')">
            <UButton
                color="neutral"
                variant="ghost"
                icon="i-mdi-swap-horizontal"
                :aria-label="$t('components.dispatch.sidebar_layout.trigger')"
            />
        </UTooltip>

        <template #content>
            <div class="w-72 max-w-[calc(100vw-2rem)] p-3">
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
                            <div class="flex items-center justify-between gap-3">
                                <div class="flex flex-col">
                                    <span class="font-medium text-highlighted">
                                        {{ $t('components.dispatch.layout_popover.map') }}
                                    </span>
                                    <span class="text-xs text-muted">
                                        {{
                                            isReversed
                                                ? $t('components.dispatch.layout_popover.right')
                                                : $t('components.dispatch.layout_popover.left')
                                        }}
                                    </span>
                                </div>

                                <UButton
                                    color="neutral"
                                    variant="soft"
                                    size="sm"
                                    :label="$t('components.dispatch.sidebar_layout.swap')"
                                    @click="swapLayout"
                                />
                            </div>

                            <div class="mt-2 flex items-center justify-between gap-3">
                                <div class="flex flex-col">
                                    <span class="font-medium text-highlighted">
                                        {{ $t('components.dispatch.layout_popover.sidebar') }}
                                    </span>
                                    <span class="text-xs text-muted">
                                        {{
                                            isReversed
                                                ? $t('components.dispatch.layout_popover.left')
                                                : $t('components.dispatch.layout_popover.right')
                                        }}
                                    </span>
                                </div>
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
