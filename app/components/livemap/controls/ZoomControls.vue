<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import { useLivemapStore } from '~/stores/livemap';

const livemapStore = useLivemapStore();
const { zoom } = storeToRefs(livemapStore);
</script>

<template>
    <LControl position="topleft">
        <UButtonGroup orientation="vertical" class="inline-flex w-full flex-col">
            <UTooltip :text="$t('common.zoom_in')" :popper="{ placement: 'right' }">
                <UButton
                    icon="i-mdi-plus-thick"
                    block
                    class="inset-0 border border-black/20 bg-clip-padding p-1.5 hover:bg-[#f4f4f4]"
                    @click="if (zoom + 1 <= 7) zoom++;"
                />
            </UTooltip>

            <UTooltip :text="$t('common.zoom')" :popper="{ placement: 'right' }">
                <UButton :label="zoom.toString()" block class="inset-0 select-none border border-black/20 bg-clip-padding" />
            </UTooltip>

            <UTooltip :text="$t('common.zoom_out')" :popper="{ placement: 'right' }">
                <UButton
                    icon="i-mdi-minus-thick"
                    block
                    class="inset-0 border border-black/20 bg-clip-padding p-1.5 hover:bg-[#f4f4f4]"
                    @click="if (zoom - 1 >= 1) zoom--;"
                />
            </UTooltip>
        </UButtonGroup>
    </LControl>
</template>
