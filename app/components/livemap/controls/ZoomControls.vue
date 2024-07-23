<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import { useLivemapStore } from '~/store/livemap';

const livemapStore = useLivemapStore();
const { zoom } = storeToRefs(livemapStore);
</script>

<template>
    <LControl position="topleft">
        <div class="inline-flex w-full flex-col">
            <UButton
                class="inset-0 inline-flex items-center justify-center rounded-md border border-black/20 bg-clip-padding p-1.5 hover:bg-[#f4f4f4]"
                @click="if (zoom + 1 <= 7) zoom++;"
            >
                <UIcon name="i-mdi-plus-thick" class="size-5" />
            </UButton>

            <UButton
                class="inset-0 inline-flex select-none items-center justify-center rounded-md border border-black/20 bg-clip-padding"
            >
                {{ zoom }}
            </UButton>

            <UButton
                class="inset-0 inline-flex items-center justify-center rounded-md border border-black/20 bg-clip-padding p-1.5 hover:bg-[#f4f4f4]"
                @click="if (zoom - 1 >= 1) zoom--;"
            >
                <UIcon name="i-mdi-minus-thick" class="size-5" />
            </UButton>
        </div>
    </LControl>
</template>
