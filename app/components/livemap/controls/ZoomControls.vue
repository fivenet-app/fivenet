<script lang="ts" setup>
import { useLivemapStore } from '~/stores/livemap';

const livemapStore = useLivemapStore();
const { zoom } = storeToRefs(livemapStore);
</script>

<template>
    <LControl position="topleft">
        <UButtonGroup class="inline-flex w-full flex-col" orientation="vertical">
            <UTooltip :text="$t('common.zoom_in')">
                <UButton
                    class="inset-0 border border-black/20 bg-clip-padding p-1.5"
                    icon="i-mdi-plus-thick"
                    block
                    @click="if (zoom + 1 <= 7) zoom++;"
                />
            </UTooltip>

            <UTooltip :text="$t('common.zoom')">
                <UButton class="inset-0 border border-black/20 bg-clip-padding select-none" :label="zoom.toString()" block />
            </UTooltip>

            <UTooltip :text="$t('common.zoom_out')">
                <UButton
                    class="inset-0 border border-black/20 bg-clip-padding p-1.5"
                    icon="i-mdi-minus-thick"
                    block
                    @click="if (zoom - 1 >= 1) zoom--;"
                />
            </UTooltip>
        </UButtonGroup>
    </LControl>
</template>
